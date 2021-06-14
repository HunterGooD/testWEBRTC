package app

import (
	"net/http"

	"github.com/HunterGooD/testWEBRTC/internal/db"
	"github.com/HunterGooD/testWEBRTC/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (a *App) signin(c *gin.Context) {
	type ReqUserData struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	reqData := new(ReqUserData)
	c.BindJSON(reqData)
	var dbUser db.User

	if err := a.DB.Model(&db.User{}).Where("login = ?", reqData.Login).First(&dbUser).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Такой записи не существует",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "При получении данных произошла ошибка.",
		})
		return
	}

	if !utils.CheckPasswordHash(reqData.Password, dbUser.Password) {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "Не верный пароль.",
		})
		return
	}

	rBytes, _ := utils.RandomBytes(32)
	token := utils.CreateToken(dbUser.ID, rBytes)

	c.JSON(http.StatusOK, map[string]interface{}{
		"surname":  dbUser.Surname,
		"name":     dbUser.Name,
		"lastname": dbUser.Lastname,
		"token":    token,
	})
}

func (a *App) logout(c *gin.Context) {

}

func (a *App) getRooms(c *gin.Context) {
	userID, _ := c.Get("userInfo")
	a.DB.Where("user_id = ?", userID)
	var u db.User
	a.DB.Model(&db.User{}).Where("id = ?", userID).Preload("Rooms").Find(&u)
	c.JSON(http.StatusOK, u)
}

func (a *App) createRoom(c *gin.Context) {

}

func (a *App) joinRoom(c *gin.Context) {

}
