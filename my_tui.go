package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	number_1 int
	number_2 int
}

func initialModel() model {
	return model{
		number_1: 0,
		number_2: 1,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "1":
			m.number_1++
		case "2":
			m.number_2++
		}
	}
	return m, nil
}

func (m model) View() string {
	instructions := "Press 1 or 2 to increase their counters :) \n"
	counters := fmt.Sprintf("Counter 1: %s \nCounter 2: %s", m.number_1, m.number_2)
	s := instructions + counters
	return s
}

func main() {
	fmt.Println("Hello world")
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("There has been an error: %v", err)
		os.Exit(1)
	}
}
