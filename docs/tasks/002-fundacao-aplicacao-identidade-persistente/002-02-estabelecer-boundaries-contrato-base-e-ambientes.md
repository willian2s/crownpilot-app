# 002-02 — Estabelecer boundaries, contrato base e ambientes

- **Ticker:** `002`
- **Número:** `02`
- **Status:** `completed`

## Objetivo e resultado esperado

Fixar Clean Architecture pragmática, boundaries do Modular Monolith,
authentication/authorization, contrato HTTP base e matriz de
Local, Preview, Staging e Production antes de provisionar dados reais. O
resultado é um boundary explícito entre React/Vite, Firebase, API ASP.NET Core,
Application/Domain/Infrastructure, PostgreSQL/Supabase e provider externo.

## Requisitos cobertos

- stack e boundaries da [ADR 004](../../decisions/004-aspnet-core-react-vite-firebase-postgresql.md);
- bearer Firebase ID Token, sem sessão cookie nesta fase;
- API REST `/api/v1` com pipeline OpenAPI gerado por `Microsoft.AspNetCore.OpenApi`;
- Swagger UI consumindo o mesmo `/openapi/v1.json`, somente Local/Staging;
- authorization obrigatória no backend;
- PostgreSQL/Supabase sem SDK central do provider;
- Preview sem login real e Staging com hostname fixo;
- backend portátil fora de Vercel;
- configuração fail-closed e secrets separados.

## Escopo incluído

- documentar dependências permitidas por camada;
- definir DTOs JSON base, ProblemDetails, códigos estáveis, status (`400`, `401`, `403`,
  `404`, `409`, `422` quando necessário, `429`, `500`, `503`) e CORS por ambiente;
- registrar códigos `invalid_player_tag`, `player_not_found`,
  `provider_rate_limited`, `provider_unavailable`, `version_conflict` e
  `reauthentication_required` sem dados sensíveis;
- fixar prefixo `/api/v1`, convenções de `/openapi/v1.json` e `/docs`; contratos
  completos de endpoint ficam na Task `002-07-implementar-casos-de-uso-e-api-v1.md`;
- fixar `version`/`expectedVersion` como concorrência JSON, sem ETag paralelo;
- fixar OpenAPI: JSON `200` em Local/Staging, `404` em Preview/Production por
  padrão; UI `/docs` somente Local/Staging;
- definir fluxo `React -> Firebase -> Bearer -> ASP.NET Core -> Application`;
- separar `FirebaseUid` externo de `CrownPilotUserId` interno;
- definir variáveis públicas, server-only e allowlists por ambiente;
- definir Local com configuração isolada, Preview sem Auth real, Staging fixo e
  Production separado; Emulator/fixtures de token pertencem à Task
  `002-04-implementar-google-sign-in-e-firebase-bearer.md`;
- definir frontend estático opcional em Vercel e API Docker/OCI portátil;
- exigir `auth_time` presente, janela de 5 minutos e tolerância de relógio de 60
  segundos para exclusão de dados CrownPilot, sem apagar Firebase;
- registrar ownership de EF migrations versus SQL RLS/grants.

## Escopo excluído

- implementar middleware Firebase ou casos de uso;
- provisionar Firebase/Supabase ou criar secrets reais;
- schema, repository, lookup, endpoints completos ou UI final;
- escolha de egress definitivo;
- liberar sync, retenção, billing ou polling.

## Dependências

- `002-01-bootstrap-toolchain.md`;
- [ADR 004](../../decisions/004-aspnet-core-react-vite-firebase-postgresql.md);
- veredito da Fase 001.

## Arquivos e símbolos prováveis

- `docs/decisions/004-aspnet-core-react-vite-firebase-postgresql.md`;
- `src/Api/Program.cs`, `CorsPolicy`, `ProblemDetailsMapping`;
- `src/Application/Abstractions/`, `AuthContext`;
- `appsettings*.json`, `.env.example`, `EnvironmentName`, `RuntimeConfig`;
- documentação de ambientes e release.

## Passos de implementação

1. Desenhar fluxo bearer e separar authentication de authorization.
2. Registrar dependências permitidas e proibidas em cada camada.
3. Definir contrato HTTP base, ProblemDetails, pipeline OpenAPI, CORS e respostas
   de authentication/authorization.
4. Definir matriz de hosts, Firebase project IDs, Supabase databases e secrets.
5. Definir Preview sem hostname autorizado para login real e Staging com hostname
   fixo/Google Sign-In.
6. Definir API containerizada no Render inicialmente, fora de Render/Vercel no
   domínio, e frontend estático relocável.
7. Definir matriz: token ausente/inválido -> `401`, token válido sem permissão ->
    `403`, recurso próprio ausente -> `404`, conflito -> `409`, provider
    indisponível -> `503`; usuário A não enumera dados de B.
8. Registrar rollout, fail-closed e rollback sem dados compartilhados.

## Testes e comandos de validação

```text
npm run typecheck --prefix frontend
npm run lint --prefix frontend
```

