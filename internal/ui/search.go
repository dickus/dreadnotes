package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/blevesearch/bleve/v2"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dickus/dreadnotes/internal/search"
	"golang.org/x/term"
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

func performSearch(idx bleve.Index, query string) tea.Cmd {
	return func() tea.Msg {
		if strings.TrimSpace(query) == "" {
			return searchResultMsg{}
		}

		res, err := search.Search(idx, query, 20)
		if err != nil {
			return searchResultMsg{err: err}
		}

		queryLower := strings.ToLower(strings.TrimSpace(query))

		items := make([]resultItem, 0, len(res.Hits))
		for _, hit := range res.Hits {
			title, _ := hit.Fields["title"].(string)
			content, _ := hit.Fields["content"].(string)

			snippet := findMatchingLine(content, queryLower)

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

func getTermWidth() int {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 { return 80 }

	return w
}

func wrapLine(s string, width int) string {
	runes := []rune(s)
	if len(runes) <= width { return s }

	var b strings.Builder
	for i, r := range runes {
		if i > 0 && i%width == 0 {
			b.WriteRune('\n')
			b.WriteString("    ")
		}
		b.WriteRune(r)
	}

	return b.String()
}

func findMatchingLine(content string, query string) string {
	for line := range strings.SplitSeq(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" { continue }

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
	query   string
	results []resultItem
	cursor  int
	err     error
	chosen  string
}

func NewSearchModel(idx bleve.Index) SearchModel {
	return SearchModel{idx: idx}
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
				return m, performSearch(m.idx, m.query)
			}

		case tea.KeyEnter:
			if len(m.results) > 0 {
				m.chosen = m.results[m.cursor].path

				return m, tea.Quit
			}

		case tea.KeySpace:
			m.query += " "
			m.cursor = 0
			return m, performSearch(m.idx, m.query)

		case tea.KeyRunes:
			m.query += string(msg.Runes)
			m.cursor = 0
			return m, performSearch(m.idx, m.query)
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

	b.WriteString(fmt.Sprintf("  Search: %s▁ \n", m.query))

	if m.err != nil {
		b.WriteString(fmt.Sprintf("  Error: %v\n", m.err))
		return b.String()
	}

	if len(m.results) == 0 && m.query != "" {
		b.WriteString("  No results.\n")
		return b.String()
	}

	for i, r := range m.results {
		cursor := "  "
		if i == m.cursor {
			cursor = "❯ "
		}
		b.WriteString(fmt.Sprintf("%s%s\n", cursor, r.title))
		if r.snippet != "" {
			b.WriteString(fmt.Sprintf("    %s\n", r.snippet))
		}
	}

	return b.String()
}

