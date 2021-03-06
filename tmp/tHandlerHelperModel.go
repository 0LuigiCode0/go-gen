package tmp

const HandlerHelperModelTmp = `package {{printf "%v_helper" (index . 0)}}
{{if eq (index . 1) "` + WS + `"}}
import "net"
{{end}}
type Handler interface { {{if eq (index . 1) "` + WS + `"}}
	GetAll() map[int64]map[string]net.Conn
	GetConn(userID int64) map[string]net.Conn 
	DeleteConn(userID int64, ipAddr string){{end}}
}`
