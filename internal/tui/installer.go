package tui

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/cjairm/devgita/pkg/utils"
)

type Step struct {
	Label string
	Run   func() error
}

type stepMsg int

type stepErrorMsg struct {
	stepIndex int
	err       error
	logs      []string
}

type stepSuccessMsg struct {
	stepIndex int
	logs      []string
}

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
	currentStep   int
	showRetry     bool
	logs          []string
}

var (
	boxWidth     = 70
	boxHeight    = 18
	logBoxHeight = 6

	fgColor = lipgloss.Color("#FFFFFF")
	bgColor = lipgloss.Color("#1E1E1E")
	blue    = lipgloss.Color("#2196F3")
	red     = lipgloss.Color("#F44336")
	green   = lipgloss.Color("#4CAF50")
)

var logo = `
    .___                .__  __
  __| _/_______  ______ |__|/  |______
 / __ |/ __ \  \/ / ___\|  \   __\__  \
/ /_/ \  ___/\   / /_/  >  ||  |  / __ \_
\____ |\___  >\_/\___  /|__||__| (____  /
     \/    \/   /_____/               \/
@cjairm
`

var welcomeMessage = `
Welcome to DevGita installer!

Use arrow keys to select Install or Cancel.
Press Enter to confirm.
Press Ctrl+C anytime to exit.
`

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
		showRetry:     false,
		logs:          []string{},
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.cancel()
			return m, tea.Quit
		case "up":
			if !m.loading && !m.showRetry && m.selectedIndex > 0 {
				m.selectedIndex--
			}
			return m, nil
		case "down":
			if !m.loading && !m.showRetry && m.selectedIndex < 1 {
				m.selectedIndex++
			}
			return m, nil
		case "enter":
			if !m.loading && !m.done && !m.showRetry {
				if m.selectedIndex == 0 {
					m.loading = true
					m.done = false
					m.err = nil
					m.currentStep = -1
					m.choice = ""
					m.logs = []string{}
					return m, m.runStepCmd(0)
				}
				m.cancel()
				return m, tea.Quit
			}
			if m.done {
				return m, tea.Quit
			}
			if m.showRetry {
				m.loading = true
				m.err = nil
				m.showRetry = false
				m.currentStep = -1
				m.logs = []string{}
				return m, m.runStepCmd(0)
			}
			return m, nil
		case "r", "R":
			if m.showRetry {
				m.loading = true
				m.err = nil
				m.showRetry = false
				m.currentStep = -1
				m.logs = []string{}
				return m, m.runStepCmd(0)
			}
			return m, nil
		case "c", "C":
			if m.showRetry {
				m.cancel()
				return m, tea.Quit
			}
			return m, nil
		}

	case stepMsg:
		m.currentStep = int(msg)
		m.logs = []string{}
		return m, nil

	case stepErrorMsg:
		m.currentStep = msg.stepIndex
		m.err = msg.err
		m.loading = false
		m.showRetry = true
		m.logs = msg.logs
		return m, nil

	case stepSuccessMsg:
		m.currentStep = msg.stepIndex
		m.logs = msg.logs
		if m.currentStep+1 < len(m.steps) {
			return m, m.runStepCmd(m.currentStep + 1)
		}
		return m, func() tea.Msg { return stepsDoneMsg{} }

	case stepsDoneMsg:
		m.loading = false
		m.done = true
		m.choice = "✅ Installation completed successfully!"
		m.currentStep = -1
		m.showRetry = false
		return m, nil

	case error:
		m.err = msg
		m.loading = false
		m.currentStep = -1
		m.showRetry = true
		m.logs = append(m.logs, "Error: "+msg.Error())
		return m, nil
	}
	return m, nil
}

