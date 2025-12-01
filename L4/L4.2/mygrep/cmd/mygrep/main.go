package main

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

func main() {
	config, err := parseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка парсинга флагов: %v\n", err)
		os.Exit(1)
	}

	if err := runGrep(config); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка выполнения grep: %v\n", err)
		os.Exit(1)
	}
}

// parseFlags парсит аргументы командной строки
func parseFlags() (*Config, error) {
	// Создаём набор флагов
	var pattern string
	flag.StringVar(&pattern, "pattern", "", "Паттерн для поиска")
	flag.StringVar(&pattern, "e", "", "Паттерн для поиска (alias для --pattern)")

	// Основные флаги grep
	color := flag.Bool("color", false, "Выделить совпадения цветом")
	invertMatch := flag.Bool("v", false, "Инвертировать совпадения")
	ignoreCase := flag.Bool("i", false, "Игнорировать регистр")
	wholeLine := flag.Bool("x", false, "Искать только полные строки")
	lineNumber := flag.Bool("n", false, "Показать номера строк")
	count := flag.Bool("c", false, "Показать только количество совпадений")
	onlyMatching := flag.Bool("o", false, "Показать только совпадающие части")

	flag.Parse()

	if pattern == "" {
		// Если паттерн не указан, берём первый не-flag аргумент
		args := flag.Args()
		if len(args) == 0 {
			return nil, fmt.Errorf("не указан паттерн для поиска")
		}
		pattern = args[0]
		args = args[1:]

		// Оставшиеся аргументы - это файлы
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
			Input:  os.Stdin,
			Output: os.Stdout,
		}
		return config, nil
	}

	// Если паттерн указан через флаг, все остальные аргументы - файлы
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
		Input:  os.Stdin,
		Output: os.Stdout,
	}
	return config, nil
}

// runGrep выполняет поиск
func runGrep(config *Config) error {
	// Если указаны файлы, обрабатываем каждый
	if len(config.Files) > 0 {
		for _, filename := range config.Files {
			if err := processFile(filename, config); err != nil {
				return fmt.Errorf("ошибка обработки файла %s: %v", filename, err)
			}
		}
	} else {
		// Читаем из stdin
		if err := processStream(os.Stdin, "stdin", config); err != nil {
			return err
		}
	}
	return nil
}

// processFile обрабатывает один файл
func processFile(filename string, config *Config) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("не удалось открыть файл %s: %v", filename, err)
	}
	defer file.Close()

	return processStream(file, filename, config)
}

// processStream обрабатывает поток данных
func processStream(reader io.Reader, sourceName string, config *Config) error {
	scanner := bufio.NewScanner(reader)
	lineNumber := 0
	matches := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		result, found := searchInLine(line, lineNumber, config)
		if found {
			matches++
			if !config.Flags.Count {
				if len(config.Files) > 1 {
					// Если несколько файлов, показываем имя файла
					fmt.Fprintf(config.Output, "%s:", sourceName)
				}
				printResult(result, config)
			}
		}
	}

	if config.Flags.Count {
		if len(config.Files) > 1 {
			fmt.Fprintf(config.Output, "%s:", sourceName)
		}
		fmt.Fprintf(config.Output, "%d\n", matches)
	}

	return nil
}

// searchInLine ищет совпадения в строке
func searchInLine(line string, lineNumber int, config *Config) (*GrepResult, bool) {
	// Создаём регулярное выражение
	pattern := regexp.QuoteMeta(config.Pattern)
	if config.Flags.WholeLine {
		pattern = "^" + pattern + "$"
	}

	// Настраиваем флаги регулярного выражения
	flags := ""
	if config.Flags.IgnoreCase {
		flags = "(?i)"
	}

	fullPattern := flags + pattern
	re, err := regexp.Compile(fullPattern)
	if err != nil {
		return nil, false
	}

	// Ищем совпадения
	matches := re.FindAllString(line, -1)

	if config.Flags.InvertMatch {
		// Инвертируем результат
		if len(matches) == 0 {
			return &GrepResult{
				LineNumber: lineNumber,
				Line:       line,
				Match:      line,
			}, true
		}
		return nil, false
	}

	// Обычный поиск
	if len(matches) > 0 {
		if config.Flags.OnlyMatching {
			// Возвращаем только совпадающие части
			for _, match := range matches {
				result := &GrepResult{
					LineNumber: lineNumber,
					Line:       match,
					Match:      match,
				}
				return result, true
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
func printResult(result *GrepResult, config *Config) {
	if config.Flags.LineNumber {
		fmt.Fprintf(config.Output, "%d:", result.LineNumber)
	}

	if config.Flags.OnlyMatching {
		fmt.Fprintf(config.Output, "%s\n", result.Match)
	} else {
		fmt.Fprintf(config.Output, "%s\n", result.Line)
	}
}
