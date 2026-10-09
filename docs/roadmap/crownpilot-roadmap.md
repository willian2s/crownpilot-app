# Roadmap — CrownPilot

> **Play the right deck. Upgrade the right cards.**

**Atualizado em:** 8 de outubro de 2026

Este documento define **o que construir**, **em qual ordem**, as principais dependências, os boundaries do produto e os grandes marcos do CrownPilot.

O roadmap **não substitui specs detalhadas**. Specs definem como uma fase será implementada. O roadmap define direção, sequência, dependências e critérios para liberar a próxima etapa.

---

# 1. Visão do produto

CrownPilot é um companion personalizado para Clash Royale.

A pergunta central do produto é:

> **Qual é o melhor deck que eu consigo jogar agora, na minha Arena, com a conta que eu realmente tenho?**

A experiência deve partir da conta do jogador, e não de um ranking genérico.

Quando os dados permitirem, a recomendação combina:

- Arena e faixa de troféus;
- coleção de cartas;
- níveis;
- Evolutions / Heroes disponíveis;
- decks observados em partidas reais;
- força do deck no meta relevante;
- tamanho e confiança da amostra;
- histórico e desempenho individual em fases posteriores.

O produto evolui em quatro camadas:

- **DISCOVER** — O que eu deveria jogar?
- **PROGRESS** — O que eu deveria melhorar?
- **OPTIMIZE** — Como melhorar o que eu já jogo?
- **COACH** — Por que isso faz sentido para mim?

## 1.1 Proposta de valor

O CrownPilot não deve responder apenas:

> Quais são os melhores decks do meta?

Ele deve responder:

> Quais são os melhores decks **para mim**?

Dois jogadores na mesma Arena podem receber recomendações diferentes porque possuem coleções, níveis, Evolutions / Heroes e progressão diferentes.

## 1.2 Hipótese principal

O produto gera valor recorrente quando consegue transformar:

**conta do jogador + contexto competitivo + decks reais + progressão**

em:

**recomendações acionáveis e explicáveis.**

O objetivo não é ser apenas um deck finder.

O objetivo é se tornar a camada de inteligência da conta do jogador.

---

# 2. Princípios arquiteturais e de produto

## 2.1 Deterministic by default

> **Deterministic by default. AI when it adds value.**

Toda feature que puder ser resolvida de forma determinística deve funcionar sem LLM.

Isso inclui inicialmente:

- Best Decks for You;
- Deck Readiness;
- Arena / Trophy Fit;
- Collection Fit;
- Upgrade Impact;
- Upgrade Planner;
- comparação de variantes;
- métricas de matchup.

IA deve interpretar e explicar resultados do engine, não substituir o engine.

## 2.2 A conta é o centro do produto

O usuário vincula sua Player Tag uma única vez.

Depois de autenticado em qualquer dispositivo, o CrownPilot deve recuperar a conta vinculada sem exigir que o usuário abra o jogo para copiar a tag novamente.

Nunca armazenar credenciais da conta Supercell.

## 2.3 Recomendações precisam de evidência

A primeira versão não deve gerar oito cartas do zero.

O recommendation engine começa com **decks reais observados** e ranqueia esses candidatos para o jogador.

Deck generation só entra em exploração depois que Discover e Optimize estiverem validados.

## 2.4 Arena é contexto, não apenas filtro

Arena / faixa de troféus influencia a relevância.

O produto não deve assumir que o deck mais forte no topo da ladder é automaticamente a melhor recomendação para um jogador em outra faixa competitiva.

## 2.5 Toda recomendação deve ser explicável

Para qualquer deck recomendado, o sistema precisa responder:

> Por que este deck está sendo recomendado para este jogador?

O score deve ser decomponível em fatores observáveis.

## 2.6 Dados vencidos não são meta

Balance changes e mudanças de temporada podem invalidar dados antigos rapidamente.

Toda estatística competitiva deve carregar:

- período;
- população;
- tamanho da amostra;
- contexto;
- confiança.

## 2.7 Infraestrutura-base e portabilidade

Decisões iniciais do projeto, atualizadas pela ADR 005 proposta:

- Firebase Authentication com Google para a conta CrownPilot;
- Go `1.27.2` com `net/http`/`ServeMux` como API da aplicação;
- React + TypeScript + Vite como frontend independente;
- PostgreSQL como banco principal, hospedado inicialmente no Supabase;
- `pgx`/`sqlc` como acesso oficial ao PostgreSQL e `goose` para migrations;
- API Render em Virgínia (`us-east`) e região Supabase `us-east-1` (Northern
  Virginia); dados pessoais ficam fora do Brasil e transferência internacional,
  DPA, backups, subprocessadores, residência e egress são gates antes de dados
  reais;
- Docker/OCI como runtime portátil da API;
- Render como hosting inicial da API Docker e alvo preferido do frontend estático;
- Vercel como alternativa de frontend estático, nunca backend obrigatório.

Vercel não deve se tornar um boundary do domínio.

O core da aplicação e integrações devem permanecer portáveis para outro runtime
sem reescrita do domínio. Serviços exclusivos da Vercel podem ser usados apenas
quando isolados atrás de adapters ou quando existir estratégia clara de
substituição. A API não depende da Vercel para executar.

Render é escolha inicial por simplicidade, custo e suporte a containers, não
dependência arquitetural. A imagem OCI deve poder executar futuramente em Azure
App Service, Azure Container Apps, AWS, GCP ou outro runtime. Azure não é hosting
obrigatório nesta fase; Kubernetes só entra mediante necessidade real.

A integração com a Clash Royale API deve permanecer separada do runtime,
especialmente por possíveis requisitos de egress/IP allowlist.

## 2.8 Segurança e ambientes fazem parte da fundação

A fundação deve nascer segura por padrão:

- segredos de integração permanecem somente no servidor;
- browser não acessa PostgreSQL ou Data API para dados CrownPilot;
- API REST versionada em `/api/v1`, com OpenAPI gerado e ProblemDetails;
- health checks liveness/readiness, structured logging, correlation IDs, métricas
  básicas e redaction desde a fundação;
- testes unit, Application, integration, persistence/RLS, authn/authz, contract,
  API/OpenAPI, frontend, E2E e smoke como gates progressivos;
