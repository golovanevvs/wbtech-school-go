package main

import (
	"os"

	"github.com/spf13/cobra"
)

type options struct {
	Column   int
	Numeric  bool
	Reverse  bool
	Unique   bool
	Month    bool
	IgnoreBl bool
	Check    bool
	Human    bool
}

func main() {
	opts := &options{}

	rootCmd := &cobra.Command{
		Use:   "gsort [file]",
		Short: "gsort - simplified UNIX sort clone",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(args, opts)
		},
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(args []string, opts *options) error {
	return nil
}
