# Image Processor

Веб-приложение для асинхронной обработки изображений.

Протестировать работу сервиса можно по ссылке:

[https://incompletely-elemental-moonfish.cloudpub.ru/image-processor_web-client/](https://incompletely-elemental-moonfish.cloudpub.ru/image-processor_web-client/)

## Возможности

- **Загрузка изображений** с выбором операций обработки
- **Операции обработки:**
  - Изменение размера
  - Наложение водяного знака
  - Создание миниатюр
- **Отслеживание статуса** обработки в реальном времени
- **Просмотр и скачивание** обработанных изображений
- **Управление изображениями** - удаление загруженных файлов
- **Адаптивный интерфейс** для всех устройств

## Архитектура

### Backend (Go)

- **Язык:** Go
- **Фреймворк:** Gin
- **База данных:** PostgreSQL
- **Очередь задач:** Apache Kafka
- **Архитектура:** Трёхслойная (handlers → service → repository)
- **Интерфейсы:** IService, IRepository
- **Обработка изображений:** `golang.org/x/image`, `github.com/nfnt/resize`
- **Файловое хранилище:** локальная файловая система
- **Миграции:** `migrate`

### Frontend (Next.js)

- **Язык:** TypeScript
- **Фреймворк:** Next.js 16 (App Router)
- **UI:** Material UI
- **Состояние:** React Hooks
- **Стили:** Emotion + MUI
- **Формы:** встроенные элементы

## API Endpoints

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/upload` | Загрузить изображение с выбором операций |
| `GET` | `/image/{id}` | Получить статус и информацию об изображении |
| `DELETE` | `/image/{id}` | Удалить изображение и файлы |
| `GET` | `/uploads/{filename}` | Получить обработанное изображение |

## Функционал веб-интерфейса

- **Интуитивная форма загрузки** с выбором операций обработки
- **Визуализация статуса** обработки: `uploading → processing → completed/failed`
- **Real-time отслеживание** статуса через WebSocket/SSE или polling
- **Галерея изображений** с центрированным адаптивным дизайном
- **Просмотр и скачивание** обработанных изображений
- **Управление карточками** изображений с возможностью удаления

## Технологии

**Backend:**

- Go
- Gin Web Framework
- PostgreSQL
- Apache Kafka
- golang.org/x/image

**Frontend:**

- Next.js 16
- TypeScript
- Material UI (MUI)
- Emotion
- React Hooks
