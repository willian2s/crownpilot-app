# 002-15 — Persistência, cutover e remoção do .NET

- **Ticker:** `002`
- **Número:** `15`
- **Status:** `pending`

## Objetivo e resultado esperado

Completar migração técnica para Go: substituir EF Core/Npgsql por `pgx`/`sqlc` e
migrations SQL, provar RLS e pooler, registrar os ajustes de `002-05` a `002-12`,
fazer cutover dos scripts/CI/OCI e remover o backend .NET após os gates da
migração estarem verdes. As tasks de produto `002-05` a `002-12` implementam seus
fluxos depois do cutover, usando os ajustes registrados aqui.

Não há migração de dados legados nem dual-write. A spec 002 continua fonte dos
requisitos funcionais; esta task transporta esses requisitos para a implementação
Go sem relaxá-los.

## Requisitos cobertos

- schema da aplicação com único dono Go e SQL de segurança separado;
- `goose` como ferramenta pinada de migrations;
- `cmd/crownpilot-migrate` como job separado do runtime;
- ordem schema migrations -> SQL de segurança/RLS/grants -> fixtures sanitizadas;
- UUID interno, `firebase_uid` único, `public_profile`/`unverified`, constraints e
  ausência de snapshot/raw/cache;
- bridge transacional `firebase_uid` -> `CrownPilotUserID`, RLS sem `BYPASSRLS` e
  isolamento A/B;
- Session pooler no projeto Supabase dev e em todos os ambientes hospedados;
  conexão direta somente no PostgreSQL local do Docker Compose;
- transaction pooler somente com prova de protocolo simples e prepared statements
  desativados;
- contratos de lookup, vínculo, replace, unlink e delete da spec 002 preservados
  como norma para as tasks de produto pós-cutover;
- migrations compartilhadas forward-only; `down` somente em desenvolvimento local;
- mudanças destrutivas em expand/contract e rollback de aplicação por reimplantação
  da imagem anterior sobre o schema atual;
- remoção de solution/projetos/artefatos .NET após paridade;
- CI, Docker, smoke, staging e frontend usando Go.

## Escopo incluído

- migrations SQL em `database/migrations/`, security/RLS em
  `database/security/`, `pgxpool`, `sqlc` e job de migration;
- repositories e transações sem HTTP nos módulos;
- adaptação dos scripts `db:*`, smoke, OpenAPI, CI, Docker e runbooks;
- teste arquitetural com `depguard` como principal e `go list` somente para lacunas;
- remoção de `.sln`, `.csproj`, código .NET, EF snapshot e tooling após gates;
- atualizar documentação de operações e handoff sem alterar esta ADR para incluir
  critérios funcionais;
- registrar os ajustes de stack das tasks `002-05` a `002-12` listados abaixo;
  implementação desses fluxos permanece nas respectivas tasks pós-cutover.

## Escopo excluído

- dados reais legados, dual-write ou rollback destrutivo automático;
- novos recursos de produto, snapshot, sync, ownership verificado ou provider live
  na CI;
- alteração de requisitos da spec 002;
- remoção de frontend React/Vite, Firebase Web ou PostgreSQL/Supabase.

## Dependências

- `002-13-bootstrap-http-config-openapi-go.md` e
  `002-14-autenticacao-firebase-go.md` com gates verdes;
- ADR 005 aprovada e perguntas pendentes fechadas ou explicitamente adiadas; ambos
  são pré-condições obrigatórias para iniciar o cutover;
- `002-03` como histórico do harness PostgreSQL/RLS;
- conexão PostgreSQL descartável e ambiente de prova Staging;
- gates reais de Docker, PostgreSQL, Firebase Emulator e CI disponíveis.

## Ajustes transportados de 002-05 a 002-12

| Task | Aplicação durante a migração Go |
|---|---|
| `002-05` lookup | port em `internal/playerlink`, adapter `platform/clashroyale` com `context.Context`, `http.Client`, timeout, retry limitado, fixtures e redaction; manter estados, encoding, host server-only e ausência de rede live. |
| `002-06` persistência/RLS | `pgxpool` + `sqlc`, migrations SQL Go, security SQL para RLS/grants, UUID interno, bridge transacional, A/B, version conflict, Data API denial e ausência de snapshot. |
| `002-07` API | handlers `net/http`/`ServeMux`, OpenAPI spec-first único, DTOs, endpoints, auth-before-authorization, statuses, `Retry-After`, `traceId`, redaction e vínculo anterior em falha. |
| `002-08` frontend | manter React/AuthProvider/API client; consumir tipos gerados de `openapi-typescript`; atualizar somente base URL, contrato e mensagens necessárias. |
| `002-09` revisão frontend | revisar o client/contrato Go e corrigir somente achados de segurança, acessibilidade ou manutenção dos fluxos existentes. |
| `002-10` observabilidade | trocar `ILogger`/`Activity` por `slog`, `context` request ID e métricas equivalentes; manter liveness/readiness, baixa cardinalidade e redaction. |
| `002-11` CI/OCI | trocar setup .NET por Go, `go test -race`, `depguard`, geração OpenAPI, `sqlc`, migrations, lint e build multi-stage; manter gates PR/main/Staging. |
| `002-12` Staging | trocar imagem e runbook para Go; manter Preview sem login real, Staging separado, migration job, E2E, smoke, promoção do mesmo digest e go/no-go. |

## Mapa de equivalência final

