package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type options struct {
	fields    string // -f
	delimiter string // -d
	separated bool   // -s
}

func main() {
	opts := &options{}

	rootCmd := &cobra.Command{
		Use:   "vcut",
		Short: "vcut - cut-like utility",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(opts)
		},
	}

	rootCmd.Flags().StringVarP(&opts.fields, "fields", "f", "", "Field (column) numbers to output")
	rootCmd.Flags().StringVarP(&opts.delimiter, "delimiter", "d", "", "Delimiter")
	rootCmd.Flags().BoolVarP(&opts.separated, "separated", "s", false, "Only rows containing the delimiter")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(opts *options) error {
	delimiter := "\t"
	if opts.delimiter != "" {
		delimiter = opts.delimiter
	}

	fieldSet, err := parseFields(opts.fields)
	if err != nil {
		return fmt.Errorf("invalid fields: %w", err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if opts.separated && !strings.Contains(line, delimiter) {
			continue
		}

		cols := strings.Split(line, delimiter)

		var out []string

		for idx := range fieldSet {
			if idx >= 0 && idx < len(cols) {
				out = append(out, cols[idx])
			}
		}

		if len(out) > 0 {
			fmt.Println(strings.Join(out, delimiter))
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan error: %w", err)
	}

	return nil
}

func parseFields(fields string) (map[int]struct{}, error) {
	result := make(map[int]struct{})

	if fields == "" {
		return result, nil
	}

	parts := strings.Split(fields, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.Contains(part, "-") {
			bounds := strings.Split(part, "-")
			if len(bounds) != 2 {
				return nil, fmt.Errorf("bad range: %s", part)
			}

			start, err1 := strconv.Atoi(bounds[0])
			end, err2 := strconv.Atoi(bounds[1])
			if err1 != nil || err2 != nil || start < 1 || end < 1 || start > end {
				return nil, fmt.Errorf("bad range: %s", part)
			}

			for i := start; i <= end; i++ {
				result[i-1] = struct{}{}
			}
		} else {
			num, err := strconv.Atoi(part)
			if err != nil || num < 1 {
				return nil, fmt.Errorf("bad field: %s", part)
			}

			result[num-1] = struct{}{}
		}
	}

	return result, nil
}
