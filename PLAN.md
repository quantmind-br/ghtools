# PLAN.md - Refatoracao do ghtools: Bash → Go

**Versao do Plano:** 3.0
**Data:** 2025-12-05
**Objetivo:** Refatorar o ghtools de um script Bash monolitico para uma aplicacao CLI moderna em Go.

---

## 1. Resumo Executivo

Este plano detalha a refatoracao completa do `ghtools` (atualmente um script Bash de ~2144 linhas) para uma aplicacao CLI em Go. A refatoracao resolve problemas estruturais de seguranca, testabilidade e portabilidade, mantendo a experiencia de usuario atual.

### Beneficios Esperados

| Aspecto | Bash (Atual) | Go (Proposto) |
|---------|--------------|---------------|
| **Seguranca** | Vulneravel a injecao de comandos | Elimina classe de vulnerabilidades RCE |
| **Testabilidade** | ~38% cobertura, dificil mockar | >80% cobertura, arquitetura testavel |
| **Portabilidade** | 5 dependencias externas | Binario unico estatico |
| **Performance** | Subprocessos, I/O lento | Compilado, goroutines nativas |
| **Distribuicao** | Script + install.sh | Single binary (cross-compile) |

---

## 2. Analise do Estado Atual

### 2.1 Estrutura do Script Bash

```
ghtools (2144 linhas)
├── Configuracao (linhas 1-92)
│   ├── Variaveis globais (VERSION, CACHE_*, MAX_JOBS)
│   ├── load_config() - carrega ~/.config/ghtools/config
│   └── init_config() - cria config padrao
│
├── UI/UX Layer (linhas 93-429)
│   ├── Cores e estilos (COLOR_*, NC, etc.)
│   ├── gum_*() - wrappers para gum CLI
│   ├── print_*() - output formatado (error, success, info, warning)
│   ├── show_header() / show_divider() - formatacao visual
│   └── show_usage() - help text
│
├── Core Infrastructure (linhas 430-543)
│   ├── check_dependencies() - verifica gh, fzf, git, jq
│   ├── check_gh_auth() - verifica autenticacao GitHub
│   ├── is_cache_valid() - TTL do cache
│   ├── fetch_repositories_json() - API GitHub com cache
│   ├── wait_for_jobs() - controle de paralelismo
│   └── truncate_text() - utilidade de string
│
├── Actions (linhas 544-1936)
│   ├── action_list() - lista repos com filtros
│   ├── action_clone() - clone paralelo
│   ├── action_sync() - sync paralelo com ff-only
│   ├── action_delete() - delecao com confirmacao
│   ├── action_create() - criacao com templates
│   ├── action_fork() - fork de repos externos
│   ├── action_explore() - busca repos externos
│   ├── action_trending() - repos trending
│   ├── action_archive() - archive/unarchive
│   ├── action_stats() - dashboard estatisticas
│   ├── action_search() - busca fuzzy local
│   ├── action_browse() - abre no browser
│   ├── action_visibility() - altera visibilidade
│   ├── action_pr() - gerencia PRs
│   └── action_status() - status de repos locais
│
├── Menu Interativo (linhas 1937-2036)
│   └── show_menu() - menu principal com gum/fzf
│
└── Entry Point (linhas 2037-2143)
    └── main() - parsing de args, routing
```

### 2.2 Dependencias Externas Atuais

| Dependencia | Uso | Substituicao em Go |
|-------------|-----|-------------------|
| `gh` (GitHub CLI) | API GitHub, auth | `go-github` + OAuth |
| `fzf` | Selecao interativa | `bubbletea` + `bubbles` |
| `gum` | UI bonita | `lipgloss` + `huh` |
| `jq` | Parsing JSON | Nativo (`encoding/json`) |
| `git` | Operacoes Git locais | `go-git` ou `os/exec` |

### 2.3 Funcionalidades por Categoria

**Repository Management:**
- `list` - Listar repos com filtros (lang, org)
- `clone` - Clone paralelo de multiplos repos
- `create` - Criar repo com templates (python, node, go)
- `delete` - Deletar repos com confirmacao
- `fork` - Fork de repos externos
- `archive` - Archive/unarchive repos
- `visibility` - Alterar public/private

**Local Operations:**
- `sync` - Sync paralelo (git pull --ff-only)
- `status` - Status de repos locais (dirty, ahead/behind)

**Discovery:**
- `search` - Busca fuzzy nos proprios repos
- `explore` - Busca repos externos (GitHub Search)
- `trending` - Repos trending por linguagem
- `browse` - Abrir repos no browser

**Pull Requests:**
- `pr list` - Listar PRs de um repo
- `pr create` - Criar PR da branch atual

**Utilities:**
- `refresh` - Limpar cache
- `config` - Inicializar/mostrar config
- `stats` - Dashboard de estatisticas

---

## 3. Arquitetura Proposta em Go

