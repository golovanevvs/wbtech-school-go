package main

import (
	"os"

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
			return run(args, opts)
		},
	}

	rootCmd.Flags().StringVarP(&opts.fields, "fields", "f", "", "Field (column) numbers to output")
	rootCmd.Flags().StringVarP(&opts.delimiter, "delimiter", "d", "", "Delimiter")
	rootCmd.Flags().BoolVarP(&opts.separated, "separated", "s", false, "Only rows containing the delimiter")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(args []string, opts *options) error {
	return nil
}
