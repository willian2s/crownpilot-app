# 002-01 — Bootstrap do toolchain

- **Ticker:** `002`
- **Número:** `01`
- **Status:** `completed with constraints`

## Objetivo e resultado esperado

Criar baseline reproduzível a partir de clone limpo com solution .NET, API ASP.NET
Core, frontend React + TypeScript + Vite, EF Core + Npgsql, testes e container
Docker. O resultado ainda não configura login real, banco remoto ou lookup.

## Requisitos cobertos

- SDK .NET pinado e solution versionada;
- .NET 10 LTS e Node.js 24 LTS pinados, com patches atualizados conforme suporte
  oficial vigente;
- camadas `API`, `Application`, `Domain` e `Infrastructure`;
- React, TypeScript strict e Vite independentes da API;
- EF Core e Npgsql referenciados somente na Infrastructure;
- testes .NET/frontend, lint, type-check, build e health;
- Docker/OCI portátil, sem dependência estrutural de Vercel;
- comandos de desenvolvimento e CI mínimo documentados.
- documentação Mac/Linux equivalente e aprendizado incremental em `docs/learning/`.

## Escopo incluído

- criar `global.json` com .NET 10, solution e projetos `.csproj`;
- fixar Node/npm/package manager por `.nvmrc`, `packageManager` ou equivalente;
- criar manifesto local de ferramentas `dotnet-ef` e registrar versões de Firebase
  CLI, Supabase CLI e Docker exigidas pelo setup;
- criar `src/Api`, `src/Application`, `src/Domain` e `src/Infrastructure`;
- configurar referências para impedir Infrastructure no Domain e API no Domain;
- criar frontend Vite separado, com `package.json`, lockfile e TypeScript strict;
- adicionar EF Core, Npgsql e tooling de migrations sem conectar a ambiente real;
- adicionar runner de testes .NET, testes frontend e endpoint health mínimo;
- definir scripts oficiais sem dependência de shell específico: `test:unit`,
  `test:contract`, `test:rls`, `test:e2e`, `smoke`, `smoke:container`,
  `openapi:check`, `typecheck`, `lint` e `build`;
- criar `Dockerfile` multi-stage para API e documentação de execução local;
- criar `.env.example`/settings sem secrets;
- criar `AGENTS.md` local com comandos, boundaries e limites de escopo;
- deixar workflow CI mínimo para restore, build, test, type-check e frontend build.

## Escopo excluído

- Google Sign-In ou validação real de Firebase ID Token;
- provisionamento de Firebase/Supabase, schema de negócio ou secrets;
- Player Tag, lookup, autorização de usuário ou UI final;
- deploy de Staging/Production;
- uso de SDK C# do Supabase.

## Dependências

- Fase 001 e [ADR 004](../../decisions/004-aspnet-core-react-vite-firebase-postgresql.md);
- nenhum runtime existente ou código legado.

## Arquivos e símbolos prováveis

- `global.json`, `CrownPilot.sln`, `src/*/*.csproj`;
- `src/Api/Program.cs`, `HealthEndpoints`, `DependencyInjection`;
- `src/Domain/`, `src/Application/`, `src/Infrastructure/`;
- `frontend/package.json`, `vite.config.ts`, `tsconfig.json`;
- `tests/Unit`, `tests/Application`, `frontend/src/**/*.test.ts`;
- `Dockerfile`, `.dockerignore`, `.env.example`, `.github/workflows/ci.yml`.

## Passos de implementação

1. Fixar .NET 10 LTS e Node.js 24 LTS, validando patches suportados na documentação
   oficial; registrar exceção somente se compatibilidade concreta exigir.
2. Criar solution/projetos e validar referências entre camadas.
3. Configurar frontend React/Vite independente, scripts e build `dist/`.
4. Adicionar EF Core/Npgsql na Infrastructure, manifesto `dotnet-ef` local e
   tooling sem migration de domínio ainda.
5. Criar health mínimo e teste unitário por camada sem provider externo.
6. Criar imagem Docker genérica da API, com configuração em runtime e filesystem
   efêmero.
