package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"tinybasic/program"
	"tinybasic/tinybasic"
)

func main() {
	s := tinybasic.NewSource()

	if len(os.Args) < 2 {
		// If file name is not provided, enter interactive mode (similar to Python REPL)
		fmt.Println("Entering interactive mode. Type 'exit' to quit.")
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print(">>> ")
			scanned := scanner.Scan()
			if !scanned {
				return
			}
			line := scanner.Text()
			if strings.ToLower(line) == "exit" {
				return
			}
			if strings.ToLower(line) == "end" {
				break
			}
			// Load the line into source and run it
			err := s.LoadLine(line)
			if err != nil {
				log.Fatalln(err)
			}
			err = program.Run(s)
			if err != nil {
				log.Fatalln(err)
			}
		}
	} else {
		fileName := os.Args[1]
		err := s.Load(fileName)
		if err != nil {
			log.Fatalln(err)
		}

		err = program.Run(s)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
