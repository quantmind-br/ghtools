package gh

import (
	"fmt"

	"github.com/diogo/ghtools/internal/types"
)

func SearchRepos(query string, sort string, language string, limit int) ([]types.SearchResult, error) {
	args := []string{"search", "repos", query, "--limit", fmt.Sprintf("%d", limit)}
	if sort != "" {
		args = append(args, "--sort", sort)
	}
	if language != "" {
		args = append(args, "--language", language)
	}
	args = append(args, "--json", "fullName,description,stargazersCount,forksCount,language,updatedAt")

	var results []types.SearchResult
	if err := RunJSON(&results, args...); err != nil {
		return nil, err
	}
	return results, nil
}