### 3.1 Estrutura de Diretorios

```
ghtools/
├── cmd/
│   └── ghtools/
│       └── main.go              # Entry point
│
├── internal/
│   ├── cli/
│   │   ├── root.go              # Root command (cobra)
│   │   ├── list.go              # list command
│   │   ├── clone.go             # clone command
│   │   ├── sync.go              # sync command
│   │   ├── create.go            # create command
│   │   ├── delete.go            # delete command
│   │   ├── fork.go              # fork command
│   │   ├── explore.go           # explore command
│   │   ├── trending.go          # trending command
│   │   ├── archive.go           # archive command
│   │   ├── stats.go             # stats command
│   │   ├── search.go            # search command
│   │   ├── browse.go            # browse command
│   │   ├── visibility.go        # visibility command
│   │   ├── pr.go                # pr command group
│   │   └── status.go            # status command
│   │
│   ├── github/
│   │   ├── client.go            # GitHub API client wrapper
│   │   ├── repository.go        # Repository operations
│   │   ├── pullrequest.go       # PR operations
│   │   ├── search.go            # Search operations
│   │   ├── auth.go              # Authentication handling
│   │   └── cache.go             # Repository cache
│   │
│   ├── git/
│   │   ├── local.go             # Local git operations
│   │   ├── sync.go              # Sync/pull logic
│   │   ├── status.go            # Status checking
│   │   └── discovery.go         # Find git repos in directory
│   │
│   ├── tui/
│   │   ├── styles.go            # lipgloss styles
│   │   ├── selector.go          # Interactive selector (bubbletea)
│   │   ├── confirm.go           # Confirmation dialogs
│   │   ├── input.go             # Text input
│   │   ├── spinner.go           # Loading spinner
│   │   ├── table.go             # Table rendering
│   │   └── menu.go              # Main menu
│   │
│   ├── config/
│   │   ├── config.go            # Configuration struct
│   │   ├── loader.go            # Load/save config
│   │   └── defaults.go          # Default values
│   │
│   ├── templates/
│   │   ├── templates.go         # Template interface
│   │   ├── python.go            # Python template
│   │   ├── node.go              # Node template
│   │   └── golang.go            # Go template
│   │
│   └── util/
│       ├── parallel.go          # Parallel execution helpers
│       ├── strings.go           # String utilities
│       └── browser.go           # Open URL in browser
│
├── pkg/
│   └── version/
│       └── version.go           # Version info
│
├── test/
│   ├── integration/             # Integration tests
│   │   ├── github_test.go
│   │   ├── git_test.go
│   │   └── cli_test.go
│   └── fixtures/                # Test fixtures
│       └── repos.json
│
├── go.mod
├── go.sum
├── Makefile
├── .goreleaser.yml              # Release automation
└── README.md
```

### 3.2 Diagrama de Dependencias

```
┌─────────────────────────────────────────────────────────────────┐
│                         cmd/ghtools                              │
│                           main.go                                │
└─────────────────────────────┬───────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                        internal/cli                              │
│   root.go ─── list.go, clone.go, sync.go, create.go, ...        │
└────────┬────────────────────┬────────────────────┬──────────────┘
         │                    │                    │
         ▼                    ▼                    ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│ internal/github │  │  internal/git   │  │  internal/tui   │
│                 │  │                 │  │                 │
│ client.go       │  │ local.go        │  │ styles.go       │
│ repository.go   │  │ sync.go         │  │ selector.go     │
│ cache.go        │  │ status.go       │  │ confirm.go      │
│ search.go       │  │ discovery.go    │  │ table.go        │
└────────┬────────┘  └────────┬────────┘  └────────┬────────┘
         │                    │                    │
         ▼                    ▼                    ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   go-github     │  │    go-git /     │  │   bubbletea     │
│   (biblioteca)  │  │    os/exec      │  │   lipgloss      │
└─────────────────┘  └─────────────────┘  └─────────────────┘
```

### 3.3 Bibliotecas Go Recomendadas

| Categoria | Biblioteca | Justificativa |
|-----------|-----------|---------------|
| **CLI Framework** | `github.com/spf13/cobra` | Padrao industria, subcomandos, flags |
| **GitHub API** | `github.com/google/go-github/v57` | Client oficial, tipado |
| **Git Local** | `github.com/go-git/go-git/v5` | Pure Go, sem dependencia de git CLI |
| **TUI Framework** | `github.com/charmbracelet/bubbletea` | Elm architecture, testavel |
| **TUI Styles** | `github.com/charmbracelet/lipgloss` | Estilos consistentes |
| **TUI Forms** | `github.com/charmbracelet/huh` | Forms interativos |
| **TUI Spinner** | `github.com/charmbracelet/bubbles` | Componentes reutilizaveis |
| **Config** | `github.com/spf13/viper` | Multiplos formatos, env vars |
| **Logging** | `log/slog` (stdlib) | Structured logging nativo |
| **Testing** | `github.com/stretchr/testify` | Assertions, mocks |

