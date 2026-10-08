# 002-03 — Preparar PostgreSQL, migrations e harness RLS

- **Ticker:** `002`
- **Número:** `03`
- **Status:** `completed`

## Objetivo e resultado esperado

Configurar PostgreSQL descartável via Docker Compose para testes/reset e registrar
projeto Supabase exclusivo de desenvolvimento via Session pooler. Preparar o
pipeline de migrations EF Core e SQL de segurança, referências de banco separadas
por ambiente, roles e harness RLS, sem tratar Supabase como backend ou implementar
o fluxo de login.

## Requisitos cobertos

- PostgreSQL Docker local descartável; projeto Supabase dev separado e persistente,
  usado via Session pooler após aprovação da ADR 005;
- referências de banco separadas para Staging e Production;
- EF Core como dono do schema da aplicação;
- Npgsql e migrations reproduzíveis;
- harness/roles para RLS e grants versionados sem duplicar schema; policies
  concretas pertencem à Task `002-06-modelar-persistencia-repositories-e-rls.md`;
- schema CrownPilot dedicado, fora da Data API;
- nenhum acesso de Preview à Production.

## Escopo incluído

- fixar PostgreSQL/Docker Compose como runner local descartável; projeto Supabase
  dev é referência persistente de desenvolvimento e usa Session pooler;
- iniciar PostgreSQL local e documentar connection string sem secrets reais;
- definir projetos/refs de banco separados, sem commitar credenciais;
- criar harness para aplicar migrations EF Core e depois SQL RLS/grants;
- preparar comando e contrato para gerar artefato de migration revisável/bundle
  one-shot depois do schema final; não migrar no startup das réplicas da API;
- documentar Render em Virgínia (`us-east`) e Supabase em `us-east-1` (Northern
  Virginia) para Staging/Production; dados pessoais ficam fora do Brasil e a
  transferência internacional deve constar na política de privacidade antes de
  dados reais;
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

- `002-01-bootstrap-toolchain.md` e
  `002-02-estabelecer-boundaries-contrato-base-e-ambientes.md`;
- acesso administrativo aos projetos, fornecido fora do Git;
- Docker disponível localmente/CI;
- [ADR 004](../../decisions/004-aspnet-core-react-vite-firebase-postgresql.md).

## Arquivos e símbolos prováveis

- `supabase/config.toml` e SQL versionado de RLS/grants;
- `src/Infrastructure/Persistence/` e `Migrations/`;
- contrato de scripts `db:*` no `package.json`, fornecido pelo bootstrap da
  Task `002-01-bootstrap-toolchain.md`;
- `EnvironmentName` e `DatabaseOptions`;
- `tests/Persistence/` e `tests/Integration/`.

## Passos de implementação

1. Fixar Docker Compose/PostgreSQL e definir o contrato cross-platform dos scripts npm
   de banco que o bootstrap fornecerá.
2. Configurar pipeline local: banco limpo, migrations EF Core, SQL RLS/grants,
   fixtures quando schema existir; nenhum SQL de segurança pode criar tabela do
   schema.
3. Validar connection strings Npgsql, SSL/configuração e fail-closed.
4. Preparar role de runtime sem `BYPASSRLS`, schema dedicado, runner de segurança
   e mecanismo transacional de contexto; não criar policies de tabelas ainda
   inexistentes. A Task `002-06-modelar-persistencia-repositories-e-rls.md` prova
    conexão direta no PostgreSQL local e Session pooler no Supabase dev/hospedado com schema
    concreto; registra transaction pooler como gate separado se for adotado.
5. Registrar referências de banco separadas para Docker Local, Supabase dev,
   Staging e Production; não commitar secrets.
6. Registrar critérios para `us-east-1`, DPA, backups, subprocessadores,
   transferência internacional, residência e egress antes de dados reais.
7. Manter provisionamento Production como etapa controlada posterior.

## Pré-requisitos Mac/Linux

