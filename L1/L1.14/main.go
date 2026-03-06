package main

import (
	"fmt"
	"reflect"
)

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

		default:
			t := reflect.TypeOf(v)
			if t.Kind() == reflect.Chan {
				fmt.Printf("Value (channel memory address): %v, type: %v %v\n",
					v, t.Kind(), t.Elem())
			} else {
				fmt.Printf("Value: %v, unknown variable type: %T\n", v, v)
			}
		}
	}
}
