package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Enter a string:")

	in := bufio.NewReader(os.Stdin)
	startStringWithSuffix, err := in.ReadString('\n')
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	startString := strings.TrimRight(startStringWithSuffix, "\r\n")

	bytes := []byte(startString)

	// string reversal
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}

	// word reversals
	start := 0
	for k := 0; k <= len(bytes); k++ {
		if k == len(bytes) || bytes[k] == ' ' {
			for i, j := start, k-1; i < j; i, j = i+1, j-1 {
				bytes[i], bytes[j] = bytes[j], bytes[i]
			}
			start = k + 1
		}
	}

	reversedString := string(bytes)

	fmt.Printf("Reversed string: %s\n", reversedString)
}
