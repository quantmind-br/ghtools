package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/diogo/ghtools/internal/types"
)

func run(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s: %s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func CheckInstalled() error {
	_, err := exec.LookPath("git")
	if err != nil {
		return fmt.Errorf("git not found in PATH")
	}
	return nil
}

func FindRepos(basePath string, maxDepth int) ([]string, error) {
	var repos []string
	err := filepath.WalkDir(basePath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // skip errors
		}
		if !d.IsDir() {
			return nil
		}

		// Calculate depth relative to basePath
		rel, _ := filepath.Rel(basePath, path)
		depth := 0
		if rel != "." {
			depth = strings.Count(rel, string(os.PathSeparator)) + 1
		}

		if depth > maxDepth {
			return filepath.SkipDir
		}

		if d.Name() == ".git" {
			repos = append(repos, filepath.Dir(path))
			return filepath.SkipDir
		}

		// Skip common non-repo directories
		if d.Name() == "node_modules" || d.Name() == ".cache" || d.Name() == "vendor" {
			return filepath.SkipDir
		}

		return nil
	})
	return repos, err
}

func CurrentBranch(dir string) string {
	out, err := run(dir, "branch", "--show-current")
	if err != nil {
		return "detached"
	}
	if out == "" {
		return "detached"
	}
	return out
}

func IsDirty(dir string) bool {
	cmd := exec.Command("git", "diff-index", "--quiet", "HEAD", "--")
	cmd.Dir = dir
	return cmd.Run() != nil
}

func HasUntracked(dir string) bool {
	out, _ := run(dir, "ls-files", "--others", "--exclude-standard")
	return out != ""
}

func AheadBehind(dir string) (ahead int, behind int) {
	outA, err := run(dir, "rev-list", "--count", "@{u}..HEAD")
	if err == nil {
		ahead, _ = strconv.Atoi(outA)
	} else {
		ahead = -1
	}

	outB, err := run(dir, "rev-list", "--count", "HEAD..@{u}")
	if err == nil {
		behind, _ = strconv.Atoi(outB)
	} else {
		behind = -1
	}

	return
}

func Pull(dir string) (string, error) {
	return run(dir, "pull", "--ff-only")
}

func Fetch(dir string) error {
	_, err := run(dir, "fetch", "--quiet")
	return err
}

func Push(dir string, branch string) error {
	_, err := run(dir, "push", "-u", "origin", branch)
	return err
}

func HasRemoteBranch(dir string, branch string) bool {
	_, err := run(dir, "ls-remote", "--exit-code", "--heads", "origin", branch)
	return err == nil
}

func IsGitRepo() bool {
	_, err := run(".", "rev-parse", "--is-inside-work-tree")
	return err == nil
}

func GetRepoStatus(dir string) types.GitRepoStatus {
	name := filepath.Base(dir)
	branch := CurrentBranch(dir)
	dirty := IsDirty(dir)
	untracked := HasUntracked(dir)
	ahead, behind := AheadBehind(dir)

	return types.GitRepoStatus{
		Path:      dir,
		Name:      name,
		Branch:    branch,
		Dirty:     dirty,
		Untracked: untracked,
		Ahead:     ahead,
		Behind:    behind,
	}
}
