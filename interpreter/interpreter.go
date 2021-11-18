package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Helper function that checks if a slice contains a string
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if a == b {
			return true
		}
	}
	return false
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// the language has 4 registers
var registers = [4]int{0, 0, 0, 0}

// 3 arrays containing the different types of instructions we allow
var registerType = []string{"add", "sub", "set", "jeq"}
var jumpType = []string{"j"}
var specialType = []string{"input", "print", "exit"}

func splitIntoLines(file []byte) []string {
	instructions := []string{}

	// split data on newline
	formattedData := strings.Split(string(file), "\n")
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
		instructions = append(instructions, line)
	}
	// instructions is now a slice of all lines without comments or extra whitespace
	return instructions
}

func instructionsAreValid(instructions []string) bool {
	for i, line := range instructions {
		//fmt.Println(i, line)
		words := strings.Split(line, " ")
		if stringInSlice(words[0], registerType) {
			fmt.Println(words[0])
			// if it's a registerType, length of the whole instructions must be 4
			if len(words) != 4 {
				fmt.Printf("Error on line %d: Incorrect instruction of length %d\n", i+1, len(words))
				return false
			} else {

				var j int
				for j = 1; j < 3; j++ {
					fmt.Println(strings.Index(string(words[j][0]), "#"))
					// Register must start with '#'
					if !strings.Contains(string(words[j][0]), "#") {
						fmt.Printf("Error on line %d: Missing '#'\n", j+1)
						return false
					}
				}
				// first arg to an registerType must start with #1, #2 or #3
				if !stringInSlice(string(words[1][1]), []string{"1", "2", "3"}) {
					fmt.Printf("Error on line %d: %s cannot be used as first register argument for r_type instructions\n", i+1, string(words[1][1]))
					return false
				}

				// immidiate can only be a 1 bit value
				if !stringInSlice(string(words[3]), []string{"0", "1"}) {
					fmt.Printf("Error on line %d: invalid immidiate. Can only be a 1 bit value\n", i+1)
					return false
				}

			}
		} else if stringInSlice(words[0], jumpType) {
			// jump instruction should only have a total length of 2
			if len(words) != 2 {
				fmt.Printf("Error on line %d: Incorrect instruction of length %d\n", i+1, len(words))
				return false
			}

			wordAsInt, err := strconv.Atoi(words[1])
			if err != nil {
				panic(err)
			}

			// 5 bit signed value range is (-16, 15)
			if wordAsInt > 15 || wordAsInt < -16 {
				fmt.Printf("Error on line %d: %d does not fit in a 5 bit (signed) value\n", i+1, wordAsInt)
				return false
			}

		} else if stringInSlice(words[0], specialType) {
			// Special instructions should only have length 1
			if len(words) != 1 {
				fmt.Printf("Error on line %d: Incorrect instruction of length {%d}\n", i+1, len(words))
				return false
			}
		} else {
			// "" means it's just an empty line
			if words[0] != "" {
				fmt.Printf("Error on line %d: unknown command: {%s}\n", i+1, words[0])
				return false
			}
		}
	}
	return true
}

// Read from stdio into #1, per specification
func readInput() {
	fmt.Print("> ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()
	input, err := strconv.Atoi(text)
	checkError(err)
	registers[1] = input
}

// Print content of register #1, per specification
func println() {
	fmt.Printf("%d\n", registers[1])
}

func parseRegisterType(currentInstructions []string) (int, int, int) {
	register1, err := strconv.Atoi(string(currentInstructions[1])[1:])
	checkError(err)
	register2, err := strconv.Atoi(string(currentInstructions[2])[1:])
	checkError(err)
	immediate, err := strconv.Atoi(string(currentInstructions[3]))
	checkError(err)

	return register1, register2, immediate
}

func parseInstructions(instructions []string) {
	idx := 0
	for idx < len(instructions) {
		currentInstructions := strings.Split(instructions[idx], " ")
		currentInstruction := currentInstructions[0]

		if currentInstruction == "input" {
			readInput()
		} else if currentInstruction == "exit" {
			break
		} else if currentInstruction == "print" {
			println()
		} else if currentInstruction == "add" {
			reg1, reg2, imm := parseRegisterType(currentInstructions)
			registers[reg1] = registers[reg1] + registers[reg2] + imm
		} else if currentInstruction == "sub" {
			reg1, reg2, imm := parseRegisterType(currentInstructions)
			registers[reg1] = registers[reg1] - registers[reg2] - imm
		} else if currentInstruction == "set" {
			reg1, reg2, imm := parseRegisterType(currentInstructions)
			registers[reg1] = registers[reg2] + imm
		} else if currentInstruction == "j" {
			instructionAsInt, err := strconv.Atoi(currentInstructions[1])
			checkError(err)
			idx += instructionAsInt
		} else if currentInstruction == "jeq" {
			reg1, reg2, imm := parseRegisterType(currentInstructions)
			// if compared regs are eq => increment with the immediate value (program counter)
			if registers[reg1] == registers[reg2] {
				idx += imm
			} else if imm == 0 {
				idx += 1
			}
		}

		idx += 1
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("\tIncorrect number of arguments, command should be\n\t'go run . <filename>'")
	} else {
		file := os.Args[1]
		fmt.Println(file)
		// check that last 5 chars is .bbvv
		if len(file) < 5 || file[len(file)-5:] != ".bbvv" {
			fmt.Println("bs file =", file)
		} else {
			// check if file exists
			fileContent, err := os.ReadFile(file)
			checkError(err)

			// being here means it's a legit bbvv file
			instructions := splitIntoLines(fileContent)
			fmt.Println(instructions)
			// check if the provided instructions are valid according to the bbvv specification
			if instructionsAreValid(instructions) {
				// parse and execute the instructions
				parseInstructions(instructions)
			} else {
				fmt.Println("Provided instructions are not valid.")
			}
		}
	}

}
