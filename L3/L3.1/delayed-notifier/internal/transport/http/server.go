package transporthttp

import (
	"fmt"
	"net/http"

	"github.com/wb-go/wbf/ginext"
)

type Server struct {
	httpsrv *http.Server
}

func New(cfg *Config) *Server {
	return &Server{
		httpsrv: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Port),
			Handler: ,
		},
	}
}

func (srv *Server) InitRoutes() http.Handler {
	rt := ginext.New().Engine

	rt.GET("/eee", srv.ddd)
	return rt
}
