package main

import (
	"errors"
	"strings"
)

type ModelAllowlist struct {
	models map[string]struct{}
}

func NewModelAllowlist(modelIDs []string) (ModelAllowlist, error) {
	if len(modelIDs) == 0 {
		return ModelAllowlist{}, errors.New("at least one model ID is required")
	}

	models := make(map[string]struct{}, len(modelIDs))

	for _, modelID := range modelIDs {
		if strings.TrimSpace(modelID) == "" {
			return ModelAllowlist{}, errors.New("model ID is required")
		}

		if strings.TrimSpace(modelID) != modelID {
			return ModelAllowlist{}, errors.New("model ID must not have surrounding whitespace")
		}

		if _, exist := models[modelID]; exist {
			return ModelAllowlist{}, errors.New("duplicate model ID")
		}

		models[modelID] = struct{}{}
	}

	return ModelAllowlist{
		models: models,
	}, nil
}

func (a ModelAllowlist) Allows(modelID string) bool {
	_, exists := a.models[modelID]
	return exists
}
