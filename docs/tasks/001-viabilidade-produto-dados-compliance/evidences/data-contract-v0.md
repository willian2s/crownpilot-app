# Data contracts v0 — CrownPilot

- **Task:** `001-07`
- **Ticker:** `001`
- **Observed at:** `2026-10-02`
- **Status:** `completed with constraints`
- **Scope:** contratos conceituais de discovery; não são interfaces TypeScript,
  schema de banco ou contrato eterno da API.
- **Authority:** [Clash Royale API developer portal](https://developer.clashroyale.com/)
- **Operational route observed:** RoyaleAPI Proxy, server-side, substituindo
  somente o host da API. Observações do proxy não são promovidas a contrato
  oficial.

## Decision summary

Fase 002 pode assumir um snapshot privado e read-only de perfil público, com
coleção, níveis observados, contexto competitivo básico e current deck separado.
Cada snapshot carrega provenance e `fetchedAt`. A ausência de um campo nunca é
convertida automaticamente em `false`, `0` ou array vazio.

O contrato preserva três contextos que não podem ser colapsados:

1. capability do catálogo (`/cards`);
2. estado/sinal retornado no perfil (`cards[]` e `supportCards[]`);
3. deployment retornado no deck (`currentDeck[]` e
   `currentDeckSupportCards[]`).

Evolution possui capability observável, mas ownership e deployment permanecem
`unresolved`. Hero não possui representação explícita observada. O contrato não
cria `heroOwned`, `heroLevel` ou `evolutionOwned` por inferência.

## Evidence boundary and vocabulary

| Classificação | Regra v0 |
| --- | --- |
| `source` | Valor copiado do campo observado, com tipo e ausência preservados. |
| `normalized` | Valor source convertido somente por regra determinística documentada; não corrige semântica desconhecida. |
| `derived` | Valor calculado de campos source/normalized; fórmula obrigatória e resultado não substitui inputs. |
| `available` | Campo observado e utilizável no snapshot, sem afirmar completude ou contrato oficial. |
| `optional` | Campo observado em parte da amostra ou sujeito a modo/temporada/shape. |
| `unresolved` | Campo existe, mas sua semântica de domínio não foi confirmada. |
| `unavailable` | Requisito não apareceu no shape observado; não deve ser sintetizado. |

Ausência usa estados explícitos no envelope de valor, ou equivalente no adapter:

| Estado | Significado | Exemplo | Não converter para |
| --- | --- | --- | --- |
| `missing` | Campo não veio na resposta ou objeto não contém a chave. | `bestTrophies` omitido. | `null`, `0` ou `false`. |
| `null` | Campo veio com valor JSON nulo. | `clan: null`. | `missing` ou objeto vazio. |
| `zero` | Campo veio numericamente com zero. | `count: 0`. | ausência ou ownership confirmado. |
| `unsupported` | Provider/shape não oferece requisito. | Hero explícito. | `false` ou `not observed`. |
| `not observed` | A amostra não permite concluir presença/semântica. | ownership semântico de `count`. | `unsupported` ou `false`. |
| `not applicable` | Campo não se aplica ao contexto. | `trophyChange` em contexto sem troféus. | `0` ou `null`. |

Se implementação futura não usar um wrapper de valor, deve manter o estado em
metadado paralelo; valor ausente não pode ser achatado.

`presenceState` é eixo de presença do valor e usa somente os seis estados de
ausência acima quando o valor não está presente. `availability` é eixo separado
de contrato (`available`, `optional`, `unresolved`, `unavailable`) e não substitui
`presenceState`. Por exemplo: Hero tem `availability: unavailable` e
`presenceState: not observed`; isso não é `unsupported` nem `false`.

## Provenance and freshness envelope

Cada `PlayerSnapshotV0`, cada coleção e cada deck devem apontar para o mesmo
contexto de obtenção, salvo quando uma operação explicitamente combinar fontes.
O envelope conceitual é:

| Campo | Origem | Semântica | Unidade/tipo | Optionalidade | Freshness | Observação/data | Derivação |
| --- | --- | --- | --- | --- | --- | --- | --- |
| `source.authority` | configuração do adapter | Autoridade lógica declarada para endpoint. | enum/string | required | estático | Portal oficial consultado em `2026-10-02`. | Nenhuma. |
| `source.providerRoute` | runtime do adapter | Rota efetiva usada na obtenção. | enum/string | required | snapshot | `royaleapi_proxy` observado em `2026-10-02`; não é endorsement. | Nenhuma. |
| `source.endpoint` | adapter | Endpoint lógico, sem incluir tag real em logs. | string | required | snapshot | `/v1/players/{playerTag}` ou `/v1/cards`. | Template sanitizado. |
| `source.fetchedAt` | relógio do adapter | Instante em que resposta foi recebida. | UTC instant | required | ponto de medição | Obrigatório em cada sync; data de evidência: `2026-10-02`. | Nenhuma. |
| `source.schemaObservedAt` | evidência versionada | Data do shape observado que sustenta normalização. | ISO date | required | versionado | `2026-10-02` para profile/collection/battle; `2026-10-01` para superfície inicial. | Nenhuma. |
| `source.sourceRefs[]` | adapter/normalizer | Referências de todas as fontes que sustentam o objeto. | array de `sourceRef` | required | snapshot | `[profile]`, `[profile, catalog]` ou referência estática explícita. | União dos refs dos campos/inputs. |
| `source.cacheControlObserved` | response header | Hint de transporte; não freshness do domínio. | string/optional | optional | transporte | `max-age` variou por endpoint e chamada. | Nenhuma. |
| `freshness.status` | política do adapter | Estado relativo ao TTL de produto escolhido posteriormente. | enum | required | derivado no sync/read | `fresh`, `stale`, `unknown` ou `static`; nenhum TTL de produção foi decidido. | `fresh` se idade ≤ TTL de produto; `stale` se exceder; `unknown` sem relógio/TTL confiável; `static` para decisão versionada. |
| `freshness.ageSeconds` | `fetchedAt` + relógio de leitura | Idade desde busca, quando calculável. | seconds | optional | leitura | Não confundir com `Cache-Control`. | `now - fetchedAt`. |

`Cache-Control` observado pelo proxy (`cards` 7–49s, profile 36–60s,
battlelog 21–60s) é somente hint operacional. Fase 002 não pode tratá-lo como
garantia de freshness, SLA ou limite de quota. `fetchedAt` é obrigatório mesmo
quando a resposta veio de cache do provider.

### Source references para joins

Quando um contrato combina profile e catálogo, o envelope carrega referências
independentes. Todo campo source/normalized recebe `sourceRef`; todo derivado
recebe `sourceRefs[]` dos inputs. Assim, um `fetchedAt` do profile nunca é
aplicado silenciosamente ao catálogo.

| `sourceRef` | Endpoint/evento lógico | Conteúdo | Envelope obrigatório |
| --- | --- | --- | --- |
| `profile` | `/v1/players/{playerTag}` | identidade, contexto, `cards[]`, `supportCards[]`, current deck | `providerRoute`, `fetchedAt`, `schemaObservedAt`, `sourceRefs: []`, freshness e cache hint opcional |
| `catalog` | `/v1/cards` | capability, metadata e catálogo de cards/support | `providerRoute`, `fetchedAt`, `schemaObservedAt`, `sourceRefs: []`, freshness e cache hint opcional |
| `input` | entrada CrownPilot | Player Tag fornecida para resolução | `fetchedAt` = instante de recebimento, `schemaObservedAt: not applicable`, `sourceRefs: []`, freshness `fresh` no recebimento e idade derivada depois, regra de normalização |
| `decision` | decisão de domínio | `public_profile`, `unverified` e políticas de ausência | `fetchedAt: not applicable`, `schemaObservedAt` = data/versão da decisão, `sourceRefs: []`, freshness `static` |
| `derived` | cálculo do adapter/domínio | freshness, coverage e outros derivados permitidos | `fetchedAt` = instante de cálculo, `schemaObservedAt` = versão do contrato, `sourceRefs[]` dos inputs, fórmula e freshness |

No mesmo snapshot, `profile` e `catalog` podem ter instantes diferentes. Se uma
fonte não foi obtida, a referência é `not observed` e o campo dependente não é
preenchido com default.

Nas matrizes abaixo, o prefixo da coluna `Origem/raw field` é o `sourceRef` do
campo (`profile`, `catalog`, `input`, `decision` ou `derived`). Quando dois
`sourceRef`s aparecem, ambos devem ser preservados; joins não podem ocultar a
origem individual dos valores.

## PlayerSnapshotV0

Representação conceitual de uma leitura de perfil público. `playerTag` identifica
o sujeito consultado; não prova que o usuário CrownPilot possui a conta.

```text
PlayerSnapshotV0
  identity
    playerTag
    subjectType = public_profile
    ownershipStatus = unverified
    playerName?
  competitiveContext?
    arena?
    trophies?
    bestTrophies?
    rankedContext?              # unresolved/optional, não generalizar
  progression?
    collectionLevel?
    kingTowerLevel?
  collection[]
  currentDeck?
  source
  freshness
```

### Identity, context and progression matrix

| Campo de domínio | Origem/raw field | Classe | Semântica | Unidade/tipo | Optionalidade | Freshness | Observado/data | Derivação |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `identity.playerTag` | tag de entrada + path `/players/{tag}` | normalized | Identificador textual do perfil resolvido. | string | required para snapshot resolvido | snapshot | `%23` confirmado no proxy em `2026-10-01`; semântica oficial direta pendente. | `trim`; garantir exatamente um `#` inicial; encode `#` como `%23` somente no path; preservar restante sem uppercase/rewrite não confirmado. |
| `identity.subjectType` | decisão de produto | normalized | Natureza do vínculo CrownPilot. | enum | required | estático | `public_profile`, decisão registrada em `2026-10-01`. | Constante; não derivar de resposta da API. |
| `identity.ownershipStatus` | ausência de endpoint oficial confirmado | normalized | Estado de ownership do sujeito. | enum | required | estático | `unverified`; ownership permanece `UNRESOLVED`. | Nunca elevar por HTTP `2xx`, nome, tag ou `count`. |
| `identity.playerName` | profile `name` | source | Nome retornado pelo provider. | string | optional | snapshot | Observado em três perfis em `2026-10-02`; omitir de evidência versionada. | Nenhuma; aplicar política de PII antes de logs/exposição. |
| `competitiveContext.arena.id` | profile `arena.id` | source | Identificador do contexto de Arena retornado. | integer/string provider-defined | optional | snapshot | Observado em três perfis em `2026-10-02`. | Não mapear para faixa própria sem tabela versionada. |
| `competitiveContext.arena.name` | profile `arena.name` | source | Nome legível do contexto retornado. | string | optional | snapshot | Observado em três perfis em `2026-10-02`. | Nenhuma. |
| `competitiveContext.arena.rawName` | profile `arena.rawName` | source | Rótulo raw adicional do provider. | string | optional | snapshot | Observado em três perfis em `2026-10-02`. | Nenhuma; não tratar como tradução estável. |
| `competitiveContext.trophies` | profile `trophies` | source | Troféus atuais retornados para o perfil. | integer | optional | snapshot | Observado em três perfis em `2026-10-02`. | Nenhuma; unidade é troféu, não Arena. |
| `competitiveContext.bestTrophies` | profile `bestTrophies` | source | Maior valor histórico retornado pelo provider. | integer | optional | snapshot | Observado em três perfis em `2026-10-02`. | Nenhuma; não usar como contexto atual. |
| `competitiveContext.rankedContext` | profile `progress`, Path of Legend, `leagueStatistics` | source/normalized | Contexto ranked/seasonal separado do trophy road. | objeto provider-defined | optional, unresolved | snapshot/season | Estruturas variaram; `leagueStatistics` ausente em A/B e presente em C em `2026-10-02`. | Preservar raw somente quando presente; não normalizar chaves dinâmicas sem novo contrato. |
| `progression.collectionLevel` | profile `collectionLevel` | source | Nível de progressão observado. | integer | optional | snapshot | Observado em três perfis em `2026-10-02`. | Nenhuma; não calcular readiness. |
| `progression.kingTowerLevel` | profile `kingTowerLevel` | source | Nível de King Tower observado. | integer | optional | snapshot | Observado em três perfis em `2026-10-02`. | Nenhuma; não substituir por Arena/trophies. |
| `collection` | profile `cards[]` | normalized container | Linhas de coleção retornadas, não garantia de catálogo completo. | array | required como resposta preservada; pode ser vazio somente se provider retornar vazio | snapshot | 73/123, 123/123 e 123/123 linhas contra catálogo em `2026-10-02`. | Normalizar cada linha sem preencher cards não retornados. |
| `currentDeck` | profile `currentDeck[]` + `currentDeckSupportCards[]` | normalized container | Deployment reportado pelo profile, separado da coleção. | objeto/arrays | optional | snapshot | 8 cartas e 1 support row observados nos três perfis em `2026-10-02`; outras amostras exibiram 7/8. | Preservar cardinalidade; não assumir deck de oito nem active slot sem confirmação. |

`clan`, `progress`, estruturas de league e
`legacyTrophyRoadHighScore` podem permanecer em raw quarantine/metadata, mas não
fazem parte do núcleo mínimo de `PlayerSnapshotV0` até contrato específico.

## CardCollectionEntryV0

Uma entrada representa uma linha retornada em `cards[]`; não representa por si
só ownership completo. A mesma carta no `currentDeck[]` é outra entrada de
contexto, mesmo quando possui o mesmo `id`.

| Campo de domínio | Origem/raw field | Classe | Semântica | Unidade/tipo | Optionalidade | Freshness | Observado/data | Derivação |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `cardId` | `profile.id` | normalized | Identidade numérica/provider-defined da carta. | integer/provider ID | required na linha | `profile` snapshot | Presente em todas as linhas observadas em `2026-10-02`. | Conversão de tipo somente se lossless. |
| `name` | `profile.name` | source | Nome retornado para identificação/apresentação. | string | required na amostra; tolerar missing futuro | `profile` snapshot | Presente nas linhas observadas em `2026-10-02`. | Nenhuma. |
| `level` | `profile.cards[].level` | source | Nível retornado no contexto de coleção. | integer | required na amostra; optional no contrato tolerante | `profile` snapshot | 73/73, 123/123 e 123/123 linhas em `2026-10-02`. | Nenhuma; não comparar raridades em escala única. |
| `maxLevel` | `profile.maxLevel` | source | Teto associado à carta/raridade no row de coleção. | integer | optional | `profile` snapshot | Presente em linhas observadas; varia por raridade em `2026-10-02`. | Nenhuma. |
| `count` | profile `cards[].count` | source | Quantidade numérica retornada pelo provider. | integer | optional | snapshot | Zero observado em `2026-10-02`; semântica de ownership não confirmada. | Não derivar `owned` de presença ou valor diferente de zero. |
| `evolution.capabilityMaxLevel` | `catalog.maxEvolutionLevel` | source | Capability máxima declarada para a carta no catálogo. | integer | optional | `catalog` snapshot | 55/123 itens de catálogo observados em `2026-10-02`. | Nenhuma. |
| `evolution.collectionSignal` | profile `cards[].evolutionLevel` | source | Sinal numérico emitido na linha de coleção. | integer | optional, unresolved | snapshot | 3/45/39 linhas nos perfis A/B/C; semântica de ownership não confirmada em `2026-10-02`. | Preservar valor; não gerar `evolutionOwned`. |
| `evolution.ownershipStatus` | ausência de semântica confirmada | normalized | Estado de ownership de Evolution. | enum | required quando Evolution é avaliada | snapshot | `unresolved` em `2026-10-02`. | Não derivar de `capabilityMaxLevel` ou `collectionSignal`. |
| `hero` | profile/catalog hero-specific fields | source/normalized | Representação de Hero associada à carta. | objeto/status | optional, unavailable | `profile`/`catalog` snapshot | Nenhum campo explícito `hero*` observado em `2026-10-02`. | `availability: unavailable`, `presenceState: not observed`; nunca sintetizar de Champion, ícone ou nome. |
| `rawPresence` | shape da resposta | normalized metadata | Estado de presença de cada campo. | enum | required para campos opcionais | snapshot | Aplicável a toda linha. | Mapear missing/null/zero sem colapsar estados. |

### Capability, ownership and deployment

| Contexto | Fonte | O que o v0 pode afirmar | O que permanece proibido |
| --- | --- | --- | --- |
| Capability | `/v1/cards.items[].maxEvolutionLevel`, `supportItems[]` | Catálogo retornou capability/metadata opcional de Evolution; support catalog é separado. | Tratar capability como carta possuída ou Hero. |
| Collection state | `/v1/players/{tag}.cards[]`, `supportCards[]` | Perfil retornou linhas, `level`, `count` e sinais opcionais; support cards são separados. | Tratar presença, `count > 0` ou `evolutionLevel` como ownership confirmado. |
| Deployment | `currentDeck[]`, `currentDeckSupportCards[]` | Perfil retornou conjunto reportado como current deck, com sinais opcionais. | Afirmar que `evolutionLevel` prova slot ativo, ou que support prova Hero deployed. |

## CompetitiveContextV0

Contexto mínimo para personalização, não distribuição representativa de meta.

```text
CompetitiveContextV0
  arena { id?, name?, rawName? }
  trophies?
  bestTrophies?
  rankedContext?       # separado; optional/unresolved
  source
```

Regras:

- `trophies` é valor atual; `bestTrophies` é high-water mark. Não trocar um pelo
  outro.
- `arena`, `trophies`, `bestTrophies`, `gameMode` e `type` de battle log não são
  intercambiáveis. Contextos `PvP`, `pathOfLegend`, `trail` e `unknown` devem
  permanecer particionados.
- Contexto de perfil não prova “meta da Arena”. A aquisição aprovada é híbrida,
  limitada e personalizada: **Best Decks for Your Collection**.
- `rankedContext` só pode ser preenchido com campos raw presentes; não inventar
  season, league ou ranking quando ausentes.

## CurrentDeckV0

Contrato opcional para deployment reportado pelo profile. Ele não é histórico de
batalhas nem prova de deck ativo em cada modo.

| Campo de domínio | Origem/raw field | Classe | Semântica | Unidade/tipo | Optionalidade | Freshness | Observado/data | Derivação |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `cards[]` | profile `currentDeck[]` | normalized container | Cartas reportadas no deck atual do profile. | array de entradas `CurrentDeckCardV0` | optional | `profile` snapshot | 8 rows nos três perfis principais; 7/8 em amostra de meta, `2026-10-02`. | Mapear linha mantendo contexto; não preencher para oito. |
| `supportCards[]` | profile `currentDeckSupportCards[]` | normalized container | Support cards reportadas separadamente. | array | optional | snapshot | 1 row nos três perfis principais em `2026-10-02`. | Nunca mesclar em `cards[]`. |
| `cards[].cardId` | profile `currentDeck[].id` | normalized | Identidade da carta no contexto de deployment. | integer/provider ID | required na linha | `profile` snapshot | Presente nas linhas observadas em `2026-10-02`. | Não usar para substituir coleção de mesmo ID. |
| `cards[].name` | profile `currentDeck[].name` | source | Nome da carta no contexto de deployment. | string | optional | `profile` snapshot | Observado nas linhas de deck em `2026-10-02`. | Nenhuma. |
| `cards[].level` | profile `currentDeck[].level` | source | Nível retornado no contexto de deployment. | integer | optional | `profile` snapshot | Presente nas amostras; divergiu da coleção no perfil A em `2026-10-02`. | Não mesclar com `cards[].level` da coleção. |
| `cards[].maxLevel` | profile `currentDeck[].maxLevel` | source | Teto retornado no contexto de deployment. | integer | optional | `profile` snapshot | Observado nas linhas de deck em `2026-10-02`. | Nenhuma. |
| `cards[].count` | profile `currentDeck[].count` | source | Quantidade retornada no contexto de deployment. | integer | optional | `profile` snapshot | Observado; não prova ownership. | Não mesclar com `cards[].count` da coleção. |
| `deployment.evolutionSignal` | deck row `evolutionLevel` | source | Sinal de Evolution na linha deployed. | integer | optional, unresolved | snapshot | 0/8, 4/8 e 7/8 linhas nos perfis A/B/C. | Preservar sinal; não converter em active deployment. |
| `mode` | nenhum campo confirmado no current deck | unavailable | Modo ao qual o deck se aplica. | enum | optional/unavailable | n/a | `availability: unavailable`, `presenceState: not observed` em `2026-10-02`. | Não derivar de `arena` ou quantidade de cartas. |
| `rawPresence` | shape da resposta | normalized metadata | Presença de array/linha/campo. | enum | required | `profile` snapshot | Aplicável a cada array/field opcional. | Preservar missing versus `[]` versus `null`. |

`CurrentDeckCardV0` é uma forma conceitual própria, com os mesmos nomes de
identidade quando aplicável, mas `sourceRef: profile` e contexto
`deployment`. Não é alias semântico de `CardCollectionEntryV0`.

Oito cartas é cardinalidade observada, não invariant. Listas de outros modos ou
versões devem ser aceitas e marcadas para revisão sem quebrar normalização.

## Derived fields and formulas

Somente derivados abaixo entram no v0. Nenhum é Fit Score, readiness score ou
prevalência de meta.

| Campo derived | Inputs | Fórmula/regra | Limite |
| --- | --- | --- | --- |
| `identity.playerTag` | tag de entrada | `trim` + exatamente um `#` inicial; `%23` apenas no transporte HTTP | Não uppercase/rewrite caracteres sem contrato oficial confirmado. |
| `freshness.ageSeconds` | `now`, `source.fetchedAt` | `max(0, now - fetchedAt)` quando relógios válidos | Idade local, não idade garantida do provider. |
| `freshness.status` | `ageSeconds`, TTL de produto | `fresh` se dentro do TTL; `stale` se exceder; `unknown` sem TTL/clock | TTL de produção ainda não definido. |
| `collection.coverage` | IDs de `collection[]`, catálogo observado | `observedCardIds / catalogCardIds` com denominador versionado | Mede cobertura da resposta, não ownership. Se catálogo ausente: `not observed`. |

Não derivar `evolutionOwned`, `heroOwned`, `heroDeployed`, `isAccountOwner`,
`arenaMeta`, `readiness`, `fitScore`, `winner`, `battleId` ou completude de
coleção com os dados disponíveis.

## Player Tag normalization and privacy

1. Aceitar input somente após `trim` de whitespace externo.
2. Representar identidade de domínio com exatamente um `#` inicial. Se a camada
   de entrada aceitar tag sem `#`, adicionar o prefixo uma única vez; não aceitar
   `##`.
3. Fazer percent-encoding do `#` (`%23`) no path HTTP. Não registrar URL com tag
   real em logs ou evidência.
4. Preservar case e demais caracteres até confirmação de regra oficial de
   canonicalização; não usar uppercase como normalização silenciosa.
5. Associar perfil a `crownpilotUserId` derivado server-side, com
   `subjectType: public_profile` e `ownershipStatus: unverified`.
6. Comparar/cachear tags somente após a normalização acima e dentro de boundary
   autorizado. Não expor a associação como ownership nem reutilizar dados de um
   usuário sem revisão de privacidade/termos.
7. Falha de lookup, provider indisponível, rate limit e snapshot stale são estados
   operacionais distintos de ownership.

## Matrix of unavailable, unresolved and excluded data

| Requisito | Estado v0 | Evidência/razão | Consequência |
| --- | --- | --- | --- |
| Ownership verificado da conta | `unresolved` | Nenhum endpoint oficial de ownership confirmado. | Somente perfil público read-only; nunca “minha conta”. |
| Coleção completa/ownership por `count` | `unresolved` | `cards[]` variou 73/123 a 123/123; semântica de `count` não fechada. | Preservar coverage e sinal; não preencher ausentes. |
| Evolution capability | `available`/optional | `maxEvolutionLevel` observado no catálogo. | Pode informar capability raw, não ownership. |
| Evolution ownership | `unresolved` | `evolutionLevel` observado, significado não confirmado. | Não expor `evolutionOwned`. |
| Evolution deployment ativo | `unresolved` | Sinal em current deck não prova slot ativo. | Não afirmar forma equipada/ativa. |
| Hero capability/ownership/deployment | `unavailable`/`unresolved` | Nenhum campo explícito `hero*`; ícone não é semântica. | Não responder Hero como fato; manter ausência explícita. |
| Current deck por modo | `unavailable` | current deck foi observado, mas modo aplicável não foi confirmado. | Não prometer deck ativo em todos os modos. |
| Battle result/winner | `not observed` | Battle log possui crowns/trophyChange, sem `result`/`winner`. | Qualquer resultado é derivação futura e rotulada. |
| Stable battle ID/cursor/backfill | `unavailable`/`unresolved` | Nenhum ID universal ou paginação observada. | Histórico é best-effort por polling/fingerprint. |
| Meta representativa por Arena | `unresolved` | Estratégias não demonstraram cobertura intermediária suficiente. | Produto usa contexto personalizado, não prevalência de Arena. |
| API data storage/redistribution permission | `unresolved` | Agreements autenticados e termos de proxy pendentes. | Retenção, exportação e redistribuição permanecem gates. |
| Freshness/SLA oficial | `unresolved` | Cache observado é proxy hint variável; sem SLA. | `fetchedAt` obrigatório; TTL/SLO futuro separado. |

## Sanitized examples

Exemplos são sintéticos e não representam payload completo nem jogador real.

### Snapshot com sinais de ausência preservados

```text
PlayerSnapshotV0 {
  identity: {
    playerTag: "<REDACTED_TAG>",
    subjectType: "public_profile",
    ownershipStatus: "unverified",
    playerName: { presenceState: "missing" }
  },
  competitiveContext: {
    arena: { id: 54000000, name: "<REDACTED_ARENA>", rawName: "<REDACTED>" },
    trophies: 0,
    bestTrophies: { presenceState: "missing" },
    rankedContext: { presenceState: "not observed" }
  },
  collection: [
    {
      cardId: 1,
      level: 10,
      count: 0,
      evolution: {
        capabilityMaxLevel: 1,
        collectionSignal: { presenceState: "missing" },
        ownershipStatus: "unresolved"
      },
      hero: { availability: "unavailable", presenceState: "not observed" }
    }
  ],
  currentDeck: { presenceState: "null" },
  source: {
    authority: "clash_royale_api",
    providerRoute: "royaleapi_proxy",
    endpoint: "/v1/players/{playerTag}",
    fetchedAt: "2026-10-02T00:00:00Z",
    schemaObservedAt: "2026-10-02"
  },
  freshness: { status: "unknown", ageSeconds: { presenceState: "not observed" } }
}
```

Neste exemplo `count: 0` não significa ausência da linha nem ownership
confirmado; `currentDeck: null` não significa `missing`; Hero tem
`availability: unavailable` e `presenceState: not observed` no shape observado;
e `trophies: 0` permanece zero se provider realmente retornar zero.

### Deck com deployment separado

```text
CurrentDeckV0 {
  cards: [
    { cardId: 1, level: 10, deployment: {
        evolutionSignal: { presenceState: "not observed" }
    } }
  ],
  supportCards: [
    { cardId: 9001, level: 3 }
  ],
  mode: { availability: "unavailable", presenceState: "not observed" }
}
```

`supportCards` não é anexado a `cards[]`; nenhum item é chamado de Hero sem
campo/semântica explícita.

## Handoff to next phases

### Fase 002 — identidade e domínio

Pode assumir:

- associação privada de uma Player Tag a `crownpilotUserId` server-side;
- `public_profile` + `unverified`, sem credencial Supercell;
- snapshot versionado com `source`, `fetchedAt`, `schemaObservedAt` e freshness;
- campos ausentes preservados por estado, não preenchidos por defaults semânticos;
- collection, current deck e support cards como contextos separados.

Não pode assumir ownership confirmado, coleção completa, Hero ou autorização para
armazenar/redistribuir API data.

### Fase 003/005/006 — sync, meta e compliance

- sync deve tolerar campos novos, arrays parciais, `null`, missing e shapes
  opcionais sem falhar o snapshot inteiro;
- adapter deve manter endpoint lógico separado de proxy/egress;
- battle log deve usar fingerprint opaco versionado, sem tratar como ID oficial;
- recomendação inicial deve ser **Best Decks for Your Collection**, contextualizada
  por trophies/Arena/mode, não “meta da Arena”;
- retenção de snapshots, dados de opponents/clans e redistribuição aguardam
  revisão de agreements, proxy e privacidade;
- billing e software/AI coaching continuam bloqueados pelos gates de compliance.

## Residual risks

1. Schema oficial autenticado, semântica de `count`/Evolution e contrato de Hero
   ainda não foram confirmados na sessão oficial.
2. Todas as observações live desta consolidação usam proxy operacional; key
   handling, retenção, limites, SLA e compatibilidade com Supercell permanecem
   pendentes.
3. `Cache-Control` e shape podem mudar; freshness do domínio exige TTL próprio e
   estado `stale`, não confiança em header upstream.
4. Normalização de case/charset de Player Tag permanece conservadora até contrato
   oficial direto; adapter deve rejeitar entradas claramente inválidas sem
   registrar o valor real.
5. Contrato v0 não resolve persistência, índices, migrations, retenção ou Fit
   Score; qualquer escolha deve ficar em spec posterior.
