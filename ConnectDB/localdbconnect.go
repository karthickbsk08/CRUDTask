package connectdb

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	POSTGRES = "postgres"
)

func LocalDBConnect(pDbType string) (lSqlDb *sql.DB, lGORMDb *gorm.DB, lErr error) {
	var Dbdetails AllDatabaseDetails
	var dsn string
	var lDialector gorm.Dialector
	Dbdetails.Init()

	switch pDbType {
	case "postgres":
		var dbConfig DBConfig
		dbConfig = Dbdetails.Postgres
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbConfig.Server, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Database)
		lDialector = postgres.New(postgres.Config{
			DSN: dsn})
	default:
		lErr = fmt.Errorf("unsupported database type: %s", pDbType)
		return lSqlDb, lGORMDb, lErr
	}

	lGORMDb, lErr = gorm.Open(lDialector, &gorm.Config{})
	if lErr != nil {
		log.Println("Error while gorm db connection " + lErr.Error())
		return lSqlDb, lGORMDb, lErr
	}

	lSqlDb, lErr = lGORMDb.DB()
	if lErr != nil {
		log.Println("Error while extracting db connection from gorm connection " + lErr.Error())
		return lSqlDb, lGORMDb, lErr
	}

	lSqlDb.SetMaxOpenConns(Dbdetails.Max_Open_Conns)
	lSqlDb.SetMaxIdleConns(Dbdetails.Max_Idle_Conns)
	lSqlDb.SetConnMaxIdleTime(time.Duration(Dbdetails.Conn_Max_Idle_Time))
	return
}
