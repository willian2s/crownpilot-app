# 002-03 — Preparar PostgreSQL, migrations e harness RLS

- **Ticker:** `002`
- **Número:** `03`
- **Status:** `pending`

## Objetivo e resultado esperado

Configurar PostgreSQL local via Supabase CLI/Docker e o pipeline de migrations EF
Core + SQL de segurança. Preparar referências de banco separadas por ambiente,
roles e harness RLS, sem tratar Supabase como backend ou implementar o fluxo de
login.

## Requisitos cobertos

- PostgreSQL/Supabase local descartável;
- referências de banco separadas para Staging e Production;
- EF Core como dono do schema da aplicação;
- Npgsql e migrations reproduzíveis;
- harness/roles para RLS e grants versionados sem duplicar schema; policies
  concretas pertencem à Task 002-06;
- schema CrownPilot dedicado, fora da Data API;
- nenhum acesso de Preview à Production.

## Escopo incluído

- fixar versão/uso de Supabase CLI e Docker;
- iniciar PostgreSQL local e documentar connection string sem secrets reais;
- definir projetos/refs de banco separados, sem commitar credenciais;
- criar harness para aplicar migrations EF Core e depois SQL RLS/grants;
- preparar comando e contrato para gerar artefato de migration revisável/bundle
  one-shot depois do schema final; não migrar no startup das réplicas da API;
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
- `EnvironmentName` e `DatabaseOptions`;
- `tests/Persistence/` e `tests/Integration/`.

## Passos de implementação

1. Fixar Supabase CLI/Docker e iniciar banco local descartável.
2. Configurar pipeline local: banco limpo, migrations EF Core, SQL RLS/grants,
   fixtures quando schema existir; nenhum SQL de segurança pode criar tabela do
   schema.
3. Validar connection strings Npgsql, SSL/configuração e fail-closed.
4. Preparar role de runtime sem `BYPASSRLS`, schema dedicado, runner de segurança
   e mecanismo transacional de contexto; não criar policies de tabelas ainda
   inexistentes. A Task 002-06 prova conexão direta/session pooler com schema
   concreto e registra transaction pooler como gate separado se for adotado.
5. Registrar referências de banco separadas para Local, Staging e Production; não
   commitar secrets.
6. Registrar critérios para `sa-east-1`, DPA, backups, subprocessadores e
   residência antes de dados reais.
7. Manter provisionamento Production como etapa controlada posterior.

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

O comando de reset deve provar a ordem EF Core -> SQL RLS/grants -> fixtures
quando o schema existir. Nesta task, validar apenas conexão, roles, fail-closed e
idempotência do pipeline; testes de usuário A/B, policies e reset de contexto
concreto pertencem à Task 002-06.

## Definição de pronto

- PostgreSQL local inicia por comando documentado;
- migrations EF Core são a única fonte de tabelas/constraints;
- SQL separado aplica somente RLS/grants/objetos permitidos;
- não cria policies de domínio antes da migration das tabelas; bridge A/B fica na
  Task 002-06;
- reset local é reproduzível e não usa produção;
- roles, grants, contexto transacional e cleanup em pool estão documentados para a
  prova concreta da Task 002-06;
- referências Staging/Production são distintas e sem secrets versionados;
- schema CrownPilot não depende da Supabase Data API;
- pipeline de bundle/job está definido, mas o artefato final pertence à Task
  002-11 e será gerado após as migrations da Task 002-06;
- região, backups e residência são gates antes de dados reais;
- nenhuma tabela futura ou payload de jogador é criada por conveniência.

## Riscos e cuidados

- Não usar `auth.uid()` de um provider diferente no RLS.
- Não rodar migrations automaticamente em múltiplas réplicas.
- Não duplicar schema entre EF Core e SQL.
- Não confundir projeto de banco Supabase com projeto Firebase.
- Não provisionar Production durante bootstrap local.
