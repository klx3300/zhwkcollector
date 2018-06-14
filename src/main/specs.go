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
