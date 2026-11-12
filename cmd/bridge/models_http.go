package main

import (
	"encoding/json"
	"net/http"
)

type ListModelsResponse struct {
	Models []string `json:"models"`
}

func newListModelsHandler(
	allowlist ModelAllowlist,
) authenticatedHandler {
	return func(
		w http.ResponseWriter,
		r *http.Request,
		principal Principal,
	) {
		response := ListModelsResponse{
			Models: allowlist.IDs(),
		}

		body, err := json.Marshal(response)
		if err != nil {
			writeSessionError(w, 500, "response_encoding_failed")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	}
}
