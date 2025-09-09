package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	input := []string{"пятка", "пятак", "тяпка", "листок", "слиток", "столик", "стол"}
	res := getAnagrams(input)
	fmt.Printf("Input: %v\n", input)
	fmt.Println("Result:")
	for k, v := range res {
		fmt.Printf("%q: %v\n", k, v)
	}
}

func normalize(s string) string {
	runes := []rune(strings.ToLower(s))
	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})
	return string(runes)
}

func getAnagrams(words []string) map[string][]string {
	groups := make(map[string][]string)
	firstWord := make(map[string]string)

	for _, w := range words {
		key := normalize(w)
		groups[key] = append(groups[key], strings.ToLower(w))
		if _, ok := firstWord[key]; !ok {
			firstWord[key] = strings.ToLower(w)
		}
	}

	result := make(map[string][]string)

	for key, group := range groups {
		if len(group) < 2 {
			continue
		}
		sort.Strings(group)
		result[firstWord[key]] = group
	}

	return result
}
