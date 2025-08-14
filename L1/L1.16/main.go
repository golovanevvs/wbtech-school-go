package main

import (
	"fmt"
	"math/rand"
)

func main() {
	n := 20
	m := 100
	dataNonSorted := make([]int, n)
	for i := range n {
		dataNonSorted[i] = rand.Intn(m)
	}

	dataSorted := quickSort(dataNonSorted)

	fmt.Printf("data non sorted: %v\n", dataNonSorted)
	fmt.Printf("data sorted: %v\n", dataSorted)
}

func quickSort(dataNonSorted []int) []int {
	if len(dataNonSorted) < 2 {
		return dataNonSorted
	}

	pivot := dataNonSorted[0]

	var less, greater []int

	for _, value := range dataNonSorted[1:] {
		if value <= pivot {
			less = append(less, value)
		} else {
			greater = append(greater, value)
		}
	}

	dataSorted := append(quickSort(less), pivot)
	dataSorted = append(dataSorted, quickSort(greater)...)

	return dataSorted
}
