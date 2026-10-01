# 001 — Viabilidade de produto, dados e compliance

- **Status geral:** in_progress
- **Spec:** [001-viabilidade-produto-dados-compliance.md](../../specs/001-viabilidade-produto-dados-compliance.md)
- **Progresso:** 1/8 subtarefas concluídas

## Objetivo

Provar, antes de construir a aplicação, que a API, os dados de meta e os
boundaries de compliance sustentam o MVP do CrownPilot.

A fase termina com um veredito explícito: **GO**, **GO WITH CONSTRAINTS** ou
**NO-GO / REDESIGN**.

## Checklist

- [x] [001-01-mapear-api-oficial-e-auth.md](001-01-mapear-api-oficial-e-auth.md)
- [ ] [001-02-validar-player-profile-e-collection.md](001-02-validar-player-profile-e-collection.md)
- [ ] [001-03-validar-battlelog-e-historico.md](001-03-validar-battlelog-e-historico.md)
- [ ] [001-04-validar-aquisicao-do-meta.md](001-04-validar-aquisicao-do-meta.md)
- [ ] [001-05-medir-operacao-cache-e-rate-limits.md](001-05-medir-operacao-cache-e-rate-limits.md)
- [ ] [001-06-validar-compliance-e-monetizacao.md](001-06-validar-compliance-e-monetizacao.md)
- [ ] [001-07-definir-contratos-de-dados-v0.md](001-07-definir-contratos-de-dados-v0.md)
- [ ] [001-08-fechar-gates-e-handoff.md](001-08-fechar-gates-e-handoff.md)

## Observações

- Ordem: `001-01` libera probes; `001-02` e `001-03` seguem; `001-04` depende
  deles; `001-05` consolida custo operacional; `001-06` pode ocorrer em paralelo;
  `001-07` congela contratos v0; `001-08` fecha gates e handoff.
- Gates da Fase 002: coleção insuficiente, Evolution/Hero sem fallback,
  ausência de meta sustentável, custo/rate incompatível, uso não permitido ou
  contrato baseado em suposição.
- Infra-base já está definida: Firebase Authentication + Cloud Firestore e
  Vercel como deploy inicial.
- Framework e modelagem física continuam abertos; esta fase não deve escolhê-los
  por acidente.
- Dependências Vercel-specific devem ficar fora do domínio para preservar
  portabilidade.
- Probes podem usar curl ou script descartável.
- Token da API e IPs nunca entram no Git.
- Evidências versionadas devem ser sanitizadas.
- `001-01` foi concluída com constraints: a rota direta continua bloqueada por allowlist, mas a nova chave produziu 2xx via RoyaleAPI Proxy para cards, profile, battlelog e locations; ownership, limites e termos third-party continuam pendentes.
- RoyaleAPI Proxy é o transporte operacional atual até decisão de troca ou egress próprio com IP fixo; probes com chave rotacionada funcionaram. Retenção/tratamento de chave, limites, segurança e compatibilidade com Supercell permanecem `UNRESOLVED`; uso continua server-side, sem browser exposure da chave, billing ou autorização comercial implícita.
- Ranking de topo não será aceito como prova de "meta da Arena" sem análise de
  cobertura e viés.
- Monetização continua bloqueada até 001-06 produzir boundary explícito.
- Fan Content Policy trata cobrança por funcionalidades como dependente de
  aprovação expressa; ads, donations e coaching permanecem exceções condicionais,
  e software/AI coaching não é automaticamente coberto.
- Player Tag pode salvar perfil público read-only com `ownershipStatus: unverified`;
  não representa conta própria, ownership, exclusividade ou autorização de ação.
