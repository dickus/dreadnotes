package doctor

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// ANSI color codes for terminal output formatting.
const (
	Bold   = "\033[1m"
	Green  = "\033[32m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
	Reset  = "\033[0m"
	Dim    = "\033[2m"
)

// PrintReport formats and prints the linter report to standard output.
// It colorizes the output, groups broken links and duplicate titles by file, and ensures consistent alphabetical sorting of the results.
func PrintReport(r Report) {
	hasIssues := false

	if len(r.EmptyNotes) > 0 {
		hasIssues = true
		fmt.Printf("%s%s▌ Empty Notes:%s\n", Bold, Yellow, Reset)

		sort.Strings(r.EmptyNotes)

		for _, note := range r.EmptyNotes {
			fmt.Printf("%s▌%s - %s\n", Yellow, Reset, filepath.Base(note))
		}

		fmt.Println()
	}

	if len(r.Duplicates) > 0 {
		hasIssues = true
		fmt.Printf("%s%s▌ Duplicate Titles:%s\n", Bold, Yellow, Reset)

		sort.Slice(r.Duplicates, func(i, j int) bool {
			return r.Duplicates[i].Title < r.Duplicates[j].Title
		})

		for _, dup := range r.Duplicates {
			fmt.Printf("%s▌%s %s%s\"%s%s%s\"%s\n", Yellow, Reset, Bold, Dim, Reset, dup.Title, Dim, Reset)

			sort.Strings(dup.Paths)
			totalPaths := len(dup.Paths)

			for i, path := range dup.Paths {
				if i == totalPaths-1 {
					fmt.Printf("%s▌%s %s╰❯%s %s\n", Yellow, Reset, Dim, Reset, filepath.Base(path))
				} else {
					fmt.Printf("%s▌%s %s├❯%s %s\n", Yellow, Reset, Dim, Reset, filepath.Base(path))
				}
			}
		}

		fmt.Println()
	}

	if len(r.BrokenLinks) > 0 {
		hasIssues = true
		fmt.Printf("%s%s▌ Broken Links:%s\n", Bold, Red, Reset)

		groupedLinks := make(map[string][]string)
		var files []string

		for _, link := range r.BrokenLinks {
			baseName := filepath.Base(link.SourceFile)

			if _, exists := groupedLinks[baseName]; !exists {
				files = append(files, baseName)
			}

			groupedLinks[baseName] = append(groupedLinks[baseName], link.TargetNote)
		}

		sort.Strings(files)

		for _, file := range files {
			fmt.Printf("%s▌%s %s%s%s\n", Red, Reset, Bold, file, Reset)

			targets := groupedLinks[file]
			sort.Strings(targets)
			totalTargets := len(targets)

			for i, target := range targets {
				if i == totalTargets-1 {
					fmt.Printf("%s▌%s %s╰❯%s [[%s%s%s]]\n", Red, Reset, Dim, Reset, Red, target, Reset)
				} else {
					fmt.Printf("%s▌%s %s├❯%s [[%s%s%s]]\n", Red, Reset, Dim, Reset, Red, target, Reset)
				}
			}
		}
	}

	if !hasIssues {
		fmt.Printf("%s%s✓ No issues found.%s\n", Bold, Green, Reset)
	} else {
		os.Exit(1)
	}
}
