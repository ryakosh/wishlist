package lib

import (
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
		LogError(LFatal, "'WISHLIST_DB' must be set", err)
	}

	DB, err = gorm.Open("postgres", dbEnv)
	if err != nil {
		LogError(LFatal, "Could not create or connect to the database", err)
	}
}
