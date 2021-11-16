package main

import (
	"fmt"
	"os"
	"strings"
)

// the language has 4 registers
var registers = [4]int{0, 0, 0, 0}

// 3 arrays containing the different types of instructions we allow
var rType = [4]string{"add", "sub", "set", "jeq"}
var jType = [1]string{"j"}
var specType = [3]string{"input", "print", "exit"}

func readLineForLine(file string) []string {
	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	instructions := []string{}

	// split data on newline
	formattedData := strings.Split(string(data), "\n")
	for _, line := range formattedData {

		// find comment or end of line
		commentStart := strings.Index(line, "//")
		var endOfLine int
		if commentStart == -1 {
			endOfLine = len(line)
		} else {
			endOfLine = commentStart
		}

		line = line[:endOfLine]
		line = strings.TrimSpace(line)
		//fmt.Printf("line #%d: %s\n", i+1, line)
		instructions = append(instructions, line)
	}
	fmt.Println("len(instructions) =", len(instructions))
	// instructions is now a slice of all lines without comments or extra whitespace
	return instructions
}

func main() {
	//fmt.Println(j_type)
	// read from stdin and print the output
	fmt.Println(len(os.Args))
	if len(os.Args) != 2 {
		fmt.Println("\tIncorrect number of arguments, command should be\n\t'go run . <filename>'")
	} else {
		file := os.Args[1]
		fmt.Println(file)
		// check that last 5 chars is .bbvv
		if len(file) < 5 || file[len(file)-5:] != ".bbvv" {
			fmt.Println("bs file =", file)
		} else {
			// it's a legit bbvv file
			instructions := readLineForLine(file)
			fmt.Println(instructions)
		}
	}
	//_ := os.Args[0]
	//_ := os.Args[1]

	// "parse" the lines
}
