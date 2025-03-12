package wsevent

import (
	"home-media/sys"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/fiber/v2"
)

func HandleConnect(cfg *sys.Config, app *fiber.App) func(*socketio.EventPayload) {
	return func(ep *socketio.EventPayload) {}
}
