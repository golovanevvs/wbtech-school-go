package main

import "fmt"

func main() {
	data := []float64{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5}

	mapGroupValue := make(map[int][]float64)

	for _, t := range data {
		mapGroupValue[(int(t)/10)*10] = append(mapGroupValue[(int(t)/10)*10], t)

	}

	for i := range mapGroupValue {
		fmt.Printf("%d:%v\n", i, mapGroupValue[i])
	}
}
