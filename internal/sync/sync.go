package sync

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/dickus/dreadnotes/internal/notes"
)

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
		return "", fmt.Errorf("Couldn't define current branch: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func remoteBranchExists(repoPath, branch string) bool {
	cmd := exec.Command("git", "ls-remote", "--heads", "origin", branch)
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil { return false }

	return strings.TrimSpace(string(output)) != ""
}

func needsPull(repoPath, branch, remote string) bool {
	return revCount(repoPath, branch + ".." + remote) > 0
}

func needsPush(repoPath, branch, remote string) bool {
	return revCount(repoPath, branch + ".." + remote) > 0
}

func revCount(repoPath, revRange string) int {
	cmd := exec.Command("git", "rev-list", "--count", revRange)
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil { return 0 }

	n, _ := strconv.Atoi(strings.TrimSpace(string(output)))

	return n
}

func Sync(repoPath string) error {
	repoPath = notes.PathParse(repoPath)

	if !IsRepo(repoPath) {
		return fmt.Errorf("Directory %s is not a git repo.\nInitialize it with 'git init %s' command.", repoPath, repoPath)
	}

	if err := run(repoPath, "git", "add", "--all"); err != nil { return err }

	if !nothingToCommit(repoPath) {
		if err := run(repoPath, "git", "commit", "-m", "update"); err != nil { return err }
	}

	if ok, _ := HasRemote(repoPath); !ok { return nil }

	branch, err := currentBranch(repoPath)
	if err != nil { return err }

	if !remoteBranchExists(repoPath, branch) {
		return run(repoPath, "git", "push", "-u", "origin", branch)
	}

	if err := run(repoPath, "git", "fetch", "origin", branch); err != nil {
		return fmt.Errorf("Couldn't fetch data from origin: %w", err)
	}

	remote := "origin/" + branch

	if needsPull(repoPath, branch, remote) {
		if err := run(repoPath, "git", "pull", "--rebase", "origin", branch); err != nil {
			return fmt.Errorf("pull conflict: %w", err)
		}
	}

	if needsPush(repoPath, branch, remote) {
		if err := run(repoPath, "git", "push", "origin", branch); err != nil { return err }
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

