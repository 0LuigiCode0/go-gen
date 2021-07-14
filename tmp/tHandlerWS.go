package tmp

const HandlerWSTmp = `package  {{printf "%v_handler" (index . 0)}}

import (
	"net"
	"sync"
	{{printf "\"%v/handlers/%v_handler/%v_helper\"" (index . 1) (index . 0) (index . 0)}}
	{{printf "\"%v/helper\"" (index . 1)}}
	{{printf "\"%v/hub/hub_helper\"" (index . 1)}}
)

type handler struct {
	hub_helper.HelperForHandler

	conn map[int64]map[string]net.Conn

	rw sync.Mutex
}

func InitHandler(hub hub_helper.HelperForHandler, conf *helper.HandlerConfig) (H {{print (index . 0)}}_helper.Handler, err error) {
	h := &handler{
		HelperForHandler: hub,
		conn:             map[int64]map[string]net.Conn{},
	}
	H = h

	helper.Wg.Add(1)
	go h.initTicker()

	h.Router().HandleFunc("/ws", h.WS)
	return
}

func (h *handler) addConn(userID int64, ipAddr string, connect net.Conn) {
	h.rw.Lock()
	defer h.rw.Unlock()
	if _, ok := h.conn[userID]; ok {
		h.conn[userID][ipAddr] = connect
	} else {
		h.conn[userID] = make(map[string]net.Conn)
		h.conn[userID][ipAddr] = connect
	}
}

func (h *handler) DeleteConn(userID int64, ipAddr string) {
	h.rw.Lock()
	defer h.rw.Unlock()
	if _, ok := h.conn[userID]; ok {
		delete(h.conn[userID], ipAddr)

		if len(h.conn[userID]) == 0 {
			delete(h.conn, userID)
		}
	}
}

func (h *handler) GetConn(userID int64) map[string]net.Conn {
	h.rw.Lock()
	defer h.rw.Unlock()
	res := h.conn[userID]
	return res
}

func (h *handler) GetAll() map[int64]map[string]net.Conn { return h.conn }`
