package tinybasic

import (
	"fmt"
)

type Variables struct {
	values map[string]int
}

func NewVariables() *Variables {
	return &Variables{
		values: make(map[string]int),
	}
}

func (v *Variables) Set(varName string, varValue int) {
	v.values[varName] = varValue
}

func (v *Variables) Get(varName string) int {
	if value, exists := v.values[varName]; exists {
		return value
	}
	return 0
}

func (v *Variables) Reset() {
	v.values = make(map[string]int)
}

func TestVariables() {
	vars := NewVariables()

	// Установка и получение значения переменной
	vars.Set("A", 10)
	fmt.Println("A =", vars.Get("A")) // Вывод: A = 10

	// Получение несуществующей переменной (должно вернуть 0)
	fmt.Println("B =", vars.Get("B")) // Вывод: B = 0

	// Сброс всех переменных
	vars.Reset()
	fmt.Println("A after reset =", vars.Get("A")) // Вывод: A after Reset = 0
}
