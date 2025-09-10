package main

import (
	"fmt"
	"io"
	"os"

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

func main() {
	opts := newOptions()

	rootCmd := &cobra.Command{
		Use:   "vgrep",
		Short: "vgrep - v grep",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(args, opts)
		},
	}

	rootCmd.Flags().IntVar(&opts.afterContextCount, "A", 0, "Print N lines of context after each matching line")
	rootCmd.Flags().IntVar(&opts.beforeContextCount, "B", 0, "Print N lines of context before each matching line")
	rootCmd.Flags().IntVar(&opts.aroundContextCount, "C", 0, "Print N lines of context around each matching line")
	rootCmd.Flags().BoolVar(&opts.onlyCount, "c", false, "Print only the count of matching lines")
	rootCmd.Flags().BoolVar(&opts.ignoreCase, "i", false, "Ignore case")
	rootCmd.Flags().BoolVar(&opts.invertMatch, "v", false, "Invert the match")
	rootCmd.Flags().BoolVar(&opts.fixedString, "F", false, "Treat the pattern as a fixed string")
	rootCmd.Flags().BoolVar(&opts.lineNumber, "n", false, "Print line numbers")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

}

func newOptions() *options {
	return &options{}
}

func run(args []string, opts *options) error {
	var reader io.Reader = os.Stdin
	if len(args) == 0 {
		return fmt.Errorf("no search pattern specified")
	}

	pattern := args[0]

	if len(args) > 1 {
		f, err := os.Open(args[1])
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer f.Close()
		reader = f
	}

	return nil
}
