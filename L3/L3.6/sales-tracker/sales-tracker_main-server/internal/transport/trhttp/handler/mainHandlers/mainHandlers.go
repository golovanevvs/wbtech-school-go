package mainHandlers

import (
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

type IService interface {
}

type MainHandlers struct {
	lg *zlog.Zerolog
	rt *ginext.Engine
	sv IService
}

func New(parentLg *zlog.Zerolog, rt *ginext.Engine, sv IService) *MainHandlers {
	lg := parentLg.With().Str("component", "ImageProcessor").Logger()
	return &MainHandlers{
		lg: &lg,
		rt: rt,
		sv: sv,
	}
}

func (hd *MainHandlers) RegisterRoutes() {

}
