# ghdelete - GitHub Repository Deletion Tool

Ferramenta interativa para deletar repositórios do GitHub com busca fuzzy e seleção múltipla.

## Instalação

### Método Rápido (Script Automático)

```bash
git clone https://github.com/quantmind-br/ghdelete.git
cd ghdelete
./install.sh
```

### Instalação Manual

```bash
git clone https://github.com/quantmind-br/ghdelete.git
cd ghdelete
mkdir -p ~/scripts
cp ghdelete ~/scripts/ghdelete
chmod +x ~/scripts/ghdelete

# Adicionar ~/scripts ao PATH (apenas necessário uma vez)
echo 'export PATH="$HOME/scripts:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

## Dependências

- `gh` (GitHub CLI)
- `fzf` (Fuzzy Finder)

Instalar dependências no Arch Linux / CachyOS:

```bash
sudo pacman -S github-cli fzf
```

ou

```bash
yay -S github-cli fzf
```

## Uso

Simplesmente execute:

```bash
ghdelete
```

## Funcionalidades

### 1. Listagem Automática
O script lista automaticamente todos os seus repositórios do GitHub.

### 2. Busca Fuzzy
Use o fzf para buscar repositórios:
- Digite qualquer parte do nome do repositório
- A busca é incremental e fuzzy

### 3. Seleção Múltipla
Atalhos do teclado:
- `TAB` - Selecionar/desselecionar repositório
- `CTRL+A` - Selecionar todos
- `CTRL+D` - Desselecionar todos
- `ENTER` - Confirmar seleção
- `ESC` - Cancelar

### 4. Confirmação de Segurança
- Antes de deletar, o script mostra todos os repositórios selecionados
- Você deve digitar "yes" (completo) para confirmar
- Qualquer outra resposta cancela a operação

### 5. Exclusão Automática
- Usa `gh repo delete --yes` para deletar sem confirmação adicional
- Executa as exclusões sequencialmente
- Mostra progresso em tempo real
- Exibe resumo final com sucessos e falhas

## Verificação de Permissões

O script verifica automaticamente se você tem o scope `delete_repo` necessário para deletar repositórios.

Se não tiver, execute:

```bash
gh auth refresh -s delete_repo
```

## Saídas Coloridas

O script usa cores para facilitar a leitura:
- VERMELHO: Avisos e erros
- VERDE: Operações bem-sucedidas
- AMARELO: Avisos importantes
- AZUL: Informações gerais

## Segurança

- Múltiplas confirmações antes de deletar
- Mensagens de aviso claras
- Validação de autenticação e permissões
- Tratamento de erros robusto
- Lista de repositórios que falham na exclusão

## Exemplo de Fluxo

```
1. Execute: ghdelete
2. Script lista todos os repositórios
3. Use busca fuzzy para filtrar (opcional)
4. Pressione TAB para selecionar repositórios
5. Pressione ENTER para confirmar seleção
6. Digite "y" para confirmar exclusão
7. Script deleta e mostra progresso
8. Resumo final exibido
```

## Estrutura do Script

```bash
~/scripts/ghdelete  # Script instalado
```

## Atualização

Para atualizar o script:

```bash
cd ghdelete
git pull
cp ghdelete ~/scripts/ghdelete
```

## Desinstalação

Para remover o script:

```bash
rm ~/scripts/ghdelete
```

## Aviso

**ESTA OPERAÇÃO NÃO PODE SER DESFEITA!**

Repositórios deletados não podem ser recuperados. Use com cautela.