| Artefato atual | Tratamento após cutover |
|---|---|
| `src/Api/` | `cmd/crownpilot-api` + `internal/httpapi`/`internal/config`. |
| `src/Application/` e `src/Domain/` | `internal/identity` e `internal/playerlink`, sem dependência de provider/HTTP. |
| `src/Infrastructure/Authentication/` | `internal/platform/firebase`; wiring explícito no composition root. |
| `src/Infrastructure/Persistence/` e EF migrations | `internal/platform/postgres`, `sqlc`, `database/migrations` e `cmd/crownpilot-migrate`. |
| `database/security/*.sql` | permanece fonte exclusiva de RLS, roles, grants, extensões e plataforma. |
| `tests/*` | packages Go, contract tests HTTP/OpenAPI, harness PostgreSQL/RLS e testes frontend conforme boundary. |
| `frontend/src/api` e auth | permanece; usa Firebase ID Token e tipos derivados, nunca banco/provider externo. |
| `.github/workflows/ci.yml` | jobs Go/frontend/OpenAPI/DB/RLS/container/smoke, sem secrets reais. |
| `Dockerfile`/`.dockerignore` | imagem multi-stage Go sem `.env`, secrets, testes ou migration automática. |
| `frontend/scripts/*.mjs` e `package.json` | interfaces públicas permanecem; processos chamam API/migration job Go. |
| `docs/operations/*` e `docs/learning/*` | runbooks atualizados para Go, secrets, migration job e composition root. |
| `cmd/crownpilot-openapi` | não existe; fonte `api/openapi/v1.json` é embutida pela API. |
| `global.json`, `.sln`, `.csproj`, `dotnet-tools.json` | removidos somente após o gate final; Go e ferramentas pinadas tornam-se fonte. |

## Paridade funcional obrigatória

Os seguintes comportamentos da spec 002 continuam contratos funcionais para as
tasks de produto pós-cutover. `002-15` não implementa nem testa esses fluxos;
cada teste pertence à task correspondente, sem relaxar a spec:

- `auth_time` para exclusão, com claim ausente, malformada, futura, fora da janela
  de cinco minutos e tolerância de relógio definida pela spec;
- normalização conservadora da Player Tag, exatamente um `#` inicial e `%23`
  somente no path HTTP;
- somente resultado `resolved` persistido como `public_profile`/`unverified`;
- `expectedVersion` ausente na criação, obrigatório no replace e `409` em conflito;
- validação do novo perfil antes de substituir o vínculo antigo;
- unlink e delete explícitos e idempotentes;
- `NOT NULL`/checks, ausência de snapshot/raw/cache e preservação do vínculo
  antigo em falha externa;
- liveness sem banco/provider e readiness sem Clash Royale;
- autorização A/B no backend, RLS como defesa adicional e nenhum UID vindo do
  request.

## Passos de implementação

1. Criar schema/migrations SQL e adaptar repositories/transações.
2. Aplicar security SQL depois do schema e provar contexto RLS em conexão direta
   local e Session pooler no Supabase dev e hospedado.
3. Provar, se necessário, transaction pooler com protocolo simples e prepared
   statements desativados; sem prova, não declarar isolamento nesse modo.
4. Validar que scripts, CI, Docker, OpenAPI e runbooks usam os boundaries Go;
   lookup, API e frontend ficam para as tasks pós-cutover.
5. Executar somente gates de migração, schema, RLS, architecture/imports e
   artefatos, sem exigir testes dos fluxos de produto posteriores.
6. Remover .NET somente após evidência verde do bootstrap, authentication,
   persistência e cutover; registrar rollback não destrutivo.

## Gate verde obrigatório

```text
go test ./...
go test -race ./...
go vet ./...
golangci-lint run
npm run typecheck --prefix frontend
npm run lint --prefix frontend
npm run test:unit --prefix frontend
npm run build --prefix frontend
npm run openapi:check --prefix frontend
npm run smoke --prefix frontend
npm run smoke:container --prefix frontend
npm run db:start --prefix frontend
npm run db:migrate --prefix frontend
npm run db:security --prefix frontend
npm run db:fixtures --prefix frontend
npm run test:rls --prefix frontend
npm run db:stop --prefix frontend
```

`test:e2e` e gates de Staging continuam pré-condicionados a ambiente real e
devem reportar bloqueio explícito quando indisponíveis. Nenhum bloqueio pode ser
substituído por no-op.

## Definição de pronto

- migration job separado aplica schema, security/RLS e fixtures na ordem correta;
- aplicação Go usa somente `CrownPilotUserID` internamente e role sem
  `BYPASSRLS`;
- RLS passa em A/B, contexto ausente/adulterado, rollback/commit e conexão
  direta local/Session pooler dev e hospedado; transaction pooler só passa com prepared
  statements desativados, protocolo simples e prova equivalente;
- teste consulta `pg_class.relrowsecurity` e falha se qualquer tabela do schema
  `crownpilot` estiver sem RLS ativo; CI local e gate antes de Staging;
- ajustes Go de `002-05` a `002-12` estão registrados sem relaxar a spec e as
  tasks futuras estão bloqueadas até este cutover;
- OpenAPI, frontend, scripts, CI e OCI usam fonte/artefatos Go sem drift;
- todos os gates locais/CI aplicáveis passam, e bloqueios externos ficam
  registrados;
- `.NET`, `.sln`, `.csproj`, EF snapshot e comandos antigos são removidos após os
  gates desta migração;
- não há dual-write, migração de dados legados ou rollback destrutivo automático;
- overview registra `3` tasks de migração pendentes/concluídas conforme evidência,
  sem checklist adicional.

## Riscos e cuidados

- Não remover .NET antes de bootstrap, authentication, persistência e smoke
  provarem paridade do baseline migrado.
- Não aplicar migrations no startup de múltiplas réplicas.
- Não confundir conexão direta local/Session pooler dev e hospedado verde com
  transaction pooler verde.
- Não declarar isolamento RLS sem prova do transporte efetivo.
- Não transformar gates bloqueados por ambiente em sucesso falso.
