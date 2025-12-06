# ghtools - GitHub Repository Management Tool

Ferramenta unificada para gerenciar repositórios do GitHub com interface interativa, busca fuzzy e seleção múltipla.

## Funcionalidades

- **Listagem de Repositórios**: Liste, filtre, ordene e exporte repositórios
- **Clone de Repositórios**: Clone múltiplos repositórios de uma vez
- **Sincronização**: Sincronize repositórios locais com remotes do GitHub
- **Criação de Repositórios**: Crie novos repositórios com templates
- **Exclusão de Repositórios**: Delete repositórios com segurança
- **Menu Interativo**: Interface amigável com fzf
- **Busca Fuzzy**: Encontre repositórios rapidamente
- **Interface Moderna**: TUI aprimorada com o uso de gum (recomendado)
- **Seleção Múltipla**: Gerencie vários repositórios simultaneamente
- **Confirmações de Segurança**: Proteção contra ações acidentais
- **Output Colorido**: Interface visual clara e intuitiva

## Instalação

### Método Rápido (Script Automático)

```bash
git clone https://github.com/quantmind-br/ghtools.git
cd ghtools
./install.sh
```

O script de instalação irá:
- Verificar dependências necessárias
- Instalar o ghtools em `~/scripts`
- Configurar o PATH automaticamente
- Detectar e remover configurações duplicadas

### Instalação Manual

```bash
git clone https://github.com/quantmind-br/ghtools.git
cd ghtools
mkdir -p ~/scripts
cp ghtools ~/scripts/ghtools
chmod +x ~/scripts/ghtools

# Adicionar ~/scripts ao PATH (apenas necessário uma vez)
echo 'export PATH="$HOME/scripts:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

## Dependências

- `gh` (GitHub CLI)
- `gum` (Terminal UI Toolkit) - Altamente recomendado para a TUI moderna
- `fzf` (Fuzzy Finder) - Essencial para fallback de seleção
- `git` (apenas para clone)

### Instalar no Arch Linux / CachyOS:

```bash
sudo pacman -S github-cli fzf git gum
```

ou

```bash
yay -S github-cli fzf git gum
```

## Uso

### Menu Interativo (Recomendado)

Simplesmente execute:

```bash
ghtools
```

Um menu interativo será exibido com as seguintes opções:
- 📋 List repositories
- 📦 Clone repositories
- 🔄 Sync repositories
- ➕ Create repository
- 🗑️  Delete repositories
- ❓ Help
- 🚪 Exit

### Comandos Diretos

```bash
ghtools list     # Liste e filtre repositórios
ghtools clone    # Clone repositórios
ghtools sync     # Sincronize repositórios locais
ghtools create   # Crie novo repositório
ghtools delete   # Delete repositórios
ghtools help     # Exibir ajuda
```

## Listagem de Repositórios

### Como Usar

```bash
ghtools list [OPTIONS]
```

### Opções Disponíveis

- `--lang LANGUAGE` - Filtre por linguagem de programação
- `--visibility public|private` - Filtre por visibilidade
- `--archived` - Inclua apenas repositórios arquivados
- `--no-archived` - Exclua repositórios arquivados (padrão)
- `--sort stars|created|updated|name` - Campo de ordenação
- `--order asc|desc` - Ordem de classificação (padrão: desc)
- `--export table|csv|json` - Formato de saída (padrão: table)
- `--limit N` - Limite de repositórios (padrão: 1000)

### Exemplos de Uso (List)

```bash
# Listar todos os repositórios
ghtools list

# Repositórios Python apenas
ghtools list --lang python

# Ordenar por estrelas
ghtools list --sort stars --order desc

# Repositórios privados criados recentemente
ghtools list --visibility private --sort created

# Exportar para CSV
ghtools list --export csv > repos.csv

# Exportar para JSON
ghtools list --export json > repos.json

# Filtros combinados
ghtools list --lang rust --sort stars --no-archived
```

### Formatos de Exportação

- **table** - Tabela formatada e colorida (padrão)
- **csv** - Formato CSV para importação em planilhas
- **json** - Formato JSON para processamento programático

## Sincronização de Repositórios

### Como Usar

```bash
ghtools sync [OPTIONS]
```

### Opções Disponíveis

- `--path DIR` - Diretório para escanear (padrão: diretório atual)
- `--max-depth N` - Profundidade máxima de busca (padrão: 5)
- `--all` - Sincronizar todos sem seleção interativa
- `--dry-run` - Mostrar o que seria feito sem executar

### Funcionalidades de Sincronização

1. **Descoberta Automática**: Encontra todos os repositórios Git do GitHub
2. **Verificação de Status**: Mostra ahead/behind/dirty para cada repo
3. **Seleção Interativa**: Escolha quais repositórios sincronizar
4. **Segurança**: Pula repositórios com mudanças não commitadas
5. **Fast-forward Only**: Usa `--ff-only` para evitar merges acidentais
6. **Resumo Detalhado**: Mostra sucessos, falhas e repositórios pulados

### Status dos Repositórios

- **✓ SYNCED** - Atualizado com o remote
- **↓ BEHIND** - Atrás do remote (precisa pull)
- **↑ AHEAD** - À frente do remote (precisa push)
- **⚠ DIRTY** - Com mudanças não commitadas

### Exemplos de Uso (Sync)

```bash
# Sincronizar repositórios no diretório atual
ghtools sync

