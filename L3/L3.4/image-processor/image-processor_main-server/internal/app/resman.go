package app

import (
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/pkg/pkgConst"
	"github.com/wb-go/wbf/zlog"
)

type resourceManager struct {
	resources []resource
}

type resource struct {
	name      string
	closeFunc func() error
}

func (rm *resourceManager) addResource(resource resource) {
	rm.resources = append(rm.resources, resource)
}

func (rm *resourceManager) closeAll() error {
	var lastErr error
	for i := len(rm.resources) - 1; i >= 0; i-- {
		if err := rm.resources[i].closeFunc(); err != nil {
			zlog.Logger.Error().Err(err).Str("resource", rm.resources[i].name).Msgf("%s failed to resource closing", pkgConst.Error)
			lastErr = err
		} else {
			zlog.Logger.Debug().Str("resource", rm.resources[i].name).Msgf("%s resource closed", pkgConst.OpSuccess)
		}
	}
	return lastErr
}
