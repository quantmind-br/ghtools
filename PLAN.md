# Plano de Correção e Refatoração: ghtools

Este documento detalha o plano de ação para corrigir vulnerabilidades de segurança críticas, bugs de lógica e problemas de compatibilidade identificados na versão 3.1.0 do script `ghtools`.

**Objetivo:** Tornar o script seguro, robusto em ambientes multi-usuário e compatível com estruturas de diretórios complexas, sem alterar a funcionalidade core existente.

**Versão do Plano:** 2.0 (Revisado)

---

## Fase 1: Segurança Crítica (Prioridade Alta)

### 1.1. Eliminação de Injeção de Comandos (`eval`) em `action_create`
**Localização:** Linhas 904-909

O uso de `eval` para executar comandos construídos via string é a vulnerabilidade mais grave. Deve ser substituído pelo uso de **Arrays do Bash**, que tratam argumentos de forma segura automaticamente.

**Código Atual (Vulnerável):**
```bash
local cmd="gh repo create $name --$vis"
[ -n "$desc" ] && cmd="$cmd --description \"$desc\""
if eval "$cmd --clone"; then
```

**Código Corrigido:**
```bash
local cmd_args=("gh" "repo" "create" "$name" "--$vis")
if [ -n "$desc" ]; then
    cmd_args+=("--description" "$desc")
fi
if "${cmd_args[@]}" --clone; then
```

### 1.2. Eliminação de Injeção de Comandos (`eval`) em `action_explore`
**Localização:** Linhas 1114-1118

**Código Atual (Vulnerável):**
```bash
local gh_cmd="gh search repos \"$search_query\" --limit $limit --sort $sort_by"
[ -n "$language" ] && gh_cmd="$gh_cmd --language \"$language\""
results=$(eval "$gh_cmd" --json fullName,description,...)
```

**Código Corrigido:**
```bash
local cmd_args=("gh" "search" "repos" "$search_query" "--limit" "$limit" "--sort" "$sort_by")
[ -n "$language" ] && cmd_args+=("--language" "$language")
results=$("${cmd_args[@]}" --json fullName,description,...)
```

### 1.3. Eliminação de Injeção de Comandos (`eval`) em `run_with_spinner`
**Localização:** Linhas 212-228

**Problema:** A função `run_with_spinner` usa `eval` tanto no fallback quanto indiretamente via `bash -c`. Isso é perigoso se comandos contiverem entrada do usuário.

**Código Atual (Vulnerável):**
```bash
run_with_spinner() {
  local title="$1"
  shift
  local cmd="$@"
  if use_gum; then
    gum spin ... -- bash -c "$cmd"
  else
    if eval "$cmd"; then
```

**Solução:** Esta função deve receber um **array** ou ser refatorada para executar comandos diretamente. Como ela é usada internamente apenas com comandos fixos, uma abordagem segura é:

```bash
run_with_spinner() {
  local title="$1"
  shift
  # Executa diretamente o comando passado como argumentos separados
  if use_gum; then
    gum spin --spinner dot --spinner.foreground "$COLOR_ACCENT" --title "$title" -- "$@"
  else
    echo -n "$title... "
    if "$@"; then
      echo "done"
    else
      echo "failed"
      return 1
    fi
  fi
}
```
**Nota:** Verificar todos os call sites desta função para garantir compatibilidade.

### 1.4. Proteção do Arquivo de Cache (Permissões)
**Localização:** Linha 488

**Problema:** O arquivo de cache é criado sem controle de permissões, potencialmente expondo informações sobre repositórios privados a outros usuários.

**Ação:** Definir `umask` restritivo antes de criar o cache ou usar `install` com permissões explícitas.

```bash
# Antes de escrever no cache:
(umask 077 && $gh_cmd --limit "$limit" --json "$fields" > "$CACHE_FILE")
```

---

## Fase 2: Segurança Alta (Hardening)

### 2.1. Validação de Arquivo de Configuração (`source`)
**Localização:** Linhas 22-27

**Problema:** O uso de `source "$CONFIG_FILE"` executa qualquer código presente no arquivo de configuração. Se um atacante conseguir modificar este arquivo, pode executar código arbitrário.

