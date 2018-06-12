package main

import (
	"configrd"
	"logger"
	"net/http"
	"os"
	"push"
	"status"
)

var notif = push.NewNotificationMgr()
var servstat = status.NewServiceStatus()

func main() {
	var confile = configrd.Config(os.Args[1])
	conf := confile.ReadConfig()
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
	http.HandleFunc("/fetch", onFetch)
	http.HandleFunc("/callback", onCallback)
	logger.Log.Logln(logger.LEVEL_WARNING, "Listen", http.ListenAndServe(":"+conf["ListenPort"], nil))
}

func onFetch(w http.ResponseWriter, r *http.Request) {

}

func onCallback(w http.ResponseWriter, r *http.Request) {

}
