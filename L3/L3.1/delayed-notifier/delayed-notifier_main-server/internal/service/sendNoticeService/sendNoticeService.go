package sendNoticeService

import (
	"context"
	"fmt"
	"sync"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgEmail"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgTelegram"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
)

type IRepository interface {
	LoadTelegramChatID(ctx context.Context, username string) (chatID int64, err error)
}

type SendNoticeService struct {
	lg            zlog.Zerolog
	tg            *pkgTelegram.Client
	em            *pkgEmail.Client
	rp            IRepository
	retryStrategy retry.Strategy
}

func New(cfg *Config, tg *pkgTelegram.Client, em *pkgEmail.Client, rp IRepository) *SendNoticeService {
	lg := zlog.Logger.With().Str("component", "service-sendNoticeService").Logger()
	return &SendNoticeService{
		lg:            lg,
		tg:            tg,
		em:            em,
		rp:            rp,
		retryStrategy: retry.Strategy(cfg.RetryStrategy),
	}
}

func (sv *SendNoticeService) SendNotice(ctx context.Context, notice model.Notice) {
	wg := sync.WaitGroup{}
	for _, ch := range notice.Channels {
		wg.Go(func() {
			switch ch.Type {
			case model.ChannelTelegram:
				sv.SendNoticeToTelegram(ctx, ch.Value, notice)
			case model.ChannelEmail:
				sv.SendNoticeToEmail(ctx, ch.Value, notice)
			}
		})
	}
	wg.Wait()
}

func (sv *SendNoticeService) SendNoticeToTelegram(ctx context.Context, username string, notice model.Notice) error {
	chatID, err := sv.rp.LoadTelegramChatID(ctx, username)
	if err != nil {
		sv.lg.Error().Err(err).Msg("failed to load telegram chat id")
		return fmt.Errorf("failed to load telegram chat id: %w", err)
	}

	fn := func() error {
		err := sv.tg.SendTo(chatID, notice.Message)
		if err != nil {
			sv.lg.Warn().Err(err).Int64("chat ID", chatID).Msg("failed to send notice to telegram")
		}
		return err
	}

	if err := retry.Do(fn, sv.retryStrategy); err != nil {
		sv.lg.Error().Err(err).Int("attempts", sv.retryStrategy.Attempts).Msg("failed to send notice to telegram after all attempts")
		return fmt.Errorf("failed to send notice to telegram: %w", err)
	}

	return nil
}

func (sv *SendNoticeService) SendNoticeToEmail(ctx context.Context, email string, notice model.Notice) error {
	fn := func() error {
		err := sv.em.SendEmail([]string{email}, "delayed-notifier", notice.Message, false)
		if err != nil {
			sv.lg.Warn().Err(err).Str("e-mail", email).Msg("failed to send notice to e-mail")
		}
		return err
	}

	if err := retry.Do(fn, sv.retryStrategy); err != nil {
		sv.lg.Error().Err(err).Int("attempts", sv.retryStrategy.Attempts).Msg("failed to send notice to e-mail after all attempts")
		return fmt.Errorf("failed to send notice to e-mail: %w", err)
	}

	return nil
}