**Mitigação Recomendada:**
*   Usar um parser mais seguro que apenas aceite atribuições de variáveis específicas.
*   Ou validar o conteúdo antes de source.

**Código Sugerido:**
```bash
load_config() {
    if [ -f "$CONFIG_FILE" ]; then
        # Valida que o arquivo contém apenas atribuições de variáveis conhecidas
        local allowed_vars="CACHE_TTL|CACHE_FILE|MAX_JOBS|DEFAULT_ORG|DEFAULT_CLONE_PATH"
        if grep -Evq "^[[:space:]]*(#.*)?$|^[[:space:]]*($allowed_vars)=" "$CONFIG_FILE"; then
            print_warning "Config file contains invalid lines. Skipping."
            return 1
        fi
        # shellcheck source=/dev/null
        source "$CONFIG_FILE"
    fi
}
```

**Alternativa Simples:** Documentar que o arquivo de config deve ter permissões 600 e adicionar verificação:
```bash
if [ -f "$CONFIG_FILE" ]; then
    local perms
    perms=$(stat -c %a "$CONFIG_FILE" 2>/dev/null || stat -f %Lp "$CONFIG_FILE")
    if [ "$perms" != "600" ]; then
        print_warning "Config file has insecure permissions. Run: chmod 600 $CONFIG_FILE"
    fi
    source "$CONFIG_FILE"
fi
```

---

## Fase 3: Confiabilidade e Ambiente (Prioridade Média)

### 3.1. Correção de Manipulação de Caminhos com Espaços
**Localização:** Linhas 694 e 1847

A lógica atual quebra caminhos de diretórios que contêm espaços ao usar `xargs` sem delimitadores nulos.

**Código Atual:**
```bash
mapfile -t git_dirs < <(find "$base_path" ... | xargs -n1 dirname)
```

**Código Corrigido:**
```bash
mapfile -t git_dirs < <(find "$base_path" -maxdepth "$max_depth" -name ".git" -type d -prune -print0 | xargs -0 -n1 dirname)
```

**Funções Afetadas:** `action_sync` e `action_status`

### 3.2. Resolução de Conflito de Permissões (Cache Multi-usuário)
**Localização:** Linha 15

O arquivo de cache fixo em `/tmp` impede que múltiplos usuários usem a ferramenta na mesma máquina.

**Código Atual:**
```bash
CACHE_FILE="/tmp/ghtools_repos.json"
```

**Código Corrigido:**
```bash
CACHE_FILE="/tmp/ghtools_repos_$(id -u).json"
```

---

## Fase 4: Lógica e Compatibilidade (Prioridade Média/Baixa)

### 4.1. Tratamento de "Detached HEAD" em PRs
**Localização:** Linha 1793-1794

Evitar falhas silenciosas ou erros confusos ao tentar criar PRs sem estar em um branch nomeado.

**Código Atual:**
```bash
local current_branch
current_branch=$(git branch --show-current)

if [ "$current_branch" = "main" ] || [ "$current_branch" = "master" ]; then
```

**Código Corrigido:**
```bash
local current_branch
current_branch=$(git branch --show-current)

if [ -z "$current_branch" ]; then
    print_error "You are in 'detached HEAD' state. Please checkout a branch first."
    return 1
fi

if [ "$current_branch" = "main" ] || [ "$current_branch" = "master" ]; then
```

### 4.2. Compatibilidade de `wait -n`
**Localização:** Linha 512

Garantir que o script não quebre em versões antigas do Bash que não suportam a flag `-n`.

**Código Atual:**
```bash
wait_for_jobs() {
  local current_jobs
  current_jobs=$(jobs -p | wc -l)
  if [ "$current_jobs" -ge "$MAX_JOBS" ]; then
    wait -n
  fi
}
```

**Código Corrigido:**
```bash
wait_for_jobs() {
    local current_jobs
    current_jobs=$(jobs -p | wc -l)
    if [ "$current_jobs" -ge "$MAX_JOBS" ]; then
        # Tenta usar wait -n (bash 4.3+), se falhar usa wait (espera todos)
        wait -n 2>/dev/null || wait
    fi
}
```

