package tmp

const HandlerWSTmp = `package  {{printf "%vHandler" (index . 0)}}

import (
	"net"
	"sync"
	{{printf "%vHelper \"%v/handlers/%v/helper\"" (index . 0) (index . 1) (index . 0)}}
	{{printf "\"%v/helper\"" (index . 1)}}
	{{printf "hubHelper \"%v/hub/helper\"" (index . 1)}}
)

type handler struct {
	hubHelper.HelperForHandler

	conn map[int64]map[string]net.Conn

	rw sync.Mutex
}

func InitHandler(hub hubHelper.HelperForHandler, conf *helper.HandlerConfig) (H {{print (index . 0)}}Helper.Handler, err error) {
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
