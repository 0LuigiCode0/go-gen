package tmp

const HubTmp = `package hub{{$module := .ModuleName}}

import (
	"fmt"
	"net/http" {{if gt (len .DBS) 0}}
	{{printf "\"%v/core/database\"" $module}}{{end}} {{range $i,$k := .Handlers}}
	{{printf "\"%v/handlers/%v_handler\"" $module $i}}
	{{printf "\"%v/handlers/%v_handler/%v_helper\"" $module $i $i}}{{end}}
	{{printf "\"%v/hub/hub_helper\"" $module}}
	{{printf "\"%v/helper\"" $module}}

	"github.com/gorilla/mux"
)

const ( {{range $i,$k := .Handlers}}
	{{printf "_%v = \"%v\"" $i $i}}{{end}}
)

type Hub interface {
	GetHandler() http.Handler
	Close()
}

type hub struct {
	helper hub_helper.Helper {{if gt (len .DBS) 0}}
	database.DB{{end}}
	router  *mux.Router
	handler http.Handler
	config  *helper.Config
	{{range $i,$k := .Handlers}}
	{{printf "_%v %v_helper.Handler" $i $i}}{{end}}
}

func InitHub({{if gt (len .DBS) 0}}db database.DB, {{end}}conf *helper.Config) (H Hub, err error) {
	hh := &hub{ {{if gt (len .DBS) 0}}
		DB: db,{{end}}
		router:       mux.NewRouter(),
		config:       conf,
	}
	hh.SetHandler(hh.router)

	hh.helper = hub_helper.InitHelper(hh)

	if err = hh.intiDefault(); err != nil {
		helper.Log.Warningf("initializing default is failed: %v", err)
		return nil, fmt.Errorf("handler not initializing: %v", err)
	}
	helper.Log.Service("initializing default")

	{{range $i,$k := .Handlers}}
	if v, ok := conf.Handlers[{{printf "_%v" $i}}]; ok {
		hh.{{printf "_%v" $i}}, err = {{print $i}}_handler.InitHandler(hh, v)
		if err != nil {
			return nil, fmt.Errorf("handler not initializing: %v", err)
		}
		helper.Log.Servicef("handler %q initializing", {{printf "_%v" $i}})
	}{{end}}

	hh.Router().PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(helper.UploadDir))))

	H = hh
	helper.Log.Service("handler initializing")
	return
}

func (h *hub) Config() *helper.Config     { return h.config }
func (h *hub) Helper() hub_helper.Helper   { return h.helper }
func (h *hub) Router() *mux.Router        { return h.router }
func (h *hub) SetHandler(hh http.Handler) { h.handler = hh }
func (h *hub) GetHandler() http.Handler    { return h.handler } {{range $i,$k := .Handlers}}
{{printf "func (h *hub) %v() %v_helper.Handler { return h._%v }" (title $i) $i $i}}{{end}}
func (h *hub) Close() {}

func (h *hub) intiDefault() error { return nil }`
