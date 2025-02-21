package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/madneal/gshark/global"
	"io"
	"net/http"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatCompletionResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

func Question(command, question string) string {
	result, err := callChatCompletion(command, question)
	if err != nil {
		//global.GVA_LOG.Error("callChatCompletion failed", zap.Error(err))
		return ""
	}
	answer, err := handleResponse(result)
	if err != nil {
		//global.GVA_LOG.Error("handleResponse failed", zap.Error(err))
		return ""
	}
	return answer
}

func callChatCompletion(command, question string) ([]byte, error) {
	var result []byte
	requestData := ChatCompletionRequest{
		Model: "deepseek-r1",
		Messages: []Message{
			{
				Role:    "system",
				Content: command,
			},
			{
				Role:    "user",
				Content: question,
			},
		},
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return result, fmt.Errorf("error marshalling payload: %v", err)
	}
	url := global.GVA_CONFIG.System.AiServer
	token := global.GVA_CONFIG.System.AiToken
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return result, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return result, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("error reading response: %v", err)
	}

	fmt.Printf("Response Status: %s\nResponse Body: %s\n", resp.Status, string(body))

	return body, err
}

func handleResponse(respBody []byte) (string, error) {
	var res ChatCompletionResponse
	if err := json.Unmarshal(respBody, &res); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if len(res.Choices) == 0 {
		return "", fmt.Errorf("no choices found in response")
	}

	// 返回第一个 choice 的 message content
	return res.Choices[0].Message.Content, nil
}
