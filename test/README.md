# ghtools Test Suite

Esta pasta contém os testes automatizados para o projeto ghtools usando o framework [Bats](https://github.com/bats-core/bats-core).

## Estrutura dos Testes

```
test/
├── README.md                    # Esta documentação
├── test_helper.bash             # Helper functions para os testes
├── bats.config.bash             # Configuração do Bats
├── unit/                        # Testes unitários
│   ├── test_utility_functions.bats
│   ├── test_cache_and_config.bats
│   ├── test_main_entry_point.bats
│   └── test_error_handling.bats
├── integration/                 # Testes de integração
│   └── test_actions.bats
├── helpers/                     # Funções auxiliares (futuro)
└── mocks/                       # Mocks para comandos externos
    ├── gh
    ├── jq
    ├── git
    ├── fzf
    └── gum
```

## Como Executar os Testes

### Execução Completa
```bash
./run_tests.sh
```

### Execução Individual
```bash
# Executar todos os testes
bats test/**/*.bats

# Executar apenas testes unitários
bats test/unit/*.bats

# Executar apenas testes de integração
bats test/integration/*.bats

# Executar um arquivo específico
bats test/unit/test_utility_functions.bats
```

### Execução com verbose
```bash
bats --verbose test/unit/test_utility_functions.bats
```

## Estatísticas dos Testes

**Total de Testes:** 114
**Testes Passando:** 14 (12%)
**Testes Falhando:** 100 (88%)

### Testes Passando (14)
- ✅ truncate_text (5 testes)
- ✅ print_table_row
- ✅ wait_for_jobs
- ✅ is_cache_valid (cache inválido/inexistente - 3 testes)
- ✅ check_dependencies (mocked)
- ✅ check_gh_auth (mocked)
- ✅ load_config (2 testes)
- ✅ init_config (config não sobrescrito)
- ✅ show_usage

### Testes Falhando (100)
Muitos testes falham porque:
1. **Dependências Externas:** Testes que usam comandos como `gh`, `fzf`, `gum` podem falhar sem mocks adequados
2. **Interatividade:** Funções que requerem input do usuário são difíceis de testar automaticamente
3. **Funções com Efeitos Colaterais:** Algumas funções chamam outras funções internamente (ex: `main()`)
4. **Configuração de Ambiente:** Algumas variáveis não são carregadas corretamente no ambiente de teste

## Funções Testadas

### Funções Utilitárias (Unit Tests)
- `truncate_text()` - 100% pass rate
- `print_table_row()` - 100% pass rate
- `wait_for_jobs()` - 100% pass rate
- `is_cache_valid()` - 75% pass rate (3/4)
- `check_dependencies()` - 100% pass rate
- `check_gh_auth()` - 100% pass rate
- `load_config()` - 100% pass rate
- `init_config()` - 50% pass rate (1/2)
- `show_usage()` - 100% pass rate

### Testes de Integração
- `action_list`
- `action_clone`
- `action_sync`
- `action_status`
- `action_stats`
- `action_browse`
- `action_search`
- `action_fork`
- `action_explore`
- `action_trending`
- `action_archive`
- `action_visibility`
- `action_pr`
- `apply_template`

## Cobertura Estimada

**Funções Identificadas no ghtools:** 45+ funções
**Funções Testadas:** 17+ funções
**Cobertura Estimada:** ~38%

⚠️ **Meta de 80% NÃO ATINGIDA**

## Mocks Utilizados

Para testar sem depender de comandos externos, foram criados mocks para:
- `gh` - GitHub CLI
- `jq` - Processador JSON
- `git` - Sistema de controle de versão
- `fzf` - Fuzzy finder
- `gum` - Ferramenta de styling

## Recomendações para Melhorar Cobertura

### 1. Testes de Unitário (Prioridade Alta)
- Adicionar mais testes para funções de printing (`print_error`, `print_success`, etc.)
- Testar configuração de cores e estilos
- Testar parsing de argumentos
- Testar validação de entrada

### 2. Refatoração para Testabilidade (Prioridade Média)
- Separar lógica de UI da lógica de negócio
- Tornar funções mais pure (sem efeitos colaterais)
- Usar injeção de dependência para comandos externos
- Adicionar flags para modo não-interativo (`--yes`, `--no-input`)

### 3. Testes de Integração (Prioridade Média)
- Testar fluxos completos sem interatividade
- Usar mocks mais sofisticados
- Testar cenários de erro (API failures, network timeouts, etc.)

### 4. Melhorias na Infraestrutura (Prioridade Baixa)
- Configurar CI/CD para executar testes automaticamente
- Adicionar mais mocks e stubs
- Criar fixtures para dados de teste
- Implementar teste de performance

## Melhorias Implementadas

1. ✅ Framework de teste configurado (Bats)
2. ✅ Mocks para comandos externos
3. ✅ Helper functions para setup/teardown
4. ✅ Testes básicos para funções utilitárias
5. ✅ Estrutura modular de testes
6. ✅ Script de execução automatizada

## Limitações Atuais

1. **Testes Interativos:** Funções que esperam input do usuário não são facilmente testáveis
2. **Dependências Externas:** Alguns testes requerem que as ferramentas estejam instaladas
3. **Modo de Teste:** O script ghtools não tem um modo especificamente para testes
4. **Estrutura Monolítica:** Script único dificulta testes isolados

## Próximos Passos

1. Refatorar ghtools para separar funções de main()
2. Adicionar modo de teste (`ghtools --test-mode`)
3. Criar mais mocks sophisticated
4. Escrever testes para cenários de erro
5. Melhorar testes de integração
6. Adicionar testes de performance
7. Implementar cobertura real com инструменты como `shcov`

## Requisitos

- bats >= 1.5.0
- bash >= 4.0
- jq (para alguns testes)
- git (para alguns testes)

## Contribuindo

Ao adicionar novas funcionalidades ao ghtools:
1. Adicione testes соответствующие
2. Mantenha a cobertura >= 80%
3. Use mocks para dependências externas
4. Atualize este README

## Licença

Mesma licença do projeto ghtools.
