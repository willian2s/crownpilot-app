# 001-05 — Medir operação, cache e rate limits

- **Ticker:** `001`
- **Número:** `05`
- **Status:** `completed with constraints`

## Requisitos cobertos

- cache, latência, erros, `Retry-After`, rate behavior e requisitos de token/IP;
- custo de requests para sync de jogador, histórico, catálogo e meta;
- impacto operacional em 1k, 10k e 100k usuários sem executar carga real.

## Objetivo e resultado esperado

Transformar a estratégia de dados em ordem de grandeza operacional antes de
escolher infraestrutura.

## Escopo incluído

- registrar headers de cache relevantes;
- observar latência básica;
- confirmar comportamento de 403/404/429/5xx quando ocorrer de forma natural ou
  com probes seguros;
- honrar `Retry-After`;
- verificar restrições de token/IP;
- validar compatibilidade do modelo de egress com Vercel;
- registrar se Static IPs da Vercel são necessários/adequados ou se um egress
  gateway substituível será necessário;
- estimar requests de:
  - player sync;
  - battle history;
  - card catalog;
  - meta ingestion;
- projetar cenários de 1k, 10k e 100k usuários ativos sem executar carga real;
- produzir `evidences/operational-findings.md`.

## Escopo excluído

- load test;
- tentativa de descobrir limite por saturação;
- bypass de quota;
- rotação de IP/token para aumentar throughput;
- troca da decisão de Vercel como deploy inicial;
- implementação definitiva de gateway/worker.

## Dependências

- 001-01 para auth, endpoints e headers básicos;
- 001-02 e 001-03 para custo de perfil e histórico;
- 001-04 para projetar custo de aquisição do meta.

## Arquivos e símbolos prováveis

- `evidences/operational-findings.md`;
- headers `Cache-Control` e `Retry-After`, códigos HTTP e variável
  `CLASH_ROYALE_API_TOKEN`;
- operações conceituais `syncPlayer`, `ingestMetaBatch` e `refreshCatalog`;
  nenhum job ou adapter existe na `main`.

## Passos de execução

1. Coletar headers de endpoints já usados nas tasks anteriores.
2. Registrar `Cache-Control`, `Retry-After` quando existir e códigos.
3. Separar requests user-driven de background ingestion.
4. Definir fórmula simples de requests por sync.
5. Projetar três cenários de escala.
6. Identificar quais resultados podem ser compartilhados/cacheados globalmente.
7. Verificar o impacto de IP allowlist no deploy inicial em Vercel.
8. Se egress estático for necessário, comparar:
   - Vercel Static IPs;
   - egress gateway externo/substituível.
9. Registrar requisitos que futura infraestrutura precisa satisfazer sem
   acoplar o domínio ao provider.

## Princípio de segurança operacional

Não provocar 429 deliberadamente para "medir o limite".

Se o limite oficial não for publicado, registrar como desconhecido e trabalhar
com observação conservadora.

## Saída mínima

`evidences/operational-findings.md` deve mostrar:

- request path;
- cache observado/documentado;
- custo por operação;
- cadence candidata;
- failure behavior;
- impacto em 1k/10k/100k;
- mitigação.

## Definição de pronto

A Fase 002 conhece as restrições de integração sem estar presa a um deploy
específico.

## Testes e comandos de validação

- reutilizar respostas das tasks anteriores, sem provocar `429` deliberadamente;
- registrar headers e status de probes seguros em evidência sanitizada;
- calcular requests por operação e projetar os três cenários de escala,
  marcando-os como estimativas;
- revisar `Retry-After`, separação user-driven/background e `git diff --check`.

## Riscos e cuidados

- IP-bound token pode exigir Vercel Static IPs ou egress gateway dedicado;
- dependência direta de networking proprietário da Vercel criaria lock-in;
- cache do provedor pode mudar;
- custo de meta ingestion pode dominar custo de user sync;
- números projetados devem ser marcados como estimativas, não SLOs.

## Registro de execução

### Execução em `2026-10-02`

