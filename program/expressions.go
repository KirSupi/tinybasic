package program

import (
	"math/rand/v2"

	"tinybasic/tinybasic"
)

func (p *Program) calculateExpression(expression []token) (res int, err error) {
	node, err := p.generateSyntaxTree(expression)
	if err != nil {
		return res, err
	}

	if node == nil {
		return res, nil
	}

	res, err = p.calculateExpressionTreeItem(*node)
	if err != nil {
		return res, err
	}

	return res, nil
}
func (p *Program) calculateExpressionTreeItem(treeItem syntaxTreeNode) (res int, err error) {
	if treeItem.left == nil && treeItem.right == nil {
		if treeItem.item.tokenType == TokenTypeValue {
			return *treeItem.item.value, nil
		}

		return res, ErrInvalidParams
	}

	if treeItem.item == nil {
		if treeItem.left != nil && treeItem.right == nil {
			return p.calculateExpressionTreeItem(*treeItem.left)
		}
		if treeItem.left == nil && treeItem.right != nil {
			return p.calculateExpressionTreeItem(*treeItem.right)
		}

		return res, ErrInvalidParams
	}

	resLeft := 0
	if treeItem.left != nil {
		resLeft, err = p.calculateExpressionTreeItem(*treeItem.left)
		if err != nil {
			return res, err
		}
	}

	resRight := 0
	if treeItem.right != nil {
		resRight, err = p.calculateExpressionTreeItem(*treeItem.right)
		if err != nil {
			return res, err
		}
	}

	// Unary operations
	if treeItem.left == nil {
		switch treeItem.item.tokenType {
		case TokenTypeOperatorPlus:
			return resRight, nil
		case TokenTypeOperatorMinus:
			return -resRight, nil
		case TokenTypeNot:
			return fromBool(!toBool(resRight)), nil
		}
	}

	switch treeItem.item.tokenType {
	case TokenTypeValue:
		return *treeItem.item.value, nil
	case TokenTypeOperatorPlus:
		return resLeft + resRight, nil
	case TokenTypeOperatorMinus:
		return resLeft - resRight, nil
	case TokenTypeOperatorDivision:
		return resLeft / resRight, nil
	case TokenTypeMod:
		return resLeft % resRight, nil
	case TokenTypeOperatorMultiply:
		return resLeft * resRight, nil
	case TokenTypeOr:
		return fromBool(toBool(resLeft) || toBool(resRight)), nil
	case TokenTypeAnd:
		return fromBool(toBool(resLeft) && toBool(resRight)), nil
	case TokenTypeXor:
		return fromBool(toBool(resLeft) != toBool(resRight)), nil
	case TokenTypeBitShiftLeft:
		return resLeft << resRight, nil
	case TokenTypeBitShiftRight:
		return resLeft >> resRight, nil
	case TokenTypeEquals:
		return fromBool(resLeft == resRight), nil
	case TokenTypeNotEquals, TokenTypeNotEquals2:
		return fromBool(resLeft != resRight), nil
	case TokenTypeLess:
		return fromBool(resLeft < resRight), nil
	case TokenTypeLessOrEquals:
		return fromBool(resLeft <= resRight), nil
	case TokenTypeGreater:
		return fromBool(resLeft > resRight), nil
	case TokenTypeGreaterOrEquals:
		return fromBool(resLeft >= resRight), nil
	default:
		return res, ErrInvalidParams
	}
}

