package gogpt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CompletionRequest struct {
	Prompt    string `json:"prompt,omitempty"`
	MaxTokens int    `json:"max_tokens,omitempty"`

	Temperature float32 `json:"temperature,omitempty"`
	TopP        float32 `json:"top_p,omitempty"`

	BestOf int `json:"best_of,omitempty"`
	N      int `json:"n,omitempty"`

	LogProbs int `json:"logprobs,omitempty"`

	Echo bool     `json:"echo,omitempty"`
	Stop []string `json:"stop,omitempty"`

	PresencePenalty  float32 `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32 `json:"frequency_penalty,omitempty"`
}

type Choice struct {
	Text         string        `json:"text"`
	Index        int           `json:"index"`
	FinishReason string        `json:"finish_reason"`
	LogProbs     LogprobResult `json:"logprobs"`
}

type LogprobResult struct {
	Tokens        []string             `json:"tokens"`
	TokenLogprobs []float32            `json:"token_logprobs"`
	TopLogprobs   []map[string]float32 `json:"top_logprobs"`
	TextOffset    []int                `json:"text_offset"`
}

type CompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created uint64   `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

// CreateCompletion — API call to create a completion. This is the main endpoint of the API. Returns new text as well as, if requested, the probabilities over each alternative token at each position.
func (c *Client) CreateCompletion(ctx context.Context, engineID string, request CompletionRequest) (response CompletionResponse, err error) {
	var reqBytes []byte
	reqBytes, err = json.Marshal(request)
	if err != nil {
		return
	}

	urlSuffix := fmt.Sprintf("/engines/%s/completions", engineID)
	req, err := http.NewRequest("POST", c.fullURL(urlSuffix), bytes.NewBuffer(reqBytes))
	if err != nil {
		return
	}

	req = req.WithContext(ctx)
	err = c.sendRequest(req, &response)
	if err != nil {
		data, _ := ioutil.ReadAll(req.Body)
		fmt.Printf("[gogpt] error response data: %v\n", err, string(data))
		return
	}
	return
}
