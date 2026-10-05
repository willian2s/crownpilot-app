# ADR 004 — ASP.NET Core, React/Vite, Firebase Auth e PostgreSQL/Supabase

- **Status:** accepted
- **Data:** 2026-10-05
- **Escopo:** stack oficial, boundaries, identidade, persistência, ambientes e deploy
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
- Camadas: `API`, `Application`, `Domain` e `Infrastructure`.
- Persistência: PostgreSQL hospedado no Supabase, acessado por EF Core + Npgsql.
- Auth: Firebase Authentication com Google Sign-In.
- Infraestrutura: Docker/OCI; Vercel é alvo opcional para frontend, nunca
  dependência estrutural do backend.

```text
React/Vite
  -> Firebase Authentication / Google Sign-In
  -> Firebase ID Token
  -> ASP.NET Core API
  -> Application / Domain / Infrastructure
  -> EF Core + Npgsql
  -> PostgreSQL hospedado no Supabase
```

### Boundaries

`Domain` contém entidades e invariantes. Não conhece Firebase, Supabase, EF
Core, Npgsql, HTTP, ASP.NET Core ou Vercel.

`Application` orquestra casos de uso e depende de ports/interfaces. Não importa
SDKs de providers.

`Infrastructure` implementa adapters para Firebase, EF Core/Npgsql, PostgreSQL,
Clash Royale API e observabilidade.

`API` concentra HTTP, DTOs, ProblemDetails, autenticação, autorização de
entrada, CORS, endpoints e composição da aplicação.

Authentication responde “quem é o usuário?”. Authorization responde “o que ele
pode fazer?” e permanece responsabilidade do backend.

### Authentication

React usa Firebase Web SDK para login, refresh e logout. Cada request envia
`Authorization: Bearer <Firebase ID Token>` por HTTPS. ASP.NET Core valida
assinatura, issuer, audience/project ID, expiração, claims temporais, `sub`,
`kid` e rotação de chaves usando integração suportada na Infrastructure.

Não haverá troca por cookie nem dois esquemas simultâneos nesta fase. UID só é
aceito de claim verificada; body, query e headers arbitrários são ignorados.

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

RLS usa role de runtime sem `BYPASSRLS`; após validar o Firebase token e resolver
o ID interno, o backend define `SET LOCAL app.crownpilot_user_id` dentro da
transação. Policies negam sem contexto e o contexto desaparece no fim da
transação. A conexão de migration/admin é separada. Compatibilidade dessa bridge
com o pooler é gate obrigatório; sem prova, RLS continua defesa deny-by-default,
mas não é declarado isolamento por usuário.

`primary_player_links.version` é token de concorrência otimista. Replace exige
precondition da versão lida e retorna `409` em mismatch; não há last-write-wins
silencioso. `EnsureCrownPilotUser` cria/recupera por `firebase_uid` de forma
idempotente e relê após unique conflict de primeiro login.

Supabase é infraestrutura de banco, não backend da aplicação. SDK C# do Supabase
não é dependência central.

### Player Tag

Player Tag é referência pública/read-only. Vínculo inicial usa
`public_profile` e `unverified`. Não armazenar credenciais Supercell, snapshot,
coleção, Arena, battle history, cache persistente ou payload raw.

Lookup fica atrás de port `IClashRoyaleClient`, implementado na Infrastructure.
Somente resultado `resolved` cria/substitui vínculo; falha externa preserva o
vínculo anterior.

### Ambientes e deploy

| Ambiente | Auth | Banco | Regra |
|---|---|---|---|
| Local | Emulator/fixtures | Supabase CLI/Docker descartável | sem produção |
| PR Preview | sem login real | sem produção | hostname efêmero |
| Staging | Firebase separado, Google real | Supabase separado | hostname fixo, E2E |
| Production | projeto próprio | banco próprio | promoção controlada via `main` |

Preview e Staging são problemas diferentes. Preview valida build/UI/smoke sem
URL fixa ou autenticação real; Staging valida integração completa.

Frontend gera `dist/` e pode rodar em Vercel, CDN, Nginx ou storage estático.
API roda em imagem Docker/OCI ASP.NET Core, com secrets em runtime, filesystem
efêmero e health endpoints. Nenhuma API proprietária da Vercel é obrigatória.

## Motivação

- portabilidade de runtime e hospedagem;
- separação clara de responsabilidades;
- menor acoplamento a providers;
- PostgreSQL tratado como PostgreSQL;
- evolução arquitetural adequada ao longo prazo;
- autorização crítica centralizada no backend.

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
- EF Core e SQL duplicando schema;
- Firebase UID como chave primária de domínio.

## Gates e revisão

- Domain e Application não importam SDK, HTTP ou framework;
- Task 01 produz API, frontend, EF Core, Npgsql, testes e Docker sem Auth real;
- tokens Firebase têm validação por ambiente e testes de rotação;
- migrations EF e SQL RLS não duplicam schema;
- PR não usa Firebase/Supabase production nem API Clash Royale live;
- frontend pode sair da Vercel e API pode iniciar fora dela;
- Fase 003 permanece bloqueada pelos gates de API data, retenção, ownership,
  egress, meta e compliance da Fase 001.

Não há dados legados ou dual-write. Mudanças incompatíveis usam expand/contract.
Reabrir ADR se Firebase, PostgreSQL, portabilidade ou limites operacionais
invalidarem a decisão.
