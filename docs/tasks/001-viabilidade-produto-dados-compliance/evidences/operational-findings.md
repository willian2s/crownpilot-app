# Operational findings — cache, rate behavior e escala

- **Task:** `001-05`
- **Observed at:** `2026-10-02`
- **Status:** `completed with constraints`
- **Authority:** [Clash Royale API developer portal](https://developer.clashroyale.com/)
- **Operational route:** RoyaleAPI Proxy, server-side, substituindo somente o host
  da API. A rota oficial direta continua rejeitando o egress atual por allowlist.
- **Evidence boundary:** números de latência, cache e status são observações de
  probes controlados e das tasks `001-01` a `001-04`; projeções são estimativas,
  não limites, SLOs ou load tests.

## Resumo executivo

1. O MVP deve tratar requests de perfil e battle history como operações
   user-driven ou de tracking explicitamente opt-in. Polling de battle log a cada
   5 minutos é uma cadência candidata para jogadores ativos, não uma promessa de
   polling para todos os usuários.
2. Catálogo e locations podem ser compartilhados globalmente. Perfil e battle log
   exigem chave de cache por Player Tag e fronteira privada; não compartilhar
   associação de usuário, nomes ou payload bruto.
3. O portal oficial documenta que tokens têm limitações de taxa, mas o número não
   foi publicado na sessão disponível. Nenhum `429` foi provocado; não há base para
   descobrir limite por saturação.
4. O egress oficial rejeitou a chave por IP não permitido. Vercel sem produto de
   egress fixo não atende allowlist; Static IPs pode atender em Pro/Enterprise, mas
   número de regiões/IPs precisa caber na allowlist da chave. Tem custo, pool
   compartilhado e escopo por projeto. Gateway externo permanece alternativa
   substituível.
5. Polling de 5 minutos para 100k jogadores produziria aproximadamente 28,8M
   requests de battle log/dia na borda da aplicação. Isso é cenário de impacto,
   não recomendação de arquitetura nem aprovação de escala.

## Request paths e observações de transporte

| Operação | Path observado | Tipo | Cache observado | Latência observada* | Custo de planejamento |
| --- | --- | --- | --- | ---: | ---: |
| `syncPlayer` | `GET /v1/players/{tag}` | user-driven; refresh controlado | `max-age=36s` via proxy em `001-02`; `60s` em `001-01` | não repetido nesta sessão; use apenas como ordem de grandeza de rede | 1 request por sync |
| `battleHistory` | `GET /v1/players/{tag}/battlelog` | background opt-in ou user-driven | `max-age=21–60s` em `001-03` | não repetido nesta sessão; observação anterior via proxy | 1 request por polling |
| `refreshCatalog` | `GET /v1/cards` | background global | `max-age=7s` em `001-02`; `49s` em `001-01` | `0,522s` total, probe autenticado via proxy | 1 request por refresh global |
| `refreshLocations` | `GET /v1/locations` | background global; útil para discovery | `max-age=138–225s` em `001-04`; `404s` em `001-01` | `0,470s` total, probe autenticado via proxy | 1 request por refresh global |
| erro de tag | `GET /v1/players/%23INVALIDTAG` | probe seguro/error path | `public max-age=600s` em `001-01` | `0,447s` total, `404` via proxy | não repetir como health check frequente |

\* Latência é uma amostra única de `curl` em 2026-10-02, incluindo conexão,
TLS, proxy e resposta. Não representa p95/p99, disponibilidade ou SLO. Os probes
de latência não salvaram body/header em arquivo versionado.

### Rota direta e proxy

| Rota | Status observado | Interpretação operacional |
| --- | --- | --- |
| API oficial, sem Authorization | `403` | auth gate alcançável; `Cache-Control: public max-age=600` observado em `001-01` |
| API oficial, com chave atual | `403`, `accessDenied.invalidIp`; `0,745s` nesta sessão | token chegou à API, mas egress não está allowlisted; não classificar como token inválido |
| Proxy, com chave atual, `/cards` | `200`; `0,522s` nesta sessão | transporte operacional funciona para probe controlado |
| Proxy, com chave atual, `/locations` | `200`; `0,470s` nesta sessão | catálogo global auxiliar funciona |
| Proxy, tag inválida | `404`; `0,447s` nesta sessão | erro de recurso observado; não inferir contrato direto oficial |

O proxy não prova origem oficial, endorsement, ownership ou autorização comercial.
Key handling, retenção, limites, SLA, segurança e compatibilidade contratual do
proxy permanecem `UNRESOLVED`.

## Cache e compartilhamento

`Cache-Control` observado varia entre chamadas e rotas. Valores são hints do
transporte/proxy, não freshness contract. O adapter deve separar dois números:

- **application request:** chamada feita pelo runtime do CrownPilot;
- **upstream/origin work:** trabalho que proxy/API efetivamente executa após seu
  próprio cache.

Sem métrica de hit/miss do proxy, não reduzir projeções de quota usando cache.
Cache compartilhado pode reduzir origin work, mas não autoriza assumir que uma
chamada do adapter não conta para limite ou custo.

| Resultado | Compartilhamento recomendado | Mitigação |
| --- | --- | --- |
| cards/catalog capability | global, por endpoint e versão de origem | refresh background; deduplicar requests concorrentes; TTL de produto separado do `max-age` observado |
| locations | global | refresh lento e sob demanda; não usar como health check frequente |
| perfil | por Player Tag, atrás de boundary server-side | não retornar cache de uma associação para outro usuário sem autorização; preservar `fetchedAt` e estado stale |
| battle log | por Player Tag e modo/uso interno | armazenar mínimo necessário; deduplicar por fingerprint; não compartilhar dados de oponentes por padrão |
| meta derivada | global somente após gate de retenção/uso | sem meta representativa de Arena nesta fase; manter estratégia híbrida limitada |

Dados derivados de um perfil podem ser reaproveitados tecnicamente entre usuários
que informem a mesma tag, mas a associação CrownPilot, acesso e retenção devem
continuar privados até revisão de privacy/terms.

## Failure behavior e rate behavior

### Evidência observada

| Condição | Resultado | Estado da evidência |
| --- | --- | --- |
| auth ausente | `403 accessDenied` | observado na rota oficial |
| egress/IP não permitido | `403 accessDenied.invalidIp` | observado na rota oficial com chave |
| Player Tag inválida | `404` com objeto de erro | observado via proxy |
| `429` | nenhum ocorrido naturalmente | desconhecido; probe deliberado proibido |
| `5xx`/manutenção | nenhum ocorrido nos probes | desconhecido |
| `Retry-After` | não apareceu nos headers observados | desconhecido; deve ser honrado se vier |

Regra futura do adapter: em `429`, aguardar exatamente `Retry-After` quando
presente; sem header, usar backoff exponencial com jitter e teto, sem retry
imediato. Em `5xx`/timeout, retry limitado e idempotente; em `403` por IP/token
ou `404` de tag, não fazer retry automático. Registrar status, rota abstrata,
latência, tentativa e motivo sanitizado, nunca token, IP, Player Tag ou payload
bruto.

O portal oficial informa qualitativamente que tokens possuem rate limitations e
que chamadas podem falhar quando excedidas. Limite numérico, janela, escopo por
token/IP e semântica de proxy não foram confirmados. Planejamento deve usar
margem conservadora e jitter; não tratar nenhuma taxa desta evidência como limite
seguro.

## Request model

### Fórmulas

Premissas explícitas, todas estimativas:

- `P` = jogadores acompanhados;
- `S` = sync inicial: `profile + battlelog = 2 × P` uma vez;
- `H` = battle history ativo a cada 5 minutos: `12 × 24 = 288 × P/dia`;
- `R` = refresh de perfil diário candidato: `1 × P/dia`;
- `C` = catálogo global a cada 6 horas candidato: `4/dia`, independente de `P`;
- `L` = locations global diário candidato: `1/dia`, independente de `P`;
- `B` = batch de meta limitado a 100 seeds/dia: `2 × B = 200/dia` para
  profile+battlelog, podendo deduplicar com requests de usuários;
- `K` = cohort de clan com 50 perfis: `3 + K = 53 requests/cohort`, conforme
  `001-04`; dois cohorts/dia adicionam `106/dia`;
- `M` = meta background bounded: `200 + 106 = 306/dia` no exemplo;
- `D(P)` = cenário de impacto: `H + R + C + L + M = 289P + 311 requests/dia`.

O refresh diário de perfil não substitui sync inicial; `S` é custo de onboarding
ou primeira sincronização. Cache hit pode reduzir trabalho upstream, mas não foi
descontado de `D(P)`.

### Operações de meta

| Operação | Fórmula | Uso permitido no MVP |
| --- | --- | --- |
| seed profile + battle log | `2 × B` | batch pequeno e explícito; não crawler |
| clan cohort | `3 + K`; `K=50 → 53` | enriquecimento opcional e limitado |
| opponent expansion de um salto | `1 + U`; `U=30 → 31` | não executar em escala; viés e retenção pendentes |
| graph de dois saltos | `1 + U + U×U`; `U=30 → 931` | fora do MVP |

Meta global por Arena não está resolvida. O batch acima mede uma estratégia
bounded de candidatos personalizados, não um dataset representativo nem uma
promessa de freshness.

## Projeção de escala

Tabela usa cenário de impacto em que todos os usuários são acompanhados com
polling de 5 minutos e recebem refresh de perfil diário. Não é carga executada.

| Usuários `P` | Sync inicial uma vez | User-driven/background por dia (`289P`) | Meta+catálogo+locations compartilhados por dia | Total estimado por dia | Total estimado em 30 dias (inclui sync inicial) |
| ---: | ---: | ---: | ---: | ---: | ---: |
| 1.000 | 2.000 | 289.000 | 311 | 289.311 | 8.681.330 |
| 10.000 | 20.000 | 2.890.000 | 311 | 2.890.311 | 86.729.330 |
| 100.000 | 200.000 | 28.900.000 | 311 | 28.900.311 | 867.209.330 |

Somente battle history a cada 5 minutos representa, em média, `P/300` chamadas
por segundo: aproximadamente `3,3`, `33,3` e `333,3` req/s para 1k, 10k e 100k.
Sem jitter, sincronização de cadência cria picos ainda maiores. Limite oficial e
capacidade do proxy são desconhecidos; o cenário de 100k não pode ser aceito como
premissa de MVP sem scheduler/fila/egress e teste de contrato fora da API oficial.

Uma política mais segura é acompanhar somente jogadores ativos/opt-in, usar
refresh de perfil sob demanda ou diário, e tornar batch de meta compartilhado e
bounded. Nesse desenho, `P` deve ser número de jogadores efetivamente tracked,
não total de contas registradas.

## Vercel, allowlist e portabilidade

### Fatos documentados

- Vercel documenta que deployments não têm IP fixo por padrão; Functions usam
  faixa dinâmica de saída.
- Vercel Static IPs está disponível em planos Pro e Enterprise, usa egress estático
  compartilhado, cobra US$100/mês por projeto Pro mais Private Data Transfer e
  atribui um par por região ativa. A documentação consultada permite até três
  regiões; três pares podem exceder a capacidade de cinco entradas CIDR observada
  na configuração da chave da API.
- Static IPs é outbound-only, aplica-se ao projeto/ambientes, não ao Edge Runtime
  ou Routing Middleware, e não fornece isolamento dedicado.
- Secure Compute fornece rede dedicada em Enterprise com preço customizado; é mais
  forte para isolamento/VPC peering, mas não deve ser requisito do domínio.
- Vercel Functions podem ser fixadas em região; isso melhora proximidade/latência,
  mas não cria IP de saída estável.

Fontes Vercel consultadas em `2026-10-02`: [Static IPs](https://vercel.com/docs/networking/static-ips),
[fixed IP guidance](https://vercel.com/kb/guide/can-i-get-a-fixed-ip-address),
[Function regions](https://vercel.com/docs/functions/configuring-functions/region).

### Comparação operacional

| Opção | IP allowlist | Vantagem | Custo/risco | Decisão nesta fase |
| --- | --- | --- | --- | --- |
| Vercel sem Static IP | não atende requisito atual | menor configuração | egress dinâmico; rota oficial retorna `invalidIp` | não usar para chamada direta |
| Vercel Static IPs | atende outbound em uma região (2 IPs) ou potencialmente duas (4); três regiões podem exceder cinco entradas permitidas pela chave | mantém Function/runtime inicial | Pro/Enterprise; US$100/projeto/mês + transferência; pool compartilhado; sem Edge/Middleware; sem isolamento dedicado | opção adequada com limite explícito de regiões e validação da allowlist |
| Vercel Secure Compute | atende com IP dedicado | isolamento e VPC peering | Enterprise/custom pricing; maior complexidade | não necessário para MVP; avaliar somente requisito concreto |
| gateway/worker externo | atende com IP controlado pelo gateway | provider substituível; isolamento de adapter; pode operar fora Vercel | componente, custo, observabilidade e SLA adicionais | fallback preferido se Static IP não for adequado |
| RoyaleAPI Proxy atual | observadamente contorna egress local | já produziu probes 2xx | key handling, retenção, limites, SLA e compatibilidade `UNRESOLVED` | transporte temporário condicional, não contrato final |

O domínio deve depender de uma porta `ClashRoyaleClient`/HTTP adapter, não de
Static IPs, Secure Compute, gateway ou proxy. Jobs conceituais (`syncPlayer`,
`ingestMetaBatch`, `refreshCatalog`) devem receber scheduler externo e manter
retry, cache, idempotência e observabilidade no boundary operacional.

## Cadence e mitigação candidata

| Fluxo | Cadência candidata | Mitigação obrigatória |
| --- | --- | --- |
| profile sync | onboarding + manual; refresh diário apenas para tracked | cache por tag, coalescing, timeout, estado `stale` explícito |
| battle history | 5m somente para jogador ativo/opt-in; menos frequente ou sob demanda fora disso | jitter, cursor não assumido, fingerprint idempotente, teto de retries |
| card catalog | 6h global como hipótese operacional | cache global, single-flight, fallback para último catálogo válido |
| locations | diário ou sob demanda | cache global; não usar para health check |
| meta ingestion | batch diário bounded; sem expansão global | deduplicação, limite de seeds/cohorts, sem dois saltos, gate de retenção/terms |
| 429/5xx | event-driven, não cadence | honrar `Retry-After`, backoff+jitter, circuit breaker e fila futura |

## Handoff e requisitos para Fase 002

Fase 002 pode assumir:

- requests server-side atrás de adapter substituível;
- catálogo/locations globalmente cacheáveis e perfil/battle log privados por tag;
- `Retry-After` como contrato de controle quando presente, sem número fixo de
  rate limit;
- polling de 5m somente como opção para tracked active players;
- meta como aquisição híbrida bounded e personalizada, não meta representativa de
  Arena;
- métricas mínimas: status, latência, cache header, retry count, stale age,
  request class (`user-driven`/`background`) e provider route, sanitizadas.

Fase 002 não pode assumir:

- IP fixo na Vercel sem habilitar produto de egress;
- limite numérico, SLA ou `Retry-After` sempre presente;
- que cache do proxy elimina quota/custo;
- polling de todos os usuários em 5m;
- retenção/redistribuição permitida de dados de oponentes, clans ou meta;
- API key, proxy ou Static IP como autorização comercial.

## Riscos residuais

1. Limite oficial numérico, janela de rate limit e comportamento de `429`/`5xx`
   continuam desconhecidos porque não houve stress test.
2. Headers de cache podem mudar e foram observados via proxy, não em origem direta.
3. Latência medida é amostra pequena; não há p95/p99, cold-start, concorrência ou
   disponibilidade de Vercel medida.
4. Proxy pode aplicar limites, retenção ou tratamento de chave diferentes da API;
   dependência permanece condicional.
5. Static IPs pode não ser economicamente adequado em escala; três regiões podem
   exceder as cinco entradas de IP observadas na chave; gateway externo ainda
   precisa de seleção, segurança, SLA e plano de saída.
6. Projeções não incluem retries, outages, duplicação, bytes, custo Vercel,
   Firestore ou transferência; números são ordem de grandeza de requests.
7. Dados de perfil, opponents e clan members mantêm gates de privacidade, termos e
   uso permitido; nenhum resultado desta task fecha `001-06`.