Revisão documental e teste arquitetural inicial devem confirmar `API → Application
→ Domain` e `Infrastructure → Application/Domain`, que Domain não conhece
providers/frameworks, que nenhum UID vindo do request é confiável e que Preview não aponta para
Firebase/Supabase de Staging/Production.

## Definição de pronto

- boundaries e dependências de camadas estão documentados;
- API JSON, bearer, ProblemDetails e authorization possuem contrato;
- contrato HTTP base e pipeline OpenAPI estão definidos; documentação completa de
  endpoints, bodies, responses, erros e headers será verificada na Task
  `002-07-implementar-casos-de-uso-e-api-v1.md`;
- paths e semântica de `version`/`expectedVersion` estão fixados;
- Firebase e Supabase têm configuração independente por ambiente;
- Preview e Staging são explicitamente diferentes;
- Vercel é opcional para frontend e não é dependência da API;
- configuração inválida falha fechado;
- Local/Staging expõem UI OpenAPI; Production segue política restrita;
- ownership de migrations e ordem EF -> SQL de segurança estão definidos;
- nenhum secret real ou provisionamento foi criado.

## Riscos e cuidados

- Não tratar authentication como autorização.
- Não usar audience/issuer de outro ambiente.
- Não permitir CORS amplo por conveniência.
- Não introduzir cookie, sessão ou segundo esquema bearer nesta fase.
- Não transformar Vercel em boundary de domínio ou API.

## Implementação e evidências

### Arquivos alterados

- `.env.example`;
- `src/Api/Program.cs`, `src/Api/Authentication/`, `src/Api/Configuration/`,
  `src/Api/OpenApi/`, `src/Api/ProblemDetails/` e configurações
  `appsettings.Development.json`, `appsettings.Preview.json`,
  `appsettings.Staging.json` e `appsettings.Production.json`;
- `src/Application/Identity/AuthContracts.cs` e
  `src/Domain/Identity/CrownPilotUserId.cs`;
- `tests/Api/HealthEndpointTests.cs`, `tests/Api/RuntimeOptionsTests.cs` e
  `tests/Contract/OpenApiContractTests.cs`;
- `docs/operations/002-02-boundaries-environments.md` e
  `docs/operations/local-development.md`;
- `docs/learning/001-minimal-api-composition-root.md` e `AGENTS.md`.

### Decisões e desvios

- Authentication e authorization foram separadas no pipeline. A API usa bearer
  contratual somente em Local, com fixtures determinísticas para provar `401`,
  `403` e acesso autorizado; Firebase real permanece na Task `002-04`.
- `RuntimeOptions` resolve `Development` como `Local`, exige correspondência
  entre `CrownPilot:Environment` e `ASPNETCORE_ENVIRONMENT`, rejeita CORS
  curinga e falha no startup para exposição OpenAPI proibida ou configuração
  inválida. Modes de provider futuros permanecem compatíveis; somente fixture
  contratual é restrita a Local.
- OpenAPI e Swagger UI continuam code-first e compartilham `/openapi/v1.json`;
  Local/Staging expõem documentação e Preview/Production respondem `404`.
- Contratos `AuthenticatedSubject`/`AuthContext` e `CrownPilotUserId` registram
  separação entre Firebase UID externo e identidade interna sem criar usuário,
  provider real, schema ou migration.
- O contrato completo de endpoints de Player Link não foi antecipado; somente
  prefixo, ProblemDetails, códigos/status, bearer e concorrência
  `version`/`expectedVersion` foram fixados.

### Comandos executados e resultados

Passaram:

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

Evidências: build .NET passou com `0` warnings e `0` errors; testes .NET
passaram com `20` testes; testes frontend passaram com `1` teste; typecheck,
lint, build, smoke local e verificação OpenAPI passaram. Testes API cobrem
OpenAPI/UI por ambiente, bearer `401`/`403`/`200`, CORS allowlist,
ProblemDetails e validação de opções. A primeira execução após atualizar o
teste de ProblemDetails encontrou título legado `Not Found`; a asserção foi
ajustada para o código estável `resource_not_found` e a suíte final passou.

### Riscos residuais

- Firebase ID Token real, Emulator/fixtures oficiais, `auth_time`, issuer,
  audience, assinatura, rotação de `kid` e autorização de casos de uso ainda
  pertencem à Task `002-04`/`002-07`.
- Firebase/Supabase não foram provisionados; origins fixos de Staging/Production
  e secrets devem ser fornecidos por runtime antes desses ambientes.
- CORS sem origins em Preview/Staging/Production falha fechado até allowlist ser
  configurada. Não há dados ou secrets reais versionados.
- Staging/Production dependem de TLS termination e forwarded headers confiáveis
  no ingress; essa prova operacional fica para validação de ambiente posterior.
- A policy do endpoint bootstrap exige autenticação; fora da fixture Local ela
  não inventa permission claims. Policies de recurso específicas serão definidas
  quando provider e casos de uso forem implementados.
