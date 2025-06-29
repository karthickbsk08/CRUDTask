package connectdb

import (
	"database/sql"

	"gorm.io/gorm"
)

type DBConfig struct {
	Server   string
	Port     int
	User     string
	Password string
	Database string
	DBType   string
	DBflag   string
}

type AllDatabaseDetails struct {
	Mysql    DBConfig
	Mssql    DBConfig
	Postgres DBConfig
	DConnection_Pooling_Configuration
}

type DConnection_Pooling_Configuration struct {
	Max_Open_Conns     int
	Max_Idle_Conns     int
	Conn_Max_Idle_Time int
}

type DBInstance struct {
	GRMPostgres *gorm.DB
	GPostgres   *sql.DB
}
