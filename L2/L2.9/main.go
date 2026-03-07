package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func main() {
	fmt.Println("Enter the string:")
	in := bufio.NewReader(os.Stdin)
	str, err := in.ReadString('\n')
	if err != nil {
		fmt.Printf("Input error: %v\n", err)
		os.Exit(1)
	}

	str = strings.TrimSpace(str)

	unpackingString, err := unpackString(str)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Unpacking string: %s\n", unpackingString)
}

func unpackString(str string) (string, error) {
	rs := []rune(str)

	var res strings.Builder
	var r rune

	for i := 0; i < len(rs); i++ {
		switch {
		case rs[i] == '\\':
			if i == len(rs)-1 {
				return "", errors.New("invalid string")
			}
			if r != 0 {
				res.WriteRune(r)
			}
			i++
			r = rs[i]

		case !unicode.IsDigit(rs[i]):
			if r != 0 {
				res.WriteRune(r)
			}
			r = rs[i]

		case unicode.IsDigit(rs[i]):
			if r == 0 {
				return "", errors.New("invalid string")
			}

			numStr := string(rs[i])
			j := i + 1
			for j < len(rs) && unicode.IsDigit(rs[j]) {
				numStr += string(rs[j])
				j++
			}

			count := 0
			for _, d := range numStr {
				count = count*10 + int(d-'0')
			}

			for range count {
				res.WriteRune(r)
			}

			r = 0
			i = j - 1
		}
	}

	if r != 0 {
		res.WriteRune(r)
	}

	return res.String(), nil
}
