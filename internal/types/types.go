package types

import "time"

type Repo struct {
	Name            string          `json:"name"`
	NameWithOwner   string          `json:"nameWithOwner"`
	Description     string          `json:"description"`
	Visibility      string          `json:"visibility"`
	PrimaryLanguage *PrimaryLang    `json:"primaryLanguage"`
	StargazerCount  int             `json:"stargazerCount"`
	ForkCount       int             `json:"forkCount"`
	DiskUsage       int             `json:"diskUsage"`
	UpdatedAt       time.Time       `json:"updatedAt"`
	CreatedAt       time.Time       `json:"createdAt"`
	IsArchived      bool            `json:"isArchived"`
	URL             string          `json:"url"`
	SSHUrl          string          `json:"sshUrl"`
}

type PrimaryLang struct {
	Name string `json:"name"`
}

func (r Repo) Lang() string {
	if r.PrimaryLanguage != nil {
		return r.PrimaryLanguage.Name
	}
	return ""
}

type SearchResult struct {
	FullName        string `json:"fullName"`
	Description     string `json:"description"`
	StargazersCount int    `json:"stargazersCount"`
	ForksCount      int    `json:"forksCount"`
	Language        string `json:"language"`
	UpdatedAt       string `json:"updatedAt"`
}

type PR struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
	Author    PRAuthor  `json:"author"`
	CreatedAt time.Time `json:"createdAt"`
}

type PRAuthor struct {
	Login string `json:"login"`
}

type GitRepoStatus struct {
	Path      string
	Name      string
	Branch    string
	Dirty     bool
	Untracked bool
	Ahead     int
	Behind    int
}