# Sincronizar em diretório específico
ghtools sync --path ~/projects

# Modo dry-run (visualizar sem executar)
ghtools sync --dry-run

# Sincronizar todos automaticamente
ghtools sync --all

# Busca rasa (apenas 2 níveis)
ghtools sync --max-depth 2
```

### Comportamento Seguro

O comando sync:
- **Nunca** faz merge forçado
- **Pula** repositórios com mudanças não commitadas
- **Pula** repositórios com conflitos
- **Usa** `--ff-only` para garantir segurança
- **Exibe** mensagens claras sobre repositórios pulados

## Criação de Repositórios

### Como Usar

```bash
ghtools create [NOME] [OPTIONS]
```

### Opções Disponíveis

- `--description TEXT` - Descrição do repositório
- `--public` / `--private` - Visibilidade (padrão: prompt)
- `--readme` / `--no-readme` - Adicionar README.md
- `--license MIT|Apache-2.0|GPL-3.0|BSD-3-Clause` - Licença
- `--gitignore Python|Node|Go|Rust|Java|C++|Web` - Template .gitignore
- `--template python|node|go|rust|web` - Template de projeto
- `--clone` / `--no-clone` - Clonar após criar
- `--default-branch NOME` - Nome da branch padrão

### Templates de Projeto

#### Python
- `requirements.txt` - Dependências
- `main.py` - Script principal executável
- `pyproject.toml` - Configuração do projeto

#### Node.js
- `package.json` - Configuração e dependências
- `index.js` - Script principal executável

#### Go
- `go.mod` - Module definition
- `main.go` - Aplicação principal

#### Rust
- `Cargo.toml` - Configuração do projeto
- `src/main.rs` - Aplicação principal

#### Web
- `index.html` - Página principal
- `style.css` - Estilos
- `script.js` - JavaScript

### Exemplos de Uso (Create)

```bash
# Modo interativo (recomendado)
ghtools create

# Criação rápida
ghtools create my-api --public --readme

# Com template Python
ghtools create my-python-project --template python --clone

# Projeto completo
ghtools create my-app \
  --description "Minha aplicação incrível" \
  --private \
  --license MIT \
  --gitignore Node \
  --template node \
  --clone