- migrations SQL e SQL de RLS/grants são versionados/testados, sem duplicar schema;
- acesso server-side usa `pgx` com secrets mínimos;
- local, preview e produção não devem compartilhar dados/segredos de forma
  acidental;
- Docker Compose/PostgreSQL pinado deve ser usado para isolamento e testes locais
  reproduzíveis; projeto Supabase dev separado atende desenvolvimento persistente;
- decisões difíceis de reverter, como região Supabase, precisam ser tomadas e
  registradas antes de provisionar produção.

A estratégia de client vs. server access fica definida na Fase 002: Firebase
cuida de authentication, enquanto API Go aplica authorization. Não é permitido
tratar authentication como authorization.

## 2.9 Estratégia de ambientes e release

O CrownPilot usa ambientes com responsabilidades diferentes.

### Local

Desenvolvimento local usa:

- Docker Compose/PostgreSQL pinado para testes/reset locais; projeto Supabase dev
  separado para desenvolvimento persistente, via Session pooler;
- Firebase Emulator ou fixtures para authentication;
- API Go e frontend Vite executados separadamente;
- dados locais/descartáveis;
- nenhum secret de produção.

### Pull Request / Preview

Branches de feature e Pull Requests podem gerar Preview Deployments efêmeros em
Render Static Site, Vercel ou outro host estático.

Preview serve para:

- validar build;
- revisar UI/UX;
- executar smoke checks que não dependem de autenticação real;
- compartilhar uma versão da mudança antes do merge.

Preview **não é ambiente Supabase completo**.

Cada PR possui hostname potencialmente efêmero. Como Firebase Authentication
exige configuração de domínios autorizados, não vamos automatizar autorização de
hosts efêmeros nem compartilhar Firebase/Supabase de produção com previews.

Portanto, Preview não precisa oferecer Google Sign-In real.

Features que exigem Firebase real, API completa ou PostgreSQL são validadas em
staging.

### Staging

A branch `staging` é o ambiente estável de pré-produção.

Deve possuir:

- hostname fixo;
- host estático próprio, Vercel opcional;
- environment variables próprias;
- Firebase project separado de produção;
- Supabase project/database separado de produção;
- redirect/authorized domains configurados no Firebase;
- Google Sign-In real;
- PostgreSQL Supabase de staging em `us-east-1` (Northern Virginia), com API
  Render em Virgínia (`us-east`) e Session pooler após aprovação da ADR 005;
- dados de teste;
- execução de testes E2E e smoke.

O gate de Staging roda por workflow protegido/manual ou promoção equivalente,
com owner definido, aprovação do ambiente, migration job controlado e smoke
obrigatório antes de qualquer promoção para `main`/Production. PR e `main` não
substituem essa validação de integração.

Fluxo de release esperado:

```text
feature/*
    |
    v
Pull Request
    |
    +--> Preview efêmero
    |      build / UI / smoke sem auth real
    |
    v
staging
    |
    +--> Firebase staging
    +--> Supabase PostgreSQL staging
    +--> E2E
    +--> smoke
    |
    v
main
    |
    v
production
```

### Production

A branch `main` é a única fonte de deploy de produção.

Production usa Firebase project, Supabase project, secrets e domínio próprios.

## 2.10 Estratégia de testes

Testes fazem parte da fundação e evoluem junto com o produto.

Camadas previstas:

- **unit** — funções puras, domínio, normalização, scoring e canonicalização;
- **integration** — adapters, Firebase/API/PostgreSQL e boundaries;
- **persistence/RLS** — migrations `goose`, policies PostgreSQL e grants testados
  contra PostgreSQL local via Docker Compose;
- **contract** — contratos de integrações como `ClashRoyaleClient` usando
  fixtures sanitizadas;
- **E2E** — fluxos críticos completos executados em staging;
- **smoke** — checks rápidos após deploy em staging e production.

Chamadas live para a Clash Royale API não devem fazer parte da CI normal.
Probes live devem ser explícitos, limitados e executados somente quando houver
motivo técnico.

Gates mínimos esperados:

- **PR:** Go restore/build/test + lint + type-check + unit + integration +
  persistence + contract + RLS + build Docker descartável;
- **staging candidate:** gates de PR + publica um digest OCI imutável e executa
  migration job, E2E e smoke;
- **main:** promove o mesmo digest aprovado, sem rebuild divergente, e executa
  health smoke;
- **production:** somente após aprovação dos gates, deploy do digest promovido e
  smoke não destrutivo.

A ferramenta específica de testes é escolhida no bootstrap da Fase 002.

## 2.11 Compliance é boundary arquitetural

CrownPilot é um companion de análise e coaching.

Não fazem parte do produto:

- automação de gameplay;
- bots;
- cliente alternativo;
- interferência no jogo;
- overlay competitivo em tempo real;
- ferramentas destinadas a obter vantagem injusta durante a partida.

Monetização também não deve ser assumida como livre de restrições.

A Fan Content Policy vigente da Supercell deve ser tratada como requisito do produto:

https://supercell.com/en/fan-content-policy/

---

# 3. Convenções de status

| Status | Significado |
|---|---|
| ✅ | Concluído e validado |
| 🚧 | Em andamento |
| ⬜ | Planejado |
| ⏸️ | Bloqueado por decisão ou dependência |
| 🧪 | Discovery / investigação |

Uma fase só é concluída quando:

- entregas previstas estão implementadas;
- critérios de aceite foram validados;
- decisões relevantes estão documentadas;
- riscos e débitos remanescentes estão explícitos;
- existe handoff claro para a próxima fase.

---

# 4. Fases

## 001 — Viabilidade de produto, dados e compliance ✅

### Objetivo

Validar se os dados e permissões disponíveis sustentam o produto antes de comprometer a arquitetura.

### Status e veredito

✅ **Fase encerrada com constraints.** O veredito é `GO WITH CONSTRAINTS /
APPROVAL DEPENDENCY`. O núcleo técnico de bootstrap, identidade CrownPilot e
vínculo privado read-only de perfil público pode avançar; rota oficial direta,
ownership, coleção completa, histórico completo, meta representativa da Arena,
retenção/redistribuição de API data e billing não estão liberados.

