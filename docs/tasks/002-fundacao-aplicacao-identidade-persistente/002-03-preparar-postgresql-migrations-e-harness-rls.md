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
  concretas pertencem à Task `002-06-modelar-persistencia-repositories-e-rls.md`;
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

1. Fixar Supabase CLI/Docker e definir o contrato cross-platform dos scripts npm
   de banco que o bootstrap fornecerá.
2. Configurar pipeline local: banco limpo, migrations EF Core, SQL RLS/grants,
   fixtures quando schema existir; nenhum SQL de segurança pode criar tabela do
   schema.
3. Validar connection strings Npgsql, SSL/configuração e fail-closed.
4. Preparar role de runtime sem `BYPASSRLS`, schema dedicado, runner de segurança
   e mecanismo transacional de contexto; não criar policies de tabelas ainda
   inexistentes. A Task `002-06-modelar-persistencia-repositories-e-rls.md` prova
   conexão direta/session pooler com schema concreto e registra transaction pooler
   como gate separado se for adotado.
5. Registrar referências de banco separadas para Local, Staging e Production; não
   commitar secrets.
6. Registrar critérios para `sa-east-1`, DPA, backups, subprocessadores e
   residência antes de dados reais.
7. Manter provisionamento Production como etapa controlada posterior.

## Pré-requisitos Mac/Linux

O setup local exige Git, .NET SDK e Node.js LTS/npm nas versões pinadas pela Task
`002-01-bootstrap-toolchain.md`, `dotnet-ef` como ferramenta local do repositório,
Docker (Docker Desktop no Mac ou Docker Engine/Compose no Linux) e Supabase CLI na
versão registrada.
Não é necessário instalar PostgreSQL no host. Os comandos usam somente banco local
descartável e não exigem credenciais de Staging ou Production.

## Interface oficial de banco

Os scripts npm são a interface pública e cross-platform. A Task
`002-01-bootstrap-toolchain.md` deve
fornecer estes nomes no `package.json`, sem adicionar task runner ou outra
dependência apenas para encadear comandos:

| Comando oficial | Contrato |
|---|---|
| `npm run db:start` | inicia PostgreSQL local descartável via Supabase CLI/Docker; falha se dependências locais não estiverem disponíveis |
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

O `package.json` e seus scripts ainda não existem porque a Task
`002-01-bootstrap-toolchain.md` não iniciou. Qualquer runner interno de shell é
detalhe de implementação encapsulado
por esses scripts; desenvolvedores e CI usam somente a interface `npm run`.

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
- referências Staging/Production são distintas e sem secrets versionados;
- schema CrownPilot não depende da Supabase Data API;
- pipeline de bundle/job está definido, mas o artefato final pertence à Task
  `002-11-automatizar-ci-oci-e-gates-de-release.md` e será gerado após as
  migrations da Task `002-06-modelar-persistencia-repositories-e-rls.md`;
- região, backups e residência são gates antes de dados reais;
- nenhuma tabela futura ou payload de jogador é criada por conveniência.

## Riscos e cuidados

- Não usar `auth.uid()` de um provider diferente no RLS.
- Não rodar migrations automaticamente em múltiplas réplicas.
- Não duplicar schema entre EF Core e SQL.
- Não confundir projeto de banco Supabase com projeto Firebase.
- Não provisionar Production durante bootstrap local.

## Registro de execução desta correção documental

- **Status:** `pending`; a Task `002-03-preparar-postgresql-migrations-e-harness-rls.md`
  não foi executada.
- **Arquivos alterados:** `docs/roadmap/crownpilot-roadmap.md`,
  `docs/specs/002-fundacao-aplicacao-identidade-persistente.md`,
  `docs/decisions/004-aspnet-core-react-vite-firebase-postgresql.md`,
  `docs/tasks/002-fundacao-aplicacao-identidade-persistente/002-00-overview.md`,
  `002-02-estabelecer-boundaries-contrato-base-e-ambientes.md`,
  `002-03-preparar-postgresql-migrations-e-harness-rls.md`,
  `002-04-implementar-google-sign-in-e-firebase-bearer.md`,
  `002-05-implementar-port-e-adapter-de-lookup.md`,
  `002-06-modelar-persistencia-repositories-e-rls.md`,
  `002-07-implementar-casos-de-uso-e-api-v1.md`,
  `002-08-entregar-frontend-de-identidade-e-vinculo.md`,
  `002-09-revisar-arquitetura-frontend-e-ux-visual.md`,
  `002-10-instrumentar-observabilidade-health-e-redaction.md`,
  `002-11-automatizar-ci-oci-e-gates-de-release.md` e
  `002-12-validar-staging-e2e-smoke-e-handoff.md`, todos em
  `docs/tasks/002-fundacao-aplicacao-identidade-persistente/`.
- **Decisões/desvios:** referências atuais da Fase 002 usam filenames canônicos;
  comandos de banco foram definidos como contrato `npm run db:*`; nenhum
  `package.json`, runtime ou runner executável foi criado porque
  `002-01-bootstrap-toolchain.md` ainda
  está pending. Shell runner permanece detalhe interno.
- **Comandos executados:** `git diff --check`; validação `rg` de referências de
  task e comandos shell; contagem de checklist com `rg`; verificador Python de
  filenames, links e targets; varredura global `rg` dos termos arquiteturais;
  inspeção de `git diff --stat`; verificador Python de estrutura SDD, ticker,
  checklist, filenames e links canônicos. Uma primeira checagem inline de links
  falhou por erro de sintaxe do próprio comando de validação; a versão corrigida
  passou. Os comandos oficiais `npm run db:*` não foram executados:
  `package.json` ainda não existe.
- **Resultados/evidências:** checks finais passaram; roadmap e spec têm 12
  links canônicos na ordem; overview tem um único checklist com 12 itens
  desmarcados e progresso `0/12`; nenhum comando direto de shell runner aparece
  nos docs da Fase 002; varredura global preserva ocorrências históricas,
  superseded, ADR antigo, Fase 001 e alternativas estáticas não conflitantes;
  todos os 12 arquivos têm ticker `002` e numeração correspondente. Revisão
  independente encontrou três referências abreviadas fora do conjunto inicial;
  foram normalizadas para filenames canônicos. Segunda revisão independente
  passou sem novos achados.
- **Riscos residuais:** scripts npm, migrations EF Core, SQL RLS/grants, fixtures,
  roles e testes ainda não existem; `002-03-preparar-postgresql-migrations-e-harness-rls.md`
  permanece bloqueada por `002-01-bootstrap-toolchain.md` e
  `002-02-estabelecer-boundaries-contrato-base-e-ambientes.md`.
