package logger

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/trace"
)

// Options 配置参数
//   - ServiceName 服务名 用于日志字段和目录
//   - Mode "file" "console" "volume"
//   - BasePath 默认 ./logs 会拼成 ./logs/{service}
//   - Level "debug" "info" "error" "severe"
//   - KeepDays 日志保留天数 仅 file volume 模式生效
//   - Encoding "json" "plain"
type Options struct {
	ServiceName string

	Mode string

	BasePath string

	Level string

	KeepDays int

	Encoding string
}

var (
	initOnce    sync.Once
	serviceName string

	// 过滤无关栈帧 标准库 测试 框架 logger 自身
	filterPrefixes = []string{
		"runtime.",
		"testing.",
		"github.com/zeromicro/go-zero/",
	}

	// 最多采集多少层调用链
	maxCallDepth = 48

	// loggerDir logger 包所在目录 用于过滤 logger 自身栈帧 避免硬编码路径
	loggerDir string
)

// init 记录 logger 包目录
func init() {
	_, file, _, ok := runtime.Caller(0)
	if ok && file != "" {
		loggerDir = filepath.ToSlash(filepath.Dir(file))
	}
}

// InitLogger 服务启动时调用一次
func InitLogger(svcName string) {
	Init(Options{
		ServiceName: svcName,
		Mode:        getenvDefault("LOG_MODE", "file"),
		BasePath:    getenvDefault("LOG_PATH", "./logs"),
		Level:       getenvDefault("LOG_LEVEL", "info"),
		KeepDays:    getenvKeepDaysDefault("LOG_KEEP_DAYS", getenvKeepDaysDefault("LOG_MAX_AGE_DAYS", 7)),
		Encoding:    getenvDefault("LOG_ENCODING", "json"),
	})
}

// Init 可配置 只初始化一次
func Init(opt Options) {
	initOnce.Do(func() {
		if strings.TrimSpace(opt.ServiceName) == "" {
			opt.ServiceName = "unknown-service"
		}
		serviceName = opt.ServiceName

		if strings.TrimSpace(opt.BasePath) == "" {
			opt.BasePath = "./logs"
		}
		if strings.TrimSpace(opt.Mode) == "" {
			opt.Mode = "file"
		}
		if strings.TrimSpace(opt.Level) == "" {
			opt.Level = "info"
		}
		if strings.TrimSpace(opt.Encoding) == "" {
			opt.Encoding = "json"
		}
		if opt.KeepDays < 0 {
			opt.KeepDays = 0
		}

		opt.Level = normalizeLevel(opt.Level)
		opt.Encoding = normalizeEncoding(opt.Encoding)
		opt.Mode = normalizeMode(opt.Mode)

		// 仅 file volume 模式需要准备目录和 Path
		path := ""
		if opt.Mode == "file" {
			path = filepath.Join(opt.BasePath, opt.ServiceName)
			if err := os.MkdirAll(path, 0o755); err != nil {
				fmt.Fprintf(os.Stderr, "[logger] mkdir logs dir failed path=%s err=%v\n", path, err)
				opt.Mode = "console"
				path = ""
			}
		} else if opt.Mode == "volume" {
			// volume 模式 Path 传 BasePath 由 go-zero 内部拼 ServiceName Hostname
			path = opt.BasePath
			if err := os.MkdirAll(path, 0o755); err != nil {
				fmt.Fprintf(os.Stderr, "[logger] mkdir logs dir failed path=%s err=%v\n", path, err)
				opt.Mode = "console"
				path = ""
			}
		}

		conf := logx.LogConf{
			ServiceName: opt.ServiceName,
			Mode:        opt.Mode,
			Encoding:    opt.Encoding,
			Level:       opt.Level,
			KeepDays:    opt.KeepDays,
		}
		if opt.Mode == "file" || opt.Mode == "volume" {
			conf.Path = path
		}

		logx.MustSetup(conf)
	})
}

// LogBizErr 逻辑层专用 打印业务错误日志
//   - ctx 请求上下文 用于提取 trace_id 把同一次请求的日志串起来排查
//   - code 业务错误码 int32 用于分类 统计 告警 例如某个 code 激增
//   - err 原始错误对象 根因 logger 只打印 err.Error() 不做 unwrap 分类
//   - extra 附加字段 可选 用于补充业务关键信息 如 repo_id, commit_id, user_id, bucket, file_key 等
func LogBizErr(ctx context.Context, code int32, err error, extra ...logx.LogField) {
	safeLog(func() {
		ctx = ensureCtx(ctx)
		err = ensureErr(err)

		fileLine, callPath, callChain := captureCallsite(2)
		traceID := traceIDFrom(ctx)

		fields := []logx.LogField{
			logx.Field("service", serviceName),
			logx.Field("trace_id", traceID),
			logx.Field("error_code", code),
			logx.Field("file_line", fileLine),
			logx.Field("call_path", callPath),
			logx.Field("call_chain", callChain),
			logx.Field("error_reason", err.Error()),
		}
		fields = append(fields, extra...)

		logx.WithContext(ctx).Errorw("business_error", fields...)
	})
}

