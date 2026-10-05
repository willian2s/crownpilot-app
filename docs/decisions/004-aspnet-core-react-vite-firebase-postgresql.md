# ADR 004 — ASP.NET Core, React/Vite, Firebase Auth e PostgreSQL/Supabase

- **Status:** accepted
- **Data:** 2026-10-05
- **Escopo:** stack, arquitetura, identidade, persistência, API, testes, ambientes e deploy
- **Fase:** 002 — Fundação da aplicação e identidade persistente
- **Substitui:** ADRs 001, 002 e 003 nas partes afetadas

## Contexto

O projeto ainda não possui runtime, dados legados ou deploy implementado. A
Fase 001 liberou bootstrap, identidade CrownPilot e vínculo read-only de perfil
público. A arquitetura intermediária ficou inconsistente com a direção oficial
e precisa ser substituída antes da implementação.

## Decisão

- Frontend: React + TypeScript + Vite, artefato estático independente.
- Backend: C# + ASP.NET Core Web API.
- Arquitetura: Clean Architecture pragmática dentro de um Modular Monolith.
- Camadas: `API`, `Application`, `Domain` e `Infrastructure`.
- Persistência: PostgreSQL hospedado no Supabase, acessado por EF Core + Npgsql.
- Auth: Firebase Authentication com Google Sign-In.
- API: REST JSON em `/api/v1`, OpenAPI code-first e ProblemDetails.
- Infraestrutura: Docker/OCI; Render é o hosting inicial da API e alvo preferido
  do frontend estático. Vercel permanece alternativa estática.

```text
React/Vite
  -> Firebase Authentication / Google Sign-In
  -> Firebase ID Token
  -> ASP.NET Core API
  -> Application / Domain / Infrastructure
  -> EF Core + Npgsql
  -> PostgreSQL hospedado no Supabase
```

### Clean Architecture pragmática e Modular Monolith

O backend começa como um processo, um deploy e uma imagem OCI. `Identity` e
`PlayerLink` são módulos funcionais internos; integrações externas são adapters
técnicos, não serviços distribuídos. Módulos compartilham somente contratos
explícitos e não acessam entidades ou repositories de outro módulo diretamente.

```text
Domain          -> nenhuma camada externa
Application     -> Domain
Infrastructure  -> Application + Domain
API             -> Application
API composition root -> Infrastructure
Frontend        -> contrato HTTP/OpenAPI
```

Não introduzir microservices, Kafka, RabbitMQ, service mesh, Kubernetes, API
Gateway adicional, worker ou queue sem necessidade operacional demonstrada.
Extração futura permanece possível por boundary interno, não por distribuição
prematura.

### Boundaries

`Domain` contém entidades e invariantes. Não conhece Firebase, Supabase, EF
Core, Npgsql, HTTP, ASP.NET Core, Render, Docker ou Vercel.

`Application` orquestra casos de uso e depende de ports/interfaces. Não importa
SDKs de providers.

`Infrastructure` implementa adapters para Firebase, EF Core/Npgsql, PostgreSQL,
Clash Royale API e observabilidade.

`API` concentra HTTP, DTOs, ProblemDetails, autenticação, CORS, endpoints e
composição. Authorization de recurso pertence aos casos de uso/policies do
backend, sem confiar no frontend.

Authentication responde “quem é o usuário?”. Authorization responde “o que ele
pode fazer?” e permanece responsabilidade do backend.

### Authentication

React usa Firebase Web SDK para login, refresh e logout. Cada request envia
`Authorization: Bearer <Firebase ID Token>` por HTTPS. ASP.NET Core valida
assinatura, issuer, audience/project ID, expiração, claims temporais, `sub`,
`kid` e rotação de chaves usando Firebase Admin SDK para .NET atrás de adapter
na Infrastructure. Não haverá validação JWT criptográfica manual.

Não haverá troca por cookie nem dois esquemas simultâneos nesta fase. UID só é
aceito de claim verificada; body, query e headers arbitrários são ignorados.
O authentication handler não cria usuário nem escreve no banco: produz somente
`AuthenticatedSubject`. Application resolve/cria a identidade interna apenas no
caso de uso que precisa dela. Falha PostgreSQL após token válido não vira `401`.

### Identidade e persistência

Firebase UID é subject externo. O domínio usa `CrownPilotUserId` interno.

```text
Firebase UID -> CrownPilot User -> Player Profile -> Player Tag
```

Schema mínimo:

```text
crownpilot_users(id uuid, firebase_uid text unique, timestamps)
primary_player_links(user_id uuid FK, player_tag text,
  subject_type = public_profile, ownership_status = unverified)
```

EF Core é dono de tabelas, colunas, constraints, índices e migrations da
aplicação. Npgsql é provider PostgreSQL. SQL separado cobre somente RLS,
grants, roles, extensões e objetos de plataforma, sem duplicar schema.

Ordem operacional: migrations EF Core, SQL de segurança/RLS, fixtures. Migrations
de produção rodam em job controlado, não no startup de múltiplas réplicas.