O runner descartável local exige Git, .NET SDK e Node.js LTS/npm nas versões pinadas pela Task
`002-01-bootstrap-toolchain.md`, `dotnet-ef` como ferramenta local do repositório e
Docker (Docker Desktop no Mac ou Docker Engine/Compose no Linux). Supabase CLI não
é dependência do runner local. Desenvolvimento integrado também pode usar o projeto
Supabase dev, com referência e secret de runtime fora do Git e Session pooler.
Não é necessário instalar PostgreSQL no host. Os comandos usam somente banco local
descartável e não exigem credenciais de Staging ou Production.

## Interface oficial de banco

Os scripts npm são a interface pública e cross-platform. A Task
`002-01-bootstrap-toolchain.md` deve
fornecer estes nomes no `package.json`, sem adicionar task runner ou outra
dependência apenas para encadear comandos:

| Comando oficial | Contrato |
|---|---|
| `npm run db:start` | inicia PostgreSQL local descartável via Docker Compose; falha se dependências locais não estiverem disponíveis |
| `npm run db:stop` | para o ambiente local; não toca Staging/Production |
| `npm run db:migrate` | aplica somente migrations EF Core no banco local configurado |
| `npm run db:security` | aplica somente SQL versionado de roles, grants e RLS/segurança, depois de `db:migrate` |
| `npm run db:fixtures` | carrega somente fixtures/seed sanitizados, depois de `db:security` |
| `npm run db:reset` | recria/resetta banco local e executa, nessa ordem, EF Core → SQL security/RLS → fixtures |

Fluxo oficial completo:

```text
npm run db:start
npm run db:reset
npm run db:stop
```

Quando for necessário inspecionar etapas separadamente, usar `db:migrate`,
`db:security` e `db:fixtures` na mesma ordem. `db:reset` é comando destrutivo
somente para o banco local descartável e deve falhar fechado se a referência de
banco apontar para Staging ou Production. O reset não aplica migrations no
startup da API nem em múltiplas réplicas.

O `package.json` e os scripts foram fornecidos por `002-01`; o runner Node agora
encapsula Docker, `dotnet ef` e `psql`. Desenvolvedores e CI usam somente a
interface `npm run`; nenhum runner interno vira contrato público.

Os comandos `db:start`, `db:reset` e `db:stop` operam somente PostgreSQL Docker
descartável. O projeto Supabase dev não é resetado por esses comandos; seu acesso
usa configuração separada e Session pooler.

## Testes e comandos de validação

```text
npm run db:start
npm run db:reset
dotnet test --filter Category=Persistence
dotnet test --filter Category=Rls
npm run db:stop
```

Os comandos acima são o contrato que será executável após o bootstrap. Nesta
task, validar apenas conexão, roles, fail-closed e idempotência do pipeline;
testes de usuário A/B, policies e reset de contexto concreto pertencem à Task
`002-06-modelar-persistencia-repositories-e-rls.md`. O reset deve provar a ordem
EF Core -> SQL RLS/grants -> fixtures quando o schema existir.

## Definição de pronto

- PostgreSQL local inicia por comando documentado;
- migrations EF Core são a única fonte de tabelas/constraints;
- SQL separado aplica somente RLS/grants/objetos permitidos;
- não cria policies de domínio antes da migration das tabelas; bridge A/B fica na
  Task `002-06-modelar-persistencia-repositories-e-rls.md`;
- reset local é reproduzível e não usa produção;
- roles, grants, contexto transacional e cleanup em pool estão documentados para a
  prova concreta da Task `002-06-modelar-persistencia-repositories-e-rls.md`;
- referências Supabase dev, Staging e Production são distintas e sem secrets
  versionados;
- schema CrownPilot não depende da Supabase Data API;
- pipeline de bundle/job está definido, mas o artefato final pertence à Task
  `002-11-automatizar-ci-oci-e-gates-de-release.md` e será gerado após as
  migrations da Task `002-06-modelar-persistencia-repositories-e-rls.md`;
- região está decidida como Render `us-east` + Supabase `us-east-1`; backups,
  subprocessadores, transferência internacional, residência e egress são gates
  antes de dados reais;
