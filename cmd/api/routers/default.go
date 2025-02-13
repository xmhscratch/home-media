package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping comment
func (route *RouteContext) Ping() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ginCtx.String(http.StatusOK, "")
	}
}

func (route *RouteContext) GetDefault() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		// json := make(map[string]interface{})
		// json["server_name"] = beego.AppConfig.String("streaming_domain_name")
		// json["server_node_name"] = beego.AppConfig.String("server_node_name")
		// json["host"] = beego.AppConfig.String("host")
		// c.Data["json"] = json
		// c.ServeJSON()
		ginCtx.String(http.StatusOK, "")
	}
}

func (ctx *RouteContext) Error(errCode int, msg string) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ginCtx.JSON(errCode, gin.H{
			"code":    errCode,
			"message": msg,
		})
	}
}

// NoRoute comment
func (route *RouteContext) NoRoute() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ginCtx.JSON(404, gin.H{
			"code":    "PAGE_NOT_FOUND",
			"message": "Page not found",
		})
	}
}
