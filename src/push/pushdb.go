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

type NotificationAck struct {
	Noti           Notification
	UnAckedNames   []string
	NamelessRemain int
}

// NotificationList is the list of notifications.
type NotificationList []NotificationAck

// This returns true self.
func (nl *NotificationList) This() []NotificationAck {
	return []NotificationAck(*nl)
}

// NewNotificationList makes new ones.
func NewNotificationList() NotificationList {
	return make([]NotificationAck, 0)
}

// Append will behaves like the append, time automatically generated.
func (nl *NotificationList) Append(heading, content string, dest []string, nameless int) {
	*nl = append(nl.This(), NotificationAck{
		Noti: Notification{
			Tm:      time.Now(),
			Heading: heading,
			Content: content},
		UnAckedNames:   make([]string, 0),
		NamelessRemain: nameless,
	})
}
