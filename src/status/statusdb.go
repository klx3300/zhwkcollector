package status

import (
	"encoding/json"
	"net/http"
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
	defer ss.AfterAccess()
	for serv := range ss.data {
		resp, err := http.Get(serv.Addr)
		if err != nil {
			// boomed.
			ss.data[serv] = Status{
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
			ss.data[serv] = Status{
				Tm:         time.Now(),
				Stat:       STATE_DCONN,
				Informator: "Invalid Response.",
			}
		}
		_, mmok := respmap["inform"]
		if !mmok {
			ss.data[serv] = Status{
				Tm:         time.Now(),
				Stat:       STATE_DCONN,
				Informator: "Invalid Response.",
			}
		}
		vint, converr := strconv.Atoi(intstr)
		if converr != nil {
			ss.data[serv] = Status{
				Tm:         time.Now(),
				Stat:       STATE_DCONN,
				Informator: "Invalid Response.",
			}
		}
		ss.data[serv] = Status{
			Tm:         time.Now(),
			Stat:       vint,
			Informator: respmap["inform"],
		}
	}
}
