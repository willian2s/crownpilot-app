# 002-06 — Modelar persistência, repositories e RLS

- **Ticker:** `002`
- **Número:** `06`
- **Status:** `pending`

## Objetivo e resultado esperado

Modelar e provar somente usuário CrownPilot e vínculo primário de perfil público
usando EF Core + Npgsql, com ID interno, isolamento por usuário e sem criar
snapshot de API. Casos de uso e endpoints ficam na Task
`002-07-implementar-casos-de-uso-e-api-v1.md`.

## Requisitos cobertos

- PostgreSQL server-side hospedado no Supabase;
- `CrownPilotDbContext` e mappings EF Core;
- migrations EF Core para schema da aplicação;
- `firebase_uid` externo separado do UUID interno;
- usuário A isolado de B por authorization server-side;
- `public_profile` + `unverified`;
- replace protegido contra concorrência e falhas por `expectedVersion` JSON;
- RLS/grants como defesa adicional, sem duplicar schema.

## Escopo incluído

- entidade `CrownPilotUser` com UUID interno e `firebase_uid` único;
- entidade `PrimaryPlayerLink` com `user_id` como primary key/FK;
- repository baseado no `CrownPilotUserId` resolvido no backend;
- bootstrap RLS transacional: `set_config('app.firebase_uid', <claim>, true)`
  permite resolver/criar somente a linha correspondente em `crownpilot_users`;
  depois `set_config('app.crownpilot_user_id', <id>, true)` limita os vínculos;
- repositories e transações necessárias aos casos de uso, sem possuir transporte
  HTTP;
- transaction/concurrency token/precondition para replace;
- timestamps e `last_validated_at`;
- migration EF Core e SQL posterior de RLS/grants;
- policies sem pressupor claims nativas de provider externo;
- validação que impede snapshot/raw/cache/coleção.

## Escopo excluído

- tabelas globais de players/decks/cards;
- payload raw, nome externo, Arena, coleção, battle history ou cache;
- unicidade global de Player Tag;
- TTL, retenção de API data ou freshness garantida;
- Data API client-side permissiva;
- migration automática no startup de múltiplas réplicas.

## Dependências

- `002-02-estabelecer-boundaries-contrato-base-e-ambientes.md` e
  `002-03-preparar-postgresql-migrations-e-harness-rls.md`;
  lookup/authentication podem ser substituídos por ports e fakes nesta task;
- contrato `PrimaryPlayerLink` da spec;
- conexão PostgreSQL local/Staging aprovada;
- [ADR 004](../../decisions/004-aspnet-core-react-vite-firebase-postgresql.md).

## Arquivos e símbolos prováveis

- `src/Domain/Identity/CrownPilotUserId.cs`;
- `src/Domain/Identity/PrimaryPlayerLink.cs`;
- `src/Application/Ports/IUserRepository.cs`;
- `src/Application/Ports/IPlayerLinkRepository.cs`;
- `src/Infrastructure/Persistence/CrownPilotDbContext.cs`;
- `src/Infrastructure/Persistence/Configurations/`;
- `src/Infrastructure/Persistence/Migrations/`;
- `database/security/` e `tests/Persistence`, `tests/Integration`, `tests/Rls`.

## Passos de implementação

1. Modelar ID interno, usuário, vínculo e invariantes no Domain.
2. Configurar mappings EF Core, FK, unique `firebase_uid`, `NOT NULL` checks,
   timestamps e token `version` de concorrência.
3. Gerar migration EF Core; não criar as mesmas tabelas em SQL de segurança.
4. Implementar repositories usando ID resolvido do `AuthContext`; ignorar UID do
   body/query.
5. Expor repositories e operações transacionais necessárias, deixando validação do
   provider e orquestração de casos de uso para a Task
   `002-07-implementar-casos-de-uso-e-api-v1.md`.
6. Aplicar e testar RLS/grants após migration EF, mantendo authorization obrigatória.
   Resolver/criar `crownpilot_users` por Firebase UID verificado dentro da mesma
   transação que define contexto RLS; nunca confiar em UID do request. Policies
   de `crownpilot_users` usam `app.firebase_uid` somente para a linha própria;
   policies de `primary_player_links` usam `app.crownpilot_user_id` e negam
   contexto ausente. `set_config(..., true)` substitui pseudo-SQL interpolado.

## Testes e comandos de validação

```text
dotnet test --filter Category=Unit
dotnet test --filter Category=Application
dotnet test --filter Category=Persistence
dotnet test --filter Category=Integration
dotnet test --filter Category=Rls
```

Cobrir usuário novo/existente, A/B, anônimo, contexto RLS ausente/adulterado,
`app.firebase_uid` divergente, UID adulterado, FK/unique/NOT NULL/check
constraints, role runtime sem `BYPASSRLS`, acesso Data API/PostgREST negado,
cleanup após pooling, reset/rollback transacional e bridge direta/session pooler.

## Definição de pronto

- schema mínimo possui somente usuário CrownPilot e vínculo primário;
- schema CrownPilot dedicado não é acessado via Supabase Data API;
- EF Core é fonte das tabelas/constraints e migration aplica em PostgreSQL;
- SQL de RLS/grants não duplica schema;
- UUID interno é usado em authorization e relacionamentos;
- A não lê/altera B;
- usuário novo resolve/cria identidade sem abrir linhas de outro usuário;
- contexto externo e interno são transacionais e não vazam em conexão pooled;
- repositories não aceitam identidade controlada pelo request;
- contexto RLS não vaza entre transações/conexões pooled;
- acesso Data API/PostgREST ao schema CrownPilot é negado;
- nenhum snapshot/raw/API data é persistido;
- testes de aplicação, persistência, integração e RLS passam.

## Riscos e cuidados

- Não usar Firebase UID como PK interna do domínio.
- Não alegar RLS como substituto da authorization da API.
- Não usar `auth.uid()` sem bridge explicitamente testada.
- Não criar índices/TTL “para o futuro”.
- Não impor unicidade global: uma tag pública pode ser vinculada por mais de um usuário.
