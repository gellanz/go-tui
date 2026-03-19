package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

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

func generate_concurrent_candles(ch chan []float64) tea.Cmd {
	// concurrent candle generator
	// ch := make(chan []float64)

	return func() tea.Msg {
		i := 0
		for {
			time.Sleep(500 * time.Millisecond)

			candle := make([]float64, 5)
			scale := rand.Float64()
			candle[0] = 5 * rand.Float64()
			candle[1] = candle[0] + 2*scale
			candle[2] = candle[0] + 4*scale
			candle[3] = candle[2] - 2*scale
			candle[4] = float64(i)
			i++
			ch <- candle
		}
	}
}

func (m model) draw_candles(candles [][]float64) {
	m.c1.Clear()
	graph.DrawXYAxis(&m.c1, m.cursor, axisStyle)

	style := redCandle
	if rand.Intn(2) == 1 {
		style = greenCandle
	}
	for _, candle := range candles {
		graph.DrawCandlestickBottomToTop(&m.c1,
			m.cursor.Add(canvas.Point{X: int(candle[4]), Y: -1}),
			candle[0],
			candle[1],
			candle[2],
			candle[3],
			style,
		)
	}
}

type candleListener []float64

func candleTeaMsg(ch chan []float64) tea.Cmd {
	return func() tea.Msg {
		return candleListener(<-ch)
	}
}

type model struct {
	c1      canvas.Model
	cursor  canvas.Point
	candle  chan []float64
	candles [][]float64
}

func (m model) Init() tea.Cmd {
	return tea.Batch(generate_concurrent_candles(m.candle), candleTeaMsg(m.candle))
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case candleListener:
		c := <-m.candle
		m.candles = append(m.candles, c)
		m.draw_candles(m.candles)
		return m, candleTeaMsg(m.candle)
	}

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
	h := 20
	c1 := canvas.New(w, h)
	cs := [][]float64{}

	m := model{c1, canvas.Point{0, h - 1}, make(chan []float64), cs}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program: ", err)
		os.Exit(1)
	}
}
