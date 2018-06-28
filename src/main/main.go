package main

import (
	"configrd"
	"ctool"
	"encoding/json"
	"logger"
	"net/http"
	"os"
	"push"
	"status"
	"strconv"
	"time"
)

var notif = push.NewNotificationMgr()
var servstat = status.NewServiceStatus()
var conf map[string]string
var timeout int
var fetchSym, cbSym, ctrlSym, pollSym string
var fetchKey, cbKey, ctrlKey, pollKey, aesIv []byte
var acceptBare = false

func main() {
	var confile = configrd.Config(os.Args[1])
	conf = confile.ReadConfig()
	_, mpok := conf["NotificationSave"]
	if !mpok {
		panic("NotificationSave has to be exist in config map!")
	}
	notif.Load(conf["NotificationSave"])
	_, mpok = conf["ServiceSave"]
	if !mpok {
		panic("ServiceSave has to be exist in config map!")
	}
	servstat.Load(conf["ServiceSave"])
	_, mpok = conf["ListenPort"]
	if !mpok {
		panic("ListenPort has to be exist in config map!")
	}
	_, mpok = conf["Timeout"]
	if !mpok {
		panic("Timeout has to be exist in config map!")
	}
	var err error
	timeout, err = strconv.Atoi(conf["Timeout"])
	if err != nil {
		panic("Timeout is not a valid integer!")
	}
	fetchSym, mpok = conf["FetchSymKey"]
	if !mpok {
		panic("FetchSymKey not speced.")
	}
	fetchKey = ctool.AESKeyRound(fetchSym)
	cbSym, mpok = conf["CallbackSymKey"]
	if !mpok {
		panic("CallbackSymKey not speced.")
	}
	cbKey = ctool.AESKeyRound(cbSym)
	ctrlSym, mpok = conf["ControlSymKey"]
	if !mpok {
		panic("ControlSymKey not speced.")
	}
	ctrlKey = ctool.AESKeyRound(ctrlSym)
	pollSym, mpok = conf["ServicePollSymKey"]
	if !mpok {
		panic("ServicePollSymKey not speced.")
	}
	pollKey = ctool.AESKeyRound(pollSym)
	var aiv string
	aiv, mpok = conf["AESIV"]
	if !mpok {
		aesIv = ctool.AESIV(aiv)
	} else {
		aesIv = ctool.AESIV("IVECTOR")
	}
	_, mpok = conf["AcceptBareConn"]
	if !mpok {
		logger.Log.Logln(logger.LEVEL_WARNING, "AcceptBareConn unspeced. default to FALSE.")
	} else {
		if conf["AcceptBareConn"] == "true" {
			acceptBare = true
		}
	}
	http.HandleFunc("/fetch", onFetch)
	http.HandleFunc("/callback", onCallback)
	http.HandleFunc("/register", onRegister)
	http.HandleFunc("/revoke", onRevoke)
	go onTimeout()
	logger.Log.Logln(logger.LEVEL_WARNING, "Listen", http.ListenAndServe(":"+conf["ListenPort"], nil))
}

func onTimeout() {
	for {
		<-time.After(time.Duration(timeout) * time.Millisecond)
		servstat.Refresh()
	}
}

func onFetch(w http.ResponseWriter, r *http.Request) {
	baremode := false
	fr := FetchRequest{}
	jdecoder := json.NewDecoder(r.Body)
	err := jdecoder.Decode(&fr)
	if err != nil {
		// check fallback permit
		if acceptBare {
			baremode = true
		} else {
			logger.Log.Logln(logger.LEVEL_WARNING, "Attempt to access with bare mode, denied.")
			return
		}
	}
	resp := PollResponse{
		Serv: make([]status.Service, 0),
		Stat: make([]status.Status, 0),
		Noti: make([]push.Notification, 0),
	}
	// some copy work here..
	mss := servstat.OnAccess()
	for k, v := range *mss {
		resp.Serv = append(resp.Serv, k)
		resp.Stat = append(resp.Stat, v)
	}
	servstat.AfterAccess()
	pmsg := notif.OnAccess()
	for _, item := range *pmsg {
		resp.Noti = append(resp.Noti, item)
	}
	notif.AfterAccess()
	// we can't simply say goodbye now.. shame
	// notif = push.NewNotificationMgr() // goodbye!
	notif.Save(conf["NotificationSave"])
	jencoder := json.NewEncoder(w)
	err := jencoder.Encode(resp)
	if err != nil {
		logger.Log.Logln(logger.LEVEL_FATAL, "Unable to response,", err)
	}
}

func onCallback(w http.ResponseWriter, r *http.Request) {
	nr := NotifyRequest{}
	jdecoder := json.NewDecoder(r.Body)
	err := jdecoder.Decode(&nr)
	logger.Log.Logln(logger.LEVEL_INFO, "NotiCb", nr)
	if err != nil {
		logger.Log.Logln(logger.LEVEL_WARNING, "Unable to unmarshal callback, ", err)
		return
	}
	nl := notif.OnAccess()
	nl.Append(nr.Heading, nr.Content)
	notif.AfterAccess()
	notif.Save(conf["NotificationSave"])
}

func onRegister(w http.ResponseWriter, r *http.Request) {
	sr := status.Service{}
	jdecoder := json.NewDecoder(r.Body)
	err := jdecoder.Decode(&sr)
	if err != nil {
		logger.Log.Logln(logger.LEVEL_WARNING, "Unable to unmarshal register", err)
		return
	}
	servstat.Watch(sr)
	servstat.Save(conf["ServiceSave"])
}

func onRevoke(w http.ResponseWriter, r *http.Request) {
	var srvname string
	jdecoder := json.NewDecoder(r.Body)
	err := jdecoder.Decode(&srvname)
	if err != nil {
		logger.Log.Logln(logger.LEVEL_WARNING, "Unable to unmarshal revoke", err)
	}
	servstat.Unwatch(srvname)
	servstat.Save(conf["ServiceSave"])
}