func (p *Program) parseExpression(scanner *tinybasic.LineScanner) (items []token, err error) {
	parser := tinybasic.NewLineParserWithScanner(scanner)

	for !scanner.IsEOL() {
		if len(items) >= 1 && items[len(items)-1].tokenType == TokenTypeThen {
			break
		}

		scanner.GetSpaces()

		rndFunction := scanner.GetString(TokenTypeRnd)
		if rndFunction != nil {
			scanner.GetSpaces()

			bracket := scanner.GetString(TokenTypeBracketOpen)
			if bracket == nil {
				return nil, ErrInvalidParams
			}

			scanner.Shift(-1) // отступаем обратно, чтоб спарсить вместе с открывающей скобкой

			nextItems, err := p.parseExpression(scanner)
			if err != nil {
				return append(items, nextItems...), err
			}

			index, err := getClosingBracketIndex(nextItems)
			if err != nil {
				return append(items, nextItems...), err
			}

			if len(nextItems) <= index || index < 1 {
				return append(items, nextItems...), ErrInvalidParams
			}

			randomMaxValue, err := p.calculateExpression(nextItems[1:index])
			if err != nil {
				return append(items, nextItems...), err
			}

			if randomMaxValue <= 0 {
				return append(items, nextItems...), ErrInvalidParams
			}

			randomValue := rand.N(randomMaxValue)

			items = append(items, token{
				tokenType: TokenTypeValue,
				value:     &randomValue,
			})

			if len(nextItems) > index+1 {
				items = append(items, nextItems[index+1:]...)
			}

			continue
		}

		thenOperator := scanner.GetString(TokenTypeThen)
		if thenOperator != nil {
			items = append(items, token{
				tokenType: TokenTypeThen,
				value:     nil,
			})

			break
		}

		variableName := parser.GetVariable()
		if variableName != nil {
			value := p.vars.Get(*variableName)
			items = append(items, token{
				tokenType: TokenTypeValue,
				value:     &value,
			})
			continue
		}

		value := scanner.GetNumber()
		if value != nil {
			items = append(items, token{
				tokenType: TokenTypeValue,
				value:     value,
			})
			continue
		}

		item := scanner.GetStrings([]string{
			TokenTypeOperatorPlus,
			TokenTypeOperatorMinus,
			TokenTypeOperatorDivision,
			TokenTypeOperatorMultiply,
			TokenTypeBracketOpen,
			TokenTypeBracketClose,
			TokenTypeOr,
			TokenTypeXor,
			TokenTypeAnd,
			TokenTypeNot,
			TokenTypeMod,
			TokenTypeBitShiftLeft,
			TokenTypeBitShiftRight,
			TokenTypeGreaterOrEquals,
			TokenTypeLessOrEquals,
			TokenTypeNotEquals,
			TokenTypeNotEquals2,
			TokenTypeEquals,
			TokenTypeLess,
			TokenTypeGreater,
		})
		if item != nil {
			items = append(items, token{
				tokenType: tokenType(*item),
				value:     nil,
			})

			continue
		}

		break
	}

	return items, nil
}

type tokenType string

const (
	TokenTypeValue            = "value"
	TokenTypeBracketOpen      = "("
	TokenTypeBracketClose     = ")"
	TokenTypeOperatorPlus     = "+"
	TokenTypeOperatorMinus    = "-"
	TokenTypeOperatorDivision = "/"
	TokenTypeOperatorMultiply = "*"
	TokenTypeOr               = "or"
	TokenTypeXor              = "xor"
	TokenTypeAnd              = "and"
	TokenTypeNot              = "not"
	TokenTypeMod              = "mod"
	TokenTypeBitShiftLeft     = "<<"
	TokenTypeBitShiftRight    = ">>"
	TokenTypeEquals           = "="
	TokenTypeNotEquals        = "<>"
	TokenTypeNotEquals2       = "><"
	TokenTypeLess             = "<"
	TokenTypeGreater          = ">"
	TokenTypeGreaterOrEquals  = ">="
	TokenTypeLessOrEquals     = "<="
	TokenTypeThen             = "THEN"
	TokenTypeRnd              = "RND"
)

type token struct {
	tokenType tokenType
	value     *int
}

func toBool(v int) bool { return v != 0 }
func fromBool(v bool) int {
	if v {
		return 1
	}

	return 0
}

type syntaxTreeNode struct {
	item  *token
	left  *syntaxTreeNode
	right *syntaxTreeNode
}

