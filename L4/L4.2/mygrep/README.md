# mygcs
# mygrep - Distributed Grep Utility

A high-performance distributed version of the classic `grep` utility with quorum-based fault tolerance and concurrent processing.

## Features

- **Full grep compatibility** - Supports all standard grep flags (`-n`, `-v`, `-i`, `-x`, `-c`, `-o`, `--color`)
- **Distributed processing** - Automatically splits files across multiple servers
- **Quorum-based fault tolerance** - Continues operation even when some servers fail
- **Concurrent processing** - Uses Go goroutines for parallel execution
- **TCP networking** - JSON-based message passing between servers

## Quick Start

### Local Mode (like regular grep)

```bash
# Basic search
./mygrep "pattern" file.txt

# With line numbers
./mygrep -n "pattern" file.txt

# Inverted search
./mygrep -v "pattern" file.txt

# Case insensitive
./mygrep -i "PATTERN" file.txt

# Count matches only
./mygrep -c "pattern" file.txt

# Multiple files
./mygrep -n "pattern" file1.txt file2.txt file3.txt
```

### Distributed Mode

```bash
# Single server (fallback to local mode)
./mygrep -n "pattern" largefile.txt -port=8080

# Two servers
Terminal 1: ./mygrep -n "test" largefile.txt -port=8080 -peers="localhost:8081"
Terminal 2: ./mygrep -n "test" largefile.txt -port=8081 -peers="localhost:8080"

# Three servers (recommended for quorum)
Terminal 1: ./mygrep -n "test" largefile.txt -port=8080 -peers="localhost:8081,localhost:8082"
Terminal 2: ./mygrep -n "test" largefile.txt -port=8081 -peers="localhost:8080,localhost:8082"
Terminal 3: ./mygrep -n "test" largefile.txt -port=8082 -peers="localhost:8080,localhost:8081"
```

## Command Line Options

### Standard grep flags
- `-n, --line-number` - Show line numbers
- `-v, --invert-match` - Show lines NOT matching pattern
- `-i, --ignore-case` - Case insensitive search
- `-x, --line-regexp` - Match whole lines only
- `-c, --count` - Show only match count
- `-o, --only-matching` - Show only matching parts
- `--color` - Highlight matches with color

### Distributed flags
- `-port PORT` - Port for TCP server (enables distributed mode)
- `-peers PEERS` - Comma-separated list of peers (e.g., "localhost:8081,localhost:8082")
- `-server-id ID` - Server identifier (default: hostname:port)

## Building

```bash
# Build the application
go build ./cmd/mygrep

# Run tests
go test ./...

# Install
go install ./cmd/mygrep
```

## How It Works

### File Distribution
1. Input file is split into chunks based on line count
2. Each server receives a chunk with start/end line numbers
3. Servers process their chunks independently in parallel
4. Results are collected and combined

### Quorum System
- System waits for responses from N/2+1 servers (where N is total servers)
- For 3 servers: waits for 2+ responses
- For 5 servers: waits for 3+ responses
- Continues operation even if some servers fail
- 30-second timeout for slow responses

## Documentation

- [Usage Guide](USAGE.md) - Detailed usage instructions
- [Testing Guide](TESTING.md) - Test scenarios and benchmarks
- [Implementation Plan](plan.md) - Development roadmap and architecture