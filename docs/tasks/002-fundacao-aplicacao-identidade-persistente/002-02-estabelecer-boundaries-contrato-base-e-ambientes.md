# 002-02 — Estabelecer boundaries, contrato base e ambientes

- **Ticker:** `002`
- **Número:** `02`
- **Status:** `pending`

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
- UI Scalar consumindo o mesmo `/openapi/v1.json`, somente Local/Staging;
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
  completos de endpoint ficam na Task 002-07;
- fixar `version`/`expectedVersion` como concorrência JSON, sem ETag paralelo;
- fixar OpenAPI: JSON `200` em Local/Staging, `404` em Preview/Production por
  padrão; UI `/docs` somente Local/Staging;
- definir fluxo `React -> Firebase -> Bearer -> ASP.NET Core -> Application`;
- separar `FirebaseUid` externo de `CrownPilotUserId` interno;
- definir variáveis públicas, server-only e allowlists por ambiente;
- definir Local com configuração isolada, Preview sem Auth real, Staging fixo e
  Production separado; Emulator/fixtures de token pertencem à Task 002-04;
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

- `002-01`;
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
npm run typecheck
npm run lint
```

Revisão documental e teste arquitetural inicial devem confirmar `API → Application
→ Domain` e `Infrastructure → Application/Domain`, que Domain não conhece
providers/frameworks, que nenhum UID vindo do request é confiável e que Preview não aponta para
Firebase/Supabase de Staging/Production.

## Definição de pronto

- boundaries e dependências de camadas estão documentados;
- API JSON, bearer, ProblemDetails e authorization possuem contrato;
- contrato HTTP base e pipeline OpenAPI estão definidos; documentação completa de
  endpoints, bodies, responses, erros e headers será verificada na Task 002-07;
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
