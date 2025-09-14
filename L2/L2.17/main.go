package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

type options struct {
	timeoutS int // --timeout in seconds
}

func main() {
	opts := &options{}

	rootCmd := &cobra.Command{
		Use:   "vtelnet",
		Short: "v telnet client",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(args[0], opts)
		},
	}

	rootCmd.Flags().IntVar(&opts.timeoutS, "timeout", 10, "Connection timeout in seconds")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(addr string, opts *options) error {
	if _, _, err := net.SplitHostPort(addr); err != nil {
		return fmt.Errorf("invalid host:port format: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(opts.timeoutS)*time.Second)
	defer cancel()

	fmt.Fprintf(os.Stderr, "Connecting to %s with timeout %d seconds...\n", addr, opts.timeoutS)
	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", addr, err)
	}
	fmt.Fprintf(os.Stderr, "Connected to %s\n", addr)

	doneCh := make(chan struct{})

	var once sync.Once
	closeConn := func() {
		once.Do(func() {
			conn.Close()
			close(doneCh)
		})
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Fprintln(os.Stderr, "\nInterrupt received, closing connection...")
		closeConn()
	}()

	// conn -> stdout
	go func() {
		if _, err := io.Copy(os.Stdout, conn); err != nil {
			fmt.Fprintf(os.Stderr, "read error: %v\n", err)
		}

		closeConn()
	}()

	// stdin -> conn
	go func() {
		if _, err := io.Copy(conn, os.Stdin); err != nil {
			fmt.Fprintf(os.Stderr, "write error: %v\n", err)
		}
		closeConn()
	}()

	<-doneCh
	fmt.Fprintln(os.Stderr, "Connection closed")
	return nil
}
