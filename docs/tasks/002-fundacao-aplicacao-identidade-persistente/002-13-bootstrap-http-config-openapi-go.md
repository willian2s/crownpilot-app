# 002-13 — Bootstrap Go, HTTP, configuração e OpenAPI

- **Ticker:** `002`
- **Número:** `13`
- **Status:** `completed`

## Objetivo e resultado esperado

Substituir o bootstrap executável .NET por um processo Go mínimo, sem ainda
implementar autenticação Firebase real ou persistência de domínio. A task entrega
toolchain, composition root, HTTP, configuração fail-closed e contrato OpenAPI
spec-first, com gate verde independente.

`002-01` a `002-04` continuam concluídas como histórico do baseline .NET; esta
task inicia migração, não reabre essas tarefas.

## Requisitos cobertos

- Go `1.27.2` pinado, `go.mod`, comandos locais/CI, lint, testes, race e imagem OCI;
- layout `cmd/` e `internal/` do Modular Monolith;
- `net/http` com `ServeMux` e padrões de método/path de Go 1.22+;
- middlewares como funções que envolvem `http.Handler`;
- configuração tipada, ambientes Local/Preview/Staging/Production e startup
  fail-closed antes do listener;
- `/health/live`, bootstrap, ProblemDetails, CORS e porta de runtime;
- `api/openapi/v1.json` como fonte única spec-first;
- contrato planejado para `DELETE /api/v1/me` documentando `503` com
  `code: authentication_unavailable` em ProblemDetails genérico, sem detalhe do
  provider;
- `oapi-codegen` para Go e `openapi-typescript` para frontend;
- artefatos gerados versionados e CI bloqueando drift;
- `/openapi/v1.json` servido diretamente da fonte embutida, sem
  `cmd/crownpilot-openapi`.

## Escopo incluído

- criar `cmd/crownpilot-api`, `internal/httpapi`, `internal/config`,
  `internal/observability` e pacote de embed sob `api/openapi`;
- transportar health, bootstrap, ProblemDetails, CORS e regras de exposição
  OpenAPI da spec 002;
- manter `/docs` somente nos ambientes permitidos, consumindo
  `/openapi/v1.json`;
- definir `api/openapi/v1.json`, geração pinada de Go/frontend e verificação de
  drift no CI;
- adaptar scripts `smoke`, `openapi:check`, desenvolvimento local e Docker para
  o processo Go;
- registrar mapa de equivalência de bootstrap, HTTP, configuração, scripts e
  documentação na seção própria desta task.

## Escopo excluído

- Firebase Admin Go e validação de bearer, pertencentes à `002-14`;
- schema, repositories, RLS, pooler, cutover final e remoção completa do .NET,
  pertencentes à `002-15`;
- casos de uso de Player Link, lookup live ou mudança de requisito funcional da
  spec 002.

## Dependências

- `002-01` a `002-04` concluídas como baseline histórico;
- [ADR 005](../../decisions/005-go-react-vite-firebase-postgresql.md);
- [spec 002](../../specs/002-fundacao-aplicacao-identidade-persistente.md);
- ADR 005 aprovada;
- perguntas pendentes fechadas ou adiadas formalmente dentro da ADR aprovada;
  aprovação da ADR é pré-condição obrigatória antes de iniciar qualquer
  implementação;
- `002-01` como evidência da porta local .NET `5080`; smoke pode usar override
  `5089`.

## Mapa de equivalência desta task

| Baseline .NET | Destino Go implementado |
|---|---|
| `global.json`, `CrownPilot.sln`, `Directory.Build.props`, `dotnet-tools.json` | `go.mod` (`go 1.27.2`, sem `.go-version`), `tools/go.mod` com `golangci-lint` e `oapi-codegen` via diretiva `tool`, `.golangci.yml` e CI. |
| `src/Api/Program.cs` | `cmd/crownpilot-api/main.go` (`run`, bind antes do `Serve`, shutdown gracioso) + `internal/httpapi/server.go` (`NewHandler`) + `internal/config`. |
| `BootstrapStatus.cs` | `BootstrapStatus` gerado em `internal/httpapi/api.gen.go`; handler em `server.go`. |
| `RuntimeEnvironment.cs`, `RuntimeOptions*.cs` | `internal/config/config.go` (`Environment`, `Config`, `Load`, `Validate`). |
| `ProblemDetailsContract.cs`, CORS e health | `internal/httpapi/problem.go`, `middleware.go`, `cors.go`, `auth.go` e `/health/live`; `internal/observability` para slog e request ID. |
| `ContractAuthentication.cs` (fixture Local) | `internal/httpapi/auth.go`, somente fixture Local protegida por config; Firebase real fica na `002-14`. |
| `src/Api/OpenApi/*` e documento code-first | `api/openapi/v1.json`, `api/openapi/openapi.go` (`go:embed`), `oapi-codegen` (models) e `openapi-typescript` (`frontend/src/api/schema.gen.ts`). |
| `cmd/crownpilot-openapi` planejado | não criado; a fonte JSON embutida é servida diretamente. |
| `frontend/scripts/process.mjs`, `smoke.mjs`, `openapi-check.mjs` | mesmos comandos; compilam o binário Go em diretório temporário e validam o contrato servido contra a fonte. |
| `Dockerfile` e workflow mínimo | `Dockerfile` multi-stage `golang:1.27.2-trixie` → `distroless/static-debian13:nonroot`; CI com Go, lint, race e drift, mantendo o .NET até a `002-15`. |

