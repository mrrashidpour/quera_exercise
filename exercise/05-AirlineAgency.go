package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	n := parseInt(scanner.Text())

	countryCodes := make(map[string]string)

	for i := 0; i < n; i++ {
		scanner.Scan()
		line := scanner.Text()
		parts := strings.Split(line, " ")

		countryName := parts[0]
		countryCode := parts[1]

		countryCodes[countryCode] = countryName
	}

	scanner.Scan()
	q := parseInt(scanner.Text())

	for i := 0; i < q; i++ {
		scanner.Scan()
		phoneNumber := scanner.Text()

		if len(phoneNumber) < 4 {
			fmt.Println("Invalid Number")
			continue
		}

		countryCode := phoneNumber[0:3]

		if countryName, exists := countryCodes[countryCode]; exists {
			fmt.Println(countryName)
		} else {
			fmt.Println("Invalid Number")
		}
	}
}

func parseInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}
