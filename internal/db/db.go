package db

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New() *gorm.DB {
	var DB *gorm.DB
	urlDB := os.Getenv("DATABASE_URL")
	if urlDB == "" {
		panic("Error url for database  not found")
	}

	spDB := strings.Split(urlDB, "://")[1]
	infU := strings.Split(spDB, ":")
	user := infU[0]
	pass := strings.Split(infU[1], "@")[0]
	infH := strings.Split(urlDB, "@")[1]
	infoHost := strings.Split(infH, ":")
	host := infoHost[0]
	port := strings.Split(infoHost[1], "/")[0]
	dbName := strings.Split(infoHost[1], "/")[1]

	var dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Europe/Samara",
		host,
		user,
		pass,
		dbName,
		port,
	)

	dbType := os.Getenv("TYPE_DB")
	if dbType == "postgres" {
		if db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
			log.Printf("%v", err)
			panic(err)
		} else {
			DB = db
		}
	} else if dbType == "sqlite" {
		if db, err := gorm.Open(sqlite.Open("db/test.db"), &gorm.Config{}); err != nil {
			panic(err)
		} else {
			DB = db
		}
	}

	if err := DB.AutoMigrate(&User{}, &Room{}); err != nil {
		log.Print("Не возможно создать базу")
		log.Fatal(err)
	}
	return DB
}
