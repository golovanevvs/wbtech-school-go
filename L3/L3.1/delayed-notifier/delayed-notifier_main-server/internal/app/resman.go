package app

import "github.com/wb-go/wbf/zlog"

type resourceManager struct {
	resources []func() error
}

func (rm *resourceManager) addResource(closeFunc func() error) {
	rm.resources = append(rm.resources, closeFunc)
}

func (rm *resourceManager) closeAll() error {
	var lastErr error
	for i := len(rm.resources) - 1; i >= 0; i-- {
		if err := rm.resources[i](); err != nil {
			lastErr = err
		} else {
			zlog.Logger.Debug().Msg("resource closed")
		}

	}
	return lastErr
}
