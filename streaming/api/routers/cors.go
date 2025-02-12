package routers

import (
	"github.com/gin-gonic/gin"
)

// CORS comment
func (route *RouteContext) CORS() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ginCtx.Header("Access-Control-Allow-Origin", "*")
		ginCtx.Header("Access-Control-Allow-Headers", "*")
		ginCtx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")

		ginCtx.Next()
	}
}
