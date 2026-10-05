# 002-02 — Fixar boundaries e ambientes

- **Ticker:** `002`
- **Número:** `02`
- **Status:** `pending`

## Objetivo e resultado esperado

Fixar contratos de camadas, authentication/authorization, API JSON e matriz de
Local, Preview, Staging e Production antes de provisionar dados reais. O
resultado é um boundary explícito entre React/Vite, Firebase, API ASP.NET Core,
Application/Domain/Infrastructure, PostgreSQL/Supabase e provider externo.

## Requisitos cobertos

- stack e boundaries da [ADR 004](../../decisions/004-aspnet-core-react-vite-firebase-postgresql.md);
- bearer Firebase ID Token, sem sessão cookie nesta fase;
- authorization obrigatória no backend;
- PostgreSQL/Supabase sem SDK central do provider;
- Preview sem login real e Staging com hostname fixo;
- backend portátil fora de Vercel;
- configuração fail-closed e secrets separados.

## Escopo incluído

- documentar dependências permitidas por camada;
- definir DTOs JSON, ProblemDetails, status (`401`, `403`, `404`, `409`) e CORS por ambiente;
- definir fluxo `React -> Firebase -> Bearer -> ASP.NET Core -> Application`;
- separar `FirebaseUid` externo de `CrownPilotUserId` interno;
- definir variáveis públicas, server-only e allowlists por ambiente;
- definir Local com emuladores/fixtures, Preview sem Auth real, Staging fixo e
  Production separado;
- definir frontend estático opcional em Vercel e API Docker/OCI portátil;
- registrar ownership de EF migrations versus SQL RLS/grants.

## Escopo excluído

- implementar middleware Firebase ou casos de uso;
- provisionar Firebase/Supabase ou criar secrets reais;
- schema, repository, lookup ou UI final;
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
3. Definir contrato HTTP JSON, ProblemDetails, CORS e respostas `401`/`403`.
4. Definir matriz de hosts, Firebase project IDs, Supabase databases e secrets.
5. Definir Preview sem hostname autorizado para login real e Staging com hostname
   fixo/Google Sign-In.
6. Definir API containerizada fora de Vercel e frontend estático relocável.
7. Definir matriz: token ausente/inválido -> `401`, token válido sem permissão ->
   `403`, recurso próprio ausente -> `404`, usuário A tentando recurso de B sem
   revelar dados -> `403` ou `404` conforme contrato anti-enumeração.
8. Registrar rollout, fail-closed e rollback sem dados compartilhados.

## Testes e comandos de validação

```text
npm run typecheck
npm run lint
```

Revisão documental deve confirmar que Domain não conhece providers/frameworks,
que nenhum UID vindo do request é confiável e que Preview não aponta para
Firebase/Supabase de Staging/Production.

## Definição de pronto

- boundaries e dependências de camadas estão documentados;
- API JSON, bearer, ProblemDetails e authorization possuem contrato;
- Firebase e Supabase têm configuração independente por ambiente;
- Preview e Staging são explicitamente diferentes;
- Vercel é opcional para frontend e não é dependência da API;
- configuração inválida falha fechado;
- ownership de migrations e ordem EF -> SQL de segurança estão definidos;
- nenhum secret real ou provisionamento foi criado.

## Riscos e cuidados

- Não tratar authentication como autorização.
- Não usar audience/issuer de outro ambiente.
- Não permitir CORS amplo por conveniência.
- Não introduzir cookie, sessão ou segundo esquema bearer nesta fase.
- Não transformar Vercel em boundary de domínio ou API.
