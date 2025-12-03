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
		a.handleJobRequest(conn, msg)
	default:
		fmt.Printf("Получено неизвестное сообщение типа: %s\n", msg.Type)
	}
}

// handleJobRequest обрабатывает запрос на выполнение задания
func (a *App) handleJobRequest(conn net.Conn, msg Message) {
	fmt.Printf("Получен запрос на выполнение задания от %s\n", msg.From)

	// TODO: реализовать выполнение задания
	// 1. Извлечь данные задания из msg.Data
	// 2. Прочитать указанные строки из файла
	// 3. Выполнить grep на этих строках
	// 4. Отправить результат обратно

	jobData, ok := msg.Data.(map[string]interface{})
	if !ok {
		fmt.Println("Неверный формат данных задания")
		return
	}

	filename, _ := jobData["filename"].(string)
	startLine, _ := jobData["start_line"].(float64)
	endLine, _ := jobData["end_line"].(float64)
	pattern, _ := jobData["pattern"].(string)

	fmt.Printf("Задание: файл=%s, строки=%d-%d, паттерн=%s\n",
		filename, int(startLine), int(endLine), pattern)

	// Заглушка - отправляем успешный ответ
	result := map[string]interface{}{
		"job_id":       jobData["job_id"],
		"server_id":    a.config.ServerID,
		"matches":      []map[string]interface{}{},
		"processed":    int(endLine - startLine + 1),
		"success":      true,
		"completed_at": time.Now(),
	}

	response := Message{
		Type:      "job_response",
		From:      a.config.ServerID,
		To:        msg.From,
		Data:      result,
		Timestamp: time.Now(),
	}

	json.NewEncoder(conn).Encode(response)
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

	if len(a.config.Files) == 0 {
		return fmt.Errorf("в distributed режиме необходимо указать файлы для обработки")
	}

	// Для каждого файла запускаем distributed обработку
	for _, filename := range a.config.Files {
		if err := a.processFileDistributed(filename); err != nil {
			return fmt.Errorf("ошибка distributed обработки файла %s: %v", filename, err)
		}
	}

	return nil
}

// processFileDistributed обрабатывает файл в distributed режиме
func (a *App) processFileDistributed(filename string) error {
	fmt.Printf("Обработка файла %s в distributed режиме\n", filename)

	// Подсчитываем количество строк в файле
	totalLines, err := a.countLinesInFile(filename)
	if err != nil {
		return err
	}

	fmt.Printf("Файл содержит %d строк\n", totalLines)

	// Определяем количество серверов (текущий + пиы)
	numServers := 1 + len(a.config.Peers)
	if numServers == 1 {
		// Только один сервер, выполняем локально
		return a.processFile(filename)
	}

	// Разбиваем файл на части и отправляем задания
	return a.distributeAndProcess(filename, totalLines, numServers)
}

// countLinesInFile подсчитывает строки в файле
func (a *App) countLinesInFile(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := 0
	for scanner.Scan() {
		lines++
	}

	return lines, scanner.Err()
}

// distributeAndProcess распределяет задания между серверами
func (a *App) distributeAndProcess(filename string, totalLines, numServers int) error {
	linesPerServer := totalLines / numServers

	// Создаём задания для каждого сервера
	jobs := make([]map[string]interface{}, 0, numServers)

	for i := 0; i < numServers; i++ {
		startLine := i*linesPerServer + 1
		endLine := (i + 1) * linesPerServer

		// Последний сервер получает остаток строк
		if i == numServers-1 {
			endLine = totalLines
		}

		job := map[string]interface{}{
			"job_id":     fmt.Sprintf("%s-%d", a.config.ServerID, i),
			"server_id":  a.config.ServerID,
			"pattern":    a.config.Pattern,
			"start_line": startLine,
			"end_line":   endLine,
			"filename":   filename,
			"flags":      a.config.Flags,
			"created_at": time.Now(),
		}
		jobs = append(jobs, job)
	}

	// Если есть пиры, отправляем им задания
	if len(a.config.Peers) > 0 {
		return a.sendJobsToPeers(jobs)
	}

	// Если нет пиров, выполняем локально
	return a.executeJobsLocally(jobs)
}

// sendJobsToPeers отправляет задания пирам
func (a *App) sendJobsToPeers(jobs []map[string]interface{}) error {
	fmt.Printf("Отправляем %d заданий %d пирам\n", len(jobs), len(a.config.Peers))

	var wg sync.WaitGroup
	errors := make(chan error, len(a.config.Peers))

	// Отправляем задания каждому пиру
	for i, peer := range a.config.Peers {
		wg.Add(1)
		go func(peerAddr string, jobIndex int) {
			defer wg.Done()

			if err := a.sendJobToPeer(peerAddr, jobs[jobIndex]); err != nil {
				errors <- fmt.Errorf("ошибка отправки задания пиру %s: %v", peerAddr, err)
			}
		}(peer, i)
	}

	// Ждём завершения всех отправок
	wg.Wait()
	close(errors)

	// Проверяем на ошибки
	for err := range errors {
		fmt.Printf("Ошибка: %v\n", err)
	}

	// Выполняем локальные задания
	return a.executeJobsLocally(jobs)
}

