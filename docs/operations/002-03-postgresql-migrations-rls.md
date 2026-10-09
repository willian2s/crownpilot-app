# PostgreSQL, migrations e harness RLS — 002-03

Esta subtarefa prepara somente a infraestrutura local e a ordem operacional.
Não implementa login, repositories, tabelas de domínio, policies concretas ou
prova A/B. Esses itens pertencem às subtarefas seguintes, principalmente
`002-06-modelar-persistencia-repositories-e-rls.md`.

## Runtime local

O caminho oficial usa Docker Compose com imagem pinada
`postgres:17.6-alpine`. Não é necessário instalar PostgreSQL no host. A
imagem executa banco descartável em `127.0.0.1:54322`, database
`crownpilot_local`, usuário `postgres` e senha fixa somente para este container.
O diretório de dados é `tmpfs`; `db:reset` remove o container e seus dados.

Docker Desktop no Mac ou Docker Engine + Compose v2 no Linux são pré-requisitos.
Supabase CLI não é dependência do runner local: Supabase é o PostgreSQL remoto,
consumido por Npgsql em job controlado. Se CLI Supabase for adotado em etapa
futura, sua versão deve ser pinada antes de substituir o contrato Docker.

## Interface cross-platform

Os scripts Node são chamados somente por npm e não carregam `.env.local`:

```text
npm run db:start --prefix frontend
npm run db:stop --prefix frontend
npm run db:migrate --prefix frontend
npm run db:security --prefix frontend
npm run db:fixtures --prefix frontend
npm run db:reset --prefix frontend
```

`db:reset` executa, nesta ordem, `down --volumes`, start, migration EF, SQL de
segurança e fixtures. Nenhum comando migra API no startup.

Todos os comandos falham fechado quando `CROWNPILOT__ENVIRONMENT` ou
`ASPNETCORE_ENVIRONMENT` resolve para `Preview`, `Staging` ou `Production`, ou
quando connection string aponta para host não local. Docker ausente ou daemon
parado retorna exit `2`; não há sucesso simulado.

## Migrations e ownership

`src/Infrastructure/Persistence/CrownPilotDbContext.cs` reserva schema
`crownpilot` e a migration inicial cria somente esse schema. EF Core é dono das
tabelas, schemas, colunas, constraints, índices e histórico. A migration de domínio com
`crownpilot_users` e `primary_player_links` será criada por `002-06`.

SQL versionado em `database/security/` configura apenas:

- roles `crownpilot_runtime` e `crownpilot_migrator` sem login, superuser ou
  `BYPASSRLS`;
- owner e grants mínimos do schema privado `crownpilot` já criado pelo EF;
- default privileges para objetos EF futuros;
- revogação de schema para `PUBLIC` e, se existirem, roles Supabase
  `anon`/`authenticated`/`service_role`.

Não há `CREATE TABLE` nem `CREATE POLICY` nesse SQL. Schema privado não é
adicionado ao conjunto exposto pela Supabase Data API. Provisionamento de login,
senha e owners reais de Staging/Production fica fora do Git.

## Referências por ambiente

| Ambiente | Referência versionada | Connection string | Operação |
|---|---|---|---|
| Local | Docker Compose, `crownpilot_local`, `127.0.0.1:54322` | derivada localmente ou `CROWNPILOT__DATABASE__CONNECTIONSTRING` | scripts `db:*` |
| Local dev | projeto Supabase exclusivo de desenvolvimento | secret de runtime separado, Session pooler, TLS obrigatório | desenvolvimento/integrações manuais após aprovação da ADR 005; fora de `db:reset` |
| Preview | nenhum banco | proibida | build/smoke sem dados |
| Staging | projeto Supabase separado; ref/região fornecidos no provisionamento | secret de runtime separado, Session pooler, TLS obrigatório após aprovação da ADR 005 | migration job manual/protegido |
| Production | projeto Supabase próprio; nunca ref de Staging | secret de runtime separado, Session pooler, TLS obrigatório após aprovação da ADR 005 | migration job aprovado |

Render usa região Virgínia (`us-east`) e Supabase usa `us-east-1` (Northern
Virginia). Dados pessoais ficam fora do Brasil; antes de dados reais, registrar
política de privacidade, DPA, subprocessadores, backups, residência e egress.
Nenhum projeto remoto é provisionado por esta subtarefa. Preview não recebe
secrets de Staging/Production.

## Contexto transacional RLS

O contrato para `002-06` usa role runtime sem `BYPASSRLS`, authorization da API
como requisito independente e contexto `SET LOCAL` dentro de uma transação:

```text
BEGIN
  set_config('app.firebase_uid', verifiedFirebaseUid, true)
  resolve/create crownpilot user
  set_config('app.crownpilot_user_id', resolvedInternalId, true)
  repository operations
COMMIT or ROLLBACK
```

Contexto ausente deve negar acesso. Nunca usar `auth.uid()` de outro provider,
UID vindo do request, variável de sessão persistente ou pseudo-SQL interpolado.
`002-06` adiciona policies concretas, testa usuário A/B, cleanup em pool,
conexão direta e session pooler. Sem essa prova, não se declara isolamento por
usuário em pooler.

## Fixtures e rollback

`database/fixtures/` contém apenas contrato/README neste ponto. `db:fixtures`
falha fechado até `002-06` fornecer fixtures sanitizadas para tabelas aprovadas;
não existe no-op que reporte sucesso. Fixtures futuras não podem conter UID
real, Player Tag real, payload externo, credencial ou dados de produção.

Reset local é destrutivo e reversível por recriação do container. Em ambientes
remotos, rollback de aplicação preserva schema e vínculos; migrations
compartilhadas são forward-only e não executam `down`. Migrations `down` ficam
restritas ao desenvolvimento local; rollback de dados usa backup/restore conforme
runbook. Bundle/one-shot final para release pertence a `002-11` depois do schema
concreto.

## Validação

Sem daemon Docker, ainda é possível validar compilação, migration metadata e
fail-closed:

```text
dotnet tool restore
dotnet restore CrownPilot.sln
dotnet build CrownPilot.sln --configuration Release
dotnet test CrownPilot.sln --configuration Release
dotnet ef migrations list --project src/Infrastructure/Infrastructure.csproj --startup-project src/Infrastructure/Infrastructure.csproj
CROWNPILOT__ENVIRONMENT=Production npm run db:reset --prefix frontend
npm run db:start --prefix frontend
```

O comando Production deve retornar `2` antes de tocar Docker. `db:start` só é
considerado validado quando daemon, container e healthcheck passarem. Depois do
daemon disponível, executar `db:reset`; até `002-06`, o estágio final de
fixtures deve continuar bloqueando explicitamente.
