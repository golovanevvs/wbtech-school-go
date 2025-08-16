package main

import (
	"fmt"
	"unicode"
)

func main() {
	fmt.Println("Enter the string:")
	var str string
	fmt.Scan(&str)

	var latinMask uint32
	var digitMask uint16
	other := make(map[rune]bool)
	unique := true

label:
	for _, r := range str {
		rLower := unicode.ToLower(r)
		switch {
		case rLower >= 'a' && rLower <= 'z':
			bit := rLower - 'a'
			if latinMask&(1<<bit) != 0 {
				unique = false
				break label
			}
			latinMask |= 1 << bit

		case rLower >= '0' && rLower <= '9':
			bit := rLower - '0'
			if digitMask&(1<<bit) != 0 {
				unique = false
				break label
			}
			digitMask |= 1 << bit

		default:
			if other[rLower] {
				unique = false
				break label
			}
			other[rLower] = true
		}
	}

	var result string
	if unique {
		result = "is unique"
	} else {
		result = "is not unique"
	}

	fmt.Printf("Result: %s\n", result)
}
