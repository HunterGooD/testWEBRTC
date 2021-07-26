package app

import (
	"encoding/json"
	"time"

	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"
)

// Вызывается при создании комнаты в виде горутины
func (r *Room) dispatchTimer() {
	for range time.NewTicker(time.Second * 3).C {
		r.dispatchKeyFrame()
	}
}

func (r *Room) dispatchKeyFrame() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for i := range r.peerConnections {
		for _, receiver := range r.peerConnections[i].peerConnection.GetReceivers() {
			if receiver.Track() == nil {
				continue
			}

			_ = r.peerConnections[i].peerConnection.WriteRTCP([]rtcp.Packet{
				&rtcp.PictureLossIndication{
					MediaSSRC: uint32(receiver.Track().SSRC()),
				},
			})
		}
	}
}

// Добавляет в список новый track и рассылает остальным соединениям
func (r *Room) addTrack(t *webrtc.TrackRemote) *webrtc.TrackLocalStaticRTP {
	r.mutex.Lock()
	defer func() {
		r.mutex.Unlock()
		r.signalPeerConnections()
	}()

	// Создает новый трек с кодеком
	trackLocal, err := webrtc.NewTrackLocalStaticRTP(t.Codec().RTPCodecCapability, t.ID(), t.StreamID())
	if err != nil {
		panic(err)
	}

	r.trackLocals[t.ID()] = trackLocal
	return trackLocal

}

// Удаляет соединение и уведомляет всех
func (r *Room) removeTrack(t *webrtc.TrackLocalStaticRTP) {
	r.mutex.Lock()
	defer func() {
		r.mutex.Unlock()
		r.signalPeerConnections()
	}()

	delete(r.trackLocals, t.ID())
}

// Обновляется с каждым подключением и обновляет все треки
func (r *Room) signalPeerConnections() {
	r.mutex.Lock()
	defer func() {
		r.mutex.Unlock()
		r.dispatchKeyFrame()
	}()

	attemptSync := func() (tryAgain bool) {
		for i := range r.peerConnections {
			if r.peerConnections[i].peerConnection.ConnectionState() == webrtc.PeerConnectionStateClosed {
				r.peerConnections = append(r.peerConnections[:i], r.peerConnections[i+1:]...)
				return true // We modified the slice, start from the beginning
			}

			// map of sender we already are seanding, so we don't double send
			existingSenders := map[string]bool{}

			for _, sender := range r.peerConnections[i].peerConnection.GetSenders() {
				if sender.Track() == nil {
					continue
				}

				existingSenders[sender.Track().ID()] = true

				// If we have a RTPSender that doesn't map to a existing track remove and signal
				if _, ok := r.trackLocals[sender.Track().ID()]; !ok {
					if err := r.peerConnections[i].peerConnection.RemoveTrack(sender); err != nil {
						return true
					}
				}
			}

			// Don't receive videos we are sending, make sure we don't have loopback
			for _, receiver := range r.peerConnections[i].peerConnection.GetReceivers() {
				if receiver.Track() == nil {
					continue
				}

				existingSenders[receiver.Track().ID()] = true
			}

			// Add all track we aren't sending yet to the PeerConnection
			for trackID := range r.trackLocals {
				if _, ok := existingSenders[trackID]; !ok {
					if _, err := r.peerConnections[i].peerConnection.AddTrack(r.trackLocals[trackID]); err != nil {
						return true
					}
				}
			}

			offer, err := r.peerConnections[i].peerConnection.CreateOffer(nil)
			if err != nil {
				return true
			}

			if err = r.peerConnections[i].peerConnection.SetLocalDescription(offer); err != nil {
				return true
			}

			offerString, err := json.Marshal(offer)
			if err != nil {
				return true
			}

			if err = r.peerConnections[i].websocket.WriteJSON(&websocketMessage{
				Event: "offer",
				Data:  string(offerString),
			}); err != nil {
				return true
			}
		}

		return
	}

	for syncAttempt := 0; ; syncAttempt++ {
		if syncAttempt == 25 {
			// Release the lock and attempt a sync in 3 seconds. We might be blocking a RemoveTrack or AddTrack
			go func() {
				time.Sleep(time.Second * 3)
				r.signalPeerConnections()
			}()
			return
		}

		if !attemptSync() {
			break
		}
	}
}
