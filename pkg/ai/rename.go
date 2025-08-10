package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/krau/SaveAny-Bot/config"
)

// OpenAI API compatible request/response structures
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
	Temperature float32       `json:"temperature,omitempty"`
}

type ChatChoice struct {
	Message ChatMessage `json:"message"`
}

type ChatResponse struct {
	Choices []ChatChoice `json:"choices"`
	Error   *APIError    `json:"error,omitempty"`
}

type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

type RenameService struct {
	config *config.AIRename
	client *http.Client
}

func NewRenameService(cfg *config.AIRename) *RenameService {
	return &RenameService{
		config: cfg,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

// GenerateFileName generates a new filename based on the message content and original filename
func (s *RenameService) GenerateFileName(ctx context.Context, messageText, originalFileName string) (string, error) {
	logger := log.FromContext(ctx)
	
	if !s.config.Enable {
		return originalFileName, nil
	}

	if messageText == "" {
		return originalFileName, nil
	}

	// Get file extension
	ext := filepath.Ext(originalFileName)
	baseName := strings.TrimSuffix(originalFileName, ext)

	// Build prompt
	prompt := s.buildPrompt(messageText, baseName)

	// Call AI API
	newName, err := s.callAI(ctx, prompt)
	if err != nil {
		logger.Errorf("AI rename failed: %v", err)
		return originalFileName, err
	}

	// Clean and validate the new name
	cleanName := s.cleanFileName(newName)
	if cleanName == "" {
		logger.Warn("AI returned empty filename, using original")
		return originalFileName, nil
	}

	// Ensure the extension is preserved
	if !strings.HasSuffix(cleanName, ext) {
		cleanName += ext
	}

	logger.Infof("AI renamed file: %s -> %s", originalFileName, cleanName)
	return cleanName, nil
}

// GenerateFolderName generates a folder name for media groups
func (s *RenameService) GenerateFolderName(ctx context.Context, messageText string, defaultName string) (string, error) {
	logger := log.FromContext(ctx)
	
	if !s.config.Enable {
		return defaultName, nil
	}

	if messageText == "" {
		return defaultName, nil
	}

	// Build prompt for folder naming
	prompt := s.buildFolderPrompt(messageText, defaultName)

	// Call AI API
	newName, err := s.callAI(ctx, prompt)
	if err != nil {
		logger.Errorf("AI folder rename failed: %v", err)
		return defaultName, err
	}

	// Clean and validate the new name
	cleanName := s.cleanFolderName(newName)
	if cleanName == "" {
		logger.Warn("AI returned empty folder name, using default")
		return defaultName, nil
	}

	logger.Infof("AI renamed folder: %s -> %s", defaultName, cleanName)
	return cleanName, nil
}

func (s *RenameService) buildPrompt(messageText, originalBaseName string) string {
	defaultPrompt := `请根据以下消息内容，为文件生成一个合适的中文文件名。要求：
1. 文件名应该简洁明了，能够反映文件内容
2. 使用中文
3. 不要包含特殊字符，只使用中文、英文字母、数字、下划线和短横线
4. 长度控制在50个字符以内
5. 只返回文件名，不要包含扩展名和其他说明

消息内容：%s
原始文件名：%s

请生成新的文件名：`

	prompt := s.config.Prompt
	if prompt == "" {
		prompt = defaultPrompt
	}

	return fmt.Sprintf(prompt, messageText, originalBaseName)
}

func (s *RenameService) buildFolderPrompt(messageText, defaultName string) string {
	defaultPrompt := `请根据以下消息内容，为相册文件夹生成一个合适的中文文件夹名。要求：
1. 文件夹名应该简洁明了，能够反映相册内容
2. 使用中文
3. 不要包含特殊字符，只使用中文、英文字母、数字、下划线和短横线
4. 长度控制在30个字符以内
5. 只返回文件夹名，不要包含其他说明

消息内容：%s
默认文件夹名：%s

请生成新的文件夹名：`

	prompt := s.config.Prompt
	if prompt == "" {
		prompt = defaultPrompt
	}

	return fmt.Sprintf(prompt, messageText, defaultName)
}

func (s *RenameService) callAI(ctx context.Context, prompt string) (string, error) {
	req := ChatRequest{
		Model: s.config.Model,
		Messages: []ChatMessage{
			{Role: "user", Content: prompt},
		},
		MaxTokens:   s.config.MaxTokens,
		Temperature: s.config.Temperature,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", s.config.Endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+s.config.APIKey)

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %d %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("API error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return strings.TrimSpace(chatResp.Choices[0].Message.Content), nil
}

func (s *RenameService) cleanFileName(name string) string {
	// Remove quotes and other unwanted characters
	name = strings.Trim(name, "\"'`")
	name = strings.TrimSpace(name)
	
	// Replace invalid characters for file names
	invalidChars := regexp.MustCompile(`[<>:"/\\|?*]`)
	name = invalidChars.ReplaceAllString(name, "_")
	
	// Replace multiple spaces/underscores with single underscore
	multiSpace := regexp.MustCompile(`[\s_]+`)
	name = multiSpace.ReplaceAllString(name, "_")
	
	// Trim underscores from start and end
	name = strings.Trim(name, "_")
	
	// Limit length
	if len(name) > 50 {
		runes := []rune(name)
		name = string(runes[:50])
	}
	
	return name
}

func (s *RenameService) cleanFolderName(name string) string {
	// Similar to cleanFileName but for folders
	name = strings.Trim(name, "\"'`")
	name = strings.TrimSpace(name)
	
	// Replace invalid characters for folder names
	invalidChars := regexp.MustCompile(`[<>:"/\\|?*]`)
	name = invalidChars.ReplaceAllString(name, "_")
	
	// Replace multiple spaces/underscores with single underscore
	multiSpace := regexp.MustCompile(`[\s_]+`)
	name = multiSpace.ReplaceAllString(name, "_")
	
	// Trim underscores from start and end
	name = strings.Trim(name, "_")
	
	// Limit length for folders
	if len(name) > 30 {
		runes := []rune(name)
		name = string(runes[:30])
	}
	
	return name
}
