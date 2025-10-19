package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())

	results := make([]string, n)

	for i := 0; i < n; i++ {
		scanner.Scan()
		teacherName := scanner.Text()

		scanner.Scan()
		percentagesStr := scanner.Text()

		percentages := strings.Split(percentagesStr, " ")
		sum := 0
		count := 0

		for _, p := range percentages {
			if p == "" {
				continue
			}
			percent, _ := strconv.Atoi(p)
			sum += percent
			count++
		}

		var average float64
		if count > 0 {
			average = float64(sum) / float64(count)
		} else {
			average = 0
		}

		var rating string
		switch {
		case average >= 80:
			rating = "Excellent"
		case average >= 60:
			rating = "Very Good"
		case average >= 40:
			rating = "Good"
		default:
			rating = "Fair"
		}

		results[i] = fmt.Sprintf("%s %s", teacherName, rating)
	}

	for _, result := range results {
		fmt.Println(result)
	}
}
