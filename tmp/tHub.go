package tmp

const HubTmp = `package hub{{$module := .ModuleName}}

import (
	"fmt"
	"net/http" {{if gt (len .DBS) 0}}
	{{printf "\"%v/core/database\"" $module}}{{end}} {{range $i,$k := .Handlers}}
	{{printf "%vHandler \"%v/handlers/%v\"" $i $module $i}}
	{{printf "%vHelper \"%v/handlers/%v/helper\"" $i $module $i}}{{end}}
	{{printf "hubHelper \"%v/hub/helper\"" $module}}
	{{printf "\"%v/helper\"" $module}}

	"github.com/gorilla/mux"
)

const ( {{range $i,$k := .Handlers}}
	{{printf "_%v = \"%v\"" $i $i}}{{end}}
)

type hub struct {
	helper hubHelper.Helper {{if gt (len .DBS) 0}}
	database.DBForHandler{{end}}
	router  *mux.Router
	handler http.Handler
	config  *helper.Config
	{{range $i,$k := .Handlers}}
	{{printf "_%v %vHelper.Handler" $i $i}}{{end}}
}

func InitHub({{if gt (len .DBS) 0}}db database.DB, {{end}}conf *helper.Config) (H http.Handler, err error) {
	hh := &hub{ {{if gt (len .DBS) 0}}
		DBForHandler: db,{{end}}
		router:       mux.NewRouter(),
		config:       conf,
	}
	hh.SetHandler(hh.router)

	{{range $i,$k := .Handlers}}
	if v, ok := conf.Handlers[{{printf "_%v" $i}}]; ok {
		hh.{{printf "_%v" $i}}, err = {{print $i}}Handler.InitHandler(hh, v)
		if err != nil {
			return nil, fmt.Errorf("handler not initializing: %v", err)
		}
		helper.Log.Servicef("handler %q initializing", {{printf "_%v" $i}})
	}{{end}}

	hh.helper = hubHelper.InitHelper(hh)

	H = hh.handler
	helper.Log.Service("handler initializing")
	return
}

func (h *hub) Config() *helper.Config     { return h.config }
func (h *hub) Helper() hubHelper.Helper   { return h.helper }
func (h *hub) Router() *mux.Router        { return h.router }
func (h *hub) SetHandler(hh http.Handler) { h.handler = hh } {{range $i,$k := .Handlers}}
{{printf "func (h *hub) %v() %vHelper.Handler { return h._%v }" (title $i) $i $i}}{{end}}`
