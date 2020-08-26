package main

// Given a list of parameters as follows:
//   ./argparse --numbers 1 2 3 --flag -o stdout
//
// more specifically, ["--numbers", "1", "2", "3", "--flag", "-o", "stdout"],
// I'd like to create a mapping from keys to values:
//   "numbers": [1, 2, 3]
//   "flag": true
//   "output": stdout
//   "input": None
//
// I'm going to skip actually specifying the formatting, though, and focus on
// string parsing as a way to learn some Go.

import (
	"fmt"
	"os"
	"strings"
)

func parseArg(parameters []string) map[string]interface{} {
	result := map[string]interface{}{} // string: any

outer:
	for i := 0; i < len(parameters); {
		str := parameters[i]
		fmt.Println(str)

		key := str

		if strings.HasPrefix(key, "--") {
			key = key[2:]
		} else if strings.HasPrefix(key, "-") {
			key = key[1:]
		}

		if strings.LastIndex(key, "s") == len(key)-1 &&
			strings.HasPrefix(str, "--") {
			fmt.Println("  plural argument")

			i += 1
			if i >= len(parameters) {
				break
			}

			next := parameters[i]
			for !strings.HasPrefix(next, "-") {
				if _, ok := result[key]; ok {
					result[key] = append(result[key].([]string), next)
				} else {
					result[key] = []string{next}
				}

				i += 1
				if i >= len(parameters) {
					break outer
				}

				next = parameters[i]
			}
		} else {
			if i+1 >= len(parameters) {
				result[key] = true
				break
			}

			next := parameters[i+1]
			if strings.HasPrefix(next, "-") { // str is just a flag
				result[key] = true
			} else {
				result[key] = next
				i++
			}

			i++
		}

		fmt.Printf("  values: %s\n", result[key])
	}

	return result
}

func main() {
	args := parseArg([]string{
		"--numbers",
		"1",
		"2",
		"3",
		"--flag",
		"-o",
		"stdout",
	})

	fmt.Println()
	fmt.Println(args["o"])
	fmt.Println(args["flag"])

	fmt.Println()
	fmt.Println(parseArg(os.Args[1:]))
}
