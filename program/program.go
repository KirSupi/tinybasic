package program

import (
	"errors"
	"fmt"
	"slices"

	"tinybasic/tinybasic"
)

var program Program

type Program struct {
	s                        *tinybasic.Source
	vars                     *tinybasic.Variables
	currentIndex             int
	gosubCallingLinesIndexes []int // Номера строк, из которых мы запустили подпрограммы (GOSUB)
	supportedOperators       []string
}

func Run(s *tinybasic.Source) (err error) {
	Operators = map[Operator]func(s *tinybasic.LineScanner) error{
		CLEAR:  program.clear,
		REM:    program.rem,
		LIST:   program.list,
		LOAD:   program.load,
		SAVE:   program.save,
		GOTO:   program.gotoOperator,
		RUN:    program.run,
		END:    program.end,
		INPUT:  program.input,
		PRINT:  program.print,
		LET:    program.let,
		GOSUB:  program.gosub,
		RETURN: program.returnOperator,
		IF:     program.ifOperator,
	}
	program.s = s
	program.vars = tinybasic.NewVariables()

	for cmd := range Operators {
		program.supportedOperators = append(program.supportedOperators, string(cmd))
	}

	line := tinybasic.Line{}

	program.currentIndex = 0

	for program.currentIndex < len(program.s.Lines) {
		line = program.s.Lines[program.currentIndex]

		scanner := tinybasic.NewLineScanner(line.Text)

		operator := scanner.GetStrings(program.supportedOperators)
		if operator == nil {
			return fmt.Errorf("unsupported operator on line %d %s", line.Label, line.Text)
		}

		if slices.Contains([]int{
			2390,
		}, line.Label) {
			fmt.Printf("running %d %s\n", line.Label, line.Text)
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
	CLEAR  Operator = "CLEAR"
	REM    Operator = "REM"
	LIST   Operator = "LIST"
	LOAD   Operator = "LOAD"
	SAVE   Operator = "SAVE"
	GOTO   Operator = "GOTO"
	RUN    Operator = "RUN"
	END    Operator = "END"
	INPUT  Operator = "INPUT"
	PRINT  Operator = "PRINT"
	LET    Operator = "LET"
	GOSUB  Operator = "GOSUB"
	RETURN Operator = "RETURN"
	IF     Operator = "IF"
)

var Operators map[Operator]func(s *tinybasic.LineScanner) error

var ErrInvalidParams = errors.New("invalid params")
