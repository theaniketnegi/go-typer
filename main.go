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
	QUOTE  = "Hello world"
)

var (
	QUOTE_STYLE = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("63")).Width(70)

	CORRECT_STYLE      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#AFE1AF"))
	INCORRECT_STYLE    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#C70039"))
	MISSING_UNDERLINED = lipgloss.NewStyle().Underline(true)
)

type model struct {
	currentQuoteWords []string
	letterTracker     []string
	cursor            int
	curWord           int
	isFinished        bool
	accumulatedLen    int
	totalKeystrokes   int
	correctKeystrokes int
	missingKeystrokes int
	accuracy          float32
}

func initialModel() model {
	currQuoteWords := strings.Split(QUOTE, " ")
	icWords := make([]string, len(currQuoteWords))

	return model{
		currentQuoteWords: currQuoteWords,
		letterTracker:     icWords,
		cursor:            0,
		curWord:           0,
		isFinished:        false,
		accumulatedLen:    0,
		totalKeystrokes:   0,
		correctKeystrokes: 0,
		missingKeystrokes: 0,
		accuracy:          0.,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen)
}

func (m model) EvaluateResult(typedLastLetter bool) model {
	for i := 0; i < len(m.currentQuoteWords); i++ {
		for idx := range m.letterTracker[i] {
			if idx >= len(m.currentQuoteWords[i]) {
				break
			}
			if m.currentQuoteWords[i][idx] == m.letterTracker[i][idx] {
				m.correctKeystrokes++
			}
		}

		letterTrackerLen := len(m.letterTracker[i])
		if typedLastLetter && i == len(m.currentQuoteWords)-1 {
			letterTrackerLen++ //+1 for last letter of last word
		}

		if len(m.currentQuoteWords[i]) > letterTrackerLen {
			m.missingKeystrokes += (len(m.currentQuoteWords[i]) - len(m.letterTracker[i]))
		}
	}
	m.isFinished = true
	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyBackspace:
			if m.isFinished {
				return m, nil
			}
			if m.cursor > 0 {
				if len(m.letterTracker[m.curWord]) > 0 {
					m.letterTracker[m.curWord] = m.letterTracker[m.curWord][:len(m.letterTracker[m.curWord])-1]
					m.cursor--
				} else if m.curWord > 0 {
					m.curWord--
					m.accumulatedLen -= max(len(m.currentQuoteWords[m.curWord]), len(m.letterTracker[m.curWord])) + 1
					m.cursor = m.accumulatedLen
					if len(m.letterTracker[m.curWord]) > 0 {
						m.cursor += len(m.letterTracker[m.curWord])
					}
				}
			}
			return m, nil

		case tea.KeySpace:
			if m.isFinished || len(m.letterTracker[m.curWord]) == 0 {
				return m, nil
			}

			if m.curWord < len(m.currentQuoteWords)-1 {
				m.cursor += m.accumulatedLen + max(len(m.letterTracker[m.curWord]), len(m.currentQuoteWords[m.curWord])) - m.cursor + 1
				m.accumulatedLen += max(len(m.letterTracker[m.curWord]), len(m.currentQuoteWords[m.curWord])) + 1
				m.curWord++
			} else {
				m = m.EvaluateResult(false)
				m.accuracy = float32(m.correctKeystrokes) / float32(m.totalKeystrokes+m.missingKeystrokes) * 100
			}

			return m, nil
		}

		if m.isFinished {
			return m, nil
		}
		if len(msg.String()) == 1 {
			if len(m.letterTracker[m.curWord]) > len(m.currentQuoteWords[m.curWord])+3 {
				return m, nil
			}
			m.totalKeystrokes++

			if m.curWord == len(m.currentQuoteWords)-1 && m.cursor == m.accumulatedLen+len(m.currentQuoteWords[m.curWord])-1 && msg.String() == string(m.currentQuoteWords[m.curWord][len(m.currentQuoteWords[m.curWord])-1]) {
				m = m.EvaluateResult(true)
				m.correctKeystrokes += 1 // For the final character
				m.accuracy = float32(m.correctKeystrokes) / float32(m.totalKeystrokes+m.missingKeystrokes) * 100
			}
			m.letterTracker[m.curWord] += msg.String()
			m.cursor++
		}
	}
	return m, nil
}

func (m model) View() string {
	var s strings.Builder
	if !m.isFinished {
		joinedStr := strings.Trim(strings.Join(m.currentQuoteWords, " "), " ")

		var formattedString strings.Builder
		excessLen := 0

		for i := 0; i <= m.curWord; i++ {
			for idx, rn := range m.letterTracker[i] {
				if idx < len(m.currentQuoteWords[i]) {
					if byte(rn) != m.currentQuoteWords[i][idx] {
						formattedString.WriteString(INCORRECT_STYLE.Render(string(m.currentQuoteWords[i][idx])))
					} else {
						formattedString.WriteString(CORRECT_STYLE.Render(string(m.currentQuoteWords[i][idx])))
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
			} else if i != m.curWord {
				formattedString.WriteString(MISSING_UNDERLINED.Render(m.currentQuoteWords[i][len(m.letterTracker[i]):]))
				formattedString.WriteString(" ")
			}
		}

		s.WriteString(QUOTE_STYLE.Render(formattedString.String() + CURSOR + joinedStr[m.cursor-excessLen:]))
	} else {
		s.WriteString(fmt.Sprintf("Accuracy: %.2f%%\nCorrect: %d/Incorrect: %d/Missing: %d", m.accuracy, m.correctKeystrokes, m.totalKeystrokes-m.correctKeystrokes, m.missingKeystrokes))
	}
	return s.String()
}

func main() {
	prog := tea.NewProgram(initialModel())
	if _, err := prog.Run(); err != nil {
		fmt.Printf("There was some error running the program: %v", err)
		os.Exit(1)
	}
}
