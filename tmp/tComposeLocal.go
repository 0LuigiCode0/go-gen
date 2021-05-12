package tmp

const ComposeLocalTmp = `version: "3"

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
