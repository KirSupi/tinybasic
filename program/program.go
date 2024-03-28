package program

import (
	"errors"
	"fmt"

	"tinybasic/tinybasic"
)

var program Program

type Program struct {
	s            *tinybasic.Source
	vars         *tinybasic.Variables
	currentIndex int
}

func Run(s *tinybasic.Source) (err error) {
	program.s = s
	program.vars = tinybasic.NewVariables()

	supportedOperators := make([]string, 0)
	for cmd := range Operators {
		supportedOperators = append(supportedOperators, string(cmd))
	}

	line := tinybasic.Line{}

	program.currentIndex = 0

	for program.currentIndex < len(program.s.Lines) {
		line = program.s.Lines[program.currentIndex]

		scanner := tinybasic.NewLineScanner(line.Text)

		operator := scanner.GetStrings(supportedOperators)
		if operator == nil {
			return fmt.Errorf("unsupported operator on line %d %s", line.Label, line.Text)
		}

		err = Operators[Operator(*operator)](scanner)
		if err != nil {
			return err
		}

		program.currentIndex++
	}

	return nil
}

type Operator string

const (
	CLEAR Operator = "CLEAR"
	REM   Operator = "REM"
	LIST  Operator = "LIST"
	LOAD  Operator = "LOAD"
	SAVE  Operator = "SAVE"
	GOTO  Operator = "GOTO"
	RUN   Operator = "RUN"
	END   Operator = "END"
	INPUT Operator = "INPUT"
	PRINT Operator = "PRINT"
	LET   Operator = "LET"
)

var Operators = map[Operator]func(s *tinybasic.LineScanner) error{
	CLEAR: program.clear,
	REM:   program.rem,
	LIST:  program.list,
	LOAD:  program.load,
	SAVE:  program.save,
	GOTO:  program.gotoOperator,
	RUN:   program.run,
	END:   program.end,
	INPUT: program.input,
	PRINT: program.print,
	LET:   program.let,
}

var ErrInvalidParams = errors.New("invalid params")
