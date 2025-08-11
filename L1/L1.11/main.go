package main

import "fmt"

func main() {
	arrA := []int{1, 2, 3}
	arrB := []int{2, 3, 4}

	mapA := make(map[int]bool)
	intersection := make([]int, 0)

	for _, v := range arrA {
		mapA[v] = true
	}

	for _, v := range arrB {
		if mapA[v] {
			intersection = append(intersection, v)
		}
	}

	fmt.Printf("set A: %v\n", arrA)
	fmt.Printf("set B: %v\n", arrB)
	fmt.Printf("Intersection of sets A, B: %v\n", intersection)
}
