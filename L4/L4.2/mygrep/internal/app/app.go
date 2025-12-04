package app

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.2/mygrep/internal/model"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.2/mygrep/internal/service"
	"github.com/golovanevvs/wbtech-school-go/tree/main/L4/L4.2/mygrep/internal/transport"
)

// App represents the main application
type App struct {
	config      *model.Config
	grepService *service.GrepService
	server      *transport.Server
	client      *transport.Client
}

// NewApp creates a new application instance
func NewApp() *App {
	return &App{}
}

// Run starts the application
func (a *App) Run() error {
	var err error
	a.config, err = a.parseFlags()
	if err != nil {
		return fmt.Errorf("flag parsing error: %v", err)
	}

	// Initialize grep service
	a.grepService = service.NewGrepService(a.config)

	// If distributed mode is enabled, start server
	if a.config.IsDistributed {
		if err := a.startDistributedMode(); err != nil {
			return err
		}
	} else {
		// Local mode
		if err := a.grepService.ExecuteGrep(); err != nil {
			return fmt.Errorf("grep execution error: %v", err)
		}
	}

	return nil
}

// parseFlags parses command line arguments
func (a *App) parseFlags() (*model.Config, error) {
	var pattern string
	flag.StringVar(&pattern, "pattern", "", "Pattern to search for")
	flag.StringVar(&pattern, "e", "", "Pattern to search for (alias for --pattern)")

	// Standard grep flags
	color := flag.Bool("color", false, "Highlight matches with color")
	invertMatch := flag.Bool("v", false, "Invert match (show lines NOT containing pattern)")
	ignoreCase := flag.Bool("i", false, "Ignore case")
	wholeLine := flag.Bool("x", false, "Match whole lines only")
	lineNumber := flag.Bool("n", false, "Show line numbers")
	count := flag.Bool("c", false, "Show only match count")
	onlyMatching := flag.Bool("o", false, "Show only matching parts")

	// Distributed flags
	port := flag.String("port", "", "Port for distributed mode")
	peers := flag.String("peers", "", "Comma-separated list of peers")
	serverID := flag.String("server-id", "", "Server identifier (default: hostname:port)")

	flag.Parse()

	if pattern == "" {
		args := flag.Args()
		if len(args) == 0 {
			return nil, fmt.Errorf("no pattern specified for search")
		}
		pattern = args[0]
		args = args[1:]

		config := &model.Config{
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

	config := &model.Config{
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

// parsePeers parses comma-separated peer list
func parsePeers(peersStr string) []string {
	if peersStr == "" {
		return nil
	}
	return strings.Split(peersStr, ",")
}

// startDistributedMode starts distributed mode operation
func (a *App) startDistributedMode() error {
	// Generate ServerID if not specified
	if a.config.ServerID == "" {
		hostname, _ := os.Hostname()
		a.config.ServerID = fmt.Sprintf("%s:%s", hostname, a.config.Port)
	}

	fmt.Printf("Starting distributed mode. ServerID: %s\n", a.config.ServerID)

	// Start TCP server with job handler adapter
	jobHandler := transport.NewJobHandlerAdapter(a.grepService)
	a.server = transport.NewServer(a.config, jobHandler)
	if err := a.server.Start(); err != nil {
		return fmt.Errorf("TCP server startup error: %v", err)
	}
	defer a.server.Stop()

	// Connect to peers
	if len(a.config.Peers) > 0 {
		a.client = transport.NewClient(a.config)
		if err := a.client.ConnectToPeers(); err != nil {
			fmt.Printf("Warning: failed to connect to peers: %v\n", err)
		}
	}

	// In distributed mode, process files or wait for commands
	if len(a.config.Files) > 0 {
		return a.runDistributedGrep()
	}

	// Wait for commands (for now just show status)
	fmt.Println("Server started. Waiting for commands...")

	// Simple wait mechanism
	quit := make(chan bool)
	<-quit
	return nil
}

// runDistributedGrep runs distributed grep operation
func (a *App) runDistributedGrep() error {
	fmt.Println("Running distributed grep...")

	// If we have peers, use distributed processing
	if len(a.config.Peers) > 0 && a.client != nil {
		// Use client's distributed processing
		return a.runDistributedProcessingWithClient()
	}

	// Fallback to local processing
	return a.grepService.ExecuteDistributedGrep()
}

// runDistributedProcessingWithClient runs distributed processing using client
func (a *App) runDistributedProcessingWithClient() error {
	// Get jobs from grep service
	// For now, this is a simplified version
	// In full implementation, this would coordinate with the transport layer

	fmt.Println("Distributed processing with client not fully implemented yet")
	fmt.Println("Falling back to local processing...")

	return a.grepService.ExecuteDistributedGrep()
}
