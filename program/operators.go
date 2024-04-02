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

		expression, err := p.parseExpression(s)
		if err != nil {
			return err
		}

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
	expression, err := p.parseExpression(s)
	if err != nil {
		return err
	}

	if len(expression) <= 1 {
		return ErrInvalidParams
	}

	if expression[0].tokenType != TokenTypeEquals {
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
func (p *Program) gosub(s *tinybasic.LineScanner) error {
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
			// Добавляем индекс текущей строки в стек вызовов, чтоб вернуться сюда после RETURN
			p.gosubCallingLinesIndexes = append(p.gosubCallingLinesIndexes, p.currentIndex)

			// Ставим i-1, т.к. потом мы увеличиваем currentIndex на 1 в цикле в функции Run
			p.currentIndex = i - 1
			return nil
		}
	}

	return errors.New("line not found")
}
func (p *Program) returnOperator(_ *tinybasic.LineScanner) error {
	if len(p.gosubCallingLinesIndexes) == 0 {
		return ErrInvalidParams
	}

	// берём индекс строки, из которой была вызвана текущая подпрограмма
	index := p.gosubCallingLinesIndexes[len(p.gosubCallingLinesIndexes)-1]

	// удаляем индекс из стека
	p.gosubCallingLinesIndexes = p.gosubCallingLinesIndexes[:len(p.gosubCallingLinesIndexes)-1]

	p.currentIndex = index
	// после завершения работы этого обработчика p.currentIndex увеличится на 1,
	// и мы окажемся на следующей строке после той, из которой вызывали

	return nil
}
func (p *Program) ifOperator(s *tinybasic.LineScanner) error {
	spaces := s.GetSpaces()
	if spaces == nil {
		return ErrInvalidParams
	}

	expression, err := p.parseExpression(s)
	if err != nil {
		return err
	}
	if len(expression) <= 1 {
		return ErrInvalidParams
	}

	if expression[len(expression)-1].tokenType != TokenTypeThen {
		return ErrInvalidParams
	}

	result, err := p.calculateExpression(expression[:len(expression)-1])
	if err != nil {
		return ErrInvalidParams
	}

	if result != 0 {
		s.GetSpaces()
		operator := s.GetStrings(program.supportedOperators)
		if operator == nil {
			return ErrInvalidParams
		}

		err = Operators[Operator(*operator)](s)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}
func (p *Program) forOperator(s *tinybasic.LineScanner) error {
	s.GetSpaces()

	cycle := cycleContext{
		startLineIndex: p.currentIndex,
		variableName:   "",
		start:          0,
		stop:           0,
		step:           1,
	}

	parser := tinybasic.NewLineParserWithScanner(s)

	variableName := parser.GetVariable()
	if variableName == nil {
		return ErrInvalidParams
	}

	cycle.variableName = *variableName

	expression, err := p.parseExpression(s)
	if err != nil {
		return err
	}

	if len(expression) < 2 {
		return ErrInvalidParams
	}

	if expression[0].tokenType != TokenTypeEquals {
		return ErrInvalidParams
	}

	tokenToIndex := 0
	tokenStepIndex := 0
	breakCycle := false

	for i, t := range expression {
		switch t.tokenType {
		case TokenTypeTo:
			tokenToIndex = i

			// Если строка FOR A TO 10
			if i == 0 {
				return errors.New("tokenToIndex == 0")
			}

			// Если строка FOR A=123 TO
			if i+1 >= len(expression) {
				return ErrInvalidParams
			}

			cycle.start, err = p.calculateExpression(expression[1:i])
			if err != nil {
				return err
			}

			p.vars.Set(cycle.variableName, cycle.start)

		case TokenTypeStep:
			tokenStepIndex = i

			if tokenToIndex == 0 {
				return errors.New("tokenToIndex == 0")
			}

			if tokenToIndex+1 >= len(expression) {
				return ErrInvalidParams
			}

			cycle.stop, err = p.calculateExpression(expression[tokenToIndex+1 : tokenStepIndex])
			if err != nil {
				return err
			}

			if i+1 >= len(expression) {
				return ErrInvalidParams
			}

			cycle.step, err = p.calculateExpression(expression[tokenStepIndex+1:])
			if err != nil {
				return err
			}

			breakCycle = true
		}

		if breakCycle {
			break
		}
	}

	// Если не нашли STEP, то всё после TO - это cycle.stop
	if tokenStepIndex == 0 {
		cycle.stop, err = p.calculateExpression(expression[tokenToIndex+1:])
		if err != nil {
			return err
		}
	}

	p.cycles = append(p.cycles, cycle)

	return nil
}
func (p *Program) next(s *tinybasic.LineScanner) error {
	s.GetSpaces()

	if len(p.cycles) == 0 {
		return errors.New("unsupported NEXT without cycle")
	}

	cycle := p.cycles[len(p.cycles)-1]

	cycleVariableValue := p.vars.Get(cycle.variableName)
	cycleVariableValue += cycle.step
	p.vars.Set(cycle.variableName, cycleVariableValue)

	if cycle.stop > cycle.start && cycleVariableValue > cycle.stop {
		p.vars.Set(cycle.variableName, 0)
		p.cycles = p.cycles[:len(p.cycles)-1]
	} else if cycle.stop < cycle.start && cycleVariableValue < cycle.stop {
		p.vars.Set(cycle.variableName, 0)
		p.cycles = p.cycles[:len(p.cycles)-1]
	} else {
		p.currentIndex = cycle.startLineIndex
	}

	return nil
}
func (p *Program) exit(_ *tinybasic.LineScanner) error {
	if len(p.cycles) == 0 {
		return errors.New("unsupported EXIT without cycle")
	}

	cycle := p.cycles[len(p.cycles)-1]

	p.vars.Set(cycle.variableName, 0)
	p.cycles = p.cycles[:len(p.cycles)-1]

	return nil
}
