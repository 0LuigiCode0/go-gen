package tmp

const (
	DirCore     = "core"
	DirCmd      = "cmd"
	DirDatabase = "database"
	DirHub      = "hub"
	DirHelper   = "helper"
	DirHandlers = "handlers"
	DirStore    = "store"
	DirSource   = "source"
	DirConfigs  = "configs"
	DirUploads  = "uploads"
)
const (
	FileMain          = "main.go"
	FileDatabase      = "database.go"
	FileServer        = "server.go"
	FileModel         = "model.go"
	FileFunctions     = "functions.go"
	FileHelper        = "helper.go"
	FileHandler       = "handler.go"
	FileHubMiddleware = "middleware.go"
	FileStore         = "store.go"
	FileHub           = "hub.go"
	FileConfigServer  = "configServer.json"
	FileMod           = "go.mod"
	FileSum           = "go.sum"
	FileDocker        = "Dockerfile"
	FileComposeBuild  = "docker-compose-build.yaml"
	FileComposeLocal  = "docker-compose-local.yaml"
	FileReadme        = "README.md"
)

type DBType string

const (
	Postgres DBType = "postgres"
	Mongodb  DBType = "mongodb"
)

type HandlerType string

const (
	TCP  HandlerType = "tcp"
	MQTT HandlerType = "mqtt"
	WS   HandlerType = "ws"
)

const (
	MainTmp = `package main

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
	ServerTmp = `package core

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
	srv http.Server {{if gt (len .DBS) 0}}
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
	s.srv.Handler, err = hub.InitHub({{if gt (len .DBS) 0}}s.db, {{end}}conf)
	if err != nil {
		return nil, fmt.Errorf("hub not initialized: %v", err)
	}

	s.srv.Addr = fmt.Sprintf("%v:%v", conf.Host, conf.Port)

	helper.Log.Service("server initialized")
	return s, nil
}

func (s *server) Start() error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	helper.Wg.Add(1)
	go func() {
		defer helper.Wg.Done()
		if err := s.srv.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				helper.Log.Service("serve stoped")
				return
			}
			helper.Log.Errorf("serve error: %v", err)
			return
		}
	}()

	helper.Log.Service("server started at address:", s.srv.Addr)
	<-c
	return nil
}

func (s *server) Close() {
	s.srv.Shutdown(helper.Ctx) {{if gt (len .DBS) 0}}
	s.db.Close() {{end}}
	helper.CloseCtx()
}`
	DatabaseTmp = `package database{{$module := .ModuleName}}

import (
	"database/sql"
	"fmt"
	{{printf "\"%v/helper\"" $module}}
	{{range $i,$k := .DBS}}
	{{printf "%vStore \"%v/store/%v\"" $i $module $i}}{{end}}
	{{if isOneMongo}}
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"{{end}}{{if isOnePostgres}}
	_ "github.com/lib/pq"{{end}}
)

const ( {{range $i,$k := .DBS}}
	{{printf "_%v = \"%v\"" $i $i}}{{end}}
)

type DB interface { {{range $i,$k := .DBS}}{{if eq $k "` + string(Mongodb) + `"}}
	{{printf "%v() *mongo.Database" (title $i)}}{{else if eq $k "` + string(Postgres) + `"}}
	{{printf "%v() *sql.DB" (title $i)}}{{end}}
	{{printf "%vStore() %vStore.Store" (title $i) $i}}{{end}}
	Close()
}

type DBForHandler interface { {{range $i,$k := .DBS}}
	{{printf "%vStore() %vStore.Store" (title $i) $i}}{{end}}
}

type db struct { {{range $i,$k := .DBS}}
	{{printf "_%v *d" $i}}{{end}}
}

type d struct {
	store  interface{}
	conn   interface{}
	dbName string
}

func InitDB(conf *helper.Config) (DB DB, err error) {
	db := &db{}
	DB = db

	var conn interface{} {{range $i,$k := .DBS}}
	if v, ok := conf.DBS[{{printf "_%v" $i}}]; ok {
		conn, err = {{if eq $k "` + string(Mongodb) + `"}}connMongo(v){{else if eq $k "` + string(Postgres) + `"}}connPostgres(v){{end}}
		if err != nil {
			return nil, fmt.Errorf("db not initializing: %v", err)
		}
		db.{{printf "_%v" $i}} = &d{
			store:  {{print $i}}Store.InitStore(),
			dbName: v.Name,
			conn:   conn,
		}
		helper.Log.Servicef("db %q initializing", {{printf "_%v" $i}})
	}{{end}}

	helper.Log.Service("db initializing")
	return
}