// BizErr 打印日志并直接返回 code return logger.BizErr(ctx, code, err)
//   - ctx 请求上下文 用于提取 trace_id 把同一次请求的日志串起来排查
//   - code 业务错误码 int32 用于分类 统计 告警 例如某个 code 激增
//   - err 原始错误对象 根因 logger 只打印 err.Error() 不做 unwrap 分类
//   - extra 附加字段 可选 用于补充业务关键信息 如 repo_id, commit_id, user_id, bucket, file_key 等
func BizErr(ctx context.Context, code int32, err error, extra ...logx.LogField) int32 {
	LogBizErr(ctx, code, err, extra...)
	return code
}

// LogInfo 打印业务关键成功节点 初始化成功
//   - ctx 请求上下文 用于提取 trace_id 把同一次请求的日志串起来排查
//   - msg 业务信息 用于描述成功事件 如 publish success init db ok 等
//   - extra 附加字段 可选 用于补充业务关键信息 如 repo_id, commit_id, user_id, cost_ms 等
func LogInfo(ctx context.Context, msg string, extra ...logx.LogField) {
	safeLog(func() {
		ctx = ensureCtx(ctx)
		traceID := traceIDFrom(ctx)

		// info 默认不采集调用链 避免高频日志导致 runtime.Callers 开销过大
		fields := []logx.LogField{
			logx.Field("service", serviceName),
			logx.Field("trace_id", traceID),
			logx.Field("msg", msg),
		}
		fields = append(fields, extra...)

		logx.WithContext(ctx).Infow("business_info", fields...)
	})
}

// LogFatal 初始化阶段专用 打印致命错误并退出 业务运行期禁止用
//   - ctx 请求上下文 用于提取 trace_id 把同一次启动流程或请求相关日志串起来排查
//   - err 致命错误对象 根因 logger 只打印 err.Error() 不做 unwrap 分类
//   - extra 附加字段 可选 用于补充关键上下文 如 component, config_key, listen_addr, db_dsn_masked 等
func LogFatal(ctx context.Context, err error, extra ...logx.LogField) {
	ctx = ensureCtx(ctx)
	err = ensureErr(err)

	fileLine, callPath, callChain := captureCallsite(2)
	traceID := traceIDFrom(ctx)

	fields := []logx.LogField{
		logx.Field("service", serviceName),
		logx.Field("trace_id", traceID),
		logx.Field("file_line", fileLine),
		logx.Field("call_path", callPath),
		logx.Field("call_chain", callChain),
		logx.Field("error_reason", err.Error()),
	}
	fields = append(fields, extra...)

	// fatal 语义是必须退出 即使日志系统异常也要退出
	safeLog(func() {
		logx.WithContext(ctx).Errorw("fatal_error", fields...)
	})
	// flush 日志避免丢失
	safeLog(func() {
		_ = logx.Close()
	})
	os.Exit(1)
}

// internals

