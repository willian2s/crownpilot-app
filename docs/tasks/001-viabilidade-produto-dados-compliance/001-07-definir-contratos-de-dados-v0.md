# 001-07 — Definir contratos de dados v0

- **Ticker:** `001`
- **Número:** `07`
- **Status:** `completed with constraints`

## Requisitos cobertos

- `PlayerSnapshotV0`, `CardCollectionEntryV0`, `CompetitiveContextV0` e
  `CurrentDeckV0` quando suportado;
- provenance, freshness, optionality e separação raw/normalized/derived;
- semântica diferenciada de capability, ownership e deployment de Evolution/Hero,
  com estados `unresolved`/`unavailable` quando a evidência não fecha o contrato.

## Objetivo e resultado esperado

Consolidar as evidências em contratos provisórios de domínio que a Fase 002
pode usar sem acoplar o CrownPilot ao payload bruto da Supercell.

## Escopo incluído

- `PlayerSnapshotV0`;
- `CardCollectionEntryV0`;
- `CompetitiveContextV0`;
- `CurrentDeckV0` quando suportado;
- provenance/freshness;
- optionality;
- classificação source/derived;
- regras de normalização de Player Tag;
- semântica de Evolution/Hero classificada explicitamente, incluindo
  `unresolved`/`unavailable` quando a evidência não fecha o contrato;
- matriz de campos indisponíveis.

## Escopo excluído

- banco/tabelas;
- TypeScript interfaces definitivas;
- migrations;
- API CrownPilot;
- persistence strategy;
- Fit Score.

## Dependências

- 001-01 a 001-06 concluídas ou com limitações explicitamente registradas;
- evidências sanitizadas disponíveis para cada campo incluído;
- decisão de infraestrutura da ADR 001 permanece fora do escopo desta task.

## Arquivos e símbolos prováveis

- `evidences/data-contract-v0.md`;
- `PlayerSnapshotV0`, `CardCollectionEntryV0`, `CompetitiveContextV0` e
  `CurrentDeckV0` como contratos conceituais;
- matriz consolidada de campos e lista de indisponíveis; não criar interfaces
  TypeScript nem persistência porque a `main` não possui runtime.

## Princípios

### Raw != domain

O payload externo não é o domínio interno.

Cada campo do contrato v0 deve ter:

- origem;
- semântica;
- unidade;
- optionality;
- freshness;
- observação de versão/data;
- regra de derivação, se houver.

### Ausência é informação

Diferenciar:

- missing;
- null;
- zero;
- unsupported;
- not observed;
- not applicable.

### Contexto de Evolution/Hero

Não reutilizar cegamente um raw field em três contextos diferentes.

Ownership, capability e deployment devem ser modelados separadamente se os
probes confirmarem essa necessidade.

## Passos de execução

1. Ler api/player/battle/meta findings.
2. Criar matriz final de requirements.
3. Remover campos sem uso no MVP.
4. Criar contratos conceituais v0.
5. Marcar todo campo derived e sua fórmula.
6. Registrar freshness/provenance.
7. Criar lista de dados indisponíveis.
8. Revisar contra as necessidades da Fase 003/005/006.

## Evidência obrigatória

`evidences/data-contract-v0.md`.

Deve conter exemplos sanitizados, não payloads completos de jogadores reais.

## Perguntas de aceite

O contrato consegue responder:

- quais cards o jogador tem?
- em qual nível?
- quais formas especiais possui?
- onde está competitivamente?
- qual deck está usando quando disponível?
- quando esses dados foram buscados?
- de onde cada informação veio?

## Definição de pronto

A Fase 002 pode modelar identidade persistente e a Fase 003 pode implementar
sync sem descobrir novamente a semântica fundamental dos dados.

## Testes e comandos de validação

- conferir cada campo contra as evidências de API, perfil, battle log e meta;
- exigir origem, semântica, unidade, optionalidade, freshness e data para cada
  campo incluído;