// sendJobToPeer отправляет одно задание пиру
func (a *App) sendJobToPeer(peerAddr string, job map[string]interface{}) error {
	conn, err := net.Dial("tcp", peerAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Создаём сообщение job_request
	msg := Message{
		Type:      "job_request",
		From:      a.config.ServerID,
		To:        peerAddr,
		Data:      job,
		Timestamp: time.Now(),
	}

	// Отправляем сообщение
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(msg); err != nil {
		return err
	}

	// Читаем ответ
	var response Message
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&response); err != nil {
		return err
	}

	if response.Type == "job_response" {
		fmt.Printf("Получен ответ от %s: успех=%v\n", peerAddr, response.Data.(map[string]interface{})["success"])
	} else {
		fmt.Printf("Неожиданный ответ от %s: тип=%s\n", peerAddr, response.Type)
	}

	return nil
}

// executeJobsLocally выполняет задания локально
func (a *App) executeJobsLocally(jobs []map[string]interface{}) error {
	fmt.Printf("Выполняем %d заданий локально\n", len(jobs))

	// Выполняем каждое задание локально
	for i, job := range jobs {
		filename := job["filename"].(string)
		startLine := int(job["start_line"].(float64))
		endLine := int(job["end_line"].(float64))
		pattern := job["pattern"].(string)

		fmt.Printf("Задание %d: файл=%s, строки=%d-%d, паттерн=%s\n",
			i+1, filename, startLine, endLine, pattern)

		// Выполняем grep на указанном диапазоне строк
		if err := a.executeJobLocally(job); err != nil {
			fmt.Printf("Ошибка выполнения задания %d: %v\n", i+1, err)
			return err
		}
	}

	return nil
}

// executeJobLocally выполняет одно задание локально
func (a *App) executeJobLocally(job map[string]interface{}) error {
	filename := job["filename"].(string)
	startLine := int(job["start_line"].(float64))
	endLine := int(job["end_line"].(float64))

	// Открываем файл
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentLine := 0
	matches := 0

	for scanner.Scan() {
		currentLine++

		// Пропускаем строки до startLine
		if currentLine < startLine {
			continue
		}

		// Останавливаемся после endLine
		if currentLine > endLine {
			break
		}

		line := scanner.Text()

		// Выполняем поиск в строке
		result, found := a.searchInLineForJob(line, currentLine, job)
		if found {
			matches++

			// Выводим результат
			if len(a.config.Files) > 1 {
				fmt.Fprintf(a.config.Output, "%s:", filename)
			}
			a.printResultForJob(result)
		}
	}

	if a.config.Flags.Count {
		if len(a.config.Files) > 1 {
			fmt.Fprintf(a.config.Output, "%s:", filename)
		}
		fmt.Fprintf(a.config.Output, "%d\n", matches)
	}

	return scanner.Err()
}

// searchInLineForJob выполняет поиск в строке для конкретного задания
func (a *App) searchInLineForJob(line string, lineNumber int, job map[string]interface{}) (*GrepResult, bool) {
	pattern, _ := job["pattern"].(string)
	flags, _ := job["flags"].(model.GrepFlags)

	// Создаём временный конфиг для поиска
	tempConfig := &Config{
		Pattern: pattern,
		Flags:   flags,
		Output:  a.config.Output,
	}

	// Используем существующую логику поиска
	result, found := a.searchInLineWithConfig(line, lineNumber, tempConfig)
	return result, found
}

// searchInLineWithConfig выполняет поиск с заданным конфигом
func (a *App) searchInLineWithConfig(line string, lineNumber int, config *Config) (*GrepResult, bool) {
	pattern := regexp.QuoteMeta(config.Pattern)
	if config.Flags.WholeLine {
		pattern = "^" + pattern + "$"
	}

	flags := ""
	if config.Flags.IgnoreCase {
		flags = "(?i)"
	}

	fullPattern := flags + pattern
	re, err := regexp.Compile(fullPattern)
	if err != nil {
		return nil, false
	}

	matches := re.FindAllString(line, -1)

	if config.Flags.InvertMatch {
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
		if config.Flags.OnlyMatching {
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

// printResultForJob выводит результат для задания
func (a *App) printResultForJob(result *GrepResult) {
	if a.config.Flags.LineNumber {
		fmt.Fprintf(a.config.Output, "%d:", result.LineNumber)
	}

	if a.config.Flags.OnlyMatching {
		fmt.Fprintf(a.config.Output, "%s\n", result.Match)
	} else {
		fmt.Fprintf(a.config.Output, "%s\n", result.Line)
	}
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
