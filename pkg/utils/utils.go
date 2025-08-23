package utils

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"

	_ "github.com/lib/pq"
)

var Connection *bun.DB

func ConnectToDB() error {
	var err error
	db, err := sql.Open("postgres", "user=admin dbname=mydb password=qwerty port=5432 host=postgres sslmode=disable")
	if err != nil {
		fmt.Println("Error while connecting to DB", err)
		return err
	}

	Connection = bun.NewDB(db, pgdialect.New())

	// Test the connection to the database
	err = Connection.Ping()
	if err != nil {
		log.Fatal("Error while checking connection", err)
		return err
	} else {
		log.Println("Successfully Connected")
	}
	return nil
}