---

## 4. Mapeamento de Funcoes Bash → Go

### 4.1 Configuracao

| Bash | Go | Pacote |
|------|-----|--------|
| `load_config()` | `config.Load()` | `internal/config` |
| `init_config()` | `config.Init()` | `internal/config` |
| Variaveis globais | `config.Config` struct | `internal/config` |

### 4.2 UI/UX

| Bash | Go | Pacote |
|------|-----|--------|
| `print_error()` | `tui.Error()` | `internal/tui` |
| `print_success()` | `tui.Success()` | `internal/tui` |
| `print_info()` | `tui.Info()` | `internal/tui` |
| `print_warning()` | `tui.Warning()` | `internal/tui` |
| `gum_confirm()` | `tui.Confirm()` | `internal/tui` |
| `gum_input()` | `tui.Input()` | `internal/tui` |
| `gum_choose()` | `tui.Select()` | `internal/tui` |
| `gum_filter()` | `tui.Filter()` | `internal/tui` |
| `show_header()` | `tui.Header()` | `internal/tui` |
| `run_with_spinner()` | `tui.Spinner()` | `internal/tui` |

### 4.3 Core

| Bash | Go | Pacote |
|------|-----|--------|
| `check_dependencies()` | Removido (binario unico) | - |
| `check_gh_auth()` | `github.CheckAuth()` | `internal/github` |
| `fetch_repositories_json()` | `github.ListRepos()` + cache | `internal/github` |
| `is_cache_valid()` | `cache.IsValid()` | `internal/github` |
| `wait_for_jobs()` | `util.RunParallel()` | `internal/util` |
| `truncate_text()` | `util.Truncate()` | `internal/util` |

### 4.4 Actions

| Bash | Go Command | Handler |
|------|------------|---------|
| `action_list()` | `ghtools list` | `cli.listCmd` |
| `action_clone()` | `ghtools clone` | `cli.cloneCmd` |
| `action_sync()` | `ghtools sync` | `cli.syncCmd` |
| `action_create()` | `ghtools create` | `cli.createCmd` |
| `action_delete()` | `ghtools delete` | `cli.deleteCmd` |
| `action_fork()` | `ghtools fork` | `cli.forkCmd` |
| `action_explore()` | `ghtools explore` | `cli.exploreCmd` |
| `action_trending()` | `ghtools trending` | `cli.trendingCmd` |
| `action_archive()` | `ghtools archive` | `cli.archiveCmd` |
| `action_stats()` | `ghtools stats` | `cli.statsCmd` |
| `action_search()` | `ghtools search` | `cli.searchCmd` |
| `action_browse()` | `ghtools browse` | `cli.browseCmd` |
| `action_visibility()` | `ghtools visibility` | `cli.visibilityCmd` |
| `action_pr()` | `ghtools pr` | `cli.prCmd` |
| `action_status()` | `ghtools status` | `cli.statusCmd` |

---

## 5. Fases de Implementacao

### Fase 0: Setup do Projeto (Pre-requisito)

**Duracao estimada:** 1 sessao de trabalho

**Tarefas:**

0.1. Criar estrutura de diretorios Go
```bash
mkdir -p cmd/ghtools internal/{cli,github,git,tui,config,templates,util} pkg/version test/{integration,fixtures}
```

0.2. Inicializar modulo Go
```bash
go mod init github.com/diogosalesdev/ghtools
```

0.3. Adicionar dependencias principais
```bash
go get github.com/spf13/cobra@latest
go get github.com/google/go-github/v57@latest
go get github.com/go-git/go-git/v5@latest
go get github.com/charmbracelet/bubbletea@latest
go get github.com/charmbracelet/lipgloss@latest
go get github.com/charmbracelet/huh@latest
go get github.com/charmbracelet/bubbles@latest
go get github.com/spf13/viper@latest
go get github.com/stretchr/testify@latest
```

0.4. Criar arquivos base
- `cmd/ghtools/main.go` - entry point minimo
- `internal/cli/root.go` - comando root cobra
- `pkg/version/version.go` - informacao de versao

0.5. Criar Makefile com targets basicos
```makefile
.PHONY: build test lint run

build:
	go build -o bin/ghtools ./cmd/ghtools

test:
	go test -v -race -cover ./...

lint:
	golangci-lint run

run: build
	./bin/ghtools
```

0.6. Verificar build inicial funciona
```bash
make build && ./bin/ghtools --version
```

---

### Fase 1: Foundation - Config e TUI Base

**Duracao estimada:** 2-3 sessoes de trabalho

**Objetivo:** Estabelecer infraestrutura de configuracao e UI base.

#### 1.1 Modulo de Configuracao (`internal/config`)

