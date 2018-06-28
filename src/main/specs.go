package main

import (
	"push"
	"status"
)

// PollResponse is the response to /fetch access.
type PollResponse struct {
	Serv []status.Service
	Stat []status.Status
	Noti []push.Notification
}

// NotifyRequest is for callback notifications.
type NotifyRequest struct {
	Heading string
	Content string
}

// EncryptedData represents encrypted data.
type EncryptedData struct {
	EncryptionType string
	Data           string
}

// FetchRequest wow~!
type FetchRequest struct {
	Name     string
	StatOnly string
}
