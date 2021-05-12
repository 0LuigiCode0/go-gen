package tmp

const HelperModelTmp = `package helper

import (
	"context"
	"sync"

	"github.com/0LuigiCode0/logger"
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