Evidência consolidada em [phase-001-verdict.md](../tasks/001-viabilidade-produto-dados-compliance/evidences/phase-001-verdict.md).

### Escopo

Investigar e documentar:

- acesso à API oficial;
- consulta de jogador por Player Tag;
- campos disponíveis no perfil;
- Arena e troféus;
- coleção e níveis;
- representação disponível de Evolutions / Heroes;
- battle log;
- cards e locations;
- rankings ou outras superfícies úteis para amostragem;
- rate limits;
- limites de histórico;
- estratégia de cache;
- restrições de uso dos dados;
- políticas aplicáveis;
- semântica de vínculo da Player Tag e possibilidade de prova de ownership.

### Perguntas críticas

1. Conseguimos reconstruir de forma confiável o estado da coleção necessário para Deck Readiness?
2. Conseguimos determinar Arena / faixa competitiva sem heurísticas frágeis?
3. Como obter uma amostra representativa de decks por faixa?
4. O battle log é suficiente para histórico contínuo por snapshots?
5. Precisaremos de fonte externa complementar para meta?
6. Quais dados econômicos não estão disponíveis e exigiriam entrada manual?
7. Quais limites de uso impactam sync e ingestão em escala?
8. Vincular uma Player Tag significa selecionar um perfil público ou precisamos
   provar ownership? Existe mecanismo oficial aplicável e ele é necessário para
   o MVP?

### Compliance comercial

Antes de qualquer assinatura, paywall ou cobrança por feature, validar explicitamente o modelo comercial contra a política vigente e demais acordos aplicáveis.

A Fan Content Policy vigente descreve Fan Content como não comercial por padrão: não é permitido cobrar taxa de qualquer tipo, incluindo por funcionalidades in-app, sem aprovação expressa da Supercell. A política lista anúncios, donations e coaching como exceções específicas, mas cada exceção possui condições próprias e não equivale a uma autorização geral para SaaS, premium features ou AI/software coaching.

**Não interpretar isso automaticamente como autorização para qualquer plano Pro.**

O gate da Fase 001 deve determinar, com fonte oficial e classificação explícita, quais modelos de monetização são permitidos para o CrownPilot e quais exigem aprovação expressa da Supercell. Assinatura/paywall por funcionalidades não pode ser assumida como permitida. Ads, donations e coaching devem ser analisados individualmente; a exceção de coaching não autoriza automaticamente coaching por software ou AI. Qualquer modelo dependente de aprovação expressa permanece bloqueado até a aprovação correspondente.

A arquitetura deve continuar funcional sem depender de monetização ainda não aprovada. O produto pode avançar tecnicamente, operar gratuitamente ou usar apenas modelo classificado como permitido sob condições, sem converter a ausência de assinatura em NO-GO automático.

### Entregas

- mapa das fontes de dados;
- matriz campo → fonte → freshness → limitação;
- spike de integração;
- modelo inicial de Player Snapshot;
- modelo inicial de Card Collection;
- proposta de aquisição do dataset de meta;
- registro dos principais riscos;
- decisão documentada sobre viabilidade do MVP.

### Critérios de aceite

Os critérios abaixo preservam comportamento funcional. Menções a EF Core,
ASP.NET Core e code-first são baseline histórico; o destino após ADR 005 usa
`goose`, Go, `depguard` e OpenAPI spec-first conforme `002-13` a `002-15`.

- [ ] Player Tag válida pode ser resolvida de forma reproduzível;
- [ ] campos do MVP estão classificados como disponíveis, derivados, opcionais ou indisponíveis;
- [ ] estratégia de meta possui uma fonte inicial viável para candidatos
  personalizados, sem alegar meta representativa da Arena;
- [ ] limites de rate e histórico estão documentados;
- [ ] boundary de compliance está documentado;
- [ ] modelos de monetização estão classificados sob a Fan Content Policy como `ALLOWED`, `ALLOWED WITH CONDITIONS`, `REQUIRES EXPLICIT APPROVAL`, `NOT ALLOWED` ou `UNRESOLVED`;
- [ ] nenhuma dependência crítica do MVP exige assinatura, paywall ou feature paga ainda não aprovada;
- [ ] semântica de vínculo/ownership da Player Tag está decidida;
- [ ] nenhuma dependência crítica do MVP continua baseada apenas em suposição.

### Handoff

A Fase 002 está **released with constraints** para bootstrap reproduzível,
Firebase Authentication, identidade CrownPilot e vínculo privado read-only de
perfil público. Pode validar a tag server-side, mas não sincroniza nem persiste
coleção, níveis, Arena, battle history ou Player Snapshot completo; isso começa
somente após os gates posteriores. Não há liberação de billing, ownership ou meta
Arena por este handoff.

---

## 002 — Fundação da aplicação e identidade persistente ⬜

**Handoff status:** `released with constraints` — `002-01` a `002-04` concluídas
como baseline .NET; migração Go `002-13` a `002-15` permanece pendente de ADR 005.

### Objetivo

Criar a fundação executável do CrownPilot, estabelecer seus boundaries de
segurança e entregar a identidade persistente que elimina o atrito de informar a
Player Tag em cada dispositivo.

### Escopo

#### Bootstrap da aplicação

- inicializar API Go e frontend React/Vite após aprovação da ADR 005;
- definir package manager e lockfile autoritativo;
- TypeScript strict e configuração de build;
- lint/format e convenções de código;
- estrutura inicial de diretórios e boundaries;
- `AGENTS.md` e comandos operacionais do repositório;
- toolchain de testes unitários, integração, contract e E2E;
- CI com gates de lint, type-check, testes e build;
- OpenAPI/ProblemDetails como contrato de API e health checks como contrato
  operacional;
- `.env.example` sem secrets;
- estratégia local / preview / staging / production;
- Preview Deployments efêmeros para PRs sem dependência de Firebase real;
- branch `staging` com hostname fixo para validação completa;
- deploy de produção somente a partir da `main`;
- frontend estático hospedável no Render Static Site, Vercel ou alternativa
  compatível, sem dependência obrigatória do backend.

#### Identidade e persistência