## Passos de implementação

1. Fixar Go `1.27.2` e ferramentas, criar módulo e composição mínima sem provider.
2. Implementar loader/validator de configuração antes de abrir listener.
3. Registrar `ServeMux`, handlers, middlewares `http.Handler`, CORS, health,
   bootstrap e ProblemDetails.
4. Criar fonte OpenAPI JSON e gerar os artefatos Go/frontend.
5. Embutir a fonte canônica com `go:embed` e servir `/openapi/v1.json` sem
   conversão ou binário próprio.
6. Adaptar scripts, Docker e CI; manter contratos públicos de porta e smoke.
7. Executar gate verde desta task e registrar evidências antes de iniciar
   `002-14`.

## Gate verde obrigatório

O gate deve passar sem Firebase real, PostgreSQL ou chamadas externas:

```text
go test ./...
go test -race ./...
go vet ./...
golangci-lint run
go build ./cmd/crownpilot-api
npm run typecheck --prefix frontend
npm run lint --prefix frontend
npm run test:unit --prefix frontend
npm run build --prefix frontend
npm run openapi:check --prefix frontend
npm run smoke --prefix frontend
npm run smoke:container --prefix frontend
```

Bloqueios de Docker ou dependências indisponíveis devem ser registrados como
bloqueio real, nunca convertidos em sucesso falso.

## Definição de pronto

- clone limpo restaura e compila o processo Go;
- porta local canônica permanece `5080`, com `5089` apenas como override de smoke;
- `ServeMux` é único router e `chi` não aparece como dependência;
- `identity` e `playerlink` não importam HTTP, providers ou banco; a regra entre
  módulos fica pronta para `depguard`;
- configuração inválida encerra antes do listener e secrets não entram em logs;
- liveness é process-only e `/docs`/OpenAPI respeitam ambiente da spec 002;
- fonte JSON gera código Go e tipos frontend determinísticos;
- `DELETE /api/v1/me` documenta `authentication_unavailable` para falha fechada
  de comunicação com Firebase durante revogação sensível, sem reutilizar
  `provider_unavailable`;
- artefatos gerados estão no Git e CI falha em drift;
- `cmd/crownpilot-openapi` não existe;
- frontend, smoke, contrato e imagem usam backend Go sem alterar contrato
  funcional da spec 002.

## Riscos e cuidados

- Não portar requisitos de autenticação ou persistência para esta task.
- Não manter JSON servido separado da fonte canônica.
- Não abrir listener antes de validar ambiente e configuração não secreta.
- Não adicionar abstração de router ou framework sem necessidade concreta.

## Implementação e evidências

- **Status:** `completed` em 2026-10-09.
- **Arquivos criados:** `go.mod`, `tools/go.mod`, `tools/go.sum`,
  `.golangci.yml`, `cmd/crownpilot-api/main.go`,
  `cmd/crownpilot-api/main_test.go`, `internal/config/config.go`,
  `internal/config/load_test.go`, `internal/config/validate_test.go`,
  `internal/httpapi/{server,auth,cors,docs,middleware,problem,generate}.go`,
  `internal/httpapi/api.gen.go`, `internal/httpapi/{server,cors,openapi}_test.go`,
  `internal/observability/observability.go`, `api/openapi/v1.json`,
  `api/openapi/openapi.go`, `api/openapi/oapi-codegen.yaml`,
  `frontend/src/api/schema.gen.ts`.
- **Arquivos alterados:** `Dockerfile`, `.dockerignore`, `.gitignore`,
  `.env.example`, `.github/workflows/ci.yml`, `AGENTS.md`,
  `docs/operations/local-development.md`, `frontend/package.json`,
  `frontend/package-lock.json`, `frontend/scripts/process.mjs`,
  `frontend/scripts/openapi-check.mjs`, `frontend/scripts/container-smoke.mjs`.
