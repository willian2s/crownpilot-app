# CrownPilot database harness

Local database contract for SDD `002-03`.

## Runtime choice and version

- Canonical local runtime: Docker Compose, not a host PostgreSQL installation.
- PostgreSQL image: `postgres:17.6-alpine`, pinned in
  `docker-compose.yml`.
- Docker Desktop (Mac) or Docker Engine plus Compose v2 (Linux) is required.
- Supabase CLI is intentionally not required by local scripts. Supabase remains
  the hosted PostgreSQL target; remote migration jobs use direct Npgsql
  connection strings supplied at runtime. A future CLI adapter must pin its
  version before replacing this Docker contract.
- The local container uses `tmpfs`; stopping/removing it does not preserve a
  database volume. This is deliberate disposable infrastructure.

The local endpoint is `127.0.0.1:54322`, database `crownpilot_local`, and the
`postgres/postgres` credential is local-only fixture data. The scripts never
load `.env.local`.

## Ownership and order

```text
EF Core migration -> database/security/*.sql -> database/fixtures/*.sql
```

`CrownPilotDbContext` currently reserves only schema `crownpilot`. EF Core owns
future application tables, columns, constraints, indexes, and migration history.
The security SQL creates only roles, schema/grants, default privileges, and
Data API denial; it does not create tables or concrete domain policies.

`crownpilot_runtime` and `crownpilot_migrator` are created without `BYPASSRLS`
and without login. Production role login/password setup stays outside Git. The
runtime role is never the owner of application tables. Concrete policies,
Firebase UID/user-ID bridge, A/B isolation, and pooler proof belong to `002-06`.

## Commands

Use npm as the cross-platform interface from repository root:

```text
npm run db:start --prefix frontend
npm run db:migrate --prefix frontend
npm run db:security --prefix frontend
npm run db:fixtures --prefix frontend
npm run db:stop --prefix frontend
```

`db:reset` is local-only and destructive:

```text
npm run db:reset --prefix frontend
```

Every command rejects `Preview`, `Staging`, and `Production`, and rejects a
non-loopback connection target. `db:reset` never runs from API startup or from
multiple replicas. It leaves the container available when a later stage fails
so the failure can be inspected; retry starts from `db:reset`.

`db:fixtures` currently fails closed because no domain fixture can be defined
before `002-06` creates its approved tables. This is intentional and is not an
empty successful stage.

## Connection configuration

Local scripts derive the disposable connection unless
`CROWNPILOT__DATABASE__CONNECTIONSTRING` is supplied. The value must point to
the local endpoint and database. `DATABASE_CONNECTION_STRING` is accepted only
as a local compatibility alias. Staging and Production use separate secret
values in the migration job; Preview has no database target.

Migration/admin and runtime connections are separate responsibilities. Runtime
context must be transaction-local:

```sql
BEGIN;
SELECT set_config('app.firebase_uid', $1, true);
SELECT set_config('app.crownpilot_user_id', $2, true);
-- authorized repository work
COMMIT; -- or ROLLBACK; both clear SET LOCAL context
```

No session-persistent context is allowed through pooling. `002-06` must prove
this with direct and session-pooler connections before claiming per-user RLS
isolation; transaction-pooler use remains a separate gate.
