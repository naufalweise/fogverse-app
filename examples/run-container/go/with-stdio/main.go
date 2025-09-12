package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// Sum numbers per line

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		println(sumNumbers(line))
	}
}

// sumNumbers returns the sums of numbers contained in a line
func sumNumbers(line string) int {
	sum := 0
	numberStrList := strings.Split(line, " ")
	for _, numberStr := range numberStrList {
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			continue
		}
		sum += number
	}
	return sum
}
