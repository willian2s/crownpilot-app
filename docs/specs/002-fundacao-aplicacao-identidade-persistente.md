# 002 — Fundação da aplicação e identidade persistente

- **Ticker:** `002`
- **Status:** `planned`
- **Roadmap:** [Fase 002](../roadmap/crownpilot-roadmap.md#002--fundacao-da-aplicacao-e-identidade-persistente)
- **Dependência:** Fase 001 — `GO WITH CONSTRAINTS / APPROVAL DEPENDENCY`

## Contexto e baseline

A Fase 001 liberou a Fase 002 somente para bootstrap reproduzível, identidade
CrownPilot e vínculo privado, read-only, de uma Player Tag que representa um
perfil público. O vínculo deve usar `subjectType: public_profile` e
`ownershipStatus: unverified`; não prova que o usuário possui a conta consultada.

O baseline real da `main`, verificado em `2026-10-02`, contém documentação,
`.env.example` e `.gitignore`, mas não contém runtime, `package.json`, lockfile,
framework, testes, CI, configuração Firebase ou Vercel. Não há comportamento de
produção para preservar.

Evidências principais:

- [baseline da Fase 001](001-viabilidade-produto-dados-compliance.md#comportamento-atual-encontrado-e-baseline-da-main);
- [ADR de infraestrutura](../decisions/001-firebase-firestore-vercel-portable.md);
- [ADR 002 de stack, deploy e região](../decisions/002-laravel-inertia-react-vite-firestore-region.md);
- [ADR 003 de Supabase Auth, PostgreSQL, JWT/JWKS, RLS e região](../decisions/003-supabase-auth-postgresql-jwks-rls-region.md);
- [veredito e handoff da Fase 001](../tasks/001-viabilidade-produto-dados-compliance/evidences/phase-001-verdict.md#dependencias-da-fase-002);
- [contratos conceituais v0](../tasks/001-viabilidade-produto-dados-compliance/evidences/data-contract-v0.md#handoff-to-next-phases).

## Problema

O produto ainda não executa. Sem fundação, não há forma reproduzível de:

1. autenticar uma identidade CrownPilot sem pedir credenciais Supercell;
2. guardar uma Player Tag primária e recuperá-la em outro dispositivo;
3. validar a tag no servidor sem expor token externo;
4. trocar ou remover o vínculo sem perder o vínculo anterior em caso de falha;
5. separar autenticação, autorização, dados do usuário e resposta de provider;
6. validar segurança, ambientes e deploy antes de iniciar sync da conta.

## Objetivo

Entregar uma aplicação web executável, reproduzível e segura por padrão, com
Supabase Auth via Google, PostgreSQL isolado por ambiente, quality gates e
vínculo persistente de uma Player Tag primária.

Ao fim da fase, um usuário deve conseguir:

```text
Google Sign-In
    -> sessão CrownPilot
    -> informar Player Tag
    -> lookup server-side de perfil público
    -> salvar vínculo privado como public_profile/unverified
    -> recuperar, trocar ou remover vínculo em outro dispositivo
```

Falhas do provider externo não podem apagar ou invalidar a identidade local nem
o vínculo anterior.

## Não objetivos e limites obrigatórios

Fase 002 não implementa nem libera:

- sync ou persistência de Player Snapshot completo;
- coleção, níveis, Arena, troféus, Evolutions, Heroes ou battle history como
  estado persistido;
- cache persistente, retenção ou redistribuição de payload da API;
- ownership verificado, exclusividade, ações sobre conta ou notificações;
- meta da Arena, recommendation engine, readiness ou Upgrade Planner;
- polling, crawler, ingestão global, jobs, filas ou scheduler;
- billing, assinatura, paywall, premium feature ou AI Coach;
- credencial ou senha Supercell;
- dependência obrigatória de Vercel KV, Blob, Cron, Queue, Workflow, Edge
  Config ou API proprietária de runtime;
- chamadas live da Clash Royale API na CI normal.

O lookup pode confirmar que o perfil público existe e retornar apenas resultado
operacional mínimo ao fluxo. Resposta raw, nome, coleção, contexto competitivo e
deck não entram no documento de vínculo.

## Requisitos consolidados

### Bootstrap

- framework web escolhido e bootstrap reproduzível a partir de clone limpo;
- Laravel + Inertia + React + TypeScript + Vite;
- PHP 8.5 com Composer e `composer.lock` para dependências PHP;
- npm e `package-lock.json` para dependências JavaScript e build Vite;
- TypeScript strict, lint, format, type-check, unit, integration, contract,
  Rules, E2E, smoke e build com comandos documentados;
- `Dockerfile.vercel`, `Caddyfile`, `vercel.json` e `.dockerignore` para deploy
  Laravel containerizado com FrankenPHP;
- estrutura que mantenha domínio/aplicação independentes de Laravel, Inertia,
  Supabase e APIs Vercel;
- `AGENTS.md`, `.env.example` e comandos operacionais versionados;
- GitHub Actions executando gates em Pull Requests e também em todo push para
  `main`.

### Supabase, segurança e ambientes

- Supabase Auth com Google para a conta CrownPilot;
- Supabase PostgreSQL como banco principal;
- região Supabase `sa-east-1` (São Paulo), condicionada à disponibilidade no
  plano/organização, igual em staging e production;
- projetos Supabase separados para staging e production;
- Supabase CLI/Docker no desenvolvimento/testes;
- preview efêmero sem Google Sign-In real e sem Supabase de produção;
- staging com hostname fixo, Google Sign-In real e projeto Supabase próprio;
- production somente pela `main`;
- adapter PHP server-side para validação Supabase JWT/JWKS e PostgreSQL/Eloquent;
- `@supabase/supabase-js` limitado ao login, refresh e logout;
- browser não acessa Data API/PostgREST para dados CrownPilot;
- RLS e grants versionados, deny-by-default e testados;
- secrets nunca presentes em bundle, logs, Git ou configuração `VITE_*`.

### Identidade e vínculo

- uma Player Tag primária por usuário CrownPilot;
- `crownpilotUserId` derivado do UUID `sub` do Supabase JWT verificado no servidor;
- normalização conservadora: `trim`, exatamente um `#` inicial, `%23` apenas no
  path HTTP e preservação de case/restante;
- lookup somente pelo adapter `ClashRoyaleClient` server-side;
- somente resultado `resolved` pode criar ou substituir vínculo;
- estados `invalid_input`, `not_found`, `provider_unavailable`, `rate_limited` e
  `provider_misconfigured` distinguíveis;
- troca valida o novo perfil antes de substituir o antigo;
- desvinculação explícita e idempotente;
- exclusão inicial dos dados próprios do usuário, sem alegar exclusão de dados
  no provider externo;
- UI com mensagem equivalente a “Perfil público salvo — ownership não
  verificado”.

## Comportamento atual encontrado

1. O repositório não possui aplicação nem comandos executáveis. A Fase 001
   registra isso em seu baseline e a inspeção atual confirma ausência de
   `package.json`, lockfile, CI, testes e configs de Firebase/Vercel.
2. A [ADR 001](../decisions/001-firebase-firestore-vercel-portable.md#decisão)
   registrava Firebase Auth Google, Firestore e Vercel como escolhas-base. A
   [ADR 003](../decisions/003-supabase-auth-postgresql-jwks-rls-region.md)
   supersede Auth/Firestore para a Fase 002 e fixa Supabase Auth/PostgreSQL,
   JWT/JWKS, RLS e região condicionada.
3. A Fase 001 observou lookup operacional via proxy, mas rota oficial direta,
   termos do proxy, limites, SLA, retenção, ownership e egress permanecem
   condicionais ou `UNRESOLVED`.
4. O contrato v0 exige provenance/freshness e ausência explícita, mas não é
   schema persistente. Não pode ser promovido a documento de snapshot na Fase
   002.
5. O único segredo exemplificado atualmente é `CLASH_ROYALE_API_TOKEN` em
   `.env.example`; `.env.local` existe localmente e não é versionado.

## Abordagem escolhida

### Stack e runtime

Usar Laravel como backend HTTP, Inertia como integração server-driven, React +
TypeScript como frontend e Vite como build/dev server. PHP 8.5 usa Composer
e `composer.lock`; frontend usa npm e `package-lock.json`.

Controllers, Form Requests e respostas Inertia devem ser adapters finos. Regras
de domínio não importam Laravel, Inertia, Supabase ou Vite. O frontend compilado
fica em `resources/js/` e não cria backend separado.

Runtime PHP serverless/host PHP first-class fica atrás de gate. Vercel continua
alvo operacional, mas Laravel só pode ser liberado nela após provar runtime PHP,
extensões, sessão, Inertia, Auth e PostgreSQL. Se Vercel falhar, migrar o host
Laravel sem alterar domínio, contratos ou frontend.

### Deploy Vercel via Docker

Usar o processo oficial [Deploy Laravel on Vercel with Docker](https://vercel.com/kb/guide/laravel-php-with-docker):

- `Dockerfile.vercel` baseado em FrankenPHP/PHP 8.5, com Composer e extensões
  necessárias;
- `Caddyfile` expondo somente `public/` e encaminhando rotas para
  `public/index.php`;
- `vercel.json` declarando service container e rewrite catch-all;
- `.dockerignore` excluindo `.env*`, `vendor`, `node_modules`, testes e arquivos
  locais da imagem;
- configuração injetada por environment variables da Vercel em runtime;
- `APP_KEY` estável e separado por ambiente, nunca versionado ou compartilhado;
- filesystem do container tratado como efêmero; vínculo persistente fica no
  PostgreSQL/Supabase, não em arquivos locais;
- sessão/cache não dependem de arquivos locais nem de estado compartilhado no
  processo. Sessão cookie criptografada é opção inicial compatível com o escopo;
- endpoint `/health` usado por smoke e validação de container.

O deploy pode usar integração Git da Vercel ou CLI (`vercel deploy`/`vercel deploy
--prod`), mas produção permanece restrita à `main`. GitHub Actions é CI
autoritativo: roda gates em Pull Requests e repete gates, build Docker e smoke
de container em push para `main`.

### Boundary de autenticação e autorização

No browser, `@supabase/supabase-js` executa Google Sign-In, refresh e logout. O
React envia Supabase access JWT a endpoint Laravel de troca; o backend valida
assinatura via JWKS, `iss`, `aud`, `exp`, `nbf`, `alg`, `kid` e project ref, deriva
o UUID `sub` e cria sessão Laravel segura. Google ID token não é enviado ao
Laravel nem ao provider Clash Royale.

No servidor, adapter PHP valida JWT/JWKS e acessa PostgreSQL via Eloquent. O
backend aplica autorização usando UID da sessão verificada e `user_id` derivado
dele. Autenticação não concede acesso a outro usuário. RLS é defesa adicional;
conexão privilegiada pode ignorá-la, portanto autorização Laravel é obrigatória.

O browser não acessa Data API/PostgREST para dados CrownPilot. RLS fica
versionado, habilitado e deny-by-default para evitar abertura acidental.

### Boundary da Clash Royale API

Criar um port de aplicação `ClashRoyaleClient` com operação estreita
`resolvePublicProfile(playerTag)`. A implementação HTTP conhece provider route,
token, timeout, retry e egress; domínio e UI conhecem somente resultado
discriminado.

O proxy observado na Fase 001 continua configuração operacional substituível,
não contrato oficial. Host e rota são configuração do servidor, nunca input do
usuário. CI usa fixtures; probes live são explícitos e limitados.

### Persistência mínima

`supabase/migrations/` é fonte única de schema e RLS. Eloquent fornece models,
casts, queries e repositories; não duplicar migrations equivalentes em
`database/migrations/`.

Persistir somente o mínimo necessário para identidade e vínculo:

```sql
crownpilot_users
  id uuid primary key references auth.users(id)
  schema_version text
  created_at timestamptz
  updated_at timestamptz

primary_player_links
  user_id uuid primary key references crownpilot_users(id)
  schema_version text
  player_tag text
  subject_type text check = 'public_profile'
  ownership_status text check = 'unverified'
  linked_at timestamptz
  updated_at timestamptz
  last_validated_at timestamptz null
```

`last_validated_at` registra validação local bem-sucedida; não é freshness
garantida pelo provider. Não criar `players` global, snapshot, cards, battlelog,
raw API, cache ou unicidade global de Player Tag. `auth.users` é gerenciado pelo
Supabase; aplicação não escreve diretamente nele.

O vínculo não possui unicidade global de Player Tag. A mesma tag pode estar
associada a usuários CrownPilot diferentes, sem exclusividade ou ownership.

### Região e ambientes

Supabase usará `sa-east-1` (South America, São Paulo) em staging e production,
somente após confirmar disponibilidade no plano/organização. Não usar região
genérica `Americas` como fallback silencioso. Mudança futura exige novo projeto e
cutover controlado.

| Ambiente | Supabase/Auth | API externa | Dados |
|---|---|---|---|
| Local | Supabase CLI/Docker | mock/fixture por padrão | descartáveis |
| Preview | nenhum Supabase real | mock/desabilitada | nenhum persistente |
| Staging | projeto separado | secret próprio, smoke limitado | teste |
| Production | projeto separado | configuração própria após gates | reais |

Configuração ausente ou com project ref/issuer incompatível deve falhar fechado.
Preview nunca recebe secrets de staging/production. Redirect allowlist de staging
usa somente hostname fixo; OAuth real não é configurado para hosts efêmeros.

Vercel permanece alvo inicial, mas deploy Laravel/Vercel é gate técnico. Se o
runtime PHP, `pdo_pgsql`, SSL, pooler, sessão ou PostgreSQL não funcionarem de
forma reproduzível, usar host PHP first-class para backend e preservar Vercel
somente onde houver compatibilidade comprovada.

## Contratos e estados

Os nomes abaixo são contratos de Fase 002, não devem ser confundidos com
`PlayerSnapshotV0`.

### Entrada e identidade

```text
PlayerTagInput = string
NormalizedPlayerTag = string       # exatamente um # inicial

AuthContext {
  crownpilotUserId: string         # derivado de Supabase JWT.sub verificado
}

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
sem requisito posterior aprovado.

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

Somente `resolved` pode persistir vínculo. O resultado não contém snapshot,
coleção, Arena, deck ou nome para o domínio de identidade.

### HTTP da aplicação

Os paths concretos devem usar rotas/controllers Laravel e páginas Inertia, mas
devem manter semântica:

| Operação | Autenticação | Resultado esperado |
|---|---|---|
| Ler vínculo primário | obrigatória | vínculo do UID atual ou ausência |
| Criar/substituir vínculo | obrigatória | lookup antes da escrita; falha não remove vínculo antigo |
| Desvincular | obrigatória | remoção idempotente do vínculo do UID atual |
| Excluir dados CrownPilot | obrigatória/reautenticação se exigida | remove usuário e vínculo próprios |

Mapeamento mínimo: `401` sem sessão válida, `400` input inválido, `404` perfil
não encontrado, `409` conflito de versão/concurrency, `429` rate limit e `503`
provider indisponível ou mal configurado. Respostas não vazam token, URL real,
provider route, IAM, payload bruto ou detalhes internos.

## Arquivos, módulos e contratos afetados

Arquivos de bootstrap esperados:

- `composer.json`, `composer.lock`, `package.json`, `package-lock.json`,
  `phpunit.xml`, `pest.php`, `vite.config.ts`, `tsconfig.json` e configs de lint/
  formatter;
- `Dockerfile.vercel`, `Caddyfile`, `vercel.json` e `.dockerignore`;
- `AGENTS.md`, `.env.example`, `.php-version` ou equivalente de pinagem PHP;
- `supabase/config.toml`, `supabase/migrations/`, RLS tests e configuração local;
- `.github/workflows/ci.yml` e configuração de deploy/ambientes sem credenciais.

Módulos prováveis, sujeitos ao layout escolhido no bootstrap:

- `app/Domain/` — `NormalizedPlayerTag`, `PrimaryPlayerLink`, estados;
- `app/Application/` — Auth, link, replace, unlink e delete;
- `app/Contracts/` — `ClashRoyaleClient`, clock e repositories;
- `app/Adapters/ClashRoyale/` — HTTP adapter, redaction, timeout/retry;
- `app/Adapters/Supabase/` — JWT/JWKS verifier, PostgreSQL repository e config;
- `app/Http/Controllers/` e `routes/` — endpoints Laravel finos;
- `resources/js/` — login, páginas Inertia, link, troca, unlink e erros;
- `tests/Unit/`, `tests/Feature/`, `tests/Contract/`, `tests/Rls/`,
  `tests/Browser/` e `tests/Smoke/`.

Decisões duradouras estão registradas na [ADR 003](../decisions/003-supabase-auth-postgresql-jwks-rls-region.md)
para Supabase Auth, PostgreSQL, JWT/JWKS, RLS e região `sa-east-1`. A ADR 002
continua válida para stack Laravel/Inertia, PHP 8.5 e deploy Vercel/Docker. A
`002-02` detalha ambientes e política de persistência; a spec registra o plano,
não finge provisionamento.

## Alternativas descartadas

| Alternativa | Motivo do descarte nesta fase |
|---|---|
| React/Vite + API separada | Dois deploys, CORS e boundary operacional sem benefício para o baseline unitário. |
| Next.js | Não atende decisão explícita de aprofundar Laravel/PHP. |
| Laravel Blade | Reduz stack, mas abandona React/Inertia escolhidos. |
| Runtime Node/Edge | Não atende backend Laravel/PHP e não resolve egress da API externa. |
| Backend separado em Cloud Functions | Cria split runtime e reduz coerência do monólito Laravel. |
| Firebase Auth + Firestore | Sem Admin SDK PHP oficial; exigiria adapters separados para Auth/Firestore e maior risco de compatibilidade. |
| Supabase Data API no browser | Aumenta superfície pública e mistura RLS/Data API com boundary Laravel. |
| Migrations Laravel + Supabase duplicadas | Duas fontes de schema criam drift; `supabase/migrations/` é fonte única. |
| Supabase multi-region/região `Americas` genérica | Pode escolher região diferente de São Paulo; custo/complexidade sem justificativa. |
| Persistir resposta do perfil | Gates de retenção, redistribuição, privacy e agreements permanecem abertos. |
| Proxy, Vercel ou egress como contrato de domínio | Impede troca de transporte e mistura operação com regra de identidade. |

## Critérios de aceite verificáveis

- [ ] clone limpo reproduz bootstrap com `composer install`, `npm ci` e build sem arquivo local
      obrigatório não documentado;
- [ ] `composer.lock` e `package-lock.json`, PHP 8.5, TypeScript strict,
      lint, formatter e type-check possuem scripts executáveis;
- [ ] comandos de unit, integration, contract, RLS, E2E, smoke e build estão
      definidos e documentados;
- [ ] GitHub Actions executa lint, type-check, unit, integration, contract, RLS
      e build em Pull Requests e em todo push para `main`, sem chamada live à
      Clash Royale API;
- [ ] `Dockerfile.vercel` constrói imagem sem secrets e smoke de `/health` passa
      em container local/CI;
- [ ] `Caddyfile` expõe somente `public/` e `vercel.json` roteia service container;
- [ ] `.dockerignore` impede `.env*`, credenciais, `vendor`, `node_modules` e
      testes de entrarem na imagem de produção;
- [ ] Supabase CLI/Docker executa Auth/PostgreSQL local e migrations reproduzíveis;
- [ ] browser usa `@supabase/supabase-js` somente para Auth e não acessa dados
      CrownPilot via Data API/PostgREST;
- [ ] backend valida JWT/JWKS, issuer, audience, expiry, algorithm, kid e deriva
      UUID `sub` server-side;
- [ ] migrations SQL e RLS são versionadas, deny-by-default e testadas contra
      tentativas diretas anônimas e autenticadas;
- [ ] autorização server-side impede usuário A de ler ou alterar rows de B;
- [ ] `sa-east-1` (São Paulo) está disponível no plano e registrada antes de
      provisionar staging/production;
- [ ] local, preview, staging e production têm projetos, secrets, hosts e dados
      separados, com falha fechada para configuração ausente/incompatível;
- [ ] preview gera build/smoke sem Google Sign-In real e sem secrets Supabase de
      produção;
- [ ] staging tem hostname fixo autorizado e Google Sign-In funcional;
- [ ] `main` é a única fonte de deploy de produção;
- [ ] deploy Vercel usa container Laravel/FrankenPHP conforme guia oficial;
- [ ] nenhum token Clash Royale, secret Supabase/PHP ou segredo chega ao
      client, bundle, logs ou Git;
- [ ] `ClashRoyaleClient` é o único boundary de lookup e não aceita host vindo do
      usuário;
- [ ] tag é normalizada com `trim`, um `#` inicial, `%23` somente no path e
      preservação de case;
- [ ] login cria/recupera identidade CrownPilot e uma sessão nova recupera o
      vínculo primário em outro dispositivo;
- [ ] somente resultado de perfil público resolvido cria ou substitui vínculo;
- [ ] vínculo persistido contém `public_profile` e `unverified`, sem alegar
      ownership ou exclusividade;
- [ ] `not_found`, `provider_unavailable`, `rate_limited`, input inválido e
      configuração inválida são estados distintos e testados;
- [ ] troca valida a nova tag antes de substituir a antiga e falha externa não
      remove vínculo existente;
- [ ] desvinculação e exclusão dos dados próprios são explícitas e idempotentes;
- [ ] nenhuma resposta externa é persistida como snapshot, coleção, Arena,
      battle history, cache ou payload raw;
- [ ] testes E2E em staging cobrem login, vínculo, reload/outro dispositivo,
      troca e desvinculação;
- [ ] smoke de staging e production passa sem expor dados de teste de outro
      ambiente;
- [ ] observabilidade registra resultado, latência e erro sem Player Tag real,
      e-mail, token, URL real ou payload externo.

## Estratégia de testes e validação

### Unit

Testar em Pest/PHPUnit e Vitest normalização e validação de Player Tag, `%23`,
estados de ausência, contratos de erro, mapeamento HTTP, redaction, clock e
invariantes de link.

### Integration

Usar Supabase CLI/Docker com Auth e PostgreSQL para verificar
criação/recuperação do usuário, validação JWT/JWKS, autorização por UUID,
replace transacional, unlink, delete idempotente e falha sem perda do vínculo
antigo.

### RLS

Testar acesso direto anônimo e autenticado às tabelas. RLS deve ser habilitado e
deny-by-default; testar separadamente que a conexão Eloquent privilegiada não é
confundida com RLS e que autorização Laravel cobre cada operação.

### Contract

Fixtures sanitizadas cobrem resposta resolvida, `404`, `403`, `429` com e sem
`Retry-After`, `5xx`, timeout, host/configuração inválida e payload incompleto.
Nenhum teste normal chama API live.

### E2E e smoke

Playwright ou ferramenta equivalente roda em staging com Google via Supabase Auth
e usuário de teste: sign-in, refresh, vínculo, recuperação, troca, unlink,
logout/login e estados de erro. Smoke pós-deploy verifica build, JWT/JWKS,
leitura controlada do vínculo e isolamento de ambiente. Production usa conta/test
data controlada e não expõe fixture nem segredo.

Comandos normativos a configurar no bootstrap:

```text
composer run lint
composer run analyse
composer run test:unit
composer run test:integration
php artisan test
npm run lint
npm run typecheck
npm run test:contract
npm run test:rls
npx supabase test db
npm run test:e2e
npm run smoke
npm run build
docker build -f Dockerfile.vercel -t crownpilot-vercel-ci .
```

## Observabilidade

Adicionar logs estruturados e métricas mínimas para:

- sign-in aceito/recusado por categoria;
- link, replace, unlink e delete por resultado;
- latência e classe de status do lookup externo;
- `not_found`, rate limit, timeout, indisponibilidade e misconfiguration;
- falhas de autorização e configuração por ambiente;
- correlation/request ID e versão do app.

Redaction obrigatória: não registrar Player Tag real, URL contendo tag, e-mail,
nome de usuário, Supabase JWT, service key, API token, IP, payload externo ou segredo. Logs
devem permitir diagnosticar falha sem reidentificar perfil público.

## Riscos e mitigação

| Risco | Impacto | Mitigação/saída |
|---|---|---|
| `sa-east-1` indisponível | Provisionamento bloqueado ou região inadequada | Confirmar plano/organização antes de criar projetos; sem fallback silencioso. |
| JWT/JWKS inválido ou key rotation | Login quebrado ou aceitação indevida | Validar issuer/audience/expiry/alg/kid, cache curto e refresh em `kid` desconhecido. |
| RLS ignorado pela conexão Eloquent | Vazamento entre usuários | Autorização Laravel obrigatória; role/configuração real e testes A/B. |
| Proxy/egress sem contrato | Lookup indisponível ou incompatível | Port `ClashRoyaleClient`, fixtures, timeout/retry limitado, configuração substituível e fallback de erro sem perda local. |
| Ownership não verificado | Comunicação enganosa ou ação indevida | Constante `unverified`, texto de UI, sem exclusividade/notificação/ação. |
| Tag ou token em logs/bundle | Exposição de dado/segredo | Redaction testada, secrets server-only, scan de bundle e revisão de logs. |
| Troca concorrente de tag | Vínculo incorreto ou perda de estado | Validar antes de escrever, precondition/transaction, `409` e testes concorrentes. |
| Preview acessa Supabase real | Corrupção/vazamento entre ambientes | Preview sem Supabase real/Auth real, project ref allowlist e startup fail-closed. |
| Persistência prematura de API data | Violação de gate de retenção/privacy | Documento mínimo de vínculo; proibição explícita de snapshot/raw/cache nesta fase. |
| Supabase/Laravel acoplados ao domínio | Migração cara | ports/adapters, controllers Inertia finos, PHP host portável e nenhuma API proprietária Vercel. |
| Container Vercel/FrankenPHP incompatível | Deploy bloqueado | Build Docker e smoke `/health` no GH Actions antes de staging/production; host PHP first-class continua fallback. |
| `pdo_pgsql`, SSL ou pooler incompatíveis | PostgreSQL/Auth sem execução | Gate técnico; validar container, conexão SSL e session pooler/direct. |
| Filesystem/container efêmero | Perda de sessão/cache/arquivos | PostgreSQL para estado durável, sessão cookie ou store externo aprovado; nenhum arquivo local como fonte de verdade. |
| `APP_KEY` ausente ou compartilhado | Cookies inválidos ou vazamento entre ambientes | Secret estável por ambiente, injetado pela Vercel, nunca na imagem/Git. |

## Rollout, rollback e migração

1. Bootstrap local e gates sem provisionar produção.
2. Registrar ADRs de stack, Supabase Auth/PostgreSQL, `sa-east-1`, JWT/JWKS, RLS
   e persistência.
3. Configurar Supabase CLI/Docker, migrations e fixtures; validar RLS e fluxos de
   identidade.
4. Configurar GitHub Actions para Pull Requests e pushes em `main`.
5. Construir e testar `Dockerfile.vercel`/FrankenPHP, incluindo `/health`.
6. Confirmar `sa-east-1` e configurar staging com projeto/host fixos e Google
   Sign-In via Supabase Auth.
7. Rodar E2E e smoke em staging; corrigir antes de qualquer produção.
8. Provisionar staging/production em `sa-east-1` após validar JWT/JWKS, PostgreSQL,
   RLS, container e secrets próprios.
9. Deploy production somente pela `main`; executar smoke controlado.

Não há migração de dados legados no baseline. Documentos devem carregar
`schemaVersion`; mudanças incompatíveis exigem migração explícita. Rollback de
aplicação deve manter documentos compatíveis e nunca apagar vínculo. Falha do
provider externo degrada link/refresh, mas não remove identidade nem vínculo
existente. Mudança de região exige novo projeto/database e cutover separado; não
é operação in-place.

## Ordem das subtarefas

1. [002-01 — bootstrap do toolchain](../tasks/002-fundacao-aplicacao-identidade-persistente/002-01-bootstrap-toolchain.md)
2. [002-02 — boundaries e ambientes](../tasks/002-fundacao-aplicacao-identidade-persistente/002-02-fixar-boundaries-e-ambientes.md)
3. [002-03 — Supabase local e projetos](../tasks/002-fundacao-aplicacao-identidade-persistente/002-03-configurar-supabase-local-e-projetos.md)
4. [002-04 — Auth Google e sessão](../tasks/002-fundacao-aplicacao-identidade-persistente/002-04-implementar-auth-google-e-sessao.md)
5. [002-05 — adapter de lookup](../tasks/002-fundacao-aplicacao-identidade-persistente/002-05-criar-adapter-de-lookup.md)
6. [002-06 — persistência e autorização](../tasks/002-fundacao-aplicacao-identidade-persistente/002-06-persistir-vinculo-com-autorizacao.md)
7. [002-07 — fluxos de vínculo e exclusão](../tasks/002-fundacao-aplicacao-identidade-persistente/002-07-entregar-fluxos-de-vinculo-e-exclusao.md)
8. [002-08 — quality gates e observabilidade](../tasks/002-fundacao-aplicacao-identidade-persistente/002-08-automatizar-quality-gates-e-observabilidade.md)
9. [002-09 — staging, deploy, E2E e handoff](../tasks/002-fundacao-aplicacao-identidade-persistente/002-09-validar-staging-deploy-e-handoff.md)

As subtarefas 01–03 estabelecem execução, decisões e Supabase. 04–07 entregam
identidade e vínculo. 08 fecha gates automatizados. 09 valida ambientes reais,
deploy e handoff. Nenhuma subtarefa libera Fase 003 sem reabrir os gates
explicitados no veredito da Fase 001.

## Premissas explícitas

- O ticker é `002`: não havia ticker explícito no pedido e `001` era o maior
  prefixo numérico existente em `docs/specs/`/`docs/tasks/`.
- Google é o único provider de login nesta fase.
- Uma Player Tag primária basta para o MVP; múltiplas tags ficam fora.
- O usuário fornece a tag; não haverá credential exchange com Supercell.
- A validação externa é best-effort e pode ficar indisponível; identidade local
  continua válida.
- `sa-east-1` (São Paulo) é região escolhida para staging e production, condicionada
  à disponibilidade no plano; provisionar somente após validar isolamento.
- Laravel + Inertia + React + TypeScript + Vite substituem Next.js; Composer e
  npm possuem lockfiles próprios.
- Vercel continua alvo inicial, com compatibilidade Laravel/PHP validada por
  container; host PHP first-class é fallback se a prova falhar.
- Deploy Vercel segue container Docker/FrankenPHP, com `Dockerfile.vercel`,
  `Caddyfile`, `vercel.json` e `.dockerignore`.
- GitHub Actions é CI obrigatório para Pull Requests e pushes em `main`.
- A operação do proxy observada na Fase 001 é configuração substituível e não
  autorização de retenção, redistribuição ou cobrança.
- O contrato v0 orienta provenance e ausência, mas não vira schema de snapshot.
- Não há dados legados a migrar e nenhum deploy existente a preservar.

## Handoff esperado

Ao concluir a fase, entregar:

- aplicação que sobe a partir de clone limpo;
- stack, comandos, ambientes e decisions documentados;
- Supabase Auth Google funcionando em staging;
- PostgreSQL/Supabase seguro, separado, com migrations e RLS testados;
- vínculo primário privado, removível e `unverified` recuperável em qualquer
  dispositivo;
- adapter de lookup server-side coberto por fixtures;
- CI, E2E e smoke com evidência;
- observabilidade sem dados sensíveis;
- riscos residuais e gates de API data, retenção, ownership, egress, meta e
  compliance explicitamente mantidos para revisão.

O próximo consumidor é a Fase 003, mas somente após reabertura e aprovação dos
gates que a Fase 001 deixou fora do escopo. A Fase 002 não pode apresentar a
existência do vínculo como ownership nem transformar autenticação em autorização.