- Firebase Authentication com Google;
- PostgreSQL como banco principal, hospedado inicialmente no Supabase;
- `pgx`/`sqlc` e `goose` para acesso e migrations do schema após o cutover;
- SQL separado somente para RLS, grants e objetos de plataforma;
- usar Supabase `us-east-1` com API Render em Virgínia (`us-east`) e validar DPA,
  backups, subprocessadores, transferência internacional, residência e egress
  antes de produção;
- Docker Compose/PostgreSQL pinado para testes e projeto Supabase dev via Session
  pooler para desenvolvimento;
- Firebase projects e Supabase databases separados para staging e production;
- authorized domains fixos de staging no Firebase;
- validar Firebase ID Token no adapter Firebase Admin Go e usar bearer;
- RLS/grants versionados, sem duplicar schema das migrations SQL Go, e testados;
- browser sem acesso ao PostgreSQL/Data API para dados CrownPilot;
- garantir que local/preview/staging não usem produção por acidente.

#### Identidade CrownPilot

- usuário CrownPilot;
- uma Player Tag primária no MVP;
- vínculo persistente com Player Tag;
- validação de existência da tag;
- aplicar a decisão da Fase 001 sobre perfil público vs. ownership verificado;
- tag inválida / jogador inexistente;
- troca e desvinculação controladas;
- recuperação do vínculo em outro dispositivo;
- estratégia inicial para exclusão dos dados da conta;
- observabilidade inicial.

#### Arquitetura e aprendizado

- Clean Architecture pragmática dentro de um Modular Monolith;
- um processo/backend e uma imagem OCI, com módulos internos e boundaries claros;
- Domain sem dependência de ASP.NET Core, Firebase, Supabase, EF Core, Npgsql ou
  providers externos;
- Application dependente de ports/interfaces, Infrastructure implementando
  adapters e API compondo transporte;
- código didático sem artificialidade, com comentários de intenção na primeira
  ocorrência de DI, middleware, authn/authz, EF Core, migrations, Npgsql,
  `async/await`, `CancellationToken`, options e lifecycle;
- documentação incremental em `docs/learning/`, somente para conceitos usados.

### Plano executável da Fase 002

As subtarefas são implementáveis e verificáveis nesta ordem:

1. [002-01-bootstrap-toolchain.md](../tasks/002-fundacao-aplicacao-identidade-persistente/002-01-bootstrap-toolchain.md);
2. [002-02-estabelecer-boundaries-contrato-base-e-ambientes.md](../tasks/002-fundacao-aplicacao-identidade-persistente/002-02-estabelecer-boundaries-contrato-base-e-ambientes.md);
3. [002-03-preparar-postgresql-migrations-e-harness-rls.md](../tasks/002-fundacao-aplicacao-identidade-persistente/002-03-preparar-postgresql-migrations-e-harness-rls.md);
4. [002-04-implementar-google-sign-in-e-firebase-bearer.md](../tasks/002-fundacao-aplicacao-identidade-persistente/002-04-implementar-google-sign-in-e-firebase-bearer.md);
5. [002-13-bootstrap-http-config-openapi-go.md](../tasks/002-fundacao-aplicacao-identidade-persistente/002-13-bootstrap-http-config-openapi-go.md);
6. [002-14-autenticacao-firebase-go.md](../tasks/002-fundacao-aplicacao-identidade-persistente/002-14-autenticacao-firebase-go.md);
7. [002-15-persistencia-cutover-remocao-dotnet.md](../tasks/002-fundacao-aplicacao-identidade-persistente/002-15-persistencia-cutover-remocao-dotnet.md);
8. [002-05-implementar-port-e-adapter-de-lookup.md](../tasks/002-fundacao-aplicacao-identidade-persistente/002-05-implementar-port-e-adapter-de-lookup.md);
9. [002-06-modelar-persistencia-repositories-e-rls.md](../tasks/002-fundacao-aplicacao-identidade-persistente/002-06-modelar-persistencia-repositories-e-rls.md);
10. [002-07-implementar-casos-de-uso-e-api-v1.md](../tasks/002-fundacao-aplicacao-identidade-persistente/002-07-implementar-casos-de-uso-e-api-v1.md);
11. [002-08-entregar-frontend-de-identidade-e-vinculo.md](../tasks/002-fundacao-aplicacao-identidade-persistente/002-08-entregar-frontend-de-identidade-e-vinculo.md);
12. [002-09-revisar-arquitetura-frontend-e-ux-visual.md](../tasks/002-fundacao-aplicacao-identidade-persistente/002-09-revisar-arquitetura-frontend-e-ux-visual.md);
13. [002-10-instrumentar-observabilidade-health-e-redaction.md](../tasks/002-fundacao-aplicacao-identidade-persistente/002-10-instrumentar-observabilidade-health-e-redaction.md);
14. [002-11-automatizar-ci-oci-e-gates-de-release.md](../tasks/002-fundacao-aplicacao-identidade-persistente/002-11-automatizar-ci-oci-e-gates-de-release.md);
15. [002-12-validar-staging-e2e-smoke-e-handoff.md](../tasks/002-fundacao-aplicacao-identidade-persistente/002-12-validar-staging-e2e-smoke-e-handoff.md).

`002-01` a `002-04` permanecem histórico concluído do baseline. A migração
`002-13` -> `002-14` -> `002-15` deve completar seus gates, com ADR 005 aprovada,
antes de qualquer task `002-05` a `002-12` iniciar, inclusive contra o baseline.
`002-09-revisar-arquitetura-frontend-e-ux-visual.md` não é uma task genérica de
melhoria: analisa o código React entregue,
seus boundaries, auth/API client, estado, configuração, testes, acessibilidade,
responsividade e estados visuais em mobile/desktop, registrando achados e
correções. `002-12-validar-staging-e2e-smoke-e-handoff.md` promove exatamente o
digest OCI validado em Staging; não há rebuild divergente.

### Boundary da fase

A Fase 002 pode consultar a Clash Royale API para validar o vínculo, mas **não**
sincroniza nem persiste coleção, níveis, Arena, battle history ou Player
Snapshot completo. Isso começa na Fase 003.

O token da Clash Royale API nunca é exposto ao browser.