**1.1.1. Criar struct de configuracao (`config.go`)**
```go
type Config struct {
    CacheTTL         time.Duration `mapstructure:"cache_ttl"`
    CacheFile        string        `mapstructure:"cache_file"`
    MaxJobs          int           `mapstructure:"max_jobs"`
    DefaultOrg       string        `mapstructure:"default_org"`
    DefaultClonePath string        `mapstructure:"default_clone_path"`
}
```

**1.1.2. Implementar loader com viper (`loader.go`)**
- Suportar arquivo TOML em `~/.config/ghtools/config.toml`
- Suportar variaveis de ambiente com prefixo `GHTOOLS_`
- Validacao de valores

**1.1.3. Implementar defaults (`defaults.go`)**
- CacheTTL: 10 minutos
- MaxJobs: 5
- CacheFile: `/tmp/ghtools_repos_<uid>.json`

**1.1.4. Testes unitarios**
- Teste de carregamento de config valida
- Teste de defaults quando arquivo nao existe
- Teste de validacao de valores invalidos
- Cobertura minima: 90%

#### 1.2 Modulo TUI Base (`internal/tui`)

**1.2.1. Definir tema de cores (`styles.go`)**
```go
var (
    ColorPrimary   = lipgloss.Color("99")   // Purple
    ColorSecondary = lipgloss.Color("39")   // Cyan
    ColorAccent    = lipgloss.Color("212")  // Pink
    ColorSuccess   = lipgloss.Color("78")   // Green
    ColorWarning   = lipgloss.Color("220")  // Yellow
    ColorError     = lipgloss.Color("196")  // Red
    ColorInfo      = lipgloss.Color("75")   // Light blue
    ColorMuted     = lipgloss.Color("240")  // Gray
)
```

**1.2.2. Implementar funcoes de output (`output.go`)**
- `Error(msg string)` - output de erro formatado
- `Success(msg string)` - output de sucesso
- `Info(msg string)` - output informativo
- `Warning(msg string)` - output de aviso
- `Header(title, subtitle string)` - cabecalho estilizado
- `Divider(title string)` - divisor de secao

**1.2.3. Implementar spinner (`spinner.go`)**
- Wrapper sobre bubbles/spinner
- Suporte a titulo e callback de conclusao

**1.2.4. Implementar confirmacao (`confirm.go`)**
- Dialog de confirmacao Yes/No
- Suporte a default value
- Fallback para input simples se terminal nao suportar

**1.2.5. Implementar input de texto (`input.go`)**
- Input com placeholder
- Suporte a valor default
- Validacao opcional

**1.2.6. Testes unitarios**
- Testes de renderizacao (verificar output esperado)
- Testes de styles aplicados
- Cobertura minima: 80%

---

### Fase 2: GitHub Client e Cache

**Duracao estimada:** 2-3 sessoes de trabalho

**Objetivo:** Implementar integracao com GitHub API com cache.

#### 2.1 Cliente GitHub (`internal/github`)

**2.1.1. Implementar autenticacao (`auth.go`)**
- Usar token do `gh auth token` (compatibilidade)
- Suporte a GITHUB_TOKEN env var
- Verificacao de autenticacao valida
- Listar scopes disponiveis

**2.1.2. Implementar client wrapper (`client.go`)**
```go
type Client struct {
    gh     *github.Client
    cache  *Cache
    config *config.Config
}

func NewClient(cfg *config.Config) (*Client, error)
func (c *Client) ListRepositories(ctx context.Context, opts ListOptions) ([]Repository, error)
func (c *Client) GetRepository(ctx context.Context, owner, repo string) (*Repository, error)
```

**2.1.3. Definir modelo Repository (`repository.go`)**
```go
type Repository struct {
    Name            string
    NameWithOwner   string
    Description     string
    Visibility      string
    PrimaryLanguage string
    Stars           int
    Forks           int
    DiskUsage       int
    UpdatedAt       time.Time
    CreatedAt       time.Time
    IsArchived      bool
    URL             string
    SSHURL          string
}
```

**2.1.4. Implementar operacoes de repositorio (`repository.go`)**
- `CreateRepository(ctx, name, opts)` - criar repo
- `DeleteRepository(ctx, owner, repo)` - deletar repo
- `ForkRepository(ctx, owner, repo, opts)` - fork
- `ArchiveRepository(ctx, owner, repo)` - archive
- `UnarchiveRepository(ctx, owner, repo)` - unarchive
- `SetVisibility(ctx, owner, repo, visibility)` - alterar visibilidade

**2.1.5. Implementar cache (`cache.go`)**
```go
type Cache struct {
    file     string
    ttl      time.Duration
    mu       sync.RWMutex
}

func (c *Cache) IsValid() bool
func (c *Cache) Load() ([]Repository, error)
func (c *Cache) Save(repos []Repository) error
func (c *Cache) Invalidate() error
```

