package tui

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	blue    = lipgloss.Color("#6B72FF")
	bgColor = lipgloss.Color("#1B1B1B")
	fgColor = lipgloss.Color("#C6C6C6")

	boxWidth  = 60
	boxHeight = 15

	logo = ` _____       _     _ _
|  ___|__ __| |___| | |___ _ _
| |_ / -_) _| / -_) | / -_) '_|
|  _|\___\__|_\___|_|_\___|_|
|_|                              `
	welcomeMessage = `Welcome to the Installer!
Please select "Install" to begin the setup.
Press "Cancel" to exit.`
)

// Step represents one installation step.
type Step struct {
	Label string
	Run   func() error
}

type stepMsg int
type stepsDoneMsg struct{}

type model struct {
	steps         []Step
	selectedIndex int
	loading       bool
	err           error
	choice        string
	ctx           context.Context
	cancel        context.CancelFunc
	width, height int
	done          bool
	currentStep   int // -1 means no step running
}

// Run starts the Bubble Tea installer program with given steps.
func Run(ctx context.Context, steps []Step) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	m := model{
		steps:         steps,
		selectedIndex: 0,
		loading:       false,
		choice:        "",
		ctx:           ctx,
		cancel:        cancel,
		currentStep:   -1,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.cancel()
			return m, tea.Quit
		case "up":
			if !m.loading && m.selectedIndex > 0 {
				m.selectedIndex--
			}
			return m, nil
		case "down":
			if !m.loading && m.selectedIndex < 1 {
				m.selectedIndex++
			}
			return m, nil
		case "enter":
			if !m.loading && !m.done {
				if m.selectedIndex == 0 {
					m.loading = true
					m.done = false
					m.err = nil
					m.currentStep = -1
					return m, m.runStepsCmd(0)
				}
				m.cancel()
				return m, tea.Quit
			}
			if m.done {
				return m, tea.Quit
			}
			return m, nil
		}

	case stepMsg:
		m.currentStep = int(msg)
		err := m.steps[m.currentStep].Run()
		if err != nil {
			m.err = err
			m.loading = false
			return m, nil
		}
		if m.currentStep+1 < len(m.steps) {
			return m, m.runStepsCmd(m.currentStep + 1)
		}
		return m, func() tea.Msg { return stepsDoneMsg{} }

	case stepsDoneMsg:
		m.loading = false
		m.done = true
		m.choice = "✅ Installation completed successfully!"
		m.currentStep = -1
		return m, nil

	case error:
		m.err = msg
		m.loading = false
		m.currentStep = -1
		return m, nil
	}

	return m, nil
}

func (m model) runStepsCmd(i int) tea.Cmd {
	return func() tea.Msg {
		return stepMsg(i)
	}
}

func (m model) View() string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(blue).
		Padding(1, 2).
		Width(boxWidth).
		Height(boxHeight).
		Background(bgColor).
		Foreground(fgColor).
		Align(lipgloss.Center)

	logoStyled := lipgloss.NewStyle().
		Foreground(blue).
		MarginBottom(1).
		Render(logo)

	welcomeStyled := lipgloss.NewStyle().
		Foreground(fgColor).
		Align(lipgloss.Left).
		MarginBottom(1).
		Render(strings.TrimSpace(welcomeMessage))

	content := logoStyled + "\n" + welcomeStyled

	if m.loading {
		content += "\n\n"
		for i, step := range m.steps {
			var status string
			if i < m.currentStep {
				status = "✅"
			} else if i == m.currentStep {
				status = "💻"
			} else {
				status = "😴"
			}
			content += fmt.Sprintf("%s %s\n", status, step.Label)
		}
		content += fmt.Sprintf("\n\n%s", m.choice)
		if m.err != nil {
			content += "\n\n" + lipgloss.NewStyle().
				Foreground(lipgloss.Color("1")).
				Render("Error: "+m.err.Error())
		}
	} else if m.done {
		content += fmt.Sprintf("\n\n%s\n\nPress any key to exit.", m.choice)
	} else {
		options := []string{"Install", "Cancel"}
		for i, opt := range options {
			cursor := " "
			if m.selectedIndex == i {
				cursor = ">"
			}
			content += fmt.Sprintf("\n%s %s", cursor, opt)
		}
	}

	termWidth, termHeight := m.width, m.height
	if termWidth == 0 || termHeight == 0 {
		termWidth, termHeight = 80, 24
	}

	return lipgloss.Place(termWidth, termHeight,
		lipgloss.Center, lipgloss.Center,
		style.Render(content))
}

// steps := []tui.Step{
// 	{
// 		Label: "Validating OS version",
// 		Run: func() error {
// 			time.Sleep(2 * time.Second)
// 			return osCmd.ValidateOSVersion()
// 		},
// 	},
// 	{
// 		Label: "Checking disk space",
// 		Run: func() error {
// 			time.Sleep(2 * time.Second)
// 			return nil
// 		},
// 	},
// 	{
// 		Label: "Setting up environment",
// 		Run: func() error {
// 			time.Sleep(2 * time.Second)
// 			return nil
// 		},
// 	},
// }
//
// if err := tui.Run(context.Background(), steps); err != nil {
// 	// fallback or error print
// 	panic(err)
// }
//