- **Status final:** `completed with constraints`.
- **Arquivos alterados:** `evidences/operational-findings.md`; este registro;
  overview para marcar somente `001-05` e recalcular progresso.
- **Cache observado:** catálogo `max-age=7–49s`, perfil `36–60s`, battle log
  `21–60s`, locations `138–404s` e erro de tag `public max-age=600s`, todos via
  probes anteriores pelo proxy. Valores são observações, não contratos de
  freshness.
- **Latência observada:** probes autenticados seguros via proxy retornaram
  `/cards=200` em `0,522s`, `/locations=200` em `0,470s` e tag inválida `404` em
  `0,447s`; rota oficial autenticada retornou `403 invalidIp` em `0,745s`. Uma
  amostra por rota não sustenta SLO.
- **Rate behavior:** portal documenta rate limitation qualitativa; limite
  numérico permanece desconhecido. Não ocorreu `429`/`5xx` naturalmente, nenhum
  probe de saturação foi executado, `Retry-After` não apareceu nos headers
  observados e a regra proposta é honrar o header quando presente.
- **Request model:** sync inicial `2/player`; histórico a cada 5m `288/player/dia`;
  refresh de perfil diário `1/player/dia`; catálogo global a cada 6h `4/dia`;
  locations global diário `1/dia`; batch de meta limitado `306/dia` no exemplo
  (100 seeds + 2 cohorts de 50). Cenário de impacto `289P+311/dia`: 289.311,
  2.890.311 e 28.900.311 requests/dia para 1k/10k/100k, sem retries/cache-hit
  descontado.
- **Infra constraints:** egress direto da chave é IP-bound e falhou na Vercel/local
  sem allowlist. Vercel Static IPs Pro/Enterprise pode atender outbound allowlist;
  uma ou duas regiões cabem nas cinco entradas CIDR observadas na chave, mas três
  regiões podem exigir seis IPs e exceder essa capacidade. Há pool compartilhado e
  custo documentado de US$100/projeto/mês mais transferência; gateway externo
  continua fallback substituível. Proxy atual é transporte temporário condicional,
  não autorização de API, dados ou billing.
- **Decisões:** catalog/locations podem ter cache global; profile/battle log ficam
  particionados por tag e boundary privado; tracking de battle log a cada 5m fica
  restrito a jogadores ativos/opt-in; meta permanece bounded e personalizada;
  domínio usa adapter independente do provider.
- **Desvios:** não houve probe específico de 429/5xx nem carga real, conforme
  princípio de segurança; latência foi medida em probes pequenos sem persistir
  token, IP, tag, body ou headers brutos.
- **Comandos e fontes:** reutilização das evidências `001-01`–`001-04`; probes
  `curl` autenticados server-side somente para status/tempo de `/cards`,
  `/locations`, tag inválida e rota oficial; `webfetch` da documentação Vercel
  Static IPs, fixed IP guidance e Function regions; `git diff --check`.
- **Resultados/evidências:** `evidences/operational-findings.md` registra paths,
  cache, latência, failure behavior, formulas, cadences, cenários 1k/10k/100k,
  comparação Static IP/gateway e handoff portátil. Não existem lint, typecheck,
  build ou testes de aplicação neste baseline documental.
- **Riscos residuais:** limite numérico e semântica de `Retry-After`; headers e
  latência variáveis do proxy; key handling/retenção/SLA do proxy; custo de Static
  IPs/gateway; ausência de p95/p99 e concorrência; privacy/terms para dados de
  oponentes e clans; projections não são SLOs.
- **Revisão independente:** primeira revisão aprovou evidência, sanitização,
  fórmulas diárias e escopo, mas solicitou qualificar compatibilidade entre até
  três regiões Static IPs e as cinco entradas CIDR da chave; correção aplicada.
  Checklists extras preexistentes na spec/001-08 foram apontados como fora do
  escopo desta subtarefa; overview continua com uma seção e oito itens. Follow-up
  independente aprovou `001-05` sem blockers. Nenhuma subtarefa seguinte foi
  iniciada.
