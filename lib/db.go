package lib

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Used for it's side effect
)

var (
	dbEnv string

	// DB is used to communicate with the database
	DB *gorm.DB
)

func init() {
	var err error

	dbEnv = os.Getenv("WISHLIST_DB")
	if len(dbEnv) == 0 {
		log.Fatal("error: 'WISHLIST_DB' must be set")
	}

	DB, err = gorm.Open("postgres", dbEnv)
	if err != nil {
		log.Fatalf("error: Could not create or connect to the database\n\treason: %s", err)
	}
}
