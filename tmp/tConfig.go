package tmp

const ConfigTmp = `{
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
