package tmp

const HubHelperTmp = `package hubHelper{{$module := .ModuleName}}

import (
	"net/http" {{if gt (len .DBS) 0}}
	{{printf "\"%v/core/database\"" $module}}{{end}} {{range $i,$k := .Handlers}}
	{{printf "%vHelper \"%v/handlers/%v/helper\"" $i $module $i }}{{end}}
	{{if isOneTCP}}
	{{printf "\"%v/helper\"" $module}}
	"github.com/gorilla/mux"{{end}}
)

type Helper interface {
}

type HelperForHandler interface { {{if gt (len .DBS) 0}}
	database.DBForHandler {{end}} {{if isOneTCP}}
	Helper() Helper
	Config() *helper.Config
	Router() *mux.Router
	SetHandler(hh http.Handler){{end}}
}

type HandlerForHelper interface { {{range $i,$k := .Handlers}}
	{{printf "%v() %vHelper.Handler" (title $i) $i}}{{end}} {{if gt (len .DBS) 0}}
	database.DBForHandler {{end}}
	Config() *helper.Config
}

type help struct {
	HandlerForHelper
}

func InitHelper(H HandlerForHelper) Helper { return &help{HandlerForHelper: H} }`
