# shortener-opt (оптимизация API сервиса с профилировкой)

## Введение

**shortener-opt** - сервис сокращения URL с оптимизацией производительности.

Проект состоит из двух частей:

- [**shortener-opt_main-server**](shortener-opt_main-server) (Go) — отвечает за API и работу с данными.  
- [**shortener-opt_web-client**](shortener-opt_web-client) (Next.js) — веб-интерфейс для работы с сервисом.

В этом проекте выполнена оптимизация HTTP API сервиса сокращения URL (shortener) по CPU и памяти с использованием:

- **pprof** - профилирование CPU и памяти
- **benchstat** - статистическое сравнение бенчмарков
- **trace** - трассировка выполнения

## Оптимизации

### 1. Оптимизация генерации короткого кода

**Проблема**: Использование `uuid.New()` создаёт множество аллокаций.

**Решение**: Заменено на `math/rand` с ручной конкатенацией байтов и `sync.Pool`.

**Было**:

```go
func (sv *AddShortURLService) generateShortCode() string {
    short := uuid.New()
    return base62.EncodeToString(short[:])[:8]
}
```

**Стало**:

```go
var byteSlicePool = sync.Pool{
    New: func() interface{} {
        b := make([]byte, 16)
        return &b
    },
}

func (sv *AddShortURLService) generateShortCode() string {
    bufPtr := byteSlicePool.Get().(*[]byte)
    buf := *bufPtr
    defer byteSlicePool.Put(bufPtr)

    binary.LittleEndian.PutUint64(buf[:8], uint64(time.Now().UnixNano()))
    binary.LittleEndian.PutUint64(buf[8:], uint64(rand.Uint64()))

    return base62.EncodeToString(buf[:])[:8]
}
```

### 2. Добавление pprof в HTTP сервер

**Реализовано**:

- Отдельный порт для pprof (по умолчанию 6060);
- Поддержка `/debug/pprof/` эндпоинтов;
- Возможность записи профилей через переменные окружения:
  - `CPU_PROFILE=/path/to/cpu.prof`;
  - `MEM_PROFILE=/path/to/mem.prof`;
  - `TRACE_FILE=/path/to/trace.out`.

## Команды профилирования

Все команды профилирования можно вводить с использованием Makefile. Команды необходимо вводить, находясь в директории `./L4/L4.5/shortener-opt/shortener-opt_main-server`.

### Benchmarks

| Команда | Описание |
| ------- | -------- |
| `make install-benchstat` | Установка benchstat |
| `make bench` | Запуск всех бенчмарков |
| `make bench-old` | Запуск бенчмарка old версии (uuid) |
| `make bench-new` | Запуск бенчмарка new версии (sync.Pool) |
| `make bench-save-old` | Запуск old бенчмарка с сохранением в /tmp/bench_old.txt |
| `make bench-save-new` | Запуск new бенчмарка с сохранением в /tmp/bench_new.txt |
| `make bench-compare` | Сравнение результатов через benchstat |

### Profiling

| Команда | Описание |
| ------- | -------- |
| `make bench-cpu` | CPU профилирование → /tmp/cpu.prof |
| `make bench-mem` | Memory профилирование → /tmp/mem.prof |
| `make bench-trace` | Trace профилирование → /tmp/trace.out |

### Просмотр профилей

| Команда | Описание |
| ------- | -------- |
| `make profile-view-cpu` | Открыть CPU профиль в интерактивном режиме (команды: top, web, list) |
| `make profile-view-mem` | Открыть memory профиль (с флагом -alloc_space) |
| `make profile-view-trace` | Открыть trace в браузере |

### Утилиты

| Команда | Описание |
| ------- | -------- |
| `make profile-clean` | Удалить все созданные профили |

### Пример использования

Установить benchstat (один раз)

```bash
make install-benchstat
```

Сравнить производительность

```bash
make bench-compare
```

Создать CPU профиль

```bash
make bench-cpu
```

Просмотреть профиль

```bash
make profile-view-cpu
```

Очистить

```bash
make profile-clean
```

## Результаты профилирования

### Результат benchstat

```bash
goos: darwin
goarch: amd64
pkg: .../addShortURL
cpu: Intel(R) Core(TM) i5-1038NG7 CPU @ 2.00GHz

                  │ /tmp/bench_old.txt │         /tmp/bench_new.txt          │
                  │       sec/op       │   sec/op     vs base                │
GenerateShortCode          632.9n ± 4%   243.1n ± 7%  -61.59% (p=0.000 n=10)

                  │ /tmp/bench_old.txt │         /tmp/bench_new.txt         │
                  │        B/op        │    B/op     vs base                │
GenerateShortCode           48.00 ± 0%   32.00 ± 0%  -33.33% (p=0.000 n=10)

                  │ /tmp/bench_old.txt │         /tmp/bench_new.txt         │
                  │     allocs/op      │ allocs/op   vs base                │
GenerateShortCode           2.000 ± 0%   1.000 ± 0%  -50.00% (p=0.000 n=10)
```

### Итоговое сравнение

| Метрика | До (uuid) | После (sync.Pool) | Улучшение |
| ------- | --------- | ----------------- | --------- |
| **CPU время** | ~633 ns/op | ~243 ns/op | **-61.6% (2.6x)** |
| **Память** | 48 B/op | 32 B/op | **-33%** |
| **Аллокации** | 2 allocs/op | 1 allocs/op | **-50%** |

Тестировалось на: Intel(R) Core(TM) i5-1038NG7 CPU @ 2.00GHz.
Количество прогонов: 10.
