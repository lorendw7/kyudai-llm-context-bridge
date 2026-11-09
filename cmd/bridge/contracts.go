package main

import (
	"context"
	"errors"
	"strings"
)

type CreateTurnRequest struct {
	SchemaVersion string `json:"schema_version"`
	RequestID     string `json:"request_id"`
	ModelID       string `json:"model_id"`
	Content       string `json:"content"`
}

func (r CreateTurnRequest) Validate() error {
	if r.SchemaVersion != "1" {
		return errors.New("unsupported schema version")
	}

	if strings.TrimSpace(r.RequestID) == "" {
		return errors.New("request ID is required")
	}

	if strings.TrimSpace(r.ModelID) == "" {
		return errors.New("model ID is required")
	}

	if strings.TrimSpace(r.Content) == "" {
		return errors.New("content is required")
	}

	return nil
}

type GenerateRequest struct {
	ModelID         string
	Messages        []Message
	MaxOutputTokens int
}

type GenerateResult struct {
	Content      string
	InputTokens  int
	OutputTokens int
}

type ModelClient interface {
	Generate(
		ctx context.Context,
		request GenerateRequest,
	) (GenerateResult, error)
}