### Requisitos

- nenhuma credencial da Supercell é solicitada;
- Firebase Authentication identifica o usuário externo;
- CrownPilot User possui ID interno separado do Firebase UID;
- Player Tag é vínculo de domínio, com semântica definida pela Fase 001;
- novo dispositivo recupera a tag vinculada após authentication;
- mudança/desvinculação da tag é explícita;
- autenticação não é tratada como autorização;
- falhas da API externa não invalidam a identidade local;
- runtime/domain não dependem de API proprietária da Vercel;
- mesma imagem Docker/OCI deve ser promovível por digest entre Staging e
  Production, sem acoplamento ao Render;
- authentication não é authorization; regras críticas ficam no backend.

### Critérios de aceite

- [ ] bootstrap pode ser reproduzido a partir do repositório limpo;
- [ ] lint, type-check, unit, integration, contract, E2E e build possuem comandos definidos;
- [ ] CI executa lint, type-check, unit, integration, contract, RLS e build nos PRs;
- [ ] staging candidate publica digest OCI e executa smoke; `main` promove o
      mesmo digest sem rebuild divergente e executa health smoke;
- [ ] migrations `goose` e testes RLS executam contra PostgreSQL local via Docker
      Compose;
- [ ] PRs geram Preview Deployments sem depender de Firebase/Google Sign-In real;
- [ ] branch `staging` possui hostname fixo e ambiente de pré-produção;
- [ ] staging possui Firebase project e Supabase database separados, com Google Sign-In funcional;
- [ ] após aprovação da ADR 005, Staging/Production usam Session pooler e conexão
      direta permanece restrita ao PostgreSQL Docker local descartável;
- [ ] E2E do fluxo crítico roda em staging;
- [ ] `main` é a única fonte de deploy de produção;
- [ ] ambientes local/preview/staging/production estão separados e documentados;
- [ ] Supabase `us-east-1` e Render Virgínia estão provisionados; DPA, backups,
      subprocessadores, transferência internacional, residência e egress estão
      validados antes do banco de produção;
- [ ] PostgreSQL/Data API não possui acesso público irrestrito;
- [ ] RLS/grants e autorização backend Go refletem o boundary escolhido e possuem
      validação;
- [ ] usuário consegue criar sessão com Google;
- [ ] usuário vincula uma Player Tag uma vez;
- [ ] outro dispositivo recupera o vínculo após login;
- [ ] usuário consegue trocar/desvincular a tag;
- [ ] erros de integração não causam perda do vínculo;
- [ ] nenhum token da Supercell, Firebase service account ou secret de banco chega ao client;
- [ ] frontend estático e API Docker funcionam sem tornar Render ou Vercel parte
      do domínio;
- [ ] OpenAPI spec-first gerado documenta endpoints, authn/authz, ProblemDetails e status;
- [ ] health liveness/readiness e observabilidade redacted possuem checks;
- [ ] smoke de Staging passa; smoke de Production é executado quando o ambiente
      for provisionado e houver go/no-go explícito.

### Dependências

- Fase 001.

### Handoff

Entregar aplicação reproduzível, identidade estável, PostgreSQL/RLS seguro e
boundaries suficientes para a Fase 003 implementar o sync da conta.

---

## 003 — Sync da conta e Player Dashboard ⬜

### Objetivo

Transformar a Player Tag em um estado de conta útil dentro do CrownPilot.

### Escopo

Sincronizar e normalizar, conforme disponibilidade confirmada na Fase 001,
sempre através de integração server-side com a Clash Royale API:

- definir política de freshness e gatilhos de refresh;
- separar payload bruto, normalização e snapshot de domínio;
- tratar concorrência/retry para impedir snapshots inconsistentes;
- definir retenção mínima necessária antes de acumular histórico;

Sincronizar:

- perfil;
- Arena;
- troféus;
- coleção;
- níveis;
- Evolutions / Heroes;
- deck atual quando disponível;
- timestamps de atualização.

### Dashboard inicial

O dashboard deve responder:

- quem é este jogador?
- onde ele está?
- qual o estado da coleção?
- quando os dados foram atualizados?
- quais dados não conseguimos observar?

### Arquitetura

Separar resposta externa de domínio interno:

**External API Response → Normalization → Player Snapshot**

O domínio do CrownPilot não deve depender diretamente do formato bruto do provedor.

### Freshness

Toda informação sincronizada possui semântica clara de freshness.

Snapshot antigo não deve parecer estado atual.

### Critérios de aceite

- [ ] sync completo do jogador funciona;
- [ ] coleção é normalizada para o recommendation engine;
- [ ] Arena / troféus são persistidos com timestamp;
- [ ] dados ausentes não quebram o perfil;
- [ ] freshness fica visível;
- [ ] sync repetido é idempotente;
- [ ] token da Clash Royale API permanece server-side;
- [ ] payload bruto e Player Snapshot normalizado possuem boundary explícito;
- [ ] política de freshness/retry está documentada.

### Dependências

- Fases 001 e 002.

### Handoff

Player Snapshot se torna o contrato de entrada da personalização.

---

## 004 — Dataset de meta e canonicalização de decks ⬜

### Objetivo

Construir uma base confiável de decks reais e contexto competitivo.

### Escopo

- estratégia de amostragem de jogadores / batalhas;
- ingestão de battle data;
- canonicalização de deck;
- representação de Evolutions / Heroes;
- agrupamento de decks equivalentes;
- frequência;
- win rate quando calculável;
- tamanho da amostra;
- faixa de Arena / troféus;
- janela temporal;
- confidence;
- decaimento após mudanças relevantes;
- versionamento de catálogo/temporada/balance context;
- estratégia de retenção e agregação;
- validação do custo e do modelo de armazenamento no Supabase PostgreSQL;
- possibilidade de componente analítico especializado se PostgreSQL deixar de
  ser adequado para dados brutos em escala.

### Regra central

Um deck só entra no recommendation engine como candidato quando existe evidência suficiente para tratá-lo como deck observado, e não combinação arbitrária.

### Estatística

Nenhuma taxa deve ser apresentada sem:

- população;
- período;
- tamanho da amostra;
- contexto competitivo.

### Critérios de aceite

