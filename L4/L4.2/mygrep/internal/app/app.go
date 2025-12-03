package app

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.2/mygrep/internal/model"
)

// Config содержит конфигурацию для локального режима
type Config struct {
	Pattern string
	Files   []string
	Flags   model.GrepFlags
	Input   io.Reader
	Output  io.Writer
}

// GrepResult результат поиска
type GrepResult struct {
	LineNumber int
	Line       string
	Match      string
}

// App представляет основное приложение
type App struct {
	config *Config
}

// NewApp создаёт новый экземпляр приложения
func NewApp() *App {
	return &App{}
}

// Run запускает приложение
func (a *App) Run() error {
	var err error
	a.config, err = a.parseFlags()
	if err != nil {
		return fmt.Errorf("ошибка парсинга флагов: %v", err)
	}

	if err := a.runGrep(); err != nil {
		return fmt.Errorf("ошибка выполнения grep: %v", err)
	}

	return nil
}

// parseFlags парсит аргументы командной строки
func (a *App) parseFlags() (*Config, error) {
	var pattern string
	flag.StringVar(&pattern, "pattern", "", "Паттерн для поиска")
	flag.StringVar(&pattern, "e", "", "Паттерн для поиска (alias для --pattern)")

	invertMatch := flag.Bool("v", false, "Инвертировать совпадения")
	ignoreCase := flag.Bool("i", false, "Игнорировать регистр")
	wholeLine := flag.Bool("x", false, "Искать только полные строки")
	lineNumber := flag.Bool("n", false, "Показать номера строк")
	count := flag.Bool("c", false, "Показать только количество совпадений")
	onlyMatching := flag.Bool("o", false, "Показать только совпадающие части")

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
				InvertMatch:  *invertMatch,
				IgnoreCase:   *ignoreCase,
				WholeLine:    *wholeLine,
				LineNumber:   *lineNumber,
				Count:        *count,
				OnlyMatching: *onlyMatching,
			},
			Input:  os.Stdin,
			Output: os.Stdout,
		}
		return config, nil
	}

	config := &Config{
		Pattern: pattern,
		Files:   flag.Args(),
		Flags: model.GrepFlags{
			InvertMatch:  *invertMatch,
			IgnoreCase:   *ignoreCase,
			WholeLine:    *wholeLine,
			LineNumber:   *lineNumber,
			Count:        *count,
			OnlyMatching: *onlyMatching,
		},
		Input:  os.Stdin,
		Output: os.Stdout,
	}
	return config, nil
}

// runGrep выполняет поиск
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
