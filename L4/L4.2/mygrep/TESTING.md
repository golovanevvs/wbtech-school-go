# Тестирование mygrep

## Тестовые файлы

### test_basic.txt
```
Hello World
This is a test file
Contains multiple lines
Some lines have the word "test"
Other lines do not match
Testing grep functionality
Another test example
Just regular text
```

### test_large.txt (генерируется автоматически)
```
Line 1: Hello world
Line 2: Test line 2
Line 3: Another test
...
Line 1000: Final test line
```

## Тесты локального режима

### Базовые тесты
```bash
# Тест 1: Простой поиск
./mygrep "test" test_basic.txt
# Ожидаемый результат: строки 2,4,6,7 содержат "test"

# Тест 2: Поиск с номерами строк
./mygrep -n "test" test_basic.txt
# Ожидаемый результат: "2:This is a test file" и т.д.

# Тест 3: Инвертированный поиск
./mygrep -v "test" test_basic.txt
# Ожидаемый результат: строки БЕЗ слова "test"

# Тест 4: Поиск без учёта регистра
./mygrep -i "HELLO" test_basic.txt
# Ожидаемый результат: строка 1 "Hello World"

# Тест 5: Подсчёт совпадений
./mygrep -c "test" test_basic.txt
# Ожидаемый результат: "4"

# Тест 6: Поиск только совпадающих частей
./mygrep -o "test" test_basic.txt
# Ожидаемый результат: только слово "test" 4 раза
```

### Тесты потоков
```bash
# Тест 7: Чтение из stdin
echo -e "hello world\ntest line" | ./mygrep "test"
# Ожидаемый результат: "test line"

# Тест 8: Множественные файлы
./mygrep -n "test" test_basic.txt test_large.txt
# Ожидаемый результат: результаты из обоих файлов с префиксами
```

## Тесты distributed режима

### Подготовка к тестированию
1. Скомпилируйте проект: `go build ./cmd/mygrep`
2. Убедитесь, что порты 8080, 8081, 8082 свободны
3. Подготовьте тестовый файл с достаточным количеством строк (например, 100+ строк)

### Distributed тесты

#### Тест 1: Один сервер (локальный режим)
```bash
# Запуск сервера на порту 8080
./mygrep -n "test" test_large.txt -port=8080
# Ожидаемый результат: как обычный grep, но с информацией о distributed режиме
```

#### Тест 2: Два сервера
```bash
# Terminal 1: Сервер 1
./mygrep -n "test" test_large.txt -port=8080 -server-id="server1:8080" -peers="localhost:8081"

# Terminal 2: Сервер 2  
./mygrep -n "test" test_large.txt -port=8081 -server-id="server2:8081" -peers="localhost:8080"
```

#### Тест 3: Три сервера (полный кворум)
```bash
# Terminal 1
./mygrep -n "test" test_large.txt -port=8080 -server-id="server1:8080" -peers="localhost:8081,localhost:8082"

# Terminal 2
./mygrep -n "test" test_large.txt -port=8081 -server-id="server2:8081" -peers="localhost:8080,localhost:8082"

# Terminal 3
./mygrep -n "test" test_large.txt -port=8082 -server-id="server3:8082" -peers="localhost:8080,localhost:8081"
```

#### Тест 4: Тестирование отказоустойчивости
```bash
# Запустите 3 сервера как в тесте 3
# Затем остановите один из серверов (Ctrl+C)
# Система должна продолжить работу с оставшимися 2 серверами
```

#### Тест 5: Только сервер (без файлов)
```bash
# Запуск сервера без файлов (ожидание заданий)
./mygrep -port=8080
# Ожидаемый результат: "Сервер запущен. Ожидание команд..."
```

## Ожидаемые результаты

### Локальный режим
- Все результаты должны совпадать с `grep`
- Поддержка всех стандартных флагов grep
- Корректная обработка stdin и множественных файлов

### Distributed режим
- Файл должен разбиваться на части между серверами
- Каждый сервер обрабатывает свою часть
- Результаты должны объединяться корректно
- Кворум должен достигаться при N/2+1 ответах
- Система должна работать при отказе части серверов

### Производительность
- Distributed режим должен быть быстрее для больших файлов
- Параллельная обработка должна уменьшать общее время
- Кворум должен срабатывать быстрее ожидания всех серверов

## Сравнение с оригинальным grep

### Тесты эквивалентности
```bash
# Сравнение результатов
./mygrep "pattern" file.txt > mygrep_result.txt
grep "pattern" file.txt > grep_result.txt
diff mygrep_result.txt grep_result.txt  # Должно быть пусто

# Тест с различными флагами
./mygrep -nivxo "pattern" file.txt > mygrep_flags.txt
grep -nivxo "pattern" file.txt > grep_flags.txt
diff mygrep_flags.txt grep_flags.txt  # Должно быть пусто
```

## Диагностика проблем

### Проверка сетевого взаимодействия
```bash
# Проверка доступности портов
netstat -an | grep 808

# Проверка логов сервера
# Ищите сообщения:
# - "TCP сервер запущен на порту"
# - "Кворум: получено X из Y голосов"
# - "Объединяем результаты от X серверов"
```

### Типичные проблемы
1. **Порт занят**: Измените порт или освободите используемый
2. **Серверы не видят друг друга**: Проверьте сетевое подключение и firewall
3. **Кворум не достигается**: Проверьте, что все серверы запущены и отвечают
4. **Неправильные результаты**: Проверьте, что файл одинаковый на всех серверах

## Производительность

### Бенчмарки
```bash
# Тест на большом файле (10000+ строк)
time ./mygrep "pattern" large_file.txt  # Локальный
time ./mygrep "pattern" large_file.txt -port=8080 -peers="localhost:8081,localhost:8082"  # Distributed
```

### Ожидаемые улучшения
- 2 сервера: ускорение ~1.5-2x
- 3 сервера: ускорение ~2-3x
- Больше серверов: улучшение до насыщения I/O