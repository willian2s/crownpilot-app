# 001-08 — Fechar gates e handoff

- **Ticker:** `001`
- **Número:** `08`
- **Status:** `planned`

## Objetivo e resultado esperado

Revisar toda a Fase 001, tomar a decisão de viabilidade e produzir o handoff
que libera — ou bloqueia — a Fase 002.

## Dependências

- 001-01 a 001-07 concluídas ou explicitamente bloqueadas com motivo.

## Escopo incluído

- revisar critérios de aceite da spec;
- reconciliar conflitos entre documentação e live probes;
- classificar riscos;
- definir veredito;
- registrar constraints;
- atualizar overview/progresso;
- atualizar roadmap se a hipótese do produto mudar;
- criar `evidences/phase-001-verdict.md`;
- registrar ADRs somente se alguma decisão arquitetural já estiver madura.

## Vereditos permitidos

### GO

Todos os gates centrais possuem solução demonstrada.

### GO WITH CONSTRAINTS

O MVP é viável, mas parte da experiência precisa respeitar limitações explícitas.

Exemplos:

- meta segmentado por trophy band em vez de Arena;
- histórico limitado;
- monetização bloqueada;
- determinado dado de progressão manual.

### NO-GO / REDESIGN

Uma hipótese central não possui caminho sustentável.

Exemplos:

- não existe coleção suficiente para readiness;
- não há estratégia aceitável de candidatos de meta;
- termos proíbem o modelo fundamental.

## Checklist final

- [ ] api surface documentada;
- [ ] player/collection provados;
- [ ] battle log provado;
- [ ] meta acquisition avaliada;
- [ ] operação/rate/caching estimados;
- [ ] compliance avaliado;
- [ ] data contracts v0 definidos;
- [ ] riscos classificados;
- [ ] veredito registrado;
- [ ] Fase 002 sabe o que pode assumir.

## Estrutura do handoff

`evidences/phase-001-verdict.md` deve conter:

### Status

GO / GO WITH CONSTRAINTS / NO-GO.

### Resultado principal

Quais hipóteses foram provadas.

### Decisões

O que passa a ser contrato do produto.

### Dados disponíveis

Resumo do que é source e derived.

### Dados indisponíveis

O que não pode ser prometido.

### Meta

Estratégia aprovada e seus vieses.

### Operação

Cache, rate e restrições de infraestrutura.

### Compliance

O que é permitido, limitado e bloqueado.

### Riscos e débitos

Itens que seguem para fases futuras.

### Dependências da Fase 002

Lista objetiva do que pode ser assumido.

### Próximo passo

Liberar, restringir ou redesenhar Fase 002.

## Validação

Fazer uma leitura crítica final procurando especificamente por:

- conclusão baseada em uma única amostra;
- secondary source tratada como autoridade;
- campo deprecated usado como verdade atual;
- "meta da Arena" sem cobertura demonstrada;
- custo de request ignorado;
- monetização assumida;
- segredo/PII em evidência.

## Definição de pronto

Não há gate crítico escondido em "descobriremos durante a implementação".

## Registro de execução

- **Status final:**
- **Veredito:**
- **Constraints:**
- **Roadmap atualizado:** yes/no + motivo.
- **ADRs criados:**
- **Fase 002:** released/blocked/redesign.