Dados CrownPilot ficam em schema PostgreSQL dedicado, não exposto pela Supabase
Data API. RLS usa role de runtime que não é owner e não possui `BYPASSRLS`.
Após validar o token, uma transação explícita define com parâmetros e `SET LOCAL`
o Firebase UID verificado para resolver/criar `crownpilot_users`; depois define
`app.crownpilot_user_id` para acessar `primary_player_links`. Policies negam sem
contexto. Commit/rollback encerra ambos os contextos; retry repete toda a
transação. Conexões de migration/admin e runtime são separadas.

Compatibilidade com conexão direta/session pooler e, somente se necessário,
transaction pooler será provada em Staging com Npgsql. Migration usa conexão
direta. Nunca trocar `SET LOCAL` por estado de sessão persistente nem desabilitar
RLS silenciosamente. RLS é defesa adicional; authorization ASP.NET Core continua
obrigatória.

`primary_player_links.version` é token de concorrência otimista no JSON. Replace
exige `expectedVersion` lida anteriormente e retorna `409` em mismatch; não há
ETag paralelo nem last-write-wins silencioso. `EnsureCrownPilotUser` executa no
primeiro vínculo, cria/recupera por `firebase_uid` de forma idempotente e relê
após unique conflict. Leituras e exclusões não criam usuário como efeito colateral.

Supabase é infraestrutura de banco, não backend da aplicação. SDK C# do Supabase
não é dependência central.

### Player Tag

Player Tag é referência pública/read-only. Vínculo inicial usa
`public_profile` e `unverified`. Não armazenar credenciais Supercell, snapshot,
coleção, Arena, battle history, cache persistente ou payload raw.

Lookup fica atrás de port `IClashRoyaleClient`, implementado na Infrastructure.
Somente resultado `resolved` cria/substitui vínculo; falha externa preserva o
vínculo anterior.

### API, OpenAPI e erros

REST JSON usa `/api/v1`. Contrato mínimo:

```text
GET    /api/v1/me/player-link
PUT    /api/v1/me/player-link
DELETE /api/v1/me/player-link
DELETE /api/v1/me
GET    /openapi/v1.json
GET    /docs                    # Local e Staging controlado
```

OpenAPI é gerado do código por `Microsoft.AspNetCore.OpenApi`; Scalar consome o
mesmo documento para navegação, sem segundo gerador. Geração em build e contract
tests detectam drift. UI fica desabilitada em Production; política de exposição
do JSON é explícita por ambiente.

Erros usam `IProblemDetailsService`, `ProblemDetails` e
`ValidationProblemDetails`, com `traceId` e código estável sem stack trace,
secret, tag, UID, URL externa ou detalhe interno. Cobrir `400`, `401`, `403`,
`404`, `409`, `422` somente quando semântica exigir, `429`, `500` e `503` para
dependência externa indisponível. Authentication inválida retorna `401`;
identidade válida sem permissão retorna `403`. API atual não recebe user ID e,
para recurso alheio futuro, responde `404` para evitar enumeração.

Exclusão de dados CrownPilot exige `auth_time` presente e dentro de 5 minutos do
relógio do servidor, com tolerância máxima de 60 segundos. Claim ausente,
malformada ou fora da janela retorna `403` com
`reauthentication_required`. Exclusão não remove conta Firebase/Google.

### Ambientes e deploy

| Ambiente | Auth | Banco | Regra |
|---|---|---|---|
| Local | Emulator/fixtures | Supabase CLI/Docker descartável | sem produção |
| PR Preview | sem login real | sem produção | hostname efêmero |
| Staging | Firebase separado, Google real | Supabase separado | hostname fixo, E2E |
| Production | projeto próprio | banco próprio | promoção controlada via `main` |

Preview e Staging são problemas diferentes. Preview valida build/UI/smoke sem
URL fixa ou autenticação real; Staging valida integração completa.

Frontend gera `dist` independente. Render Static Site é alvo inicial preferido;
Vercel, CDN, Nginx ou storage estático permanecem alternativas. API roda no
Render Web Service a partir de imagem Docker/OCI ASP.NET Core, com secrets em
runtime, filesystem efêmero, `$PORT`, `/health/live` e `/health/ready`.

Render é hosting inicial por simplicidade, custo e suporte a containers, não
dependência da aplicação. Mesma imagem imutável, identificada por digest, deve
rodar localmente e poder migrar para Azure App Service, Azure Container Apps,
AWS, GCP ou outro runtime. Kubernetes é possibilidade futura, não requisito.
Limites e planos do Render são consultados antes de provisionar; números não são
congelados nesta ADR.

EF migrations são publicadas como bundle/job one-shot separado da imagem de
runtime e aplicadas antes do SQL versionado de segurança. Produção promove o
mesmo digest validado em Staging; rollback volta ao digest anterior e não executa
`Down` destrutivo automaticamente.

### Qualidade, observabilidade e aprendizado

