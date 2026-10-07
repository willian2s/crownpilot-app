# 002 — Fundação da aplicação e identidade persistente

- **Status geral:** pending
- **Spec:** [002-fundacao-aplicacao-identidade-persistente.md](../../specs/002-fundacao-aplicacao-identidade-persistente.md)
- **Progresso:** 2/12 subtarefas concluídas

## Objetivo

Preparar frontend React + TypeScript + Vite, API ASP.NET Core + C#, Firebase
Authentication com Google e PostgreSQL no Supabase via EF Core + Npgsql,
entregando identidade CrownPilot e vínculo de Player Tag público/read-only com
`ownershipStatus: unverified`. Fundação usa Clean Architecture pragmática em um
Modular Monolith, REST `/api/v1`, OpenAPI, ProblemDetails, Docker/OCI, Render
inicial portátil, testes e observabilidade.

## Checklist

- [x] [002-01-bootstrap-toolchain.md](002-01-bootstrap-toolchain.md)
- [x] [002-02-estabelecer-boundaries-contrato-base-e-ambientes.md](002-02-estabelecer-boundaries-contrato-base-e-ambientes.md)
- [ ] [002-03-preparar-postgresql-migrations-e-harness-rls.md](002-03-preparar-postgresql-migrations-e-harness-rls.md)
- [ ] [002-04-implementar-google-sign-in-e-firebase-bearer.md](002-04-implementar-google-sign-in-e-firebase-bearer.md)
- [ ] [002-05-implementar-port-e-adapter-de-lookup.md](002-05-implementar-port-e-adapter-de-lookup.md)
- [ ] [002-06-modelar-persistencia-repositories-e-rls.md](002-06-modelar-persistencia-repositories-e-rls.md)
- [ ] [002-07-implementar-casos-de-uso-e-api-v1.md](002-07-implementar-casos-de-uso-e-api-v1.md)
- [ ] [002-08-entregar-frontend-de-identidade-e-vinculo.md](002-08-entregar-frontend-de-identidade-e-vinculo.md)
- [ ] [002-09-revisar-arquitetura-frontend-e-ux-visual.md](002-09-revisar-arquitetura-frontend-e-ux-visual.md)
- [ ] [002-10-instrumentar-observabilidade-health-e-redaction.md](002-10-instrumentar-observabilidade-health-e-redaction.md)
- [ ] [002-11-automatizar-ci-oci-e-gates-de-release.md](002-11-automatizar-ci-oci-e-gates-de-release.md)
- [ ] [002-12-validar-staging-e2e-smoke-e-handoff.md](002-12-validar-staging-e2e-smoke-e-handoff.md)

## Observações

- Fase 001 liberou somente bootstrap, identidade e vínculo privado read-only.
- `API → Application → Domain`; `Infrastructure → Application` e `Infrastructure →
  Domain`. `Domain` não conhece Firebase, Supabase, EF Core, Npgsql, HTTP, ASP.NET
  Core ou Vercel; `Application` usa abstrações; `Infrastructure` implementa
  adapters. `Identity` e `PlayerLink` são módulos funcionais; providers ficam em
  `Infrastructure`.
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
- Swagger UI foi adicionada em `/docs` para Development/Staging, consumindo
  `/openapi/v1.json`; permanece apenas camada de visualização, sem alterar
  contrato OpenAPI.
- Render executa imagem Docker da API e é alvo inicial preferido do frontend
  estático; Vercel é alternativa, Azure é destino futuro possível. Mesmo digest
  OCI deve poder ser promovido entre ambientes.
- Mac/Linux usam toolchain pinado e comandos comuns; código didático explica
  intenção, boundaries e comportamento não óbvio sem comentários artificiais.
- `002-04-implementar-google-sign-in-e-firebase-bearer.md` entrega validação
  Firebase e contrato/fakes de `EnsureCrownPilotUser`;
  `002-06-modelar-persistencia-repositories-e-rls.md` integra o Ensure ao schema,
  repositories e RLS reais.
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
- `002-04-implementar-google-sign-in-e-firebase-bearer.md` e
  `002-05-implementar-port-e-adapter-de-lookup.md` podem avançar em paralelo
  depois de `002-01-bootstrap-toolchain.md` e
  `002-02-estabelecer-boundaries-contrato-base-e-ambientes.md`;
  `002-06-modelar-persistencia-repositories-e-rls.md` depende do banco e dos
  contratos; `002-07-implementar-casos-de-uso-e-api-v1.md` integra auth, lookup e
  persistência; `002-08-entregar-frontend-de-identidade-e-vinculo.md` entrega UI
  antes da revisão arquitetural/visual em
  `002-09-revisar-arquitetura-frontend-e-ux-visual.md`.
  `002-10-instrumentar-observabilidade-health-e-redaction.md` e
  `002-11-automatizar-ci-oci-e-gates-de-release.md` fecham observabilidade e gates
  antes de `002-12-validar-staging-e2e-smoke-e-handoff.md`.
- `002-01-bootstrap-toolchain.md` permanece primeiro na ordem e foi concluída
  com restrição de Docker local. Os comandos de banco de
  `002-03-preparar-postgresql-migrations-e-harness-rls.md` continuam contrato
  documental para essa subtarefa posterior, não são executados pelo bootstrap.
- Staging exige workflow protegido/manual, owner, aprovação, migration job e
  smoke antes de promoção. Production é etapa controlada posterior, não requisito
  para bootstrap local.
- `002-09-revisar-arquitetura-frontend-e-ux-visual.md` deve revisar código React
  existente e comportamento visual em viewports
  mobile/desktop; não é uma task genérica de “melhorar frontend”.
- Preview consome somente build/smoke sem login real; Staging e Production
  promovem o mesmo digest OCI, sem rebuild divergente.
- Decisões pendentes de implementação: versão exata de packages, runner .NET,
  controller versus Minimal API, transporte final do pooler, provider/egress live
  e região/backups antes de dados reais.
- Fase 003 permanece bloqueada até reabertura dos gates de API data, retenção,
  ownership, egress, meta e compliance definidos no veredito da Fase 001.
- `002-01` concluiu implementação e validações locais; smoke OCI passou após
  recuperação do Docker Desktop. RLS e E2E seguem bloqueados por pertencerem a
  subtarefas posteriores.
- `002-02` concluiu boundaries, contrato HTTP base, ProblemDetails, CORS,
  exposição OpenAPI por ambiente e fixture bearer local. Firebase real,
  provisionamento e endpoints de negócio seguem bloqueados para subtarefas
  posteriores.
