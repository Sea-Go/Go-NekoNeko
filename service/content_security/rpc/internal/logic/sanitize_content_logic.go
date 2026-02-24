package logic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"sea-try-go/service/content_security/rpc/internal/svc"
	"sea-try-go/service/content_security/rpc/pb/sea-try-go/service/content-security/rpc/pb"

	"github.com/microcosm-cc/bluemonday"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/text/unicode/norm"
)

type SanitizeContentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSanitizeContentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SanitizeContentLogic {
	return &SanitizeContentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// normalizeUnicode 将文本转换为 NFC 标准化形式
func normalizeUnicode(text string) string {
	return norm.NFC.String(text)
}

// normalizeWhitespace 统一空白字符处理
func normalizeWhitespace(text string) string {
	// 统一换行符为 \n
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	// 压缩连续空格为单个空格
	var result strings.Builder
	result.Grow(len(text))
	inSpace := false

	for _, char := range text {
		if char == ' ' || char == '\t' || char == '\n' {
			if !inSpace {
				result.WriteRune(' ')
				inSpace = true
			}
		} else {
			result.WriteRune(char)
			inSpace = false
		}
	}

	// 修剪首尾空白
	return strings.TrimSpace(result.String())
}

// sanitizeHTML 使用 bluemonday UGCPolicy 净化 HTML
func sanitizeHTML(text string, allowedTags []string) string {
	p := bluemonday.UGCPolicy()

	// 如果配置了允许的标签，则只保留这些标签
	if len(allowedTags) > 0 {
		// 首先清除所有标签
		p = bluemonday.NewPolicy()
		// 然后添加允许的标签
		for _, tag := range allowedTags {
			switch tag {
			case "b", "strong":
				p.AllowAttrs("class").OnElements("b", "strong")
			case "i", "em":
				p.AllowAttrs("class").OnElements("i", "em")
			case "u":
				p.AllowAttrs("class").OnElements("u")
			case "s":
				p.AllowAttrs("class").OnElements("s")
			case "sub":
				p.AllowAttrs("class").OnElements("sub")
			case "sup":
				p.AllowAttrs("class").OnElements("sup")
			case "blockquote":
				p.AllowAttrs("class").OnElements("blockquote")
			case "code":
				p.AllowAttrs("class").OnElements("code")
			case "pre":
				p.AllowAttrs("class").OnElements("pre")
			case "ul", "ol", "li":
				p.AllowAttrs("class").OnElements("ul", "ol", "li")
			case "dl", "dt", "dd":
				p.AllowAttrs("class").OnElements("dl", "dt", "dd")
			case "a":
				p.AllowAttrs("href", "target", "rel", "class").OnElements("a")
			case "img":
				p.AllowAttrs("src", "alt", "title", "width", "height", "class").OnElements("img")
			case "p", "br":
				p.AllowAttrs("class").OnElements("p", "br")
			}
		}
	}

	return p.Sanitize(text)
}

// detectAd 调用外部 AI 模型服务进行广告检测
func (l *SanitizeContentLogic) detectAd(text string) (bool, float32, error) {
	cfg := l.svcCtx.Config.AdDetection

	// 构建请求体
	requestBody := map[string]interface{}{
		"model": "qwen3-max",
		"input": map[string]interface{}{
			"messages": []map[string]string{
				{
					"role":    "system",
					"content": "你是一个广告检测专家。请分析以下文本是否包含广告内容。返回 JSON 格式：{\"is_ad\": boolean, \"confidence\": float}，其中 confidence 是 0.0 到 1.0 之间的置信度。",
				},
				{
					"role":    "user",
					"content": text,
				},
			},
		},
		"parameters": map[string]interface{}{
			"temperature": 0.1,
			"top_p":       0.8,
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		l.Errorf("Failed to marshal ad detection request: %v", err)
		return false, 0.0, err
	}

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(l.ctx, "POST", cfg.ApiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		l.Errorf("Failed to create ad detection request: %v", err)
		return false, 0.0, err
	}

	// 设置请求头
	req.Header.Set("Authorization", "Bearer "+cfg.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	// 设置超时
	client := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		l.Errorf("Failed to send ad detection request: %v", err)
		return false, 0.0, err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		l.Errorf("Failed to read ad detection response: %v", err)
		return false, 0.0, err
	}

	// 解析响应
	var result struct {
		Output struct {
			Choices []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			} `json:"choices"`
		} `json:"output"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		l.Errorf("Failed to unmarshal ad detection response: %v, response body: %s", err, string(body))
		return false, 0.0, err
	}

	if len(result.Output.Choices) == 0 {
		l.Errorf("No choices in ad detection response: %s", string(body))
		return false, 0.0, fmt.Errorf("no choices in response")
	}

	// 解析 AI 返回的 JSON
	var adResult struct {
		IsAd       bool    `json:"is_ad"`
		Confidence float64 `json:"confidence"`
	}

	content := result.Output.Choices[0].Message.Content
	if err := json.Unmarshal([]byte(content), &adResult); err != nil {
		l.Errorf("Failed to parse AI response content as JSON: %v, content: %s", err, content)
		// 如果解析失败，尝试提取布尔值和数字
		isAd := strings.Contains(strings.ToLower(content), "true") || strings.Contains(strings.ToLower(content), "是")
		confidence := 0.5
		if isAd {
			confidence = 0.8
		}
		return isAd, float32(confidence), nil
	}

	return adResult.IsAd, float32(adResult.Confidence), nil
}
func (l *SanitizeContentLogic) SanitizeContent(in *pb.SanitizeContentRequest) (*pb.SanitizeContentResponse, error) {
	if in == nil {
		return &pb.SanitizeContentResponse{
			Success:      false,
			ErrorMessage: "request is nil",
		}, nil
	}

	originalText := in.GetText()
	options := in.GetOptions()

	// 如果没有提供文本，直接返回成功
	if originalText == "" {
		return &pb.SanitizeContentResponse{
			SanitizedText: "",
			Success:       true,
		}, nil
	}

	processedText := originalText

	// 1. Unicode 标准化
	if options.GetEnableUnicodeNormalization() {
		processedText = normalizeUnicode(processedText)
	}

	// 2. HTML 净化
	if options.GetEnableHtmlSanitization() {
		processedText = sanitizeHTML(processedText, l.svcCtx.Config.HtmlSanitization.AllowedTags)
	}

	// 3. 空白字符处理
	if options.GetEnableWhitespaceNormalization() {
		processedText = normalizeWhitespace(processedText)
	}

	// 4. 广告检测
	var isAd bool = false
	var adConfidence float32 = 0.0

	if options.GetEnableAdDetection() && l.svcCtx.Config.AdDetection.ApiKey != "" {
		var err error
		isAd, adConfidence, err = l.detectAd(processedText)
		if err != nil {
			l.Errorf("Ad detection failed: %v", err)
			// 广告检测失败时继续处理，但记录日志
			isAd = false
			adConfidence = 0.0
		}

		// 如果置信度超过阈值，标记为广告
		if adConfidence >= float32(l.svcCtx.Config.AdDetection.Threshold) {
			isAd = true
		}
	}

	return &pb.SanitizeContentResponse{SanitizedText: processedText,
		IsAd:         isAd,
		AdConfidence: adConfidence,
		Success:      true}, nil
}
