package main

import (
	"fmt"
)

func main() {
	var s = []string{"1", "2", "3"}
	modifySlice(s)
	fmt.Println(s) // 3,2,3
}

func modifySlice(i []string) {
	i[0] = "3"         // 3,2,3
	i = append(i, "4") // 3,2,3,4
	i[1] = "5"         // 3,5,3,4
	//lint:ignore SA4006 this is expected behavior in this example
	i = append(i, "6") // 3,5,3,4,6
}
