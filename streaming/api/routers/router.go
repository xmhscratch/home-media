package routers

import (
	"fmt"
	"net/http"
	"time"

	"home-media/streaming/core"

	"github.com/gin-gonic/gin"
	expirable "github.com/hashicorp/golang-lru/v2/expirable"
)

// NewRoute comment
func NewRoute(cfg *core.Config) (ctx *RouteContext, err error) {
	sessionKeyVault := expirable.NewLRU[string, string](5000, nil, 15*60*time.Second)
	ctx = &RouteContext{
		Config:          cfg,
		SessionKeyVault: sessionKeyVault,
	}
	return ctx, err
}

// Init comment
func (ctx *RouteContext) Init(router *gin.Engine) {
	router.Use(gin.Logger())
	router.Use(gin.CustomRecovery(func(ginCtx *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			ginCtx.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		ginCtx.AbortWithStatus(http.StatusInternalServerError)
	}))

	router.Use(ctx.CORS())
	router.OPTIONS("/*all", ctx.Ping())

	router.PUT("/*filePath", ctx.CreateSession())
	// router.POST("/:ssid/*filePath", ctx.Download())
	router.HEAD("/:ssid/*filePath", ctx.CheckProgress())
	router.GET("/:ssid/files", ctx.GetFiles())
	router.GET("/:ssid/1/*filePath", ctx.DownloadDirect())
	router.GET("/:ssid/2/*filePath", ctx.DownloadTorrent())

	// default
	router.NoRoute(ctx.NoRoute())
}
