package main

import (
	"home-media/cmd/file/middleware"
	"home-media/sys"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	// "github.com/gofiber/fiber/v2/middleware/filesystem"
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
		middleware.NewStorageHandler(cfg.DataPath, middleware.StorageConfig{
			Compress:      false,
			ByteRange:     true,
			Browse:        false,
			Download:      false,
			CacheDuration: 10 * time.Second,
			MaxAge:        3600,
		}),
	)

	app.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	middleware.AttachSocket(cfg, app)

	app.Get("*", func(c *fiber.Ctx) error {
		return c.SendString("File not found!")
	})

	go app.Listen(":4150")

	exit := make(chan struct{})
	SignalC := make(chan os.Signal, 4)

	signal.Notify(
		SignalC,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	go func() {
		for s := range SignalC {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				close(exit)
				return
			}
		}
	}()

	<-exit
	os.Exit(0)
}
