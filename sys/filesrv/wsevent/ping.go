package wsevent

import (
	"context"
	"encoding/json"
	"fmt"
	"home-media/sys"
	"home-media/sys/filesrv"
	"home-media/sys/session"
	"sync"
	"time"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"

	"github.com/redis/go-redis/v9"
)

func HandlePing(cfg *sys.Config, app *fiber.App) func(*socketio.EventPayload) {
	return func(ep *socketio.EventPayload) {
		var (
			err  error
			sock *SocketConnection
		)

		if s, ok := socketManager.Get(ep.Kws.UUID); !ok {
			ep.Kws.Fire(socketio.EventClose, []byte{})
			return
		} else {
			sock = s.(*sync.Pool).Get().(*SocketConnection)
		}

		go func() {
		stopMessage:
			for range time.Tick(time.Duration(10) * time.Millisecond) {
				if sock.Subscriber == nil {
					break stopMessage
				}

				if !ep.Kws.IsAlive() {
					break stopMessage
				}

				var (
					msgChan <-chan *redis.Message = sock.Subscriber.Channel()
					msg     *redis.Message        = <-msgChan
				)

				// fmt.Println(msg)

				if msg == nil {
					break stopMessage
				}

				if err := ep.Kws.Conn.WriteJSON(&filesrv.SocketMessage{
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
		}).SendToSocket(rds, sock.FileKey); err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(time.Duration(1) * time.Second)

		// fmt.Println(sys.BuildString(keyName, ":files"), fileKey)
		var (
			metaJSON string
			meta     *session.FileMetaInfo
			keyName  string = session.GetKeyName(sock.SessionId)
		)
		if metaJSON, err = rds.HGet(context.TODO(), sys.BuildString(keyName, ":files"), sock.FileKey).Result(); err != nil {
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
			}).SendToSocket(rds, sock.FileKey); err != nil {
				fmt.Println(err)
			}
		default:
			if err = (&session.FSMessage{
				Stage:   0,
				Message: "",
			}).SendToSocket(rds, sock.FileKey); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
