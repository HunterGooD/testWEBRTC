package app

import (
	"log"
	"sync"

	"github.com/HunterGooD/testWEBRTC/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/pion/webrtc/v3"
	"gorm.io/gorm"
)

type App struct {
	addr  string
	Rooms []*Room
	DB    *gorm.DB
}

func New(addr string, DB *gorm.DB) *App {
	return &App{
		addr:  addr,
		DB:    DB,
		Rooms: make([]*Room, 0),
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
	middleware.POST("/room/:name/delete", a.deleteRoom)
	router.Any("/room/:name/join", a.joinRoom)

	//TODO: потом удалить
	a.seedRandomDB()

	if err := a.StartExistsRooms(); err != nil {
		log.Println("Ошибка при запуске комнат")
	}

	router.Run(a.addr)
}

func (a *App) StartExistsRooms() error {
	rows, err := a.DB.Model(&db.Room{}).Rows()
	if err != nil {
		return err
	}
	for rows.Next() {
		var room db.Room
		if err := a.DB.ScanRows(rows, &room); err != nil {
			return err
		}
		var r Room
		r = Room{
			Name:            room.Name,
			mutex:           &sync.Mutex{},
			peerConnections: []peerConnectionState{},
			trackLocals:     map[string]*webrtc.TrackLocalStaticRTP{},
		}
		a.Rooms = append(a.Rooms, &r)
	}
	return nil
}
