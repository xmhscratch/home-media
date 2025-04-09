package main

import (
	"net"
	"time"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	// "github.com/gofiber/fiber/v2/middleware/filesystem"

	"home-media/sys"
	"home-media/sys/filesrv"
	"home-media/sys/filesrv/wsevent"
)

func main() {
	cfg, err := sys.NewConfig("../")
	if err != nil {
		panic(err)
	}

	app := fiber.New()

	app.Use(cors.New())
	app.Use(
		"/:nodeId<[0-9a-z]{24}>/:fileKey<^[0-9a-z]+(\\.[0-9a-z]{1,20}|)(\\.vtt|\\.mp4)$>",
		filesrv.NewStorageHandler(cfg.DataPath, filesrv.StorageConfig{
			Compress:      false,
			ByteRange:     true,
			Browse:        false,
			Download:      false,
			CacheDuration: 10 * time.Second,
			MaxAge:        3600,
		}),
	)

	// socketio.On(socketio.EventConnect, wsevent.HandleConnect(cfg, app))
	// socketio.On(socketio.EventDisconnect, wsevent.HandleDisconnect(cfg, app))
	socketio.On(socketio.EventClose, wsevent.HandleClose(cfg, app))
	// socketio.On(socketio.EventError, wsevent.HandleError(cfg, app))
	socketio.On(socketio.EventPing, wsevent.HandlePing(cfg, app))

	app.Use("/ws/*", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/ws/:sessionId/:fileKey", socketio.New(wsevent.NewSocketRoute(cfg, app)))

	app.Get("*", func(c *fiber.Ctx) error {
		return c.SendString("File not found!")
	})

	_, port, err := net.SplitHostPort(cfg.EndPoint["file"])
	if err != nil {
		panic(err)
	}
	go app.Listen(net.JoinHostPort("0.0.0.0", port))

	sys.WaitTermination()
}
