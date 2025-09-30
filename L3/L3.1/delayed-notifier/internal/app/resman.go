package app

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
		}
	}
	return lastErr
}
