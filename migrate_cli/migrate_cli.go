package main

import (
	"log"

	"github.com/yneee/exsongs/models"
	"github.com/yneee/exsongs/settings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func main() {
	db, err := gorm.Open("mysql", settings.DbInfo)
	if err != nil {
		log.Fatalln(err)
	}

	db.Set("gorm:table_options", "ENGINE=InnoDB")

	db.AutoMigrate(&models.Log{})
	db.AutoMigrate(&models.Song{})
	db.AutoMigrate(&models.FailedPage{})
}
