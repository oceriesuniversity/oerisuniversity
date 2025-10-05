package dbcode

import (
	"database/sql"
	"log"
)

type DBType struct {
	DB *sql.DB
}

func SqlRead() DBType {
	db, err := sql.Open("sqlite3", "./data/ucms.db")

	if err != nil {
		log.Fatal(err)
	}

	Sql := DBType{
		DB: db,
	}
	return Sql
}
