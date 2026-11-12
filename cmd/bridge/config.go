package main

import (
	"os"
	"strings"
)

func authFromEnvironment() (TokenAuthenticator, error) {
	return NewTokenAuthenticator(map[string]Principal{
		os.Getenv("BRIDGE_TOKEN"): {
			LabID:  os.Getenv("BRIDGE_LAB_ID"),
			UserID: os.Getenv("BRIDGE_USER_ID"),
		},
	})
}
func modelsFromEnvironment() (ModelAllowlist, error) {
	raw := os.Getenv("BRIDGE_MODEL_IDS")
	modelIDs := strings.Split(raw, ",")

	return NewModelAllowlist(modelIDs)
}