**2.1.6. Implementar search (`search.go`)**
- `SearchRepositories(ctx, query, opts)` - busca global
- `SearchTrending(ctx, lang, since)` - repos trending

**2.1.7. Testes unitarios**
- Mock do cliente GitHub
- Testes de cache (save, load, invalidate, TTL)
- Testes de conversao de modelos
- Cobertura minima: 85%

---

### Fase 3: Git Local Operations

**Duracao estimada:** 1-2 sessoes de trabalho

**Objetivo:** Implementar operacoes Git locais.

#### 3.1 Modulo Git (`internal/git`)

**3.1.1. Implementar discovery (`discovery.go`)**
```go
func FindRepositories(basePath string, maxDepth int) ([]string, error)
```
- Buscar diretorios .git recursivamente
- Suporte a caminhos com espacos
- Limite de profundidade configuravel

**3.1.2. Implementar status (`status.go`)**
```go
type RepoStatus struct {
    Path       string
    Branch     string
    IsDirty    bool
    HasUntracked bool
    Ahead      int
    Behind     int
}

func GetStatus(repoPath string) (*RepoStatus, error)
```

**3.1.3. Implementar sync (`sync.go`)**
```go
type SyncResult struct {
    Path    string
    Success bool
    Message string
    Error   error
}

func Sync(repoPath string, dryRun bool) (*SyncResult, error)
func SyncParallel(repos []string, maxJobs int, dryRun bool) []SyncResult
```
- Pull com --ff-only
- Detectar dirty state antes de pull
- Suporte a dry-run

**3.1.4. Implementar clone (`local.go`)**
```go
func Clone(url, destPath string) error
func CloneParallel(repos []string, destPath string, maxJobs int) []CloneResult
```

**3.1.5. Testes unitarios**
- Testes com repositorios Git de teste
- Testes de paralelismo
- Cobertura minima: 80%

---

### Fase 4: TUI Avancado - Selectors e Menus

**Duracao estimada:** 2 sessoes de trabalho

**Objetivo:** Implementar componentes TUI interativos.

#### 4.1 Selector Interativo (`internal/tui`)

**4.1.1. Implementar selector simples (`selector.go`)**
```go
type SelectorOption struct {
    Label       string
    Value       string
    Description string
}

func Select(header string, options []SelectorOption) (string, error)
func SelectMultiple(header string, options []SelectorOption) ([]string, error)
```
- Navegacao com setas
- Busca fuzzy
- Multi-select com TAB

**4.1.2. Implementar selector de repositorios (`repo_selector.go`)**
- Formatacao especial para repos (nome, descricao, lang, stars)
- Preview de detalhes
- Multi-select para operacoes em lote

**4.1.3. Implementar menu principal (`menu.go`)**
- Menu com todas as opcoes
- Categorias visuais
- Loop interativo

**4.1.4. Implementar table rendering (`table.go`)**
```go
type Table struct {
    Headers []string
    Rows    [][]string
    Widths  []int
}

func (t *Table) Render() string
```

**4.1.5. Testes**
- Testes de renderizacao
- Testes de navegacao (simulado)
- Cobertura minima: 75%

---

### Fase 5: CLI Commands - Parte 1 (Core)

**Duracao estimada:** 3-4 sessoes de trabalho

**Objetivo:** Implementar comandos CLI principais.

#### 5.1 Comando Root (`internal/cli/root.go`)

**5.1.1. Setup cobra root command**
```go
var rootCmd = &cobra.Command{
    Use:     "ghtools",
    Short:   "GitHub repository management tool",
    Long:    `ghtools is a unified CLI for managing GitHub repositories...`,
    Version: version.Version,
}
```

**5.1.2. Flags globais**
- `--verbose, -V` - output detalhado
- `--quiet, -q` - output minimo
- `--config` - path para config alternativo

**5.1.3. Persistent pre-run**
- Carregar configuracao
- Verificar autenticacao GitHub
- Inicializar cliente

#### 5.2 Comando `list` (`internal/cli/list.go`)

**5.2.1. Implementacao**
```go
var listCmd = &cobra.Command{
    Use:   "list",
    Short: "List repositories",
    RunE:  runList,
}

func init() {
    listCmd.Flags().Bool("refresh", false, "Force refresh cache")
    listCmd.Flags().String("lang", "", "Filter by language")
    listCmd.Flags().String("org", "", "Filter by organization")
}
```

**5.2.2. Funcionalidade**
- Listar repos com formatacao tabular
- Aplicar filtros (lang, org)
- Usar cache ou forcar refresh
- Mostrar: nome, descricao, visibilidade, linguagem, data

#### 5.3 Comando `clone` (`internal/cli/clone.go`)

**5.3.1. Implementacao**
- Selector interativo de repos
- Multi-select
- Clone paralelo
- Barra de progresso

**5.3.2. Flags**
- `--path` - diretorio destino

#### 5.4 Comando `sync` (`internal/cli/sync.go`)

