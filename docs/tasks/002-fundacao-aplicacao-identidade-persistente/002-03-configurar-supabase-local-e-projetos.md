# 002-03 — Configurar PostgreSQL/Supabase local e projetos

- **Ticker:** `002`
- **Número:** `03`
- **Status:** `pending`

## Objetivo e resultado esperado

Configurar PostgreSQL local via Supabase CLI/Docker, projetos separados de banco
por ambiente e o pipeline de migrations EF Core + SQL de segurança. Preparar
Firebase Emulator/fixtures e configuração de projetos sem implementar o fluxo de
login.

## Requisitos cobertos

- PostgreSQL/Supabase local descartável;
- projetos de banco separados para Staging e Production;
- EF Core como dono do schema da aplicação;
- Npgsql e migrations reproduzíveis;
- RLS/grants/policies versionados sem duplicar schema;
- Firebase project IDs por ambiente e Emulator/fixtures;
- nenhum acesso de Preview à Production.

## Escopo incluído

- fixar versão/uso de Supabase CLI e Docker;
- iniciar PostgreSQL local e documentar connection string sem secrets reais;
- configurar projeto Firebase Emulator e fixtures de token;
- definir projetos/refs de banco separados, sem commitar credenciais;
- criar harness para aplicar migrations EF Core e depois SQL RLS/grants;
- documentar `sa-east-1` como preferência condicionada para Staging/Production;
- documentar criação manual/provisionamento controlado sem executar Production;
- definir seed/fixtures sanitizados e reset local.

## Escopo excluído

- login Google real ou middleware de autenticação;
- repository de vínculo e casos de uso;
- dados reais de jogadores;
- schema de snapshot, cache, coleção ou payload externo;
- SDK C# do Supabase ou Data API no browser;
- provisionamento irreversível de Production antes dos gates.

## Dependências

- `002-01` e `002-02`;
- acesso administrativo aos projetos, fornecido fora do Git;
- Docker disponível localmente/CI;
- [ADR 004](../../decisions/004-aspnet-core-react-vite-firebase-postgresql.md).

## Arquivos e símbolos prováveis

- `supabase/config.toml` e SQL versionado de RLS/grants;
- `src/Infrastructure/Persistence/` e `Migrations/`;
- scripts `database/reset-local.*` e `database/apply-security.*`;
- `firebase.json`, configuração do Emulator e fixtures;
- `EnvironmentName`, `DatabaseOptions`, `FirebaseProjectOptions`;
- `tests/Persistence/` e `tests/Integration/`.

## Passos de implementação

1. Fixar Supabase CLI/Docker e iniciar banco local descartável.
2. Documentar Firebase Emulator/fixtures sem criar credenciais reais.
3. Configurar pipeline local: banco limpo, migrations EF Core, SQL RLS/grants,
   fixtures; nenhum SQL de segurança pode criar tabela do schema.
4. Validar connection strings Npgsql, SSL/configuração e fail-closed.
5. Configurar role de runtime sem `BYPASSRLS`, contexto transacional
   `app.crownpilot_user_id`, reset por transação e policies deny-by-default;
   provar pooler compatível ou registrar RLS como gate não comprovado.
6. Registrar IDs e allowlists de projetos separados para Local, Staging e
   Production; não commitar secrets.
7. Registrar critérios para `sa-east-1`, DPA, backups, subprocessadores e
   residência antes de dados reais.
8. Manter provisionamento Production como etapa controlada posterior.

## Testes e comandos de validação

```text
supabase start
dotnet ef database update
dotnet test --filter Category=Persistence
dotnet test --filter Category=Rls
./database/apply-security.sh
./database/load-fixtures.sh
supabase stop
```

O comando de reset deve provar a ordem EF Core -> SQL RLS/grants -> fixtures.
Testar conexão anônima, contexto ausente, usuário A, usuário B e reset do
contexto ao fim da transação.

## Definição de pronto

- PostgreSQL local inicia por comando documentado;
- migrations EF Core são a única fonte de tabelas/constraints;
- SQL separado aplica somente RLS/grants/objetos permitidos;
- reset local é reproduzível e não usa produção;
- Firebase Emulator/fixtures estão disponíveis para testes posteriores;
- projetos Staging/Production são distintos e sem secrets versionados;
- região, backups e residência são gates antes de dados reais;
- nenhuma tabela futura ou payload de jogador é criada por conveniência.

## Riscos e cuidados

- Não usar `auth.uid()` de um provider diferente no RLS.
- Não rodar migrations automaticamente em múltiplas réplicas.
- Não duplicar schema entre EF Core e SQL.
- Não confundir projeto de banco Supabase com projeto Firebase.
- Não provisionar Production durante bootstrap local.
