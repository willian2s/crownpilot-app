# 002 — Fundação da aplicação e identidade persistente

- **Status geral:** pending
- **Spec:** [002-fundacao-aplicacao-identidade-persistente.md](../../specs/002-fundacao-aplicacao-identidade-persistente.md)
- **Progresso:** 0/9 subtarefas concluídas

## Objetivo

Preparar frontend React + TypeScript + Vite, API ASP.NET Core + C#, Firebase
Authentication com Google e PostgreSQL no Supabase via EF Core + Npgsql,
entregando identidade CrownPilot e vínculo de Player Tag público/read-only com
`ownershipStatus: unverified`. Fundação usa Clean Architecture pragmática em um
Modular Monolith, REST `/api/v1`, OpenAPI, ProblemDetails, Docker/OCI, Render
inicial portátil, testes e observabilidade.

## Checklist

- [ ] [002-01-bootstrap-toolchain.md](002-01-bootstrap-toolchain.md)
- [ ] [002-02-estabelecer-arquitetura-contrato-e-ambientes.md](002-02-estabelecer-arquitetura-contrato-e-ambientes.md)
- [ ] [002-03-preparar-postgresql-local-e-migrations.md](002-03-preparar-postgresql-local-e-migrations.md)
- [ ] [002-04-implementar-google-sign-in-e-bearer.md](002-04-implementar-google-sign-in-e-bearer.md)
- [ ] [002-05-criar-adapter-de-lookup.md](002-05-criar-adapter-de-lookup.md)
- [ ] [002-06-persistir-vinculo-com-autorizacao.md](002-06-persistir-vinculo-com-autorizacao.md)
- [ ] [002-07-entregar-fluxos-de-vinculo-e-exclusao.md](002-07-entregar-fluxos-de-vinculo-e-exclusao.md)
- [ ] [002-08-automatizar-quality-gates-e-observabilidade.md](002-08-automatizar-quality-gates-e-observabilidade.md)
- [ ] [002-09-validar-staging-deploy-e2e-smoke-handoff.md](002-09-validar-staging-deploy-e2e-smoke-handoff.md)

## Observações

- Fase 001 liberou somente bootstrap, identidade e vínculo privado read-only.
- `Domain` não conhece Firebase, Supabase, EF Core, Npgsql, HTTP, ASP.NET Core
  ou Vercel; `Application` usa abstrações; `Infrastructure` implementa adapters.
- Supabase é provedor do PostgreSQL, não backend da aplicação. EF Core é dono do
  schema; SQL separado cobre somente RLS/grants/objetos de plataforma.
- RLS usa role sem `BYPASSRLS` e contexto transacional de `CrownPilotUserId`; a
  bridge com pooler é gate explícito, não suposição.
- Replace usa `expectedVersion` e `409` em conflito; não há ETag paralelo;
  `subject_type` e
  `ownership_status` são invariantes `NOT NULL`.
- Frontend envia Firebase ID Token bearer à API; authentication e authorization
  permanecem responsabilidades distintas do backend.
- OpenAPI code-first via `Microsoft.AspNetCore.OpenApi` será contrato único;
  UI navegável não cria especificação concorrente.
- Render executa imagem Docker da API e é alvo inicial preferido do frontend
  estático; Vercel é alternativa, Azure é destino futuro possível. Mesmo digest
  OCI deve poder ser promovido entre ambientes.
- Mac/Linux usam toolchain pinado e comandos comuns; código didático explica
  intenção, boundaries e comportamento não óbvio sem comentários artificiais.
- `002-04` entrega validação Firebase e contrato/fakes de `EnsureCrownPilotUser`;
  `002-06` integra o Ensure ao schema, repositories e RLS reais.
- Authentication handler não escreve no banco; leitura/exclusão não criam usuário
  implicitamente. Criação ocorre no caso de uso de vínculo.
- RLS usa schema dedicado e contexto transacional; prova com pooler é gate, nunca
  premissa. Sem prova, não declarar isolamento por usuário.
- Exclusão de dados exige autenticação recente; não remove Firebase/Google.
- Testes unit, Application, integration, persistence/RLS, authentication/
  authorization, contract, API, E2E e smoke têm gates próprios; Clash Royale
  live fica fora da CI normal.
- Preview e Staging são ambientes diferentes: Preview não depende de URL fixa ou
  login real; Staging possui hostname fixo, Firebase e banco separados.
- Não persistir snapshot, coleção, Arena, battle history, cache ou payload raw.
- Não alegar ownership; usar `public_profile` + `unverified`.
- `002-05` pode avançar em paralelo depois de `002-01` e `002-02`; `002-08`
  estabelece gates mínimos no bootstrap e fecha automação após as features.
- Staging exige workflow protegido/manual, owner, aprovação, migration job e
  smoke antes de promoção. Production é etapa controlada posterior, não requisito
  para bootstrap local.
- Decisões pendentes de implementação: versão exata de packages, runner .NET,
  controller versus Minimal API, transporte final do pooler, provider/egress live
  e região/backups antes de dados reais.
- Fase 003 permanece bloqueada até reabertura dos gates de API data, retenção,
  operação, egress, privacidade e meta.
