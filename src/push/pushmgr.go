package push

import (
	"encoding/json"
	"logger"
	"os"
	"sync"
)

// NotificationMgr provides managements for pushdb.
type NotificationMgr struct {
	mu sync.Mutex
	nl NotificationList
}

// NewNotificationMgr returns new ones.
func NewNotificationMgr() *NotificationMgr {
	return &NotificationMgr{
		nl: NewNotificationList(),
	}
}

// OnAccess should be called on want to access.
func (nm *NotificationMgr) OnAccess() *NotificationList {
	nm.mu.Lock()
	return &(nm.nl)
}

// AfterAccess should be defered after the call to OnAccess.
func (nm *NotificationMgr) AfterAccess() {
	nm.mu.Unlock()
}

// Save save the content to file. for ft.
func (nm *NotificationMgr) Save(filename string) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logger.Log.Logln(logger.LEVEL_WARNING, "Unable to write to noti save file", filename, err)
		return
	}
	defer file.Close()
	jencoder := json.NewEncoder(file)
	nm.OnAccess()
	defer nm.AfterAccess()
	jencoder.Encode(nm.nl)
}

// Load load the previously saved content into memory.
func (nm *NotificationMgr) Load(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		logger.Log.Logln(logger.LEVEL_WARNING, "Unable to open previous noti save file", filename, err)
		return
	}
	defer file.Close()
	nm.OnAccess()
	defer nm.AfterAccess()
	jdecoder := json.NewDecoder(file)
	err = jdecoder.Decode(&(nm.nl))
	if err != nil {
		logger.Log.Logln(logger.LEVEL_WARNING, "Unable to unmarshal previous noti save file", filename, err)
		return
	}
}
