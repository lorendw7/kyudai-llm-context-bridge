package main

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

type recordingModelClient struct {
	result  GenerateResult
	request GenerateRequest
	called  bool
}

func (f *recordingModelClient) Generate(
	ctx context.Context,
	request GenerateRequest,
) (GenerateResult, error) {
	if err := ctx.Err(); err != nil {
		return GenerateResult{}, err
	}

	f.called = true
	f.request = request

	return f.result, nil
}

var _ ModelClient = (*recordingModelClient)(nil)

func TestCreateTurnRequestJSONContract(t *testing.T) {
	request := CreateTurnRequest{
		SchemaVersion: "1",
		RequestID:     "synthetic-turn-003",
		ModelID:       "local-model-b",
		Content:       "Correction: the deadline is November 15.",
	}

	if err := request.Validate(); err != nil {
		t.Fatalf("valid CreateTurnRequest was rejected: %v", err)
	}

	encoded, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("marshal CreateTurnRequest: %v", err)
	}

	var got map[string]string
	if err := json.Unmarshal(encoded, &got); err != nil {
		t.Fatalf("decode marshaled CreateTurnRequest: %v", err)
	}

	want := map[string]string{
		"schema_version": "1",
		"request_id":     "synthetic-turn-003",
		"model_id":       "local-model-b",
		"content":        "Correction: the deadline is November 15.",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected JSON contract: got %#v, want %#v", got, want)
	}
}

func TestCreateTurnRequestRejectsInvalidInput(t *testing.T) {
	valid := CreateTurnRequest{
		SchemaVersion: "1",
		RequestID:     "synthetic-turn-003",
		ModelID:       "local-model-b",
		Content:       "hello",
	}

	tests := []struct {
		name   string
		mutate func(*CreateTurnRequest)
	}{
		{
			name: "unsupported schema version",
			mutate: func(request *CreateTurnRequest) {
				request.SchemaVersion = "2"
			},
		},
		{
			name: "blank request ID",
			mutate: func(request *CreateTurnRequest) {
				request.RequestID = "   "
			},
		},
		{
			name: "blank model ID",
			mutate: func(request *CreateTurnRequest) {
				request.ModelID = "   "
			},
		},
		{
			name: "blank content",
			mutate: func(request *CreateTurnRequest) {
				request.Content = "   "
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := valid
			tt.mutate(&request)

			if err := request.Validate(); err == nil {
				t.Fatal("Validate returned nil for invalid request")
			}
		})
	}
}

func TestModelClientUsesProjectTypes(t *testing.T) {
	request := GenerateRequest{
		ModelID: "local-model-b",
		Messages: []Message{
			{
				ID:        "message-1",
				SessionID: "session-1",
				UserID:    "user-1",
				Role:      MessageRoleUser,
				Content:   "hello",
			},
		},
		MaxOutputTokens: 256,
	}
	want := GenerateResult{
		Content:      "synthetic answer",
		InputTokens:  12,
		OutputTokens: 3,
	}
	fake := &recordingModelClient{result: want}

	got, err := callModel(context.Background(), fake, request)
	if err != nil {
		t.Fatalf("callModel returned an unexpected error: %v", err)
	}

	if !fake.called {
		t.Fatal("ModelClient.Generate was not called")
	}
	if !reflect.DeepEqual(fake.request, request) {
		t.Fatalf("unexpected GenerateRequest: got %#v, want %#v", fake.request, request)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected GenerateResult: got %#v, want %#v", got, want)
	}
}

func TestModelClientReceivesCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	fake := &recordingModelClient{}
	_, err := callModel(ctx, fake, GenerateRequest{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("callModel error = %v, want context.Canceled", err)
	}
	if fake.called {
		t.Fatal("canceled request reached model generation")
	}
}

func callModel(
	ctx context.Context,
	client ModelClient,
	request GenerateRequest,
) (GenerateResult, error) {
	return client.Generate(ctx, request)
}
