package main

import (
	"push"
	"status"
)

// PollResponse is the response to /fetch access.
type PollResponse struct {
	Serv map[status.Service]status.Status
	Noti []push.Notification
}

// NotifyRequest is for callback notifications.
type NotifyRequest struct {
	Heading string
	Content string
}
