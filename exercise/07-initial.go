package main

import "fmt"

func ConvertToDigitalFormat(hour, minute, second int) string {
	return fmt.Sprintf("%02d:%02d:%02d", hour, minute, second)
}

func ExtractTimeUnits(seconds int) (int, int, int) {
	hours := seconds / 3600
	seconds %= 3600
	minutes := seconds / 60
	seconds %= 60
	return hours, minutes, seconds
}

func main() {
	fmt.Println(ConvertToDigitalFormat(2, 23, 4))   // 02:23:04
	fmt.Println(ConvertToDigitalFormat(18, 3, 0))   // 18:03:00
	fmt.Println(ConvertToDigitalFormat(0, 0, 0))    // 00:00:00
	fmt.Println(ConvertToDigitalFormat(23, 59, 59)) // 23:59:59

	h, m, s := ExtractTimeUnits(5500)
	fmt.Printf("%d %d %d\n", h, m, s) // 1 31 40

	h, m, s = ExtractTimeUnits(3661)
	fmt.Printf("%d %d %d\n", h, m, s) // 1 1 1

	h, m, s = ExtractTimeUnits(0)
	fmt.Printf("%d %d %d\n", h, m, s) // 0 0 0
}
