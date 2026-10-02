# 001-04 — Validar aquisição do meta

- **Ticker:** `001`
- **Número:** `04`
- **Status:** `completed with constraints`

## Requisitos cobertos

- fonte de candidatos de decks reais e contexto competitivo;
- cobertura de faixas intermediárias, viés, freshness, custo e compliance;
- fallback de produto quando “meta da Arena” não puder ser sustentado.

## Objetivo e resultado esperado

Provar que o CrownPilot consegue obter candidatos de deck com contexto
competitivo suficiente para recomendar **Best Decks for You**, especialmente
para jogadores que não estão no topo da ladder.

Esta é a principal prova de viabilidade do produto.

## Escopo incluído

Avaliar no mínimo estas estratégias:

1. rankings globais/regionais;
2. Path of Legend/ranked quando aplicável;
3. clan search + clan members;
4. expansão controlada via opponents do battle log;
5. fontes third-party somente quando houver API/licença/termos claros;
6. combinação híbrida das anteriores.

Para cada estratégia medir:

- cobertura por faixa de trophies/Arena;
- viés;
- freshness;
- request cost;
- facilidade de deduplicação;
- estabilidade;
- dependência externa;
- compliance.

## Escopo excluído

- crawler de produção;
- scraping de sites;
- coleta massiva;
- benchmark de milhões de batalhas;
- algoritmo final de meta strength.

## Dependências

- 001-02;
- 001-03.

## Arquivos e símbolos prováveis

- `evidences/meta-strategy-comparison.md`;
- superfícies de rankings, Path of Legend/ranked, clans/members e opponents;
- conceitos `DeckCandidate`, `CompetitiveContextV0` e estratégia de ingestão;
  nenhum pipeline ou dataset existe na `main`.

## Passos de execução

1. Definir quais contextos precisam existir para o MVP:
   - Arena/faixa intermediária;
   - ladder avançada;
   - ranked quando relevante.
2. Testar rankings oficiais e registrar cobertura real.
3. Testar discovery por clans/members com amostra pequena.
4. Verificar quantos novos players/opponents surgem de battle logs.
5. Projetar custo de expansão controlada sem executá-la em escala.
6. Identificar third-party APIs candidatas apenas se necessário.
7. Comparar as estratégias em matriz.
8. Escolher:
   - estratégia preferida;
   - fallback;
   - contextos que não podem ser suportados honestamente.

## Evidência obrigatória

`evidences/meta-strategy-comparison.md` com:

| Strategy | Coverage | Mid-ladder | Bias | Requests | Freshness | Terms | Dependency | Verdict |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |

## Gate de produto

Se nenhuma estratégia conseguir gerar contexto razoável para faixas
intermediárias, a task **não pode** marcar "meta por Arena" como resolvido.

Nesse caso o veredito deve recomendar uma destas saídas:

- reduzir escopo para "best decks for your collection" sem alegar meta da Arena;
- segmentar por outra dimensão observável;
- usar fonte licenciada;
- redesenhar aquisição de dados.

## Definição de pronto

Existe pelo menos uma estratégia defendível ou uma limitação de produto
explicitamente aceita para a Fase 005.

## Testes e comandos de validação

- executar somente amostras pequenas e documentadas por estratégia;
- preencher a matriz comparativa com cobertura observada, não cobertura inferida
  de ranking global;
- calcular requests projetados por estratégia sem coleta massiva;
- revisar termos/licença de qualquer fonte externa e executar `git diff --check`.

## Riscos e cuidados

- top rankings têm forte selection bias;
- clan sampling também pode enviesar atividade/geografia;
- graph expansion por opponents pode crescer exponencialmente;
- third-party grátis não significa permitido/reutilizável;
- não transformar probe em crawler.

## Registro de execução

### Execução em `2026-10-02`

- **Status final:** `completed with constraints`. A estratégia de aquisição de
  candidatos foi defendida, mas o gate de “meta por Arena” não foi resolvido.
