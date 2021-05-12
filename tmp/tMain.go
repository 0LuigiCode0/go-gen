package tmp

const MainTmp = `package main

import (
	{{printf "\"%v/core\"" .ModuleName}}
	{{printf "\"%v/helper\"" .ModuleName}}
)

func main() {
	helper.InitLogger()
	helper.InitCtx()
	conf, err := helper.ParseConfig()
	if err != nil {
		helper.Log.Fatalf("config parse invalid: %v", err)
	}
	srv, err := core.InitServer(conf)
	if err != nil {
		helper.Log.Fatalf("server not initialized: %v", err)
	}
	if err := srv.Start(); err != nil {
		helper.Log.Fatalf("server not started: %v", err)
	}
	srv.Close()
	helper.Wg.Wait()
}`
