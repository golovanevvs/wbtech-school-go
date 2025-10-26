package sendNoticeService

import (
	"context"
	"sync"

	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/model"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgConst"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgEmail"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgErrors"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgRetry"
	"github.com/golovanevvs/wbtech-school-go/L3/L3.1/delayed-notifier/delayed-notifier_main-server/internal/pkg/pkgTelegram"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
)

type IRepository interface {
	LoadTelegramChatID(ctx context.Context, username string) (chatID int64, err error)
}

type SendNoticeService struct {
	lg *zlog.Zerolog
	rs *pkgRetry.Retry
	tg *pkgTelegram.Client
	em *pkgEmail.Client
	rp IRepository
}

func New(
	parentLg *zlog.Zerolog,
	rs *pkgRetry.Retry,
	tg *pkgTelegram.Client,
	em *pkgEmail.Client,
	rp IRepository,
) *SendNoticeService {
	lg := parentLg.With().Str("component", "SendNoticeService").Logger()
	return &SendNoticeService{
		lg: &lg,
		rs: rs,
		tg: tg,
		em: em,
		rp: rp,
	}
}

func (sv *SendNoticeService) SendNotice(ctx context.Context, notice model.Notice) {
	lg := sv.lg.With().Str("method", "SendNotice").Logger()
	lg.Trace().Msgf("%s method starting", pkgConst.Start)
	defer lg.Trace().Msgf("%s method stopped", pkgConst.Stop)

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
	lg := sv.lg.With().Str("method", "SendNoticeToTelegram").Logger()
	lg.Trace().Msgf("%s method starting", pkgConst.Start)
	defer lg.Trace().Msgf("%s method stopped", pkgConst.Stop)

	lg.Trace().Str("username", username).Int("notice ID", notice.ID).Msgf("%s loading chat ID from repository...", pkgConst.OpStart)
	chatID, err := sv.rp.LoadTelegramChatID(ctx, username)
	if err != nil {
		return pkgErrors.Wrapf(err, "load telegram chat id from repository;  username: %s, notice ID: %d", username, notice.ID)
	}
	lg.Trace().Str("username", username).Int("notice ID", notice.ID).Int64("chat ID", chatID).Msgf("%s chat ID loaded from repository successfully", pkgConst.OpSuccess)

	fn := func() error {
		lg.Trace().Str("username", username).Int("notice ID", notice.ID).Msgf("%s sending message to Telegram...", pkgConst.OpStart)
		err := sv.tg.SendTo(chatID, notice.Message)
		if err != nil {
			sv.lg.Warn().Err(err).Int64("chat ID", chatID).Int("notice ID", notice.ID).Msgf("%s failed to send notice to telegram", pkgConst.Warn)
			return err
		}
		lg.Trace().Str("username", username).Int("notice ID", notice.ID).Int64("chat ID", chatID).Msgf("%s message sended to Telegram successfully", pkgConst.OpSuccess)
		return nil
	}

	lg.Trace().Str("username", username).Int("notice ID", notice.ID).Msgf("%s sending message to Telegram with retry starting...", pkgConst.OpStart)
	if err := retry.Do(fn, retry.Strategy(*sv.rs)); err != nil {
		return pkgErrors.Wrapf(err,
			"send notice to telegram after all ettempts; chat ID: %d, notice ID: %d, attempts: %d",
			chatID, notice.ID, sv.rs.Attempts)
	}
	lg.Debug().Str("username", username).Int("notice ID", notice.ID).Int64("chat ID", chatID).Msgf("%s sending message to Telegram with retry completed", pkgConst.OpSuccess)

	return nil
}

func (sv *SendNoticeService) SendNoticeToEmail(ctx context.Context, email string, notice model.Notice) error {
	lg := sv.lg.With().Str("method", "SendNoticeToEmail").Logger()
	lg.Trace().Msgf("%s method starting", pkgConst.Start)
	defer lg.Trace().Msgf("%s method stopped", pkgConst.Stop)

	fn := func() error {
		lg.Trace().Str("e-mail", email).Int("notice ID", notice.ID).Msgf("%s sending message to e-mail...", pkgConst.OpStart)
		err := sv.em.SendEmail([]string{email}, "delayed-notifier", notice.Message, false)
		if err != nil {
			sv.lg.Warn().Err(err).Str("e-mail", email).Int("notice ID", notice.ID).Msgf("%s failed to send notice to e-mail", pkgConst.Warn)
		} else {
			lg.Trace().Str("e-mail", email).Int("notice ID", notice.ID).Msgf("%s message sended to e-mail successfully", pkgConst.OpSuccess)
		}
		return err
	}

	lg.Trace().Str("e-mail", email).Int("notice ID", notice.ID).Msgf("%s sending message to e-mail with retry starting...", pkgConst.OpStart)
	if err := retry.Do(fn, retry.Strategy(*sv.rs)); err != nil {
		return pkgErrors.Wrapf(err,
			"send notice to e-mail after all ettempts; e-mail: %s, notice ID: %d, attempts: %d",
			email, notice.ID, sv.rs.Attempts)
	}
	lg.Debug().Str("e-mail", email).Int("notice ID", notice.ID).Msgf("%s sending message to e-mail with retry completed", pkgConst.OpSuccess)

	return nil
}
