package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/blevesearch/bleve/v2"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dickus/dreadnotes/internal/search"
	"golang.org/x/term"
)

var (
	activeTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")).
			Bold(true)

	inactiveTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	snippetStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("250")).
			PaddingLeft(4)

	separator = lipgloss.NewStyle().
			Foreground(lipgloss.Color("238")).
			Render("  ──────────────────────────────")

	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("69")).
			Bold(true)
)

type searchResultMsg struct {
	items []resultItem
	err   error
}

type resultItem struct {
	title   string
	path    string
	score   float64
	snippet string
}

func performSearch(idx bleve.Index, query string, tagMode bool) tea.Cmd {
	return func() tea.Msg {
		if strings.TrimSpace(query) == "" {
			return searchResultMsg{}
		}

		var res *bleve.SearchResult
		var err error

		if tagMode {
			res, err = search.SearchByTag(idx, query, 20)
		} else {
			res, err = search.Search(idx, query, 20)
		}
		if err != nil {
			return searchResultMsg{err: err}
		}

		queryLower := strings.ToLower(strings.TrimSpace(query))
		items := make([]resultItem, 0, len(res.Hits))

		for _, hit := range res.Hits {
			title, _ := hit.Fields["title"].(string)
			content, _ := hit.Fields["content"].(string)

			var snippet string
			if tagMode {
				snippet = contentPreview(content, 2)
			} else {
				snippet = findMatchingLine(content, queryLower)

				if snippet == "" {
					snippet = contentPreview(content, 2)
				}
			}

			items = append(items, resultItem{
				title:   title,
				path:    hit.ID,
				score:   hit.Score,
				snippet: snippet,
			})
		}

		return searchResultMsg{items: items}
	}
}

func contentPreview(content string, n int) string {
	var lines []string
	for line := range strings.SplitSeq(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		runes := []rune(trimmed)
		if len(runes) > 120 {
			trimmed = string(runes[:120]) + "…"
		}

		width := getTermWidth() - 4
		lines = append(lines, wrapLine(trimmed, width))

		if len(lines) >= n {
			break
		}
	}

	if len(lines) == 0 {
		return ""
	}

	return strings.Join(lines, "\n\n")
}

func getTermWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		return 80
	}

	return w
}

func wrapLine(s string, width int) string {
	runes := []rune(s)
	if len(runes) <= width {
		return s
	}

	var b strings.Builder
	for i, r := range runes {
		if i > 0 && i%width == 0 {
			b.WriteRune('\n')
		}
		b.WriteRune(r)
	}

	return b.String()
}

func findMatchingLine(content string, query string) string {
	for line := range strings.SplitSeq(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if strings.Contains(strings.ToLower(trimmed), query) {
			runes := []rune(trimmed)
			if len(runes) > 120 {
				trimmed = string(runes[:120]) + "…"
			}
			width := getTermWidth() - 4

			return wrapLine(trimmed, width)
		}
	}

	return ""
}

type SearchModel struct {
	idx     bleve.Index
	tagMode bool
	query   string
	results []resultItem
	cursor  int
	err     error
	chosen  string
}

func NewSearchModel(idx bleve.Index, tagMode bool) SearchModel {
	return SearchModel{
		idx:     idx,
		tagMode: tagMode,
	}
}

func (m SearchModel) Init() tea.Cmd { return nil }

func (m SearchModel) Chosen() string { return m.chosen }

func (m SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if msg.Type == tea.KeyRunes && msg.Alt {
			switch string(msg.Runes) {
			case "j":
				if len(m.results) > 0 {
					m.cursor = (m.cursor + 1) % len(m.results)
				}

				return m, nil

			case "k":
				if len(m.results) > 0 {
					m.cursor = (m.cursor - 1 + len(m.results)) % len(m.results)
				}

				return m, nil
			}
		}

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyBackspace:
			if len(m.query) > 0 {
				runes := []rune(m.query)
				m.query = string(runes[:len(runes)-1])

				return m, performSearch(m.idx, m.query, m.tagMode)
			}

		case tea.KeyEnter:
			if len(m.results) > 0 {
				m.chosen = m.results[m.cursor].path

				return m, tea.Quit
			}

		case tea.KeySpace:
			m.query += " "
			m.cursor = 0

			return m, performSearch(m.idx, m.query, m.tagMode)

		case tea.KeyRunes:
			m.query += string(msg.Runes)
			m.cursor = 0

			return m, performSearch(m.idx, m.query, m.tagMode)
		}

	case searchResultMsg:
		m.err = msg.err
		m.results = msg.items
		if m.cursor >= len(m.results) {
			m.cursor = max(0, len(m.results)-1)
		}
	}

	return m, nil
}

func (m SearchModel) View() string {
	var b strings.Builder

	label := "Search"
	if m.tagMode {
		label = "Tag"
	}

	b.WriteString(promptStyle.Render(fmt.Sprintf("  %s: ", label)))
	b.WriteString(fmt.Sprintf("%s▁\n\n", m.query))

	if m.err != nil {
		b.WriteString(fmt.Sprintf("  Error: %v\n", m.err))

		return b.String()
	}

	if len(m.results) == 0 && m.query != "" {
		b.WriteString("  No results.\n")

		return b.String()
	}

	for i, r := range m.results {
		if i == m.cursor {
			b.WriteString("❯ " + activeTitle.Render(r.title) + "\n")

			if r.snippet != "" {
				b.WriteString(snippetStyle.Render(r.snippet) + "\n")
			}

			b.WriteString(separator + "\n")
		} else {
			b.WriteString("  " + inactiveTitle.Render(r.title) + "\n")
		}
	}

	return b.String()
}
