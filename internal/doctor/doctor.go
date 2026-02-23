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
// It ignores aliases if there are any, so it will only work for the actual links
var wikilinkRe = regexp.MustCompile(`\[\[([^\]|]+)(?:\|[^\]]+)?\]\]`)

// Run reads directory and checks for problems
func Run(notesPath string) (Report, error) {
	report := Report{}
	resolvedNotesPath := utils.PathParse(notesPath)
	baseDir := filepath.Dir(resolvedNotesPath)
	filesPath := filepath.Join(baseDir, "files")
	existingTargets := make(map[string]struct{})
	titlesMap := make(map[string][]string)
	var collectedLinks []linkRef

	fileEntries, err := os.ReadDir(filesPath)
	if err == nil {
		for _, entry := range fileEntries {
			if !entry.IsDir() {
				normFileName := strings.ToLower(strings.TrimSpace(entry.Name()))
				existingTargets[normFileName] = struct{}{}
			}
		}
	}

	noteEntries, err := os.ReadDir(resolvedNotesPath)
	if err != nil {
		return report, fmt.Errorf("reading notes dir for linting: %w", err)
	}

	for _, entry := range noteEntries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		fullPath := filepath.Join(resolvedNotesPath, entry.Name())
		doc, err := frontmatter.ParseFile(fullPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Linter warning: skipping invalid note %s: %v\n", fullPath, err)

			continue
		}

		lowerName := strings.ToLower(strings.TrimSpace(entry.Name()))
		baseName := strings.TrimSuffix(lowerName, ".md")

		existingTargets[lowerName] = struct{}{}
		existingTargets[baseName] = struct{}{}

		if len(bytes.TrimSpace(doc.Content)) == 0 {
			report.EmptyNotes = append(report.EmptyNotes, fullPath)
		}

		title := strings.TrimSpace(doc.Meta.Title)
		if title != "" {
			titlesMap[title] = append(titlesMap[title], fullPath)
		}

		matches := wikilinkRe.FindAllSubmatch(doc.Content, -1)
		for _, match := range matches {
			if len(match) > 1 {
				target := string(match[1])
				target = strings.TrimSpace(target)

				cleanTarget := filepath.Base(target)

				collectedLinks = append(collectedLinks, linkRef{
					sourceFile: fullPath,
					rawTarget:  target,
					normTarget: strings.ToLower(cleanTarget),
				})
			}
		}
	}

	for title, paths := range titlesMap {
		if len(paths) > 1 {
			report.Duplicates = append(report.Duplicates, DuplicateTitle{
				Title: title,
				Paths: paths,
			})
		}
	}

	for _, link := range collectedLinks {
		if _, exists := existingTargets[link.normTarget]; !exists {
			report.BrokenLinks = append(report.BrokenLinks, BrokenLink{
				SourceFile: link.sourceFile,
				TargetNote: link.rawTarget,
			})
		}
	}

	return report, nil
}
