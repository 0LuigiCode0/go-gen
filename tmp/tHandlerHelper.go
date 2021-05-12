package tmp

const HandlerHelperTmp = `package {{printf "%vHandler" (index . 0)}}
{{if eq (index . 1) "` + WS + `"}}
import "net"
{{end}}
type Handler interface { {{if eq (index . 1) "` + WS + `"}}
	GetConn(userID int64) map[string]net.Conn {{end}}
}`
