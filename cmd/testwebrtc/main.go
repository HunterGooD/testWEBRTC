package main

import (
	"bytes"
	"log"
	"os"

	"github.com/HunterGooD/testWEBRTC/internal/app"
	"github.com/HunterGooD/testWEBRTC/internal/db"
	"github.com/gobuffalo/packr/v2"
	"github.com/joho/godotenv"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Не назначен порт для запуска")
	}

	tDB := os.Getenv("TYPE_DB")
	if tDB == "" {
		loadEnv()
	}

	DB := db.New()
	application := app.New(":"+port, DB)

	application.Start()
}

func loadEnv() {
	envFile := packr.New("env", ".env")
	buff, err := envFile.Find(".env")
	if err != nil {
		log.Println(err)
	}

	if envVar, err := godotenv.Parse(bytes.NewReader(buff)); err != nil {
		log.Println(err)
		panic("Не удается считать переменные окружения")
	} else {
		for key, val := range envVar {
			os.Setenv(key, val)
		}
	}

}
