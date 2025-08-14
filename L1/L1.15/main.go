package main

import (
	"fmt"
	"runtime"
	"strings"
)

var justString string
var justStringNew string

func someFunc() {
	v := createHugeString(1 << 20)
	justString = v[:100]
}
func someFuncNew() {
	v := createHugeString(1 << 20)
	justStringNew = strings.Clone(v[:100])
}

func createHugeString(n int) string {
	res := make([]rune, n)
	for i := range n {
		res[i] = 'n'
	}
	return string(res)
}

func printMem(name string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("%s: Alloc = %v KB\n",
		name, m.Alloc/1024)
}

func main() {
	printMem("someFunc before")
	someFunc()
	printMem("someFunc after")

	fmt.Println("")

	printMem("someFuncNew before")
	someFuncNew()
	printMem("someFuncNew after")
}