func (m model) runStepCmd(i int) tea.Cmd {
	return func() tea.Msg {
		var buf bytes.Buffer
		utils.Logger = &buf

		// Log reset here ensures clean logs for each step
		fmt.Fprintf(&buf, "🔧 Running step %d: %s\n", i+1, m.steps[i].Label)

		err := m.steps[i].Run()

		utils.Logger = nil

		logs := strings.Split(strings.TrimSpace(buf.String()), "\n")

		if err != nil {
			logs = append(logs, fmt.Sprintf("❌ Error on step %d: %s", i+1, err.Error()))
			return stepErrorMsg{stepIndex: i, err: err, logs: logs}
		}

		logs = append(logs, fmt.Sprintf("✅ Step %d completed: %s", i+1, m.steps[i].Label))
		return stepSuccessMsg{stepIndex: i, logs: logs}
	}
}

func (m model) View() string {
	termWidth, termHeight := m.width, m.height
	if termWidth == 0 || termHeight == 0 {
		termWidth, termHeight = 80, 24
	}

	mainBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(blue).
		Padding(1, 2).
		Width(boxWidth).
		Height(boxHeight).
		Background(bgColor).
		Foreground(fgColor).
		Align(lipgloss.Left)

	logBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(blue).
		Padding(1, 2).
		Width(termWidth - 2).
		Height(logBoxHeight).
		Background(bgColor).
		Foreground(fgColor).
		Align(lipgloss.Left)

	if !(m.loading || m.done || m.showRetry) {
		logoStyled := lipgloss.NewStyle().
			Foreground(blue).
			MarginBottom(1).
			Render(logo)

		welcomeStyled := lipgloss.NewStyle().
			Foreground(fgColor).
			Align(lipgloss.Left).
			MarginBottom(1).
			Render(strings.TrimSpace(welcomeMessage))

		content := logoStyled + "\n" + welcomeStyled + "\n"
		options := []string{"Install", "Cancel"}
		for i, opt := range options {
			cursor := " "
			if m.selectedIndex == i {
				cursor = ">"
			}
			content += fmt.Sprintf("%s %s\n", cursor, opt)
		}

		return lipgloss.Place(termWidth, termHeight,
			lipgloss.Center, lipgloss.Center, content)
	}

	var stepsContent strings.Builder
	for i, step := range m.steps {
		var status string
		if i < m.currentStep {
			status = "✅"
		} else if i == m.currentStep {
			status = "💻"
		} else {
			status = "😴"
		}
		stepsContent.WriteString(fmt.Sprintf("%s %s\n", status, step.Label))
	}

	if m.choice != "" {
		stepsContent.WriteString("\n" + m.choice + "\n")
	}

	if m.err != nil && !m.showRetry {
		errMsg := lipgloss.NewStyle().Foreground(red).Bold(true).Render("Error: " + m.err.Error())
		stepsContent.WriteString("\n" + errMsg + "\n")
	}

	if m.showRetry {
		errMsg := lipgloss.NewStyle().Foreground(red).Bold(true).Render("Error: " + m.err.Error())
		stepsContent.WriteString("\n" + errMsg + "\n\nPress R to retry or C to cancel.\n")
	}

	mainBox := mainBoxStyle.Render(stepsContent.String())
	logsContent := strings.Join(m.logs, "\n")
	logBox := logBoxStyle.Render(logsContent)

	finalView := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.Place(termWidth, boxHeight+4, lipgloss.Center, lipgloss.Center, mainBox),
		logBox)

	return lipgloss.Place(termWidth, termHeight,
		lipgloss.Center, lipgloss.Center, finalView)
}

// steps := []tui.Step{
// 	{
// 		Label: "Validating OS version",
// 		Run: func() error {
// 			time.Sleep(1 * time.Second)
// 			utils.Log("- Pre-install steps")
// 			return osCmd.ValidateOSVersion()
// 		},
// 	},
// 	{
// 		Label: "Checking disk space",
// 		Run: func() error {
// 			time.Sleep(1 * time.Second)
// 			return nil
// 		},
// 	},
// 	{
// 		Label: "Setting up environment",
// 		Run: func() error {
// 			time.Sleep(2 * time.Second)
// 			return fmt.Errorf("this is my error")
// 		},
// 	},
// }
//
// if err := tui.Run(context.Background(), steps); err != nil {
// 	// fallback or error print
// 	panic(err)
// }
