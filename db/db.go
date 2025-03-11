package db

import (
	"database/sql"
	"os"
	"path/filepath"
)

var dbConn *sql.DB

func GetDbConnection() (*sql.DB, error) {
	dbFilePath := getDbFilePath()
	if dbConn == nil {
		db, err := sql.Open("sqlite", dbFilePath)
		if err != nil {
			return nil, err
		}
		dbConn = db
	}
	return dbConn, nil
}

func getDbFilePath() string {
	dbFileName := os.Getenv("TODO_DBFILE")
	if dbFileName == "" {
		dbFileName = "scheduler.db"
	}
	return dbFileName
}

func CheckDb() error {
	dbFilePath := getDbFilePath()

	dbFile := filepath.Join(filepath.Dir("./"), dbFilePath)
	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	if install {
		dbConn, err = GetDbConnection()
		if err != nil {
			return err
		}
		_, err = dbConn.Exec("create table scheduler ( " +
			"id      integer constraint scheduler_pk primary key autoincrement unique, " +
			"date    char(8), " +
			"title   char(255), " +
			"comment text, " +
			"repeat  char(128)); " +
			"create index scheduler_date_index on scheduler (date);")
		if err != nil {
			return err
		}
	}
	return nil
}