- **Decisões:**
  - ferramentas Go em módulo separado (`tools/`) para não incluir as
    dependências do `golangci-lint` no grafo da aplicação; `go list -m all`
    lista somente o módulo principal e `chi` não aparece;
  - `depguard` bloqueia frameworks de router e mantém `identity`/`playerlink`
    livres de HTTP, config, adapters e da importação `identity -> playerlink`;
    regra provada com arquivo temporário importando `net/http`;
  - `CROWNPILOT_ENVIRONMENT` obrigatório, com nomes exatos e sem alias
    `Development`; `Load` aplica defaults por ambiente e sempre termina em
    `Validate`, que acumula todas as violações com `errors.Join`;
  - origens CORS aceitas somente na forma serializada exata do header `Origin`
    (minúsculas, sem barra final), comparadas por igualdade;
  - `traceId` gerado pelo servidor (`crypto/rand.Text`), nunca lido do cliente;
  - `oapi-codegen` gera somente models; `ProblemCode` gerado tipa
    `writeProblem`; `TestRoutesMatchContract` liga spec e `ServeMux`, e
    `x-crownpilot-planned` marca `DELETE /api/v1/me` até a implementação;
  - `/openapi/v1.json` serve `api/openapi/v1.json` byte a byte; `/docs` usa
    Swagger UI `5.33.0` do jsDelivr com SRI, CSP com hash do script inline e
    `validatorUrl: "none"`, para não enviar a URL da spec a
    `validator.swagger.io`; validação manual no navegador feita em Local;
  - imagem final distroless `nonroot`, binário estático, sem `.env`, testes ou
    tools; sem default de ambiente na imagem, que falha fechado sem ele;
  - comentários de código Go em português.
- **Desvios:**
  - `method_not_allowed` adicionado ao enum para `405` (o .NET devolvia
    `internal_error`); `defaultCode` cobre só status do transporte, então `503`
    nunca vira `provider_unavailable` implicitamente;
  - `/health/live` fora da spec, como no baseline .NET;
  - fixtures de contrato compiladas no binário e bloqueadas por
    `config.Validate` fora de Local; removê-las do binário de produção é
    critério da `002-14`;
  - `misspell` removido do lint por acusar os comentários em português.
- **Comandos executados e resultados (2026-10-09):**
  - `go mod verify` — passou.
  - `go test ./...` — passou; 101 casos incluindo subtestes.
  - `go test -race ./...` — passou.
  - `go vet ./...` — passou.
  - `go tool -modfile=tools/go.mod golangci-lint run` — 0 issues.
  - `go build ./cmd/crownpilot-api` — passou.
  - `go generate ./...`, `npm run openapi:generate --prefix frontend` e
    `git diff --exit-code` dos gerados — sem drift; geração repetida produz os
    mesmos hashes; alteração proposital da spec fez o diff falhar (exit 1).
  - `npm run typecheck --prefix frontend` — passou.
  - `npm run lint --prefix frontend` — passou.
  - `npm run test:unit --prefix frontend` — passou, 3 arquivos/8 testes.
  - `npm run build --prefix frontend` — passou.
  - `npm run openapi:check --prefix frontend` — passou; documento servido
    idêntico à fonte e `503 authentication_unavailable` documentado.
  - `npm run smoke --prefix frontend` — passou; processo filho encerrado e
    porta `5089` liberada.
  - `npm run smoke:container --prefix frontend` — passou com Docker `29.8.2`;
    container sem ambiente falha fechado, `/health/live` responde e
    `/openapi/v1.json` retorna `404` em Production; imagem com 15,5 MB.
  - `dotnet build CrownPilot.sln --configuration Release` — passou.
  - `node scripts/check-boundaries.mjs` — passou.
  - `dotnet test CrownPilot.sln --configuration Release` — Unit, Architecture,
    Persistence e Application passaram; Contract (1) e Api (20) falharam com
    `Firebase authentication configuration is invalid`. A mesma falha ocorre no
    commit `4141097`, anterior a esta task, então não é regressão: decorre da
    configuração Firebase local carregada pelo baseline .NET em `Development`.
    Os User Secrets não foram lidos por conterem credenciais. Erro conhecido
    pelo responsável e desconsiderado: o baseline .NET será removido na `002-15`.
- **Riscos residuais:**
  - os testes .NET dependentes de host falham nesta máquina por configuração
    Firebase local até a correção dos User Secrets ou a remoção do .NET na
    `002-15`; o CI não possui esses secrets;
  - `/docs` depende do jsDelivr em Local/Staging; alternativa offline é embutir
    `swagger-ui-dist` com `embed.FS`;
  - o primeiro job de CI compila `golangci-lint` a partir do código-fonte
    (cache por `tools/go.sum`).
