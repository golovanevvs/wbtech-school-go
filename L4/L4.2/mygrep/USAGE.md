# Инструкции по установке и запуску mygrep

## Установка

### Требования
- Go 1.19 или выше
- Операционная система: Linux, macOS, Windows

### Сборка проекта
```bash
# Перейдите в директорию проекта
cd L4/L4.2/mygrep

# Скомпилируйте проект
go build ./cmd/mygrep

# Проверьте, что исполняемый файл создался
ls -la mygrep*
```

## Запуск

### Локальный режим (аналог grep)

#### Базовое использование
```bash
# Поиск в файле
./mygrep "pattern" filename.txt

# Поиск в stdin
echo "hello world" | ./mygrep "hello"

# Поиск в нескольких файлах
./mygrep "pattern" file1.txt file2.txt file3.txt
```

#### Использование флагов
```bash
# С номерами строк
./mygrep -n "pattern" filename.txt

# Инвертированный поиск (строки БЕЗ pattern)
./mygrep -v "pattern" filename.txt

# Без учёта регистра
./mygrep -i "PATTERN" filename.txt

# Только полные строки
./mygrep -x "^full line match$" filename.txt

# Подсчёт совпадений
./mygrep -c "pattern" filename.txt

# Только совпадающие части
./mygrep -o "pattern" filename.txt

# Комбинация флагов
./mygrep -niv "pattern" filename.txt
```

### Distributed режим

#### Подготовка
1. Убедитесь, что порты свободны (по умолчанию используются 8080, 8081, 8082)
2. Подготовьте файл для обработки (желательно большой, 100+ строк)
3. Откройте несколько терминалов для запуска серверов

#### Запуск одного сервера
```bash
# Простой запуск на порту 8080
./mygrep -n "pattern" largefile.txt -port=8080

# С указанием server-id
./mygrep -n "pattern" largefile.txt -port=8080 -server-id="my-server:8080"
```

#### Запуск нескольких серверов

**Сценарий 1: Два сервера**
```bash
# Terminal 1
./mygrep -n "test" largefile.txt -port=8080 -server-id="server1:8080" -peers="localhost:8081"

# Terminal 2
./mygrep -n "test" largefile.txt -port=8081 -server-id="server2:8081" -peers="localhost:8080"
```

**Сценарий 2: Три сервера (рекомендуется)**
```bash
# Terminal 1
./mygrep -n "test" largefile.txt -port=8080 -server-id="server1:8080" -peers="localhost:8081,localhost:8082"

# Terminal 2
./mygrep -n "test" largefile.txt -port=8081 -server-id="server2:8081" -peers="localhost:8080,localhost:8082"

# Terminal 3
./mygrep -n "test" largefile.txt -port=8082 -server-id="server3:8082" -peers="localhost:8080,localhost:8081"
```

#### Запуск сервера без обработки (только приём заданий)
```bash
# Сервер будет ждать заданий от других серверов
./mygrep -port=8080
```

## Параметры командной строки

### Основные флаги (как в grep)
- `-n, --line-number` - показать номера строк
- `-v, --invert-match` - инвертировать совпадения (показать строки БЕЗ pattern)
- `-i, --ignore-case` - игнорировать регистр
- `-x, --line-regexp` - искать только полные строки
- `-c, --count` - показать только количество совпадений
- `-o, --only-matching` - показать только совпадающие части
- `--color` - выделить совпадения цветом

### Distributed флаги
- `-port PORT` - порт для TCP сервера (включает distributed режим)
- `-peers PEERS` - список пиров через запятую (например: "localhost:8081,localhost:8082")
- `-server-id ID` - идентификатор сервера (по умолчанию: hostname:port)

### Флаги паттерна
- `-pattern PATTERN` - паттерн для поиска
- `-e PATTERN` - альтеряс для -pattern

## Примеры использования

### Пример 1: Поиск в логах
```bash
# Поиск ошибок в лог-файле
./mygrep -n "ERROR" application.log

# Поиск по времени с инвертированным匹配
./mygrep -iv "debug" production.log | grep -v "INFO"
```

### Пример 2: Distributed поиск по большому файлу
```bash
# Подготовка: создание большого файла
for i in {1..1000}; do
    echo "Line $i: This is test line number $i with some content" >> bigfile.txt
done

# Запуск distributed поиска
./mygrep -n "test line" bigfile.txt -port=8080 -peers="localhost:8081,localhost:8082"
```

### Пример 3: Поиск в коде
```bash
# Поиск TODO комментариев
./mygrep -n "TODO\|FIXME\|XXX" *.go

# Поиск функций
./mygrep -x "^func " *.go

# Подсчёт строк кода
./mygrep -c "^[a-zA-Z]" *.go
```

### Пример 4: Комбинирование с другими командами
```bash
# Поиск и сортировка результатов
./mygrep "error" logfile.txt | sort | uniq -c

# Поиск и подсчёт
./mygrep -o "[0-9]+\.[0-9]+\.[0-9]+" access.log | sort | uniq -c

# Поиск в архивах
zcat access.log.gz | ./mygrep "404"
```

## Диагностика

### Проверка работы
```bash
# Проверка версии и справки (если реализовано)
./mygrep --help
./mygrep --version

# Проверка сетевых портов
netstat -an | grep 808

# Проверка процессов
ps aux | grep mygrep
```

### Логи и отладка
При запуске distributed режима вы увидите логи:
```
Запуск distributed режима. ServerID: hostname:8080
TCP сервер запущен на порту 8080
Запуск кворума: 3 серверов, нужно 2 голосов
Кворум: получено 1 из 3 голосов (нужно 2)
Кворум: получено 2 из 3 голосов (нужно 2)
Кворум достигнут!
Объединяем результаты от 3 серверов
```

### Решение проблем

#### Порт занят
```bash
# Найти процесс, использующий порт
lsof -i :8080
# или
netstat -an | grep 8080

# Изменить порт
./mygrep -port=8083 -peers="localhost:8081,localhost:8082"
```

#### Серверы не соединяются
1. Проверьте firewall
2. Убедитесь, что используете правильные адреса
3. Проверьте, что все серверы видят один и тот же файл

#### Неправильные результаты
1. Убедитесь, что файл одинаковый на всех серверах
2. Проверьте, что паттерн корректный
3. Сравните с результатами обычного grep

## Производительность

### Рекомендации
- Для файлов меньше 1000 строк используйте локальный режим
- Distributed режим эффективен для файлов 10,000+ строк
- Оптимальное количество серверов: 2-5 (больше не всегда лучше)
- Убедитесь, что серверы имеют доступ к одинаковым файлам

### Бенчмарки
```bash
# Сравнение производительности
time ./mygrep "pattern" largefile.txt  # Локальный
time ./mygrep "pattern" largefile.txt -port=8080 -peers="localhost:8081,localhost:8082"  # Distributed
```

## Безопасность

### Рекомендации
- Не запускайте серверы на публичных портах без аутентификации
- Используйте файрвол для ограничения доступа
- Убедитесь, что файлы доступны только авторизованным серверам
- Мониторьте сетевой трафик в продакшене

## Поддержка

### Получение справки
```bash
# Справка по флагам
./mygrep --help  # если реализовано

# Примеры использования
./mygrep  # без аргументов покажет справку
```

### Сообщение об ошибках
При возникновении проблем:
1. Проверьте логи сервера
2. Убедитесь, что используете последнюю версию
3. Попробуйте воспроизвести проблему с минимальным примером
4. Проверьте совместимость с оригинальным grep