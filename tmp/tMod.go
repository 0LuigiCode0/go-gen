package tmp

const ModTmp = `module {{print .ModuleName}}

go {{print .GoVersion}}

require (
	github.com/0LuigiCode0/logger v0.0.1 {{if isOneMQTT}}
	github.com/eclipse/paho.mqtt.golang v1.3.4 // indirect {{end}} {{if isOneWS}}
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.0.4 // indirect {{end}} {{if isOneTCP}}
	github.com/gorilla/handlers v1.5.1 // indirect
	github.com/gorilla/mux v1.8.0 {{end}} {{if isOnePostgres}}
	github.com/lib/pq v1.10.1 {{end}} {{if isOneMongo}}
	go.mongodb.org/mongo-driver v1.5.1 {{end}}
)`
