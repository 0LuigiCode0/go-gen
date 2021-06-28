package tmp

const DatabaseTmp = `package database{{$module := .ModuleName}}

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
			store:  {{print $i}}Store.InitStore({{if eq $k "` + string(Mongodb) + `"}}conn.(*mongo.Client).Database(v.Name){{else if eq $k "` + string(Postgres) + `"}}conn.(*sql.DB){{end}}),
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
{{printf "func (d *db) %vStore() %vStore.Store { return d._%v.store.(%vStore.Store)}" (title $i) $i $i $i}}{{end}}
{{if isOnePostgres}}
func connPostgres(v *helper.DbConfig) (conn *sql.DB, err error) {
	conn, err = sql.Open("postgres", fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v sslmode=disable", v.User, v.Password, v.Host, v.Port, v.Name))
	if err != nil {
		return conn, fmt.Errorf("db not connected: %v", err)
	}
	if err = conn.PingContext(helper.Ctx); err != nil {
		return conn, fmt.Errorf("db not pinged: %v", err)
	}
	return
}{{end}}
{{if isOneMongo}}
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
}{{end}}`
