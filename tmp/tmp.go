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
