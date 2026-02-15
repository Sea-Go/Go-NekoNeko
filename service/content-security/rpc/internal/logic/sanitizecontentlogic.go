package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"sea-try-go/service/common/logger"
	"sea-try-go/service/content-security/rpc/internal/config"
	"sea-try-go/service/content-security/rpc/internal/svc"
	"sea-try-go/service/content-security/rpc/pb"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/text/unicode/norm"
)

type SanitizeContentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSanitizeContentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SanitizeContentLogic {
	return &SanitizeContentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SanitizeContentLogic) SanitizeContent(in *pb.SanitizeContentRequest) (*pb.SanitizeContentResponse, error) {
	if in == nil || in.Text == "" {
		logger.LogInfo(l.ctx, "收到空内容请求")
		return &pb.SanitizeContentResponse{
			Success:      false,
			ErrorMessage: "输入文本为空",
		}, nil
	}

	text := in.GetText()
	options := in.GetOptions()
	if options == nil {
		// 默认启用所有选项
		options = &pb.SanitizeOptions{
			EnableHtmlSanitization:        true,
			EnableAdDetection:             true,
			EnableUnicodeNormalization:    true,
			EnableWhitespaceNormalization: true,
		}
	}

	logger.LogInfo(l.ctx, "开始内容安全清洗")

	// 1. Unicode 标准化
	if options.EnableUnicodeNormalization {
		text = norm.NFC.String(text)
	}

	// 2. 空白处理
	if options.EnableWhitespaceNormalization {
		text = normalizeWhitespace(text)
	}

	// 3. HTML 净化
	if options.EnableHtmlSanitization {
		text = sanitizeHTML(text, l.svcCtx.Config)
	}

	// 4. 广告检测
	var isAd bool
	var adConfidence float32
	if options.EnableAdDetection {
		var err error
		isAd, adConfidence, err = l.detectAd(text)
		if err != nil {
			logger.LogBusinessErr(l.ctx, 500, err, logger.WithArticleID("unknown"))
			// 降级策略：跳过广告检测，返回安全结果
			isAd = false
			adConfidence = 0
		}
	}

	logger.LogInfo(l.ctx, "内容清洗完成")
	return &pb.SanitizeContentResponse{
		SanitizedText: text,
		IsAd:          isAd,
		AdConfidence:  adConfidence,
		Success:       true,
	}, nil
}

// normalizeWhitespace 标准化空白字符
func normalizeWhitespace(text string) string {
	// 统一换行符为 \n
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	// 压缩连续空格
	spaceRegex := regexp.MustCompile(`\s+`)
	text = spaceRegex.ReplaceAllString(text, " ")

	// 去除首尾空白
	return strings.TrimSpace(text)
}

// sanitizeHTML 使用 bluemonday 净化 HTML
func sanitizeHTML(text string, config config.Config) string {
	p := bluemonday.NewPolicy()

	// 允许基本的格式化标签
	p.AllowStandardURLs()
	p.AllowAttrs("href").OnElements("a")
	p.AllowAttrs("src").OnElements("img")
	p.AllowAttrs("alt").OnElements("img")
	p.AllowAttrs("title").OnElements("abbr", "acronym", "cite", "code", "dfn", "em", "strong", "q")

	// 允许的标签
	allowedTags := []string{"b", "i", "p", "br", "strong", "em", "u", "s", "sub", "sup", "blockquote", "code", "pre", "ul", "ol", "li", "dl", "dt", "dd"}

	// 如果配置中指定了允许的标签，则使用配置
	if len(config.HtmlSanitization.AllowedTags) > 0 {
		allowedTags = config.HtmlSanitization.AllowedTags
	}

	for _, tag := range allowedTags {
		p.AllowElements(tag)
	}

	return p.Sanitize(text)
}

// detectAd 调用外部 AI 模型进行广告检测
func (l *SanitizeContentLogic) detectAd(text string) (bool, float32, error) {
	config := l.svcCtx.Config.AdDetection

	if config.ApiEndpoint == "" {
		return false, 0, fmt.Errorf("广告检测 API 端点未配置")
	}

	// 构建请求体
	requestBody := map[string]interface{}{
		"model": "qwen3-max",
		"input": map[string]interface{}{
			"messages": []map[string]string{
				{
					"role":    "system",
					"content": "你是一个广告检测器。请分析以下内容是否包含广告，并返回置信度（0-1之间）。只返回JSON格式：{\"is_ad\": boolean, \"confidence\": float}",
				},
				{
					"role":    "user",
					"content": text,
				},
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return false, 0, fmt.Errorf("构建请求体失败: %w", err)
	}

	// 创建 HTTP 客户端
	client := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
	}
	if config.Timeout == 0 {
		client.Timeout = 10 * time.Second
	}

	// 发送请求
	req, err := http.NewRequest("POST", config.ApiEndpoint, strings.NewReader(string(jsonData)))
	if err != nil {
		return false, 0, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if config.ApiKey != "" {
		req.Header.Set("Authorization", "Bearer "+config.ApiKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, 0, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, 0, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return false, 0, fmt.Errorf("API 返回错误状态码: %d, body: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var result struct {
		IsAd       bool    `json:"is_ad"`
		Confidence float64 `json:"confidence"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		// 尝试解析可能的嵌套响应
		var wrapper struct {
			Output struct {
				Choices []struct {
					Message struct {
						Content string `json:"content"`
					} `json:"message"`
				} `json:"choices"`
			} `json:"output"`
		}

		if err2 := json.Unmarshal(body, &wrapper); err2 == nil && len(wrapper.Output.Choices) > 0 {
			// 解析 AI 返回的 JSON 字符串
			content := wrapper.Output.Choices[0].Message.Content
			if err3 := json.Unmarshal([]byte(content), &result); err3 == nil {
				// 成功解析
			} else {
				return false, 0, fmt.Errorf("无法解析 AI 响应: %s", content)
			}
		} else {
			return false, 0, fmt.Errorf("无法解析 API 响应: %w", err)
		}
	}

	// 判断是否为广告（基于阈值）
	threshold := config.Threshold
	if threshold == 0 {
		threshold = 0.7 // 默认阈值
	}

	isAd := result.Confidence > threshold
	confidence := float32(result.Confidence)

	return isAd, confidence, nil
}
