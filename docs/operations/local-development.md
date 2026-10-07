# Execução local

Guia para executar base CrownPilot em Mac/Linux. Local usa fixture contratual de
bearer; Firebase real, banco remoto e lookup pertencem a subtarefas posteriores.

## Pré-requisitos

- .NET SDK `10.0.401`;
- Node.js `24.21.0`;
- npm `11.19.0`;
- Docker Desktop ou Docker Engine, para execução OCI.

Confira versões:

```text
dotnet --info
node --version
npm --version
docker info
```

## Preparar clone

Na raiz do repositório:

```text
dotnet tool restore
dotnet restore CrownPilot.sln
npm ci --prefix frontend
```

Não leia nem copie `.env.local`. Arquivo é ignorado e pode conter valores
locais sensíveis.

## Executar API sem Docker

Terminal 1:

```text
dotnet run --project src/Api/Api.csproj
```

API local usa `http://localhost:5080`.

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
  --env ASPNETCORE_ENVIRONMENT=Development \
  -p 8080:8080 \
  crownpilot-api:bootstrap
```

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
  --env PORT=8081 \
  -p 8081:8081 \
  crownpilot-api:bootstrap
```

## Validação

```text
dotnet build CrownPilot.sln --configuration Release
dotnet test CrownPilot.sln --configuration Release
npm run typecheck --prefix frontend
npm run lint --prefix frontend
npm run test:unit --prefix frontend
npm run build --prefix frontend
npm run smoke:container --prefix frontend
```

`test:rls` e `test:e2e` permanecem bloqueados até subtarefas que forneçam,
respectivamente, PostgreSQL/RLS descartável e ambiente Staging.