// ensureCtx 确保 ctx 非空
func ensureCtx(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

// ensureErr 确保 err 非空
func ensureErr(err error) error {
	if err == nil {
		return errors.New("nil error")
	}
	return err
}

// traceIDFrom 提取 trace id
func traceIDFrom(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	// go-zero trace
	if tid := trace.TraceIDFromContext(ctx); tid != "" {
		return tid
	}

	// 兼容一些自定义注入
	if v := ctx.Value("trace_id"); v != nil {
		return fmt.Sprint(v)
	}
	if v := ctx.Value("request_id"); v != nil {
		return fmt.Sprint(v)
	}
	return ""
}

// captureCallsite
//	- file_line xxx.go:123
// 	- call_path xxx.go:Struct.Method
// 	- call_chain Entry -> ... -> Struct.Method

func captureCallsite(skip int) (fileLine string, callPath string, callChain string) {
	// runtime.Callers 的 skip 语义 跳过 runtime.Callers 自身 + captureCallsite
	// 所以这里额外 +2
	pcs := make([]uintptr, maxCallDepth)
	// TODO call_chain 依赖 runtime.Callers 开销较大,如果后续业务错误日志变成高频 需要加开关控制或降级为 runtime.Caller 单层定位
	n := runtime.Callers(skip+2, pcs)
	if n <= 0 {
		return "unknown:0", "unknown:unknown", "unknown"
	}

	frames := runtime.CallersFrames(pcs[:n])

	type fr struct {
		file string
		line int
		fn   string
	}
	var (
		chain []string
		first *fr // 离业务最近的第一帧
		last  *fr // 离业务最远的入口帧
	)

	for {
		f, more := frames.Next()
		if !isFiltered(f.Function) && !isLoggerFrame(f.File) {
			sf := fr{file: f.File, line: f.Line, fn: f.Function}

			if first == nil {
				tmp := sf
				first = &tmp
			}
			tmp2 := sf
			last = &tmp2

			chain = append(chain, shortFuncName(f.Function))
		}
		if !more {
			break
		}
	}

	if first == nil {
		return "unknown:0", "unknown:unknown", "unknown"
	}

	// call_path file_line 用错误点 最近帧
	baseFile := filepath.Base(first.file)
	fileLine = fmt.Sprintf("%s:%d", baseFile, first.line)
	callPath = fmt.Sprintf("%s:%s", baseFile, shortFuncName(first.fn))

	// call_chain 展示入口 -> ... -> 错误点
	// runtime.CallersFrames 迭代是从近到远 因此需要反转
	reverse(chain)
	if last != nil && len(chain) == 0 {
		chain = []string{shortFuncName(last.fn)}
	}
	callChain = strings.Join(chain, " -> ")
	if callChain == "" {
		callChain = "unknown"
	}
	return
}

// isFiltered 过滤无关函数
func isFiltered(fn string) bool {
	fn = strings.ToLower(fn)
	for _, p := range filterPrefixes {
		if strings.HasPrefix(fn, p) {
			return true
		}
	}
	return false
}

// isLoggerFrame 判断是否为 logger 自身栈帧
func isLoggerFrame(file string) bool {
	// 动态过滤 logger 包自身目录 避免硬编码路径
	f := strings.ToLower(filepath.ToSlash(file))
	if loggerDir == "" {
		return false
	}
	base := strings.ToLower(loggerDir)
	return strings.HasPrefix(f, base+"/")
}

// shortFuncName 把完整函数名压缩成更可读的 Type.Method 或 pkg.Func
func shortFuncName(full string) string {
	if full == "" {
		return "unknown"
	}

	// 去掉路径
	if idx := strings.LastIndex(full, "/"); idx >= 0 {
		full = full[idx+1:]
	}

	parts := strings.Split(full, ".")
	if len(parts) <= 2 {
		return full
	}
	// 保留最后两段 例如 (*UserService).Login
	return strings.Join(parts[len(parts)-2:], ".")
}

// reverse 反转切片
func reverse(a []string) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

// safeLog 防御日志异常
func safeLog(fn func()) {
	defer func() {
		// 日志系统绝不能影响业务
		_ = recover()
	}()
	fn()
}

// getenvDefault 读取字符串环境变量
func getenvDefault(k, def string) string {
	v := strings.TrimSpace(os.Getenv(k))
	if v == "" {
		return def
	}
	return v
}

// getenvIntDefault 读取正整数环境变量
func getenvIntDefault(k string, def int) int {
	v := strings.TrimSpace(os.Getenv(k))
	if v == "" {
		return def
	}
	var x int
	_, err := fmt.Sscanf(v, "%d", &x)
	if err != nil || x <= 0 {
		return def
	}
	return x
}

// getenvKeepDaysDefault 读取日志保留天数
func getenvKeepDaysDefault(k string, def int) int {
	if def < 0 {
		def = 0
	}
	v := strings.TrimSpace(os.Getenv(k))
	if v == "" {
		return def
	}
	var x int
	_, err := fmt.Sscanf(v, "%d", &x)
	if err != nil || x < 0 {
		return def
	}
	return x
}

// getenvBoolDefault 读取布尔环境变量
func getenvBoolDefault(k string, def bool) bool {
	v := strings.TrimSpace(os.Getenv(k))
	if v == "" {
		return def
	}
	switch strings.ToLower(v) {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		return def
	}
}

// normalizeMode 规范化日志模式
func normalizeMode(m string) string {
	switch strings.ToLower(strings.TrimSpace(m)) {
	case "file", "console", "volume":
		return strings.ToLower(strings.TrimSpace(m))
	default:
		return "console"
	}
}

// normalizeLevel 规范化日志等级
func normalizeLevel(l string) string {
	switch strings.ToLower(strings.TrimSpace(l)) {
	case "debug", "info", "error", "severe":
		return strings.ToLower(strings.TrimSpace(l))
	default:
		return "info"
	}
}

// normalizeEncoding 规范化日志编码
func normalizeEncoding(e string) string {
	switch strings.ToLower(strings.TrimSpace(e)) {
	case "json", "plain":
		return strings.ToLower(strings.TrimSpace(e))
	default:
		return "json"
	}
}
