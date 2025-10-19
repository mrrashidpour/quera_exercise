package main

import "fmt"

func main() {
	var income int
	var tax int = 0

	fmt.Scan(&income)
	remaining := income

	if remaining > 100 {
		tax += 100 * 5 / 100
		remaining -= 100
	} else {
		tax += remaining * 5 / 100
		remaining = 0
	}

	if remaining > 400 {
		tax += 400 * 10 / 100
		remaining -= 400
	} else {
		tax += remaining * 10 / 100
		remaining = 0
	}

	if remaining > 500 {
		tax += 500 * 15 / 100
		remaining -= 500
	} else {
		tax += remaining * 15 / 100
		remaining = 0
	}

	tax += remaining * 20 / 100

	fmt.Println(tax)
}
