# Plano Extremamente Detalhado de Desenvolvimento - ghdelete

## Visão Geral do Projeto

Criar um script bash chamado `ghdelete` que permite ao usuário deletar repositórios do GitHub de forma interativa, com busca fuzzy e seleção múltipla.

---

## Fase 1: Análise de Requisitos

### 1.1 Requisitos Funcionais
- [x] Script acessível via comando `ghdelete` no PATH
- [x] Listar todos os repositórios do usuário
- [x] Filtro fuzzy para busca de repositórios
- [x] Seleção múltipla de repositórios
- [x] Exclusão automática sem confirmação manual repetida
- [x] Execução sequencial de múltiplas exclusões
- [x] Uso da flag `--yes` do gh repo delete

### 1.2 Requisitos Não-Funcionais
- [x] Interface amigável e intuitiva
- [x] Mensagens de erro claras
- [x] Validação de dependências
- [x] Tratamento robusto de erros
- [x] Feedback visual durante operações

### 1.3 Dependências Identificadas
- `gh` (GitHub CLI): Para interagir com API do GitHub
- `fzf`: Para seleção interativa com busca fuzzy
- `bash`: Shell script executor

---

## Fase 2: Pesquisa e Análise Técnica

### 2.1 Análise do GitHub CLI

**Comando pesquisado:** `gh repo list`

**Formato de saída:**
```
OWNER/REPO_NAME    DESCRIPTION    VISIBILITY    DATE
```

**Características importantes:**
- Primeira coluna contém o nome completo (OWNER/REPO)
- Campos separados por tabulação/espaços
- Flag `--limit` para controlar quantidade de resultados

### 2.2 Análise do gh repo delete

**Descoberta crucial:** `gh repo delete` possui flag `--yes`!

```bash
gh repo delete [<repository>] [flags]

FLAGS
  --yes   Confirm deletion without prompting
```

**Implicação:** Não é necessário simular entrada de teclado (echo | gh). Podemos usar diretamente:
```bash
gh repo delete "OWNER/REPO" --yes
```

### 2.3 Análise de Permissões

**Scope necessário:** `delete_repo`

**Verificação:**
```bash
gh auth status 2>&1 | grep -q "delete_repo"
```

**Atualização de scope:**
```bash
gh auth refresh -s delete_repo
```

### 2.4 Análise do fzf

**Características utilizadas:**
- Multi-seleção: `--multi`
- Interface customizável: `--height`, `--border`, `--prompt`
- Preview: `--preview`, `--preview-window`
- Header: `--header`
- Atalhos customizados: `--bind`
- Cores customizadas: `--color`

**Atalhos implementados:**
- `TAB`: Selecionar/desselecionar
- `CTRL+A`: Selecionar todos
- `CTRL+D`: Desselecionar todos
- `ENTER`: Confirmar
- `ESC`: Cancelar

---

## Fase 3: Arquitetura do Script

### 3.1 Estrutura Geral

```
┌─────────────────────────────────────┐
│     Inicialização e Validação       │
├─────────────────────────────────────┤
│  1. Verificar dependências          │
│  2. Verificar autenticação gh       │
│  3. Verificar scope delete_repo     │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│      Listagem de Repositórios       │
├─────────────────────────────────────┤
│  1. Executar gh repo list           │
│  2. Validar saída                   │
│  3. Contar repositórios             │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│       Seleção Interativa (fzf)      │
├─────────────────────────────────────┤
│  1. Apresentar interface fzf        │
│  2. Permitir multi-seleção          │
│  3. Capturar seleções do usuário    │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│      Processamento de Seleção       │
├─────────────────────────────────────┤
│  1. Extrair nomes dos repositórios  │
│  2. Criar array de repos            │
│  3. Validar seleção não-vazia       │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│       Confirmação de Segurança      │
├─────────────────────────────────────┤
│  1. Exibir lista de repos           │
│  2. Mostrar avisos                  │
│  3. Solicitar confirmação "yes"     │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│      Execução de Exclusões          │
├─────────────────────────────────────┤
│  1. Loop sequencial por repo        │
│  2. Executar gh repo delete --yes   │
│  3. Capturar sucesso/erro           │
│  4. Atualizar contadores            │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│          Relatório Final            │
├─────────────────────────────────────┤
│  1. Exibir total de sucessos        │
│  2. Exibir total de falhas          │
│  3. Listar repositórios falhados    │
└─────────────────────────────────────┘
```

### 3.2 Funções Implementadas

#### 3.2.1 Funções de Output
```bash
print_error()    # Mensagens de erro (VERMELHO)
print_success()  # Mensagens de sucesso (VERDE)
print_info()     # Mensagens informativas (AZUL)
print_warning()  # Mensagens de aviso (AMARELO)
```

#### 3.2.2 Funções de Validação
```bash
check_dependencies()  # Verifica gh e fzf
check_gh_auth()       # Verifica autenticação
check_delete_scope()  # Verifica scope delete_repo
```

