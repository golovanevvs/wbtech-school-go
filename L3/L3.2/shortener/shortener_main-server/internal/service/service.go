package service

// import (
// 	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgEmail"
// 	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgRabbitmq"
// 	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgRetry"
// 	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/pkg/pkgTelegram"
// 	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/service/addNoticeService"
// 	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/service/consumeNoticeService"
// 	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/service/deleteNoticeService"
// 	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/service/getNoticeService"
// 	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/service/sendNoticeService"
// 	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/service/telegramStartService"
// 	"github.com/golovanevvs/wbtech-school-go/tree/main/L3/L3.2/shortener/shortener_main-server/internal/service/updateNoticeService"
// 	"github.com/wb-go/wbf/zlog"
// )

// type iRepository interface {
// 	addNoticeService.ISaveNoticeRepository
// 	deleteNoticeService.IDelRepository
// 	getNoticeService.IRepository
// 	sendNoticeService.IRepository
// 	updateNoticeService.IUpdateNoticeRepository
// 	telegramStartService.IRepository
// }

// type Service struct {
// 	*addNoticeService.AddNoticeService
// 	*deleteNoticeService.DeleteNoticeService
// 	*getNoticeService.GetNoticeService
// 	*telegramStartService.TelegramStartService
// 	*consumeNoticeService.ConsumeNoticeService
// 	*sendNoticeService.SendNoticeService
// 	*updateNoticeService.UpdateNoticeService
// }

// func New(rs *pkgRetry.Retry, rp iRepository, rb *pkgRabbitmq.Client, tg *pkgTelegram.Client, em *pkgEmail.Client) *Service {
// 	lg := zlog.Logger.With().Str("layer", "service").Logger()
// 	getNotSv := getNoticeService.New(&lg, rp)
// 	updNotSv := updateNoticeService.New(&lg, rp)
// 	delNotSv := deleteNoticeService.New(&lg, rp, getNotSv, updNotSv)
// 	sendNotSv := sendNoticeService.New(&lg, rs, tg, em, rp)
// 	return &Service{
// 		AddNoticeService:     addNoticeService.New(&lg, rb, delNotSv, rp),
// 		DeleteNoticeService:  delNotSv,
// 		GetNoticeService:     getNotSv,
// 		TelegramStartService: telegramStartService.New(&lg, tg, rp),
// 		ConsumeNoticeService: consumeNoticeService.New(&lg, rb, delNotSv, sendNotSv, getNotSv, updNotSv),
// 		SendNoticeService:    sendNotSv,
// 		UpdateNoticeService:  updNotSv,
// 	}
// }
