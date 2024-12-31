package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	CURSOR = "|"
	QUOTE  = "No problem can be solved from the same level of consciousness that created it. We must see the world anew. No problem can be solved from the same level of consciousness that created it. We must see the world anew. No problem can be solved from the same level of consciousness that created it. We must see the world anew."
)

var (
	QUOTE_STYLE = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")).Width(70)

	CORRECT_STYLE   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#AFE1AF"))
	INCORRECT_STYLE = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#C70039"))
)

type model struct {
	currentQuoteWords []string
	letterTracker     []string
	cursor            int
	curWord           int
}

func initialModel() model {
	currQuoteWords := strings.Split(QUOTE, " ")
	icWords := make([]string, len(currQuoteWords))

	return model{
		currentQuoteWords: currQuoteWords,
		letterTracker:     icWords,
		cursor:            0,
		curWord:           0,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyBackspace:
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		case tea.KeySpace:
			if m.curWord < len(m.currentQuoteWords)-1 {
				accumulatedLen := 0
				for i := m.curWord - 1; i >= 0; i-- {
					if len(m.letterTracker[i]) > len(m.currentQuoteWords[i]) {
						accumulatedLen += len(m.letterTracker[i]) + 1
					} else {
						accumulatedLen += len(m.currentQuoteWords[i]) + 1
					}
				}
				if len(m.letterTracker[m.curWord]) > len(m.currentQuoteWords[m.curWord]) {
					accumulatedLen += len(m.letterTracker[m.curWord])
				} else {
					accumulatedLen += len(m.currentQuoteWords[m.curWord])
				}
				m.cursor += accumulatedLen - m.cursor + 1
				m.curWord++
			}
			return m, nil
		}
		if len(msg.String()) == 1 {
			if len(m.letterTracker[m.curWord]) > len(m.currentQuoteWords[m.curWord])+3 {
				return m, nil
			}
			m.cursor++
			m.letterTracker[m.curWord] += msg.String()
		}
	}
	return m, nil
}

func (m model) View() string {
	var s strings.Builder
	joinedStr := strings.Trim(strings.Join(m.currentQuoteWords, " "), " ")

	var formattedString strings.Builder
	excessLen := 0

	for i := 0; i <= m.curWord; i++ {
		for idx, rn := range m.letterTracker[i] {
			if idx < len(m.currentQuoteWords[i]) {
				if byte(rn) != m.currentQuoteWords[i][idx] {
					formattedString.WriteString(INCORRECT_STYLE.Render(string(rn)))
				} else {
					formattedString.WriteString(CORRECT_STYLE.Render(string(rn)))
				}
			} else {
				formattedString.WriteString(INCORRECT_STYLE.Render(string(rn)))
			}
		}
		if len(m.letterTracker[i]) >= len(m.currentQuoteWords[i]) {
			excessLen += len(m.letterTracker[i]) - len(m.currentQuoteWords[i])
			if i != m.curWord {
				formattedString.WriteString(" ")
			}
		}

	}

	fmt.Print(excessLen)
	fmt.Print(formattedString.String())
	s.WriteString(QUOTE_STYLE.Render(formattedString.String() + CURSOR + joinedStr[m.cursor-excessLen:]))
	return s.String()
}

func main() {
	prog := tea.NewProgram(initialModel())
	if _, err := prog.Run(); err != nil {
		fmt.Printf("There was some error running the program: %v", err)
		os.Exit(1)
	}
}
