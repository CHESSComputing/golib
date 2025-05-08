package ollama

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Config struct {
	Host     string
	Port     int
	Model    string
	Timeout  time.Duration
	MaxRetry int
	Client   *http.Client
}

type Client struct {
	config Config
	url    string
}

type Request struct {
	Model   string `json:"model"`
	Prompt  string `json:"prompt"`
	Stream  bool   `json:"stream,omitempty"`
	NumPred int    `json:"num_predict,omitempty"`
}

type Response struct {
	Response string `json:"response"`
	Done     bool   `json:"done,omitempty"`
	Error    string `json:"error,omitempty"`
}

func NewClient(config Config) *Client {
	if config.Client == nil {
		config.Client = &http.Client{}
	}
	url := fmt.Sprintf("http://%s:%d/api/generate", config.Host, config.Port)
	return &Client{config: config, url: url}
}

func (c *Client) SendPrompt(ctx context.Context, prompt string) (string, error) {
	return c.sendWithRetry(ctx, prompt, false, nil)
}

func (c *Client) StreamPrompt(ctx context.Context, prompt string, callback func(string)) error {
	_, err := c.sendWithRetry(ctx, prompt, true, callback)
	return err
}

func (c *Client) sendWithRetry(ctx context.Context, prompt string, stream bool, callback func(string)) (string, error) {
	var lastErr error
	for attempt := 0; attempt <= c.config.MaxRetry; attempt++ {
		resp, err := c.SendRequest(ctx, prompt, stream, callback)
		if err == nil {
			return resp, nil
		}
		lastErr = err
		time.Sleep(time.Second * time.Duration(attempt))
	}
	return "", fmt.Errorf("all retries failed: %w", lastErr)
}

func (c *Client) SendRequest(ctx context.Context, prompt string, stream bool, callback func(string)) (string, error) {
	reqBody := Request{
		Model:  c.config.Model,
		Prompt: prompt,
		Stream: stream,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequestWithContext(ctx, "POST", c.url, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.config.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama server error: %s", string(body))
	}

	if stream {
		var out []string
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return "", err
			}
			var res Response
			if err := json.Unmarshal(line, &res); err != nil {
				continue
			}
			if callback != nil {
				callback(res.Response)
			} else {
				out = append(out, res.Response)
			}
			if res.Done {
				break
			}
		}
		return strings.Join(out, ""), nil
	}

	var res Response
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}
	if res.Error != "" {
		return "", errors.New(res.Error)
	}

	return res.Response, nil
}
