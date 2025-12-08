# Analyzator - Утилита анализа GC и памяти

Утилита для мониторинга и анализа сборщика мусора (GC) и использования памяти в Go приложениях. Предоставляет HTTP-endpoint в формате Prometheus для получения метрик.

## Возможности

- **Prometheus метрики**: Экспорт метрик в формате Prometheus
- **GC мониторинг**: Отслеживание количества сборок мусора и времени их выполнения
- **Memory анализ**: Мониторинг использования памяти (heap, stack, другие типы)
- **HTTP профилирование**: Встроенная поддержка pprof для профилирования
- **Настройка GC**: Возможность изменения процента GC
- **Health checks**: Endpoints для проверки состояния сервиса

## Установка и запуск

### Требования
- Go 1.25.4 или выше

### Сборка
```bash
go build ./cmd/analyzator
```

### Запуск с параметрами по умолчанию
```bash
./analyzator
```

### Запуск с пользовательскими параметрами
```bash
./analyzator \
  --port=8080 \
  --gc=100 \
  --read-timeout=10s \
  --write-timeout=10s \
  --idle-timeout=60s
```

### Параметры командной строки

- `--port` (по умолчанию: 8080) - Порт для запуска сервера
- `--gc` (по умолчанию: -1) - Процент для GC (-1 означает использовать настройки по умолчанию)
- `--read-timeout` (по умолчанию: 10s) - Таймаут чтения HTTP запросов
- `--write-timeout` (по умолчанию: 10s) - Таймаут записи HTTP ответов
- `--idle-timeout` (по умолчанию: 60s) - Таймаут простоя HTTP соединений

## Endpoints

### Основные метрики
- `GET /metrics` - Все метрики в формате Prometheus

### API endpoints
- `GET /api/v1/health` - Проверка здоровья сервиса
- `GET /api/v1/memory` - Информация о памяти
- `POST /api/v1/gc/trigger` - Принудительный запуск сборки мусора
- `GET /api/v1/gc/status` - Статус GC

### Профилирование (pprof)
- `GET /debug/pprof/` - Индекс профилирования
- `GET /debug/pprof/cmdline` - Команда запуска
- `GET /debug/pprof/profile` - CPU профиль (30 секунд по умолчанию)
- `GET /debug/pprof/symbol` - Символы
- `GET /debug/pprof/trace` - Трассировка
- `GET /debug/pprof/goroutine` - Goroutines
- `GET /debug/pprof/heap` - Heap профиль
- `GET /debug/pprof/allocs` - Allocations профиль
- `GET /debug/pprof/block` - Block профиль
- `GET /debug/pprof/mutex` - Mutex профиль

## Примеры запросов

### Получение всех метрик
```bash
curl http://localhost:8080/metrics
```

### Проверка здоровья
```bash
curl http://localhost:8080/api/v1/health
```

### Информация о памяти
```bash
curl http://localhost:8080/api/v1/memory
```

### Принудительный запуск GC
```bash
curl -X POST http://localhost:8080/api/v1/gc/trigger
```

### Статус GC
```bash
curl http://localhost:8080/api/v1/gc/status
```

### Получение heap профиля
```bash
curl http://localhost:8080/debug/pprof/heap > heap.prof
```

### Получение CPU профиля (30 секунд)
```bash
curl http://localhost:8080/debug/pprof/profile > cpu.prof
```

## Метрики Prometheus

### HTTP метрики
- `http_requests_total` - Общее количество HTTP запросов (метки: status_code, method, path)
- `http_request_duration_seconds` - Время выполнения HTTP запросов (метки: method, path)

### Memory метрики
- `go_memory_alloc_bytes` - Количество байт, выделенных в данный момент (метки: type)
- `go_memory_sys_bytes` - Количество байт, полученных от системы (метки: type)

### GC метрики
- `go_gc_count_total` - Общее количество сборок мусора (метки: gc_type)
- `go_gc_pause_seconds` - Время пауз GC (метки: gc_type)
- `go_gc_last_duration_seconds` - Время последней сборки мусора (метки: gc_type)

## Мониторинг с Prometheus

Пример конфигурации Prometheus для мониторинга:

```yaml
scrape_configs:
  - job_name: 'analyzator'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 15s
```

## Пример использования с Grafana

1. Настройте Prometheus для сбора метрик
2. Импортируйте дашборд для Go приложений в Grafana
3. Используйте метрики для мониторинга:
   - `go_memory_alloc_bytes` для отслеживания использования памяти
   - `go_gc_count_total` для мониторинга частоты GC
   - `go_gc_pause_seconds` для анализа времени пауз GC
   - `http_requests_total` для мониторинга HTTP нагрузки

## Отладка и профилирование

### Анализ heap профиля
```bash
curl http://localhost:8080/debug/pprof/heap > heap.prof

go tool pprof ./analyzator heap.prof
```

### Анализ CPU профиля
```bash
curl http://localhost:8080/debug/pprof/profile > cpu.prof

go tool pprof ./analyzator cpu.prof
```

### Анализ goroutines
```bash
curl http://localhost:8080/debug/pprof/goroutine > goroutines.prof
go tool pprof ./analyzator goroutines.prof
```
