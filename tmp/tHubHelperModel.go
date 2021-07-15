package tmp

const HubHelperModelTmp = `package hub_helper{{$module := .ModuleName}}

import (
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
}`
