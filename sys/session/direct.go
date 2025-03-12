package session

import (
	"fmt"
	"home-media/sys"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	// "github.com/sanity-io/litter"
)

func (ctx *File[FileDirect]) InitDirect() (map[string]interface{}, error) {
	req, err := http.NewRequest("HEAD", ctx.SourceURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return map[string]interface{}{
		"fileSize": res.Header.Get("Content-Length"),
		"fileType": res.Header.Get("Content-Type"),
	}, err
}

func (ctx *File[FileDirect]) DownloadDirect(ginCtx *gin.Context, filePath string) (err error) {
	filePath = GetFilePath(filePath)

	fileKey := sys.GenerateID(ctx.NodeID, filePath)
	defer ctx.notify(fileKey, 0)

	req, err := http.NewRequest("GET", ctx.SourceURL, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	req.Header.Add("Accept", ginCtx.Request.Header.Get("Accept"))
	req.Header.Add("Accept-Encoding", ginCtx.Request.Header.Get("Accept-Encoding"))
	req.Header.Add("Accept-Language", ginCtx.Request.Header.Get("Accept-Language"))
	req.Header.Add("Accept-Charset", ginCtx.Request.Header.Get("Accept-Charset"))
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.87 Safari/537.36")

	if ginCtx.Request.Header.Get("Range") != "" {
		ginCtx.Writer.WriteHeader(http.StatusPartialContent)
		req.Header.Add("Range", ginCtx.Request.Header.Get("Range"))
	}
	defer ginCtx.Request.Body.Close()

	res, err := HTTPClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	ginCtx.Header("Accept-Ranges", "bytes")

	if ginCtx.Query("mime") == "" {
		fileName := ctx.GetFileName(filePath)
		fileExt := ctx.GetFileExt(filePath)
		ginCtx.Header("Content-Disposition", "attachment; filename="+fileName+fileExt)
	}

	for key, value := range res.Header {
		if key == "Content-Disposition" {
			continue
		}
		ginCtx.Header(key, value[0])
	}

	io.Copy(ginCtx.Writer, ratelimit.Reader(res.Body, DownloadBucket))
	return err
}
