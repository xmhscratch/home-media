package main

import (
	"errors"

	"github.com/gin-gonic/gin"

	"home-media/streaming/api/routers"
	"home-media/streaming/core"
	"log"
)

// Server comment
type Server struct {
	Config *core.Config
	Router *gin.Engine
}

// NewServer comment
func NewServer(cfg *core.Config) (srv *Server, err error) {
	router := gin.Default()
	router.MaxMultipartMemory = 32 << 20 // 32 MiB
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false
	router.RemoveExtraSlash = true

	route, err := routers.NewRoute(cfg)
	if err != nil {
		return nil, err
	}
	route.Init(router)

	srv = &Server{
		Config: cfg,
		Router: router,
	}

	return srv, err
}

// Start comment
func (ctx *Server) Start() (err error) {
	if ctx.Router == nil {
		return errors.New("server is uninitialized")
	}
	portNumber := core.BuildString(":", ctx.Config.Port)
	return ctx.Router.Run(portNumber)
}

func main() {
	defer func() {
		if rec := recover(); rec != nil {
			log.Println(rec)
			return
		}
	}()

	cfg, err := core.NewConfig("../")
	if err != nil {
		panic(err)
	}

	srv, err := NewServer(cfg)
	if err != nil {
		panic(err)
	}

	if err := srv.Start(); err != nil {
		panic(err)
	}
}
