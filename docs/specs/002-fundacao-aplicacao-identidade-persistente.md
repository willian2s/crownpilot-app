# 002 — Fundação da aplicação e identidade persistente

- **Ticker:** `002`
- **Status:** `planned`
- **Roadmap:** [Fase 002](../roadmap/crownpilot-roadmap.md#002--fundacao-da-aplicacao-e-identidade-persistente)
- **ADR canônica:** [ADR 004](../decisions/004-aspnet-core-react-vite-firebase-postgresql.md)
- **Dependência:** Fase 001 — `GO WITH CONSTRAINTS / APPROVAL DEPENDENCY`

## Contexto e baseline

A Fase 001 liberou somente bootstrap reproduzível, identidade CrownPilot e
vínculo privado, read-only, de uma Player Tag que representa perfil público. O
vínculo usa `subjectType: public_profile` e `ownershipStatus: unverified`; não
prova que o usuário possui a conta consultada.

O baseline da `main`, verificado em 2026-10-02, contém documentação e arquivos de
configuração, mas não contém runtime, dependências instaladas, testes, CI ou
dados de produção. Não há comportamento de aplicação para preservar. A
documentação anterior da Fase 002 foi substituída pela [ADR 004](../decisions/004-aspnet-core-react-vite-firebase-postgresql.md).

## Problema

Sem fundação executável, não há forma reproduzível de:

1. autenticar uma identidade CrownPilot sem pedir credenciais Supercell;
2. guardar uma Player Tag primária e recuperá-la em outro dispositivo;
3. validar a tag no servidor sem expor token externo;
4. trocar ou remover o vínculo sem perder o vínculo anterior em caso de falha;
5. separar authentication, authorization, domínio, persistência e providers;
6. provar segurança, ambientes, migrations, testes e deploy antes do sync.

## Objetivo

Preparar uma aplicação web portátil com React + TypeScript + Vite no frontend e
ASP.NET Core + C# no backend, usando Firebase Authentication para Google
Sign-In e PostgreSQL hospedado no Supabase através de EF Core + Npgsql.

Ao fim da fase, um usuário deve conseguir:

```text
React/Vite
    -> Firebase Authentication / Google Sign-In
    -> Firebase ID Token
    -> ASP.NET Core API
    -> identificação e autorização CrownPilot
    -> lookup server-side de perfil público
    -> EF Core + Npgsql -> PostgreSQL/Supabase
    -> salvar vínculo public_profile/unverified
```

Falha do provider externo não pode apagar ou invalidar identidade local nem
vínculo anterior.

## Não objetivos e limites obrigatórios

Fase 002 não implementa nem libera:

- sync ou persistência de Player Snapshot completo;
- coleção, níveis, Arena, troféus, Evolutions, Heroes ou battle history;
- cache persistente, retenção ou redistribuição de payload externo;
- ownership verificado, exclusividade, ações sobre conta ou notificações;
- recommendation engine, readiness ou Upgrade Planner;
- polling, crawler, ingestão global, jobs, filas ou scheduler;
- billing, assinatura, paywall, premium feature ou AI Coach;
- credencial ou senha Supercell;
- acesso direto do browser ao PostgreSQL/Data API para dados CrownPilot;
- dependência obrigatória de APIs proprietárias de hospedagem;
- chamadas live da Clash Royale API na CI normal.

O lookup pode confirmar que o perfil público existe e retornar resultado
operacional mínimo. Nome, coleção, contexto competitivo, deck e payload raw não
entram no documento de vínculo.

## Requisitos consolidados

### Stack e bootstrap

- React + TypeScript + Vite no frontend, separado da API;
- ASP.NET Core Web API em C# no backend;
- solution e SDK .NET versionados/pinados;
- camadas `API`, `Application`, `Domain` e `Infrastructure`;
- EF Core + Npgsql para PostgreSQL;
- Docker/OCI para a API, sem runtime dependente de Vercel;
- TypeScript strict, lint, format, type-check, unit, integration, contract,
  persistence, authentication/authorization, E2E, smoke e build com comandos
  documentados;
- `dotnet restore`, `dotnet build`, `dotnet test`, `npm ci`, `npm run build` e
  scripts equivalentes reproduzíveis em clone limpo;
- CI em Pull Requests e push para `main`, sem secrets reais.

### Authentication, authorization e ambientes

- Firebase Authentication com Google Sign-In é responsável por autenticação e
  emissão/refresh do Firebase ID Token;
- ASP.NET Core valida o token no servidor e deriva o Firebase UID somente de
  claims verificadas;
- authorization de recursos é responsabilidade da API/Application, não do
  frontend;
- frontend envia `Authorization: Bearer <Firebase ID Token>` por HTTPS;
- não há troca por cookie nem sessão paralela nesta fase;
- Firebase e Supabase têm projetos separados por ambiente;
- local usa Firebase Emulator ou fixtures de token e PostgreSQL descartável;
- Preview efêmero não depende de URL fixa ou login real;
- Staging usa hostname fixo, Google Sign-In real e banco separado;
- Production usa configuração, secrets, Firebase e banco próprios;
- Preview não recebe dados/secrets de Staging ou Production.

### Persistência

- PostgreSQL no Supabase é infraestrutura de banco, não backend da aplicação;
- EF Core é dono das tabelas, colunas, constraints, índices e histórico de
  migrations da aplicação;
- Npgsql é o provider de acesso PostgreSQL;
- SQL separado versiona somente RLS, grants, roles, extensões e objetos de
  plataforma, sem duplicar schema do EF Core;
- aplicação não usa SDK C# do Supabase como dependência central;
- migrations de produção são job controlado, nunca startup automático em
  múltiplas réplicas;
- execução documentada: migrations EF Core, depois SQL de segurança/RLS, depois
  fixtures sanitizadas.

### Identidade e vínculo

- Firebase UID é subject externo, não ID de domínio;
- `CrownPilotUserId` interno é estável e independente do provider;
- uma Player Tag primária por usuário CrownPilot;
- normalização conservadora: `trim`, exatamente um `#` inicial, `%23` somente no
  path HTTP e preservação de case/restante;
- lookup somente pelo port `IClashRoyaleClient` server-side;
- somente `resolved` cria ou substitui vínculo;
- estados `invalid_input`, `not_found`, `provider_unavailable`, `rate_limited` e
  `provider_misconfigured` são distinguíveis;
- troca valida o novo perfil antes de substituir o antigo;
- desvinculação e exclusão de dados próprios são explícitas e idempotentes;
- UI comunica “Perfil público salvo — ownership não verificado”.

## Comportamento atual encontrado

1. O repositório não possui aplicação, comandos executáveis, testes ou CI; a
   Fase 001 registra esse baseline.
2. Existem ADRs históricas e uma especificação anterior com stack divergente; a
   ADR 004 é a autoridade atual para a Fase 002.
3. A Fase 001 observou lookup operacional via proxy, mas rota oficial direta,
   termos, limites, SLA, retenção, ownership e egress permanecem condicionais
   ou `UNRESOLVED`.
4. O contrato de dados v0 exige provenance/freshness e ausência explícita, mas
   não é schema persistente e não deve virar snapshot nesta fase.
5. `.env.example` exemplifica somente configuração de integração externa; não há
   credenciais reais versionadas.

## Abordagem escolhida

### Arquitetura em camadas

```text
API
  HTTP, DTOs, ProblemDetails, auth pipeline, endpoints e composição
Application
  casos de uso, autorização de aplicação e ports
Domain
  entidades, value objects, invariantes e resultados sem dependências externas
Infrastructure
  Firebase, EF Core, Npgsql, PostgreSQL, Clash Royale e observabilidade
```

`Domain` não conhece Firebase, Supabase, EF Core, Npgsql, HTTP, ASP.NET Core ou
Vercel. `Application` depende de abstrações e não de SDKs de providers. `API`
converte HTTP para casos de uso; controllers/endpoints devem permanecer finos.

### Authentication e authorization

```text
React Firebase SDK
    -> Google Sign-In / refresh / logout
    -> Firebase ID Token
    -> Authorization: Bearer ...
    -> ASP.NET Core authentication middleware
    -> Firebase UID verificado
    -> CrownPilot User interno
    -> authorization por recurso
```

Validar assinatura, issuer, audience/project ID, expiração, claims temporais,
`sub`, `kid` e rotação de chaves por integração suportada na Infrastructure.
Não implementar validação JWT manual. Nunca aceitar UID de body, query ou header
arbitrário. A API responde `401`, sem redirecionamento.

### Persistência e migrations

Schema inicial conceitual:

```sql
crownpilot_users
  id uuid primary key
  firebase_uid text unique not null
  schema_version text not null
  created_at timestamptz not null
  updated_at timestamptz not null

primary_player_links
  user_id uuid primary key references crownpilot_users(id)
  schema_version text not null
  player_tag text not null
  subject_type text not null check = 'public_profile'
  ownership_status text not null check = 'unverified'
  version bigint not null default 1
  linked_at timestamptz not null
  updated_at timestamptz not null
  last_validated_at timestamptz null
```

`firebase_uid` é referência externa e não chave primária de domínio. Não há
referência a tabelas de identidade gerenciadas pelo provider. `player_tag` não
possui unicidade global; a mesma tag pode ser referência pública de vários
usuários.

`version` é token de concorrência otimista da linha de vínculo; não é a versão
do schema. Read retorna `version`/ETag e replace exige precondition da versão
lida. Mismatch retorna `409` sem sobrescrever o vínculo. Primeiro vínculo usa
precondition ausente; corrida nessa criação também retorna `409` e pode ser
repetida após nova leitura.

EF Core gera e aplica migrations de schema. SQL de RLS/grants deve ser aplicado
depois, com runner explícito, e não pode recriar tabelas. RLS é defesa adicional:
authorization ASP.NET Core continua obrigatória. Policies não podem presumir
claims nativas de outro provider; qualquer contexto de usuário usado por RLS
precisa ser estabelecido pelo backend e validado no pooler escolhido. O desenho
previsto usa role de runtime sem `BYPASSRLS`, `SET LOCAL
app.crownpilot_user_id = <CrownPilotUserId>` dentro da transação autenticada e
policies que negam quando o contexto está ausente. O contexto expira com a
transação; conexão de migration/admin é separada. Se o pooler não preservar essa
bridge, RLS continua deny-by-default, mas não pode ser alegado como isolamento
por usuário até existir solução comprovada.

### Deploy e ambientes

Frontend Vite gera `dist/` e pode ser hospedado em Vercel, CDN, Nginx ou storage
estático. A API gera imagem Docker/OCI ASP.NET Core, com secrets em runtime,
filesystem efêmero e `/health/live` e `/health/ready`. Nenhum componente de
backend depende de Vercel.

| Ambiente | API/frontend | Auth | Banco | Validação |
|---|---|---|---|---|
| Local | Vite + API local/container | Emulator/fixtures | Supabase CLI/Docker descartável | unit/integration |
| Preview | build efêmero, API opcional | sem login real | sem produção | build/UI/smoke |
| Staging | hostname fixo | Firebase separado + Google real | Supabase separado | E2E/smoke |
| Production | promoção controlada via `main` | projeto próprio | projeto próprio | smoke controlado |

## Contratos e estados

### Identity

```text
AuthenticatedSubject { firebaseUid: string }       # somente boundary de auth
AuthContext { crownpilotUserId: CrownPilotUserId } # application/domain
```

`AuthenticatedSubject` nunca é aceito diretamente em operações de domínio. O
mapping para `CrownPilotUserId` ocorre no backend após validação do token.
`EnsureCrownPilotUser` cria/recupera a linha por `firebase_uid` de forma
idempotente; unique conflict de primeiro login é relido, não gera usuário
duplicado. Falha de banco não autentica parcialmente o request.

```text
PrimaryPlayerLink {
  playerTag: NormalizedPlayerTag
  subjectType: "public_profile"
  ownershipStatus: "unverified"
  linkedAt: Instant
  updatedAt: Instant
  lastValidatedAt?: Instant
  schemaVersion: string
}
```

Não persistir e-mail, nome Google, nome do perfil externo ou payload do provider
sem requisito aprovado.

### Lookup

```text
resolvePublicProfile(playerTag):
  | { kind: "resolved"; normalizedPlayerTag: string; validatedAt: Instant }
  | { kind: "invalid_input" }
  | { kind: "not_found" }
  | { kind: "provider_unavailable" }
  | { kind: "rate_limited"; retryAfterSeconds?: number }
  | { kind: "provider_misconfigured" }
```

Somente `resolved` persiste vínculo. O resultado não contém snapshot, coleção,
Arena, deck ou nome.

### HTTP

Os paths concretos pertencem à API ASP.NET Core e retornam JSON/ProblemDetails.

| Operação | Auth | Resultado |
|---|---|---|
| Ler vínculo primário | obrigatória | vínculo do usuário atual ou ausência |
| Criar/substituir | obrigatória | lookup antes da escrita; vínculo antigo preservado em falha |
| Desvincular | obrigatória | remoção idempotente do usuário atual |
| Excluir dados CrownPilot | obrigatória/reautenticação se aprovada | remove somente dados próprios |

Mapeamento mínimo: `401` token ausente/inválido, `403` token válido sem
autorização para a operação, `400` input inválido, `404` perfil inexistente ou
ausência do vínculo próprio, `409` precondition/concurrency, `429` rate limit e
`503` provider indisponível ou mal configurado. Usuário A nunca recebe dados de
B. Não vazar token, URL, rota externa, IAM, payload bruto ou detalhes internos.

## Arquivos, módulos e contratos afetados

Arquivos esperados após implementação, sujeitos ao bootstrap:

- `global.json`, solution e projetos `.csproj`;
- `src/Api/`, `src/Application/`, `src/Domain/`, `src/Infrastructure/`;
- `src/Infrastructure/Persistence/Migrations/` para migrations EF Core;
- SQL versionado de RLS/grants, sem schema duplicado;
- `frontend/`, `vite.config.ts`, `tsconfig.json`, `package.json` e lockfile;
- `Dockerfile`, `.dockerignore`, `appsettings*.json` sem secrets;
- `.github/workflows/ci.yml`, `.env.example` e documentação operacional;
- `tests/Unit`, `tests/Application`, `tests/Integration`, `tests/Persistence`,
  `tests/Contract`, `tests/Api`, `tests/E2E` e `tests/Smoke`.

Símbolos prováveis: `CrownPilotUserId`, `PrimaryPlayerLink`,
`IClashRoyaleClient`, `IUserRepository`, `IPlayerLinkRepository`,
`FirebaseAuthentication`, `CrownPilotDbContext`, handlers de authorization,
controllers/endpoints de vínculo e `ProblemDetails` mapping.

## Alternativas descartadas

| Alternativa | Motivo |
|---|---|
| Backend server-driven acoplado ao frontend | A decisão oficial exige React/Vite consumindo API ASP.NET Core. |
| Next.js | Não é frontend oficial; Vite permanece a ferramenta escolhida. |
| Supabase como backend completo | Supabase fornece PostgreSQL; regras, API e autorização pertencem ao backend CrownPilot. |
| SDK C# do Supabase como integração central | Acesso deve ser PostgreSQL via EF Core + Npgsql. |
| Supabase Auth | Firebase Auth/Google é a decisão oficial de identidade. |
| Cookie e bearer simultâneos | Duplica lifecycle e superfície de segurança sem requisito. |
| EF Core e SQL duplicando schema | Drift e rollback ambíguo. |
| Firebase UID como ID interno | Acoplamento do domínio ao provider externo. |
| Backend obrigatório na Vercel | Contradiz portabilidade e limita escolha operacional. |

## Critérios de aceite verificáveis

- [ ] clone limpo restaura solution .NET e frontend, compila e executa testes;
- [ ] Task 01 contém ASP.NET Core, React, TypeScript, Vite, EF Core, Npgsql,
      Docker, Node/npm pinados, `AGENTS.md` e toolchain de testes, sem Auth real
      ou provisionamento;
- [ ] Domain não referencia Firebase, Supabase, EF Core, Npgsql, HTTP, ASP.NET
      Core ou Vercel;
- [ ] Application usa ports/interfaces, sem SDK direto de provider;
- [ ] API expõe JSON/ProblemDetails e valida autenticação antes da autorização;
- [ ] React usa Firebase Auth/Google e envia somente Firebase ID Token bearer;
- [ ] backend valida issuer, audience, expiração, assinatura, `kid` e project ID
      por ambiente, sem aceitar UID arbitrário;
- [ ] authentication e authorization possuem testes separados;
- [ ] `crownpilot_users` usa UUID interno e `firebase_uid` único;
- [ ] vínculo referencia o ID interno e não qualquer identidade do provider;
- [ ] `subject_type` e `ownership_status` são `NOT NULL` com checks invariantes;
- [ ] replace usa token/version precondition e retorna `409` em conflito;
- [ ] bridge RLS define role, contexto, reset e testes sem contexto, A/B e usuário correto;
- [ ] EF Core migrations criam schema; SQL posterior cobre somente RLS/grants;
- [ ] migrations não são aplicadas automaticamente por múltiplas réplicas;
- [ ] acesso A/B é autorizado no backend e RLS é testado como defesa adicional;
- [ ] local, Preview, Staging e Production não compartilham dados/secrets;
- [ ] workflow protegido/manual de Staging possui owner, aprovação e smoke antes
      da promoção para `main`;
- [ ] Preview não depende de hostname fixo ou login real; Staging possui hostname
      fixo, Firebase separado e Google Sign-In funcional;
- [ ] frontend pode ser hospedado fora de Vercel e API inicia fora de Vercel;
- [ ] `IClashRoyaleClient` é o único boundary de lookup e não aceita host do usuário;
- [ ] somente `resolved` cria/substitui vínculo `public_profile/unverified`;
- [ ] falha externa não remove vínculo anterior; replace concorrente é protegido;
- [ ] unlink/delete são explícitos e idempotentes;
- [ ] nenhum token, secret, Player Tag real ou payload externo aparece em bundle/log;
- [ ] E2E Staging cobre login, vínculo, reload/outro dispositivo, replace e unlink;
- [ ] Fase 003 continua bloqueada pelos gates de API data, retenção, ownership,
      egress, meta e compliance da Fase 001.

## Estratégia de testes e validação

- **Unit:** xUnit ou runner .NET equivalente para invariantes de domínio,
  normalização, estados, concorrência e redaction; Vitest/React Testing Library
  para frontend;
- **Application/use-case:** casos de link, replace, unlink, delete e mapping de
  autorização com ports fake;
- **Integration:** `WebApplicationFactory`/host ASP.NET e PostgreSQL local;
- **Persistence:** EF Core migrations, Npgsql, constraints, transações e RLS;
- **Contract:** fixtures sanitizadas do `IClashRoyaleClient`, sem rede live;
- **Authentication/authorization:** token ausente, inválido, expirado, issuer/
  audience errados, project errado, A/B e UID adulterado;
- **E2E:** Playwright ou equivalente em Staging com Google Sign-In real;
- **Smoke:** health, configuração, build e fluxo mínimo pós-deploy.

Gates de PR: restore/build/analyzers .NET, `dotnet test`, npm ci, lint,
type-check, testes frontend, migrations contra PostgreSQL descartável, RLS,
contract e build. Push em `main` adiciona imagem e smoke. Staging adiciona
Firebase/Supabase reais e E2E. Nenhum gate normal chama Clash Royale live ou usa
produção.

## Observabilidade

Registrar somente categoria de resultado, latência, status, ambiente,
correlation/request ID e versão da aplicação. Redaction obrigatória para Player
Tag, e-mail, UID, token, URL com tag, payload externo, senha e secrets.

## Riscos e mitigação

| Risco | Mitigação |
|---|---|
| issuer/audience Firebase incorretos | allowlist por ambiente, validação server-side e testes de rotação |
| confundir Firebase com identidade de domínio | UUID interno e mapping exclusivo no backend |
| drift entre EF e SQL | EF é dono do schema; SQL somente RLS/grants; gate de ordem |
| RLS incompatível com pooler | testar bridge de contexto; não alegar isolamento não comprovado |
| Preview atingir produção | sem secrets reais, allowlists e startup fail-closed |
| Vercel virar dependência backend | frontend estático opcional e API OCI portátil |
| replace concorrente perder vínculo | transação/precondition e resposta `409` |
| lookup externo indisponível | port, fixtures, retry limitado e preservação do vínculo |
| exposição de dados/segredos | server-only, HTTPS, CORS exato, redaction e scans |

## Rollout, rollback e migração

1. Bootstrap local sem provisionar Production.
2. Fixar boundaries, ambientes e contratos.
3. Configurar PostgreSQL/Supabase local, projetos separados, EF migrations e RLS.
4. Implementar e testar Firebase bearer authentication e autorização.
5. Implementar lookup e persistência mínima.
6. Entregar API/React e fluxos de vínculo.
7. Fechar CI, observabilidade e smoke de container.
8. Validar Staging com Google Sign-In real, E2E e smoke.
9. Promover Production somente por `main`, com migration job controlado.

Não há migração de dados legados nem dual-write. Rollback de aplicação deve
preservar schema e vínculo; mudanças incompatíveis usam expand/contract.

## Ordem das subtarefas

1. [002-01 — bootstrap do toolchain](../tasks/002-fundacao-aplicacao-identidade-persistente/002-01-bootstrap-toolchain.md)
2. [002-02 — boundaries e ambientes](../tasks/002-fundacao-aplicacao-identidade-persistente/002-02-fixar-boundaries-e-ambientes.md)
3. [002-03 — PostgreSQL/Supabase local, Firebase por ambiente e migrations](../tasks/002-fundacao-aplicacao-identidade-persistente/002-03-configurar-supabase-local-e-projetos.md)
4. [002-04 — Firebase Auth bearer](../tasks/002-fundacao-aplicacao-identidade-persistente/002-04-implementar-auth-google-e-sessao.md)
5. [002-05 — adapter de lookup](../tasks/002-fundacao-aplicacao-identidade-persistente/002-05-criar-adapter-de-lookup.md)
6. [002-06 — persistência e autorização](../tasks/002-fundacao-aplicacao-identidade-persistente/002-06-persistir-vinculo-com-autorizacao.md)
7. [002-07 — fluxos de vínculo e exclusão](../tasks/002-fundacao-aplicacao-identidade-persistente/002-07-entregar-fluxos-de-vinculo-e-exclusao.md)
8. [002-08 — quality gates e observabilidade](../tasks/002-fundacao-aplicacao-identidade-persistente/002-08-automatizar-quality-gates-e-observabilidade.md)
9. [002-09 — Staging, deploy, E2E e handoff](../tasks/002-fundacao-aplicacao-identidade-persistente/002-09-validar-staging-deploy-e2e-smoke-handoff.md)

`002-05` pode ser implementada em paralelo após `002-01` e `002-02`, pois é
fixture-driven e não depende de authentication. `002-08` deve deixar gates
mínimos no bootstrap, embora seu fechamento dependa das features anteriores.

## Premissas explícitas

- ticker `002` e slug existentes permanecem; não há nova task de implementação;
- React + TypeScript + Vite continuam oficiais;
- Firebase Authentication com Google é o único provider de login desta fase;
- Supabase hospeda PostgreSQL, mas não é backend da aplicação nem provider de
  identidade;
- EF Core + Npgsql são a integração oficial de persistência;
- bearer Firebase é o único mecanismo de autenticação API nesta fase;
- uma Player Tag primária basta para o MVP;
- usuário fornece a tag; não há credential exchange com Supercell;
- `sa-east-1` permanece preferência condicionada à disponibilidade e revisão de
  residência, backups e subprocessadores;
- Vercel pode hospedar frontend estático, mas não é requisito do backend;
- ferramentas específicas só entram quando justificadas pelo bootstrap;
- não há dados legados a migrar e nenhum deploy existente a preservar.

## Handoff esperado

Entregar aplicação reproduzível, API portátil, frontend independente, identidade
estável, PostgreSQL seguro, migrations EF Core/RLS testadas, vínculo primário
removível e `unverified`, lookup server-side com fixtures, CI/E2E/smoke com
evidência e observabilidade sem dados sensíveis. A Fase 003 só pode começar após
reabrir e aprovar os gates de API data, retenção, ownership, egress, meta e
compliance deixados pela Fase 001.
