# 001-04 — Validar aquisição do meta

- **Ticker:** `001`
- **Número:** `04`
- **Status:** `pending`

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

- **Status final:**
- **Estratégia preferida:**
- **Fallback:**
- **Cobertura não resolvida:**
- **Estimativa de requests:**
- **Riscos residuais:**
