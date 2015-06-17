package main

import (
	"log"
	"os"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

const (
	DBFile = "ps2bot-sql.db"
)

var (
	db *sql.DB
)

func init() {
	tmpDB, err := sql.Open("sqlite3", DBFile)
	if err != nil {
		log.Printf("Error opening database: %v", err)
		os.Exit(1)
	}
	db = tmpDB

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		os.Exit(1)
	}

	tx.Exec("CREATE TABLE IF NOT EXISTS oldmentions(id TEXT)")
	tx.Exec("CREATE INDEX IF NOT EXISTS mentionindex on oldmentions(id)")

	err = tx.Commit()
	if err != nil {
		log.Printf("Failed to execute transaction: %v", err)
		os.Exit(1)
	}
}
