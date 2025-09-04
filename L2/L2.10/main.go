package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type options struct {
	Column    int
	Numeric   bool
	Reverse   bool
	Unique    bool
	Month     bool
	IgnoreTB  bool
	Check     bool
	Human     bool
	humanMult map[string]float64
}

func main() {
	opts := newOptions()

	rootCmd := &cobra.Command{
		Use:   "vsort [file]",
		Short: "vsort - file content sorter",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(args, opts)
		},
	}

	rootCmd.Flags().IntVarP(&opts.Column, "key", "k", 0, "sort by column number")
	rootCmd.Flags().BoolVarP(&opts.Numeric, "numeric", "n", false, "sort by numeric value")
	rootCmd.Flags().BoolVarP(&opts.Reverse, "reverse", "r", false, "sort in reverse order")
	rootCmd.Flags().BoolVarP(&opts.Unique, "unique", "u", false, "do not output duplicate lines")
	rootCmd.Flags().BoolVarP(&opts.Month, "month", "M", false, "sort by month name")
	rootCmd.Flags().BoolVarP(&opts.IgnoreTB, "ignore-trailing-blanks", "b", false, "ignore trailing blanks")
	rootCmd.Flags().BoolVarP(&opts.Check, "check", "c", false, "check whether input is sorted; do not sort")
	rootCmd.Flags().BoolVarP(&opts.Human, "human-numeric-sort", "h", false, "compare human readable numbers (2K, 1M)")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newOptions() *options {
	return &options{
		humanMult: map[string]float64{
			"K": 1024,
			"M": 1024 * 1024,
			"G": 1024 * 1024 * 1024,
			"T": 1024 * 1024 * 1024 * 1024,
		},
	}
}

func run(args []string, opts *options) error {
	var reader io.Reader = os.Stdin
	if len(args) > 0 {
		f, err := os.Open(args[0])
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer f.Close()
		reader = f
	}

	lines, err := readLines(reader, *opts)
	if err != nil {
		return err
	}

	if opts.Check {

	}

	return nil
}

func readLines(reader io.Reader, opts options) ([]string, error) {
	var lines []string

	s := bufio.NewScanner(reader)

	for s.Scan() {
		line := s.Text()
		if opts.IgnoreTB {
			line = strings.TrimRight(line, "\t")
		}
		lines = append(lines, line)
	}

	return lines, s.Err()
}

func isSorted(lines []string, opts *options) bool {
	return sort.SliceIsSorted(lines, func(i, j int) bool {
		vi := getKeyString(lines[i], opts.Column)
		vj := getKeyString(lines[j], opts.Column)

		switch {
		case opts.Month:
			return opts.monthCompare(vi, vj)
		case opts.Human:
			return opts.humanCompare(vi, vj)
		case opts.Numeric:
		default:
			return vi < vj
		}
	})
}

func getKeyString(s string, column int) string {
	if column <= 0 {
		return s
	}

	fields := strings.Split(s, "\t")
	if column-1 < len(fields) {
		return fields[column-1]
	}

	return ""
}

func (opts *options) monthCompare(a, b string) bool {
	ai, errA := parseMonth(a)
	bi, errB := parseMonth(b)

	switch {
	case errA == nil && errB == nil:
		return ai < bi
	case errA == nil && errB != nil:
		return true
	case errA != nil && errB == nil:
		return false
	default:
		return a < b
	}
}

func parseMonth(s string) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty")
	}

	if t, err := time.Parse("Jan", strings.ToUpper(s[:1])+strings.ToLower(s[1:])); err == nil {
		return int(t.Month()), nil
	}

	return 0, fmt.Errorf("not a month")
}

func (opts *options) humanCompare(a, b string) bool {
	ai, errA := opts.parseHuman(a)
	bi, errB := opts.parseHuman(b)

	switch {
	case errA == nil && errB == nil:
		return ai < bi
	case errA == nil && errB != nil:
		return true
	case errA != nil && errB == nil:
		return false
	default:
		return a < b
	}
}

func (opts *options) parseHuman(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty")
	}

	sign := 1.0
	if s[0] == '+' || s[0] == '-' {
		if s[0] == '-' {
			sign = -1.0
		}
		s = s[1:]
		s = strings.TrimSpace(s)
		if s == "" {
			return 0, fmt.Errorf("no digits")
		}
	}

	numEnd := 0
	for numEnd < len(s) {
		c := s[numEnd]
		if (c >= '0' && c <= '9') || c == '.' {
			numEnd++
			continue
		}
		break
	}
	if numEnd == 0 {
		return 0, fmt.Errorf("no numeric prefix")
	}

	numStr := s[:numEnd]
	rest := strings.TrimSpace(s[numEnd:])

	baseVal, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, fmt.Errorf("bad number: %w", err)
	}

	if rest == "" {
		return sign * baseVal, nil
	}

	rs := strings.ToUpper(rest)

	if strings.HasSuffix(rs, "B") {
		rs = strings.TrimSuffix(rs, "B")
	}

	if len(rs) == 0 {
		return sign * baseVal, nil
	}

	first := rs[0:1]
	mult, ok := opts.humanMult[first]
	if !ok {
		return 0, fmt.Errorf("unknown suffix: %s", rest)
	}

	return sign * baseVal * mult, nil
}
