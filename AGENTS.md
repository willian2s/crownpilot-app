# CrownPilot repository guide

## Scope and sources

- Work is organized as SDD tasks under
  `docs/tasks/002-fundacao-aplicacao-identidade-persistente/`; check
  `002-00-overview.md` and the assigned task before changing code. ADR 005
  (Go) supersedes ADR 004 for stack decisions; neither is permission to
  implement later tasks. The .NET baseline stays until the `002-15` cutover.
- Repository state is authoritative when planning docs lag. Update task status
  and implementation evidence when completing an SDD task.
- Never read, copy, or embed ignored `.env.local`; it may contain secrets.
  `VITE_*` values are public bundle configuration. Server credentials belong
  only in runtime secret storage.

## Toolchain and commands

- Required versions are pinned: Go `1.27.2` (`go` line in `go.mod`), .NET SDK
  `10.0.401` (`global.json`), Node `>=24.20.0 <25`, npm `>=11.19.0 <12`
  (`frontend/package.json`).
- Go tools (`golangci-lint`, `oapi-codegen`) are pinned in the separate
  `tools/go.mod` and run with `go tool -modfile=tools/go.mod <tool>`; do not
  add them to the main `go.mod` or install them globally.
- Initial setup from repository root:

```text
go mod download
dotnet tool restore
dotnet restore CrownPilot.sln
npm ci --prefix frontend
```

- Main verification (Go API first; .NET until the `002-15` cutover):

```text
go vet ./...
go test ./...
go test -race ./...
go tool -modfile=tools/go.mod golangci-lint run
go build ./cmd/crownpilot-api
dotnet build CrownPilot.sln --configuration Release
dotnet test CrownPilot.sln --configuration Release
npm run typecheck --prefix frontend
npm run lint --prefix frontend
npm run test:unit --prefix frontend
npm run build --prefix frontend
npm run smoke --prefix frontend
npm run openapi:check --prefix frontend
```

- Run the Go API locally with
  `CROWNPILOT_ENVIRONMENT=Local go run ./cmd/crownpilot-api` (port `5080`).
- `smoke` and `openapi:check` build the Go API into a temporary directory and
  start it as Local. Override their port with `CROWNPILOT_SMOKE_PORT` if `5089`
  is occupied.
- After editing `api/openapi/v1.json`, regenerate and commit both artifacts;
  CI fails on drift:

```text
go generate ./...
npm run openapi:generate --prefix frontend
```

- Focus one Go test with `go test ./internal/httpapi/ -run 'TestCORS/allowed'`;
  focus one .NET test with `dotnet test <test-project> --filter
  "FullyQualifiedName~ClassName.TestName"`; focus frontend with
  `npm run test:unit --prefix frontend -- src/App.test.tsx`.
- Write Go code comments in Brazilian Portuguese; identifiers, error messages
  and logs stay in English.
- `npm run smoke:container --prefix frontend` builds/runs the Go Docker image
  and requires a live daemon. `test:rls`, `test:e2e`, `db:*`, and unavailable container
  prerequisites intentionally exit `2` as blocked gates; do not replace them
  with no-op passes.

## Architecture

- Backend is one modular monolith process. The Go API (`cmd/crownpilot-api`
  composition root, `internal/config`, `internal/httpapi`,
  `internal/observability`, `api/openapi`) is the executable bootstrap; the
  ASP.NET Core baseline (`src/`, `tests/`) remains until `002-15`. Frontend is
  an independent static Vite artifact.
- Go routing uses only `net/http` `ServeMux` with method/path patterns;
  middlewares are `func(http.Handler) http.Handler`. `depguard` in
  `.golangci.yml` forbids router frameworks and keeps `internal/identity` and
  `internal/playerlink` free of HTTP, config, platform adapters, and the
  `identity -> playerlink` import.
- .NET dependency direction is enforced by architecture tests: `Domain` has no
  outward dependency; `Application -> Domain`; `Infrastructure -> Application
  + Domain`; `Api -> Application + Infrastructure` only for composition.
- Keep HTTP, DTOs, ProblemDetails, authentication, CORS, health, and OpenAPI in
  `Api`; use cases/ports in `Application`; invariants in `Domain`; provider,
  EF Core, and Npgsql adapters in `Infrastructure`.
- Firebase UID is an authenticated external subject, never domain identity or
  trusted request input. Domain identity is `CrownPilotUserId`. Authentication
  and resource authorization remain separate backend checks.
- Supabase is PostgreSQL hosting, not application backend. EF Core owns schema
  and migrations; versioned SQL owns only RLS, grants, roles, extensions, and
  platform objects. Deployment order is EF migration, security/RLS SQL, then
  sanitized fixtures; never migrate on multi-replica startup.

## HTTP and environments

- REST contract lives under `/api/v1`. OpenAPI is spec-first:
  `api/openapi/v1.json` is the single source, embedded with `go:embed` and
  served unchanged at `/openapi/v1.json`. Go models (`oapi-codegen`) and
  frontend types (`openapi-typescript`) are generated from it and versioned.
  `/docs` (Swagger UI) must consume `/openapi/v1.json`; do not add YAML, a
  second schema generator, or `cmd/crownpilot-openapi`. Mark documented but
  unimplemented operations with `x-crownpilot-planned: true`.
- `/health/live` is process-only and must not contact providers. OpenAPI JSON/UI
  are enabled only in Local and Staging; Preview and Production return `404`.
- The Go API requires `CROWNPILOT_ENVIRONMENT` with an exact name (`Local`,
  `Preview`, `Staging`, `Production`); missing or invalid configuration stops
  the process before it listens. Local port is `5080`; hosting sets `PORT`.
- In the .NET baseline, ASP.NET Core `Development` maps to runtime `Local`; if
  `CrownPilot:Environment` is supplied, it must match `ASPNETCORE_ENVIRONMENT`
  or startup fails.
- CORS is an exact allowlist: no wildcard; non-Local origins require HTTPS.
  Local contract bearer tokens prove `401`/`403`/`200` pipeline behavior only;
  they are not Firebase token validation and must never be enabled elsewhere.
- Errors use ProblemDetails with stable `code` and `traceId`; never expose token,
  Firebase UID, Player Tag, external URL, raw payload, stack, or infrastructure
  detail.
