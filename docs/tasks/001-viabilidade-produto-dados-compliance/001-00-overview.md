# 001 — Viabilidade de produto, dados e compliance

- **Status geral:** planned
- **Spec:** [001-viabilidade-produto-dados-compliance.md](../../specs/001-viabilidade-produto-dados-compliance.md)
- **Progresso:** 0/8 subtarefas concluídas

## Objetivo

Provar, antes de construir a aplicação, que a API, os dados de meta e os
boundaries de compliance sustentam o MVP do CrownPilot.

A fase termina com um veredito explícito: **GO**, **GO WITH CONSTRAINTS** ou
**NO-GO / REDESIGN**.

## Checklist

- [ ] [001-01-mapear-api-oficial-e-auth.md](001-01-mapear-api-oficial-e-auth.md)
- [ ] [001-02-validar-player-profile-e-collection.md](001-02-validar-player-profile-e-collection.md)
- [ ] [001-03-validar-battlelog-e-historico.md](001-03-validar-battlelog-e-historico.md)
- [ ] [001-04-validar-aquisicao-do-meta.md](001-04-validar-aquisicao-do-meta.md)
- [ ] [001-05-medir-operacao-cache-e-rate-limits.md](001-05-medir-operacao-cache-e-rate-limits.md)
- [ ] [001-06-validar-compliance-e-monetizacao.md](001-06-validar-compliance-e-monetizacao.md)
- [ ] [001-07-definir-contratos-de-dados-v0.md](001-07-definir-contratos-de-dados-v0.md)
- [ ] [001-08-fechar-gates-e-handoff.md](001-08-fechar-gates-e-handoff.md)

## Sequência e paralelismo

- 001-01 libera os probes autenticados.
- 001-02 e 001-03 podem avançar em sequência curta após 001-01.
- 001-04 depende do entendimento de player/battle data.
- 001-05 consolida custo operacional das estratégias observadas.
- 001-06 pode avançar em paralelo com 001-02 a 001-05.
- 001-07 só congela contratos v0 depois das evidências técnicas.
- 001-08 fecha a fase e atualiza o handoff.

## Gates que bloqueiam a Fase 002

- coleção insuficiente para Deck Readiness;
- Evolution/Hero ownership não observável e sem fallback aceitável;
- ausência de estratégia sustentável para candidatos de meta;
- custos/rate limits incompatíveis com a coleta mínima;
- uso pretendido incompatível com políticas aplicáveis sem alternativa;
- contrato de dados ainda baseado em suposição.

## Observações

- Infra-base já está definida: Firebase Authentication + Cloud Firestore e
  Vercel como deploy inicial.
- Framework e modelagem física continuam abertos; esta fase não deve escolhê-los
  por acidente.
- Dependências Vercel-specific devem ficar fora do domínio para preservar
  portabilidade.
- Probes podem usar curl ou script descartável.
- Token da API e IPs nunca entram no Git.
- Evidências versionadas devem ser sanitizadas.
- Ranking de topo não será aceito como prova de "meta da Arena" sem análise de
  cobertura e viés.
- Monetização continua bloqueada até 001-06 produzir boundary explícito.
