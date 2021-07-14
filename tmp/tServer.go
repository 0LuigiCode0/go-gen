package tmp

const ServerTmp = `package core

import (
	"fmt"
	"net/http"
	"os"
	"os/signal" {{if gt (len .DBS) 0}}
	{{printf "\"%v/core/database\"" .ModuleName}}{{end}}
	{{printf "\"%v/helper\"" .ModuleName}}
	{{printf "\"%v/hub\"" .ModuleName}}
)

type Server interface {
	Start() error
	Close()
}

type server struct {
	srv http.Server
	hub hub.Hub {{if gt (len .DBS) 0}}
	db  database.DB{{end}}
}

func InitServer(conf *helper.Config) (S Server, err error) {
	s := &server{}
	S = s {{if gt (len .DBS) 0}}
	s.db, err = database.InitDB(conf)
	if err != nil {
		s.srv.Close()
		return nil, fmt.Errorf("db not initialized: %v", err)
	} {{end}}
	s.hub, err = hub.InitHub({{if gt (len .DBS) 0}}s.db, {{end}}conf)
	if err != nil {
		return nil, fmt.Errorf("hub not initialized: %v", err)
	}

	s.srv.Handler = s.hub.GetHandler()
	s.srv.Addr = fmt.Sprintf("%v:%v", conf.Host, conf.Port)

	helper.Log.Service("server initialized")
	return s, nil
}

func (s *server) Start() error {
	signal.Notify(helper.C, os.Interrupt)

	helper.Wg.Add(1)
	go s.loop()

	helper.Wg.Add(1)
	go func() {
		defer helper.Wg.Done()
		if err := s.srv.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				helper.Log.Service("serve stoped")
				helper.C <- os.Interrupt
				return
			}
			helper.Log.Errorf("serve error: %v", err)
			helper.C <- os.Interrupt
			return
		}
	}()

	helper.Log.Service("server started at address:", s.srv.Addr)
	<-helper.C
	return nil
}

func (s *server) loop() {
	defer helper.Wg.Done()
	for {
		select {
		case <-helper.Ctx.Done():
			return
		default:
			if err := recover(); err != nil {
				helper.Log.Errorf("critical damage: %v", err)
			}
		}
	}
}

func (s *server) Close() {
	s.srv.Shutdown(helper.Ctx)
	s.hub.Close() {{if gt (len .DBS) 0}}
	s.db.Close() {{end}}
	helper.CloseCtx()
}`
