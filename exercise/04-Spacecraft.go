package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Spaceship struct {
	Name  string
	Count int
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())

	spaceships := make([]Spaceship, n)

	for i := 0; i < n; i++ {
		scanner.Scan()
		line := scanner.Text()
		parts := strings.Split(line, " ")

		name := parts[0]
		fuelData := make([]int, len(parts)-1)

		for j := 1; j < len(parts); j++ {
			fuel, _ := strconv.Atoi(parts[j])
			fuelData[j-1] = fuel
		}

		count := 0
		n := len(fuelData)

		for length := 3; length <= n; length++ {
			for start := 0; start <= n-length; start++ {

				arr := fuelData[start : start+length]

				isAdded := true

				if len(arr) < 3 {
					isAdded = false
				}

				d := arr[1] - arr[0]

				for i := 2; i < len(arr); i++ {
					if arr[i]-arr[i-1] != d {
						isAdded = false
					}
				}

				if isAdded {
					count++
				}

			}
		}

		spaceships[i] = Spaceship{Name: name, Count: count}
	}

	sort.Slice(spaceships, func(i, j int) bool {
		if spaceships[i].Count == spaceships[j].Count {
			return spaceships[i].Name < spaceships[j].Name
		}
		return spaceships[i].Count > spaceships[j].Count
	})

	for _, ship := range spaceships {
		fmt.Printf("%s %d\n", ship.Name, ship.Count)
	}
}
