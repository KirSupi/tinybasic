package main

import (
	"log"

	"tinybasic/program"
	"tinybasic/tinybasic"
)

func main() {
	s := tinybasic.NewSource()

	//err := s.Load("./test.bas")
	err := s.Load("./euphoria.bas")
	if err != nil {
		log.Fatalln(err)
	}

	err = program.Run(s)
	if err != nil {
		log.Fatalln(err)
	}
}
