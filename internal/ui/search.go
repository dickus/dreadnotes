// Package ui provides a user interface for searching.
package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/blevesearch/bleve/v2"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dickus/dreadnotes/internal/search"
	"golang.org/x/term"
)

const visibleResults = 6

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
			Foreground(lipgloss.Color("4")).
			Bold(true)

	activeInputStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("15"))

	placeholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240"))

	cursorStyle = lipgloss.NewStyle().
			Bold(true).
			Underline(false)
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

func performSearch(m SearchModel) tea.Cmd {
	return func() tea.Msg {
		hasValidDateFilter := !m.tagMode && len(m.dateStart) >= 10
		queryEmpty := strings.TrimSpace(m.query) == ""

		var res *bleve.SearchResult
		var err error

		limit := 100

		if m.tagMode {
			if queryEmpty {
				req := bleve.NewSearchRequestOptions(bleve.NewMatchAllQuery(), limit, 0, false)
				req.Fields = []string{"title", "content"}
				res, err = m.idx.Search(req)
			} else {
				res, err = search.SearchByTag(m.idx, m.query, limit)
			}
		} else {
			if hasValidDateFilter {
				start, end, parseErr := parseDateRange(m.dateStart, m.dateEnd)
				if parseErr != nil {
					return searchResultMsg{err: fmt.Errorf("invalid date format: use YYYY-MM-DD")}
				}

				dateField := "created"
				if m.searchUpdated {
					dateField = "updated"
				}

				res, err = search.SearchWithDateFilter(m.idx, m.query, start, end, dateField, limit)
			} else {
				if queryEmpty {
					req := bleve.NewSearchRequestOptions(bleve.NewMatchAllQuery(), limit, 0, false)
					req.Fields = []string{"title", "content"}
					res, err = m.idx.Search(req)
				} else {
					res, err = search.Search(m.idx, m.query, limit)
				}
			}
		}

		if err != nil {
			return searchResultMsg{err: err}
		}

		queryLower := strings.ToLower(strings.TrimSpace(m.query))
		items := make([]resultItem, 0, len(res.Hits))

		for _, hit := range res.Hits {
			title, _ := hit.Fields["title"].(string)
			content, _ := hit.Fields["content"].(string)

			var snippet string
			if !m.tagMode && !queryEmpty {
				snippet = findMatchingLine(content, queryLower)
			}

			if snippet == "" {
				snippet = contentPreview(content, 5)
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

// parseDateRange is a helper function to properly parse dates.
func parseDateRange(startStr, endStr string) (time.Time, time.Time, error) {
	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	var end time.Time
	if strings.TrimSpace(endStr) == "" || len(endStr) < 10 {
		end = start.Add(24*time.Hour - time.Nanosecond)
	} else {
		end, err = time.Parse("2006-01-02", endStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}

		end = end.Add(24*time.Hour - time.Nanosecond)
	}

	return start, end, nil
}

func contentPreview(content string, n int) string {
	var lines []string
	var contentCount int
	var lastWasEmpty bool

	for line := range strings.SplitSeq(content, "\n") {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			if lastWasEmpty || len(lines) == 0 {
				continue
			}

			lines = append(lines, "")
			lastWasEmpty = true

			continue
		}

		lastWasEmpty = false

		runes := []rune(trimmed)
		if len(runes) > 120 {
			trimmed = string(runes[:120]) + "…"
		}

		width := getTermWidth() - 4
		lines = append(lines, wrapLine(trimmed, width))

		contentCount++
		if contentCount >= n {
			break
		}
	}

	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	if len(lines) == 0 {
		return ""
	}

	return strings.Join(lines, "\n")
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

	dateStart     string
	dateEnd       string
	focusIndex    int
	searchUpdated bool

	results       []resultItem
	cursor        int
	err           error
	chosen        string
	viewportStart int
}

func NewSearchModel(idx bleve.Index, tagMode bool) SearchModel {
	return SearchModel{
		idx:     idx,
		tagMode: tagMode,
	}
}

func (m SearchModel) Init() tea.Cmd { return performSearch(m) }

func (m SearchModel) Chosen() string { return m.chosen }

func (m SearchModel) updateViewport() SearchModel {
	if len(m.results) == 0 {
		m.viewportStart = 0

		return m
	}

	centerOffset := visibleResults / 2
	maxStart := max(0, len(m.results)-visibleResults)

	m.viewportStart = max(0, min(m.cursor-centerOffset, maxStart))

	return m
}

func (m SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if msg.Alt && msg.Type == tea.KeyRunes {
			switch string(msg.Runes) {
			case "j":
				if len(m.results) > 0 {
					m.cursor = (m.cursor + 1) % len(m.results)
					m = m.updateViewport()
				}

				return m, nil

			case "k":
				if len(m.results) > 0 {
					m.cursor = (m.cursor - 1 + len(m.results)) % len(m.results)
					m = m.updateViewport()
				}

				return m, nil

			case "d":
				if !m.tagMode {
					m.searchUpdated = !m.searchUpdated

					return m, performSearch(m)
				}
			}
		}

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyTab:
			if !m.tagMode {
				m.focusIndex = (m.focusIndex + 1) % 3

				return m, nil
			}

		case tea.KeyEnter:
			if len(m.results) > 0 {
				m.chosen = m.results[m.cursor].path

				return m, tea.Quit
			}

		case tea.KeyBackspace:
			changed := false
			if m.focusIndex == 0 && len(m.query) > 0 {
				m.query = string([]rune(m.query)[:len([]rune(m.query))-1])
				changed = true
			} else if m.focusIndex == 1 && len(m.dateStart) > 0 {
				m.dateStart = string([]rune(m.dateStart)[:len([]rune(m.dateStart))-1])
				changed = true
			} else if m.focusIndex == 2 && len(m.dateEnd) > 0 {
				m.dateEnd = string([]rune(m.dateEnd)[:len([]rune(m.dateEnd))-1])
				changed = true
			}
			if changed {
				m.cursor = 0
				m.viewportStart = 0

				return m, performSearch(m)
			}

		case tea.KeySpace:
			if m.focusIndex == 0 {
				m.query += " "
				m.cursor = 0
				m.viewportStart = 0

				return m, performSearch(m)
			}

		case tea.KeyRunes:
			if m.focusIndex == 0 {
				m.query += string(msg.Runes)
				m.cursor = 0
				m.viewportStart = 0

				return m, performSearch(m)
			} else {
				char := string(msg.Runes)
				if strings.ContainsAny(char, "0123456789-") {
					if m.focusIndex == 1 {
						m.dateStart += char
					} else if m.focusIndex == 2 {
						m.dateEnd += char
					}

					m.cursor = 0
					m.viewportStart = 0

					return m, performSearch(m)
				}
			}
		}

	case searchResultMsg:
		m.err = msg.err
		if msg.err == nil {
			m.results = msg.items
		} else {
			m.results = nil
		}

		if m.cursor >= len(m.results) {
			m.cursor = max(0, len(m.results)-1)
		}

		m = m.updateViewport()
	}

	return m, nil
}

func (m SearchModel) View() string {
	var b strings.Builder

	cursor := cursorStyle.Render("_")

	if m.tagMode {
		b.WriteString(promptStyle.Render("  Tag: "))
		if m.focusIndex == 0 {
			b.WriteString(activeInputStyle.Render(m.query) + cursor + "\n\n")
		} else {
			b.WriteString(m.query + "\n\n")
		}
	} else {
		targetName := "Created"
		if m.searchUpdated {
			targetName = "Updated"
		}

		b.WriteString(promptStyle.Render("  Search : "))
		if m.focusIndex == 0 {
			b.WriteString(activeInputStyle.Render(m.query) + cursor + "\n")
		} else {
			b.WriteString(m.query + "\n")
		}

		renderDate := func(val string, focused bool) string {
			if val == "" {
				if focused {
					firstY := lipgloss.NewStyle().
						Bold(true).
						Underline(true).
						Render(" ")

					return firstY + placeholderStyle.Render("YYY-MM-DD")
				}

				return placeholderStyle.Render("YYYY-MM-DD")
			}
			if focused {
				return activeInputStyle.Render(val) + cursor
			}

			return activeInputStyle.Render(val)
		}

		b.WriteString(promptStyle.Render(fmt.Sprintf("  %-7s: ", targetName)))
		b.WriteString(renderDate(m.dateStart, m.focusIndex == 1) + "\n")

		b.WriteString(promptStyle.Render("  To     : "))
		b.WriteString(renderDate(m.dateEnd, m.focusIndex == 2) + "\n\n")
	}

	if m.err != nil {
		b.WriteString(fmt.Sprintf("  Error: %v\n", m.err))

		return b.String()
	}

	if len(m.results) == 0 {
		b.WriteString("  No results.\n")

		return b.String()
	}

	start := m.viewportStart
	end := min(start+visibleResults, len(m.results))
	slice := m.results[start:end]

	for i, r := range slice {
		actualIndex := start + i
		if actualIndex == m.cursor {
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
