package tinybasic

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
