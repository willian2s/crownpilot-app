# Veredito da Fase 001 — viabilidade de produto, dados e compliance

- **Ticker:** `001`
- **Observed at:** `2026-10-02`
- **Status:** `completed with constraints`
- **Veredito:** `GO WITH CONSTRAINTS / APPROVAL DEPENDENCY`
- **Escopo:** fechamento documental dos gates da Fase 001 e handoff controlado
  para a Fase 002; não constitui implementação nem parecer jurídico.

## Status

A Fase 001 está encerrada com constraints. O núcleo técnico de bootstrap,
identidade CrownPilot e vínculo privado de perfil público é viável, mas não há
base para liberar sync completo, retenção/redistribuição de dados, ownership
verificado, meta representativa da Arena ou billing.

O resultado combina dois fatos que não podem ser colapsados:

1. houve acesso operacional autenticado a cards, perfil, battle log e locations
   por RoyaleAPI Proxy server-side;
2. a rota oficial direta continua bloqueada por allowlist de IP, e contrato
   autenticado completo, key handling, retenção, limites, SLA e compatibilidade
   do proxy permanecem condicionais ou `UNRESOLVED`.

### Classificação dos gates

| Gate | Classificação final | Critério/evidência | Efeito |
| --- | --- | --- | --- |
| API access e transporte | `PASS WITH CONSTRAINTS` | 2xx observado para cards, profile, battle log e locations via proxy; rota direta retornou `403 invalidIp`. [api-surface.md](api-surface.md) | Somente server-side, atrás de adapter substituível; proxy não é contrato oficial nem autorização comercial. |
| Identidade e vínculo | `PASS WITH CONSTRAINTS` | Três perfis foram comparados; vínculo permitido é `public_profile` com `ownershipStatus: unverified`. [player-field-matrix.md](player-field-matrix.md) | Fase 002 pode salvar associação privada read-only; não pode afirmar “minha conta”, ownership ou controle. |
| Perfil, coleção e progressão | `PASS WITH CONSTRAINTS` | `cards[]` variou entre 73/123, 123/123 e 123/123; níveis e contexto básico foram observados. Semântica de `count`, completude e Evolution ownership ficaram abertas. [player-field-matrix.md](player-field-matrix.md) | Preservar raw/ausência/proveniência; não preencher cards ausentes nem sintetizar ownership. |
| Battle history | `PASS WITH CONSTRAINTS` | Foram observadas 6/30/30 entradas; não houve cursor, paginação, ID estável, resultado explícito ou Hero explícito. [battlelog-findings.md](battlelog-findings.md) | Histórico é best-effort; polling de 5m é candidato somente para jogadores ativos/tracked com opt-in. |
| Meta da Arena | `BLOCKED AS PRODUCT CLAIM` | Nenhuma estratégia demonstrou cobertura intermediária suficiente; rankings e Path of Legend retornaram 404 via proxy, clans foram top-biased e opponents são enviesados. [meta-strategy-comparison.md](meta-strategy-comparison.md) | Não alegar prevalência ou “melhores decks da Arena”; usar fallback personalizado bounded. |
| Operação e custo | `PASS WITH CONSTRAINTS` | Rate numérico é desconhecido; não houve 429 deliberado; `Retry-After` não apareceu; cenário de impacto 5m projeta 289.311/2.890.311/28.900.311 requests/dia para 1k/10k/100k. [operational-findings.md](operational-findings.md) | Sem polling global; cache é hint, não quota; separar user-driven, opt-in e ingestão bounded. |
| Compliance e sustentabilidade | `APPROVAL DEPENDENCY / BILLING BLOCKED` | App gratuito de guia/análise é condicional; cobrança, SaaS, premium e software/AI coaching exigem aprovação; AI Coach gratuito e storage/redistribuição de API data são `UNRESOLVED`. [compliance-findings.md](compliance-findings.md) | Fase 002 segue sem billing e sem tratar API access como autorização comercial. |
| Contratos e handoff | `PASS WITH CONSTRAINTS` | Contratos v0 preservam provenance, freshness, ausência e separação de contexts; vários campos permanecem `unresolved`/`unavailable`. [data-contract-v0.md](data-contract-v0.md) | Fase 002 pode usar somente contrato conceitual e boundaries explicitados abaixo. |

