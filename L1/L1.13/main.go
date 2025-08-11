package main

import "fmt"

func main() {
	fmt.Println("Tuple assignment:")
	a := 1
	b := 2
	fmt.Printf("a=%d, b=%d\n", a, b)
	a, b = b, a
	fmt.Printf("a=%d, b=%d\n", a, b)

	fmt.Println("Arithmetic swap:")
	a = 1
	b = 2
	fmt.Printf("a=%d, b=%d\n", a, b)
	a = a + b
	b = a - b
	a = a - b
	fmt.Printf("a=%d, b=%d\n", a, b)

	fmt.Println("XOR swap:")
	a = 1
	b = 2
	fmt.Printf("a=%d, b=%d\n", a, b)
	a = a ^ b
	b = b ^ a
	a = a ^ b
	fmt.Printf("a=%d, b=%d\n", a, b)
}
