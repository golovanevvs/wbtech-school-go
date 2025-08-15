package main

import (
	"fmt"
	"math/rand"
	"sort"
)

func main() {
	n := 20
	m := 100
	data := make([]int, 0, n)
	seen := make(map[int]struct{})
	for len(data) < n {
		v := rand.Intn(m)
		if _, exists := seen[v]; !exists {
			seen[v] = struct{}{}
			data = append(data, v)
		}
	}
	sort.Ints(data)
	fmt.Printf("data: %v\n", data)

	var targetValue int
	fmt.Println("Enter the target value:")
	fmt.Scan(&targetValue)

	index := search(data, targetValue)
	if index == -1 {
		fmt.Printf("Target value %d not found\n", targetValue)
		return
	}

	fmt.Printf("The index of the target value %d is %d\n", targetValue, index)
}

func search(data []int, targetValue int) int {
	left, right := 0, len(data)-1

	for left <= right {
		mid := left + (right-left)/2

		switch {
		case data[mid] == targetValue:
			return mid
		case data[mid] > targetValue:
			right = mid - 1
		default:
			left = mid + 1
		}
	}

	return -1
}
