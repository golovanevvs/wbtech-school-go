package handler

import (
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type Handler struct {
	lg     zlog.Zerolog
	Router *ginext.Engine
}

func New(cfg *Config) *Handler {
	rt := ginext.New(cfg.GinMode)
	hd := &Handler{
		Router: rt,
	}

	hd.Router.GET("/notify", hd.createNotice)

	return hd
}
