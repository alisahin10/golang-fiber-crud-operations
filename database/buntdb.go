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

// InitDB initializes the BuntDB database with the given path.
func InitDB(dbPath string) {
	once.Do(func() {
		var err error
		DB, err = buntdb.Open(dbPath) // Use the provided database path.
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
