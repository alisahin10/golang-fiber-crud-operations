package database

import (
	"log"
	"sync"

	"github.com/tidwall/buntdb"
)

var (
	DB   *buntdb.DB
	once sync.Once // Opens the DB only once.
)

// InitDB starts the db provide connection.
func InitDB() {
	once.Do(func() {
		var err error
		DB, err = buntdb.Open("./buntdb.db") // DB file location.
		if err != nil {
			log.Fatal(err)
		}
	})
}

// CloseDB closes the DB connection.
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
