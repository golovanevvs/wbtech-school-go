package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	fmt.Println(err)
	//lint:ignore SA4023 this is expected behavior in this example
	fmt.Println(err == nil)
}