- [ ] decks possuem identidade canônica;
- [ ] dataset distingue contexto competitivo;
- [ ] métricas carregam janela temporal;
- [ ] amostras pequenas são penalizadas ou descartadas;
- [ ] balance changes podem reduzir a relevância dos dados antigos;
- [ ] recommendation engine consulta candidatos eficientemente;
- [ ] mudanças de temporada/balance podem separar ou invalidar amostras;
- [ ] retenção e agregação possuem estratégia explícita;
- [ ] custo/query pattern foi validado para Supabase/PostgreSQL ou existe boundary para
  store analítico especializado.

### Dependências

- Fase 001;
- Fase 002 para fundação executável, configuração e boundaries de integração.

### Handoff

Entregar catálogo de candidatos com evidência suficiente para personalização.

---

## 005 — Recommendation Engine v1: Best Decks for You ⬜

### Objetivo

Entregar a primeira feature que define o CrownPilot.

> **Best Decks for You**

### Entrada

- Player Snapshot;
- Deck Candidates;
- Competitive Context.

### Saída

Lista ordenada de decks acompanhada do motivo da recomendação.

### Fit Score

O score deve começar simples e decomponível.

Componentes conceituais:

- Meta Strength;
- Collection Compatibility;
- Level Readiness;
- Arena Relevance;
- Confidence.

A fórmula final não é definida neste roadmap. A spec da fase decide pesos, normalização e forma de combinação com dados reais.

### Readiness

Cada recomendação deve possuir classificação compreensível:

- READY TO PLAY;
- 1 UPGRADE AWAY;
- 2 UPGRADES AWAY;
- MISSING EVOLUTION / HERO;
- NOT READY.

### Regras

- não usar LLM para calcular ranking;
- não gerar deck arbitrário;
- mesma entrada produz mesma saída;
- ausência de dados reduz confiança;
- ranking global não substitui contexto do jogador.

### Critérios de aceite

- [ ] engine recebe Player Snapshot normalizado;
- [ ] candidatos vêm de decks observados;
- [ ] score é determinístico;
- [ ] score é decomponível;
- [ ] readiness é calculada;
- [ ] usuário entende por que um deck foi recomendado;
- [ ] Arena / faixa competitiva influencia a relevância quando os dados suportarem isso.

### Dependências

- Fases 003 e 004.

### Resultado

O CrownPilot passa a responder sua pergunta central.

---

## 006 — Progress: Upgrade Planner ⬜

### Objetivo

Responder:

> **O que eu deveria melhorar agora?**

### Experiência

- **NOW** — decks que posso jogar hoje;
- **NEXT** — decks que ficam prontos com pequenas melhorias;
- **TARGET** — decks relevantes que fazem sentido perseguir como objetivo.

### Upgrade Impact

Cada upgrade deve ser avaliado pelo impacto que gera.

Exemplo conceitual:

- uma carta sobe de nível;
- quatro decks passam para READY;
- sete decks melhoram Level Readiness.

Isso deve pesar mais do que um upgrade que melhora apenas um deck pouco relevante.

### Boundary

Se Gold, Wild Cards ou outros recursos econômicos não estiverem disponíveis na fonte de dados, o engine não inventa disponibilidade.

Entrada manual pode ser considerada em spec futura.

### Critérios de aceite

- [ ] cada recomendação de upgrade possui impacto explicável;
- [ ] NOW / NEXT / TARGET derivam do estado real da conta;
- [ ] nenhuma recomendação assume recursos desconhecidos;
- [ ] usuário navega de upgrade para decks impactados;
- [ ] usuário navega de deck para upgrades necessários.

### Dependências

- Fase 005.

### Resultado

O CrownPilot deixa de ser apenas discovery e passa a guiar progressão.

---

## 007 — MVP Experience e Beta ⬜

### Objetivo

Transformar os engines em uma experiência coerente e utilizável regularmente.

### Onboarding

**Sign in → Link Player Tag → Sync Account → Best Decks for You**

### Home

Priorizar:

- estado atual da conta;
- Arena / troféus;
- Best Decks for You;
- NOW / NEXT / TARGET;
- próximo upgrade recomendado;
- freshness.

### Produto

- mobile-first / responsive;
- testes unitários e de integração para novas regras de domínio;
- E2E dos fluxos críticos mantidos em staging;
- loading e erro consistentes;
- estados vazios;
- cache;
- sync controlado;
- deep links quando fizer sentido;
- analytics de produto;
- logs e observabilidade;
- privacidade;
- fluxo de exclusão de conta/dados;
- revisão de autorização e RLS/grants;
  - manter dados CrownPilot fora do acesso client-side à Supabase Data API;
- disclaimer de conteúdo não oficial.

### Métricas iniciais

Medir:

- Player Tag vinculada;
- sync concluído;
- recomendações geradas;
- deck visualizado;
- deck copiado / aberto quando aplicável;
- upgrade recomendado visualizado;
- retorno após mudança de progressão.

Metas numéricas entram na spec de beta depois de existir baseline.

### Critérios de aceite

- [ ] fluxo completo funciona em dispositivo novo após login;
- [ ] recomendações aparecem sem reentrada da tag;
- [ ] falha de provedor degrada de forma controlada;
- [ ] observabilidade mínima existe;
- [ ] eventos críticos são medidos;
- [ ] compliance e disclaimer estão visíveis;
- [ ] usuário possui caminho funcional para excluir sua conta e dados próprios;
- [ ] revisão de segurança do acesso ao PostgreSQL/Supabase foi concluída;
- [ ] MVP funciona sem IA.

### Dependências

- Fases 002–006.

### Marco

Primeira versão capaz de validar utilidade real com usuários.

---

## 008 — Histórico, Performance e Matchups ⬜

### Objetivo

Adicionar contexto do próprio jogador às recomendações.

### Escopo

- snapshots contínuos de battle history;
- retenção de histórico próprio quando permitido;
- classificação de decks / arquétipos;
- performance por deck;
- performance por matchup;
- tendências;
- período selecionável;
- confiança por volume.

### Boundary

A análise é pós-partida.

Não existe objetivo de assistência competitiva em tempo real.

### Critérios de aceite

