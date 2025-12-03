package app

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.2/mygrep/internal/model"
)

// Config содержит конфигурацию для локального режима
type Config struct {
	Pattern string
	Files   []string
	Flags   model.GrepFlags
	Input   io.Reader
	Output  io.Writer

	// Distributed режим
	IsDistributed bool
	Port          string
	Peers         []string
	ServerID      string
}

// GrepResult результат поиска
type GrepResult struct {
	LineNumber int
	Line       string
	Match      string
}

// Message структура для сетевого обмена
type Message struct {
	Type      string      `json:"type"`
	From      string      `json:"from"`
	To        string      `json:"to,omitempty"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

// Server структура для хранения информации о сервере
type Server struct {
	ID      string
	Address string
	Status  string // "online", "offline", "processing"
	Conn    net.Conn
}

// App представляет основное приложение
type App struct {
	config      *Config
	server      net.Listener
	servers     map[string]*Server
	serverMutex sync.RWMutex
	quit        chan bool
}

// NewApp создаёт новый экземпляр приложения
func NewApp() *App {
	return &App{
		servers: make(map[string]*Server),
		quit:    make(chan bool),
	}
}

// Run запускает приложение
func (a *App) Run() error {
	var err error
	a.config, err = a.parseFlags()
	if err != nil {
		return fmt.Errorf("ошибка парсинга флагов: %v", err)
	}

	// Если включен distributed режим, запускаем сервер
	if a.config.IsDistributed {
		if err := a.startDistributedMode(); err != nil {
			return err
		}
	} else {
		// Локальный режим
		if err := a.runGrep(); err != nil {
			return fmt.Errorf("ошибка выполнения grep: %v", err)
		}
	}

	return nil
}

// parseFlags парсит аргументы командной строки
func (a *App) parseFlags() (*Config, error) {
	var pattern string
	flag.StringVar(&pattern, "pattern", "", "Паттерн для поиска")
	flag.StringVar(&pattern, "e", "", "Паттерн для поиска (alias для --pattern)")

	color := flag.Bool("color", false, "Выделить совпадения цветом")
	invertMatch := flag.Bool("v", false, "Инвертировать совпадения")
	ignoreCase := flag.Bool("i", false, "Игнорировать регистр")
	wholeLine := flag.Bool("x", false, "Искать только полные строки")
	lineNumber := flag.Bool("n", false, "Показать номера строк")
	count := flag.Bool("c", false, "Показать только количество совпадений")
	onlyMatching := flag.Bool("o", false, "Показать только совпадающие части")

	// Distributed флаги
	port := flag.String("port", "", "Порт для distributed режима")
	peers := flag.String("peers", "", "Список пиров (через запятую)")
	serverID := flag.String("server-id", "", "ID сервера (по умолчанию - hostname:port)")

	flag.Parse()

	if pattern == "" {
		args := flag.Args()
		if len(args) == 0 {
			return nil, fmt.Errorf("не указан паттерн для поиска")
		}
		pattern = args[0]
		args = args[1:]

		config := &Config{
			Pattern: pattern,
			Files:   args,
			Flags: model.GrepFlags{
				Color:        *color,
				InvertMatch:  *invertMatch,
				IgnoreCase:   *ignoreCase,
				WholeLine:    *wholeLine,
				LineNumber:   *lineNumber,
				Count:        *count,
				OnlyMatching: *onlyMatching,
			},
			Input:         os.Stdin,
			Output:        os.Stdout,
			IsDistributed: *port != "" || *peers != "",
			Port:          *port,
			Peers:         parsePeers(*peers),
			ServerID:      *serverID,
		}
		return config, nil
	}

	config := &Config{
		Pattern: pattern,
		Files:   flag.Args(),
		Flags: model.GrepFlags{
			Color:        *color,
			InvertMatch:  *invertMatch,
			IgnoreCase:   *ignoreCase,
			WholeLine:    *wholeLine,
			LineNumber:   *lineNumber,
			Count:        *count,
			OnlyMatching: *onlyMatching,
		},
		Input:         os.Stdin,
		Output:        os.Stdout,
		IsDistributed: *port != "" || *peers != "",
		Port:          *port,
		Peers:         parsePeers(*peers),
		ServerID:      *serverID,
	}
	return config, nil
}

// parsePeers парсит список пиров
func parsePeers(peersStr string) []string {
	if peersStr == "" {
		return nil
	}
	return strings.Split(peersStr, ",")
}

// startDistributedMode запускает distributed режим
func (a *App) startDistributedMode() error {
	// Генерируем ServerID если не указан
	if a.config.ServerID == "" {
		hostname, _ := os.Hostname()
		a.config.ServerID = fmt.Sprintf("%s:%s", hostname, a.config.Port)
	}

	fmt.Printf("Запуск distributed режима. ServerID: %s\n", a.config.ServerID)

	// Запускаем TCP сервер
	if err := a.startTCPServer(); err != nil {
		return fmt.Errorf("ошибка запуска TCP сервера: %v", err)
	}

	// Подключаемся к пирам
	if len(a.config.Peers) > 0 {
		if err := a.connectToPeers(); err != nil {
			fmt.Printf("Предупреждение: не удалось подключиться к пирам: %v\n", err)
		}
	}

	// В distributed режиме ждём команды или запускаем обработку
	if len(a.config.Files) > 0 {
		return a.runDistributedGrep()
	}

	// Ждём команды (пока просто выводим статус)
	fmt.Println("Сервер запущен. Ожидание команд...")
	<-a.quit
	return nil
}

// startTCPServer запускает TCP сервер
func (a *App) startTCPServer() error {
	listener, err := net.Listen("tcp", ":"+a.config.Port)
	if err != nil {
		return err
	}
	a.server = listener

	go func() {
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-a.quit:
					return
				default:
					continue
				}
			}

			go a.handleConnection(conn)
		}
	}()

	fmt.Printf("TCP сервер запущен на порту %s\n", a.config.Port)
	return nil
}

// handleConnection обрабатывает входящее соединение
func (a *App) handleConnection(conn net.Conn) {
	defer conn.Close()

	var msg Message
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&msg); err != nil {
		return
	}

	switch msg.Type {
	case "status_request":
		a.sendStatusResponse(conn, msg.From)
	case "job_request":
		// TODO: реализовать обработку заданий
		fmt.Println("Получен запрос на выполнение задания (пока не реализован)")
	default:
		fmt.Printf("Получено неизвестное сообщение типа: %s\n", msg.Type)
	}
}

// sendStatusResponse отправляет ответ со статусом
func (a *App) sendStatusResponse(conn net.Conn, from string) {
	status := map[string]interface{}{
		"server_id": a.config.ServerID,
		"status":    "online",
		"timestamp": time.Now(),
	}

	response := Message{
		Type:      "status_response",
		From:      a.config.ServerID,
		To:        from,
		Data:      status,
		Timestamp: time.Now(),
	}

	json.NewEncoder(conn).Encode(response)
}

// connectToPeers подключается к пирам
func (a *App) connectToPeers() error {
	for _, peer := range a.config.Peers {
		conn, err := net.Dial("tcp", peer)
		if err != nil {
			fmt.Printf("Не удалось подключиться к %s: %v\n", peer, err)
			continue
		}

		// Отправляем запрос статуса
		msg := Message{
			Type:      "status_request",
			From:      a.config.ServerID,
			To:        peer,
			Data:      nil,
			Timestamp: time.Now(),
		}

		json.NewEncoder(conn).Encode(msg)
		conn.Close()
	}
	return nil
}

// runDistributedGrep выполняет distributed grep
func (a *App) runDistributedGrep() error {
	fmt.Println("Запуск distributed grep...")
	// TODO: реализовать распределённую обработку
	fmt.Println("Distributed grep ещё не реализован")
	return nil
}

// runGrep выполняет поиск (локальный режим)
func (a *App) runGrep() error {
	if len(a.config.Files) > 0 {
		for _, filename := range a.config.Files {
			if err := a.processFile(filename); err != nil {
				return fmt.Errorf("ошибка обработки файла %s: %v", filename, err)
			}
		}
	} else {
		if err := a.processStream(os.Stdin, "stdin"); err != nil {
			return err
		}
	}
	return nil
}

// processFile обрабатывает один файл
func (a *App) processFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл %s: %v", filename, err)
	}
	defer file.Close()

	return a.processStream(file, filename)
}

// processStream обрабатывает поток данных
func (a *App) processStream(reader io.Reader, sourceName string) error {
	scanner := bufio.NewScanner(reader)
	lineNumber := 0
	matches := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		result, found := a.searchInLine(line, lineNumber)
		if found {
			matches++
			if !a.config.Flags.Count {
				if len(a.config.Files) > 1 {
					fmt.Fprintf(a.config.Output, "%s:", sourceName)
				}
				a.printResult(result)
			}
		}
	}

	if a.config.Flags.Count {
		if len(a.config.Files) > 1 {
			fmt.Fprintf(a.config.Output, "%s:", sourceName)
		}
		fmt.Fprintf(a.config.Output, "%d\n", matches)
	}

	return nil
}

// searchInLine ищет совпадения в строке
func (a *App) searchInLine(line string, lineNumber int) (*GrepResult, bool) {
	pattern := regexp.QuoteMeta(a.config.Pattern)
	if a.config.Flags.WholeLine {
		pattern = "^" + pattern + "$"
	}

	flags := ""
	if a.config.Flags.IgnoreCase {
		flags = "(?i)"
	}

	fullPattern := flags + pattern
	re, err := regexp.Compile(fullPattern)
	if err != nil {
		return nil, false
	}

	matches := re.FindAllString(line, -1)

	if a.config.Flags.InvertMatch {
		if len(matches) == 0 {
			return &GrepResult{
				LineNumber: lineNumber,
				Line:       line,
				Match:      line,
			}, true
		}
		return nil, false
	}

	if len(matches) > 0 {
		if a.config.Flags.OnlyMatching {
			for _, match := range matches {
				return &GrepResult{
					LineNumber: lineNumber,
					Line:       match,
					Match:      match,
				}, true
			}
		}

		return &GrepResult{
			LineNumber: lineNumber,
			Line:       line,
			Match:      matches[0],
		}, true
	}

	return nil, false
}

// printResult выводит результат
func (a *App) printResult(result *GrepResult) {
	if a.config.Flags.LineNumber {
		fmt.Fprintf(a.config.Output, "%d:", result.LineNumber)
	}

	if a.config.Flags.OnlyMatching {
		fmt.Fprintf(a.config.Output, "%s\n", result.Match)
	} else {
		fmt.Fprintf(a.config.Output, "%s\n", result.Line)
	}
}
