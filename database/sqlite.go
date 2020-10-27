package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"portScan_demo/forms"
)

var Db *gorm.DB

func InitDb() {
	db, err := gorm.Open("sqlite3", fmt.Sprintf("./database/database.db"))
	if err != nil {
		log.Fatalln(err)
	}

	Db = db
	Db.AutoMigrate(&forms.ScannerDb{})
}