- nenhuma tabela futura ou payload de jogador é criada por conveniência.

## Riscos e cuidados

- Não usar `auth.uid()` de um provider diferente no RLS.
- Não rodar migrations automaticamente em múltiplas réplicas.
- Não duplicar schema entre EF Core e SQL.
- Não confundir projeto de banco Supabase com projeto Firebase.
- Não provisionar Production durante bootstrap local.

## Registro de execução

- **Status:** `completed`; pipeline local EF Core → SQL security foi executado
  com Docker PostgreSQL descartável.
- **Arquivos alterados:** `.env.example`, `CrownPilot.sln`,
  `frontend/scripts/db-gate.mjs`, `src/Api/Program.cs`,
  `src/Infrastructure/Infrastructure.csproj`,
  `src/Infrastructure/Persistence/`, `database/`,
  `tests/Persistence/` e `docs/operations/002-03-postgresql-migrations-rls.md`.
- **Decisões/desvios:** Docker Compose com `postgres:17.6-alpine` é o runner local
  pinado; Supabase CLI não foi introduzida como dependência adicional porque o
  contrato local é Docker + Npgsql e Supabase remoto será operado por migration
  job controlado. EF Core cria somente o schema `crownpilot`; tabelas e policies
  de domínio permanecem em `002-06`. `db:fixtures` falha fechado enquanto não
  houver fixtures de domínio.
- **Revisão independente:** apontou validação fraca da migration e aceitação de
  TLS sem validação de certificado. Corrigido: `db:security`/`db:fixtures` exigem
  migration conhecida, SQL não cria schema, e Staging/Production aceitam apenas
  `VerifyCA`/`VerifyFull` sem certificado confiado.
- **Comandos executados:** `dotnet build CrownPilot.sln --configuration Release`;
  `dotnet test tests/Persistence/Persistence.Tests.csproj --configuration Release
  --no-build`; `dotnet test CrownPilot.sln --configuration Release --no-build`;
  `dotnet ef migrations list --project src/Infrastructure/Infrastructure.csproj
  --startup-project src/Infrastructure/Infrastructure.csproj --configuration
  Release`; `npm run lint --prefix frontend`; `npm run typecheck --prefix frontend`;
  `npm run test:unit --prefix frontend`; `npm run build --prefix frontend`;
  `CROWNPILOT__ENVIRONMENT=Production npm run db:reset --prefix frontend`;
  `npm run db:start --prefix frontend`; `npm run db:migrate --prefix frontend`;
  `npm run db:security --prefix frontend`; `npm run db:fixtures --prefix frontend`;
  `npm run db:reset --prefix frontend`; segunda execução de
  `npm run db:security --prefix frontend` seguida de consulta PostgreSQL para
  roles/schema/history; `npm run db:stop --prefix frontend`.
- **Resultados/evidências:** build .NET passou sem warnings/erros; suíte .NET
  passou com 30 testes; Persistence passou com 10 testes; lint, typecheck, testes
  unitários frontend e build frontend passaram. `dotnet ef migrations list`
  descobriu `20261007000000_PrepareCrownPilotSchema` antes do banco local estar
  disponível. Com Docker ativo, `db:start`, `db:migrate` e `db:security` passaram;
  `db:fixtures` e etapa final de `db:reset` retornaram exit `2` deliberadamente,
  pois fixtures de domínio pertencem a `002-06`. Reexecução de security passou;
  consulta confirmou `crownpilot_runtime|false` e
  `crownpilot_migrator|false`, schema `crownpilot` e migration registrada.
  `db:reset` em Production recusou exit `2` antes de tocar Docker; `db:stop`
  passou.
- **Riscos residuais:** `db:fixtures` só será liberado após tabelas e fixtures
  sanitizadas de `002-06`; policies concretas, repositories, bridge A/B e pooler
  continuam fora desta subtarefa. A migration local não deve ser aplicada no
  startup de réplicas.
