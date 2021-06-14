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

	middleware := router.Group("/")

	middleware.Use(a.middleware)

	router.POST("/signin", a.signin)
	middleware.POST("/logout", a.logout)
	middleware.GET("/rooms", a.getRooms)
	middleware.POST("/room/create", a.createRoom)
	middleware.Any("/room/join", a.joinRoom)

	//TODO: потом удалить
	a.seedRandomDB()

	router.Run(a.addr)
}
