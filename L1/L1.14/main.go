package main

import "fmt"

func main() {
	data := []any{1, "string", true, make(chan int), 1.1}

	for _, v := range data {
		switch v.(type) {
		case int:
			fmt.Printf("Value: %d, variable type: %T\n", v, v)

		case string:
			fmt.Printf("Value: %s, variable type: %T\n", v, v)

		case bool:
			fmt.Printf("Value: %v, variable type: %T\n", v, v)

		case chan int:
			fmt.Printf("Value: make(chan int), variable type: %T\n", v)

		default:
			fmt.Printf("Value: %v, unknown variable type\n", v)
		}
	}
}
