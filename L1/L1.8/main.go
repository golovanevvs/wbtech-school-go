package main

import (
	"fmt"
)

func main() {
	var number int64
	fmt.Println("Enter a number:")
	fmt.Scan(&number)

	numberStr := fmt.Sprintf("%b", number)
	fmt.Printf("The number %d in binary is: %s\n", number, numberStr)

	var idx int
	fmt.Println("Enter the index of the i-th bit to set (1 = LSB):")
	fmt.Scan(&idx)

	if idx < 1 || idx > 64 {
		fmt.Printf("Invalid bit index: %d. Must be between 1 and 64\n", idx)
		return
	}

	bitPos := idx - 1

	var bitValue int
	fmt.Println("Enter the bit value (0 or 1):")
	fmt.Scan(&bitValue)

	if bitValue != 0 && bitValue != 1 {
		fmt.Println("Invalid bit value")
		return
	}

	mask := int64(1) << bitPos
	var result int64
	if bitValue == 1 {
		result = number | mask
	} else {
		result = number &^ mask
	}

	fmt.Printf("Modified number: %d (%b)\n", result, result)
}
