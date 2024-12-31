package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	CURSOR = "█"
)

var (
	QUOTE_STYLE = lipgloss.NewStyle().Bold(true).BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")).Width(70)

	CORRECT_STYLE   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#AFE1AF"))
	INCORRECT_STYLE = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#C70039"))
	CURSOR_STYLE    = lipgloss.NewStyle().Blink(true)
)

type model struct {
	curQuote             string
	cursor               int
	incorrectIndexLength [][2]int
}

func initialModel() model {
	quote := `No problem can be solved from the same level of consciousness that created it. We must see the world anew. No problem can be solved from the same level of consciousness that created it. We must see the world anew. No problem can be solved from the same level of consciousness that created it. We must see the world anew.`
	return model{
		curQuote:             quote,
		cursor:               0,
		incorrectIndexLength: [][2]int{},
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen)
}

func (m model) View() string {
	var s strings.Builder

	if len(m.incorrectIndexLength) > 0 {
		var tmp strings.Builder
		startIdx := 0
		for _, item := range m.incorrectIndexLength {
			index := item[0]
			length := item[1]
			tmp.WriteString(CORRECT_STYLE.Render(m.curQuote[startIdx:index]) + INCORRECT_STYLE.Render(m.curQuote[index:index+length]))
			startIdx = index + length
		}

		s.WriteString(QUOTE_STYLE.Render(tmp.String() + CORRECT_STYLE.Render(m.curQuote[startIdx:m.cursor]) + CURSOR_STYLE.Render(CURSOR) + m.curQuote[min(len(m.curQuote)-1, m.cursor+1):]))
	} else {
		if m.cursor > 0 {
			s.WriteString(QUOTE_STYLE.Render(CORRECT_STYLE.Render(m.curQuote[:m.cursor]) + CURSOR_STYLE.Render(CURSOR) + m.curQuote[min(len(m.curQuote)-1, m.cursor+1):]))
		} else {
			s.WriteString(QUOTE_STYLE.Render(m.curQuote[:m.cursor] + CURSOR_STYLE.Render(CURSOR) + m.curQuote[min(len(m.curQuote)-1, m.cursor+1):]))
		}
	}
	return s.String()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c":
			return m, tea.Quit
		case "backspace":
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		}
		if len(msg.String()) > 1 {
			return m, nil
		}

		if byte(msg.String()[0]) != m.curQuote[m.cursor] {
			if len(m.incorrectIndexLength) == 0 {
				m.incorrectIndexLength = append(m.incorrectIndexLength, [2]int{m.cursor, 1})
			} else {
				if m.incorrectIndexLength[len(m.incorrectIndexLength)-1][0] == m.cursor-1 {
					m.incorrectIndexLength[len(m.incorrectIndexLength)-1][1]++
				} else {
					m.incorrectIndexLength = append(m.incorrectIndexLength, [2]int{m.cursor, 1})
				}
			}
		}
		m.cursor++
		if m.cursor == len(m.curQuote) {
			return m, tea.Quit
		}
	}

	return m, nil
}

func main() {
	prog := tea.NewProgram(initialModel())
	if _, err := prog.Run(); err != nil {
		fmt.Printf("There was some error running the program: %v", err)
		os.Exit(1)
	}
}
