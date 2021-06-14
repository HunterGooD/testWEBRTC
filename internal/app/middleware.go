package app

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/HunterGooD/testWEBRTC/internal/db"
	"github.com/HunterGooD/testWEBRTC/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (a *App) middleware(c *gin.Context) {
	bearer := c.GetHeader("Authorization")
	authToken := strings.Split(bearer, " ")
	if len(authToken) != 2 {
		c.AbortWithStatusJSON(http.StatusForbidden, map[string]interface{}{
			"error": "Не верное значение Authorization",
		})
		return
	}
	tk, err := utils.VerifyToken(authToken[1])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.Set("userInfo", tk.UserId)
}

func (a *App) seedRandomDB() {
	var u1 db.User
	if err := a.DB.Model(&db.User{}).Where("id = ?", 1).First(&u1).Error; err != nil {
		if gorm.ErrRecordNotFound != err {
			log.Println("Произошла ошибка в полечении записи")
			return
		}
	}

	if u1.Login != "" {
		return
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
