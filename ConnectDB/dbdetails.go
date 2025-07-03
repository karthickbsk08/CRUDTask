package connectdb

import (
	"log"
	"strconv"
	"tasks/common"
	"tasks/tomlutil"
)

func (db *AllDatabaseDetails) Init() {
	log.Println("AllDatabaseDetails.Init(+)")

	lDbconfig := tomlutil.ReadTomlConfig("toml/dbconfig.toml")

	//DB Connection pool configuration
	db.Max_Open_Conns, _ = strconv.Atoi(tomlutil.GetKeyVal(lDbconfig, "Max_Open_Conns"))
	db.Max_Idle_Conns, _ = strconv.Atoi(tomlutil.GetKeyVal(lDbconfig, "Max_Idle_Conns"))
	db.Conn_Max_Idle_Time, _ = strconv.Atoi(tomlutil.GetKeyVal(lDbconfig, "Conn_Max_Idle_Time"))

	db.Postgres.Server = tomlutil.GetKeyVal(lDbconfig, "Db_Server")
	db.Postgres.Port, _ = strconv.Atoi(tomlutil.GetKeyVal(lDbconfig, "Db_Port"))
	db.Postgres.User = tomlutil.GetKeyVal(lDbconfig, "Db_User")
	db.Postgres.Password = tomlutil.GetKeyVal(lDbconfig, "Db_Password")
	db.Postgres.Database = tomlutil.GetKeyVal(lDbconfig, "Db_Database")
	db.Postgres.DBType = tomlutil.GetKeyVal(lDbconfig, "DBType")

	log.Println("db : ", common.DoMarshall(db))

	log.Println("AllDatabaseDetails.Init(-)")

}