**5.4.1. Implementacao**
- Descobrir repos locais
- Selector ou --all
- Sync paralelo
- Dry-run mode

**5.4.2. Flags**
- `--path` - diretorio base
- `--dry-run` - simular apenas
- `--all` - sync todos sem selecao
- `--max-depth` - profundidade de busca

#### 5.5 Comando `status` (`internal/cli/status.go`)

**5.5.1. Implementacao**
- Descobrir repos locais
- Mostrar status de cada um
- Tabela formatada

**5.5.2. Flags**
- `--path` - diretorio base
- `--max-depth` - profundidade

#### 5.6 Comando `stats` (`internal/cli/stats.go`)

**5.6.1. Implementacao**
- Dashboard com estatisticas
- Total repos, public/private/archived
- Stars, forks, size total
- Top linguagens
- Top repos por stars
- Repos recentes

---

### Fase 6: CLI Commands - Parte 2 (CRUD)

**Duracao estimada:** 2-3 sessoes de trabalho

**Objetivo:** Implementar comandos de criacao/delecao.

#### 6.1 Comando `create` (`internal/cli/create.go`)

**6.1.1. Implementacao interativa**
- Input de nome
- Input de descricao
- Selecao de visibilidade
- Selecao de template
- Confirmacao

**6.1.2. Templates**
- Implementar templates em `internal/templates`
- Python: README, .gitignore, main.py
- Node: README, package.json, index.js, .gitignore
- Go: README, go.mod, main.go

**6.1.3. Flags**
- `--name` - nome do repo (skip prompt)
- `--description` - descricao
- `--public/--private` - visibilidade
- `--template` - template a aplicar

#### 6.2 Comando `delete` (`internal/cli/delete.go`)

**6.2.1. Implementacao**
- Verificar scope delete_repo
- Selector multi-select
- Confirmacao explicita (digitar DELETE)
- Dry-run por default

**6.2.2. Seguranca**
- Listar repos a deletar
- Requer confirmacao dupla
- Invalidar cache apos

#### 6.3 Comando `archive` (`internal/cli/archive.go`)

**6.3.1. Implementacao**
- Selector de repos nao-arquivados (ou --unarchive para arquivados)
- Confirmacao
- Operacao em lote

**6.3.2. Flags**
- `--unarchive` - modo unarchive

#### 6.4 Comando `visibility` (`internal/cli/visibility.go`)

**6.4.1. Implementacao**
- Selector de repos
- Confirmacao com preview de mudanca
- Operacao em lote

**6.4.2. Flags**
- `--public` - tornar publico
- `--private` - tornar privado

---

### Fase 7: CLI Commands - Parte 3 (Discovery)

**Duracao estimada:** 2 sessoes de trabalho

**Objetivo:** Implementar comandos de descoberta.

#### 7.1 Comando `search` (`internal/cli/search.go`)

**7.1.1. Implementacao**
- Busca fuzzy nos proprios repos
- Selector com resultados
- Acoes pos-selecao (clone, browse, delete)

#### 7.2 Comando `explore` (`internal/cli/explore.go`)

**7.2.1. Implementacao interativa**
- Input de query
- Selecao de ordenacao (stars, forks, updated)
- Filtro de linguagem opcional
- Resultados com preview
- Acoes (clone, fork, browse, star)

**7.2.2. Flags**
- `--sort` - ordenacao
- `--lang` - linguagem
- `--limit` - limite de resultados

#### 7.3 Comando `trending` (`internal/cli/trending.go`)

**7.3.1. Implementacao**
- Buscar repos trending
- Filtro por linguagem
- Resultados com preview
- Acoes (clone, fork, browse, star)

**7.3.2. Flags**
- `--lang` - linguagem

#### 7.4 Comando `fork` (`internal/cli/fork.go`)

**7.4.1. Implementacao**
- Input de query de busca
- Selector de resultados
- Fork com opcao de clone

**7.4.2. Flags**
- `--clone` - clonar apos fork

#### 7.5 Comando `browse` (`internal/cli/browse.go`)

**7.5.1. Implementacao**
- Selector de repos
- Abrir no browser default

---

### Fase 8: Pull Requests e Menu

**Duracao estimada:** 1-2 sessoes de trabalho

**Objetivo:** Implementar comandos de PR e menu interativo.

#### 8.1 Comando `pr` (`internal/cli/pr.go`)

**8.1.1. Subcomando `pr list`**
- Selector de repo
- Listar PRs (numero, titulo, estado, autor, data)

**8.1.2. Subcomando `pr create`**
- Detectar branch atual
- Validar nao esta em main/master
- Validar nao esta em detached HEAD
- Push branch se necessario
- Input de titulo
- Opcao de draft

#### 8.2 Comando `refresh` (`internal/cli/refresh.go`)

**8.2.1. Implementacao**
- Invalidar cache
- Mensagem de confirmacao

