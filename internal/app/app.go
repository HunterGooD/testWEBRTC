package app

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	addr  string
	Rooms []Room
	DB    *gorm.DB
}

func New(addr string, DB *gorm.DB) *App {
	return &App{
		addr: addr,
		DB:   DB,
	}
}

func (a *App) Start() {
	router := gin.Default()

	router.POST("/signin", a.signin)
	router.POST("/logout", a.logout)
	router.GET("/rooms", a.getRooms)
	router.POST("/room/create", a.createRoom)
	router.Any("/room/join", a.joinRoom)

	a.seedRandomDB()

	router.Run()
}
