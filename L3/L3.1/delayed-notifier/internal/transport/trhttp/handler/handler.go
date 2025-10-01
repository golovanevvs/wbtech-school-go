package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
)

type Handler struct {
	Router *ginext.Engine
}

func New(cfg *Config) *Handler {
	gin.SetMode(cfg.GinMode)
	rt := ginext.New()
	hd := &Handler{
		Router: rt,
	}

	hd.Router.GET("/create", hd.create)

	return hd
}

func (hd *Handler) create(c *ginext.Context) {

}