func getClosingBracketIndex(items []token) (index int, err error) {
	level := 0

	for i, item := range items {
		switch item.tokenType {
		case TokenTypeBracketOpen:
			level++
		case TokenTypeBracketClose:
			level--

			if level == 0 {
				return i, nil
			}
			if level < 0 {
				return index, ErrInvalidParams
			}
		}
	}

	return index, ErrInvalidParams
}

var operationsPriority = map[tokenType]int{
	TokenTypeEquals:          1,
	TokenTypeNotEquals:       1,
	TokenTypeNotEquals2:      1,
	TokenTypeLess:            1,
	TokenTypeGreater:         1,
	TokenTypeGreaterOrEquals: 1,
	TokenTypeLessOrEquals:    1,

	TokenTypeOperatorPlus:  2,
	TokenTypeOperatorMinus: 2,
	TokenTypeOr:            2,
	TokenTypeXor:           2,

	TokenTypeOperatorMultiply: 3,
	TokenTypeOperatorDivision: 3,
	TokenTypeBitShiftLeft:     3,
	TokenTypeBitShiftRight:    3,
	TokenTypeMod:              3,
	TokenTypeAnd:              3,

	TokenTypeNot: 4,
}

func (p *Program) generateSyntaxTree(items []token) (*syntaxTreeNode, error) {
	ops := []*token(nil)
	nodes := []*syntaxTreeNode(nil)
	// A + 10 * 11 + 1
	for i := 0; i < len(items); i++ {
		item := items[i]
		if item.tokenType == TokenTypeValue {
			nodes = append(nodes, &syntaxTreeNode{
				item:  &items[i],
				left:  nil,
				right: nil,
			})
		} else if priority, isOperator := operationsPriority[item.tokenType]; isOperator {
			if len(ops) == 0 || priority >= operationsPriority[ops[len(ops)-1].tokenType] {
				ops = append(ops, &items[i])
				continue
			}

			node := &syntaxTreeNode{
				item:  ops[len(ops)-1],
				left:  nil,
				right: nil,
			}
			ops = ops[:len(ops)-1]

			if len(nodes) >= 1 {
				node.right = nodes[len(nodes)-1]
				nodes = nodes[:len(nodes)-1]
			}

			if len(nodes) >= 1 {
				node.left = nodes[len(nodes)-1]
				nodes = nodes[:len(nodes)-1]
			}

			nodes = append(nodes, node)

			for len(ops) != 0 && priority < operationsPriority[ops[len(ops)-1].tokenType] {
				node = &syntaxTreeNode{
					item:  ops[len(ops)-1],
					left:  nil,
					right: nil,
				}
				ops = ops[:len(ops)-1]

				if len(nodes) >= 1 {
					node.right = nodes[len(nodes)-1]
					nodes = nodes[:len(nodes)-1]
				}

				if len(nodes) >= 1 {
					node.left = nodes[len(nodes)-1]
					nodes = nodes[:len(nodes)-1]
				}

				nodes = append(nodes, node)
			}

			ops = append(ops, &items[i])
		} else if item.tokenType == TokenTypeBracketOpen {
			index, err := getClosingBracketIndex(items[i:])
			if err != nil {
				return nil, err
			}

			node, err := p.generateSyntaxTree(items[i+1 : i+index])
			if err != nil {
				return nil, err
			}

			i = i + index

			nodes = append(nodes, node)
		}
	}

	// ops: +(1), +(2)
	// nodes: "A", "11*10", "1"

	for len(ops) != 0 {
		node := &syntaxTreeNode{
			item:  ops[len(ops)-1],
			left:  nil,
			right: nil,
		}

		ops = ops[:len(ops)-1]

		if len(nodes) >= 1 {
			node.right = nodes[len(nodes)-1]
			nodes = nodes[:len(nodes)-1]
		}

		if len(nodes) >= 1 {
			node.left = nodes[len(nodes)-1]
			nodes = nodes[:len(nodes)-1]
		}

		nodes = append(nodes, node)
	}

	if len(nodes) == 1 {
		return nodes[0], nil
	}

	return nil, ErrInvalidParams
}
