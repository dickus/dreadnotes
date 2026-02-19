package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	pickerSelected = lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")).
			Bold(true)

	pickerNormal = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	pickerTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("69")).
			Bold(true).
			MarginBottom(1)
)

type templateItem struct {
	name string
	path string
}

type TemplatePickerModel struct {
	items    []templateItem
	cursor   int
	chosen   string
	quitting bool
}

func NewTemplatePicker(templatesDir string) (TemplatePickerModel, error) {
	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		return TemplatePickerModel{}, fmt.Errorf("Couldn't read templates dir: %w", err)
	}

	var items []templateItem
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}

		items = append(items, templateItem{
			name: strings.TrimSuffix(e.Name(), ".md"),
			path: filepath.Join(templatesDir, e.Name()),
		})
	}

	if len(items) == 0 {
		return TemplatePickerModel{}, fmt.Errorf("No templates found in %s", templatesDir)
	}

	return TemplatePickerModel{items: items}, nil
}

func (m TemplatePickerModel) Init() tea.Cmd { return nil }

func (m TemplatePickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			m.quitting = true

			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}

		case "enter":
			m.chosen = m.items[m.cursor].path

			return m, tea.Quit
		}
	}

	return m, nil
}

func (m TemplatePickerModel) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	b.WriteString(pickerTitle.Render("  Template:"))
	b.WriteString("\n")

	for i, item := range m.items {
		cursor := "  "
		style := pickerNormal

		if i == m.cursor {
			cursor = "❯ "
			style = pickerSelected
		}

		b.WriteString(cursor + style.Render(item.name) + "\n")
	}

	b.WriteString("\n")

	return b.String()
}

func RunTemplatePicker(templatesDir string) (string, error) {
	conf, _ := os.UserConfigDir()
	confDir := conf + "/" + templatesDir

	m, err := NewTemplatePicker(confDir)
	if err != nil {
		return "", err
	}

	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("Picker error: %w", err)
	}

	final := result.(TemplatePickerModel)
	if final.chosen == "" {
		return "", fmt.Errorf("No template selected")
	}

	return final.chosen, nil
}