#### 8.3 Comando `config` (`internal/cli/config.go`)

**8.3.1. Implementacao**
- Criar config default se nao existir
- Mostrar path do arquivo

#### 8.4 Menu Interativo

**8.4.1. Implementar no root command**
- Quando executado sem subcomando, mostrar menu
- Loop ate Exit
- Categorias visuais

---

### Fase 9: Testes e Qualidade

**Duracao estimada:** 2-3 sessoes de trabalho

**Objetivo:** Garantir qualidade e cobertura de testes.

#### 9.1 Testes Unitarios

**9.1.1. Meta de cobertura**
- `internal/config`: 90%+
- `internal/github`: 85%+
- `internal/git`: 80%+
- `internal/tui`: 75%+
- `internal/cli`: 70%+

**9.1.2. Ferramentas**
- `testify` para assertions
- Table-driven tests
- Mocks para GitHub client

#### 9.2 Testes de Integracao

**9.2.1. Implementar em `test/integration`**
- Testes E2E com GitHub API real (skip em CI)
- Testes de CLI com input simulado
- Testes de cache

#### 9.3 Linting e Formatacao

**9.3.1. Configurar golangci-lint**
```yaml
# .golangci.yml
linters:
  enable:
    - gofmt
    - govet
    - errcheck
    - staticcheck
    - gosec
    - revive
```

**9.3.2. Pre-commit hooks**
- `go fmt`
- `golangci-lint run`
- `go test`

---

### Fase 10: Build, Release e Documentacao

**Duracao estimada:** 1-2 sessoes de trabalho

**Objetivo:** Preparar para distribuicao.

#### 10.1 Build System

**10.1.1. Makefile completo**
```makefile
VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -X github.com/diogosalesdev/ghtools/pkg/version.Version=$(VERSION)

build:
	go build -ldflags "$(LDFLAGS)" -o bin/ghtools ./cmd/ghtools

build-all:
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/ghtools-linux-amd64 ./cmd/ghtools
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/ghtools-darwin-amd64 ./cmd/ghtools
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o bin/ghtools-darwin-arm64 ./cmd/ghtools
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o bin/ghtools-windows-amd64.exe ./cmd/ghtools

install:
	go install -ldflags "$(LDFLAGS)" ./cmd/ghtools
```

#### 10.2 GoReleaser

**10.2.1. Configurar `.goreleaser.yml`**
- Builds para Linux, macOS (amd64, arm64), Windows
- Checksums
- Release notes automaticas
- Homebrew tap (futuro)

#### 10.3 Documentacao

**10.3.1. README.md**
- Descricao
- Instalacao (go install, binarios)
- Uso basico
- Comandos disponiveis
- Configuracao
- Migracao do Bash

**10.3.2. CHANGELOG.md**
- Versao inicial 4.0.0 (Go rewrite)
- Breaking changes do Bash

**10.3.3. MIGRATION.md**
- Guia de migracao do script Bash
- Diferencas de comportamento
- Novos features

---

## 6. Estrategia de Testes

### 6.1 Piramide de Testes

```
           /\
          /  \   E2E Tests (5%)
         /----\  - Full CLI execution
        /      \ - Real GitHub API (optional)
       /--------\
      /  Integra-\ Integration Tests (25%)
     /   tion     \ - Multiple packages
    /--------------\ - Mocked external deps
   /                \
  /   Unit Tests     \ Unit Tests (70%)
 /    (isolated)      \ - Single package
/______________________\ - Fast execution
```

### 6.2 Casos de Teste Criticos

| Categoria | Caso | Prioridade |
|-----------|------|------------|
| Auth | Token invalido | Alta |
| Auth | Sem token | Alta |
| Cache | TTL expirado | Alta |
| Cache | Arquivo corrompido | Media |
| Clone | Path com espacos | Alta |
| Sync | Repo dirty | Alta |
| Sync | Detached HEAD | Media |
| Delete | Sem scope delete_repo | Alta |
| TUI | Terminal nao-interativo | Media |
| Config | Arquivo invalido | Media |

---

## 7. Riscos e Mitigacoes

| Risco | Probabilidade | Impacto | Mitigacao |
|-------|--------------|---------|-----------|
| go-git nao suporta todas operacoes | Media | Alto | Fallback para git CLI via os/exec |
| Bubbletea dificil de testar | Media | Medio | Separar logica de apresentacao |
| Auth diferente do gh CLI | Baixa | Alto | Reusar token do gh auth |
| Breaking changes para usuarios | Alta | Alto | Documentar MIGRATION.md |
| Performance pior que Bash em alguns casos | Baixa | Baixo | Profile e otimizar |

---

## 8. Criterios de Aceitacao

### 8.1 Funcionais

- [ ] Todos os comandos do Bash implementados em Go
- [ ] Paridade de funcionalidade com versao Bash
- [ ] Cache funcional com TTL configuravel
- [ ] Operacoes paralelas para clone/sync
- [ ] Menu interativo equivalente

