# delayed-notifier

**delayed-notifier** - отложенные уведомления через очереди.

Проект состоит из двух частей:

- [**delayed-notifier_main-server**](delayed-notifier_main-server) (Go) — отвечает за API, хранение и отправку уведомлений.  
- [**delayed-notifier_web-client**](delayed-notifier_web-client) (Next.js) — веб-интерфейс для работы с сервисом.

Работу сервиса можно протестировать по ссылке:

[https://incompletely-elemental-moonfish.cloudpub.ru/delayed-notifier_web-client](https://incompletely-elemental-moonfish.cloudpub.ru/delayed-notifier_web-client) (доступно по мере работы сервера).

Для получения уведомления на Telegram предварительно необходимо подписаться на telegram-бота. Аккаунт Telegram необходимо указывать **без** символа "@".
