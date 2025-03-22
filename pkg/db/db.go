package db

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

func connect(database_name string) *sql.DB {

    cfg := mysql.Config {

        User:   "root",
        Passwd: "root",
        Net:    "tcp",
        Addr:   "localhost:3306",
        DBName: database_name,
        MultiStatements: true,

    }

    db, err := sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        log.Fatal(err)
    }

    pingErr := db.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
    }

    log.Println("Succesfully connected to database: " + database_name)

    return db

}
