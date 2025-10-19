package main

import (
	"strconv"
)

type FilterFunc func(int) bool
type MapperFunc func(int) int

func IsSquare(x int) bool {
	if x < 0 {
		return false
	}

	for i := 0; i*i <= x; i++ {
		if i*i == x {
			return true
		}
	}
	return false
}

func IsPalindrome(x int) bool {

	if x < 0 {
		x = -x
	}

	s := strconv.Itoa(x)
	length := len(s)

	for i := 0; i < length/2; i++ {
		if s[i] != s[length-1-i] {
			return false
		}
	}
	return true
}

func Abs(num int) int {
	if num < 0 {
		return -num
	}
	return num
}

func Cube(num int) int {
	return num * num * num
}

func Filter(input []int, f FilterFunc) []int {
	if input == nil {
		return nil
	}

	result := make([]int, 0, len(input))

	for _, num := range input {
		if f(num) {
			result = append(result, num)
		}
	}

	return result
}

func Map(input []int, m MapperFunc) []int {
	if input == nil {
		return nil
	}

	result := make([]int, len(input))

	for i, num := range input {
		result[i] = m(num)
	}

	return result
}
