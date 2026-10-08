# 002-06 — Modelar persistência, repositories e RLS

- **Ticker:** `002`
- **Número:** `06`
- **Status:** `pending`

Esta task permanece bloqueada até ADR 005 ser aprovada e `002-13` -> `002-14` ->
`002-15` concluírem gates verdes. Referências EF/.NET abaixo descrevem baseline
histórico; implementação deve seguir `pgx`/`sqlc` e o ajuste Go de `002-15`.

## Objetivo e resultado esperado

Modelar e provar somente usuário CrownPilot e vínculo primário de perfil público
usando `pgx`/`sqlc` e migrations SQL com `goose`, com ID interno, isolamento por usuário e sem criar
snapshot de API. Casos de uso e endpoints ficam na Task
`002-07-implementar-casos-de-uso-e-api-v1.md`.

## Requisitos cobertos

- PostgreSQL server-side hospedado no Supabase;
- repositories `pgx`/`sqlc` e transações Go;
- migrations SQL da aplicação executadas por `goose`;
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
- migrations SQL e SQL posterior de RLS/grants;
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

- `002-02-estabelecer-boundaries-contrato-base-e-ambientes.md`,
  `002-03-preparar-postgresql-migrations-e-harness-rls.md` e gates de
  `002-13` a `002-15`;
  lookup/authentication podem ser substituídos por ports e fakes nesta task;
- contrato `PrimaryPlayerLink` da spec;
- conexão PostgreSQL local/Staging aprovada;
- [ADR 004](../../decisions/004-aspnet-core-react-vite-firebase-postgresql.md) é
  canônica até aprovação da [ADR 005](../../decisions/005-go-react-vite-firebase-postgresql.md);
  o destino Go desta task permanece condicionado à aprovação da ADR 005.

## Arquivos e símbolos prováveis

- `internal/identity/` e `internal/playerlink/`;
- `internal/platform/postgres/`;
- `database/migrations/`, `database/security/` e testes Go de persistência/RLS.

## Passos de implementação

1. Modelar ID interno, usuário, vínculo e invariantes no Domain.
2. Configurar queries/mappings `sqlc`, FK, unique `firebase_uid`, `NOT NULL` checks,
   timestamps e token `version` de concorrência.
3. Criar migrations SQL e executá-las com `goose`; não criar as mesmas tabelas em
   SQL de segurança.
4. Implementar repositories usando ID resolvido do `AuthContext`; ignorar UID do
   body/query.
5. Expor repositories e operações transacionais necessárias, deixando validação do
   provider e orquestração de casos de uso para a Task
   `002-07-implementar-casos-de-uso-e-api-v1.md`.
6. Aplicar e testar RLS/grants após migration `goose`, mantendo authorization obrigatória.
   Resolver/criar `crownpilot_users` por Firebase UID verificado dentro da mesma
   transação que define contexto RLS; nunca confiar em UID do request. Policies
   de `crownpilot_users` usam `app.firebase_uid` somente para a linha própria;
   policies de `primary_player_links` usam `app.crownpilot_user_id` e negam
   contexto ausente. `set_config(..., true)` substitui pseudo-SQL interpolado.

## Testes e comandos de validação

```text
go test ./internal/identity/... ./internal/playerlink/...
go test ./internal/platform/postgres/...
npm run test:rls --prefix frontend
```

Cobrir usuário novo/existente, A/B, anônimo, contexto RLS ausente/adulterado,
`app.firebase_uid` divergente, UID adulterado, FK/unique/NOT NULL/check
constraints, role runtime sem `BYPASSRLS`, acesso Data API/PostgREST negado,
cleanup após pooling, reset/rollback transacional e bridge direta no PostgreSQL
local descartável/Session pooler no Supabase dev e hospedado. Consultar
`pg_class.relrowsecurity` e falhar se
qualquer tabela do schema `crownpilot` estiver sem RLS ativo.

## Definição de pronto

- schema mínimo possui somente usuário CrownPilot e vínculo primário;
- schema CrownPilot dedicado não é acessado via Supabase Data API;
- migrations SQL com `goose` são fonte das tabelas/constraints e aplicam em PostgreSQL;
- SQL de RLS/grants não duplica schema;
- UUID interno é usado em authorization e relacionamentos;
- A não lê/altera B;
- usuário novo resolve/cria identidade sem abrir linhas de outro usuário;
- contexto externo e interno são transacionais e não vazam em conexão pooled;
- repositories não aceitam identidade controlada pelo request;
- contexto RLS não vaza entre transações/conexões pooled;
- teste de catálogo `pg_class.relrowsecurity` falha para qualquer tabela do schema
  `crownpilot` sem RLS ativo; esse teste roda no CI local e é gate antes de Staging;
- acesso Data API/PostgREST ao schema CrownPilot é negado;
- nenhum snapshot/raw/API data é persistido;
- testes de aplicação, persistência, integração e RLS passam.

## Riscos e cuidados

- Não usar Firebase UID como PK interna do domínio.
- Não alegar RLS como substituto da authorization da API.
- Não usar `auth.uid()` sem bridge explicitamente testada.
- Não criar índices/TTL “para o futuro”.
- Não impor unicidade global: uma tag pública pode ser vinculada por mais de um usuário.
