package handler

import (
	"github.com/wb-go/wbf/ginext"
)

type Handler struct {
	Router *ginext.Engine
}

func New(cfg *Config) *Handler {
	rt := ginext.New(cfg.GinMode)
	hd := &Handler{
		Router: rt,
	}

	hd.Router.GET("/create", hd.create)

	return hd
}

func (hd *Handler) create(c *ginext.Context) {

}
