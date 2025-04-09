package main

import (
	"errors"
	"log"
	"net"

	"github.com/gin-gonic/gin"

	"home-media/cmd/api/routers"
	"home-media/sys"
)

// Server comment
type Server struct {
	Config *sys.Config
	Router *gin.Engine
}

// NewServer comment
func NewServer(cfg *sys.Config) (srv *Server, err error) {
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

	_, port, err := net.SplitHostPort(ctx.Config.EndPoint["api"])
	if err != nil {
		return err
	}

	return ctx.Router.Run(net.JoinHostPort("0.0.0.0", port))
}

func main() {
	defer func() {
		if rec := recover(); rec != nil {
			log.Println(rec)
			return
		}
	}()

	cfg, err := sys.NewConfig("../")
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
