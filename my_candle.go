package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/NimbleMarkets/ntcharts/canvas"
	"github.com/NimbleMarkets/ntcharts/canvas/graph"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var defaultStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("63"))

var greenCandle = lipgloss.NewStyle().Foreground(lipgloss.Color("40"))

var redCandle = lipgloss.NewStyle().Foreground(lipgloss.Color("160"))

var axisStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("0"))

func generate_candles() [][]float64 {
	// Simple atm
	matrix := make([][]float64, 15) // initialize an slice of 15 elements

	for i := 0; i <= 14; i++ {
		matrix[i] = make([]float64, 4)
		// low, body_low, high, bofy_high
		matrix[i][0] = float64(rand.Intn(20))
		matrix[i][1] = matrix[i][0] + float64(rand.Intn(3))
		matrix[i][2] = matrix[i][0] + float64(rand.Intn(10))
		matrix[i][3] = matrix[i][2] - float64(rand.Intn(3))
	}

	return matrix
}

func (m model) draw_candles() {
	m.c1.Clear()
	graph.DrawXYAxis(&m.c1, m.cursor, axisStyle)

	candles_slice := generate_candles()

	for i, candle_body := range candles_slice {

		style := redCandle
		if rand.Intn(2) == 1 {
			style = greenCandle
		}
		graph.DrawCandlestickBottomToTop(&m.c1,
			m.cursor.Add(canvas.Point{X: i, Y: -1}),
			candle_body[0],
			candle_body[1],
			candle_body[2],
			candle_body[3],
			style,
		)
	}
}

type model struct {
	c1     canvas.Model
	cursor canvas.Point
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	m.draw_candles()
	return m, nil
}

func (m model) View() string {
	s := "My Candle\n\n"
	s += lipgloss.JoinHorizontal(
		lipgloss.Top,
		defaultStyle.Render(m.c1.View()),
	) + "\n\n"
	return s
}

func main() {
	w := 40
	h := 30
	c1 := canvas.New(w, h)

	m := model{c1, canvas.Point{0, h - 1}}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program: ", err)
		os.Exit(1)
	}
}
