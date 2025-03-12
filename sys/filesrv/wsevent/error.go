package wsevent

import (
	"home-media/sys"

	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/fiber/v2"
)

func HandleError(cfg *sys.Config, app *fiber.App) func(*socketio.EventPayload) {
	return func(ep *socketio.EventPayload) {
		// 	var fileKey string = ep.Kws.GetStringAttribute("file_key")
		// 	errMsg := fmt.Sprintf("Error event %s: %s", fileKey, ep.Error)
		// 	fmt.Println(errMsg)
	}
}
