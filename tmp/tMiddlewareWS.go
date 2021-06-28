package tmp

const MiddleWareWSTmp = `package {{printf "%vHandler" (index . 0)}}

import (
	"net/http"
	{{printf "\"%v/helper\"" (index . 1)}}
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func (h *handler) WS(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		helper.Log.Warningf(helper.KeyErrorUpdate+": ws connected: %v", err)
		return
	}
	ipAddr := conn.RemoteAddr().String()
	_, _, err = wsutil.ReadClientData(conn)
	if err != nil {
		helper.Log.Warning(helper.KeyErrorRead+": ws msg: %v", err)
		conn.Close()
		return
	}

	h.addConn(0, ipAddr, conn)
}

func (h *handler) initTicker() {
	defer helper.Wg.Done()

	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-helper.Ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			for u, conns := range h.conn {
				for ip, conn := range conns {
					if err := wsutil.WriteServerText(conn, []byte{'1'}); err != nil {
						h.DeleteConn(u, ip)
					}
				}
			}
		}
	}
}`
