package tmp

const HubHelperTmp = `package hub_helper{{$module := .ModuleName}}

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"net/http" {{if gt (len .DBS) 0}}
	{{printf "\"%v/core/database\"" $module}}{{end}} {{range $i,$k := .Handlers}}
	{{printf "\"%v/handlers/%v_handler/%v_helper\"" $module $i $i }}{{end}}
	{{printf "\"%v/helper\"" $module}}
	
	"github.com/gorilla/mux"
)

type Helper interface {
	Hash(data string) string
}

type HelperForHandler interface { {{if gt (len .DBS) 0}}
	database.DBForHandler {{end}}
	Helper() Helper
	Config() *helper.Config
	Router() *mux.Router
	SetHandler(hh http.Handler)
}

type HandlerForHelper interface { {{range $i,$k := .Handlers}}
	{{printf "%v() %v_helper.Handler" (title $i) $i}}{{end}} {{if gt (len .DBS) 0}}
	database.DBForHandler {{end}}
	Config() *helper.Config
}

type help struct {
	HandlerForHelper
}

func InitHelper(H HandlerForHelper) Helper { return &help{HandlerForHelper: H} }

func (h *help) Hash(data string) string {
	hash1 := sha256.New()
	hash2 := md5.New()
	hash1.Write([]byte(data))
	hash2.Write([]byte(data))
	hash1.Write([]byte(string(hash1.Sum(nil)) + string(hash2.Sum(nil))))
	return hex.EncodeToString(hash1.Sum(nil))
}`
