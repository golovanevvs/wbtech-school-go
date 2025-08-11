package main

import "fmt"

func main() {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Println("Original array:")
	fmt.Println(arr)

	inCh := generator(arr...)
	outCh := multiply(inCh)

	fmt.Println("Result:")
	for res := range outCh {
		fmt.Print(res, " ")
	}

}
func generator(arr ...int) chan int {
	inCh := make(chan int)

	go func() {
		defer close(inCh)
		for _, data := range arr {
			inCh <- data
		}
	}()

	return inCh
}

func multiply(inCh <-chan int) chan int {
	outCh := make(chan int)

	go func() {
		defer close(outCh)
		for data := range inCh {
			outCh <- data * 2
		}
	}()

	return outCh
}
