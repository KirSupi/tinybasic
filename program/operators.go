package program

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"tinybasic/tinybasic"
)

func (p *Program) clear(_ *tinybasic.LineScanner) error {
	p.s.Clear()
	p.currentIndex = -1

	return nil
}
func (p *Program) rem(_ *tinybasic.LineScanner) error { return nil }
func (p *Program) list(_ *tinybasic.LineScanner) error {
	for _, line := range p.s.Lines {
		fmt.Println(line.Label, line.Text)
	}

	return nil
}
func (p *Program) load(s *tinybasic.LineScanner) error {
	spaces := s.GetSpaces()
	if spaces == nil {
		return ErrInvalidParams
	}

	parser := tinybasic.NewLineParser(s.GetTail())

	param := parser.GetQuotedString()
	if param == nil {
		return errors.New("invalid param")
	}

	err := p.s.Load(*param)
	if err != nil {
		return err
	}

	return nil
}
func (p *Program) save(s *tinybasic.LineScanner) error {
	spaces := s.GetSpaces()
	if spaces == nil {
		return ErrInvalidParams
	}

	parser := tinybasic.NewLineParser(s.GetTail())

	param := parser.GetQuotedString()
	if param == nil {
		return ErrInvalidParams
	}

	err := p.s.Load(*param)
	if err != nil {
		return err
	}

	return nil
}
func (p *Program) gotoOperator(s *tinybasic.LineScanner) error {
	spaces := s.GetSpaces()
	if spaces == nil {
		return ErrInvalidParams
	}

	param := s.GetNumber()
	if param == nil {
		return ErrInvalidParams
	}

	for i, line := range p.s.Lines {
		if line.Label == *param {
			p.currentIndex = i - 1 // Ставим i-1, т.к. потом мы увеличиваем currentIndex на 1 в цикле в функции Run
			return nil
		}
	}

	return errors.New("line not found")
}
func (p *Program) run(s *tinybasic.LineScanner) error {
	program.vars = tinybasic.NewVariables()

	spaces := s.GetSpaces()
	if spaces == nil {
		p.currentIndex = -1

		return nil
	}

	parser := tinybasic.NewLineParser(s.GetTail())

	param := parser.GetVariable()
	if param == nil {
		p.currentIndex = -1

		return nil
	}

	lineLabel, err := strconv.Atoi(*param)
	if err != nil {
		return err
	}

	for i, line := range p.s.Lines {
		if line.Label == lineLabel {
			p.currentIndex = i - 1 // Ставим i-1, т.к. потом мы увеличиваем currentIndex на 1 в цикле в функции Run
			return nil
		}
	}

	return errors.New("line not found")
}
func (p *Program) end(_ *tinybasic.LineScanner) error {
	p.currentIndex = len(p.s.Lines)
	program.vars = tinybasic.NewVariables()

	return nil
}
func (p *Program) input(s *tinybasic.LineScanner) (err error) {
	spaces := s.GetSpaces()
	if spaces == nil {
		return ErrInvalidParams
	}

	variablesListStr := s.GetTail()

	variablesList := strings.Split(variablesListStr, ", ")
	if len(variablesList) == 0 {
		return ErrInvalidParams
	}

	for _, variableName := range variablesList {
		variableValue := 0

		for {
			_, err = fmt.Scanf("%d", &variableValue)
			if err != nil {
				fmt.Println("invalid value")
				continue
			}

			break
		}

		p.vars.Set(variableName, variableValue)
	}

	return nil
}
func (p *Program) print(s *tinybasic.LineScanner) error {
	// 10 PRINT
	// 10 PRINT "abc abc"    , A , "BBB";

	parser := tinybasic.NewLineParserWithScanner(s)

	newLine := true

	for !s.IsEOL() {
		s.GetSpaces()

		quoted := parser.GetQuotedString()
		if quoted != nil {
			fmt.Print(*quoted)
			continue
		}

		comma := s.GetChar(',')
		if comma != nil {
			fmt.Print("\t")

			if s.IsEOL() {
				newLine = false
			}

			continue
		}

		semicolon := s.GetChar(';')
		if semicolon != nil {
			fmt.Print(" ")

			if s.IsEOL() {
				newLine = false
			}

			continue
		}

		expression := p.parseExpression(s)
		if len(expression) == 0 {
			return ErrInvalidParams
		}

		result, err := p.calculateExpression(expression)
		if err != nil {
			return ErrInvalidParams
		}

		fmt.Print(result)
	}

	if newLine {
		fmt.Println()
	}

	return nil
}
func (p *Program) let(s *tinybasic.LineScanner) error {
	// 10 LET A = B + 10

	// " A = B + 10"
	spaces := s.GetSpaces()
	if spaces == nil {
		return ErrInvalidParams
	}

	parser := tinybasic.NewLineParserWithScanner(s)

	variableName := parser.GetVariable()
	if variableName == nil {
		return ErrInvalidParams
	}

	// " = B + 10"
	expression := p.parseExpression(s)
	if len(expression) <= 1 {
		return ErrInvalidParams
	}

	if expression[0].itemType != ExpressionItemTypeEquals {
		return ErrInvalidParams
	}

	// парсим ["=", "B", "+", "10"]
	result, err := p.calculateExpression(expression[1:])
	if err != nil {
		return ErrInvalidParams
	}

	p.vars.Set(*variableName, result)

	return nil
}
