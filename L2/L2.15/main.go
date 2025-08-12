package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			fmt.Println("\nexit")
			return
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) > 1 {
			runPipeline(ctx, parts)
			continue
		}

		args := strings.Fields(line)
		runCommand(ctx, args, os.Stdout, os.Stderr)
	}
}

func runPipeline(ctx context.Context, cmds []string) {
	var commands []*exec.Cmd
	for _, c := range cmds {
		args := strings.Fields(strings.TrimSpace(c))
		if len(args) == 0 {
			return
		}
		cmd := exec.CommandContext(ctx, args[0], args[1:]...)
		commands = append(commands, cmd)
	}

	for i := 0; i < len(commands)-1; i++ {
		r, w := io.Pipe()
		commands[i].Stdout = w
		commands[i+1].Stdin = r
	}

	commands[len(commands)-1].Stdout = os.Stdout

	for _, cmd := range commands {
		if err := cmd.Start(); err != nil {
			fmt.Println("error:", err)
			return
		}
	}

	for _, cmd := range commands {
		cmd.Wait()
	}
}

func runCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		return 0
	}

	switch args[0] {
	case "cd":
		return builtin_cd(args)
	case "pwd":
		return builtin_pwd(stdout)
	case "echo":
		return builtin_echo(args, stdout)
	case "kill":
		return builtin_kill(ctx, args)
	case "ps":
		return builtin_ps(ctx, args, stdout)
	case "exit":
		os.Exit(0)
	}

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func builtin_cd(args []string) int {
	if len(args) < 2 {
		fmt.Println("cd: missing operand")
		return 1
	}
	path := args[1]
	if !filepath.IsAbs(path) {
		cwd, _ := os.Getwd()
		path = filepath.Join(cwd, path)
	}
	if err := os.Chdir(path); err != nil {
		fmt.Println("cd:", err)
		return 1
	}
	return 0
}

func builtin_pwd(w io.Writer) int {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(w, "pwd:", err)
		return 1
	}
	fmt.Fprintln(w, dir)
	return 0
}

func builtin_echo(args []string, w io.Writer) int {
	fmt.Fprintln(w, strings.Join(args[1:], " "))
	return 0
}
