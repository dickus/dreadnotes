// Package doctor anaylzes notes for empty content, incorrect links and duplicates.
package doctor

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dickus/dreadnotes/internal/frontmatter"
	"github.com/dickus/dreadnotes/internal/utils"
)

// Report gathers all problems in one structure
type Report struct {
	BrokenLinks []BrokenLink
	EmptyNotes  []string
	Duplicates  []DuplicateTitle
}

type BrokenLink struct {
	SourceFile string
	TargetNote string
}

type DuplicateTitle struct {
	Title string
	Paths []string
}

type linkRef struct {
	sourceFile string
	rawTarget  string
	normTarget string
}

// Regex for [[]]
var (
	// wikilinkRe ignores aliases if there are any, so it will only work for the actual links
	wikilinkRe = regexp.MustCompile(`\[\[([^\]|]+)(?:\|[^\]]+)?\]\]`)

	// codeBlockRe ignores [[]] in multiline codeblocks
	codeBlockRe = regexp.MustCompile("(?s)```.*?```")

	// inlineCodeRe ignores [[]] in inline code
	inlineCodeRe = regexp.MustCompile("`[^`]*`")
)

type analyzer struct {
	existingTargets map[string]struct{}
	titlesMap       map[string][]string
	collectedLinks  []linkRef
	emptyNotes      []string
}

func newAnalyzer(filesPath string) *analyzer {
	a := &analyzer{
		existingTargets: make(map[string]struct{}),
		titlesMap:       make(map[string][]string),
	}

	a.loadExistingFiles(filesPath)

	return a
}

func (a *analyzer) loadExistingFiles(filesPath string) {
	fileEntries, err := os.ReadDir(filesPath)
	if err != nil {
		return
	}

	for _, entry := range fileEntries {
		if !entry.IsDir() {
			normFileName := strings.ToLower(strings.TrimSpace(entry.Name()))
			a.existingTargets[normFileName] = struct{}{}
		}
	}
}

func (a *analyzer) processNote(fullPath string, entry os.DirEntry) {
	doc, err := frontmatter.ParseFile(fullPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Linter warning: skipping invalid note %s: %v\n", fullPath, err)

		return
	}

	lowerName := strings.ToLower(strings.TrimSpace(entry.Name()))
	baseName := strings.TrimSuffix(lowerName, ".md")

	a.existingTargets[lowerName] = struct{}{}
	a.existingTargets[baseName] = struct{}{}

	if len(bytes.TrimSpace(doc.Content)) == 0 {
		a.emptyNotes = append(a.emptyNotes, fullPath)
	}

	title := strings.TrimSpace(doc.Meta.Title)
	if title != "" {
		a.titlesMap[title] = append(a.titlesMap[title], fullPath)
	}

	a.extractLinks(fullPath, doc.Content)
}

func (a *analyzer) extractLinks(sourcePath string, content []byte) {
	cleanContent := codeBlockRe.ReplaceAll(content, nil)

	cleanContent = inlineCodeRe.ReplaceAll(cleanContent, nil)

	matches := wikilinkRe.FindAllSubmatch(cleanContent, -1)
	for _, match := range matches {
		if len(match) > 1 {
			target := strings.TrimSpace(string(match[1]))

			if target == "" {
				continue
			}

			cleanTarget := filepath.Base(target)

			a.collectedLinks = append(a.collectedLinks, linkRef{
				sourceFile: sourcePath,
				rawTarget:  target,
				normTarget: strings.ToLower(cleanTarget),
			})
		}
	}
}

func (a *analyzer) generateReport() Report {
	report := Report{
		EmptyNotes: a.emptyNotes,
	}

	for title, paths := range a.titlesMap {
		if len(paths) > 1 {
			report.Duplicates = append(report.Duplicates, DuplicateTitle{
				Title: title,
				Paths: paths,
			})
		}
	}

	for _, link := range a.collectedLinks {
		if _, exists := a.existingTargets[link.normTarget]; !exists {
			report.BrokenLinks = append(report.BrokenLinks, BrokenLink{
				SourceFile: link.sourceFile,
				TargetNote: link.rawTarget,
			})
		}
	}

	return report
}

// Run reads directory and checks for problems
func Run(notesPath string) (Report, error) {
	resolvedNotesPath := utils.PathParse(notesPath)
	baseDir := filepath.Dir(resolvedNotesPath)
	filesPath := filepath.Join(baseDir, "files")

	noteEntries, err := os.ReadDir(resolvedNotesPath)
	if err != nil {
		return Report{}, fmt.Errorf("reading notes dir for linting: %w", err)
	}

	anz := newAnalyzer(filesPath)

	for _, entry := range noteEntries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		fullPath := filepath.Join(resolvedNotesPath, entry.Name())
		anz.processNote(fullPath, entry)
	}

	return anz.generateReport(), nil
}
