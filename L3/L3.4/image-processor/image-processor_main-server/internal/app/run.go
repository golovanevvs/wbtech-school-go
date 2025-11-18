package app

import (
	"context"
	"fmt"
	"time"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.4/image-processor/image-processor_main-server/internal/pkg/pkgKafka"
)

func (a *App) Run(cancel context.CancelFunc) error {
	a.deps.tr.HTTP.RunServer(cancel)
	time.Sleep(500 * time.Millisecond)

	if err := a.deps.tr.HTTP.WaitForServer(a.cfg.tr.TrHTTP.PublicHost); err != nil {
		return fmt.Errorf("failed to start http server: %w", err)
	}

	go func() {
		a.lg.Debug().Msgf("%s Starting Kafka consumer...", pkgConst.Info)
		err := a.deps.kc.StartConsuming(context.Background(), func(msg pkgKafka.ProcessImageMessage) error {
			a.lg.Debug().Any("msg", msg).Msgf("%s received message", pkgConst.OpSuccess)
			return a.deps.sv.ProcessImage(context.Background(), msg.ImageID)
		})
		if err != nil {
			a.lg.Error().Err(err).Msgf("%s consumer error", pkgConst.Error)
		}
	}()

	return nil
}
