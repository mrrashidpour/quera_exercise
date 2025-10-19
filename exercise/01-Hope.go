package main

import "fmt"

func main() {

	var p int
	var q int

	fmt.Scanf("%d %d", &p, &q)

	for i := 1; i <= q; i++ {

		if i%p == 0 {

			for j := 0; j < i/p; j++ {
				fmt.Print("Hope ")
			}
			fmt.Print("\n")
		} else {
			fmt.Println(i)
		}

	}
}
