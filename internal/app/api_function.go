package app

import (
	"log"
	"strconv"

	"github.com/HunterGooD/testWEBRTC/internal/db"
	"github.com/HunterGooD/testWEBRTC/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (a *App) seedRandomDB() {
	var u1 db.User
	if err := a.DB.Model(&db.User{}).Where("id = ?", 1).First(&u1).Error; err != nil {
		if gorm.ErrRecordNotFound != err {
			log.Println("Произошла ошибка в полечении записи")
			return
		}
	}

	var rooms = make([]db.Room, 2)

	for i := 0; i < 2; i++ {
		var r db.Room
		r = db.Room{
			Name:     "room_" + strconv.Itoa(i),
			Password: "",
		}

		if err := a.DB.Create(&r).Error; err != nil {
			log.Printf("Не удается создать комнату %v", err)
		}
		rooms[i] = r
	}

	for i := 0; i < 5; i++ {
		var u db.User
		var r []db.Room
		pass, err := utils.HashPassword("test_" + strconv.Itoa(i))
		if err != nil {
			panic(err)
		}

		if i%2 == 0 {
			r = rooms
		}

		u = db.User{
			Login:    "test_" + strconv.Itoa(i),
			Password: pass,
			Avatar:   "https://image.flaticon.com/icons/png/512/18/18601.png",
			Rooms:    r,
		}

		if err := a.DB.Create(&u).Error; err != nil {
			log.Printf("Не удается создать пользователя %v", err)
		}
	}
}

func (a *App) signin(ctx *gin.Context) {
	type ReqUserData struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	reqData := new(ReqUserData)
	ctx.BindJSON(reqData)
}

func (a *App) logout(ctx *gin.Context) {

}
func (a *App) getRooms(ctx *gin.Context) {

}

func (a *App) createRoom(ctx *gin.Context) {

}

func (a *App) joinRoom(ctx *gin.Context) {

}