### Critérios finais

| Critério | Resultado final |
| --- | --- |
| Superfície, autenticação, encoding e erros | Documentados com fonte/data; 2xx e `%23` foram observados no proxy, enquanto origem direta e schema oficial autenticado permanecem pendentes. |
| Player Tag, perfil, coleção, níveis e contexto | Provados em três perfis sanitizados, com completude, `count`, Evolution e Hero classificados sem inferência. |
| Battle log e histórico | Medidos em três perfis, com janela dependente de atividade, modos particionados e perda potencial explicitada. |
| Aquisição de meta | Avaliada por estratégia; somente candidatos personalizados/contextuais bounded foram aprovados. |
| Cache, rate, custo e infraestrutura | Estimados sem stress test; limites numéricos, SLA e comportamento de `429`/`5xx` permanecem desconhecidos. |
| Policy, agreements e monetização | Modelos possuem classificação, fonte, condição e ação; agreements autenticados não disponíveis mantêm storage/redistribuição em `UNRESOLVED`. |
| API access versus autorização comercial | Separação explícita; key, endpoint, IP allowlist ou `2xx` não liberam cobrança. |
| Contratos v0, riscos e handoff | Registrados em [data-contract-v0.md](data-contract-v0.md), com próximo passo e limites de Fase 002 definidos. |

### Revisão crítica de fechamento

| Verificação | Conclusão |
| --- | --- |
| Amostra única | Não usada como prova de contrato: perfil e battle log usaram três amostras; a pequena amostra de meta foi tratada como limitação e não como prevalência. |
| Fonte secundária como autoridade | Não ocorreu: portal oficial e políticas são autoridade; RoyaleAPI Proxy serviu somente como transporte operacional observado. Nenhum SDK, wrapper ou third-party dataset foi promovido a contrato. |
| Campo deprecated | Nenhum campo foi confirmado como deprecated. `legacyTrophyRoadHighScore` ficou `unresolved/legacy` e não entra no núcleo do produto. |
| “Meta da Arena” | Não foi declarada resolvida; rankings, clans e opponents mantêm viés e a estratégia híbrida só sustenta candidatos personalizados. |
| Custo de request | Incluído nas fórmulas e cenários de escala; projeções não são SLO, limite oficial nem autorização para polling global. |
| Monetização | Não foi presumida. Billing, paywall, assinatura e features pagas continuam bloqueados quando dependem de aprovação. |
| AI/software coaching | Não foi tratado como coaching humano automaticamente permitido; software/AI pago exige aprovação e AI Coach gratuito permanece `UNRESOLVED`. |
| Fase 012 | Não é autorização implícita para cobrança; deve consumir aprovação expressa ou operar gratuitamente sob condições. |
| API access | Mantido separado de Fan Content/data use e de autorização comercial; resposta 2xx não altera esse gate. |
| Segredo/PII | Evidências novas são sanitizadas: não contêm token, valor de Player Tag, IP, e-mail, nome ou payload bruto. |

