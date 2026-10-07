# CrownPilot repository guide

## Scope and sources

- Work is organized as SDD tasks under
  `docs/tasks/002-fundacao-aplicacao-identidade-persistente/`; check
  `002-00-overview.md` and the assigned task before changing code. ADR 004
  describes the target architecture, not permission to implement later tasks.
- Repository state is authoritative when planning docs lag. Update task status
  and implementation evidence when completing an SDD task.
- Never read, copy, or embed ignored `.env.local`; it may contain secrets.
  `VITE_*` values are public bundle configuration. Server credentials belong
  only in runtime secret storage.

## Toolchain and commands

- Required versions are pinned: .NET SDK `10.0.401` (`global.json`), Node
  `>=24.20.0 <25`, npm `>=11.19.0 <12` (`frontend/package.json`).
- Initial setup from repository root:

```text
dotnet tool restore
dotnet restore CrownPilot.sln
npm ci --prefix frontend
```

- Main verification:

```text
dotnet build CrownPilot.sln --configuration Release
dotnet test CrownPilot.sln --configuration Release
npm run typecheck --prefix frontend
npm run lint --prefix frontend
npm run test:unit --prefix frontend
npm run build --prefix frontend
npm run smoke --prefix frontend
npm run openapi:check --prefix frontend
```

- `smoke` and `openapi:check` start the Release API with `--no-build
  --no-restore`; run restore/build first. Override their port with
  `CROWNPILOT_SMOKE_PORT` if `5089` is occupied.
- Focus one .NET test with `dotnet test <test-project> --filter
  "FullyQualifiedName~ClassName.TestName"`; focus frontend with
  `npm run test:unit --prefix frontend -- src/App.test.tsx`.
- `npm run smoke:container --prefix frontend` builds/runs Docker and requires a
  live daemon. `test:rls`, `test:e2e`, `db:*`, and unavailable container
  prerequisites intentionally exit `2` as blocked gates; do not replace them
  with no-op passes.

## Architecture

- Backend is one ASP.NET Core modular monolith. `src/Api/Program.cs` is HTTP and
  composition root; frontend is an independent static Vite artifact.
- Dependency direction is enforced by architecture tests: `Domain` has no
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

- REST contract lives under `/api/v1`; generate OpenAPI code-first through
  `Microsoft.AspNetCore.OpenApi`. `/docs` must consume `/openapi/v1.json`; do not
  add handwritten YAML or a second schema generator.
- `/health/live` is process-only and must not contact providers. OpenAPI JSON/UI
  are enabled only in Local and Staging; Preview and Production return `404`.
- ASP.NET Core `Development` maps to runtime `Local`; `Preview`, `Staging`, and
  `Production` are explicit. If `CrownPilot:Environment` is supplied, it must
  match `ASPNETCORE_ENVIRONMENT` or startup fails.
- CORS is an exact allowlist: no wildcard; non-Local origins require HTTPS.
  Local contract bearer tokens prove `401`/`403`/`200` pipeline behavior only;
  they are not Firebase token validation and must never be enabled elsewhere.
- Errors use ProblemDetails with stable `code` and `traceId`; never expose token,
  Firebase UID, Player Tag, external URL, raw payload, stack, or infrastructure
  detail.
