package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CreateSessionRequest struct {
	SchemaVersion string `json:"schema_version"`
	PolicyID      string `json:"policy_id"`
}

func (r CreateSessionRequest) Validate() error {
	if r.SchemaVersion != "1" {
		return errors.New("unsupported schema version")
	}

	if r.PolicyID != "full-history-v1" {
		return errors.New("unsupported policy")
	}

	return nil
}

type CreateSessionResponse struct {
	ID        string    `json:"id"`
	PolicyID  string    `json:"policy_id"`
	Revision  uint64    `json:"revision"`
	CreatedAt time.Time `json:"created_at"`
}

func writeSessionError(
	w http.ResponseWriter,
	status int,
	code string,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, "{\"code\":%q}\n", code)
}

func decodeCreateSessionRequest(
	r *http.Request,
) (CreateSessionRequest, error) {
	var request CreateSessionRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		return CreateSessionRequest{}, err
	}

	var extra any
	err := decoder.Decode(&extra)

	if err != io.EOF {
		if err == nil {
			err = errors.New("request body must contain one JSON value")
		}
		return CreateSessionRequest{}, err
	}

	if err := request.Validate(); err != nil {
		return CreateSessionRequest{}, err
	}

	return request, nil
}

func writeSessionBodyError(
	w http.ResponseWriter,
	err error,
) {
	var tooLarge *http.MaxBytesError

	if errors.As(err, &tooLarge) {
		writeSessionError(w, 413, "request_too_large")
		return
	}

	writeSessionError(w, 422, "invalid_request")
}

func newSessionID() (string, error) {
	var randomBytes [16]byte

	if _, err := rand.Read(randomBytes[:]); err != nil {
		return "", err
	}

	return hex.EncodeToString(randomBytes[:]), nil
}

func newOwnedSession(principal Principal,
	request CreateSessionRequest) (Session, error) {
	sessionID, err := newSessionID()

	if err != nil {
		return Session{}, err
	}

	return NewSession(
		sessionID,
		principal.LabID,
		principal.UserID,
		request.PolicyID,
		time.Now().UTC(),
	)
}

func newCreateSessionHandler(
	store *MemorySessionStore,
) authenticatedHandler {
	return func(
		w http.ResponseWriter,
		r *http.Request,
		principal Principal,
	) {
		r.Body = http.MaxBytesReader(w, r.Body, 4096)

		request, err := decodeCreateSessionRequest(r)

		if err != nil {
			writeSessionBodyError(w, err)
			return
		}

		session, err := newOwnedSession(principal, request)
		if err != nil {
			writeSessionError(w, 500, "session_creation_failed")
			return
		}

		if err := store.Add(session); err != nil {
			writeSessionError(w, 500, "session_creation_failed")
			return
		}

		response := CreateSessionResponse{
			ID:        session.ID,
			PolicyID:  session.PolicyID,
			Revision:  session.Revision,
			CreatedAt: session.CreatedAt,
		}

		body, err := json.Marshal(response)
		if err != nil {
			writeSessionError(w, 500, "response_encoding_failed")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", "/v1/sessions/"+session.ID)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(body)
	}
}

func newGetSessionHandler(
	store *MemorySessionStore,
) authenticatedHandler {
	return func(w http.ResponseWriter, r *http.Request, principal Principal) {
		id := r.PathValue("id")
		session, exists := store.Get(id)

		if !exists || !principal.CanAccessSession(session) {
			writeSessionError(w, 404, "session_not_found")
			return
		}

		response := CreateSessionResponse{
			ID:        session.ID,
			PolicyID:  session.PolicyID,
			Revision:  session.Revision,
			CreatedAt: session.CreatedAt,
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

type ListSessionsResponse struct {
	Sessions []CreateSessionResponse `json:"sessions"`
}

func buildListSessionsResponse(
	sessions []Session,
) ListSessionsResponse {
	response := ListSessionsResponse{
		Sessions: make([]CreateSessionResponse, 0),
	}

	for _, session := range sessions {
		response.Sessions = append(response.Sessions, CreateSessionResponse{
			ID:        session.ID,
			PolicyID:  session.PolicyID,
			Revision:  session.Revision,
			CreatedAt: session.CreatedAt,
		})
	}

	return response
}

func newListSessionsHandler(
	store *MemorySessionStore,
) authenticatedHandler {
	return func(
		w http.ResponseWriter,
		r *http.Request,
		principal Principal,
	) {
		sessions := store.ListOwned(principal)
		response := buildListSessionsResponse(sessions)

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
