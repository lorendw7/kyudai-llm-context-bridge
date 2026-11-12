package main

import (
	"errors"
	"net/http"
	"strings"
	"unicode"
)

type Principal struct {
	LabID  string
	UserID string
}

func NewPrincipal(labID string, userID string) (Principal, error) {
	if strings.TrimSpace(labID) == "" {
		return Principal{}, errors.New("lab ID is required")
	}

	if strings.TrimSpace(userID) == "" {
		return Principal{}, errors.New("user ID is required")
	}

	return Principal{
		LabID:  labID,
		UserID: userID,
	}, nil
}

func (p Principal) CanAccessSession(session Session) bool {
	return session.OwnedBy(p.LabID, p.UserID)
}

type TokenAuthenticator struct {
	principals map[string]Principal
}

func NewTokenAuthenticator(
	principals map[string]Principal,
) (TokenAuthenticator, error) {
	if len(principals) == 0 {
		return TokenAuthenticator{}, errors.New("at least one token is required")
	}

	copied := make(map[string]Principal, len(principals))

	for token, principal := range principals {
		if token == "" {
			return TokenAuthenticator{}, errors.New("token is required")
		}

		if strings.IndexFunc(token, unicode.IsSpace) >= 0 {
			return TokenAuthenticator{}, errors.New("token must not contain whitespace")
		}

		validated, err := NewPrincipal(principal.LabID, principal.UserID)
		if err != nil {
			return TokenAuthenticator{}, err
		}

		copied[token] = validated
	}

	return TokenAuthenticator{
		principals: copied,
	}, nil
}

func (a TokenAuthenticator) Authenticate(
	r *http.Request,
) (Principal, error) {
	values := r.Header.Values("Authorization")
	if len(values) != 1 {
		return Principal{}, errors.New("invalid authorization")
	}

	parts := strings.Fields(values[0])
	if len(parts) != 2 {
		return Principal{}, errors.New("invalid authorization")
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return Principal{}, errors.New("invalid authorization")
	}

	principal, exists := a.principals[parts[1]]
	if !exists {
		return Principal{}, errors.New("invalid authorization")
	}

	return principal, nil
}
