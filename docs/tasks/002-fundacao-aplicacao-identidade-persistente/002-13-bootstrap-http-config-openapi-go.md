# 002-13 — Bootstrap Go, HTTP, configuração e OpenAPI

- **Ticker:** `002`
- **Número:** `13`
- **Status:** `pending`

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

| Baseline .NET | Destino Go da migração |
|---|---|
| `global.json`, `CrownPilot.sln`, `Directory.Build.props`, `dotnet-tools.json` | `go.mod`, manifesto de ferramentas Go, `.go-version`/equivalente e CI. |
| `src/Api/Program.cs` | `cmd/crownpilot-api/main.go` + `internal/httpapi` + `internal/config`. |
| `BootstrapStatus.cs` | DTO e handler Go de bootstrap. |
| `RuntimeEnvironment.cs`, `RuntimeOptions*.cs` | tipos e `Load`/`Validate` em `internal/config`. |
| `ProblemDetailsContract.cs`, CORS e health | handlers/middlewares `net/http` em `internal/httpapi`. |
| `src/Api/OpenApi/*` e documento code-first | `api/openapi/v1.json`, `oapi-codegen`, `openapi-typescript` e pacote de embed. |
| `cmd/crownpilot-openapi` planejado | não criar; a fonte JSON será embutida e servida diretamente. |
| `frontend/scripts/process.mjs`, `smoke.mjs`, `openapi-check.mjs` | mesmos comandos e contratos, iniciando/validando Go. |
| `Dockerfile` e workflow mínimo | build multi-stage Go, sem `.env`, secrets ou testes na imagem. |

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
