package wsevent

import (
	"fmt"
	"home-media/sys"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/fiber/v2"
)

func HandleError(cfg *sys.Config, app *fiber.App) func(*socketio.EventPayload) {
	return func(ep *socketio.EventPayload) {
		fmt.Println(ep.Error)
	}
}