- [ ] histórico pode persistir além da janela externa quando permitido;
- [ ] partidas são classificáveis;
- [ ] métricas mostram sample size;
- [ ] usuário consegue distinguir tendência de ruído.

### Dependências

- Fases 001, 003 e 004.

---

## 009 — Optimize My Deck ⬜

### Objetivo

Responder:

> **Como posso melhorar o deck que eu já jogo?**

### Estratégia

A primeira versão procura **variantes comprovadas**.

A recomendação considera:

- nível da alternativa na conta;
- evidência da variante;
- matchups relevantes;
- custo de oportunidade;
- contexto competitivo.

### Não fazer

Não pedir a um LLM para simplesmente inventar uma substituição.

### Critérios de aceite

- [ ] variante possui evidência;
- [ ] substituição é compatível com a coleção;
- [ ] motivo da troca é explicável;
- [ ] comparação mostra trade-offs.

### Dependências

- Fases 004, 005 e 008.

---

## 010 — Personalização por desempenho do jogador ⬜

### Objetivo

Fazer contas semelhantes receberem recomendações diferentes quando o histórico justificar isso.

### Escopo

- performance pessoal;
- arquétipos mais utilizados;
- consistência;
- matchups;
- confiança mínima;
- ajuste limitado do ranking.

### Regra

Personalização histórica não domina o recommendation engine com amostra pequena.

Sem evidência suficiente, prevalece:

**account + meta**

Com evidência suficiente:

**account + meta + player performance**

### Dependências

- Fases 008 e 009.

---

## 011 — AI Coach ⬜

### Objetivo

Adicionar linguagem natural e coaching sobre resultados já calculados pelo CrownPilot.

### Arquitetura

**Player Data → Recommendation / Analytics Engines → Structured Insight → AI Coach**

O modelo não deve reconstruir o domínio a partir de dados brutos quando o backend já consegue produzir resposta estruturada.

### Casos de uso

- Por que este deck é melhor para mim?
- Por que devo melhorar esta carta?
- O que mudou nas minhas últimas partidas?
- Por que estou sofrendo contra determinado arquétipo?
- Qual é meu plano para chegar ao próximo objetivo?

### Saúde operacional

A feature deve possuir:

- budget por chamada;
- limite de contexto;
- cache quando seguro;
- telemetria de custo;
- limites de uso;
- fallback sem IA;
- isolamento do provedor.

### Monetização

Qualquer plano pago, crédito, franquia de IA ou feature premium depende do gate comercial da Fase 001 e da decisão posterior da Fase 012. O resultado da Fase 001 deve ser consumido como constraint; a Fase 012 não substitui aprovação exigida pela Supercell nem autoriza cobrança por inferência.

Não construir billing primeiro e procurar justificativa depois.

### Critérios de aceite

- [ ] AI Coach usa resultados estruturados;
- [ ] resposta pode ser rastreada até dados do engine;
- [ ] indisponibilidade do modelo não quebra o core;
- [ ] custo por interação é observável;
- [ ] hallucination não altera scores ou fatos persistidos;
- [ ] modelo comercial, se existir, está validado.

### Dependências

- Fases 005–010 conforme o caso;
- gate de compliance.

---

## 012 — Comercialização e escala ⬜

### Objetivo

Definir um modelo sustentável somente depois de provar utilidade e validar o que é permitido.

### Gate obrigatório

Antes de cobrar:

- revisar Fan Content Policy vigente;
- revisar acordos da API / developer;
- confirmar enquadramento das features;
- registrar a decisão;
- buscar autorização adicional quando necessária.

Este gate consome a decisão de compliance produzida na Fase 001. A Fase 012 não é autorização implícita para cobrar: qualquer modelo classificado como `REQUIRES EXPLICIT APPROVAL` permanece bloqueado até existir aprovação expressa e rastreável da Supercell. Se a decisão da Fase 001 exigir aprovação ou mantiver o enquadramento ambíguo, a Fase 012 deve planejar operação gratuita ou modelo já permitido sob condições, sem criar billing como premissa.

### Princípio

Arquitetura de custos deve existir mesmo se o produto continuar gratuito.

Especialmente para IA, observar:

- cost per active user;
- cost per AI interaction;
- cache hit rate;
- provider spend;
- usage distribution.

### Regra

Não assumir no roadmap que assinatura SaaS tradicional, paywall por funcionalidades, AI Coach pago ou venda de analytics são permitidos. Ads, donations, coaching humano, software coaching e AI/software coaching devem respeitar a classificação e as condições registradas na Fase 001.

### Dependências

- utilidade comprovada no beta;
- métricas reais;
- gate jurídico / policy;
- custos medidos.

---

# 5. Explorações futuras

Estas ideias não fazem parte do MVP.

## 5.1 Build Me a Deck

Geração de novos candidatos considerando:

- cartas obrigatórias;
- sinergias;
- função das cartas;
- custo de elixir;
- matchups;
- coleção;
- níveis;
- dados históricos.

Qualquer deck gerado deve ser validado contra dados antes de ser apresentado como competitivo.

## 5.2 Alertas de progressão

Exemplos:

- nova Arena detectada;
- novos decks ficaram disponíveis;
- upgrade alterou ranking;
- mudança de meta afetou recomendações.

## 5.3 Planejamento avançado

- múltiplos objetivos;
- custo de progressão;
- caminhos alternativos;
- comparação entre estratégias de upgrade.

Depende da disponibilidade real de dados econômicos.

---

# 6. Mapa de dependências

Fluxo principal:

**001 Viabilidade / Dados / Compliance**

A partir dela:

**002 Fundação + Identidade**

Depois da fundação, 003 e 004 permanecem bloqueadas até reabertura explícita dos
gates de API data, retenção, ownership, egress, meta e compliance definidos no
veredito da Fase 001. Quando liberadas,
podem avançar em paralelo:

- **003 Player Sync**
- **004 Meta Dataset**

As Fases 003 e 004 convergem em:

**005 Best Decks for You → 006 Upgrade Planner → 007 MVP / Beta**

Depois:

**008 History & Matchups → 009 Optimize My Deck → 010 Player Personalization → 011 AI Coach**

Comercialização depende de:

