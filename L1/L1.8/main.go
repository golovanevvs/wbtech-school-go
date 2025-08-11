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
	fmt.Println("Enter the index of the i-th bit to set (0 = LSB, negative from LSB):")
	fmt.Scan(&idx)

	bitLen := len(numberStr)

	if idx < 0 {
		idx = bitLen + idx
	}

	if idx < 0 || idx >= bitLen {
		fmt.Printf("There is no such bit index in the binary representation of the number %d\n", number)
		return
	}

	var bitValue int
	fmt.Println("Enter the bit value (0 or 1):")
	fmt.Scan(&bitValue)

	if bitValue != 0 && bitValue != 1 {
		fmt.Println("Invalid bit value")
		return
	}

	mask := int64(1) << idx
	var result int64
	if bitValue == 1 {
		result = number | mask
	} else {
		result = number &^ mask
	}

	fmt.Printf("Modified number: %d (%b)\n", result, result)
}