Rastreabilidade principal: [spec da Fase 001](../../../specs/001-viabilidade-produto-dados-compliance.md), [overview](../001-00-overview.md), [ADR 001](../../../decisions/001-firebase-firestore-vercel-portable.md) e os sete artefatos listados em [Resultado principal](#resultado-principal).

## Qualificação comercial

Existe `APPROVAL DEPENDENCY`. A receita dependente de aprovação expressa fica
bloqueada; isso não bloqueia um MVP gratuito sob condições nem transforma esta
classificação em aprovação da Supercell.

| Modelo/uso | Status final | Limite para o produto |
| --- | --- | --- |
| Guia/análise gratuito e Fan Content app | `ALLOWED WITH CONDITIONS` | Não comercial por padrão; revisar assets, marca, dados, disclaimer e policies. |
| Ads, donations, coaching humano e sponsorship | `ALLOWED WITH CONDITIONS` | Cada modelo exige revisão própria; donations não podem desbloquear benefício e nenhum uso pode sugerir endorsement. |
| Assets não modificados e nome `crownpilot` | `ALLOWED WITH CONDITIONS` | Usar somente para exibir/identificar/discutir produtos, sem imitar branding; disclaimer legível. |
| Assets modificados ou domínio/handle com trademark Supercell/nome de jogo | `REQUIRES EXPLICIT APPROVAL` | Não modificar, registrar ou publicar sem permissão/acordo escrito aplicável. |
| Assinatura, paywall, billing, premium features e analytics pago | `REQUIRES EXPLICIT APPROVAL` | Não implementar, prometer ou colocar feature essencial atrás de cobrança. |
| SaaS pago, software coaching, AI/software coaching e AI Coach pago | `REQUIRES EXPLICIT APPROVAL` | Exceção de coaching humano não cobre automaticamente software/IA. |
| AI Coach gratuito | `UNRESOLVED` | Não liberar como capacidade automaticamente permitida. |
| Storage, retenção ou redistribuição de API data | `UNRESOLVED` | Não prometer histórico, exportação ou redistribuição antes de agreements, proxy e privacy review. |
| Bots, mods, automação, private servers, account trading e software não autorizado | `NOT ALLOWED` | Fora do produto, suporte, conteúdo e marketing. |

O disclaimer deve aparecer de forma legível em conexão com Fan Content: “This
material is unofficial and is not endorsed by Supercell. For more information
see Supercell's Fan Content Policy: www.supercell.com/fan-content-policy.”

## Resultado principal

As hipóteses foram separadas em vez de tratadas como uma única prova de produto:

- é possível consultar um perfil público e obter contexto básico, coleção
  observada, níveis e deck atual por rota operacional server-side;
- a coleção não pode ser tratada como universalmente completa, e Evolution/Hero
  não sustentam ownership ou deployment completo;
- battle log pode alimentar histórico próprio best-effort, não histórico completo;
- nenhum método testado prova meta representativa por Arena;
- operação exige adapter, cache privado por tag, opt-in/active-only e limites
  conservadores;
- MVP gratuito de guia/análise permanece possível sob condições, mas monetização
  depende de aprovação expressa.

Rastreabilidade por evidência:

| Evidência | Decisão consumida neste veredito |
| --- | --- |
| [api-surface.md](api-surface.md) | Proxy server-side operacional; rota direta, schema autenticado, ownership e termos continuam condicionais. |
| [player-field-matrix.md](player-field-matrix.md) | Perfil público, contexto básico e separação capability/collection/deployment. |
| [battlelog-findings.md](battlelog-findings.md) | Janela observada, modos, fingerprint opaco e histórico best-effort. |
| [meta-strategy-comparison.md](meta-strategy-comparison.md) | Fallback `Best Decks for Your Collection`; sem meta Arena. |
| [operational-findings.md](operational-findings.md) | Cadência, cache, rate desconhecido, custo e egress substituível. |
| [compliance-findings.md](compliance-findings.md) | Classificações comerciais, billing bloqueado e gates de API data/privacy. |
| [data-contract-v0.md](data-contract-v0.md) | Contract v0, provenance/freshness, ausência explícita e handoff de dados. |

## Decisões

1. **Vínculo:** Fase 002 pode associar uma Player Tag primária, fornecida pelo
   usuário, ao `crownpilotUserId` derivado server-side. O domínio usa
   `subjectType: public_profile` e `ownershipStatus: unverified`; associação é
   privada, read-only, removível e não exclusiva.
2. **Transporte:** RoyaleAPI Proxy é a rota operacional atual, condicional e
   substituível. O adapter deve separar endpoint lógico, provider route e domínio;
   a chave nunca chega ao browser.
3. **Dados:** `PlayerSnapshotV0`, `CardCollectionEntryV0`,
   `CompetitiveContextV0` e `CurrentDeckV0` são contratos de discovery, não
   schema persistente nem garantia eterna da API.
4. **Ausência:** `missing`, `null`, `zero`, `unsupported`, `not observed` e
   `not applicable` permanecem distinguíveis; ausência não vira `false`, `0` ou
   array vazio por conveniência.
5. **Meta:** o produto não usará “meta da Arena” como fato. Candidatos devem ser
   personalizados e contextuais, com provenance, modo e freshness; o fallback
   aprovado é **Best Decks for Your Collection**.
6. **Operação:** polling de 5m é somente opção para tracked active players com
   opt-in. Não há autorização para crawler, expansão global ou polling de todas
   as contas.
7. **Arquitetura:** nenhuma decisão arquitetural nova foi madura o suficiente
   para ADR. A [ADR 001](../../../decisions/001-firebase-firestore-vercel-portable.md)
   continua válida.

## Dados disponíveis

### Source observável, com restrições

- `arena`, `trophies`, `bestTrophies`, `collectionLevel` e `kingTowerLevel` como
  valores raw de contexto/progressão básica;
- linhas observadas de `cards[]`, com `id`, `name`, `level`, `maxLevel`, `count`
  e sinais opcionais de Evolution; `count` permanece sem semântica de ownership;
- `currentDeck[]`, `supportCards[]` e `currentDeckSupportCards[]` em contextos
  separados da coleção;
- capability opcional do catálogo, como `maxEvolutionLevel`;
- campos observados de battle log: tempo, tipo/modo, arena, participantes,
  cartas, support cards, crowns, trophy change, torres e contexto competitivo;
- provenance, endpoint lógico, provider route, `fetchedAt`, `schemaObservedAt`
  e hints de cache, conforme o contrato v0.

### Derived permitido

- normalização conservadora da Player Tag para transporte, sem reescrever case
  não confirmado;
- `freshness.ageSeconds`/status somente com TTL de produto explicitamente
  escolhido;
- coverage observada da coleção contra catálogo disponível, sem chamá-la de
  ownership ou completude;
- partição de modos e candidatos contextualizados;
- fingerprint opaco versionado para idempotência de battle log, sem chamá-lo de
  battle ID oficial.

Nenhum desses derivados pode gerar `heroOwned`, `evolutionOwned`,
`heroDeployed`, `isAccountOwner`, `arenaMeta`, `readiness`, `fitScore`,
`winner` ou completude universal.

## Dados indisponíveis

| Requisito | Estado final | Não prometer |
| --- | --- | --- |
| Ownership verificado | `UNRESOLVED` | “Minha conta”, exclusividade, ações, notificações ou controle da conta. |
| Coleção completa e ownership por `count` | `UNRESOLVED` | Que cards ausentes não existem ou que `count` prova posse. |
| Evolution ownership/deployment ativo | `UNRESOLVED` | `evolutionLevel` como posse/equipamento confirmado. |
| Hero explícito, ownership ou deployment | `UNAVAILABLE / UNRESOLVED` | Hero derivado de Champion, ícone, support card ou nome. |
| Current deck por modo | `UNAVAILABLE` | Que deck atual se aplica a todo modo. |
| Battle ID, cursor, paginação e backfill | `UNAVAILABLE / UNRESOLVED` | Histórico completo ou recuperação após eviction. |
| Result/winner explícito | `NOT OBSERVED` | Resultado como campo oficial; qualquer cálculo é derivado e rotulado. |
| Meta representativa por Arena | `UNRESOLVED` | Prevalência, “best deck da Arena” ou distribuição intermediária. |
| Rate numérico, SLA e freshness oficial | `UNRESOLVED` | Limite seguro, SLO ou `Cache-Control` como TTL de domínio. |
| Storage/retention/redistribution de API data | `UNRESOLVED` | Retenção, exportação pública, redistribuição ou histórico de terceiros. |

## Meta

O gate de meta da Arena permanece fechado. Rankings de players e Path of Legend
retornaram 404 na rota proxy observada; cohorts de clans ficaram concentrados em
14.000 trophies; opponents refletem atividade, seleção e contexto dos jogadores
semente. O fato de uma pequena amostra PvP intermediária existir não demonstra
cobertura populacional.

A estratégia aprovada para próximo desenho de produto é híbrida e limitada:
perfil + battle log do jogador, com enriquecimento opcional e pequeno de clans,
sem misturar `PvP`, `pathOfLegend`, `trail` e `unknown`, e sem expansão global.
O produto deve dizer **Best Decks for Your Collection**, contextualizado por
trophies/Arena/mode e encontros recentes. Não deve dizer “meta da Arena” nem
apresentar prevalência global.

Alternativas futuras — fonte licenciada, segmentação observável adicional ou
redesenho de aquisição — continuam abertas, mas não liberam a alegação atual.

## Operação

- Toda chamada fica server-side atrás de `ClashRoyaleClient`/adapter equivalente;
  provider route, cache e egress são substituíveis.
- Catálogo e locations podem ter cache global. Perfil e battle log usam cache
  privado por tag e não compartilham associação CrownPilot, nomes ou payload bruto
  por padrão.
- `Cache-Control` observado é hint de transporte. Não reduz quota/custo nas
  projeções sem métrica de hit/miss upstream.
- O limite numérico é desconhecido; não houve stress test, 429 deliberado ou
  `Retry-After` observado. Se `Retry-After` surgir, deve ser honrado; sem ele,
  backoff com jitter e teto é requisito futuro do adapter.
- Polling de battle log a cada 5m fica limitado a jogadores ativos/tracked com
  opt-in. Sync de perfil deve ser onboarding, sob demanda ou refresh controlado.
- Cenário de impacto com `P` jogadores tracked é `289P + 311` requests/dia:
  289.311, 2.890.311 e 28.900.311 para 1k, 10k e 100k. São projeções sem
  retries/cache-hit/outage descontados, não SLO nem escopo liberado.
- A rota oficial direta exige egress permitido. Vercel Static IPs ou gateway
  externo são opções futuras; nenhuma vira dependência do domínio.

## Compliance

O app gratuito de guia/análise pode avançar `ALLOWED WITH CONDITIONS`, com
minimização de dados, associação privada, assets dentro da política, ausência de
automação e disclaimer legível. Ads, donations, coaching humano e sponsorship
também são condicionais e exigem revisão própria.

Billing, assinatura, paywall, premium features, analytics pago, SaaS pago,
software coaching e AI Coach pago ficam `REQUIRES EXPLICIT APPROVAL`. AI Coach
gratuito e storage/redistribuição de API data ficam `UNRESOLVED`. Bots, mods,
automação, private servers, account trading e software não autorizado ficam
`NOT ALLOWED`.

API access, Fan Content/data use e autorização comercial são gates independentes.
Uma API key, endpoint, allowlist de IP, proxy ou resposta 2xx não autoriza
cobrança, redistribuição nem tratamento comercial. A Fase 012 deverá consumir
esta dependência; não poderá convertê-la em autorização implícita.

## Riscos e débitos

| Risco/débito | Estado | Saída exigida |
| --- | --- | --- |
| Schema oficial autenticado, rota direta e ownership | Pendente | Revisão oficial autenticada ou artefato oficial sanitizado; manter adapter e `unverified`. |
| RoyaleAPI Proxy | Condicional | Revisar key handling, segurança, retenção, limites, SLA, termos e compatibilidade; manter egress próprio como fallback. |
| Rate, `Retry-After`, cache e freshness | Desconhecido/variável | Definir contrato do adapter, TTL de produto, backoff, métricas e circuit breaker sem stress test. |
| Collection, `count`, Evolution e Hero | Incompleto | Tolerar ausência e novas shapes; não elevar `unresolved` sem nova evidência. |
| Battle history | Best-effort | Validar retenção/privacy antes de persistir; não prometer backfill ou histórico completo. |
| Meta | Sem cobertura Arena demonstrada | Manter fallback personalizado ou obter dataset/licença validado; não fazer crawler/graph expansion. |
| API data, opponents e clan members | `UNRESOLVED` | Revisar agreements, proxy terms, privacidade, retenção, deleção e redistribuição. |
| Monetização e AI/software coaching | Bloqueado por aprovação | Obter aprovação expressa e rastreável antes de billing, marketing ou feature paga. |
| Privacy do vínculo Player Tag | Pendente | Minimizar dados, definir retenção/deleção e publicar privacy policy antes de operação real. |
| Toolchain de aplicação | Inexistente no baseline | Fase 002 deve escolher stack, testes, CI e ambientes; nenhuma validação de runtime é alegada nesta fase. |

## Dependências da Fase 002

Fase 002 está **released with constraints**: pode iniciar bootstrap, identidade
CrownPilot e vínculo privado read-only, mas não pode antecipar sync/persistência
completa.

### Pode assumir

> Os itens de infraestrutura abaixo registram o handoff original da Fase 001.
> Para implementação da Fase 002, consultar ADR 004; somente constraints de
> produto, dados e compliance continuam normativas.

- Firebase Authentication com Google, Cloud Firestore e Vercel como decisões-base
  da [ADR 001](../../../decisions/001-firebase-firestore-vercel-portable.md),
  mantendo domínio e integração portáveis;
- um usuário CrownPilot e uma Player Tag primária, com associação privada
  server-side a `crownpilotUserId`;
- semântica `public_profile` + `ownershipStatus: unverified`, com texto de UI
  equivalente a “Perfil público salvo — ownership não verificado”;
- lookup read-only server-side para validar existência/estado do perfil, com erros
  de provider, rate limit, stale e ownership mantidos como estados distintos;
- boundary de adapter para Clash Royale API, sem expor chave no browser;
- contrato conceitual v0 como referência de proveniência, ausência e separação de
  coleção/deck/support, sem convertê-lo em schema persistente final.

### Não pode assumir

- ownership verificado, coleção completa, semântica de `count`, Evolution/Hero
  ownership/deployment ou current deck por modo;
- sync ou persistência completa de coleção, níveis, Arena, battle history ou
  Player Snapshot; isso fica para fase posterior após os gates correspondentes;
- retenção, redistribuição, exportação ou compartilhamento de API data,
  opponents, clans ou meta derivada;
- proxy como contrato oficial, limite numérico, SLA, freshness garantida ou IP
  fixo sem produto de egress validado;
- polling de todos os usuários, meta Arena, billing, assinatura, paywall,
  premium feature, AI/software coaching pago ou qualquer autorização comercial;
- tratar Fase 012 como autorização implícita.

Antes de liberar sync/persistência da fase seguinte, devem ser fechados pelo menos
agreements/API e proxy, política de retenção/privacy, contrato oficial aplicável,
rate/freshness, egress, semântica dos campos necessários e estratégia de meta
compatível com a promessa escolhida.

## Próximo passo

Liberar somente a Fase 002 para bootstrap reproduzível, quality gates, identidade
Google/Firebase e associação privada read-only de perfil público. O lookup deve
ser server-side, sem credencial Supercell, sem sync completo e sem persistência
de resposta além do mínimo de vínculo autorizado.

Não iniciar Fase 003/004, billing ou AI Coach como consequência deste fechamento.
Depois de Fase 002, reabrir os gates de API data, retenção, operação, meta e
compliance antes de transformar o perfil em snapshot persistente ou promessa de
meta.

## Atualização posterior de infraestrutura

Em `2026-10-05`, a decisão de infraestrutura intermediária foi substituída para
a implementação da Fase 002 pela [ADR 004](../../../decisions/004-aspnet-core-react-vite-firebase-postgresql.md):
ASP.NET Core + C#, React + TypeScript + Vite, Firebase Authentication, EF Core +
Npgsql e PostgreSQL hospedado no Supabase. As constraints de produto,
compliance, ownership não verificado, lookup server-side, retenção e ausência de
sync completo permanecem inalteradas. Este registro é append-only; evidências
históricas da Fase 001 não foram reescritas.