# Criação sem clonar
ghtools create test-repo --public --no-clone
```

### Fluxo Interativo

1. Nome do repositório (validado)
2. Descrição (opcional)
3. Visibilidade (public/private)
4. Adicionar README.md?
5. Selecionar licença
6. Selecionar .gitignore template
7. Usar template de projeto?
8. Resumo e confirmação
9. Criação no GitHub
10. Aplicação de template (se selecionado)
11. Opção de clonar localmente

## Clone de Repositórios

### Como Usar

```bash
ghtools clone
```

### Funcionalidades do Clone

1. **Listagem Automática**: Lista todos os seus repositórios do GitHub
2. **Busca Fuzzy**: Filtre repositórios digitando qualquer parte do nome
3. **Seleção Múltipla**: Use TAB para selecionar múltiplos repositórios
4. **Verificação de Existência**: Pula repositórios já clonados
5. **Clone Paralelo**: Clona múltiplos repositórios sequencialmente
6. **Resumo Final**: Exibe sucessos, falhas e repositórios pulados

### Atalhos do Teclado (Clone)

- `TAB` - Selecionar/desselecionar repositório
- `CTRL+A` - Selecionar todos
- `CTRL+D` - Desselecionar todos
- `ENTER` - Confirmar seleção
- `ESC` - Cancelar

### Exemplo de Uso (Clone)

```
1. Execute: ghtools clone
2. Script lista todos os repositórios
3. Use busca fuzzy para filtrar (opcional)
4. Pressione TAB para selecionar repositórios
5. Pressione ENTER para confirmar
6. Digite Y para confirmar clonagem
7. Repositórios são clonados no diretório atual
8. Resumo final é exibido
```

## Exclusão de Repositórios

### Como Usar

```bash
ghtools delete
```

### Funcionalidades de Exclusão

1. **Verificação de Permissões**: Verifica scope `delete_repo`
2. **Listagem Automática**: Lista todos os repositórios
3. **Busca Fuzzy**: Filtre repositórios facilmente
4. **Seleção Múltipla**: Selecione múltiplos para deletar
5. **Confirmação Dupla**: Requer confirmação explícita antes de deletar
6. **Exclusão Segura**: Executa com tratamento de erros robusto
7. **Resumo Final**: Exibe sucessos e falhas

### Atalhos do Teclado (Delete)

- `TAB` - Selecionar/desselecionar repositório
- `CTRL+A` - Selecionar todos
- `CTRL+D` - Desselecionar todos
- `ENTER` - Confirmar seleção
- `ESC` - Cancelar

### Verificação de Permissões

O script verifica automaticamente se você tem o scope `delete_repo` necessário.

Se não tiver, execute:

```bash
gh auth refresh -s delete_repo
```

### Exemplo de Uso (Delete)

```
1. Execute: ghtools delete
2. Script verifica permissões
3. Lista todos os repositórios
4. Use busca fuzzy para filtrar (opcional)
5. Pressione TAB para selecionar repositórios
6. Pressione ENTER para confirmar seleção
7. Digite Y para confirmar exclusão
8. Repositórios são deletados
9. Resumo final é exibido
```

## Saídas Coloridas

O script usa cores para facilitar a leitura:

- **VERMELHO**: Avisos de exclusão e erros
- **VERDE**: Operações bem-sucedidas (clone, create, public)
- **AMARELO**: Avisos importantes (dirty repos, skip, private)
- **AZUL**: Informações gerais (ahead repos, processing)
- **CIANO**: Menu interativo, list, sync, headings

## Segurança

### Para List:
- Somente leitura, sem modificações
- Suporta exportação segura para CSV/JSON
- Filtros validados

### Para Clone:
- Verificação de diretórios existentes
- Pula repositórios já clonados
- Tratamento de erros robusto
- Lista detalhada de falhas

### Para Sync:
- **Nunca usa --force ou --hard**
- Usa `--ff-only` para evitar merges acidentais
- Pula repos com mudanças não commitadas
- Pula repos com conflitos
- Confirmação antes de executar
- Modo dry-run disponível

### Para Create:
- Validação de nome de repositório
- Confirmação antes de criar
- Templates testados e seguros
- Opção de não clonar localmente

### Para Delete:
- Múltiplas confirmações antes de deletar
- Mensagens de aviso claras e em vermelho
- Validação de autenticação e permissões
- Tratamento de erros robusto
- Lista de repositórios que falham na exclusão
- Confirmação explícita (Y/y required)

## Estrutura do Projeto

```
ghtools/
├── ghtools         # Script principal
├── install.sh      # Script de instalação
└── README.md       # Este arquivo
```

## Atualização

Para atualizar o script:

```bash
cd ghtools
git pull
./install.sh  # ou copie manualmente: cp ghtools ~/scripts/ghtools
```

## Desinstalação

Para remover o script:

```bash
rm ~/scripts/ghtools
```

Para remover também do PATH, edite seu arquivo de configuração do zsh (~/.zshrc, ~/.zshrc_custom, etc.) e remova a linha:

```bash
export PATH="$HOME/scripts:$PATH"
```

## Solução de Problemas

### "Command not found: ghtools"

Verifique se `~/scripts` está no seu PATH:

```bash
echo $PATH | grep scripts
```

Se não estiver, adicione ao seu ~/.zshrc:

```bash
export PATH="$HOME/scripts:$PATH"
source ~/.zshrc
```

### "Missing required dependencies"

Instale as dependências:

```bash
sudo pacman -S github-cli fzf git
```

### "Not authenticated with GitHub CLI"

Execute:

```bash
gh auth login
```

### "delete_repo scope missing"

Execute:

```bash
gh auth refresh -s delete_repo
```

## Contribuindo

Contribuições são bem-vindas! Sinta-se à vontade para:
- Reportar bugs
- Sugerir novas funcionalidades
- Enviar pull requests

## Aviso Importante

**EXCLUSÃO DE REPOSITÓRIOS NÃO PODE SER DESFEITA!**

Repositórios deletados não podem ser recuperados. Use a funcionalidade de delete com cautela e sempre verifique cuidadosamente os repositórios selecionados antes de confirmar.

## Licença

Este projeto é disponibilizado como está, sem garantias.

## Autor

Desenvolvido para facilitar o gerenciamento de repositórios GitHub via linha de comando.
