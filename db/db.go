package db

import (
	"os"

	_ "github.com/ncruces/go-sqlite3/embed"
	"github.com/ncruces/go-sqlite3/gormlite"
	"gorm.io/gorm"

	"github.com/tinkerbaj/redditjobs/models"
	"github.com/tinkerbaj/redditjobs/utils"
)


//Global db for easier interaction with database
var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open(gormlite.Open("gorm.db"), &gorm.Config{})
	// DB, err = sql.Open("sqlite3", "file:demo.db?_pragma=foreign_keys(1)")
	if err != nil {
		utils.ColorError("Failed to open DB:", err)
		os.Exit(1)
	}
	err = DB.AutoMigrate(&models.Subreddit{}, &models.Flair{}, &models.Post{})
	if err != nil {
		utils.ColorError("Failed to migrate the tables:", err)

	}
}
