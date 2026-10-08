# Boundaries, contrato HTTP e ambientes

Task `002-02` fixa contratos; não provisiona Firebase/Supabase nem valida Firebase
ID Token real. Fixture de bearer existe somente em Local para provar pipeline
authentication → authorization. Task `002-04` substituirá fixture por adapter
Firebase sem mudar contrato HTTP.

## Boundaries

```text
React/Vite
  -> Firebase SDK (login/refresh/logout; Task 002-04)
  -> Authorization: Bearer <Firebase ID Token>
  -> API ASP.NET Core
  -> Application ports/use cases
  -> Domain invariants
  -> Infrastructure adapters (Firebase, EF Core/Npgsql, provider externo)
  -> PostgreSQL/Supabase
```

| Camada | Responsabilidade | Dependências proibidas |
|---|---|---|
| `Domain` | entidades, value objects e invariantes | HTTP, ASP.NET Core, EF Core, Npgsql, Firebase, Supabase, Vercel |
| `Application` | casos de uso, autorização de recurso e ports | SDK de provider, HTTP e infraestrutura concreta |
| `Infrastructure` | adapters Firebase, EF Core/Npgsql, PostgreSQL e provider externo | lógica de transporte HTTP |
| `Api` | HTTP, DTOs, ProblemDetails, auth pipeline, CORS, OpenAPI e composição | acesso direto a entidades/repositories de outro módulo |

Authentication responde quem é o subject. Authorization responde o que subject
pode fazer. `AuthenticatedSubject { FirebaseUid }` nasce somente após validação de
claims no boundary de authentication; Application converte para
`AuthContext { CrownPilotUserId }`. UID vindo de body, query ou header não é
confiável. `CrownPilotUserId` é ID interno e não depende do provider.

Não há cookie, sessão ou segundo esquema bearer. `Domain` e `Application` não
validam JWT. Task `002-04` implementará issuer, audience/project, assinatura,
expiração, `kid` e claims temporais com SDK suportado.

## Contrato HTTP base

- prefixo REST: `/api/v1`;
- JSON e HTTPS fora de Local;
- `Authorization: Bearer <Firebase ID Token>`;
- OpenAPI code-first gerado por `Microsoft.AspNetCore.OpenApi`;
- `/openapi/v1.json` e `/docs` compartilham mesmo documento;
- UI `/docs` é visualização, não contrato concorrente;
- concorrência usa `version` retornada e `expectedVersion` no JSON, com `409`;
  não há ETag/If-Match paralelo.

Rotas de negócio completas pertencem à Task `002-07`. Contratos mínimos futuros:

```text
GET    /api/v1/me/player-link
PUT    /api/v1/me/player-link
DELETE /api/v1/me/player-link
DELETE /api/v1/me
```

`ProblemDetails` contém `status`, `title`, `type`, `traceId` e `code` estável.
Não contém token, UID, Player Tag, URL externa, stack trace, payload bruto ou
detalhe de infraestrutura.

| Situação | Status | Código/convenção |
|---|---:|---|
| sintaxe/input inválido | `400` | `invalid_request` ou `invalid_player_tag` |
| bearer ausente/inválido | `401` | `authentication_required` |
| bearer válido sem permissão | `403` | `authorization_forbidden` |
| recurso próprio ausente | `404` | `resource_not_found` ou `player_not_found` |
| precondition/version conflict | `409` | `version_conflict` |
| validação semântica necessária | `422` | `semantic_validation_failed` |
| provider limitado | `429` | `provider_rate_limited`, `Retry-After` quando aplicável |
| lookup Clash Royale indisponível/mal configurado | `503` | `provider_unavailable` ou `provider_misconfigured` |
| Firebase indisponível durante revogação sensível | `503` | `authentication_unavailable` em ProblemDetails genérico, sem detalhe do provider |
| falha inesperada | `500` | `internal_error` e `traceId` |

Exclusão de dados CrownPilot exigirá `auth_time` presente, janela de cinco
minutos e tolerância de 60 segundos; falha retorna `403`/
`reauthentication_required`. Não exclui Firebase/Google.

## Ambientes e configuração

ASP.NET Core `Development` mapeia para runtime `Local`. `Preview`, `Staging` e
`Production` exigem nomes explícitos. `CrownPilot:Environment` deve coincidir
com `ASPNETCORE_ENVIRONMENT`; divergência falha no startup. CORS aceita somente
origins HTTPS configuradas (HTTP apenas localhost em Local), nunca `*`.

| Ambiente | Auth | Banco/secrets | OpenAPI | CORS e host |
|---|---|---|---|---|
| Local | fixture agora; Emulator depois | descartável/local, sem produção | JSON/UI habilitados | localhost allowlist |
| Preview | sem login real | nenhum secret/dado de Staging/Production | `404` | efêmero, sem hostname fixo |
| Staging | Firebase separado + Google real (Task 004) | Supabase separado | JSON/UI habilitados para equipe | hostname fixo allowlisted |
| Production | projeto Firebase próprio | banco/secrets próprios | `404` por padrão | hostname fixo allowlisted |

Variáveis `VITE_*` são públicas e entram no bundle; Firebase service account,
database connection string, provider token e chaves privadas são server-only e
entram somente no runtime. Preview nunca recebe secrets de Staging/Production.
`.env.local` é ignorado e não deve ser lido ou embutido em bundle.

Staging e Production devem terminar TLS em ingress confiável antes de encaminhar
requests à API; bearer nunca deve atravessar rede plaintext. O ingress deve
preservar o esquema HTTPS por forwarded headers configurados pelo hosting e
restringir hostnames às allowlists do ambiente. Local pode usar HTTP localhost.

## Persistência e deploy

EF Core é dono de tabelas, colunas, constraints, índices e migrations. SQL
versionado cobre somente roles, grants, extensões e RLS; não recria schema EF.
Ordem controlada: migration EF → SQL de segurança/RLS → fixtures sanitizadas.
Migrations não rodam no startup de múltiplas réplicas.

Frontend produz `frontend/dist` estático. API produz imagem OCI portátil; Render
é hosting inicial, não boundary nem dependência do backend. Vercel é alternativa
somente para frontend. Mesmo digest deve ser promovido entre Staging e Production;
rollback retorna ao digest anterior sem `Down` destrutivo automático.
