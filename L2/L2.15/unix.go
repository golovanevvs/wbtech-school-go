//go:build !windows

package main

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"syscall"
)

func builtin_ps(ctx context.Context, _ []string, w io.Writer) int {
	cmd := exec.CommandContext(ctx, "ps", "-eo", "pid,ppid,stat,comm")
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
	pid, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("kill: invalid pid")
		return 1
	}
	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		fmt.Println("kill:", err)
		return 1
	}
	return 0
}
