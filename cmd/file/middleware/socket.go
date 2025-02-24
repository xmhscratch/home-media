package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"home-media/sys"
	"home-media/sys/session"
	"time"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"

	"github.com/redis/go-redis/v9"
)

type SocketMessage struct {
	Event   int    `json:"event"`
	Payload string `json:"payload"`
}

func AttachSocket(cfg *sys.Config, app *fiber.App) {
	// var (
	// 	sockCtx context.Context = context.Background()
	// )

	// socketio.On(socketio.EventConnect, func(ep *socketio.EventPayload) { })

	// socketio.On(socketio.EventDisconnect, func(ep *socketio.EventPayload) {
	// 	// ep.Kws.Fire(socketio.EventClose, []byte{})
	// })

	socketio.On(socketio.EventClose, func(ep *socketio.EventPayload) {
		var (
			fileKey string = ep.Kws.GetStringAttribute("file_key")
		)

		if ep.Kws.GetAttribute("subscriber") != nil {
			var subscriber *redis.PubSub = ep.Kws.GetAttribute("subscriber").(*redis.PubSub)
			if err := subscriber.Unsubscribe(context.TODO(), fileKey); err != nil {
			}
			if err := subscriber.Close(); err != nil {
			}
			ep.Kws.SetAttribute("subscriber", nil)
		}
		if ep.Kws.GetAttribute("redis") != nil {
			var rds *redis.Client = ep.Kws.GetAttribute("redis").(*redis.Client)
			rds.Close()
			ep.Kws.SetAttribute("rds", nil)
		}
		ep.Kws.Conn.Close()

		// errMsg := fmt.Sprintf("Close event %s: %s", fileKey, ep.Error)
		// fmt.Println(errMsg)
	})

	// socketio.On(socketio.EventError, func(ep *socketio.EventPayload) {
	// 	var fileKey string = ep.Kws.GetStringAttribute("file_key")
	// 	errMsg := fmt.Sprintf("Error event %s: %s", fileKey, ep.Error)
	// 	fmt.Println(errMsg)
	// })

	socketio.On(socketio.EventPing, func(ep *socketio.EventPayload) {
		var (
			err        error
			fileKey    string = ep.Kws.GetStringAttribute("file_key")
			sessionId  string = ep.Kws.GetStringAttribute("session_id")
			keyName    string = session.GetKeyName(sessionId)
			rds        *redis.Client
			subscriber *redis.PubSub
		)

		rds = sys.NewClient(cfg)
		ep.Kws.SetAttribute("rds", rds)

		if subscriber == nil {
			subscriber = rds.Subscribe(context.TODO(), fileKey)
			ep.Kws.SetAttribute("subscriber", subscriber)
		}

		go func() {
		stopMessage:
			for {
				time.Sleep(time.Duration(50) * time.Millisecond)

				if subscriber == nil || ep.Kws.GetAttribute("subscriber") == nil {
					break stopMessage
				}

				if !ep.Kws.IsAlive() {
					break stopMessage
				}

				var msgChan <-chan *redis.Message = subscriber.Channel()
				var msg *redis.Message = <-msgChan

				// fmt.Println(msg)

				if msg == nil {
					break stopMessage
				}

				if err := ep.Kws.Conn.WriteJSON(&SocketMessage{
					Event:   websocket.TextMessage,
					Payload: msg.Payload,
				}); err != nil {
					fmt.Println(err)
					break stopMessage
					// ep.Kws.Fire(socketio.EventError, []byte(err.Error()))
				}
				// fmt.Println(ep.Kws.UUID, ep.Kws.IsAlive())
			}
		}()

		if err = (&session.FSMessage{
			Stage:   1,
			Message: "Initializing...",
		}).SendToSocket(rds, fileKey); err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(time.Duration(1) * time.Second)

		// fmt.Println(sys.BuildString(keyName, ":files"), fileKey)
		var (
			metaJSON string
			meta     *session.FileMetaInfo
		)
		if metaJSON, err = rds.HGet(context.TODO(), sys.BuildString(keyName, ":files"), fileKey).Result(); err != nil {
			fmt.Println(err)
			return
		} else {
			if err = json.Unmarshal([]byte(metaJSON), &meta); err != nil {
				fmt.Println(err)
				return
			}
		}

		// fmt.Println(metaJSON)
		switch meta.SourceReady {
		case int(1):
			if err = (&session.FSMessage{
				Stage:   5,
				Message: metaJSON,
			}).SendToSocket(rds, fileKey); err != nil {
				fmt.Println(err)
			}
		default:
			if err = (&session.FSMessage{
				Stage:   0,
				Message: "",
			}).SendToSocket(rds, fileKey); err != nil {
				fmt.Println(err)
				return
			}
		}
	})

	app.Get("/ws/:sessionId/:fileKey", socketio.New(func(kws *socketio.Websocket) {
		var (
			// err       error
			sessionId string              = kws.Params("sessionId")
			fileKey   string              = kws.Params("fileKey")
			msgChan   chan *SocketMessage = make(chan *SocketMessage)
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
					msg *SocketMessage
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

			var msg *SocketMessage
			if msg = <-msgChan; msg == nil {
				msg.Event = websocket.CloseMessage
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
	}))
}
