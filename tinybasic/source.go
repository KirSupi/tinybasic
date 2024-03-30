package tinybasic

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Line struct {
	Label int
	Text  string
}
type Source struct {
	Lines []Line
}

func NewSource() *Source {
	return &Source{}
}

func (s *Source) Clear() {
	s.Lines = nil
}
func (s *Source) Save(filename string) error {
	var builder strings.Builder
	for _, line := range s.Lines {
		builder.WriteString(fmt.Sprintf("%d %s\n", line.Label, line.Text))
	}
	return os.WriteFile(filename, []byte(builder.String()), 0644)
}
func (s *Source) Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("error on file.Close():", err.Error())
		}
	}(file)

	scanner := bufio.NewScanner(file)
	s.Lines = nil
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), " ", 2)
		if len(parts) < 2 {
			continue
		}

		label, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}

		s.Insert(Line{Label: label, Text: parts[1]})
	}
	return scanner.Err()
}
func (s *Source) LoadLine(line string) error {
	parts := strings.SplitN(line, " ", 2)
	if len(parts) < 2 {
		return fmt.Errorf("invalid line format: %s", line)
	}

	label, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("invalid label format: %s", parts[0])
	}

	s.Insert(Line{Label: label, Text: parts[1]})
	return nil
}
func (s *Source) Insert(newLine Line) {
	index := sort.Search(len(s.Lines), func(i int) bool {
		return s.Lines[i].Label >= newLine.Label
	})
	if index < len(s.Lines) && s.Lines[index].Label == newLine.Label {
		s.Lines[index] = newLine
	} else {
		s.Lines = append(s.Lines[:index], append([]Line{newLine}, s.Lines[index:]...)...)
	}
}

func TestSource() {
	source := NewSource()

	err := source.Load("example.txt")
	if err != nil {
		fmt.Println("error on source.Load:", err.Error())
	}

	source.Insert(Line{Label: 15, Text: "LET a=1"})
	fmt.Println(source.Lines)

	err = source.Save("modified_example.txt")
	if err != nil {
		fmt.Println("error on source.Save:", err.Error())
	}

	if err != nil {
		fmt.Println("error on source.Load:", err.Error())
	}

}
