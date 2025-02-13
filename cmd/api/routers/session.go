package routers

import (
	"encoding/json"
	"home-media/sys/session"

	"github.com/gin-gonic/gin"
)

func (ctx *RouteContext) CreateSession() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			err       error
			ss        *session.Session[session.FileSourceType]
			result    []byte
			sessionId string
		)

		sourceURL, _ := ginCtx.GetPostForm("data_source")
		sourceType, _ := ginCtx.GetPostForm("data_source_type")
		title, _ := ginCtx.GetPostForm("title")
		nodeId, _ := ginCtx.GetPostForm("id")
		rootId, _ := ginCtx.GetPostForm("root")

		switch sourceType {
		case session.FILE_SOURCE_TYPE_DIRECT.String():
			if ss, err = session.InitDirect(
				ctx.Config,
				sessionId,
				sourceURL,
				nodeId, rootId, title,
			); err != nil {
				ctx.Error(200, err.Error())(ginCtx)
				return
			}
		case session.FILE_SOURCE_TYPE_TORRENT.String():
			if ss, err = session.InitTorrent(
				ctx.Config,
				sessionId,
				sourceURL,
				nodeId, rootId, title,
			); err != nil {
				ctx.Error(200, err.Error())(ginCtx)
				return
			}
		default:
			break
		}
		// defer ctx.CheckProgress()(ginCtx)

		if result, err = json.Marshal(ss); err != nil {
			ctx.Error(200, err.Error())(ginCtx)
			return
		}
		ginCtx.Data(200, "application/json", result)
	}
}

func (ctx *RouteContext) GetFiles() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			sessionId string = ginCtx.Param("ssid")
			files     []string
		)

		if result, err := session.GetFiles(ctx.Config, sessionId); err != nil {
			ctx.Error(200, err.Error())(ginCtx)
			return
		} else {
			for _, f := range result {
				files = append(files, f)
			}
		}

		ginCtx.JSON(200, files)
	}
}

func (ctx *RouteContext) CheckProgress() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var sessionId string
		filePath := ginCtx.Param("filePath")
		if sessionId == "" {
			sessionId = ginCtx.Param("ssid")
		}

		if err := session.CreateDownload(
			ctx.Config,
			sessionId,
			filePath,
		); err != nil {
			ctx.Error(200, err.Error())(ginCtx)
			return
		}

		ginCtx.JSON(200, gin.H{})
	}
}

func (ctx *RouteContext) DownloadDirect() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			err error
			ss  *session.Session[session.FileSourceType]
		)

		if ss, err = session.InitDirect(
			ctx.Config,
			ginCtx.Param("ssid"),
		); err != nil {
			ctx.Error(200, err.Error())(ginCtx)
			return
		}

		if !ss.IsDownloadable() {
			ctx.Error(200, "This download is unacceptable")(ginCtx)
			return
		}

		// ginCtx.Request.URL.String()
		if err = ss.File.DownloadDirect(ginCtx, ginCtx.Param("filePath")); err != nil {
			ctx.Error(200, err.Error())(ginCtx)
			return
		}
	}
}

func (ctx *RouteContext) DownloadTorrent() gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		var (
			err error
			ss  *session.Session[session.FileSourceType]
		)

		if ss, err = session.InitTorrent(
			ctx.Config,
			ginCtx.Param("ssid"),
		); err != nil {
			ctx.Error(200, err.Error())(ginCtx)
			return
		}

		if !ss.IsDownloadable() {
			ctx.Error(200, "This download is unacceptable")(ginCtx)
			return
		}

		if err = ss.File.DownloadTorrent(ginCtx, ginCtx.Param("filePath")); err != nil {
			ctx.Error(200, err.Error())(ginCtx)
			return
		}
	}
}
