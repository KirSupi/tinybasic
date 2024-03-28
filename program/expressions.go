package program

import (
	"tinybasic/tinybasic"
)

func (p *Program) calculateExpression(expression []expressionItem) (res int, err error) {
	node, err := p.generateTree(expression)
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
func (p *Program) calculateExpressionTreeItem(treeItem expressionTreeItem) (res int, err error) {
	if treeItem.left == nil && treeItem.right == nil {
		if treeItem.item.itemType == ExpressionItemTypeValue {
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
		switch treeItem.item.itemType {
		case ExpressionItemTypeOperatorPlus:
			return resRight, nil
		case ExpressionItemTypeOperatorMinus:
			return -resRight, nil
		case ExpressionItemTypeNot:
			return fromBool(!toBool(resRight)), nil
		}
	}

	switch treeItem.item.itemType {
	case ExpressionItemTypeValue:
		return *treeItem.item.value, nil
	case ExpressionItemTypeOperatorPlus:
		return resLeft + resRight, nil
	case ExpressionItemTypeOperatorMinus:
		return resLeft - resRight, nil
	case ExpressionItemTypeOperatorDivision:
		return resLeft / resRight, nil
	case ExpressionItemTypeMod:
		return resLeft % resRight, nil
	case ExpressionItemTypeOperatorMultiply:
		return resLeft * resRight, nil
	case ExpressionItemTypeOr:
		return fromBool(toBool(resLeft) || toBool(resRight)), nil
	case ExpressionItemTypeAnd:
		return fromBool(toBool(resLeft) && toBool(resRight)), nil
	case ExpressionItemTypeXor:
		return fromBool(toBool(resLeft) != toBool(resRight)), nil
	case ExpressionItemTypeBitShiftLeft:
		return resLeft << resRight, nil
	case ExpressionItemTypeBitShiftRight:
		return resLeft >> resRight, nil
	case ExpressionItemTypeEquals:
		return fromBool(resLeft == resRight), nil
	case ExpressionItemTypeNotEquals:
		return fromBool(resLeft != resRight), nil
	case ExpressionItemTypeLess:
		return fromBool(resLeft < resRight), nil
	case ExpressionItemTypeLessOrEquals:
		return fromBool(resLeft <= resRight), nil
	case ExpressionItemTypeGreater:
		return fromBool(resLeft > resRight), nil
	case ExpressionItemTypeGreaterOrEquals:
		return fromBool(resLeft >= resRight), nil
	default:
		return res, ErrInvalidParams
	}
}

func (p *Program) parseExpression(scanner *tinybasic.LineScanner) (items []expressionItem) {
	parser := tinybasic.NewLineParserWithScanner(scanner)

	for !scanner.IsEOL() {
		scanner.GetSpaces()

		variableName := parser.GetVariable()
		if variableName != nil {
			value := p.vars.Get(*variableName)
			items = append(items, expressionItem{
				itemType: ExpressionItemTypeValue,
				value:    &value,
			})
			continue
		}

		value := scanner.GetNumber()
		if value != nil {
			items = append(items, expressionItem{
				itemType: ExpressionItemTypeValue,
				value:    value,
			})
			continue
		}

		item := scanner.GetStrings([]string{
			ExpressionItemTypeOperatorPlus,
			ExpressionItemTypeOperatorMinus,
			ExpressionItemTypeOperatorDivision,
			ExpressionItemTypeOperatorMultiply,
			ExpressionItemTypeBracketOpen,
			ExpressionItemTypeBracketClose,
			ExpressionItemTypeOr,
			ExpressionItemTypeXor,
			ExpressionItemTypeAnd,
			ExpressionItemTypeNot,
			ExpressionItemTypeMod,
			ExpressionItemTypeBitShiftLeft,
			ExpressionItemTypeBitShiftRight,
			ExpressionItemTypeEquals,
			ExpressionItemTypeNotEquals,
			ExpressionItemTypeLess,
			ExpressionItemTypeGreater,
			ExpressionItemTypeGreaterOrEquals,
			ExpressionItemTypeLessOrEquals,
		})
		if item != nil {
			items = append(items, expressionItem{
				itemType: expressionItemType(*item),
				value:    nil,
			})
			continue
		}
	}

	return items
}

type expressionItemType string

const (
	ExpressionItemTypeValue            = "value"
	ExpressionItemTypeBracketOpen      = "("
	ExpressionItemTypeBracketClose     = ")"
	ExpressionItemTypeOperatorPlus     = "+"
	ExpressionItemTypeOperatorMinus    = "-"
	ExpressionItemTypeOperatorDivision = "/"
	ExpressionItemTypeOperatorMultiply = "*"
	ExpressionItemTypeOr               = "or"
	ExpressionItemTypeXor              = "xor"
	ExpressionItemTypeAnd              = "and"
	ExpressionItemTypeNot              = "not"
	ExpressionItemTypeMod              = "mod"
	ExpressionItemTypeBitShiftLeft     = "<<"
	ExpressionItemTypeBitShiftRight    = ">>"
	ExpressionItemTypeEquals           = "="
	ExpressionItemTypeNotEquals        = "<>"
	ExpressionItemTypeLess             = "<"
	ExpressionItemTypeGreater          = ">"
	ExpressionItemTypeGreaterOrEquals  = ">="
	ExpressionItemTypeLessOrEquals     = "<="
)

type expressionItem struct {
	itemType expressionItemType
	value    *int
}

func toBool(v int) bool { return v != 0 }
func fromBool(v bool) int {
	if v {
		return 1
	}

	return 0
}

type expressionTreeItem struct {
	item  *expressionItem
	left  *expressionTreeItem
	right *expressionTreeItem
}

func getClosingBracketIndex(items []expressionItem) (index int, err error) {
	level := 0

	for i, item := range items {
		switch item.itemType {
		case ExpressionItemTypeBracketOpen:
			level++
		case ExpressionItemTypeBracketClose:
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

var operationsPriority = map[expressionItemType]int{
	ExpressionItemTypeEquals:          1,
	ExpressionItemTypeNotEquals:       1,
	ExpressionItemTypeLess:            1,
	ExpressionItemTypeGreater:         1,
	ExpressionItemTypeGreaterOrEquals: 1,
	ExpressionItemTypeLessOrEquals:    1,

	ExpressionItemTypeOperatorPlus:  2,
	ExpressionItemTypeOperatorMinus: 2,
	ExpressionItemTypeOr:            2,
	ExpressionItemTypeXor:           2,

	ExpressionItemTypeOperatorMultiply: 3,
	ExpressionItemTypeOperatorDivision: 3,
	ExpressionItemTypeBitShiftLeft:     3,
	ExpressionItemTypeBitShiftRight:    3,
	ExpressionItemTypeMod:              3,
	ExpressionItemTypeAnd:              3,

	ExpressionItemTypeNot: 4,
}

func (p *Program) generateTree(items []expressionItem) (*expressionTreeItem, error) {
	ops := []*expressionItem(nil)
	nodes := []*expressionTreeItem(nil)

	for i, item := range items {
		if item.itemType == ExpressionItemTypeValue {
			nodes = append(nodes, &expressionTreeItem{
				item:  &items[i],
				left:  nil,
				right: nil,
			})
		} else if priority, isOperator := operationsPriority[item.itemType]; isOperator {
			if len(ops) == 0 || priority >= operationsPriority[ops[len(ops)-1].itemType] {
				ops = append(ops, &items[i])
				continue
			}

			node := &expressionTreeItem{
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

			for len(ops) != 0 && priority < operationsPriority[ops[len(ops)-1].itemType] {
				node = &expressionTreeItem{
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
		}
	}

	for len(ops) != 0 {
		node := &expressionTreeItem{
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

//func (p *Program) getExpressionTreeItem(items []expressionItem) (res *expressionTreeItem, err error) {
//	if len(items) == 0 {
//		return nil, nil
//	}
//
//	// value
//	if len(items) == 1 {
//		return &expressionTreeItem{
//			item:  &items[0],
//			left:  nil,
//			right: nil,
//		}, nil
//	}
//
//	// unary operation
//	if len(items) == 2 {
//		res = &expressionTreeItem{
//			item: &items[0],
//			left: nil,
//			right: &expressionTreeItem{
//				item:  &items[1],
//				left:  nil,
//				right: nil,
//			},
//		}
//
//		return res, nil
//	}
//
//	bracketsLevel := 0
//
//	for i := 0; i < len(items); i++ {
//		switch items[i].itemType {
//		case ExpressionItemTypeEquals2,
//			ExpressionItemTypeNotEquals,
//			ExpressionItemTypeLess,
//			ExpressionItemTypeLessOrEquals,
//			ExpressionItemTypeGreater,
//			ExpressionItemTypeGreaterOrEquals:
//			if i == 0 {
//				return nil, ErrInvalidParams
//			}
//
//			res.item = &items[i]
//
//			res.left, err = p.getExpressionTreeItem(items[:i])
//			if err != nil {
//				return res, err
//			}
//
//			res.right, err = p.getExpressionTreeItem(items[i+1:])
//			if err != nil {
//				return res, err
//			}
//
//			return res, nil
//		case ExpressionItemTypeOperatorPlus,
//			ExpressionItemTypeOperatorMinus:
//			if i == 0 {
//				return nil, ErrInvalidParams
//			}
//
//			res.item = &items[i]
//
//			res.left, err = p.getExpressionTreeItem(items[:i])
//			if err != nil {
//				return res, err
//			}
//
//			res.right, err = p.getExpressionTreeItem(items[i+1:])
//			if err != nil {
//				return res, err
//			}
//
//			return res, nil
//		case ExpressionItemTypeOr,
//			ExpressionItemTypeXor:
//		case ExpressionItemTypeBracketOpen:
//			bracketsLevel++
//
//			j := 0
//
//			j, err = getClosingBracketIndex(items[i:])
//			if err != nil {
//				return res, err
//			}
//
//			if j <= i {
//				return res, ErrInvalidParams
//			}
//
//			treeItemInBrackets := (*expressionTreeItem)(nil)
//
//			treeItemInBrackets, err = p.getExpressionTreeItem(items[i+1 : j])
//			if err != nil {
//				return res, err
//			}
//
//			i = j // перепрыгиваем к закрывающей скобке
//
//			if res.left == nil {
//				// например, если парсим (A+B)+C
//				res.left = treeItemInBrackets
//			} else {
//				// например, если парсим C+(A+B) или C+E/(A+B)-D
//			}
//			bracketsLevel++
//		case ExpressionItemTypeBracketClose:
//			if bracketsLevel <= 0 {
//				return res, ErrInvalidParams
//			}
//			bracketsLevel--
//		}
//	}
//
//	return res, nil
//}
