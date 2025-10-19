package main

import (
	"fmt"
	"strconv"
	"unicode"
)

func main() {
	var s string
	fmt.Scan(&s)

	sum := extractAndSumNumbers(s)

	if isArmstrong(sum) {
		fmt.Println("YES")
	} else {
		fmt.Println("NO")
	}
}

func extractAndSumNumbers(s string) int {
	sum := 0
	numStr := ""

	for _, ch := range s {
		if unicode.IsDigit(ch) {
			numStr += string(ch)
		} else {
			if numStr != "" {
				num, _ := strconv.Atoi(numStr)
				sum += num
				numStr = ""
			}
		}
	}

	if numStr != "" {
		num, _ := strconv.Atoi(numStr)
		sum += num
	}

	return sum
}

func isArmstrong(num int) bool {
	if num == 0 {
		return true
	}
	s := strconv.Itoa(num)
	n := len(s)

	total := 0
	for _, ch := range s {
		d := int(ch - '0')
		pow := 1
		for i := 0; i < n; i++ {
			pow *= d
		}
		total += pow
	}

	return total == num
}
