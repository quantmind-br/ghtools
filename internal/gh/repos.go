package gh

import (
	"fmt"
	"strings"

	"github.com/diogo/ghtools/internal/cache"
	"github.com/diogo/ghtools/internal/types"
)

const repoFields = "name,nameWithOwner,description,visibility,primaryLanguage,stargazerCount,forkCount,diskUsage,updatedAt,createdAt,isArchived,url,sshUrl"

func FetchRepos(forceRefresh bool, cacheTTL int, org string) ([]types.Repo, error) {
	if !forceRefresh && cache.IsValid(cacheTTL) {
		repos, err := cache.Read()
		if err == nil {
			return repos, nil
		}
	}

	args := []string{"repo", "list"}
	if org != "" {
		args = append(args, org)
	}
	args = append(args, "--limit", "1000", "--json", repoFields)

	var repos []types.Repo
	if err := RunJSON(&repos, args...); err != nil {
		return nil, fmt.Errorf("failed to fetch repositories: %w", err)
	}

	_ = cache.Write(repos)
	return repos, nil
}

func CloneRepo(nameWithOwner, targetDir string) error {
	args := []string{"repo", "clone", nameWithOwner}
	if targetDir != "" {
		args = append(args, targetDir)
	}
	_, err := Run(args...)
	return err
}

func DeleteRepo(nameWithOwner string) error {
	_, err := Run("repo", "delete", nameWithOwner, "--yes")
	return err
}

func CreateRepo(name, description, visibility string, clone bool) (string, error) {
	args := []string{"repo", "create", name, "--" + visibility}
	if description != "" {
		args = append(args, "--description", description)
	}
	if clone {
		args = append(args, "--clone")
	}
	return Run(args...)
}

func ArchiveRepo(nameWithOwner string) error {
	_, err := Run("repo", "archive", nameWithOwner, "--yes")
	return err
}

func UnarchiveRepo(nameWithOwner string) error {
	_, err := Run("repo", "unarchive", nameWithOwner, "--yes")
	return err
}

func SetVisibility(nameWithOwner, visibility string) error {
	_, err := Run("repo", "edit", nameWithOwner, "--visibility", visibility)
	return err
}

func ForkRepo(nameWithOwner string, clone bool) error {
	args := []string{"repo", "fork", nameWithOwner}
	if clone {
		args = append(args, "--clone")
	} else {
		args = append(args, "--clone=false")
	}
	_, err := Run(args...)
	return err
}

func BrowseRepo(nameWithOwner string) error {
	_, err := Run("browse", "-R", nameWithOwner)
	return err
}

func StarRepo(nameWithOwner string) error {
	_, err := Run("api", "-X", "PUT", "user/starred/"+nameWithOwner)
	return err
}

func ViewRepo(nameWithOwner string) (string, error) {
	return Run("repo", "view", nameWithOwner)
}

func CheckDeleteScope() bool {
	out, err := Run("auth", "status")
	if err != nil {
		return false
	}
	return strings.Contains(out, "delete_repo")
}

func RefreshDeleteScope() error {
	_, err := Run("auth", "refresh", "-s", "delete_repo")
	return err
}
