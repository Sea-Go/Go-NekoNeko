package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
	z_trace "github.com/zeromicro/go-zero/core/trace"
)

// Logger 日志工具结构体（替代全局变量，支持多实例/并发安全）
type Logger struct {
	serviceName string    // 服务名（用于日志标识）
	once        sync.Once // 确保初始化只执行一次
}

// 全局实例（方便业务层快速调用，也可按需创建独立实例）
var globalLogger = &Logger{}

var (
	// 过滤无关栈帧，仅保留业务代码调用链路
	filterPrefixes = []string{
		"runtime.",
		"testing.",
		"github.com/zeromicro/go-zero/",
	}
)

// Init 初始化日志工具（服务启动时调用）
// svcName：服务名（如 user-api/order-rpc）
// 注意：日志的基础配置（Mode/Path/Level等）由go-zero配置文件加载，此处仅初始化工具包
func Init(svcName string) {
	globalLogger.once.Do(func() {
		if svcName == "" {
			panic("service name cannot be empty")
		}
		globalLogger.serviceName = svcName
	})
}

// LogBusinessErr 打印业务错误日志（核心函数）
// ctx：上下文（用于提取TraceID）
// code：业务错误码
// err：原始错误对象
func LogBusinessErr(ctx context.Context, code int, err error) {
	globalLogger.logBusinessErr(ctx, code, err)
}

// LogInfo 打印业务信息日志
// ctx：上下文
// msg：信息内容
func LogInfo(ctx context.Context, msg string) {
	globalLogger.logInfo(ctx, msg)
}

// LogFatal 打印致命错误日志并终止程序（仅用于初始化失败）
// ctx：上下文
// err：致命错误对象
func LogFatal(ctx context.Context, err error) {
	globalLogger.logFatal(ctx, err)
}

// logBusinessErr 私有实现：业务错误日志
func (l *Logger) logBusinessErr(ctx context.Context, code int, err error) {
	// 1. 校验服务名是否初始化
	if l.serviceName == "" {
		panic("logger not initialized, call logger.Init() first")
	}
	errorReason := "no return,missing task error"
	if err != nil {
		errorReason = err.Error()
	}

	// 2. 采集单步调用信息（文件行号、函数名）
	pc, file, line, ok := runtime.Caller(2) //获取业务调用方法
	fileLine := ""
	callPath := ""
	if ok {
		fileName := filepath.Base(file)
		fileLine = fmt.Sprintf("%s:%d", fileName, line)
		funcName := getFuncName(pc)
		callPath = fmt.Sprintf("%s:%s", fileName, funcName)
	}

	// 3. 采集完整调用链路
	callChain := getCallChain(3) // skip=3：跳过当前方法+LogBusinessErr+业务直接调用方法

	// 4. 提取TraceID
	traceID := z_trace.TraceIDFromContext(ctx)
	if traceID == "" {
		traceID = "unknown"
	}

	// 5. 输出标准化日志
	logx.WithContext(ctx).Errorw(
		"business_error",
		logx.Field("service", l.serviceName),
		logx.Field("trace_id", traceID),
		logx.Field("error_code", code),
		logx.Field("file_line", fileLine),
		logx.Field("call_path", callPath),
		logx.Field("call_chain", callChain),
		logx.Field("error_reason", errorReason),
	)
}

// logInfo 私有实现：业务信息日志
func (l *Logger) logInfo(ctx context.Context, msg string) {
	// 校验初始化
	if l.serviceName == "" {
		panic("logger not initialized, call logger.Init() first")
	}

	// 采集调用位置
	pc, file, _, ok := runtime.Caller(2) // skip=2：跳过当前方法+LogInfo，获取业务调用方
	position := "unknown"
	if ok {
		fileName := filepath.Base(file)
		funcName := getFuncName(pc)
		position = fmt.Sprintf("%s:%s", fileName, funcName)
	}

	// 提取TraceID
	traceID := z_trace.TraceIDFromContext(ctx)
	if traceID == "" {
		traceID = "unknown"
	}

	// 输出日志
	logx.WithContext(ctx).Infow(
		"business_info",
		logx.Field("service", l.serviceName),
		logx.Field("trace_id", traceID),
		logx.Field("position", position),
		logx.Field("msg", msg),
	)
}

// logFatal 私有实现：致命错误日志
func (l *Logger) logFatal(ctx context.Context, err error) {
	// 校验初始化
	if l.serviceName == "" {
		panic("logger not initialized, call logger.Init() first")
	}
	errorReason := "no return,missing task error "
	if err != nil {
		errorReason = err.Error()
	}

	// 采集调用信息
	pc, file, line, ok := runtime.Caller(2) // skip=2：跳过当前方法+LogFatal，获取业务调用方
	fileLine := "unknown/file.go:0"
	callPath := "unknown"
	if ok {
		fileName := filepath.Base(file)
		fileLine = fmt.Sprintf("%s:%d", fileName, line)
		funcName := getFuncName(pc)
		callPath = fmt.Sprintf("%s:%s", fileName, funcName)
	}

	// 采集调用链路
	callChain := getCallChain(3)

	// 提取TraceID
	traceID := z_trace.TraceIDFromContext(ctx)
	if traceID == "" {
		traceID = "unknown"
	}

	// 输出致命错误日志（内部已调用os.Exit(1)）
	logx.WithContext(ctx).Errorw(
		"fatal_error",
		logx.Field("service", l.serviceName),
		logx.Field("trace_id", traceID),
		logx.Field("file_line", fileLine),
		logx.Field("call_path", callPath),
		logx.Field("call_chain", callChain),
		logx.Field("error_reason", errorReason),
	)

	// 兜底退出（极端场景）
	os.Exit(1)
}

func getFuncName(pc uintptr) string {
	if pc == 0 {
		return "unknown"
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}
	fullName := fn.Name()
	parts := strings.Split(fullName, ".")
	if len(parts) >= 2 {
		// 处理匿名函数/闭包场景
		if strings.Contains(parts[len(parts)-2], "func") {
			return parts[len(parts)-1]
		}
		return strings.Join(parts[len(parts)-2:], ".")
	}
	return filepath.Base(fullName)
}

// getCallChain 采集并格式化完整调用链路（过滤无关栈帧）
func getCallChain(skip int) string {
	const maxDepth = 32 // 限制最大采集深度，避免栈溢出
	pcs := make([]uintptr, maxDepth)
	n := runtime.Callers(skip, pcs)
	if n == 0 {
		return "unknown"
	}

	var chain []string
	for _, pc := range pcs[:n] {
		funcName := getFuncName(pc)
		if isFiltered(funcName) {
			continue
		}
		chain = append(chain, funcName)
	}

	// 反转链路，呈现「入口函数→下游函数」的调用顺序
	for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
		chain[i], chain[j] = chain[j], chain[i]
	}

	if len(chain) == 0 {
		return "unknown"
	}
	return strings.Join(chain, " -> ")
}

// isFiltered 判断函数是否为无关栈帧，需要过滤
func isFiltered(funcName string) bool {
	for _, prefix := range filterPrefixes {
		if strings.HasPrefix(funcName, prefix) {
			return true
		}
	}
	return false
}
