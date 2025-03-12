package wsevent

import (
	"encoding/json"
	"fmt"
	"home-media/sys"
	"home-media/sys/filesrv"
	"time"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func NewSocketRoute(cfg *sys.Config, app *fiber.App) func(*socketio.Websocket) {
	return func(kws *socketio.Websocket) {
		var (
			// err       error
			sessionId string                      = kws.Params("sessionId")
			fileKey   string                      = kws.Params("fileKey")
			msgChan   chan *filesrv.SocketMessage = make(chan *filesrv.SocketMessage)
		)

		kws.UUID = sys.GenerateV5(sessionId, fileKey, "3df0bd87e067452aa952d3915e319461")

		kws.SetAttribute("file_key", fileKey)
		kws.SetAttribute("session_id", sessionId)

	stopMessage:
		for {
			time.Sleep(time.Duration(50) * time.Millisecond)

			go func() {
				var (
					err error
					tp  int
					pl  []byte
					msg *filesrv.SocketMessage
				)

				if tp, pl, err = kws.Conn.ReadMessage(); err != nil {
					// kws.Fire(socketio.EventError, []byte(err.Error()))
					fmt.Println(err.Error())
					// msgChan <- nil
					return
				}

				if tp != 1 {
					// msgChan <- nil
					return
				}

				if err = json.Unmarshal(pl, &msg); err != nil {
					// kws.Fire(socketio.EventError, []byte{})
					fmt.Println(err.Error())
					// msgChan <- nil
					return
				}
				msgChan <- msg
			}()

			var msg *filesrv.SocketMessage
			if msg = <-msgChan; msg == nil {
				msg = &filesrv.SocketMessage{
					Event: websocket.CloseMessage,
				}
			}

			// fmt.Println(msg)

			switch msg.Event {
			case websocket.PingMessage:
				kws.Fire(socketio.EventPing, []byte{})
			case websocket.TextMessage:
				kws.Fire(socketio.EventMessage, []byte{})
			case websocket.CloseMessage:
				kws.Fire(socketio.EventClose, []byte{})
				break stopMessage
			// kws.Fire(socketio.EventPing, []byte{})
			// case websocket.BinaryMessage:
			// 	if err := kws.Conn.WriteJSON(string(p)); err != nil {
			// 		kws.Fire(socketio.EventError, []byte(err.Error()))
			// 	}
			default:
				kws.Fire(socketio.EventClose, []byte{})
				break stopMessage
			}
		}
	}
}
