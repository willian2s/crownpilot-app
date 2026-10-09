# Execução local

Guia para executar base CrownPilot em Mac/Linux. Local usa fixture contratual de
bearer na API Go; Firebase real, banco remoto e lookup pertencem a subtarefas posteriores.

## Pré-requisitos

- Go `1.27.2` (linha `go` do `go.mod`; com `GOTOOLCHAIN=auto` um Go mais
  antigo baixa essa versão);
- .NET SDK `10.0.401`, somente para o baseline .NET até o cutover da `002-15`;
- Node.js `24.21.0`;
- npm `11.19.0`;
- Docker Desktop ou Docker Engine, para execução OCI.

Confira versões:

```text
go version
dotnet --info
node --version
npm --version
docker info
```

## Preparar clone

Na raiz do repositório:

```text
go mod download
dotnet tool restore
dotnet restore CrownPilot.sln
npm ci --prefix frontend
```

Ferramentas Go (`golangci-lint`, `oapi-codegen`) ficam fixadas em
`tools/go.mod` e rodam com `go tool -modfile=tools/go.mod <ferramenta>`; não
precisam de instalação global.

Não leia nem copie `.env.local`. Arquivo é ignorado e pode conter valores
locais sensíveis.

## Executar API sem Docker

Terminal 1:

```text
CROWNPILOT_ENVIRONMENT=Local go run ./cmd/crownpilot-api
```

API local usa `http://localhost:5080`. Sem `CROWNPILOT_ENVIRONMENT` o processo
encerra antes de abrir a porta. As demais variáveis Go estão documentadas em
`.env.example`. O baseline .NET ainda pode ser executado com
`dotnet run --project src/Api/Api.csproj` até a `002-15`.

Terminal 2, frontend:

```text
npm run dev --prefix frontend
```

Frontend Vite informa URL local no terminal, normalmente
`http://localhost:5173`.

Verifique API:

```text
curl --fail http://localhost:5080/health/live
curl --fail http://localhost:5080/openapi/v1.json
```

Endpoint protegido de contrato:

```text
curl --fail -H 'Authorization: Bearer contract-authorized-token' \
  http://localhost:5080/api/v1/bootstrap
```

`contract-authenticated-token` prova authentication sem permissão e retorna
`403`; token ausente ou inválido retorna `401`. Fixtures são aceitas somente em
Local e não representam validação Firebase.

Swagger UI:

```text
http://localhost:5080/docs
```

## Executar API com Docker

Construir imagem:

```text
docker build -t crownpilot-api:bootstrap .
```

Iniciar container:

```text
docker run --rm -d \
  --name crownpilot-api-bootstrap \
  --env CROWNPILOT_ENVIRONMENT=Local \
  -p 8080:8080 \
  crownpilot-api:bootstrap
```

A imagem Go é multi-stage, final distroless sem shell e roda como usuário
`nonroot`. Ela não define `CROWNPILOT_ENVIRONMENT`: o hosting precisa informar
o ambiente, e Preview/Production escondem OpenAPI e `/docs`.

Testar:

```text
curl --fail http://localhost:8080/health/live
curl --fail http://localhost:8080/openapi/v1.json
```

Swagger UI local no container:

```text
http://localhost:8080/docs
```

Parar:

```text
docker stop crownpilot-api-bootstrap
```

Imagem recebe configuração em runtime. Para hosting que fornece outra porta,
defina `PORT` e publique mesma porta:

```text
docker run --rm -d \
  --name crownpilot-api-bootstrap \
  --env CROWNPILOT_ENVIRONMENT=Local \
  --env PORT=8081 \
  -p 8081:8081 \
  crownpilot-api:bootstrap
```

## Validação

```text
go vet ./...
go test -race ./...
go tool -modfile=tools/go.mod golangci-lint run
go build ./cmd/crownpilot-api
npm run smoke --prefix frontend
npm run openapi:check --prefix frontend
dotnet build CrownPilot.sln --configuration Release
dotnet test CrownPilot.sln --configuration Release
npm run typecheck --prefix frontend
npm run lint --prefix frontend
npm run test:unit --prefix frontend
npm run build --prefix frontend
npm run smoke:container --prefix frontend
```

Depois de editar `api/openapi/v1.json`, regenere e versione os artefatos:
`go generate ./...` e `npm run openapi:generate --prefix frontend`. O CI falha
se eles divergirem da fonte.

`test:rls` e `test:e2e` permanecem bloqueados até subtarefas que forneçam,
respectivamente, PostgreSQL/RLS descartável e ambiente Staging.
