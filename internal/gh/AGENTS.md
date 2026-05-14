# GH KNOWLEDGE BASE

**Generated:** 2025-05-14

## OVERVIEW

GitHub API wrapper - shells out to `gh` CLI rather than using REST API directly.

## WHERE TO LOOK

| Task | File | Notes |
|------|------|-------|
| List repos | `repos.go` | FetchRepos (with caching) |
| Clone repo | `repos.go` | CloneRepo |
| Create repo | `repos.go` | CreateRepo |
| Delete repo | `repos.go` | DeleteRepo + CheckDeleteScope |
| Archive/Unarchive | `repos.go` | ArchiveRepo / UnarchiveRepo |
| Set visibility | `repos.go` | SetVisibility |
| Fork repo | `repos.go` | ForkRepo |
| Browse repo | `repos.go` | BrowseRepo |
| Star repo | `repos.go` | StarRepo |
| View repo | `repos.go` | ViewRepo |
| List PRs | `pr.go` | PRList |
| Create PR | `pr.go` | PRCreate |
| Search repos | `search.go` | GitHub search API wrapper |
| General gh calls | `gh.go` | Low-level gh CLI execution |

## CONVENTIONS

- All functions return `(error)` or `(T, error)`
- Use `repoFields` constant for `gh repo list` field selection
- Cache `FetchRepos` results via `internal/cache/`
- Wrap errors: `fmt.Errorf("description: %w", err)`

## ANTI-PATTERNS

- Never call `gh` CLI directly from cmd/ - always use these wrappers
- Never bypass the cache for repo listing (use FetchRepos)
- Never ignore `gh` CLI exit codes - always check and wrap