### 4.3. Push Automático Após Template (UX)
**Localização:** Linha 915

**Problema:** Após criar um repositório com template, o script faz push automático sem confirmação do usuário.

**Código Atual:**
```bash
(cd "$name" && git add . && git commit -m "Initial commit (Template: $tpl)" && git push origin HEAD) &>/dev/null
```

**Sugestão:** Adicionar confirmação ou flag `--no-push`:
```bash
if gum_confirm "Push initial commit to origin?" "yes"; then
    (cd "$name" && git add . && git commit -m "Initial commit (Template: $tpl)" && git push origin HEAD) &>/dev/null
    print_success "Template applied and pushed"
else
    (cd "$name" && git add . && git commit -m "Initial commit (Template: $tpl)") &>/dev/null
    print_success "Template applied (not pushed)"
fi
```

---

## Fase 5: Melhorias de Qualidade (Baixa Prioridade)

### 5.1. Refatorar `fetch_repositories_json` para usar Array
**Localização:** Linhas 475-488

**Problema:** `$gh_cmd` é construído como string e expandido sem quotes em alguns lugares.

**Código Atual:**
```bash
local gh_cmd="gh repo list"
if [ -n "$org_filter" ]; then
    gh_cmd="gh repo list $org_filter"
fi
...
if ! $gh_cmd --limit "$limit" --json "$fields" >"$CACHE_FILE"
```

**Código Corrigido:**
```bash
local cmd_args=("gh" "repo" "list")
if [ -n "$org_filter" ]; then
    cmd_args+=("$org_filter")
fi
...
if ! "${cmd_args[@]}" --limit "$limit" --json "$fields" >"$CACHE_FILE"
```

---

## Fase 6: Plano de Verificação

Após as alterações, os seguintes testes devem ser executados:

### Testes de Segurança
1.  **Teste de Injeção em `create`:** Tentar criar um repositório com o nome `test; echo VULNERAVEL`. O script deve tentar criar um repo com esse nome literal e falhar (validado pelo GitHub), mas **não** deve imprimir "VULNERAVEL" no terminal.

2.  **Teste de Injeção em `explore`:** Buscar com query `"; rm -rf /; echo "`. Deve buscar literalmente essa string, sem executar comandos.

3.  **Teste de Permissões do Cache:** Verificar que o arquivo de cache é criado com permissões `600` (apenas owner pode ler/escrever).

### Testes de Robustez
4.  **Teste de Espaços:** Criar uma pasta `"/tmp/test space/repo"`, iniciar um git nela e rodar `ghtools sync --path "/tmp/test space"`. O repositório deve ser encontrado.

5.  **Teste Multi-usuário:** Verificar se o arquivo criado em `/tmp` possui o ID do usuário no nome (ex: `ghtools_repos_1000.json`).

### Testes de Lógica
6.  **Teste Detached HEAD:** Fazer checkout de um commit específico (`git checkout <hash>`) e rodar `ghtools pr create`. Deve exibir erro amigável.

7.  **Teste wait -n Fallback:** Executar em container com Bash 4.2 (sem suporte a `wait -n`) e verificar que `ghtools sync` funciona sem erros.

---

## Resumo de Prioridades

| Prioridade | Item | Risco se não corrigido |
|------------|------|------------------------|
| **P0** | 1.1-1.3 Eliminação de `eval` | RCE (Execução remota de código) |
| **P0** | 1.4 Permissões do cache | Vazamento de dados privados |
| **P1** | 2.1 Validação de config | RCE via arquivo de config |
| **P1** | 3.1 Caminhos com espaços | Falha silenciosa em repos |
| **P2** | 3.2 Cache multi-usuário | Conflito entre usuários |
| **P2** | 4.1 Detached HEAD | UX ruim |
| **P3** | 4.2 Compatibilidade wait -n | Falha em Bash antigo |
| **P3** | 4.3 Push automático | UX ruim |
| **P4** | 5.1 Refatorar fetch | Código mais limpo |
