package tui

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/tergel-sama/kssh/internal/models"
)

type formModel struct {
	title    string
	inputs   []textinput.Model
	focused  int
	done     bool
	save     bool
}

func newFormModel(title string, host models.HostConfig) formModel {
	inputs := make([]textinput.Model, 5)

	inputs[0] = textinput.New()
	inputs[0].SetValue(host.Name)
	inputs[0].Focus()

	inputs[1] = textinput.New()
	inputs[1].SetValue(host.Hostname)

	inputs[2] = textinput.New()
	inputs[2].SetValue(host.User)

	inputs[3] = textinput.New()
	inputs[3].SetValue(strconv.Itoa(host.Port))

	inputs[4] = textinput.New()
	inputs[4].SetValue(host.Key)

	return formModel{
		title:    title,
		inputs:   inputs,
		focused:  0,
	}
}

func (m formModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m formModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.done = true
			m.save = false
			return m, nil
		case "enter":
			m.done = true
			m.save = true
			return m, nil
		case "tab", "shift+tab", "up", "down":
			s := msg.String()

			if s == "up" || s == "shift+tab" {
				m.focused--
			} else {
				m.focused++
			}

			if m.focused >= len(m.inputs) {
				m.focused = 0
			}
			if m.focused < 0 {
				m.focused = len(m.inputs) - 1
			}

			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focused {
					m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}

			return m, nil
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *formModel) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd
	for i := range m.inputs {
		newModel, cmd := m.inputs[i].Update(msg)
		m.inputs[i] = newModel
		cmds = append(cmds, cmd)
	}
	return tea.Batch(cmds...)
}

func (m formModel) View() string {
	var b strings.Builder

	b.WriteString(m.title + "\n\n")

	labels := []string{"Name", "Hostname", "User", "Port", "Key"}

	for i, input := range m.inputs {
		b.WriteString(labels[i] + ": ")
		
		b.WriteString(input.View() + "\n")
	}

	b.WriteString("\nenter to save · esc to cancel\n")

	return b.String()
}

func (m formModel) getHost() models.HostConfig {
	port, _ := strconv.Atoi(m.inputs[3].Value())
	return models.HostConfig{
		Name:     m.inputs[0].Value(),
		Hostname: m.inputs[1].Value(),
		User:     m.inputs[2].Value(),
		Port:     port,
		Key:      m.inputs[4].Value(),
	}
}