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

	runes := []rune(startString)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	reversedString := string(runes)

	fmt.Printf("Reversed string: %s\n", reversedString)
}
