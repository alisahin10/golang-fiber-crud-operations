package database

import (
	"github.com/tidwall/buntdb"
	"log"
	"sync"
)

var (
	DB   *buntdb.DB
	once sync.Once
)

// InitDB initializes the database.
func InitDB() {
	once.Do(func() {
		var err error
		DB, err = buntdb.Open("./buntdb.db") // DB file location.
		if err != nil {
			log.Fatal(err)
		}
	})
}

// CloseDB closes the database connection.
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
