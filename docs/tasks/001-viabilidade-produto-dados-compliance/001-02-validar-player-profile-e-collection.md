# 001-02 — Validar player profile e collection

- **Ticker:** `001`
- **Número:** `02`
- **Status:** `pending`

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

- **Status final:**
- **Perfis/amostra:** registrar apenas descrição sanitizada.
- **Campos confirmados:**
- **Campos deprecated/instáveis:**
- **Campos indisponíveis:**
- **Riscos residuais:**
