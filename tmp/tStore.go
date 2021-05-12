package tmp

const StoreTmp = `package {{printf "%vStore" (index . 0)}}
{{if eq (index . 1) "` + string(Mongodb) + `"}}
import "go.mongodb.org/mongo-driver/mongo"{{else if eq (index . 1) "` + string(Postgres) + `"}}
import "database/sql"{{end}}

type Store interface { {{if eq (index . 1) "` + string(Postgres) + `"}}
	Begin() (*sql.Tx, error){{end}}
}

type store struct{ {{if eq (index . 1) "` + string(Postgres) + `"}}
	db *sql.DB{{end}}
}

func InitStore({{if eq (index . 1) "` + string(Mongodb) + `"}}db *mongo.Database{{else if eq (index . 1) "` + string(Postgres) + `"}}db *sql.DB{{end}}) Store {
	return &store{ {{if eq (index . 1) "` + string(Postgres) + `"}}
		db: db,{{end}}
	}
}
{{if eq (index . 1) "` + string(Postgres) + `"}}
func (s *store) Begin() (*sql.Tx, error) {
	return s.db.Begin()
}{{end}}`
