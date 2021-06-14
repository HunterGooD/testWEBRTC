package app

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/HunterGooD/testWEBRTC/internal/db"
	"github.com/HunterGooD/testWEBRTC/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/pion/webrtc/v3"
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
	c.JSON(http.StatusOK, map[string]interface{}{
		"rooms": u.Rooms,
	})
}

func (a *App) createRoom(c *gin.Context) {

}

func (a *App) deleteRoom(c *gin.Context) {

}

func (a *App) joinRoom(ctx *gin.Context) {

	nameRoom := ctx.Param("name")

	var room *Room

	for _, roomIter := range a.Rooms {
		log.Print(roomIter.Name)
		if roomIter.Name == nameRoom {
			room = roomIter
			break
		}
	}

	if room == nil {
		ctx.JSON(http.StatusNotFound, map[string]string{
			"error": "Не найденна комната",
		})
		return
	}

	// Upgrade HTTP request to Websocket
	unsafeConn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	c := &threadSafeWriter{unsafeConn, sync.Mutex{}}

	// When this frame returns close the Websocket
	defer c.Close() //nolint

	// Create new PeerConnection
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		log.Print(err)
		return
	}

	// When this frame returns close the PeerConnection
	defer peerConnection.Close() //nolint

	// Accept one audio and one video track incoming
	for _, typ := range []webrtc.RTPCodecType{webrtc.RTPCodecTypeVideo, webrtc.RTPCodecTypeAudio} {
		if _, err := peerConnection.AddTransceiverFromKind(typ, webrtc.RTPTransceiverInit{
			Direction: webrtc.RTPTransceiverDirectionRecvonly,
		}); err != nil {
			log.Print(err)
			return
		}
	}

	// Add our new PeerConnection to global list
	room.mutex.Lock()
	room.peerConnections = append(room.peerConnections, peerConnectionState{peerConnection, c})
	room.mutex.Unlock()

	// Trickle ICE. Emit server candidate to client
	peerConnection.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			return
		}

		candidateString, err := json.Marshal(i.ToJSON())
		if err != nil {
			log.Println(err)
			return
		}

		if writeErr := c.WriteJSON(&websocketMessage{
			Event: "candidate",
			Data:  string(candidateString),
		}); writeErr != nil {
			log.Println(writeErr)
		}
	})

	// If PeerConnection is closed remove it from global list
	peerConnection.OnConnectionStateChange(func(p webrtc.PeerConnectionState) {
		switch p {
		case webrtc.PeerConnectionStateFailed:
			if err := peerConnection.Close(); err != nil {
				log.Print(err)
			}
		case webrtc.PeerConnectionStateClosed:
			room.signalPeerConnections()
		}
	})

	peerConnection.OnTrack(func(t *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		// Create a track to fan out our incoming video to all peers
		trackLocal := room.addTrack(t)
		defer room.removeTrack(trackLocal)

		buf := make([]byte, 1500)
		for {
			i, _, err := t.Read(buf)
			if err != nil {
				return
			}

			if _, err = trackLocal.Write(buf[:i]); err != nil {
				return
			}
		}
	})

	// Signal for the new PeerConnection
	room.signalPeerConnections()

	message := &websocketMessage{}
	for {
		_, raw, err := c.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		} else if err := json.Unmarshal(raw, &message); err != nil {
			log.Println(err)
			return
		}

		switch message.Event {
		case "candidate":
			candidate := webrtc.ICECandidateInit{}
			if err := json.Unmarshal([]byte(message.Data), &candidate); err != nil {
				log.Println(err)
				return
			}

			if err := peerConnection.AddICECandidate(candidate); err != nil {
				log.Println(err)
				return
			}
		case "answer":
			answer := webrtc.SessionDescription{}
			if err := json.Unmarshal([]byte(message.Data), &answer); err != nil {
				log.Println(err)
				return
			}

			if err := peerConnection.SetRemoteDescription(answer); err != nil {
				log.Println(err)
				return
			}
		}
	}
}
