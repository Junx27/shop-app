package main

import "fmt"

func isDecimal(s string) int {
	result := 0
	prevValue := 0
	value := 0

	for _, c := range s {
		switch string(c) {
		case "M":
			value = 1000
		case "D":
			value = 500
		case "C":
			value = 100
		case "L":
			value = 50
		case "X":
			value = 10
		case "V":
			value = 5
		default:
			value = 1
		}

		if prevValue < value {
			result += value - 2*prevValue
		} else {
			result += value
		}

		prevValue = value
	}

	return result
}

func main() {

	text := "MCMXCIV"
	fmt.Println(isDecimal(text))
}
