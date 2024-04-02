package program

import (
	"errors"
	"fmt"

	"tinybasic/tinybasic"
)

var program Program

type Program struct {
	s                        *tinybasic.Source
	vars                     *tinybasic.Variables
	currentIndex             int
	gosubCallingLinesIndexes []int          // Номера строк, из которых мы запустили подпрограммы (GOSUB)
	cycles                   []cycleContext // Стек из циклов, в которых находится программа
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
		FOR:    program.forOperator,
		NEXT:   program.next,
		EXIT:   program.exit,
	}
	program.s = s
	program.vars = tinybasic.NewVariables()
	program.currentIndex = 0

	for cmd := range Operators {
		program.supportedOperators = append(program.supportedOperators, string(cmd))
	}

	line := tinybasic.Line{}

	// Проходим по всем строкам программы
	for program.currentIndex < len(program.s.Lines) {
		line = program.s.Lines[program.currentIndex]

		scanner := tinybasic.NewLineScanner(line.Text)

		// Пропускаем пробелы, если они есть
		scanner.GetSpaces()

		// Считываем оператор, который надо выполнить
		operator := scanner.GetStrings(program.supportedOperators)
		if operator == nil {
			return fmt.Errorf("unsupported operator on line %d %s", line.Label, line.Text)
		}

		// Находим обработчик для оператора
		handler := Operators[Operator(*operator)]

		// Запускаем обработчик для оператора
		if line.Label == 500 {
			fmt.Println(program.vars.Get("I"))
		}
		err = handler(scanner)
		if err != nil {
			return err
		}

		// Переходим на следующую строку
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
	FOR    Operator = "FOR"
	NEXT   Operator = "NEXT"
	EXIT   Operator = "EXIT"
)

var Operators map[Operator]func(s *tinybasic.LineScanner) error

var ErrInvalidParams = errors.New("invalid params")

// Контекст, который добавляется в стек при выполнении цикла
type cycleContext struct {
	startLineIndex int
	variableName   string
	start          int
	stop           int
	step           int
}