### 8.2 Nao-Funcionais

- [ ] Binario unico sem dependencias externas
- [ ] Build para Linux, macOS, Windows
- [ ] Cobertura de testes > 80% geral
- [ ] Tempo de startup < 100ms
- [ ] Memoria < 50MB para operacoes normais

### 8.3 Seguranca

- [ ] Sem vulnerabilidades de command injection
- [ ] Token armazenado de forma segura
- [ ] Cache com permissoes 600
- [ ] Validacao de input do usuario

### 8.4 UX

- [ ] Output colorido equivalente ao Bash
- [ ] Mensagens de erro claras
- [ ] Help text para todos comandos
- [ ] Fallback gracioso quando terminal nao suporta TUI

---

## 9. Cronograma Resumido

| Fase | Descricao | Sessoes Estimadas |
|------|-----------|-------------------|
| 0 | Setup do Projeto | 1 |
| 1 | Config e TUI Base | 2-3 |
| 2 | GitHub Client e Cache | 2-3 |
| 3 | Git Local Operations | 1-2 |
| 4 | TUI Avancado | 2 |
| 5 | CLI Commands Core | 3-4 |
| 6 | CLI Commands CRUD | 2-3 |
| 7 | CLI Commands Discovery | 2 |
| 8 | PRs e Menu | 1-2 |
| 9 | Testes e Qualidade | 2-3 |
| 10 | Build e Release | 1-2 |
| **Total** | | **19-27 sessoes** |

---

## 10. Proximos Passos Imediatos

1. **Aprovar este plano** - Revisar e confirmar escopo
2. **Criar TASKS.md** - Checklist detalhado para Fase 0
3. **Iniciar Fase 0** - Setup do projeto Go
4. **Manter script Bash** - Em paralelo ate Go estar completo
5. **Testes manuais** - Validar paridade durante desenvolvimento

---

## Apendice A: Exemplo de Codigo

### A.1 Entry Point (`cmd/ghtools/main.go`)

```go
package main

import (
	"os"

	"github.com/diogosalesdev/ghtools/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
```

### A.2 Root Command (`internal/cli/root.go`)

```go
package cli

import (
	"github.com/spf13/cobra"
	"github.com/diogosalesdev/ghtools/internal/config"
	"github.com/diogosalesdev/ghtools/internal/github"
	"github.com/diogosalesdev/ghtools/pkg/version"
)

var (
	cfg      *config.Config
	ghClient *github.Client
	verbose  bool
	quiet    bool
)

var rootCmd = &cobra.Command{
	Use:     "ghtools",
	Short:   "GitHub repository management tool",
	Long:    `ghtools is a unified CLI for managing GitHub repositories with a beautiful TUI.`,
	Version: version.Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.Load()
		if err != nil {
			return err
		}
		ghClient, err = github.NewClient(cfg)
		return err
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Show interactive menu when no subcommand
		return runMenu()
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "V", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Suppress non-error output")

	// Add subcommands
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(cloneCmd)
	rootCmd.AddCommand(syncCmd)
	// ... other commands
}

func Execute() error {
	return rootCmd.Execute()
}
```

### A.3 TUI Styles (`internal/tui/styles.go`)

```go
package tui

import "github.com/charmbracelet/lipgloss"

var (
	ColorPrimary   = lipgloss.Color("99")
	ColorSecondary = lipgloss.Color("39")
	ColorAccent    = lipgloss.Color("212")
	ColorSuccess   = lipgloss.Color("78")
	ColorWarning   = lipgloss.Color("220")
	ColorError     = lipgloss.Color("196")
	ColorInfo      = lipgloss.Color("75")
	ColorMuted     = lipgloss.Color("240")

	StyleError = lipgloss.NewStyle().
		Foreground(ColorError).
		Bold(true).
		Render

	StyleSuccess = lipgloss.NewStyle().
		Foreground(ColorSuccess).
		Bold(true).
		Render

	StyleInfo = lipgloss.NewStyle().
		Foreground(ColorInfo).
		Render

	StyleWarning = lipgloss.NewStyle().
		Foreground(ColorWarning).
		Bold(true).
		Render

	StyleHeader = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary).
		Foreground(ColorSecondary).
		Align(lipgloss.Center).
		Width(60).
		Padding(1, 2)
)

func Error(msg string) {
	fmt.Fprintln(os.Stderr, StyleError("✗ ERROR")+" "+msg)
}

func Success(msg string) {
	fmt.Println(StyleSuccess("✓ SUCCESS") + " " + msg)
}

func Info(msg string) {
	fmt.Println(StyleInfo("ℹ INFO") + " " + msg)
}

func Warning(msg string) {
	fmt.Println(StyleWarning("⚠ WARNING") + " " + msg)
}
```

---

**Fim do Plano**
