package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func main() {
	n := 20
	m := 100
	data := make([]int, n)
	for i := range n {
		data[i] = rand.Intn(m)
	}

	fmt.Printf("Slice: %v\n", data)

	in := bufio.NewReader(os.Stdin)

	var index int

loop:
	for {
		fmt.Print("Enter the index of the element to remove from the slice: ")
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

		index = v

		switch {
		case index < 0:
			fmt.Println("The index of the element cannot be negative")
		case index > len(data)-1:
			fmt.Println("The element index exceeds the slice length")
		default:
			break loop
		}
	}

	copy(data[index:], data[index+1:])
	data = data[:len(data)-1]

	fmt.Printf("Result: %v\n", data)
}
