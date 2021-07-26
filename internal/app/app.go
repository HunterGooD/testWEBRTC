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

	router.POST("/api/signin", a.signin)
	middleware.POST("/api/logout", a.logout)
	middleware.GET("/api/rooms", a.getRooms)
	middleware.POST("/api/room/create", a.createRoom)
	middleware.POST("/api/room/:name/delete", a.deleteRoom)
	router.Any("/api/room/:name/join", a.joinRoom)

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
		go r.dispatchTimer()
		a.Rooms = append(a.Rooms, &r)
	}
	return nil
}
