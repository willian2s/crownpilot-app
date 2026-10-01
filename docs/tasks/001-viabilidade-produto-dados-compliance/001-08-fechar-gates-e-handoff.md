# 001-08 — Fechar gates e handoff

- **Ticker:** `001`
- **Número:** `08`
- **Status:** `pending`

## Requisitos cobertos

- revisão dos critérios da Fase 001 e classificação dos riscos;
- veredito `GO`, `GO WITH CONSTRAINTS`, `GO WITH CONSTRAINTS / APPROVAL DEPENDENCY`
  ou `NO-GO / REDESIGN`;
- handoff explícito para a Fase 002 e atualização consistente do overview.

## Objetivo e resultado esperado

Revisar toda a Fase 001, tomar a decisão de viabilidade e produzir o handoff
que libera — ou bloqueia — a Fase 002.

## Dependências

- 001-01 a 001-07 concluídas ou explicitamente bloqueadas com motivo.

## Arquivos e símbolos prováveis

- `evidences/phase-001-verdict.md`;
- `docs/tasks/001-viabilidade-produto-dados-compliance/001-00-overview.md`;
- `docs/roadmap/crownpilot-roadmap.md` somente se escopo, sequência ou promessa
  do produto mudar;
- estados de fase `GO`, `GO WITH CONSTRAINTS`, `NO-GO / REDESIGN` e release da
  Fase 002; não há código de produção a alterar.

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

## Passos de execução

1. Conferir conclusão ou bloqueio explícito de 001-01 a 001-07.
2. Reconciliar documentação, evidências e critérios da spec.
3. Classificar gates, riscos, constraints e débitos remanescentes.
4. Escolher um único veredito permitido e registrar sua justificativa.
5. Criar o handoff sanitizado e atualizar overview/roadmap somente quando
   necessário.

## Vereditos permitidos

### GO

Todos os gates centrais possuem solução demonstrada.

### GO WITH CONSTRAINTS

O MVP é viável, mas parte da experiência precisa respeitar limitações explícitas.

Exemplos:

- meta segmentado por trophy band em vez de Arena;
- histórico limitado;
- monetização bloqueada ou dependente de aprovação;
- determinado dado de progressão manual.

Quando viabilidade técnica, dados e operação forem suficientes, mas a
sustentabilidade depender de assinatura, paywall, AI/software coaching ou outro
modelo classificado como `REQUIRES EXPLICIT APPROVAL`, o veredito deve ser
`GO WITH CONSTRAINTS / APPROVAL DEPENDENCY`. Isso não é autorização presumida e
não libera billing. A Fase 001 pode continuar tecnicamente viável com MVP gratuito
ou modelo permitido sob condições.

`GO WITH CONSTRAINTS / APPROVAL DEPENDENCY` é qualificador formal de
`GO WITH CONSTRAINTS`: dados, operação e core técnico podem avançar, mas a
sustentabilidade depende de aprovação expressa para um modelo comercial. Não
libera billing nem converte a dependência em autorização.

### NO-GO / REDESIGN

Uma hipótese central não possui caminho sustentável.

Exemplos:

- não existe coleção suficiente para readiness;
- não há estratégia aceitável de candidatos de meta;
- termos proíbem o modelo fundamental.

## Critérios finais

- [ ] api surface documentada;
- [ ] player/collection provados;
- [ ] battle log provado;
- [ ] meta acquisition avaliada;
- [ ] operação/rate/caching estimados;
- [ ] compliance avaliado;
- [ ] cada modelo de monetização possui status, fonte, condições e ação;
- [ ] API access foi separado de autorização comercial;
- [ ] assinatura/paywall, premium features, AI Coach/software coaching, ads,
  donations, coaching humano, sponsorship e SaaS foram analisados;
- [ ] sustentabilidade foi avaliada também para operação sem monetização ainda
  não aprovada;
- [ ] data contracts v0 definidos;
- [ ] riscos classificados;
- [ ] veredito registrado;
- [ ] Fase 002 sabe o que pode assumir.

## Estrutura do handoff

`evidences/phase-001-verdict.md` deve conter:

### Status

GO / GO WITH CONSTRAINTS / NO-GO, com qualificador comercial opcional
`APPROVAL DEPENDENCY`.

### Qualificação comercial

Registrar separadamente se existe `APPROVAL DEPENDENCY`, qual modelo depende dela,
qual receita fica bloqueada e qual caminho gratuito ou permitido sob condições
permanece disponível.

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
- AI/software coaching tratado como automaticamente coberto por coaching;
- Fase 012 tratada como autorização implícita para cobrança;
- API access tratado como autorização comercial;
- segredo/PII em evidência.

## Definição de pronto

Não há gate crítico escondido em "descobriremos durante a implementação".

## Riscos e cuidados

- não liberar Fase 002 com evidência ausente ou contrato inferido;
- não transformar bloqueio comercial em aprovação;
- não marcar subtarefa concluída apenas porque seu arquivo de evidência existe;
- preservar a distinção entre estado planejado e resultado observado.

## Testes e comandos de validação

- conferir cada critério de aceite contra evidência rastreável ou marcar como
  bloqueado;
- executar revisão crítica para amostra única, fonte secundária, campo
  deprecated, meta sem cobertura, custo ignorado, monetização presumida e
  segredo/PII;
- confirmar que checklist do overview continua com exatamente oito itens e que
  progresso corresponde aos itens marcados;
- executar `git diff --check` e verificar links relativos antes do handoff.

## Registro de execução

- **Status final:**
- **Veredito:**
- **Constraints:**
- **Roadmap atualizado:** yes/no + motivo.
- **ADRs criados:**
- **Fase 002:** released/blocked/redesign.
