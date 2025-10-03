// cmd/lorem/main.go
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type openAIResponse struct {
	// Primary convenience field returned by the Responses API
	OutputText string `json:"output_text"`
	// Fallback detailed structure
	Output []struct {
		Content []struct {
			Type string `json:"type"` // e.g., "output_text"
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
}

func main() {
	model := flag.String("model", "gpt-5-mini", "OpenAI model (e.g. gpt-5-mini)")
	effort := flag.String("effort", "minimal", "Reasoning effort: minimal|medium|high")
	verbosity := flag.String("verbosity", "low", "Verbosity: low|medium|high")
	flag.Parse()

	// topic/prompt from args or stdin
	prompt := strings.TrimSpace(strings.Join(flag.Args(), " "))
	if prompt == "" {
		if b, _ := io.ReadAll(os.Stdin); len(b) > 0 {
			prompt = strings.TrimSpace(string(b))
		}
	}
	if prompt == "" {
		fmt.Fprintln(os.Stderr, "usage: lorem [flags] <prompt>")
		os.Exit(2)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "error: OPENAI_API_KEY not set")
		os.Exit(1)
	}

	// Build the payload. With the Responses API you can pass a simple string to `input`.
	payload := map[string]any{
		"model": *model,
		"input": prompt,
		"text": map[string]any{
			"verbosity": *verbosity,
		},
		"reasoning": map[string]any{
			"effort": *effort,
		},
		// optional: you can add temperature etc. if you ever want
		// "temperature": 0.6,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "marshal error: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/responses", bytes.NewReader(body))
	if err != nil {
		fmt.Fprintf(os.Stderr, "request error: %v\n", err)
		os.Exit(1)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "http error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "API error: %s\n%s\n", resp.Status, b)
		os.Exit(1)
	}

	var out openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		fmt.Fprintf(os.Stderr, "decode error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	// Prefer the convenience field
	if strings.TrimSpace(out.OutputText) != "" {
		fmt.Println(out.OutputText)
		fmt.Println()
		return
	}

	// Fallback: scan the structured output for text chunks
	for _, item := range out.Output {
		for _, c := range item.Content {
			if c.Type == "output_text" && strings.TrimSpace(c.Text) != "" {
				fmt.Println(c.Text)
				fmt.Println()
				return
			}
		}
	}

	fmt.Fprintln(os.Stderr, "no text in response")
	os.Exit(1)
}