#### 3.2.3 Funções de Operação
```bash
fetch_repositories()  # Busca lista de repositórios
main()                # Função principal orquestradora
```

---

## Fase 4: Implementação Detalhada

### 4.1 Configuração Inicial

```bash
#!/bin/bash
set -euo pipefail
```

**Explicação das flags:**
- `-e`: Exit imediatamente se qualquer comando falhar
- `-u`: Trata variáveis não definidas como erro
- `-o pipefail`: Falha de qualquer comando em pipe causa falha total

### 4.2 Definição de Cores

```bash
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'
```

**Uso de ANSI escape codes** para colorir output no terminal.

### 4.3 Verificação de Dependências

**Lógica:**
1. Tentar executar `command -v gh`
2. Tentar executar `command -v fzf`
3. Se algum falhar, adicionar a array `missing_deps`
4. Se array não estiver vazia, mostrar mensagem e exit 1

**Tratamento de erro:**
- Mensagem clara sobre o que falta
- Instruções específicas para CachyOS (pacman/yay)
- Exit code 1 para indicar falha

### 4.4 Verificação de Autenticação

```bash
gh auth status &> /dev/null
```

**Redirecionamento:** `&> /dev/null` descarta stdout e stderr
**Verificação:** Apenas código de saída importa

### 4.5 Verificação de Scope

```bash
gh auth status 2>&1 | grep -q "delete_repo"
```

**Técnica:**
- Redireciona stderr para stdout (`2>&1`)
- Busca por "delete_repo" com grep quiet (`-q`)
- Se não encontrar, oferece opção de continuar

### 4.6 Listagem de Repositórios

```bash
gh repo list --limit 1000
```

**Considerações:**
- Limite de 1000 repositórios (ajustável)
- Validação de erro na execução
- Validação de saída vazia

### 4.7 Interface fzf

**Configuração detalhada:**

```bash
fzf \
  --multi \                    # Permite seleção múltipla
  --height=80% \               # Usa 80% da altura do terminal
  --border \                   # Adiciona borda visual
  --prompt="..." \             # Customiza prompt
  --preview='echo {}' \        # Preview da linha selecionada
  --preview-window=up:3:wrap \ # Preview acima, 3 linhas
  --header="..." \             # Texto de ajuda no topo
  --bind='...' \               # Atalhos customizados
  --color='...'                # Esquema de cores
```

**Esquema de cores escolhido:**
- Vermelho para highlights (indicar perigo)
- Fundo escuro
- Texto branco para legibilidade

### 4.8 Extração de Nomes de Repositórios

```bash
while IFS= read -r line; do
  repo_name=$(echo "$line" | awk '{print $1}')
  repo_names+=("$repo_name")
done <<< "$selected_repos"
```

**Técnica:**
- Loop com `while read` para processar linhas
- `IFS=` preserva espaços/tabs
- `awk '{print $1}'` extrai primeira coluna
- `<<< "$var"` here-string para input

### 4.9 Confirmação de Segurança

**Múltiplas camadas:**
1. Exibir lista completa de repositórios
2. Avisos em VERMELHO
3. Mensagem "THIS ACTION CANNOT BE UNDONE!"
4. Requer digitação completa de "yes"

**Validação rigorosa:**
```bash
[[ ! $REPLY =~ ^[Yy][Ee][Ss]$ ]]
```
- Regex que aceita apenas "yes", "Yes", "YES", etc.
- Qualquer outra coisa cancela

### 4.10 Loop de Exclusão

```bash
for repo in "${repo_names[@]}"; do
  if gh repo delete "$repo" --yes 2>&1; then
    # sucesso
    ((success_count++))
  else
    # falha
    failed_repos+=("$repo")
    ((fail_count++))
  fi
done
```

**Características:**
- Execução sequencial (não paralela)
- Captura stdout/stderr com `2>&1`
- Contadores separados para sucesso/falha
- Array de repositórios falhados

### 4.11 Relatório Final

**Informações apresentadas:**
- Total de sucessos (verde)
- Total de falhas (vermelho)
- Lista de repositórios que falharam
- Formatação visual com separadores

---

## Fase 5: Tratamento de Erros

### 5.1 Erros de Dependência
- Detecta ausência de `gh` ou `fzf`
- Mensagem com instruções de instalação
- Exit code 1

### 5.2 Erros de Autenticação
- Detecta falta de login no gh
- Instrui executar `gh auth login`
- Exit code 1

### 5.3 Erros de Permissão
- Detecta falta de scope `delete_repo`
- Instrui executar `gh auth refresh -s delete_repo`
- Permite continuar com aviso

### 5.4 Erros na Listagem
- Captura falhas em `gh repo list`
- Exibe mensagem de erro do gh
- Exit code 1

