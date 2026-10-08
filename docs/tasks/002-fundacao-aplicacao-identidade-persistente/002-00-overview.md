# 002 — Fundação da aplicação e identidade persistente

- **Status geral:** pending
- **Spec:** [002-fundacao-aplicacao-identidade-persistente.md](../../specs/002-fundacao-aplicacao-identidade-persistente.md)
- **Progresso:** 4/15 subtarefas concluídas

## Objetivo

Concluir frontend React + TypeScript + Vite e backend portátil da fundação,
entregando identidade CrownPilot e vínculo de Player Tag público/read-only com
`ownershipStatus: unverified`. O baseline ASP.NET Core + EF Core permanece
histórico; a migração Go planejada em `002-13` a `002-15` preserva REST
`/api/v1`, OpenAPI, ProblemDetails, Docker/OCI, Render inicial portátil, testes e
observabilidade.

## Checklist

- [x] [002-01-bootstrap-toolchain.md](002-01-bootstrap-toolchain.md)
- [x] [002-02-estabelecer-boundaries-contrato-base-e-ambientes.md](002-02-estabelecer-boundaries-contrato-base-e-ambientes.md)
- [x] [002-03-preparar-postgresql-migrations-e-harness-rls.md](002-03-preparar-postgresql-migrations-e-harness-rls.md)
- [x] [002-04-implementar-google-sign-in-e-firebase-bearer.md](002-04-implementar-google-sign-in-e-firebase-bearer.md)
- [ ] [002-13-bootstrap-http-config-openapi-go.md](002-13-bootstrap-http-config-openapi-go.md)
- [ ] [002-14-autenticacao-firebase-go.md](002-14-autenticacao-firebase-go.md)
- [ ] [002-15-persistencia-cutover-remocao-dotnet.md](002-15-persistencia-cutover-remocao-dotnet.md)
- [ ] [002-05-implementar-port-e-adapter-de-lookup.md](002-05-implementar-port-e-adapter-de-lookup.md)
- [ ] [002-06-modelar-persistencia-repositories-e-rls.md](002-06-modelar-persistencia-repositories-e-rls.md)
- [ ] [002-07-implementar-casos-de-uso-e-api-v1.md](002-07-implementar-casos-de-uso-e-api-v1.md)
- [ ] [002-08-entregar-frontend-de-identidade-e-vinculo.md](002-08-entregar-frontend-de-identidade-e-vinculo.md)
- [ ] [002-09-revisar-arquitetura-frontend-e-ux-visual.md](002-09-revisar-arquitetura-frontend-e-ux-visual.md)
- [ ] [002-10-instrumentar-observabilidade-health-e-redaction.md](002-10-instrumentar-observabilidade-health-e-redaction.md)
- [ ] [002-11-automatizar-ci-oci-e-gates-de-release.md](002-11-automatizar-ci-oci-e-gates-de-release.md)
- [ ] [002-12-validar-staging-e2e-smoke-e-handoff.md](002-12-validar-staging-e2e-smoke-e-handoff.md)

## Observações

- ADR 005 está `accepted` em 2026-10-08. `002-13` é a próxima task liberada;
  `002-14` e `002-15` seguem em sequência, cada uma com gate verde próprio. Não
  iniciar `002-05` a `002-12` antes da conclusão de `002-13` a `002-15` e do
  cutover definido em `002-15`, independentemente de runtime antigo.
- `002-01` a `002-04` permanecem concluídas como histórico do baseline .NET; não
  desmarcar retroativamente.
- Fase 001 liberou somente bootstrap, identidade e vínculo privado read-only.
- Baseline histórico usou `API → Application → Domain` e `Infrastructure →
  Application/Domain`; na migração Go, `identity` e `playerlink` são módulos
  funcionais, `playerlink` pode importar `identity`, e `identity` não importa
  `playerlink`.
- Supabase é provedor do PostgreSQL, não backend da aplicação. Durante o baseline
  EF Core foi dono do schema; `002-15` transfere ownership para migrations SQL Go,
  mantendo SQL separado somente para RLS/grants/objetos de plataforma.
- RLS usa role sem `BYPASSRLS` e contexto transacional de `CrownPilotUserId`; a
  bridge com pooler é gate explícito, não suposição.
- Replace usa `expectedVersion` e `409` em conflito; não há ETag paralelo;
  `subject_type` e
  `ownership_status` são invariantes `NOT NULL`.
- Frontend envia Firebase ID Token bearer à API; authentication e authorization
  permanecem responsabilidades distintas do backend.
- O baseline usou OpenAPI code-first via `Microsoft.AspNetCore.OpenApi`; `002-13`
  migra para spec-first `api/openapi/v1.json`, `oapi-codegen` e
  `openapi-typescript`, sem especificação concorrente.
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
- `002-01` a `002-04` permanecem histórico concluído do baseline. As tasks de
  migração `002-13` -> `002-14` -> `002-15` devem completar seus gates antes de
  `002-05` iniciar no backend Go; `002-06-modelar-persistencia-repositories-e-rls.md`
  depende do banco e dos
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
- A ADR 005, agora aceita, registra Go `1.27.2`, `goose`, porta `5080`, Swagger UI,
  `openapi-check.mjs`, Session pooler no projeto Supabase dev e em hosting,
  conexão direta somente no PostgreSQL Docker local descartável e ausência de
  transaction pooler neste corte. API Render fica em Virgínia (`us-east`) e
  Supabase em `us-east-1` (Northern Virginia); dados pessoais ficam fora do Brasil
  e a transferência internacional precisa constar na política de privacidade antes
  de dados reais. A implementação Go está liberada pela aprovação da ADR 005.
  Pendências operacionais restantes — provider/egress live e backups do
  PostgreSQL/Supabase — continuam gates obrigatórios de `002-12` antes de dados
  reais.
  Controller versus Minimal API e runner .NET são decisões históricas do
  baseline.
- Fase 003 permanece bloqueada até reabertura dos gates de API data, retenção,
  ownership, egress, meta e compliance definidos no veredito da Fase 001.
- `002-01` concluiu implementação e validações locais; smoke OCI passou após
  recuperação do Docker Desktop. RLS e E2E seguem bloqueados por pertencerem a
  subtarefas posteriores.
- `002-02` concluiu boundaries, contrato HTTP base, ProblemDetails, CORS,
  exposição OpenAPI por ambiente e fixture bearer local. Firebase real,
  provisionamento e endpoints de negócio seguem bloqueados para subtarefas
  posteriores.
- `002-03` concluiu runner Docker local, migration EF Core inicial, SQL de
  roles/grants/schema e fail-closed por ambiente. Fixtures de domínio e policies
  concretas permanecem responsabilidades de `002-06`.
- `002-04` concluiu Firebase Web/Google Sign-In, bearer server-side via Firebase
  Admin SDK, contrato Application de resolução de identidade e fixtures locais.
  Emulator/Google real e configuração de credenciais de Staging permanecem
  validação operacional de `002-12`.
