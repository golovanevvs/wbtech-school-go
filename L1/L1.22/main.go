package main

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)

	a := new(big.Int)
	fmt.Print("Enter the first large number (a): ")
	aStr, err := in.ReadString('\n')
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	aStr = strings.TrimSpace(aStr)
	a.SetString(aStr, 10)

	b := new(big.Int)
	zero := big.NewInt(0)
	for {
		fmt.Print("Enter the second large number (b): ")
		bStr, err := in.ReadString('\n')
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}
		bStr = strings.TrimSpace(bStr)
		b.SetString(bStr, 10)
		if b.Cmp(zero) != 0 {
			break
		}
		fmt.Println("[b] cannot be equal to 0!")
	}

	product := new(big.Int).Mul(a, b)
	fmt.Printf("%s * %s = %s\n", a, b, product)

	aFloat := new(big.Float).SetInt(a)
	bFloat := new(big.Float).SetInt(b)
	quotient := new(big.Float).Quo(aFloat, bFloat)
	fmt.Printf("%s / %s = %s\n", a, b, quotient.Text('f', 2))

	sum := new(big.Int).Add(a, b)
	fmt.Printf("%s + %s = %s\n", a, b, sum)

	diff := new(big.Int).Sub(a, b)
	fmt.Printf("%s - %s = %s\n", a, b, diff)
}
