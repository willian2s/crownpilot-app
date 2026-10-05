# 001 — Viabilidade de produto, dados e compliance

- **Status geral:** completed
- **Spec:** [001-viabilidade-produto-dados-compliance.md](../../specs/001-viabilidade-produto-dados-compliance.md)
- **Progresso:** 8/8 subtarefas concluídas

## Objetivo

Provar, antes de construir a aplicação, que a API, os dados de meta e os
boundaries de compliance sustentam o MVP do CrownPilot.

A fase termina com um veredito explícito: **GO**, **GO WITH CONSTRAINTS** ou
**NO-GO / REDESIGN**.

## Checklist

- [x] [001-01-mapear-api-oficial-e-auth.md](001-01-mapear-api-oficial-e-auth.md)
- [x] [001-02-validar-player-profile-e-collection.md](001-02-validar-player-profile-e-collection.md)
- [x] [001-03-validar-battlelog-e-historico.md](001-03-validar-battlelog-e-historico.md)
- [x] [001-04-validar-aquisicao-do-meta.md](001-04-validar-aquisicao-do-meta.md)
- [x] [001-05-medir-operacao-cache-e-rate-limits.md](001-05-medir-operacao-cache-e-rate-limits.md)
- [x] [001-06-validar-compliance-e-monetizacao.md](001-06-validar-compliance-e-monetizacao.md)
- [x] [001-07-definir-contratos-de-dados-v0.md](001-07-definir-contratos-de-dados-v0.md)
- [x] [001-08-fechar-gates-e-handoff.md](001-08-fechar-gates-e-handoff.md)

## Observações

- Ordem: `001-01` libera probes; `001-02` e `001-03` seguem; `001-04` depende
  deles; `001-05` consolida custo operacional; `001-06` pode ocorrer em paralelo;
  `001-07` congela contratos v0; `001-08` fecha gates e handoff.
- Gates da Fase 002: coleção insuficiente, Evolution/Hero sem fallback,
  ausência de meta sustentável, custo/rate incompatível, uso não permitido ou
  contrato baseado em suposição.
- Infra-base do handoff original era Firebase Authentication + Cloud Firestore e
  Vercel; essa decisão histórica foi supersedida para a Fase 002 pela ADR 004.
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
- `001-02` concluída com constraints: três perfis foram comparados via proxy;
  `cards[]` variou entre 73 e 123 itens contra 123 no catálogo, Evolution
  capability foi observada de forma opcional, mas ownership/deployment de
  Evolution e Hero permaneceu `unresolved`; campo explícito de Hero ficou
  `unavailable` no shape observado. Matriz em
  `evidences/player-field-matrix.md`.
- `001-02` não confirmou semântica de ownership para `count` nem completude
  universal de `cards[]`; `currentDeck[]` e coleção devem permanecer separados.
- `001-03` concluída com constraints: três perfis retornaram 6/30/30 entradas;
  observou-se janela dependente de atividade, sem cursor/paginação/ID estável no
  shape observado; contrato oficial de paginação permanece unresolved;
  polling de 5m é recomendação MVP para jogadores ativos, mas histórico completo
  não pode ser prometido. Evidência sanitizada em
  `evidences/battlelog-findings.md`; rota proxy, semântica de modos, Hero e
  retenção permanecem riscos pendentes.
- `001-04` concluída com constraints: nenhuma estratégia testada demonstrou
  cobertura intermediária suficiente para sustentar “meta por Arena”. Rankings de
  players e Path of Legend retornaram `404` via proxy; cohorts de clans ficaram
  concentrados em 14.000 trophies; battle log fornece candidatos observados, mas
  com viés de atividade/oponente. Estratégia preferida é híbrida e limitada para
  **Best Decks for Your Collection**; meta de Arena permanece não resolvido.
  Evidência sanitizada em `evidences/meta-strategy-comparison.md`.
- `001-05` concluída com constraints: cache, latência, status e ausência de
  `Retry-After` foram consolidados sem provocar `429`; limite numérico permanece
  desconhecido. Cenário de impacto com polling de 5m para 1k/10k/100k projeta
  289.311/2.890.311/28.900.311 requests por dia, antes de retries e sem desconto
  de cache upstream. Egress direto continua bloqueado por IP; Vercel Static IPs
  pode atender em Pro/Enterprise com limite de regiões compatível com a allowlist,
  mas gateway externo continua fallback substituível. Ver
  `evidences/operational-findings.md`.
- `001-06` concluída com constraints: MVP gratuito de guia/análise é condicional;
  cobrança, assinatura, paywall, premium features, analytics pago, SaaS e
  software/AI coaching exigem aprovação expressa; AI Coach gratuito e storage/
  redistribuição de API data permanecem `UNRESOLVED`. Billing continua bloqueado;
  evidência em `evidences/compliance-findings.md`.
- Bloqueio residual de `001-06`: agreements/API terms autenticados, retenção e
  redistribuição de API data e aprovação externa para modelos comerciais ainda
  precisam ser revisados antes de qualquer cobrança; não impede discovery gratuito
  sob constraints, mas não é autorização comercial.
- `001-07` concluída com constraints: contratos conceituais v0 preservam
   provenance/freshness, estados de ausência, separação raw/normalized/derived,
   coleção versus deployment e capability/ownership/deployment de Evolution/Hero.
   Ownership verificado, Hero explícito, completude de coleção e freshness oficial
   continuam indisponíveis ou unresolved; evidência em
   `evidences/data-contract-v0.md`.
- `001-08` concluída com constraints: veredito
  `GO WITH CONSTRAINTS / APPROVAL DEPENDENCY`. Fase 002 foi liberada somente para
  bootstrap, identidade CrownPilot e vínculo privado read-only de perfil público
  (`public_profile` + `ownershipStatus: unverified`). Sync/persistência completa,
  retenção/redistribuição, billing, ownership e meta Arena permanecem bloqueados
  ou unresolved; fallback de produto é `Best Decks for Your Collection`. Veredito
  em `evidences/phase-001-verdict.md`.
- A Fase 001 está encerrada, mas constraints não são autorização comercial nem
  confirmação de contratos oficiais; Fase 002 não deve tratar o handoff como
  liberação de Fase 003/004.