**001 (decisão de compliance) + evidência do Beta + policy validation → 012 Commercialization & Scale**

Fase 012 consome esse resultado; não transforma dependência de aprovação em autorização.

---

# 7. Sequência crítica do MVP

1. **001** — provar dados e boundaries;
2. **002** — bootstrap, segurança e identidade persistente;
3. **003 + 004 em paralelo** — Player Snapshot e dataset competitivo;
4. **005** — Best Decks for You;
5. **006** — Upgrade Planner;
6. **007** — MVP / Beta.

O maior risco técnico não é UI.

É provar que conseguimos obter e manter dados suficientes para:

1. compreender a conta;
2. construir candidatos de meta confiáveis;
3. contextualizar por Arena / faixa;
4. explicar a recomendação.

Por isso a Fase 001 vem antes de congelar framework, modelagem física,
pipelines e demais detalhes de implementação.

---

# 8. Grandes marcos

## Marco A — Evidence

**Fase 001**

Sabemos quais dados temos, o que não temos, e que o caminho inicial honesto é
`Best Decks for Your Collection`; meta representativa da Arena continua dependente
de nova evidência ou fonte licenciada.

## Marco B — Connected Account

**Fases 002–003**

Aplicação possui fundação reproduzível; usuário entra em qualquer dispositivo e
o CrownPilot conhece e sincroniza sua conta.

## Marco C — Competitive Dataset

**Fase 004**

Temos decks reais, contexto e confiança suficientes para alimentar o engine.

## Marco D — Best Decks for You

**Fase 005**

A proposta central funciona.

## Marco E — Progression Intelligence

**Fase 006**

O produto recomenda não apenas o que jogar, mas o que melhorar.

## Marco F — Usable Beta

**Fase 007**

Existe experiência completa capaz de gerar aprendizado real.

## Marco G — Player Intelligence

**Fases 008–010**

Histórico e comportamento passam a melhorar recomendações.

## Marco H — Coach

**Fase 011**

IA adiciona interpretação sem substituir o engine.

## Marco I — Sustainable Product

**Fase 012**

Modelo operacional e comercial está validado técnica e legalmente.

---

# 9. Boundaries do MVP

Este mapa preserva escopo funcional. A integração EF Core/Npgsql abaixo é
baseline histórico; destino aprovado pela ADR 005 usa `pgx`/`sqlc`/`goose` após
cutover.

O MVP inclui:

- bootstrap reproduzível e quality gates;
- Firebase Authentication com Google;
- PostgreSQL/RLS seguro via `pgx`/`sqlc`/`goose` e ambientes separados;
- frontend estático e API Docker com hosting portátil;
- identidade CrownPilot;
- Player Tag persistente;
- sync da conta;
- Arena / troféus;
- coleção normalizada;
- dataset inicial de decks;
- Best Decks for You;
- Deck Readiness;
- Upgrade Planner;
- NOW / NEXT / TARGET;
- experiência responsiva;
- observabilidade básica.

O MVP não inclui:

- múltiplas Player Tags por usuário;
- AI Coach;
- geração livre de deck;
- overlay;
- tracking em tempo real durante partida;
- automação;
- billing;
- assinatura;
- social;
- clãs;
- torneios;
- app nativo obrigatório.

---

# 10. Roadmap, Specs, ADRs e repositório

## Roadmap

Responde:

> **O quê e em qual ordem?**

Este arquivo é o mapa de evolução.

## Specs

Respondem:

> **Como esta fase será implementada?**

Uma spec deve registrar, no mínimo:

- baseline;
- problema;
- objetivo;
- não objetivos;
- decisões;
- arquitetura;
- contratos;
- persistência;
- tasks ordenadas;
- critérios de aceite;
- riscos;
- testes;
- observabilidade;
- rollout;
- dependências;
- handoff.

## ADRs

Registram decisões arquiteturais que precisam sobreviver à implementação.

Exemplos futuros:

- estratégia de aquisição de meta;
- Player Snapshot;
- algoritmo do Fit Score;
- freshness;
- retenção de battle history;
- isolamento do provider de IA.

## Repositório

A implementação na main é a fonte final de verdade sobre o que existe.

Se roadmap, spec e código divergirem, investigar antes de continuar.

Não atualizar documentação para fingir que algo foi entregue.

---

# 11. Workflow por fase

Antes de planejar:

1. ler este roadmap;
2. inspecionar a main;
3. ler specs anteriores relevantes;
4. ler ADRs relacionados;
5. ler o handoff anterior;
6. confirmar baseline real;
7. identificar spikes necessários.

Durante a fase:

1. manter escopo explícito;
2. registrar decisões;
3. validar por checkpoints;
4. evitar antecipar fases futuras sem necessidade;
5. manter critérios de aceite verificáveis.

Ao finalizar:

1. validar testes;
2. validar comportamento real;
3. atualizar documentação;
4. registrar riscos e débitos;
5. produzir handoff.

---

# 12. Handoff padrão

Toda fase concluída termina com:

### Status

O que foi concluído e o que permaneceu fora.

### Resultado principal

Qual capacidade nova existe.

### Decisões

Quais decisões importantes foram tomadas.

### Arquitetura

Quais contratos e boundaries passam a importar.

### Validação

Quais testes, checks e evidências confirmam a entrega.

### Riscos e débitos

O que ainda pode afetar a próxima fase.

### Dependências

O que a próxima fase pode assumir como verdade.

### Próximo passo

Qual fase está liberada e por quê.

---

# 13. Próximo passo

A próxima fase ativa é:

> **002 — Fundação da aplicação e identidade persistente**

A Fase 001 foi encerrada com constraints. A Fase 002 pode começar pelo bootstrap,
pelos quality gates e pela identidade persistente, aplicando o vínculo de perfil
público sem alegar ownership.

O próximo deliverable técnico deve respeitar os gates ainda abertos e preparar a
validação futura de:

1. isolamento de identidade, vínculo e autorização;
2. boundary server-side e adapter substituível para a API;
3. condições para liberar sync/persistência somente após revisão dos débitos.

O fallback de produto permanece **Best Decks for Your Collection** até existir
evidência suficiente para qualquer alegação de meta da Arena.