func (d *db) Close() { {{range $i,$k := .DBS}}{{if eq $k "` + string(Mongodb) + `"}}
	{{printf "d._%v.conn.(*mongo.Client).Disconnect(helper.Ctx)" $i}}{{else if eq $k "` + string(Postgres) + `"}}
	{{printf "d._%v.conn.(*sql.DB).Close()" $i}}{{end}}
	{{printf "helper.Log.Servicef(%v, _%v)" "\"db %q stoped\"" $i}}{{end}}
}
{{range $i,$k := .DBS}}{{if eq $k "` + string(Mongodb) + `"}}
{{printf "func (d *db) %v() *mongo.Database { return d._%v.conn.(*mongo.Client).Database(d._%v.dbName)" (title $i) $i $i}}}{{else if eq $k "` + string(Postgres) + `"}}
{{printf "func (d *db) %v() *sql.DB { return d._%v.conn.(*sql.DB)}" (title $i) $i}}{{end}}
{{printf "func (d *db) %vStore() %vStore.Store { return d._%v.store}" (title $i) $i $i}}{{end}}

func connPostgres(v *helper.DbConfig) (conn *sql.DB, err error) {
	conn, err = sql.Open("postgres", fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v sslmode=disable", v.User, v.Password, v.Host, v.Port, v.Name))
	if err != nil {
		return conn, fmt.Errorf("db not connected: %v", err)
	}
	if err = conn.PingContext(helper.Ctx); err != nil {
		return conn, fmt.Errorf("db not pinged: %v", err)
	}
	return
}

func connMongo(v *helper.DbConfig) (conn *mongo.Client, err error) {
	opt := options.Client().ApplyURI(fmt.Sprintf("mongodb://%v:%v", v.Host, v.Port)).SetAuth(options.Credential{AuthMechanism: "SCRAM-SHA-256", Username: v.User, Password: v.Password})
	conn, err = mongo.Connect(helper.Ctx, opt)
	if err != nil {
		return conn, fmt.Errorf("db not connected: %v", err)
	}
	if err = conn.Ping(helper.Ctx, nil); err != nil {
		return conn, fmt.Errorf("db not pinged: %v", err)
	}
	return
}`

	HelperFuncTmp = `package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/0LuigiCode0/library/logger"
)

func ParseConfig() (*Config, error) {
	_, err := os.Stat(ConfigDir + ConfigFiel)
	if err != nil {
		return nil, fmt.Errorf(KeyErrorNotFound+": file: %v", ConfigDir+ConfigFiel)
	}
	file, err := os.Open(ConfigDir + ConfigFiel)
	if err != nil {
		return nil, fmt.Errorf(KeyErrorOpen+": file: %v", err)
	}
	defer file.Close()
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf(KeyErrorRead+": body: %v", err)
	}
	data := new(Config)
	err = json.Unmarshal(buf, data)
	if err != nil {
		return nil, fmt.Errorf(KeyErrorParse+": json: %v", err)
	}

	return data, err
}

func InitCtx() {
	Ctx, CloseCtx = context.WithCancel(context.Background())
}
func InitLogger() {
	Log = logger.InitLogger("")
}`

	HelperModelTmp = `package helper

import (
	"context"
	"sync"

	"github.com/0LuigiCode0/library/logger"
)

var Ctx context.Context
var CloseCtx context.CancelFunc
var Log *logger.Logger
var Wg sync.WaitGroup

