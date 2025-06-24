package wsevent

import (
	"context"
	"encoding/json"
	"fmt"
	"home-media/sys"
	"home-media/sys/filesrv"
	"sync"
	"time"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/groupcache/lru"
	"github.com/redis/go-redis/v9"
)

var rds *redis.Client
var socketManager *lru.Cache

func NewSocketConnection(cfg *sys.Config, sessionId string, fileKey string) *sync.Pool {
	if rds == nil {
		rds = sys.NewClient(cfg)
	}

	var sockUUID string = sys.GenerateV5(sessionId, fileKey, WSEVENT_NAMESPACE)
	ctx := &sync.Pool{
		New: func() interface{} {
			return &SocketConnection{
				UUID:       sockUUID,
				SessionId:  sessionId,
				FileKey:    fileKey,
				Subscriber: rds.Subscribe(context.TODO(), fileKey),
			}
		},
	}

	if socketManager == nil {
		socketManager = lru.New(10)
		socketManager.OnEvicted = func(key lru.Key, value interface{}) {
			// sockUUID := key.(string)
			sock := (value.(*sync.Pool)).Get().(*SocketConnection)
			// fmt.Println(sockUUID, sock)
			if err := sock.Subscriber.Unsubscribe(context.TODO(), fileKey); err != nil {
			}
			if err := sock.Subscriber.Close(); err != nil {
			}
			(value.(*sync.Pool)).Put(sock)
		}
	}

	socketManager.Add(sockUUID, ctx)
	return ctx
}

func NewSocketRoute(cfg *sys.Config, app *fiber.App) func(*socketio.Websocket) {
	return func(kws *socketio.Websocket) {
		var (
			err       error
			sessionId string = kws.Params("sessionId")
			fileKey   string = kws.Params("fileKey")
		)

		kws.SetAttribute("file_key", fileKey)
		kws.SetAttribute("session_id", sessionId)

		soc := NewSocketConnection(cfg, sessionId, fileKey)
		kws.UUID = soc.Get().(*SocketConnection).UUID

	checkNewMessage:
		for range time.Tick(time.Duration(10) * time.Millisecond) {
			var msgChan chan *filesrv.SocketMessage = make(chan *filesrv.SocketMessage)
			defer close(msgChan)

			go func() {
				var (
					tp int
					pl []byte
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

				{
					var msg *filesrv.SocketMessage
					if err = json.Unmarshal(pl, &msg); err != nil {
						// kws.Fire(socketio.EventError, []byte{})
						fmt.Println(err.Error())
						// msgChan <- nil
						return
					}
					msgChan <- msg
				}
			}()

		readMsg:
			select {
			case msg := <-msgChan:
				if msg == nil {
					msg = &filesrv.SocketMessage{
						Event: websocket.CloseMessage,
					}
				}

				switch msg.Event {
				case websocket.PingMessage:
					kws.Fire(socketio.EventPing, []byte{})
				case websocket.TextMessage:
					kws.Fire(socketio.EventMessage, []byte{})
				case websocket.CloseMessage:
					kws.Fire(socketio.EventClose, []byte{})
					break checkNewMessage
				// kws.Fire(socketio.EventPing, []byte{})
				// case websocket.BinaryMessage:
				// 	if err := kws.Conn.WriteJSON(string(p)); err != nil {
				// 		kws.Fire(socketio.EventError, []byte(err.Error()))
				// 	}
				default:
					kws.Fire(socketio.EventClose, []byte{})
					break checkNewMessage
				}

				break readMsg
			case <-time.After(TIMEOUT_READ_MESSAGE):
				fmt.Println("read message timeout")
				break readMsg
			}
		}
	}
}
