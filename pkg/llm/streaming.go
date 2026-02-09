// Package llm provides streaming support for the OpenAI API.
package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// StreamingCallback is called for each chunk of the response
type StreamingCallback func(chunk string, done bool)

// StreamingDelta represents a streaming response delta
type StreamingDelta struct {
	Content   string     `json:"content,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// StreamingChoice represents a streaming choice
type StreamingChoice struct {
	Index        int            `json:"index"`
	Delta        StreamingDelta `json:"delta"`
	FinishReason string         `json:"finish_reason"`
}

// StreamingResponse represents a streaming response chunk
type StreamingResponse struct {
	ID      string            `json:"id"`
	Object  string            `json:"object"`
	Created int64             `json:"created"`
	Model   string            `json:"model"`
	Choices []StreamingChoice `json:"choices"`
}

// ChatWithStreaming sends a chat request and streams the response
func (c *Client) ChatWithStreaming(messages []Message, tools []ToolDefinition, temperature *float64, topP *float64, maxTokens *int, callback StreamingCallback) (*ChatResponse, error) {
	reqBody := ChatRequest{
		Model:       c.model,
		Messages:    messages,
		Tools:       tools,
		Temperature: temperature,
		MaxTokens:   maxTokens,
		Stream:      true,
	}

	if len(tools) > 0 {
		reqBody.ToolChoice = "auto"
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "text/event-stream")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse SSE stream
	return c.parseSSEStream(resp.Body, callback)
}

// parseSSEStream parses Server-Sent Events and calls callback for each chunk
func (c *Client) parseSSEStream(body io.Reader, callback StreamingCallback) (*ChatResponse, error) {
	reader := bufio.NewReader(body)

	var fullContent strings.Builder
	var toolCalls []ToolCall
	var lastResponse StreamingResponse

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}

		// Check for data line
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")

		// Check for stream end
		if data == "[DONE]" {
			if callback != nil {
				callback("", true)
			}
			break
		}

		// Parse JSON
		var streamResp StreamingResponse
		if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
			log.Printf("[WARN] Skipping malformed streaming chunk: %v (data: %s)", err, data)
			continue // Skip malformed chunks
		}

		lastResponse = streamResp

		if len(streamResp.Choices) > 0 {
			choice := streamResp.Choices[0]

			// Handle content
			if choice.Delta.Content != "" {
				fullContent.WriteString(choice.Delta.Content)
				if callback != nil {
					callback(choice.Delta.Content, false)
				}
			}

			// Handle tool calls
			if len(choice.Delta.ToolCalls) > 0 {
				for _, tc := range choice.Delta.ToolCalls {
					// Merge tool call deltas
					if tc.ID != "" {
						toolCalls = append(toolCalls, tc)
					} else if len(toolCalls) > 0 {
						// Append to last tool call's arguments
						last := &toolCalls[len(toolCalls)-1]
						last.Function.Arguments += tc.Function.Arguments
					}
				}
			}
		}
	}

	// Construct final response
	finalResponse := &ChatResponse{
		ID:      lastResponse.ID,
		Object:  "chat.completion",
		Created: lastResponse.Created,
		Model:   lastResponse.Model,
		Choices: []Choice{
			{
				Index: 0,
				Message: Message{
					Role:      "assistant",
					Content:   fullContent.String(),
					ToolCalls: toolCalls,
				},
				FinishReason: "stop",
			},
		},
	}

	return finalResponse, nil
}

// SimpleQueryStreaming sends a simple query with streaming
func (c *Client) SimpleQueryStreaming(prompt string, temperature *float64, topP *float64, maxTokens *int, callback StreamingCallback) (string, error) {
	messages := []Message{
		{Role: "user", Content: prompt},
	}

	resp, err := c.ChatWithStreaming(messages, nil, temperature, topP, maxTokens, callback)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	return resp.Choices[0].Message.Content, nil
}