//Config модель конфига
type Config struct {
	Host     string                    ` + "`json:\"host\"`" + `
	Port     int                       ` + "`json:\"port\"`" + `
	DBS      map[string]*DbConfig      ` + "`json:\"dbs\"`" + `
	Handlers map[string]*HandlerConfig ` + "`json:\"handlers\"`" + `
}

type DbConfig struct {
	Type     DBType ` + "`json:\"type\"`" + `
	Host     string ` + "`json:\"host\"`" + `
	Port     int    ` + "`json:\"port\"`" + `
	Name     string ` + "`json:\"name\"`" + `
	User     string ` + "`json:\"user\"`" + `
	Password string ` + "`json:\"password\"`" + `
}

type HandlerConfig struct {
	Type     HandlerType ` + "`json:\"type\"`" + `
	Host     string      ` + "`json:\"host\"`" + `
	Port     int         ` + "`json:\"port\"`" + `
	User     string      ` + "`json:\"user\"`" + `
	Password string      ` + "`json:\"password\"`" + `
	IsTSL    bool        ` + "`json:\"is_tsl\"`" + `
}

//ErrLocal модель локализации ошибок
type ErrLocal struct {
	TitleEn string ` + "`bson:\"title_en\" json:\"title_en\"`" + `
	TitleRu string ` + "`bson:\"title_ru\" json:\"title_ru\"`" + `
}

//Модель для отправки по ws
type WsType string

const ()

type SendModel struct {
	Type WsType      ` + "`json:\"type\"`" + `
	Data interface{} ` + "`json:\"data\"`" + `
}

//ResponseError модель ошибки
type ResponseError struct {
	Code ErrCode ` + "`json:\"code\"`" + `
	Msg  string  ` + "`json:\"msg\"`" + `
}

//ResponseModel модель ответа
type ResponseModel struct {
	Success bool        ` + "`json:\"success\"`" + `
	Result  interface{} ` + "`json:\"result\"`" + `
}

type DBType string

const (
	Postgres DBType = "postgres"
	Mongodb  DBType = "mongodb"
)

type HandlerType string

const (
	TCP  HandlerType = "tcp"
	MQTT HandlerType = "mqtt"
	WS   HandlerType = "ws"
)

//Основные константы
const (
	UploadDir  = "./source/uploads/"
	ConfigDir  = "./source/configs/"
	ConfigFiel = "configServer.json"
	Secret     = "RB6SmP4h"
)

//Названеия коллекции

type Collection string

const ()

//Ключи контекста

type CtxKey int

const (
	CtxKeyValue CtxKey = iota
)

//Коды ошибок

type ErrCode byte

const (
	ErrorNotFound ErrCode = iota
	ErrorExist
	ErrorSave
	ErrorUpdate
	ErrorDelete
	ErorrAccessDenied
	ErrorInvalidParams
	ErrorParse
	ErrorWrite
	ErrorRead
	ErrorGenerate
	ErrorSend
)
const (
	KeyErrorNotFound      = "not found"
	KeyErrorExist         = "already exists"
	KeyErrorSave          = "create is failed"
	KeyErrorUpdate        = "update is failed"
	KeyErrorDelete        = "delete is failed"
	KeyErorrAccessDenied  = "access denied"
	KeyErrorInvalidParams = "invalid params"
	KeyErrorParse         = "parse is failed"
	KeyErrorOpen          = "open is failed"
	KeyErrorClose         = "close is failed"
	KeyErrorRead          = "read is failed"
	KeyErrorWrite         = "write is faleid"
	KeyErrorGenerate      = "generate is falied"
	KeyErrorSend          = "send is faied"
)

//Роли пользователей

type Role byte

const (
	RoleSuperAdmin Role = iota + 1
	RoleAdmin
	RoleUser
)

//Статусы

//Значения сортировки

//Типы полей`

	ConfigTmp = `{
	"host": "",
	"port": 8090,
	"dbs": { {{range $i,$k := .DBS}}{{if eq $k "` + string(Postgres) + `"}}
		"{{print $i}}": {
			"type": "postgres",
			"host": "0.0.0.0",
			"port": 5432,
			"name": "",
			"user": "postgres",
			"password": "1"
		}{{else if eq $k "` + string(Mongodb) + `"}}
		"{{print $i}}": {
			"type": "mongodb",
			"host": "0.0.0.0",
			"port": 27017,
			"name": "",
			"user": "mongodb",
			"password": "1"
		}{{end}},{{end}}
	},
	"handlers": { {{range $i,$k := .Handlers}}{{if eq $k "` + string(TCP) + `"}}
		"{{print $i}}": {
			"type": "tcp",
			"is_tsl": false
		}{{else if eq $k "` + string(MQTT) + `"}}
		"{{print $i}}": {
			"type": "mqtt",
			"host": "0.0.0.0",
			"port": 1883,
			"user": "",
			"password": "",
			"is_tsl": false
		}{{else if eq $k "` + string(WS) + `"}}
		"{{print $i}}": {
			"type": "ws",
			"password": "1",
			"is_tsl": false
		}{{end}},{{end}}
	}
}`
	DockerTmp = `FROM golang:{{print .GoVersion}} AS builder

RUN mkdir -p /build

WORKDIR /build
ADD . .
RUN cd core/cmd && go build

FROM ubuntu:18.04 AS {{print .ModuleName}}
WORKDIR /server
RUN apt update && apt install -y git ca-certificates && update-ca-certificates
COPY --from=builder /build/core/cmd/cmd .

CMD [ "/server/cmd" ]`
	ComposeBuildTmp = `version: "3"

services:
	{{print .ModuleName}}:
		image: localhost:5000/{{print .ModuleName}}:latest
		build:
			context: ./{{print .WorkDir}}
			dockerfile: ./Dockerfile`
	ComposeLocalTmp = `version: "3"

services: {{range $i,$k := .DBS}}{{if eq $k "` + string(Mongodb) + `"}}
	{{print $i}}:
		container_name: {{print $i}}
		image: mongo
		restart: always
		environment:
			MONGO_INITDB_DATABASE:
		volumes:
			- ./source/db:/data/db
			- ./source/dbres:/data/res
		networks:
			- net
		command: mongod --auth {{else if eq $k "` + string(Postgres) + `"}}
	{{print $i}}:
		container_name: {{print $i}}
		image: postgres
		restart: always
		environment:
			 POSTGRES_PASSWORD:
			 POSTGRES_USER:
			 POSTGRES_DB:
		ports:
			 - "5432:5432"
		volumes:
			 - ../source/db:/var/lib/postgresql/data
		networks:
			 - net {{end}}
		{{end}}
	server:
		container_name: {{print .ModuleName}}
		image: localhost:5000/{{print .ModuleName}}:latest
		restart: always
		expose: 8090
		networks:
			- net {{if gt (len .DBS) 0}}
		depends_on: {{range $i,$k := .DBS}}
			{{printf "- %v" $i}}{{end}}{{end}}
		volumes:
			- ./source/uploads:/server/source/uploads
			- ./source/configs:/server/source/configs

networks:
		net:
			driver: bridge`

	ModTmp = `module {{print .ModuleName}}

go {{print .GoVersion}}

require github.com/0LuigiCode0/library v0.0.11`

	StoreTmp = `package {{printf "%vStore" .}}

type Store interface {
}

type store struct{}

func InitStore() Store {
	return &store{}
}`

	HubTmp = `package hub{{$module := .ModuleName}}

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
	hubHelper.Helper {{if gt (len .DBS) 0}}
	database.DBForHandler{{end}}
	router  *mux.Router
	handler http.Handler
	{{range $i,$k := .Handlers}}
	{{printf "_%v %vHelper.Handler" $i $i}}{{end}}
}

func InitHub({{if gt (len .DBS) 0}}db database.DB, {{end}}conf *helper.Config) (H http.Handler, err error) {
	hh := &hub{ {{if gt (len .DBS) 0}}
		DBForHandler: db,{{end}}
		router:       mux.NewRouter(),
	}
	hh.SetHandler(hh.router)
	H = hh.handler
	{{range $i,$k := .Handlers}}
	if v, ok := conf.Handlers[{{printf "_%v" $i}}]; ok {
		hh.{{printf "_%v" $i}}, err = {{print $i}}Handler.InitHandler(hh, v)
		if err != nil {
			return nil, fmt.Errorf("handler not initializing: %v", err)
		}
		helper.Log.Servicef("handler %q initializing", {{printf "_%v" $i}})
	}{{end}}

	hh.Helper = hubHelper.InitHelper(hh)

	helper.Log.Service("handler initializing")
	return
}

func (h *hub) Router() *mux.Router        { return h.router }
func (h *hub) SetHandler(hh http.Handler) { h.handler = hh } {{range $i,$k := .Handlers}}
{{printf "func (h *hub) %v() %vHelper.Handler { return h._%v }" (title $i) $i $i}}{{end}}`

	HubHelperTmp = `package hubHelper{{$module := .ModuleName}}

import (
	"net/http" {{if gt (len .DBS) 0}}
	{{printf "\"%v/core/database\"" $module}}{{end}} {{range $i,$k := .Handlers}}
	{{printf "%vHelper \"%v/handlers/%v/helper\"" $i $module $i }}{{end}}
	{{if isOneTCP}}
	"github.com/gorilla/mux"{{end}}
)

type Helper interface {
}

type HelperForHandler interface { {{if gt (len .DBS) 0}}
	database.DBForHandler {{end}} {{if isOneTCP}}
	Router() *mux.Router
	SetHandler(hh http.Handler){{end}}
}

type HandlerForHelper interface { {{range $i,$k := .Handlers}}
	{{printf "%v() %vHelper.Handler" (title $i) $i}}{{end}}
}

type helper struct {
	HandlerForHelper
}

func InitHelper(H HandlerForHelper) Helper { return &helper{HandlerForHelper: H} }`

	HandlerTCPTmp = `package {{printf "%vHandler" (index . 0)}}

import (
	"encoding/json"
	"net/http"
	{{printf "%vHelper \"%v/handlers/%v/helper\"" (index . 0) (index . 1) (index . 0)}}
	{{printf "\"%v/helper\"" (index . 1)}}
	{{printf "hubHelper \"%v/hub/helper\"" (index . 1)}}
)

type handler struct {
	hubHelper.HelperForHandler
}

func InitHandler(hub hubHelper.HelperForHandler, conf *helper.HandlerConfig) (H {{print (index . 0)}}Helper.Handler, err error) {
	h := &handler{HelperForHandler: hub}

	h.Router().Use(h.middleware)
	h.SetHandler(applyCORS(h.Router()))
	return
}

func (h *handler) respOk(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	resp := &helper.ResponseModel{
		Success: true,
		Result:  data,
	}
	buf, err := json.Marshal(resp)
	if err != nil {
		helper.Log.Warningf(helper.KeyErrorParse+": josn: %v", err)
		h.respError(w, helper.ErrorParse, helper.KeyErrorParse+": josn")
		return
	}
	_, err = w.Write(buf)
	if err != nil {
		helper.Log.Warningf(helper.KeyErrorWrite+": response: %v", err)
		h.respError(w, helper.ErrorWrite, helper.KeyErrorWrite+": response")
		return
	}
}

func (h *handler) respError(w http.ResponseWriter, code helper.ErrCode, msg string) {
	w.Header().Set("Content-Type", "application/json")
	resp := &helper.ResponseModel{
		Success: false,
		Result: &helper.ResponseError{
			Code: code,
			Msg:  msg,
		},
	}
	buf, _ := json.Marshal(resp)
	_, err := w.Write(buf)
	if err != nil {
		helper.Log.Warningf(helper.KeyErrorWrite+": response: %v", err)
	}
}`
	HandlerMQTTTmp = `package {{printf "%vHandler" (index . 0)}}

import (
	"encoding/json"
	"fmt"
	"strings"
	{{printf "%vHelper \"%v/handlers/%v/helper\"" (index . 0) (index . 1) (index . 0)}}
	{{printf "\"%v/helper\"" (index . 1)}}
	{{printf "hubHelper \"%v/hub/helper\"" (index . 1)}}
	
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type handler struct {
	hubHelper.HelperForHandler

	client mqtt.Client
	err    chan error
}

func InitHandler(hub hubHelper.HelperForHandler, conf *helper.HandlerConfig) (H {{print (index . 0)}}Helper.Handler, err error) {
	h := &handler{HelperForHandler: hub}
	H = h

	client := mqtt.NewClient(mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%v:%v", conf.Host, conf.Port)).
		SetUsername(conf.User).
		SetPassword(conf.Password).
		SetOnConnectHandler(h.onConn).
		SetConnectionLostHandler(h.onLost).
		SetConnectRetry(true).
		SetResumeSubs(true))
	h.client = client
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	helper.Wg.Add(1)
	go h.loop()

	return
}

func (h *handler) haron(_ mqtt.Client, msg mqtt.Message) {
	var topic, route string
	str := strings.Split(msg.Topic(), "/")
	if len(str) == 3 {
		route = strings.TrimSpace(str[1])
		topic = strings.TrimSpace(str[2])
	} else {
		return
	}
	switch route {
	}
	fmt.Println(topic)
}

func (h *handler) onConn(client mqtt.Client) {
	helper.Log.Service("mqtt is conected")
	if token := client.Subscribe("back/#", 0, h.haron); token.Wait() && token.Error() != nil {
		h.err <- token.Error()
	}
}

func (h *handler) onLost(client mqtt.Client, err error) {
	helper.Log.Servicef("mqtt is disconected: %v", err)
	token := client.Unsubscribe("back/#")
	token.Wait()
}

func (h *handler) loop() {
	defer helper.Wg.Done()

	for {
		select {
		case e := <-h.err:
			helper.Log.Errorf("mqtt erorr: %v", e)
			h.client.Disconnect(1000)
			return
		case <-helper.Ctx.Done():
			helper.Log.Service("mqtt stoped")
			h.client.Disconnect(1000)
			return
		}
	}
}

func (h *handler) respOk(topic string, data interface{}) {
	resp := &helper.ResponseModel{
		Success: true,
		Result:  data,
	}
	buf, err := json.Marshal(resp)
	if err != nil {
		helper.Log.Warningf(helper.KeyErrorParse+": json: %v", err)
		return
	}
	if token := h.client.Publish(topic, 0, false, buf); token.Wait() && token.Error() != nil {
		helper.Log.Warningf(helper.KeyErrorSend+": mqtt: %v", token.Error())
	}
}

func (h *handler) respError(topic string, code helper.ErrCode, msg string) {
	resp := &helper.ResponseModel{
		Success: false,
		Result: &helper.ResponseError{
			Code: code,
			Msg:  msg,
		},
	}
	buf, _ := json.Marshal(resp)
	if token := h.client.Publish(topic, 0, false, buf); token.Wait() && token.Error() != nil {
		helper.Log.Warningf(helper.KeyErrorSend+": mqtt: %v", token.Error())
	}
}`
	HandlerWSTmp = `package  {{printf "%vHandler" (index . 0)}}

import (
	"net"
	"sync"
	{{printf "%vHelper \"%v/handlers/%v/helper\"" (index . 0) (index . 1) (index . 0)}}
	{{printf "\"%v/helper\"" (index . 1)}}
	{{printf "hubHelper \"%v/hub/helper\"" (index . 1)}}
)

type handler struct {
	hubHelper.HelperForHandler

	conn map[int64]map[string]net.Conn

	rw sync.Mutex
}

func InitHandler(hub hubHelper.HelperForHandler, conf *helper.HandlerConfig) (H {{print (index . 0)}}Helper.Handler, err error) {
	h := &handler{
		HelperForHandler: hub,
		conn:             map[int64]map[string]net.Conn{},
	}
	H = h

	helper.Wg.Add(1)
	go h.initTicker()

	h.Router().HandleFunc("/ws", h.WS)
	return
}

func (h *handler) addConn(userID int64, ipAddr string, connect net.Conn) {
	h.rw.Lock()
	defer h.rw.Unlock()
	if _, ok := h.conn[userID]; ok {
		h.conn[userID][ipAddr] = connect
	} else {
		h.conn[userID] = make(map[string]net.Conn)
		h.conn[userID][ipAddr] = connect
	}
}

func (h *handler) deleteConn(userID int64, ipAddr string) {
	h.rw.Lock()
	defer h.rw.Unlock()
	if _, ok := h.conn[userID]; ok {
		delete(h.conn[userID], ipAddr)

		if len(h.conn[userID]) == 0 {
			delete(h.conn, userID)
		}
	}
}

func (h *handler) GetConn(userID int64) map[string]net.Conn {
	h.rw.Lock()
	defer h.rw.Unlock()
	res := h.conn[userID]
	return res
}`

	MiddlewareTCPTmp = `package {{printf "%vHandler" (index . 0)}}

import (
	"net/http"

	"github.com/gorilla/handlers"
)

func (h *handler) middleware(conn http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn.ServeHTTP(w, r)
	})
}

func applyCORS(handler http.Handler) http.Handler {
	headersOk := handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "Access-Control-Allow-Origin"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	return handlers.CORS(headersOk, originsOk, methodsOk)(handler)
}`
	MiddleWareMQTTTmp = `package {{printf "%vHandler" (index . 0)}}`
	MiddleWareWSTmp   = `package {{printf "%vHandler" (index . 0)}}

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
			h.rw.Lock()
			for u, conns := range h.conn {
				for ip, conn := range conns {
					if err := wsutil.WriteServerText(conn, []byte{'1'}); err != nil {
						h.deleteConn(u, ip)
					}
				}
			}
			h.rw.Unlock()
		}
	}
}`

	HandlerHelperTmp = `package {{printf "%vHandler" (index . 0)}}
{{if eq (index . 1) "` + WS + `"}}
import "net"
{{end}}
type Handler interface { {{if eq (index . 1) "` + WS + `"}}
	GetConn(userID int64) map[string]net.Conn {{end}}
}`
)
