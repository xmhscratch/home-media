package wsevent

import (
	"context"
	"home-media/sys"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/fiber/v2"

	"github.com/redis/go-redis/v9"
)

func HandleClose(cfg *sys.Config, app *fiber.App) func(*socketio.EventPayload) {
	return func(ep *socketio.EventPayload) {
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
	}
}
