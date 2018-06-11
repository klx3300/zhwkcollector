package push

import (
	"time"
)

// Notification represents a notification that to be displayed at smartphones.
type Notification struct {
	Tm      time.Time
	Heading string
	Content string
}

// NotificationList is the list of notifications.
type NotificationList []Notification

// This returns true self.
func (nl *NotificationList) This() []Notification {
	return []Notification(*nl)
}

// NewNotificationList makes new ones.
func NewNotificationList() NotificationList {
	return make([]Notification, 0)
}

// Append will behaves like the append, time automatically generated.
func (nl *NotificationList) Append(heading, content string) {
	*nl = append(nl.This(), Notification{
		Tm:      time.Now(),
		Heading: heading,
		Content: content,
	})
}
