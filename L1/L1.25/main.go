package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func sleep(seconds int) {
	<-time.After(time.Duration(seconds) * time.Second)
}

func main() {
	in := bufio.NewReader(os.Stdin)

	var seconds int

	for {
		fmt.Print("Enter the time in seconds: ")
		str, err := in.ReadString('\n')
		if err != nil {
			fmt.Printf("Read error: %v\n", err)
			continue
		}

		str = strings.TrimSpace(str)

		v, err := strconv.Atoi(str)
		if err != nil {
			fmt.Printf("Value error: %v\n", err)
			continue
		}

		seconds = v

		break
	}

	start := time.Now()
	fmt.Printf("Waiting for %d seconds...", seconds)
	sleep(seconds)
	fmt.Printf("\nPassed: %v", time.Since(start).Seconds())
}
