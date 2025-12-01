╭────────────────────────────────────────────────────────────────────────────╮
│                                                                            │
│    Refactoring/Design Plan: Modernização da Interface TUI do  ghtools      │
│   (Implementação com  gum )                                                │
│                                                                            │
│   ## 1. Executive Summary & Goals                                          │
│                                                                            │
│   O objetivo primário deste plano é modernizar a interface de linha de     │
│   comando (TUI) do script  ghtools , adotando os padrões atuais de         │
│   Terminal User Interface (TUI) e melhorando a Usabilidade e               │
│   Experiência do Usuário (UX). O script já utiliza  gum  e  fzf , mas      │
│   o plano visa maximizar o uso do  gum  para criar um look-and-feel        │
│   mais                                                                     │
│   coeso e moderno.                                                         │
│                                                                            │
│   * Meta 1: Substituir todas as interações e saídas de texto não           │
│   estilizadas ou com estilos de terminal ( echo -e "${COLOR}..." ) por     │
│   componentes  gum  (e.g.,  gum choose ,  gum input ,  gum style ).        │
│   * Meta 2: Centralizar e aprimorar a lógica de "Fallback" ( use_gum       │
│   function) para garantir a compatibilidade em terminais sem  gum ,        │
│   utilizando  fzf  e cores ANSI Legacy de forma mais consistente.          │
│   * Meta 3: Reestruturar as ações que ainda utilizam  fzf  diretamente     │
│   ( action_list ,  action_sync ,  action_status ) para usar as funções     │
│   de abstração  gum_choose / gum_filter  aprimoradas.                      │
│                                                                            │
│   --------                                                                 │
│                                                                            │
│   ## 2. Current Situation Analysis                                         │
│                                                                            │
│   * Visão Geral: O script  ghtools  é um monolito em Bash com 10.0KB,      │
│   que utiliza  gh ,  fzf ,  jq , e introduziu recentemente o  gum          │
│   (conforme evidenciado pelas funções  use_gum ,  gum_style ,              │
│   show_header , etc.).                                                     │
│   * Pontos Fortes (TUI):                                                   │
│     * Existe uma camada de abstração para cores e utilities ( print_*      │
│     e  gum_*  helpers).                                                    │
│     * O script já detecta a presença do  gum  ( use_gum ).                 │
│     * Algumas ações ( action_create ,  show_menu ) já usam  gum            │
│     interativamente.                                                       │
│   * Pontos Fracos (TUI):                                                   │
│     * A lógica de  action_list  ainda usa a combinação  jq -r  +           │
│     printf  + cores ANSI Legacy, o que é mais complexo e menos             │
│     "bonito" que uma tabela  gum .                                         │
│     * Ações cruciais como  action_sync ,  action_delete , e                │
│     action_search  ainda dependem do  fzf  puro em sua lógica central      │
│     de seleção, com cabeçalhos não estilizados em modo fallback.           │
│     * A saída de status ( action_status ) e o  show_usage  (help)          │
│     ainda são predominantemente texto ASCII/ANSI Legacy, não               │
│     aproveitando o poder de estilização do  gum .                          │
│                                                                            │
│                                                                            │
│   --------                                                                 │
│                                                                            │
│   ## 3. Proposed Solution / Refactoring Strategy                           │
│                                                                            │
│   O refactoring se concentrará em consolidar todas as interações e         │
│   saídas de dados na camada de abstração  gum_*  e garantir que o          │
│   fallback para  fzf /ANSI seja robusto. O uso da tabela  gum  será        │
│   priorizado para listagens de dados.                                      │
│                                                                            │
│   ### 3.1. High-Level Design / Architectural Overview                      │
│                                                                            │
│   A estratégia é promover o  gum  de uma ferramenta opcional para o        │
│   componente primário de UI/UX, mantendo  fzf  como a principal            │
│   ferramenta de fallback para busca/seleção fuzzy, e as cores ANSI         │
│   Legacy como fallback final para simples mensagens.                       │
│                                                                            │
│   Target Architecture (Conceptual):                                        │
│                                                                            │
│                                                                            │
│     ----------                                                             │
│     graph TD                                                               │
│         A[ghtools Main] --> B{use_gum?};                                   │
│         B -- Yes --> C[Gum Interface];                                     │
│         B -- No --> D[Fallback Interface];                                 │
│         C --> C1[gum choose/style/table];                                  │
│         D --> D1[fzf/ANSI Legacy];                                         │
│         C1 --> E(gh CLI + jq Logic);                                       │
│         D1 --> E;                                                          │
│         E --> F[Repository Management];                                    │
│     ----------                                                             │
│                                                                            │
│   ### 3.2. Key Components / Modules                                        │
│                                                                            │
│    Componente        | Responsabilidade At… | Modificação Proposta         │
│   -------------------+----------------------+-----------------------       │
│     print_*  Helpers | Mensagens de status  | Aprimorar: Garantir          │
│                      | (INFO, ERROR, etc.)  | que o  gum style  use        │
│                      |                      | largura e alinhamento        │
│                      |                      | corretos para coesão         │
│                      |                      | visual.                      │
│     gum_*  Utilities | Abstração para       | Criar  gum_table :           │
│                      | componentes          | Nova função para             │
│                      |  gum / fzf           | renderizar dados em          │
│                      |                      | formato tabular com          │
│                      |                      |  gum table  (ou              │
│                      |                      |  printf  estilizado          │
│                      |                      | em fallback).                │
│     action_list      | Listagem principal   | Refatorar: Usar              │
│                      | de repositórios      |  gum_table  para a           │
│                      |                      | saída de dados,              │
│                      |                      | substituindo a lógica        │
│                      |                      | complexa de                  │
│                      |                      |  printf /ANSI.               │
│     action_sync      | Seleção de           | Integrar: Usar               │
│                      | repositórios para    |  gum filter --no-            │
│                      | Sync                 | limit  ou a nova             │
│                      |                      | função  gum_filter           │
│                      |                      | aprimorada para a            │
│                      |                      | seleção de múltiplos.        │
│     show_usage       | Mensagem de ajuda    | Estilizar: Usar              │
│                      |                      |  gum style  para             │
│                      |                      | estilizar cabeçalhos         │
│                      |                      | e seções, tornando-a         │
│                      |                      | mais moderna.                │
│                                                                            │
│   ### 3.3. Detailed Action Plan / Phases                                   │
│                                                                            │
│   #### Phase 1: Consolidação da Camada de Abstração  gum  (S)              │
│                                                                            │
│   * Objetivo(s): Criar a função  gum_table  e revisar os helpers de        │
│   mensagens.                                                               │
│   * Prioridade: High                                                       │
│                                                                            │
│    Task               | Rationale/Goal      | | Deliverable/Criter…        │
│   --------------------+---------------------+-+---------------------       │
│    1.1: Revisar       | Garantir que as     | |  show_header               │
│     print_*  e        | mensagens de status | | utiliza  gum style         │
│     show_header       | ( print_error ,     | | de forma otimizada;        │
│                       |  print_success ) e  | | as mensagens               │
│                       | o  show_header      | |  print_*  usam um          │
│                       | sejam visualmente   | | formato padronizado        │
│                       | impactantes e       | | (e.g.,                     │
│                       | consistentes em     | |  [ICON] MESSAGE ).         │
│                       | ambas as condições  | |                            │
│                       | ( use_gum  e        | |                            │
│                       |  fallback ).        | |                            │
│    1.2: Criar         | Padronizar a saída  | | Nova função                │
│     gum_table         | de dados            | |  gum_table(title, h        │
│    Function           | estruturados (e.g., | | eaders, data, separ        │
│                       | List, Stats).       | | ator)  que usa             │
│                       |                     | |  gum table  se             │
│                       |                     | | disponível ou              │
│                       |                     | |  printf  estilizado        │
│                       |                     | | como fallback.             │
│    1.3: Aprimorar     | Garantir que        | | A função                   │
│     gum_filter  para  |  gum_filter  lide   | |  gum_filter  lida          │
│    Seleção Múltipla   | corretamente com    | | com o argumento  --        │
│                       | seleção múltipla,   | | multi  de forma            │
│                       | substituindo o uso  | | transparente para          │
│                       | direto do  fzf  nos | | as ações que o             │
│                       | comandos.           | | invocam.                   │
│                                                                            │
│   #### Phase 2: Refatoramento das Ações de Listagem e Status (L)           │
│                                                                            │
│   * Objetivo(s): Aplicar a  gum_table  nas ações de visualização de        │
│   dados.                                                                   │
│   * Prioridade: High                                                       │
│                                                                            │
│    Task              | Rationale/Goal     | E… | Deliverable/Crite…        │
│   -------------------+--------------------+----+--------------------       │
│    2.1: Refatorar    | Substituir a       | L  |  action_list              │
│     action_list      | lógica complexa de |    | exibe uma tabela          │
│                      |  jq + printf  por  |    | moderna com  gum          │
│                      |  gum_table . O     |    | e uma tabela ANSI         │
│                      | JSON output do     |    | limpa em fallback.        │
│                      |  jq  deve ser      |    |                           │
│                      | processado para um |    |                           │
│                      | formato de entrada |    |                           │
│                      |  gum table .       |    |                           │
│    2.2: Refatorar    | Converter a saída  | M  |  action_stats  usa        │
│     action_stats     | de estatísticas    |    |  gum style  para          │
│                      | para usar os       |    | todas as caixas de        │
│                      | blocos e títulos   |    | informação e              │
│                      | estilizados do     |    |  gum style                │
│                      |  gum style  e, se  |    | simples ou                │
│                      | viável, uma        |    |  gum table  para          │
│                      |  gum_table  para a |    | listas.                   │
│                      | parte de           |    |                           │
│                      | Linguagens/Top     |    |                           │
│                      | Repos.             |    |                           │
│    2.3: Refatorar    | Usar a  gum_table  | M  |  action_status            │
│     action_status    | ou um  printf      |    | apresenta o status        │
│                      | fortemente         |    | local ( dirty ,           │
│                      | estilizado (com    |    |  ahead ,  behind )        │
│                      | base na            |    | em formato tabular        │
│                      |  gum_table         |    | moderno e legível.        │
│                      | fallback) para o   |    |                           │
│                      | output de status   |    |                           │
│                      | dos repositórios   |    |                           │
│                                                                            │
│   #### Phase 3: Integração Completa de Interações  gum  (M)                │
│                                                                            │
│   * Objetivo(s): Eliminar o uso direto de  read -p  e  fzf  puro em        │
│   favor dos wrappers  gum_* .                                              │
│   * Prioridade: Medium                                                     │
│                                                                            │
│    Task               | Rationale/Goal      | | Deliverable/Criter…        │
│   --------------------+---------------------+-+---------------------       │
│    3.1: Refatorar     | Usar  gum_confirm   | | O fluxo de                 │
│     action_delete     | aprimorado para a   | |  action_delete  é          │
│                       | confirmação de      | | totalmente                 │
│                       | exclusão, e o       | | interativo via             │
│                       |  gum input  forçado | |  gum  (ou                  │
│                       | para a confirmação  | |  fzf / read -p             │
│                       | de texto "DELETE".  | | robustos no                │
│    3.2: Refatorar     | Substituir o        | |  check_delete_scope        │
│     check_delete_scop |  read -p  dentro do | |   usa o helper             │
│    e                  |  check_delete_scope | |  gum_confirm .             │
│                       |   por  gum_confirm  | |                            │
│                       | para manter a       | |                            │
│                       | consistência do UX. | |                            │
│    3.3: Refatorar     | Garantir que a      | |  action_sync  e            │
│     action_sync  e    | seleção de          | |  action_search             │
│     action_search     | repositórios use    | | usam  gum_filter  e        │
│                       |  gum_filter         | | as ações de                │
│                       | (seleção multi) de  | | acompanhamento usam        │
│                       | forma consistente.  | |  gum_choose .              │
│                                                                            │
│   #### Phase 4: Aprimoramento da Mensagem de Ajuda e UX (S)                │
│                                                                            │
│   * Objetivo(s): Modernizar a saída de ajuda e garantir a coesão final.    │
│   * Prioridade: Low                                                        │
│                                                                            │
│    Task              | Rationale/Goal     | E… | Deliverable/Crite…        │
│   -------------------+--------------------+----+--------------------       │
│    4.1: Estilizar    | Estilizar o output | S  |  show_usage  é            │
│     show_usage       | do  show_usage     |    | visualmente               │
│                      | (help) usando      |    | segmentado e usa          │
│                      |  gum style  para   |    | um esquema de             │
│                      | cabeçalhos e       |    | cores moderno.            │
│                      | blocos de texto,   |    |                           │
│                      | melhorando a       |    |                           │
│                      | apresentação no    |    |                           │
│                      | terminal.          |    |                           │
│    4.2: Atualizar    | Mencionar e        | S  | Documentação              │
│     README.md        | promover o  gum    |    | reflete o novo            │
│                      | como a dependência |    | design TUI.               │
│                      | recomendada para a |    |                           │
│                      | TUI moderna.       |    |                           │
│                                                                            │
│   --------                                                                 │
│                                                                            │
│   ## 4. Key Considerations & Risk Mitigation                               │
│                                                                            │
│   ### 4.1. Technical Risks & Challenges                                    │
│                                                                            │
│   * Risco: Falha de Instalação do  gum : O  install.sh  já alerta          │
│   sobre a ausência do  gum , mas a ausência pode prejudicar                │
│   significativamente o UX prometido.                                       │
│     * Mitigação: Fortalecer a lógica de fallback. Garantir que o  fzf      │
│     e os estilos ANSI Legacy sejam a segunda melhor experiência, não       │
│     uma experiência quebrada.                                              │
│   * Risco: Quebra de Compatibilidade Bash/POSIX: O script é Bash, mas      │
│   o uso intensivo de strings complexas com  gum style  e  jq  pode         │
│   levar a problemas em ambientes não-Bash (embora improvável, dado o       │
│   escopo).                                                                 │
│     * Mitigação: Testar em diferentes shells (zsh/bash/fish - este         │
│     último apenas para validação, pois o script é Bash). Manter a          │
│     seção de lógica de UI em  ghtools  separada da lógica de negócios.     │
│                                                                            │
│                                                                            │
│   ### 4.2. Dependencies                                                    │
│                                                                            │
│   * Interna:                                                               │
│     * Phase 2 depende da conclusão das novas funções de abstração TUI      │
│     da Phase 1.                                                            │
│     * A refatoração de cada ação é isolada, minimizando o risco de         │
│     quebra de outras ações.                                                │
│   * Externa:                                                               │
│     *  gum  (altamente recomendado para UX moderna).                       │
│     *  fzf  (essencial para o fallback de busca fuzzy).                    │
│                                                                            │
│                                                                            │
│   ### 4.3. Non-Functional Requirements (NFRs) Addressed                    │
│                                                                            │
│   * Usabilidade (NFR Principal): Melhorada significativamente pelo uso     │
│   consistente de componentes TUI modernos ( gum choose ,  gum input ,      │
│   gum table ), que são mais intuitivos e fáceis de usar do que prompts     │
│   de texto e comandos de confirmação.                                      │
│   * Aparência (Aesthetic): O uso do esquema de cores  gum  Purple/Cyan     │
│   e os componentes estilizados satisfazem o requisito de tornar a          │
│   interface "mais moderna e bonita".                                       │
│   * Manutenibilidade: A criação da função  gum_table  centraliza a         │
│   lógica de visualização tabular, removendo a lógica complexa de           │
│   printf /ANSI das funções de ação, melhorando a coesão e a facilidade     │
│   de manutenção.                                                           │
│   * Confiabilidade: O robusto sistema de fallback ( use_gum  com  fzf      │
│   /ANSI) garante que a ferramenta permaneça funcional em ambientes         │
│   onde a dependência principal (gum) possa estar ausente.                  │
│                                                                            │
│   --------                                                                 │
│                                                                            │
│   ## 5. Success Metrics / Validation Criteria                              │
│                                                                            │
│   * O comando  ghtools  (sem argumentos) exibe o menu principal usando     │
│   gum choose  com cores e cabeçalho modernos (validado visualmente).       │
│   * A saída de  ghtools list  exibe uma tabela estruturada, colorida e     │
│   formatada corretamente, tanto com  gum  quanto em fallback (usando a     │
│   nova lógica de  gum_table ).                                             │
│   * Todas as interações de prompt, confirmação e escolha de menu (em       │
│   ações como  clone ,  delete ,  create ) utilizam os componentes  gum     │
│   (e.g.,  gum_confirm ,  gum_input ).                                      │
│   * O shellcheck passa limpo após todas as alterações.                     │
│                                                                            │
│   --------                                                                 │
│                                                                            │
│   ## 6. Assumptions Made                                                   │
│                                                                            │
│   * Assume-se que o pacote  gum  está disponível para instalação na        │
│   maioria dos ambientes alvo (Arch Linux, ou através do  go install ),     │
│   conforme implícito na estrutura do  install.sh .                         │
│   * A refatoração da lógica de visualização de dados pode exigir uma       │
│   mudança no formato de saída do  jq  para alimentar as funções de         │
│   tabela, mas isso é considerado um detalhe de implementação dentro        │
│   das tarefas.                                                             │
│   * O Bash shell subjacente tem suporte a cores (o que é padrão).          │
│                                                                            │
│   --------                                                                 │
│                                                                            │
│   ## 7. Open Questions / Areas for Further Investigation                   │
│                                                                            │
│   * Design de  gum_table : Qual é o formato exato de entrada/saída         │
│   mais ergonômico para a nova função  gum_table ? (Ex: Deve aceitar        │
│   JSON e formatar, ou aceitar uma lista de colunas delimitadas por         │
│   TSV?).                                                                   │
│   * Integração do  gum style  com  show_usage : Devemos usar um            │
│   HEREDOC grande com tags de estilo  gum  para o  show_usage  ou           │
│   processar o texto usando um pipeline de  sed / gum style ?               │
│   * Configuração de Cores: O esquema de cores (Purple/Cyan) deve ser       │
│   configurável via  config  file? (Recomendado para futuras melhorias,     │
│   mas fora do escopo do plano atual).                                      │
╰────────────────────────────────────────────────────────────────────────────╯