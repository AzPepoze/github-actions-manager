package components

import (
	"bufio"
	"github-actions-manager/internal/ui/shared"
	"io"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type TaskModel struct {
	title    string
	message  string
	cmd      *exec.Cmd
	onDone   func() tea.Msg
	logs     []string
	finished bool
	program  *tea.Program
}

type logLineMsg string
type taskFinishedMsg struct{}

func NewTaskModel() *TaskModel {
	return &TaskModel{}
}

func (m *TaskModel) SetTask(title string, message string, cmd *exec.Cmd, onDone func() tea.Msg) {
	m.title = title
	m.message = message
	m.cmd = cmd
	m.onDone = onDone
	m.logs = []string{}
	m.finished = false
}

func (m *TaskModel) SetProgram(p *tea.Program) {
	m.program = p
}

func (m *TaskModel) Init() tea.Cmd {
	return m.runAndStream()
}

func (m *TaskModel) runAndStream() tea.Cmd {
	return func() tea.Msg {
		stdout, _ := m.cmd.StdoutPipe()
		stderr, _ := m.cmd.StderrPipe()

		if err := m.cmd.Start(); err != nil {
			m.program.Send(logLineMsg("Error: " + err.Error()))
			return taskFinishedMsg{}
		}

		scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
		for scanner.Scan() {
			m.program.Send(logLineMsg(scanner.Text()))
		}
		_ = m.cmd.Wait()
		return taskFinishedMsg{}
	}
}

func (m *TaskModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case logLineMsg:
		m.logs = append(m.logs, string(msg))
		return m, nil
	case taskFinishedMsg:
		m.finished = true
		return m, nil
	case tea.KeyMsg:
		if m.finished {
			if m.onDone != nil {
				return m, func() tea.Msg { return m.onDone() }
			}
			return m, func() tea.Msg { return shared.NavigateToMsg{Screen: shared.ScreenDashboard} }
		}
	}
	return m, nil
}

func (m *TaskModel) View() string {
	var body strings.Builder
	title := "Task Progress"
	if m.title != "" {
		title = m.title
	}

	body.WriteString(shared.SubPromptStyle.Render(m.message))
	body.WriteString("\n\n")

	start := 0
	if len(m.logs) > 10 {
		start = len(m.logs) - 10
	}
	for _, line := range m.logs[start:] {
		body.WriteString(shared.LogStyle.Render(line) + "\n")
	}

	if m.finished {
		body.WriteString("\n" + shared.SuccessStyle.Render("Task Complete!") + "\n")
		body.WriteString(shared.HintStyle.Render("Press any key to return to Dashboard"))
	} else {
		body.WriteString("\n" + shared.HintStyle.Render("Running command..."))
	}

	return shared.RenderPage(title, body.String(), "")
}
