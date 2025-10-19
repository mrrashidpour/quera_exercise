package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// خواندن لیست‌ها
	scanner.Scan()
	coats := parseLine(scanner.Text())
	scanner.Scan()
	shirts := parseLine(scanner.Text())
	scanner.Scan()
	pants := parseLine(scanner.Text())
	scanner.Scan()
	caps := parseLine(scanner.Text())
	scanner.Scan()
	jackets := parseLine(scanner.Text())
	scanner.Scan()
	season := strings.ToUpper(strings.TrimSpace(scanner.Text()))

	// اعمال قوانین فصل
	if season == "SUMMER" {
		coats = []string{}
		jackets = []string{}
	} else if season == "SPRING" {
		jackets = []string{}
	} else if season == "FALL" {
		jackets = []string{}
		allowedCoats := []string{}
		for _, c := range coats {
			lc := strings.ToLower(c)
			if lc != "yellow" && lc != "orange" {
				allowedCoats = append(allowedCoats, c)
			}
		}
		coats = allowedCoats
	}

	// چاپ ترکیب‌ها
	if season == "WINTER" {
		printWinterCombinations(coats, shirts, pants, jackets)
	} else {
		printCombinations(coats, shirts, pants, caps, season)
	}
}

func parseLine(line string) []string {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) < 2 {
		return []string{}
	}
	items := strings.Fields(parts[1])
	return items
}

func printCombinations(coats, shirts, pants, caps []string, season string) {
	for _, coat := range append([]string{""}, coats...) {
		for _, shirt := range shirts {
			for _, pant := range pants {
				if season == "SUMMER" {
					// کلاه حتماً باید باشد، پس حالت بدون کلاه حذف شد
					for _, cap := range caps {
						fmt.Printf("SHIRT: %s PANTS: %s CAP: %s\n", shirt, pant, cap)
					}
				} else if season == "SPRING" || season == "FALL" {
					// کت اختیاری، کلاه اختیاری
					for _, cap := range append([]string{""}, caps...) {
						line := ""
						if coat != "" {
							line += fmt.Sprintf("COAT: %s ", coat)
						}
						line += fmt.Sprintf("SHIRT: %s PANTS: %s", shirt, pant)
						if cap != "" {
							line += fmt.Sprintf(" CAP: %s", cap)
						}
						fmt.Println(line)
					}
				}
			}
		}
	}
}

func printWinterCombinations(coats, shirts, pants, jackets []string) {
	for _, shirt := range shirts {
		for _, pant := range pants {
			// فقط یکی از کت یا ژاکت
			for _, coat := range coats {
				fmt.Printf("COAT: %s SHIRT: %s PANTS: %s\n", coat, shirt, pant)
			}
			for _, jacket := range jackets {
				fmt.Printf("SHIRT: %s PANTS: %s JACKET: %s\n", shirt, pant, jacket)
			}
		}
	}
}
