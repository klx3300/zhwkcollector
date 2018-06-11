package ssmgr

import (
	"encoding/json"
	"logger"
	"os"
)

// SubscriptionManager provides i/o functionality of subscription list.
type SubscriptionManager string

// This returns true self.
func (sm *SubscriptionManager) This() string {
	return string(*sm)
}

// ReadSubscription read subscription content from ss file.
func (sm *SubscriptionManager) ReadSubscription() SubscriptionList {
	file, err := os.Open(sm.This())
	if err != nil {
		panic("Unable to read subscription file" + sm.This() + err.Error())
	}
	defer file.Close()
	jdecoder := json.NewDecoder(file)
	cfg := make([]string, 0)
	err = jdecoder.Decode(&cfg)
	if err != nil {
		panic("Unable to unmarshal json while reading ss" + sm.This() + err.Error())
	}
	return cfg
}

// WriteSubscription write ss out to ss file.
func (sm *SubscriptionManager) WriteSubscription(cfg SubscriptionList) {
	file, err := os.OpenFile(sm.This(), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logger.Log.Logln(logger.LEVEL_WARNING, "Unable to write out to ", sm, err)
		return
	}
	defer file.Close()
	jencoder := json.NewEncoder(file)
	jencoder.Encode(cfg)
}
