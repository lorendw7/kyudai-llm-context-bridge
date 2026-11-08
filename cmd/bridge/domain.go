package main

import "time"

type Message struct {
	ID        string
	SessionID string
	UserID    string
	Role      string
	Content   string
	CreatedAt time.Time
}

func NewMessage(
	id string,
	sessionID string,
	userID string,
	role string,
	content string,
) Message {
	return Message{
		ID:        id,
		SessionID: sessionID,
		UserID:    userID,
		Role:      role,
		Content:   content,
		CreatedAt: time.Now(),
	}
}
