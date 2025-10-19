package main

func AddElement(numbers *[]int, element int) {
	*numbers = append(*numbers, element)
}

func FindMin(numbers *[]int) int {
	if numbers == nil || len(*numbers) == 0 {
		return 0
	}

	min := (*numbers)[0]
	for _, num := range *numbers {
		if num < min {
			min = num
		}
	}
	return min
}

func ReverseSlice(numbers *[]int) {
	if numbers == nil || len(*numbers) <= 1 {
		return
	}

	slice := *numbers
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func SwapElements(numbers *[]int, i, j int) {
	if numbers == nil {
		return
	}

	slice := *numbers
	if i < 0 || j < 0 || i >= len(slice) || j >= len(slice) {
		return
	}

	slice[i], slice[j] = slice[j], slice[i]
}

func main() {
	// مثال برای AddElement
	nums := []int{1, 2, 3}
	AddElement(&nums, 4)
	// nums = [1, 2, 3, 4]

	// مثال برای FindMin
	//min := FindMin(&nums)
	// min = 1

	// مثال برای ReverseSlice
	ReverseSlice(&nums)
	// nums = [4, 3, 2, 1]

	// مثال برای SwapElements
	SwapElements(&nums, 0, 3)
	// nums = [1, 3, 2, 4]
}
