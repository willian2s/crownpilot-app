# 001-07 — Definir contratos de dados v0

- **Ticker:** `001`
- **Número:** `07`
- **Status:** `planned`

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
- semântica de Evolution/Hero confirmada;
- matriz de campos indisponíveis.

## Escopo excluído

- banco/tabelas;
- TypeScript interfaces definitivas;
- migrations;
- API CrownPilot;
- persistence strategy;
- Fit Score.

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

## Riscos e cuidados

- v0 é contrato de discovery, não garantia eterna da Supercell;
- campos novos devem ser tolerados;
- campos ausentes não podem quebrar toda a normalização;
- não projetar persistência antes de escolher stack.

## Registro de execução

- **Status final:**
- **Contratos definidos:**
- **Derived fields:**
- **Dados indisponíveis:**
- **Questões abertas:**
- **Riscos residuais:**