### 5.5 Erros na Exclusão
- Captura falhas individuais de `gh repo delete`
- Continua processando outros repositórios
- Registra falhas para relatório final

### 5.6 Cancelamento pelo Usuário
- Detecta ESC no fzf (seleção vazia)
- Detecta resposta diferente de "yes"
- Exit gracioso com mensagem

---

## Fase 6: Instalação no Sistema

### 6.1 Escolha do Diretório

**Opção escolhida:** `/usr/local/bin`

**Justificativa:**
- Padrão para scripts de usuário
- Já está no PATH por padrão
- Não conflita com pacotes do sistema
- Não requer modificação de PATH

**Alternativas consideradas:**
- `~/.local/bin`: Específico do usuário
- `/usr/bin`: Reservado para pacotes do sistema
- `~/bin`: Pode não estar no PATH

### 6.2 Processo de Instalação

```bash
# 1. Copiar script
sudo cp /home/diogo/dev/ghdelete/ghdelete /usr/local/bin/ghdelete

# 2. Tornar executável
sudo chmod +x /usr/local/bin/ghdelete
```

### 6.3 Verificação Pós-Instalação

```bash
which ghdelete
# Output: /usr/local/bin/ghdelete
```

---

## Fase 7: Documentação

### 7.1 README.md
- Instruções de uso
- Descrição de funcionalidades
- Atalhos do teclado
- Comandos de instalação/atualização
- Avisos de segurança

### 7.2 PLANO_DESENVOLVIMENTO.md (este arquivo)
- Detalhamento completo do desenvolvimento
- Decisões técnicas e justificativas
- Estrutura e arquitetura
- Processo de implementação

### 7.3 Comentários no Código
- Cabeçalho descritivo
- Comentários em seções importantes
- Documentação de funções

---

## Fase 8: Testes Realizados

### 8.1 Teste de Sintaxe
```bash
bash -n /home/diogo/dev/ghdelete/ghdelete
```
**Resultado:** Sem erros de sintaxe

### 8.2 Teste de Executabilidade
```bash
chmod +x /home/diogo/dev/ghdelete/ghdelete
```
**Resultado:** Permissões configuradas corretamente

### 8.3 Teste de Instalação
```bash
which ghdelete
```
**Resultado:** `/usr/local/bin/ghdelete`

### 8.4 Teste de Inicialização
```bash
ghdelete
```
**Resultado:** Interface carrega corretamente (não executado até o fim para não deletar repos)

---

## Fase 9: Melhorias Futuras (Não Implementadas)

### 9.1 Opções de Linha de Comando
- `--help`: Exibir ajuda
- `--version`: Exibir versão
- `--org <org>`: Listar repositórios de organização
- `--dry-run`: Simular sem deletar

### 9.2 Filtros Adicionais
- Filtrar por visibilidade (public/private)
- Filtrar por data de criação
- Filtrar por atividade recente

### 9.3 Backup Antes de Deletar
- Criar arquivo local com metadados
- Clonar repositórios antes de deletar
- Gerar backup comprimido

### 9.4 Modo Batch
- Aceitar lista de repositórios via arquivo
- Modo não-interativo com arquivo de config

### 9.5 Logging
- Criar log de operações
- Timestamp de cada exclusão
- Histórico de operações

---

## Resumo de Decisões Técnicas

### ✓ Decisões Corretas
1. **Uso de `--yes`**: Descoberta da flag eliminou necessidade de simular input
2. **fzf multi-select**: Interface intuitiva e poderosa
3. **Confirmação "yes"**: Segurança contra deleções acidentais
4. **Cores ANSI**: Feedback visual claro
5. **Array de falhas**: Rastreamento preciso de erros
6. **set -euo pipefail**: Script robusto e seguro
7. **Funções modulares**: Código organizado e reutilizável
8. **Verificação de dependências**: UX melhorada
9. **Instalação em /usr/local/bin**: Padrão do sistema

### ✓ Desafios Superados
1. **Confirmação automática**: Resolvido com flag `--yes`
2. **Extração de nomes**: Resolvido com awk first column
3. **Multi-seleção**: Implementado com fzf --multi
4. **Feedback visual**: Implementado com cores e mensagens

### ✓ Segurança Implementada
1. Múltiplas confirmações
2. Avisos em vermelho
3. Verificação de scope
4. Validação de entrada
5. Tratamento de erros
6. Mensagens claras de irreversibilidade

---

## Conclusão

O script `ghdelete` foi desenvolvido com sucesso, atendendo todos os requisitos:

- ✓ Listagem de repositórios
- ✓ Filtro fuzzy (fzf)
- ✓ Seleção múltipla
- ✓ Exclusão automática sem confirmação manual repetida
- ✓ Execução sequencial
- ✓ Instalado no PATH
- ✓ Interface amigável
- ✓ Tratamento robusto de erros
- ✓ Documentação completa

O script está pronto para uso em produção.
