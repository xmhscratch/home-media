package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func main() {
	app := fiber.New()

	// app.Use(cache.New(cache.Config{
	// 	Expiration:   30 * time.Minute,
	// 	CacheControl: true,
	// }))
	app.Use(cors.New())

	app.Use("/*", static.New("../public", static.Config{
		Compress:      false,
		ByteRange:     true,
		Browse:        false,
		Download:      false,
		CacheDuration: 10 * time.Second,
		MaxAge:        3600,
	}))

	app.Get("*", func(c fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	log.Fatal(app.Listen(":4150"))
}
