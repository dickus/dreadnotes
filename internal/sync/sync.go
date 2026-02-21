// Package sync provides simplified git versioning to version and store notes.
package sync

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/dickus/dreadnotes/internal/utils"
)

// IsRepo checks if the specified path is a valid git repository.
func IsRepo(path string) bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = path

	return cmd.Run() == nil
}

func nothingToCommit(repoPath string) bool {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = repoPath

	output, err := cmd.Output()

	return err == nil && strings.TrimSpace(string(output)) == ""
}

// HasRemote determines if the repository has any configured remotes.
func HasRemote(repoPath string) (bool, error) {
	cmd := exec.Command("git", "remote")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(string(output)) != "", nil
}

func currentBranch(repoPath string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("couldn't define current branch: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func remoteBranchExists(repoPath, branch string) bool {
	cmd := exec.Command("git", "ls-remote", "--heads", "origin", branch)
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(output)) != ""
}

func needsPull(repoPath, branch, remote string) bool {
	// Commits in remote but not in local
	return revCount(repoPath, branch+".."+remote) > 0
}

func needsPush(repoPath, branch, remote string) bool {
	// Commits in local but not in remote (Fixed logic)
	return revCount(repoPath, remote+".."+branch) > 0
}

func revCount(repoPath, revRange string) int {
	cmd := exec.Command("git", "rev-list", "--count", revRange)
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return 0
	}

	n, _ := strconv.Atoi(strings.TrimSpace(string(output)))

	return n
}

// Sync adds all changes, commits them, and synchronizes the local repository with the remote (origin) via fetch, rebase, and push.
func Sync(repoPath string) error {
	repoPath = strings.TrimSuffix(utils.PathParse(repoPath), "/notes")

	if !IsRepo(repoPath) {
		return fmt.Errorf("directory %s is not a git repo. Initialize it with 'git init %s'", repoPath, repoPath)
	}

	// Stage all changes
	if err := run(repoPath, "git", "add", "--all"); err != nil {
		return err
	}

	// Commit if there are staged changes
	if !nothingToCommit(repoPath) {
		if err := run(repoPath, "git", "commit", "-m", "update"); err != nil {
			return err
		}
	}

	// Stop here if there's no remote configured
	if ok, _ := HasRemote(repoPath); !ok {
		return nil
	}

	branch, err := currentBranch(repoPath)
	if err != nil {
		return err
	}

	// If the branch doesn't exist on remote, just push -u and we're done
	if !remoteBranchExists(repoPath, branch) {
		return run(repoPath, "git", "push", "-u", "origin", branch)
	}

	// Fetch latest changes to calculate pull/push needs
	if err := run(repoPath, "git", "fetch", "origin", branch); err != nil {
		return fmt.Errorf("couldn't fetch data from origin: %w", err)
	}

	remote := "origin/" + branch

	if needsPull(repoPath, branch, remote) {
		if err := run(repoPath, "git", "pull", "--rebase", "origin", branch); err != nil {
			return fmt.Errorf("pull conflict: %w", err)
		}
	}

	if needsPush(repoPath, branch, remote) {
		if err := run(repoPath, "git", "push", "origin", branch); err != nil {
			return err
		}
	}

	return nil
}

func run(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
