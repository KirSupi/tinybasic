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

func TestArithmetics() {
	fmt.Println("Binary(1, '+', 2):", Binary(1, '+', 2))
	fmt.Println("Binary(5, '/', 2):", Binary(5, '/', 2))
	fmt.Println("Unary('-', 1):", Unary('-', 1))

	// Пример с операцией, не поддерживаемой функцией Unary
	fmt.Println("Unary('+', 1):", Unary('+', 1))

	// Пример с недопустимой операцией в Binary
	fmt.Println("Binary(1, '%', 2):", Binary(1, '%', 2))
}
