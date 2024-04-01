package arithmetics

import (
	"fmt"
)

// Binary выполняет бинарные арифметические операции (+, -, *, /) над целыми числами
func Binary(a int, op rune, b int) int {
	switch op {
	case '+':
		return a + b
	case '-':
		return a - b
	case '*':
		return a * b
	case '/':
		// При делении используем целочисленное деление, дробная часть отбрасывается
		if b == 0 {
			fmt.Println("Error: division by zero")
			return 0
		}
		return a / b
	default:
		fmt.Println("Error: unsupported operation")
		return 0
	}
}

// Unary выполняет унарные операции, в данном случае поддерживается только унарный минус
func Unary(op rune, a int) int {
	if op == '-' {
		return -a
	}
	fmt.Println("Error: unsupported operation")
	return 0
}