- **Estratégia preferida:** híbrida e limitada: perfil + battle log do jogador,
  com enriquecimento opcional de pequenos cohorts de clan. Particionar `PvP`,
  `pathOfLegend`, `trail` e `unknown`; não misturar modos.
- **Fallback:** reduzir promessa para **Best Decks for Your Collection**,
  contextualizado por trophies/Arena/mode e encontros recentes. Alternativas
  futuras: segmentação observável, fonte licenciada ou redesenho de aquisição.
  Nenhuma permite declarar “meta da Arena” nesta fase.
- **Cobertura não resolvida:** rankings de players global/regional retornaram
  `404` no proxy; Path of Legend leaderboard também retornou `404`. Clan
  rankings retornaram 999 rows, mas os 11 perfis de membros amostrados tinham
  14.000 trophies, mostrando concentração de topo consistente com viés de
  seleção. Battle logs produziram 6/30/30 oponentes únicos; o sample PvP
  intermediário cobriu 2.960–3.110 trophies em apenas seis entradas, sem
  prevalência representativa de Arena.
- **Estimativa de requests:** perfil + battle log inicial = 2/player; cohort
  de clan = `3 + K` requests, aproximadamente 53 para 50 perfis; expansão de
  opponents = `1 + U`, aproximadamente 31 com `U=30`; dois saltos podem chegar
  a 931 antes de deduplicação. Polling a cada 5m projeta 2.016/4.032/8.640
  chamadas por player em 7/14/30 dias, antes de retries e outages.
- **Arquivos alterados:**
  `evidences/meta-strategy-comparison.md`; este registro; overview para marcar
  somente `001-04` e recalcular progresso.
- **Decisões:** não tratar ranking de topo como meta de Arena; não usar
  RoyaleAPI Proxy como fonte de meta ou autorização; manter third-party sem API,
  licença e termos claros fora do MVP; limitar expansão via opponents; aceitar
  apenas recomendação contextual/personalizada.
- **Desvios:** probes autenticados usaram RoyaleAPI Proxy porque egress direto
  permanece bloqueado por allowlist. A documentação pública do endpoint oficial
  continua sem schema autenticado; observações do proxy não foram promovidas a
  contrato oficial. Nenhum crawler, coleta em escala ou expansão de dois saltos
  foi executado.
- **Comandos executados:** probes Python/urllib server-side via
  `proxy.royaleapi.dev` para locations, rankings, Path of Legend, clan search,
  clan members, profiles e battle logs; `webfetch` das fontes oficiais e docs do
  proxy; validação estrutural Python de ticker/checklist/matriz; `git diff --check`.
- **Resultados/evidências:** `evidences/meta-strategy-comparison.md` registra
  fonte/data, amostras sanitizadas, matriz obrigatória, viés, freshness, custo,
  deduplicação por estratégia, estabilidade, dependências, compliance, fallback
  e gate. Locations retornou 262 rows;
  rankings de clans 999 rows com paging; clan search 640 rows; members 50 rows;
  cache observado entre 6s e 225s conforme endpoint; todos os 11 perfis de
  membros amostrados retornaram `200`, com 14.000 trophies.
- **Testes e validação:** `git diff --check` passou. Não existem lint,
  typecheck, build ou testes de aplicação neste baseline documental; probe
  controlado, revisão de sanitização e validação estrutural foram a evidência
  pertinente.
- **Riscos residuais:** retenção e uso permitido de dados de opponents/members;
  key handling, limites, SLA e compatibilidade do proxy; schema oficial e
  disponibilidade de rankings; amostra pequena sem prova de prevalência ou SLA
  de freshness; mudanças futuras de endpoint/shape.
- **Revisão independente:** primeira revisão encontrou correções necessárias em
  request cost, cobertura por faixa e deduplicação/estabilidade; correções foram
  aplicadas. Follow-up independente aprovou a subtarefa sem blockers; ressalva
  sobre checkboxes preexistentes na spec/001-08 foi classificada fora do escopo.
  Nenhuma subtarefa seguinte foi iniciada.
