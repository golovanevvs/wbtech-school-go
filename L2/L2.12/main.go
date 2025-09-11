package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

type options struct {
	afterContextCount  int  // -A N
	beforeContextCount int  // -B N
	aroundContextCount int  // -C N
	onlyCount          bool // -c
	ignoreCase         bool // -i
	invertMatch        bool // -v
	fixedString        bool // -F
	lineNumber         bool // -n
}

type line struct {
	n int
	t string
}

func main() {
	opts := &options{}

	rootCmd := &cobra.Command{
		Use:   "vgrep PATTERN [FILE]",
		Short: "vgrep - grep-like utility",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(args, opts)
		},
	}

	rootCmd.Flags().IntVarP(&opts.afterContextCount, "after-context", "A", 0, "Print N lines after match")
	rootCmd.Flags().IntVarP(&opts.beforeContextCount, "before-context", "B", 0, "Print N lines before match")
	rootCmd.Flags().IntVarP(&opts.aroundContextCount, "context", "C", 0, "Print N lines before/after match")
	rootCmd.Flags().BoolVarP(&opts.onlyCount, "count", "c", false, "Print count of matching lines")
	rootCmd.Flags().BoolVarP(&opts.ignoreCase, "ignore-case", "i", false, "Ignore case")
	rootCmd.Flags().BoolVarP(&opts.invertMatch, "invert-match", "v", false, "Invert match")
	rootCmd.Flags().BoolVarP(&opts.fixedString, "fixed-strings", "F", false, "Match as fixed string")
	rootCmd.Flags().BoolVarP(&opts.lineNumber, "line-number", "n", false, "Show line numbers")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(args []string, opts *options) error {
	if len(args) == 0 {
		return fmt.Errorf("no search pattern specified")
	}
	pattern := args[0]

	if opts.aroundContextCount > 0 {
		opts.afterContextCount = opts.aroundContextCount
		opts.beforeContextCount = opts.aroundContextCount
	}

	var reader io.Reader = os.Stdin
	if len(args) > 1 {
		f, err := os.Open(args[1])
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer f.Close()
		reader = f
	}

	var lines []line
	s := bufio.NewScanner(reader)
	n := 1
	for s.Scan() {
		lines = append(lines, line{n, s.Text()})
		n++
	}
	if err := s.Err(); err != nil {
		return fmt.Errorf("scan error: %w", err)
	}

	var re *regexp.Regexp
	var err error
	if !opts.fixedString {
		p := pattern
		if opts.ignoreCase {
			p = "(?i)" + p
		}
		// \| -> |
		p = strings.ReplaceAll(p, `\|`, `|`)
		re, err = regexp.Compile(p)
		if err != nil {
			return fmt.Errorf("invalid regexp: %w", err)
		}
	}
	if opts.fixedString && opts.ignoreCase {
		pattern = strings.ToLower(pattern)
	}

	// match searching
	var matchIdxs []int
	for i, l := range lines {
		text := l.t
		if opts.fixedString && opts.ignoreCase {
			text = strings.ToLower(text)
		}
		isMatch := opts.fixedString && strings.Contains(text, pattern) || (!opts.fixedString && re.MatchString(l.t))
		if opts.invertMatch {
			isMatch = !isMatch
		}
		if isMatch {
			matchIdxs = append(matchIdxs, i)
		}
	}

	if opts.onlyCount {
		fmt.Println(len(matchIdxs))
		return nil
	}

	// Union of overlapping contexts
	type interval struct{ start, end int }
	var intervals []interval

	for _, idx := range matchIdxs {
		start := idx - opts.beforeContextCount
		if start < 0 {
			start = 0
		}
		end := idx + opts.afterContextCount
		if end >= len(lines) {
			end = len(lines) - 1
		}

		if len(intervals) > 0 && start <= intervals[len(intervals)-1].end+1 {
			if end > intervals[len(intervals)-1].end {
				intervals[len(intervals)-1].end = end
			}
		} else {
			intervals = append(intervals, interval{start, end})
		}
	}

	// printing
	for i, iv := range intervals {
		if i > 0 {
			fmt.Println("--")
		}
		for j := iv.start; j <= iv.end; j++ {
			prefix := ""
			if opts.lineNumber {
				prefix = fmt.Sprintf("%d:", lines[j].n)
			}
			fmt.Printf("%s%s\n", prefix, lines[j].t)
		}
	}

	return nil
}
