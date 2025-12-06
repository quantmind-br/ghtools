# Checklist de Melhorias para Testes - ghtools

## Status Atual
- ✅ Framework Bats configurado
- ✅ 114 testes implementados (14 passando, 100 falhando)
- ✅ Cobertura atual: ~38%
- ❌ **Meta de 80% NÃO ATINGIDA**

## Testes Funcionando (14)

### Funções Utilitárias
- ✅ `truncate_text()` - 5/5 testes passing
- ✅ `print_table_row()` - 1/1 teste passing
- ✅ `wait_for_jobs()` - 1/1 teste passing
- ✅ `check_dependencies()` - 1/1 teste passing
- ✅ `check_gh_auth()` - 1/1 teste passing
- ✅ `show_usage()` - 1/1 teste passing

### Cache e Configuração
- ✅ `is_cache_valid()` - 3/4 testes passing
- ✅ `load_config()` - 2/2 testes passing
- ✅ `init_config()` - 1/2 teste passing

## Prioridades de Melhoria

### 🔥 CRÍTICO - Refatoração Estrutural

#### [ ] Separar funções de main()
**Problema:** `main()` é chamado automaticamente ao fazer source
**Solução:**
```bash
# Criar ghtools_core.sh com apenas funções
# Manter ghtools como wrapper que chama main()
# Tests sourceiam ghtools_core.sh
```

#### [ ] Adicionar modo de teste (`--test-mode`)
**Problema:** Muitas funções são interativas
**Solução:**
```bash
# Adicionar flag global para desabilitar interatividade
ghtools --test-mode list
ghtools --test-mode clone --path /tmp/test
```

#### [ ] Tornar funções pure (sem efeitos colaterais)
**Problema:** Funções dependem de estado global
**Solução:**
```bash
# Passar parâmetros explicitamente
# Retornar valores em vez de imprimir
# Usar variáveis locais
```

### 🟡 ALTA - Melhorar Mocks

#### [ ] Mock sofisticado do `gh`
**Atual:** Mock básico que sempre retorna sucesso
**Necessário:**
```bash
# Suportar diferentes comandos
# Retornar dados realistas
# Simular erros (401, 403, 404, 500)
# Suportar flags --json
```

#### [ ] Mock do `fzf` com cenários
**Atual:** Sempre retorna primeira linha
**Necessário:**
```bash
# Modo multi-select
# Cancelamento (ESC)
# Busca fuzzy real
# Preview
```

#### [ ] Mock do `gum` completo
**Atual:** Mock básico
**Necessário:**
```bash
# gum choose com seleção customizada
# gum input com defaults
# gum confirm com diferentes respostas
# gum style com cores
```

### 🟡 ALTA - Adicionar Testes

#### [ ] Testes para funções de printing (8 funções)
- [ ] `print_error()`
- [ ] `print_success()`
- [ ] `print_info()`
- [ ] `print_warning()`
- [ ] `print_verbose()`
- [ ] `show_header()`
- [ ] `show_divider()`
- [ ] `run_with_spinner()`

#### [ ] Testes para funções Gum (5 funções)
- [ ] `gum_confirm()`
- [ ] `gum_input()`
- [ ] `gum_choose()`
- [ ] `gum_filter()`
- [ ] `gum_write()`

#### [ ] Testes para parsing de argumentos
- [ ] `main()` com diferentes flags
- [ ] Parsing de --help, --version, --verbose, --quiet
- [ ] Validação de argumentos
- [ ] Combinação de flags

### 🟢 MÉDIA - Infraestrutura

#### [ ] CI/CD Integration
**GitHub Actions:**
```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Run tests
        run: ./run_tests.sh
```

#### [ ] Coverage Real
**Usar shcov:**
```bash
# Instalar shcov
# Gerar relatório de cobertura HTML
# Integrar com GitHub Pages
```

#### [ ] Test Fixtures
**Criar dados de teste:**
```bash
test/fixtures/
├── repos.json          # Mock de repositórios
├── pr_data.json        # Mock de PRs
└── user_data.json      # Mock de dados do usuário
```

### 🟢 MÉDIA - Documentação

#### [ ] Adicionar exemplos nos testes
```bash
# Documentar cada teste com:
# Descrição do que está sendo testado
# Cenário de teste
# Resultado esperado
# Links para documentação
```

#### [ ] Guia de contribuição
```bash
# Como escrever novos testes
# Como executar testes localmente
# Como adicionar mocks
# Padrões e convenções
```

## Plano de Ação (Próximas 4 Semanas)

### Semana 1: Refatoração
- [ ] Dia 1-2: Extrair funções para ghtools_core.sh
- [ ] Dia 3-4: Adicionar --test-mode flag
- [ ] Dia 5-7: Testar refatoração

### Semana 2: Mocks Sophisticados
- [ ] Dia 1-3: Mock completo do gh
- [ ] Dia 4-5: Mock completo do fzf
- [ ] Dia 6-7: Mock completo do gum

### Semana 3: Mais Testes
- [ ] Dia 1-2: Testes para printing functions
- [ ] Dia 3-4: Testes para Gum functions
- [ ] Dia 5-7: Testes para parsing

### Semana 4: Infraestrutura
- [ ] Dia 1-3: Setup CI/CD
- [ ] Dia 4-5: Coverage real com shcov
- [ ] Dia 6-7: Documentação final

## Métricas de Sucesso

| Métrica | Atual | Meta | Ações |
|---------|-------|------|-------|
| Testes Passing | 14 | 100+ | Escrever mais testes |
| Cobertura | 38% | 80%+ | Refatorar + mais testes |
| Funções Testadas | 17/45 | 36/45 | Cobrir funções restantes |
| Mock Quality | Básico | Avançado | Melhorar mocks |

## Comandos Úteis

```bash
# Ver estatísticas dos testes
./run_tests.sh

# Executar apenas testes passing
bats test/unit/test_utility_functions.bats

# Executar com verbose
bats --verbose test/unit/test_cache_and_config.bats

# Ver coverage por arquivo
bats --coverage test/unit/*.bats

# Gerar relatório HTML
bats --html-report report.html test/

# Executar teste específico
bats test/unit/test_utility_functions.bats -d "truncate_text returns original text"

# Ver apenas falhas
bats test/... 2>&1 | grep "not ok"

# Contar passing vs failing
bats test/... 2>&1 | grep -E "^(ok|not ok)" | wc -l
```

## Links Úteis

- [Bats Documentation](https://bats-core.readthedocs.io/)
- [Bats GitHub](https://github.com/bats-core/bats-core)
- [Shell Script Best Practices](https://google.github.io/styleguide/shellguide.html)
- [Advanced Bash Scripting Guide](https://tldp.org/LDP/abs/html/)

---

**Atualizado em:** 2025-12-05
**Responsável:** Equipe de Desenvolvimento
**Próxima Revisão:** 2025-12-12