Testes fazem parte da fundação: unit/Application, architecture, integration com
PostgreSQL, persistence/RLS, authentication/authorization, contract por fixtures,
API/OpenAPI, frontend, E2E em Staging e smoke. CI normal não chama Clash Royale
live. PR produz apenas build verificável; staging candidate constrói/publica um
digest OCI imutável, executa E2E/smoke, `main` promove esse mesmo digest sem
rebuild e Production recebe-o somente após aprovação e smoke. Falha de migration,
RLS, OpenAPI, imagem ou smoke bloqueia promoção.

Structured logging, `Activity`/request ID, métricas de baixa cardinalidade,
redaction e health checks são obrigatórios. Liveness verifica o processo;
readiness verifica configuração e PostgreSQL, nunca Clash Royale. Tokens, service
account, senha, Firebase UID, Player Tag, URL com tag e payload externo não entram
em logs ou labels.

Código futuro será didático sem abstrações artificiais. Primeira ocorrência de
DI, middleware, authentication/authorization, EF Core, migrations, Npgsql,
async/await, `CancellationToken`, options e lifecycle explica intenção e boundary.
Comentários óbvios são evitados. `docs/learning/` cresce apenas a partir de
conceitos efetivamente usados.

## Motivação

- portabilidade de runtime e hospedagem;
- separação clara de responsabilidades;
- menor acoplamento a providers;
- PostgreSQL tratado como PostgreSQL;
- evolução arquitetural adequada ao longo prazo;
- autorização crítica centralizada no backend.
- aprendizado real de C#/.NET sem enfraquecer design ou segurança.

## Trade-offs

- maior número de componentes e providers;
- integração Firebase ↔ ASP.NET Core;
- gestão de migrations EF Core;
- RLS sem pressupor `auth.uid()` de outro provider;
- complexidade operacional maior que usar BaaS como backend completo.

## Alternativas descartadas

- backend server-driven acoplado ao frontend;
- Supabase Auth como identidade da aplicação;
- SDK C# do Supabase como acesso principal;
- Firestore;
- cookies e bearer simultâneos;
- backend obrigatório em Vercel;
- Render-native sem imagem OCI portátil;
- EF Core e SQL duplicando schema;
- Firebase UID como chave primária de domínio;
- microservices, brokers, Kubernetes e gateway sem requisito;
- GraphQL ou gRPC nesta fase;
- YAML OpenAPI manual ou dois geradores concorrentes.

## Gates e revisão

- Domain e Application não importam SDK, HTTP ou framework;
- Task 01 usa .NET 10 LTS e Node.js 24 LTS pinados, produz API, frontend, EF
  Core, Npgsql, testes e Docker sem Auth real;
- tokens Firebase têm validação por ambiente e testes de rotação;
- migrations EF e SQL RLS não duplicam schema;
- PR não usa Firebase/Supabase production nem API Clash Royale live;
- OpenAPI/ProblemDetails e imagem OCI possuem contract/smoke gates;
- frontend pode sair do Render/Vercel e API pode iniciar fora do Render;
- Fase 003 permanece bloqueada pelos gates de API data, retenção, ownership,
  egress, meta e compliance da Fase 001.

Não há dados legados ou dual-write. Mudanças incompatíveis usam expand/contract.
Reabrir ADR se Firebase, PostgreSQL, portabilidade ou limites operacionais
invalidarem a decisão.

## Fontes oficiais consultadas

Consulta em 2026-10-05; versões e limites devem ser revalidados na Task 002-01 ou
antes de provisionar:

- [.NET support policy](https://dotnet.microsoft.com/en-us/platform/support/policy/dotnet-core);
- [ASP.NET Core OpenAPI](https://learn.microsoft.com/en-us/aspnet/core/fundamentals/openapi/overview?view=aspnetcore-10.0) e [UI do documento gerado](https://learn.microsoft.com/en-us/aspnet/core/fundamentals/openapi/using-openapi-documents?view=aspnetcore-10.0);
- [ASP.NET Core ProblemDetails](https://learn.microsoft.com/en-us/aspnet/core/fundamentals/error-handling-api?view=aspnetcore-10.0) e [health checks](https://learn.microsoft.com/en-us/aspnet/core/host-and-deploy/health-checks?view=aspnetcore-10.0);
- [EF Core — applying migrations](https://learn.microsoft.com/en-us/ef/core/managing-schemas/migrations/applying);
- [Firebase Admin SDK](https://firebase.google.com/docs/admin/setup) e [verify ID tokens](https://firebase.google.com/docs/auth/admin/verify-id-tokens);
- [Supabase — connect to Postgres](https://supabase.com/docs/guides/database/connecting-to-postgres);
- [Npgsql EF Core provider](https://www.npgsql.org/efcore/);
- [Render Docker](https://render.com/docs/docker) e [planos/limites atuais](https://render.com/docs/free);
- [Node.js releases](https://nodejs.org/en/about/previous-releases).
