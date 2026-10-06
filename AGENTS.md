# CrownPilot implementation boundaries

## Current scope

This checkout is executing only SDD task `002-01-bootstrap-toolchain.md`.
Keep changes limited to reproducible toolchain, Clean Architecture project
boundaries, bootstrap health/OpenAPI, frontend build/test tooling, local test
harnesses, OCI packaging, and minimum CI.

Do not implement Firebase authentication, Google Sign-In, Supabase provisioning,
remote PostgreSQL access, domain schema/migrations, player lookup, authorization
rules, Player Tag behavior, or product identity flows here. Those belong to later
`002-*` tasks.

## Dependency boundaries

- `Domain` has no framework, HTTP, database, Firebase, Supabase, Npgsql, or
  provider dependency.
- `Application` references `Domain` and owns future ports/use cases; it does not
  reference provider SDKs.
- `Infrastructure` references `Application`/`Domain` and is the only layer that
  carries EF Core/Npgsql packages. Bootstrap registers no connection or migration.
- `Api` owns HTTP, ProblemDetails, health, OpenAPI, and composition; it may
  reference `Application` and `Infrastructure`.
- Frontend remains an independent static Vite artifact. Firebase SDK is not part
  of bootstrap.

## Official commands

Mac and Linux use the same commands from repository root:

```text
dotnet tool restore
dotnet restore CrownPilot.sln
dotnet build CrownPilot.sln --configuration Release
dotnet test CrownPilot.sln --configuration Release
npm ci --prefix frontend
npm run typecheck --prefix frontend
npm run lint --prefix frontend
npm run test:unit --prefix frontend
npm run build --prefix frontend
npm run smoke --prefix frontend
npm run openapi:check --prefix frontend
```

The scripts `test:contract`, `test:e2e`, `test:rls`, `smoke:container`, and
`db:*` are real cross-platform Node entry points. They fail with an explicit
blocked status when their later-task environment is absent; they must not be
replaced with Bash-only wrappers or empty tests.

## Validation expectations

`/health/live` is process-only and must not contact a provider. OpenAPI is
generated from Minimal API code by `Microsoft.AspNetCore.OpenApi`; no hand-written
YAML is allowed. Do not read or embed ignored `.env.local`. Never commit secrets.
