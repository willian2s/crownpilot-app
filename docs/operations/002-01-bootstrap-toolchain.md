# Bootstrap toolchain — 002-01

## Reproducible setup

Install Git, the pinned .NET SDK, Node.js 24 LTS, and npm. Mac and Linux use the
same repository commands; Docker Desktop (Mac) or Docker Engine/Compose (Linux)
is needed only for container/database gates. PostgreSQL is not required on host.

```text
dotnet --info
dotnet tool restore
dotnet restore CrownPilot.sln
dotnet build CrownPilot.sln --configuration Release
dotnet test CrownPilot.sln --configuration Release
npm ci --prefix frontend
npm run typecheck --prefix frontend
npm run lint --prefix frontend
npm run test:unit --prefix frontend
npm run build --prefix frontend
```

The same commands are exposed through the frontend `package.json` for bootstrap
contract checks:

```text
npm run test:contract --prefix frontend
npm run openapi:check --prefix frontend
npm run smoke --prefix frontend
npm run smoke:container --prefix frontend
npm run test:rls --prefix frontend
npm run test:e2e --prefix frontend
```

`test:rls`, `test:e2e`, and `db:*` intentionally return a blocked status until
their owning task supplies disposable PostgreSQL, RLS, or staging authentication.
They are not skipped silently. `smoke:container` returns a blocked status when
Docker is absent.

## Layer and runtime boundaries

The solution is one deployable ASP.NET Core process with four projects:

```text
Api -> Application -> Domain
Api composition root -> Infrastructure -> Application/Domain
```

Infrastructure carries EF Core, Npgsql, and the local `dotnet-ef` manifest. No
connection string, DbContext model, migration, Firebase Admin SDK, Supabase SDK,
lookup client, authentication handler, or business authorization is included in
bootstrap. Frontend is an independent `dist/` artifact.

## Version evidence recorded 2026-10-06

| Tool | Pin/observed state | Evidence |
|---|---|---|
| .NET | SDK `10.0.401`; runtime `10.0.12`; `global.json` pins SDK | `dotnet --info`; [.NET support policy](https://dotnet.microsoft.com/en-us/platform/support/policy/dotnet-core) lists .NET 10 LTS latest patch `10.0.12` on 2026-09-08 |
| Node | `.nvmrc` pins `24.21.0` | [Node.js releases](https://nodejs.org/en/about/previous-releases) lists v24.21.0 as LTS; local shell observed v24.20.0, so local patch is one behind pin |
| npm | `11.19.0` in `packageManager` and local observation | `npm --version` |
| EF CLI | `10.0.12` in `dotnet-tools.json` | `dotnet tool search dotnet-ef --take 1`; command was absent before manifest restore |
| Firebase CLI | `15.31.0` installed | `firebase --version` |
| Supabase CLI | not installed | `supabase --version` returned command not found; no version invented |
| Docker | not installed/daemon unavailable | `docker --version` returned `spawn docker ENOENT`; no version invented |
| OpenAPI/EF/Npgsql packages | Microsoft packages `10.0.12`; Npgsql EF provider `10.0.3` | NuGet package pages consulted; Npgsql `10.0.3` requires EF Core `>=10.0.4` |

The Node pin follows the current official LTS patch, while the existing local
`24.20.0` installation remains usable for local validation because package engines
accept Node 24.20+. CI reads `.nvmrc` and uses the current pin.

## Environment and secrets

`.env.example` contains only a placeholder for a later provider integration. The
bootstrap process does not load `.env.local`, and Docker ignores environment files.
No Firebase, Supabase, database, or production variables are versioned.
