# 001-02 — Validar player profile e collection

- **Ticker:** `001`
- **Número:** `02`
- **Status:** `completed with constraints`

## Requisitos cobertos

- perfil público versus ownership verificado;
- Arena, trophies, progressão, coleção, níveis, current deck e
  Evolution/Hero ownership/deployment;
- classificação de disponibilidade, derivação, optionalidade e depreciação.

## Objetivo e resultado esperado

Provar quais dados de conta realmente existem hoje e se eles sustentam
`Deck Readiness`, `Best Decks for You` e futuro `Upgrade Planner`.

## Escopo incluído

- consultar múltiplos perfis intencionalmente diferentes;
- comparar `cards[]`, `currentDeck[]` e catálogo `/cards`;
- observar:
  - card id/name;
  - level/maxLevel;
  - count quando existir;
  - evolution fields;
  - hero representation;
  - Arena/trophies;
  - collection/king tower fields atuais;
- identificar campos que parecem deprecated;
- validar diferenças de semântica entre ownership e deployment;
- produzir `evidences/player-field-matrix.md`.

## Escopo excluído

- calcular Fit Score;
- decidir schema de banco;
- persistir perfis;
- coletar inventário econômico que a API não oferece;
- inferir campo ausente sem evidência.

## Dependências

- 001-01 concluída;
- token válido;
- amostra pública permitida de perfis.

## Arquivos e símbolos prováveis

- `evidences/player-field-matrix.md`;
- campos `cards[]`, `currentDeck[]`, catálogo `/v1/cards`, `arena`, `trophies`
  e campos de Evolution/Hero observados;
- contratos conceituais `PlayerSnapshotV0` e `CardCollectionEntryV0`; sem módulo
  de normalização implementado na `main`.

## Estratégia de amostragem

Usar pelo menos 3 perfis e, quando possível, cobrir:

- faixa intermediária;
- faixa avançada;
- Evolution desbloqueada;
- Hero desbloqueado;
- perfil sem clan ou com dados opcionais ausentes.

Se a amostra não cobrir alguma semântica, registrar `unresolved` em vez de
forçar conclusão.

## Passos de execução

1. Buscar catálogo de cards e registrar shape estático.
2. Buscar perfil A e mapear campos.
3. Repetir com perfis B/C.
4. Comparar level/maxLevel entre raridades.
5. Comparar `cards[]` vs `currentDeck[]`.
6. Verificar campos relacionados a Evolution/Hero.
7. Separar:
   - capability da carta;
   - ownership do jogador;
   - deployment no deck.
8. Classificar cada campo:
   - available;
   - derived;
   - optional;
   - deprecated;
   - unavailable;
   - unresolved.
9. Sanitizar tags/nomes nas evidências.

## Evidência obrigatória

`evidences/player-field-matrix.md` com colunas mínimas:

| Requirement | Endpoint | Raw field | Semantics | Optionality | Freshness | MVP status | Evidence |
| --- | --- | --- | --- | --- | --- | --- | --- |

## Validação

A task só passa se conseguirmos responder objetivamente:

- possuímos todos os cards do jogador?
- sabemos o nível utilizável de cada card?
- sabemos se Evo/Hero necessário está disponível?
- sabemos o contexto competitivo básico?
- quais informações de progressão ficam indisponíveis?
- precisamos provar ownership para alguma feature do MVP ou o vínculo público é
  suficiente?

## Definição de pronto

A matriz permite desenhar `PlayerSnapshot v0` sem depender de campos
imaginados.

## Testes e comandos de validação

- repetir os mesmos probes em pelo menos três perfis sanitizados;
- comparar presença, tipo e semântica dos campos em tabela, sem inferência por
  nome histórico;
- validar `git diff --check` e revisar que tags, nomes e payloads pessoais foram
  removidos das evidências;
- marcar como `unresolved` qualquer requisito sem cobertura observada.

## Riscos e cuidados

- updates de 2026 podem ter tornado campos históricos misleading;
- `evolutionLevel` pode mudar de semântica conforme o array;
- não usar média de level como prova de readiness nesta fase;
- não publicar identificadores reais desnecessários.

## Registro de execução

### Execução em `2026-10-02`

- **Status final:** `completed with constraints`.
- **Perfis/amostra:** três perfis públicos sanitizados: A intermediário sem
  clan; B/C avançados com clan. Tags, nomes, clan values e payloads brutos não
  foram registrados.
- **Arquivos alterados:**
  `evidences/player-field-matrix.md`; este registro; overview para marcar somente
  `001-02` e recalcular progresso.
- **Campos confirmados:** `cards[]`, `currentDeck[]`, catalog `items[]`,
  `id`, `name`, `level`, `maxLevel`, `count`, `maxEvolutionLevel`,
  `evolutionLevel`, `arena`, `trophies`, `bestTrophies`, `collectionLevel`,
  `kingTowerLevel`, `supportCards[]` e `currentDeckSupportCards[]`.
- **Campos deprecated/instáveis:** nenhum confirmado como deprecated;
  `legacyTrophyRoadHighScore`, `progress`, Path of Legend e league structures foram
  marcados como instáveis ou dependentes de contrato adicional.
- **Campos indisponíveis:** campo explícito de Hero não apareceu no shape
  observado; ownership/deployment de Hero permanece `unresolved`; sem prova de
  semântica de ownership para `count` ou presença em `cards[]`; sem garantia de
  que `cards[]` sempre cubra catálogo completo.
- **Decisões:** separar capacidade do catálogo, estado de coleção e deployment
  do deck; preservar campos raw antes de normalizar raridades; manter Hero como
  `unresolved`; continuar usando associação `public_profile` com
  `ownershipStatus: unverified`.
- **Desvios:** probes foram executados via RoyaleAPI Proxy porque rota direta
  segue limitada por allowlist de IP. Observações do proxy não foram tratadas
  como contrato oficial.
- **Comandos executados:** probes autenticados server-side com `curl` para
  `/v1/cards` e três `/v1/players/{tag}` usando segredo carregado somente de
  `.env.local`; resumo estrutural Python em diretório temporário; revisão de
  headers selecionados; `git diff --check`.
- **Resultados/evidências:** quatro requests retornaram `200`; catálogo tinha
  123 items e 4 supportItems; perfis tinham 73, 123 e 123 cards, oito cards no
  current deck e suporte separado. `Cache-Control` observado: 7s no catálogo e
  36s nos perfis. Matriz criada em
  `evidences/player-field-matrix.md`.
- **Riscos residuais:** semântica oficial de `count`, completude de `cards[]`,
  Hero ownership/deployment, progressões dinâmicas, contrato direto da API e
  riscos de retenção/key handling/SLA do proxy continuam pendentes.
- **Revisão independente:** solicitada após implementação; confirmou correções
  de separação entre capability, ownership e deployment de Evolution, taxonomia
  controlada da matriz e sanitização. Ressalva estrutural sobre checklists
  preexistentes na spec/001-08 ficou fora do escopo desta subtarefa; nenhuma
  subtarefa seguinte foi iniciada.
