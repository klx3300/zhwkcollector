package status

import (
	"encoding/json"
	"logger"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// Service represents a service server.
type Service struct {
	Name string
	Addr string
}

// Status represents the status of a subscribed status server.
type Status struct {
	Tm         time.Time
	Stat       int
	Informator string
}

// Possible int status list
const (
	STATE_GOOD  = 0
	STATE_PROB  = 1
	STATE_WARN  = 2
	STATE_FAIL  = 3
	STATE_DCONN = 4
)

// ServiceStatus represents overall status of servers.
type ServiceStatus struct {
	data map[Service]Status
	mu   sync.Mutex
}

// NewServiceStatus returns new ones.
func NewServiceStatus() ServiceStatus {
	return ServiceStatus{
		data: make(map[Service]Status),
	}
}

// OnAccess keeps sync.
func (ss *ServiceStatus) OnAccess() *map[Service]Status {
	ss.mu.Lock()
	return &(ss.data)
}

// AfterAccess should be atleast defered.
func (ss *ServiceStatus) AfterAccess() {
	ss.mu.Unlock()
}

// Watch append a new watched services, or update already existed.
func (ss *ServiceStatus) Watch(newserv Service) {
	ss.OnAccess()
	defer ss.AfterAccess()
	for serv := range ss.data {
		if serv.Addr == newserv.Addr {
			delete(ss.data, serv)
			ss.data[newserv] = Status{}
			return
		}
	}
	ss.data[newserv] = Status{}
}

// Unwatch remove from watching list.
func (ss *ServiceStatus) Unwatch(servname string) {
	ss.OnAccess()
	defer ss.AfterAccess()
	for serv := range ss.data {
		if serv.Name == servname {
			delete(ss.data, serv)
		}
	}
}

// Refresh update all informations from services
func (ss *ServiceStatus) Refresh() {
	ss.OnAccess()
	// copy the map
	buffer := make(map[Service]Status)
	for k, v := range ss.data {
		buffer[k] = v
	}
	ss.AfterAccess()
	for serv := range buffer {
		resp, err := http.Get(serv.Addr)
		if err != nil {
			// boomed.
			buffer[serv] = Status{
				Tm:         time.Now(),
				Stat:       STATE_DCONN,
				Informator: "HTTP failed: " + err.Error(),
			}
			continue
		}
		defer resp.Body.Close()
		jdecoder := json.NewDecoder(resp.Body)
		respmap := make(map[string]string)
		jdecoder.Decode(&respmap)
		intstr, mok := respmap["status"]
		if !mok {
			buffer[serv] = Status{
				Tm:         time.Now(),
				Stat:       STATE_DCONN,
				Informator: "Invalid Response.",
			}
			continue
		}
		_, mmok := respmap["inform"]
		if !mmok {
			buffer[serv] = Status{
				Tm:         time.Now(),
				Stat:       STATE_DCONN,
				Informator: "Invalid Response.",
			}
			continue
		}
		vint, converr := strconv.Atoi(intstr)
		if converr != nil {
			buffer[serv] = Status{
				Tm:         time.Now(),
				Stat:       STATE_DCONN,
				Informator: "Invalid Response.",
			}
			continue
		}
		buffer[serv] = Status{
			Tm:         time.Now(),
			Stat:       vint,
			Informator: respmap["inform"],
		}
	}
	ss.OnAccess()
	defer ss.AfterAccess()
	for k, v := range buffer {
		ss.data[k] = v
	}
}

// Save to file.
func (ss *ServiceStatus) Save(fn string) {
	ss.OnAccess()
	defer ss.AfterAccess()
	buffer := make([]Service, 0)
	for k := range ss.data {
		buffer = append(buffer, k)
	}
	file, err := os.OpenFile(fn, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		logger.Log.Logln(logger.LEVEL_WARNING, "Unable to open service save file", fn, err)
		return
	}
	defer file.Close()
	jencoder := json.NewEncoder(file)
	jencoder.Encode(file)
}

// Load from file. Need refresh.
func (ss *ServiceStatus) Load(fn string) {
	ss.OnAccess()
	defer ss.AfterAccess()
	buffer := make([]Service, 0)
	ss.data = make(map[Service]Status)
	file, err := os.Open(fn)
	if err != nil {
		logger.Log.Logln(logger.LEVEL_WARNING, "Unable to open service save file", fn, err)
		return
	}
	defer file.Close()
	jdecoder := json.NewDecoder(file)
	err = jdecoder.Decode(&buffer)
	if err != nil {
		logger.Log.Logln(logger.LEVEL_WARNING, "Unable to unmarshal service save file", fn, err)
		return
	}
	for _, x := range buffer {
		ss.data[x] = Status{}
	}
}