7. Documentar comandos POSIX comuns a Mac/Linux e deixar CI mínimo executável em
   Pull Request. Scripts que precisem de lógica usam Node ou .NET, não Bash-only.

## Testes e comandos de validação

```text
dotnet --info
dotnet restore
dotnet build --configuration Release
dotnet test --configuration Release
npm ci
npm run typecheck
npm run lint
npm run build
docker build -t crownpilot-api:bootstrap .
docker run --rm -d --name crownpilot-api-bootstrap -p 8080:8080 crownpilot-api:bootstrap
curl --fail http://localhost:8080/health/live
docker stop crownpilot-api-bootstrap
```

Confirmar clone limpo, ausência de secrets na imagem/bundle e que Domain não
importa ASP.NET Core, EF Core, Npgsql, Firebase ou HTTP. `dotnet build` deve
gerar/verificar OpenAPI quando a API já possuir endpoints; a fundação não deve
manter YAML OpenAPI manual. `/health/live` não testa banco ou provider.

## Definição de pronto

- solution e SDK pinado restauram em clone limpo;
- Node/npm/package manager são pinados e reproduzidos em CI;
- API e frontend compilam separadamente;
- EF Core/Npgsql existem somente na camada prevista;
- testes .NET/frontend, type-check, lint e build possuem comandos reais;
- `/health/live` responde no processo e no container;
- Docker inicia sem secret embutido e sem depender de Vercel;
- CI mínimo bloqueia falha de restore/build/test/type-check/build;
- `AGENTS.md` documenta comandos e não permite avançar Auth/persistência nesta task;
- setup Mac/Linux usa fonte de verdade versionada (`global.json`, Node pinado,
  lockfile e `dotnet-tools.json`);
- Task registra primeiro conceito de aprendizado usado, sem criar curso paralelo;
- nenhum Auth, banco remoto, lookup ou fluxo de negócio é implementado nesta task.

## Riscos e cuidados

- Não criar abstração vazia sem consumidor previsto.
- Não adicionar pacote de provider ao Domain/Application.
- Não adicionar Firebase Admin ou configuração real antes da Task 04.
- Não aplicar migrations automaticamente no startup do container.
- Não criar arquivos de teste vazios apenas para satisfazer scripts.

## Implementação e evidências

### Arquivos alterados

- `global.json`, `CrownPilot.sln`, `Directory.Build.props` e `dotnet-tools.json`;
- `src/Api`, `src/Application`, `src/Domain` e `src/Infrastructure`;
- `tests/Unit`, `tests/Application`, `tests/Architecture`, `tests/Api` e
  `tests/Contract`;
- `frontend/` com React, TypeScript strict, Vite, Vitest, ESLint, lockfile e
  scripts cross-platform;
- `Dockerfile`, `.dockerignore`, `.nvmrc`, `.env.example`, `AGENTS.md`,
  `.github/workflows/ci.yml` e `scripts/check-boundaries.mjs`;
- `docs/operations/002-01-bootstrap-toolchain.md` e
  `docs/learning/001-minimal-api-composition-root.md`.

### Decisões e desvios

- Minimal API escolhida para manter composição, health e OpenAPI code-first
  mínimos; decisão registrada em `docs/learning/001-minimal-api-composition-root.md`.
- API registra `ProblemDetails`, `/health/live`, `/api/v1/bootstrap` e OpenAPI
  gerado; não registra conexão, migration, authentication, lookup ou provider.
- Swagger UI usa o documento OpenAPI gerado, fica disponível em `/docs` somente
  em Development/Staging e não cria contrato paralelo.
- `Infrastructure` é única camada com EF Core/Npgsql; `dotnet-ef` fica no
  manifesto local em versão `10.0.12`.
- SDK .NET `10.0.401`, runtime/pacotes Microsoft `10.0.12`, Node pinado em
  `.nvmrc` `24.21.0` e npm pinado em `packageManager` `npm@11.19.0`.
- Ambiente local tinha Node `24.20.0`, compatível com engines do frontend, mas
  está um patch atrás do pin versionado. Firebase CLI observado em `15.31.0`;
  Supabase CLI e Docker não estavam instalados.
