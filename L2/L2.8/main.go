package main

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	t, err := getTime()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("Current time: %v\n", t)
}

func getTime() (time.Time, error) {
	t, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}
