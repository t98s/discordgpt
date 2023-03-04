package gpt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	MessageRoleSystem    = "system"
	MessageRoleUser      = "user"
	MessageRoleAssistant = "assistant"
)

// Message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionReq
type ChatCompletionReq struct {
	Model            string         `json:"model"`
	Messages         []Message      `json:"messages"`
	MaxTokens        int            `json:"max_tokens,omitempty"`
	Temperature      float32        `json:"temperature,omitempty"`
	TopP             float32        `json:"top_p,omitempty"`
	N                int            `json:"n,omitempty"`
	Stream           bool           `json:"stream,omitempty"`
	Stop             []string       `json:"stop,omitempty"`
	PresencePenalty  float32        `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32        `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int `json:"logit_bias,omitempty"`
	User             string         `json:"user,omitempty"`
}

// Choice
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatCompletionRes
type ChatCompletionRes struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// CreateChatCompletion
func CreateChatCompletion(ctx context.Context, chatCompletionReq ChatCompletionReq) (ChatCompletionRes, error) {
	var chatCompletionRes ChatCompletionRes

	var reqBytes []byte
	reqBytes, err := json.Marshal(chatCompletionReq)
	if err != nil {
		return chatCompletionRes, err
	}

	err = sendRequest(ctx, http.MethodPost, "/chat/completions", bytes.NewBuffer(reqBytes), &chatCompletionRes)
	if err != nil {
		return chatCompletionRes, err
	}

	return chatCompletionRes, nil
}

// sendRequest
func sendRequest(ctx context.Context, method string, endpoint string, body io.Reader, v interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s%s", "https://api.openai.com/v1", endpoint), body)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("GPT_SECRET_KEY")))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("OpenAI-Organization", "")

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= http.StatusBadRequest {
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
		return fmt.Errorf("status code: %d, response: %s", res.StatusCode, string(b))
	}

	if v != nil {
		if err = json.NewDecoder(res.Body).Decode(v); err != nil {
			return err
		}
	}

	return nil
}
