//go:build windows

package main

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strconv"
)

func builtin_ps(ctx context.Context, _ []string, w io.Writer) int {
	cmd := exec.CommandContext(ctx, "tasklist")
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintln(w, "ps:", err)
		return 1
	}
	w.Write(out)
	return 0
}

func builtin_kill(_ context.Context, args []string) int {
	if len(args) < 2 {
		fmt.Println("kill: missing pid")
		return 1
	}
	pid := args[1]
	if _, err := strconv.Atoi(pid); err != nil {
		fmt.Println("kill: invalid pid")
		return 1
	}
	cmd := exec.Command("taskkill", "/PID", pid, "/F")
	if err := cmd.Run(); err != nil {
		fmt.Println("kill:", err)
		return 1
	}
	return 0
}