- revisar estados `missing`, `null`, `zero`, `unsupported`, `not observed` e
  `not applicable` sem colapsá-los;
- validar exemplos sanitizados e executar `git diff --check`.

## Riscos e cuidados

- v0 é contrato de discovery, não garantia eterna da Supercell;
- campos novos devem ser tolerados;
- campos ausentes não podem quebrar toda a normalização;
- não projetar persistência antes de escolher stack.

## Registro de execução

### Execução em `2026-10-02`

- **Status final:** `completed with constraints`.
- **Arquivos alterados:** `evidences/data-contract-v0.md`; este registro;
  `001-00-overview.md` para marcar somente `001-07` e recalcular progresso.
- **Contratos definidos:** `PlayerSnapshotV0`, `CardCollectionEntryV0`,
  `CompetitiveContextV0` e `CurrentDeckV0` como contratos conceituais; envelope
  de provenance/freshness; estados explícitos de ausência; matriz de dados
  indisponíveis/unresolved.
- **Derived fields:** normalização conservadora de Player Tag; idade/status de
  freshness; coverage observada de coleção; nenhum Fit Score, readiness,
  ownership de Evolution/Hero ou meta prevalence.
- **Decisões:** separar capability de catálogo, estado de coleção e deployment;
  preservar `cards[]`, `currentDeck[]`, `supportCards[]` e
  `currentDeckSupportCards[]`; manter `public_profile`/`unverified`; tratar
  `maxEvolutionLevel` como capability e `evolutionLevel` como sinal unresolved;
  não sintetizar Hero; usar `fetchedAt` obrigatório, referências de provenance por
  fonte e `Cache-Control` apenas como hint; manter matriz própria para campos de
  deployment do current deck.
- **Dados indisponíveis:** ownership verificado; coleção completa e semântica de
  `count`; ownership/deployment confirmado de Evolution; Hero explícito;
  current deck por modo; battle ID/cursor/backfill; meta representativa por Arena;
  permissão de storage/redistribuição de API data.
- **Questões abertas:** contrato oficial autenticado direto; semântica de
  canonicalização de case/charset de Player Tag; semântica de `count` e
  Evolution; Hero; agreements de API/proxy, retenção, privacidade e freshness
  oficial.
- **Desvios:** rota oficial direta continua bloqueada por allowlist; consolidação
  usou evidências via RoyaleAPI Proxy sem promovê-las a contrato oficial. Não
  foram criadas interfaces TypeScript, tabelas, migrations, API ou persistência.
- **Comandos executados:** leitura cruzada da spec, overview, tasks `001-01` a
  `001-06`, `001-08`, ADR 001 e seis evidências; `git diff --check`; validação
  estrutural posterior de ticker, status, checklist e estados de ausência.
- **Resultados/evidências:** `evidences/data-contract-v0.md` registra origem,
  semântica, unidade/tipo, optionalidade, freshness, data/versão e derivação de
  cada campo incluído; referências de provenance distinguem profile, catalog,
  input, decisão e derived; exemplos são sintéticos e não contêm payload, tag,
  nome, token, IP ou e-mail reais. Não existe toolchain de aplicação neste
  baseline, portanto lint, typecheck, build e testes de runtime não se aplicam.
- **Riscos residuais:** shape oficial e semânticas de ownership podem mudar;
  proxy, agreements, retenção, privacidade, limites e SLA seguem pendentes;
  freshness depende de TTL futuro; contrato não decide persistência nem autoriza
  retenção/redistribuição de dados.
- **Revisão independente:** primeira revisão encontrou ambiguidade em estados de
  ausência, provenance multi-fonte e CurrentDeck; correções adicionaram
  `presenceState` separado de availability, envelope `sourceRef` por profile,
  catalog, input, decision e derived, e matriz própria de deployment. Follow-ups
  confirmaram freshness explícita para cada sourceRef e aprovaram `001-07` sem
  blockers. Nenhuma subtarefa seguinte foi iniciada.
