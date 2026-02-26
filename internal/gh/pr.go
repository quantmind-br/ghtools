package gh

import (
	"fmt"

	"github.com/diogo/ghtools/internal/types"
)

func PRList(nameWithOwner string, limit int) ([]types.PR, error) {
	args := []string{"pr", "list", "-R", nameWithOwner, "--limit", fmt.Sprintf("%d", limit),
		"--json", "number,title,state,author,createdAt"}
	var prs []types.PR
	if err := RunJSON(&prs, args...); err != nil {
		return nil, err
	}
	return prs, nil
}

func PRCreate(title string, body string, draft bool) error {
	args := []string{"pr", "create", "--title", title, "--body", body}
	if draft {
		args = append(args, "--draft")
	}
	_, err := Run(args...)
	return err
}
