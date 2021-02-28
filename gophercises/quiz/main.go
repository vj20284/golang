package main

import (
	"flag"
	"bufio"
	"os"
	"fmt"
	"strings"
)

func evalAnswer(input chan string, p Problem) bool {
	return true
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "A csv file with questions and answers")
	flag.Parse()
	
	reader := bufio.NewReader(os.Stdin)
	var score = 0
	input := make(chan string)
	problems, err := filereader(*csvFilename)
	if err != nil {
		fmt.Printf("Failed to open file %s\n", *csvFilename)
		os.Exit(1)
	}
	for i := 0; i < len(problems); i++ {
		select {
			case 
			result := evalAnswer(input, problems[i])
		if result {
			score += 1
		}
		}
	}
}
