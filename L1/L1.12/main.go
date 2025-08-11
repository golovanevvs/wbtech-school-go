package main

import (
	"fmt"
)

func main() {
	arr := []string{"cat", "cat", "dog", "cat", "tree"}

	mapSet := make(map[string]bool)
	sets := make([]string, 0)

	for _, v := range arr {
		mapSet[v] = true
	}

	for v := range mapSet {
		sets = append(sets, v)
	}

	fmt.Printf("Sequence: %v\n", arr)
	fmt.Printf("Sets: %v", sets)
}
