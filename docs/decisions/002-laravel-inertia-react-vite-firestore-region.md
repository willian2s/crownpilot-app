# ADR 002 — Laravel, Inertia, React, Vite e região do Firestore

- **Status:** Accepted — Firebase/Firestore portions superseded by ADR 003
- **Data:** 2026-10-02
- **Escopo:** stack da aplicação, deploy e localização do Firestore
- **Fase:** 002 — Fundação da aplicação e identidade persistente
- **Relacionada:** [ADR 001](001-firebase-firestore-vercel-portable.md)
- **Supersession:** [ADR 003](003-supabase-auth-postgresql-jwks-rls-region.md)

As decisões de Laravel/PHP, Inertia/React/Vite, PHP 8.5 e Vercel via Docker
continuam válidas. Auth, persistência e região Firebase/Firestore foram
substituídas por Supabase Auth/PostgreSQL, conforme ADR 003.

## Contexto

O repositório ainda não possui runtime. A Fase 002 precisa criar uma aplicação
executável, aprofundar PHP/Laravel e manter os boundaries já definidos para
Firebase, Firestore, Clash Royale API e portabilidade de deploy.

O deploy Laravel na Vercel será baseado no processo oficial
[Deploy Laravel on Vercel with Docker](https://vercel.com/kb/guide/laravel-php-with-docker),
que usa imagem containerizada com FrankenPHP, Caddy e configuração de service
container.

## Decisão

### Stack da aplicação

- Laravel como backend e camada HTTP;
- PHP 8.5 como versão padronizada local, CI e produção;
- Inertia como integração server-driven;
- React + TypeScript como frontend;
- Vite para desenvolvimento e build frontend;
- Composer + `composer.lock` para dependências PHP;
- npm + `package-lock.json` para dependências JavaScript;
- Pest/PHPUnit, Laravel Pint, PHPStan/Larastan, Vitest/React Testing Library e
  Playwright conforme a camada de teste.

Domínio e aplicação permanecem independentes de controllers Laravel, Inertia,
Firebase e Vite. Controllers, Form Requests e páginas Inertia são adapters
finos.

### Firebase e identidade

Firebase Authentication continua provider de identidade CrownPilot com Google.
Firebase Web SDK roda no React somente para login. O Laravel verifica o Firebase
ID token através de adapter PHP server-side e estabelece sessão Laravel segura.
O token não é enviado ao provider Clash Royale.

Firestore continua banco principal, acessado somente pelo backend através de
adapter PHP compatível com Firebase/IAM. A implementação deve validar a
compatibilidade das bibliotecas escolhidas e `ext-grpc` durante o bootstrap, sem
alterar os contratos do domínio.

### Região

Firestore usará região regional `southamerica-east1` em staging e production.
Staging usará a mesma região para manter paridade. Região é decisão difícil de
reverter; mudança futura exige novo projeto/database e cutover controlado.

### Deploy Vercel

Vercel continua plataforma inicial. Laravel será empacotado conforme guia oficial
com:

- `Dockerfile.vercel` baseado em FrankenPHP/PHP 8.5;
- `Caddyfile` expondo somente `public/` e encaminhando para `public/index.php`;
- `vercel.json` declarando service container e rewrite catch-all;
- `.dockerignore` excluindo `.env*`, credenciais, `vendor`, `node_modules` e
  testes;
- environment variables injetadas pela Vercel em runtime;
- `APP_KEY` estável e distinto por ambiente;
- `/health` para smoke e validação do container.

Filesystem local do container é efêmero. Estado durável fica no Firestore.
Sessões não podem depender de arquivos locais ou memória de uma instância;
sessão cookie criptografada é opção inicial compatível com o escopo.

Deploy production permanece restrito à `main`. Preview pode usar `vercel deploy`
sem `--prod`; production usa integração Git ou `vercel deploy --prod` somente
após CI verde e gates de staging.

### CI

GitHub Actions é CI autoritativo e executa em:

- `pull_request`;
- `push` para `main`.

Gates incluem Composer, npm, lint, análise estática, type-check, testes unitários,
integração, contrato, Rules, build frontend, build de `Dockerfile.vercel` e smoke
de `/health`. CI não chama a Clash Royale API live.

## Alternativas consideradas

### Next.js

Não adotado. Laravel/PHP foi escolhido para aprofundamento de arquitetura e
mantém Inertia como boundary server-driven.

### Laravel Blade

Não adotado. Abandonaria React, TypeScript e Vite definidos para frontend.

### SPA React/Vite com backend separado

Não adotado. Adicionaria deploy, CORS e boundary operacional sem necessidade no
baseline monolítico.

### Firebase client SDK para Firestore

Não adotado. Aumentaria superfície pública e complexidade de autorização. Browser
usa Firebase Web SDK somente para Auth.

### Firestore multi-region

Não adotado para Fase 002. Custo e complexidade não justificam a escolha antes de
volume/SLO reais.

### Runtime PHP não containerizado na Vercel

Não adotado. O processo oficial containerizado com FrankenPHP fornece runtime
explícito e reproduzível.

O guia de referência usa exemplo PHP 8.4. CrownPilot adapta a imagem para PHP
8.5; a disponibilidade da tag FrankenPHP/PHP 8.5 e suas extensões é gate
obrigatório antes de staging/production.

## Consequências

### Positivas

- prática de PHP/Laravel sem abandonar React/TypeScript;
- domínio continua portável para outro host PHP;
- Vercel recebe imagem reproduzível e smoke verificável;
- Composer e npm tornam dependências determinísticas;
- região e paridade staging/production ficam explícitas;
- CI valida código e container antes de release.

### Trade-offs

- dois ecossistemas de dependências e testes;
- build Docker exige Docker disponível no CI;
- runtime FrankenPHP exige validação de extensões, especialmente `ext-grpc`;
- filesystem efêmero impede estado local compartilhado;
- cold starts e limites de container precisam ser observados;
- mudança futura de região será migração controlada.

## Gates de implementação

- `composer install` e `npm ci` reproduzem aplicação;
- PHP 8.5 é confirmado em local, CI e imagem de produção;
- tag FrankenPHP/PHP 8.5 e extensões necessárias são resolvidas no build;
- container Vercel inicia e `/health` retorna `200`;
- Caddy expõe somente `public/`;
- imagem não contém `.env`, token, credencial, testes ou dependência local;
- Laravel verifica Firebase ID token;
- Firestore grava transação em `southamerica-east1`;
- Emulator Suite funciona com testes PHP/JS;
- sessão segura persiste entre invocações;
- staging usa projeto Firebase separado;
- GitHub Actions passa em Pull Request e push para `main`;
- production só recebe deploy após gates verdes.

## Revisão

Revisar ADR se incompatibilidade demonstrada de FrankenPHP, Vercel, extensões
PHP, Firebase, segurança, custo, compliance ou egress exigir mudança. Preferência
por outra stack, isoladamente, não invalida esta decisão.