- Scripts dependentes de PostgreSQL/RLS, staging E2E e Docker retornam bloqueio
  explícito (`2`), nunca sucesso falso ou teste vazio.
- Docker restore usa somente `src/Api/Api.csproj`, porque testes ficam fora da
  imagem runtime; a API usa `PORT` válido fornecido pelo hosting e mantém `8080`
  como fallback local/container.
- `test:e2e` permanece bloqueado durante o bootstrap mesmo com URL configurada;
  preflight de health não é reportado como E2E de identidade.
- `frontend/scripts/build.mjs` e `process.mjs` foram ajustados para executar
  `tsc` e Vite no diretório frontend, preservando comandos oficiais com
  `--prefix frontend`.

### Comandos executados e resultados

Passaram:

```text
dotnet tool restore
dotnet restore CrownPilot.sln
dotnet build CrownPilot.sln --configuration Release
dotnet test CrownPilot.sln --configuration Release
node scripts/check-boundaries.mjs
npm ci --prefix frontend
npm run typecheck --prefix frontend
npm run lint --prefix frontend
npm run test:unit --prefix frontend
npm run build --prefix frontend
npm run test:contract --prefix frontend
npm run smoke --prefix frontend
npm run openapi:check --prefix frontend
PORT=5091 ASPNETCORE_ENVIRONMENT=Development dotnet run \
  --project src/Api/Api.csproj --configuration Release --no-build --no-restore \
  --no-launch-profile & api_pid=$!; \
  curl --fail http://localhost:5091/health/live; kill "$api_pid"; wait "$api_pid" || true
dotnet ef --version                 # 10.0.12
```

Resultados relevantes: build .NET sem warnings/erros; testes .NET `8` passaram;
teste frontend `1` passou; build frontend gerou `frontend/dist`; smoke local
confirmou `Healthy`; contrato confirmou `/openapi/v1.json` e
`/api/v1/bootstrap`; boundary check passou; build executou
`GenerateOpenApiDocuments` e escreveu `src/Api/obj/Api.json`; smoke manual com
`PORT=5091` confirmou `Healthy` na porta fornecida pelo runtime.

Bloqueados de forma esperada:

```text
npm run test:rls --prefix frontend        # exit 2; banco/RLS pertence a 002-03
npm run test:e2e --prefix frontend        # exit 2; staging não configurado
```

### Riscos residuais

- Supabase CLI, PostgreSQL/RLS, authentication real e E2E permanecem gates das
  subtarefas posteriores; nenhum foi simulado como concluído.

## Atualização posterior — Swagger UI

- **Status:** `completed with constraints` mantido; documentação navegável
  adicionada sem liberar authentication ou endpoints de negócio.
- **Arquivos alterados:** `src/Api/Api.csproj`, `src/Api/Program.cs`,
  `tests/Api/HealthEndpointTests.cs` e `docs/operations/local-development.md`.
- **Decisão:** `Swashbuckle.AspNetCore.SwaggerUI` `10.2.3` serve UI em `/docs`
  somente para Development/Staging e aponta para `/openapi/v1.json`; OpenAPI
  code-first continua fonte única. Desvio explícito da intenção anterior de
  Scalar, solicitado para disponibilizar Swagger UI local.
- **Validação:** `dotnet build CrownPilot.sln --configuration Release` passou
  sem warnings/erros; `dotnet test CrownPilot.sln --configuration Release`
  passou com `11` testes; smoke HTTP confirmou Swagger UI em `/docs`.
- **Risco residual:** UI permanece desabilitada fora de Development/Staging;
  OAuth, try-it-out autenticado e endpoints de negócio ainda não existem.

### Validação Docker posterior

- **Comando:** `docker info && npm run smoke:container --prefix frontend`.
- **Primeira tentativa:** bloqueada por `commit failed: input/output error` em
  `WORKDIR /src`; Docker Desktop estava com metadata read-only.
- **Reexecução:** Docker `29.8.2` restaurado; build multi-stage, restore, publish
  e execução passaram; `Container smoke passed: /health/live`.
- **Evidência:** imagem `crownpilot-api:bootstrap` criada e container smoke
  encerrado pelo script sem erro.
