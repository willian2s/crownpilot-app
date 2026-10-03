# 002-01 — Bootstrap do toolchain

- **Ticker:** `002`
- **Número:** `01`
- **Status:** `pending`

## Objetivo e resultado esperado

Criar aplicação web reproduzível a partir de clone limpo, com Laravel, Inertia,
React, TypeScript, Vite, PHP 8.5, Composer/npm e comandos básicos de
qualidade. O resultado é um baseline executável, ainda sem integração real de
Supabase Auth, PostgreSQL ou Clash Royale.

## Requisitos cobertos

- bootstrap reproduzível;
- Laravel/Inertia/React/Vite, Composer e npm com lockfiles;
- TypeScript strict, build, lint, format e type-check;
- estrutura inicial de diretórios e boundaries;
- `AGENTS.md` e comandos operacionais;
- nenhum uso obrigatório de API proprietária Vercel;
- deploy Laravel/FrankenPHP via container Vercel conforme guia oficial;
- gate técnico explícito para build e smoke do container.

## Escopo incluído

- inicializar Laravel e Inertia com React/TypeScript/Vite;
- adicionar `composer.json`, `composer.lock`, `package.json`, `package-lock.json`,
  config TypeScript, lint e formatter;
- definir scripts Composer e npm para `dev`, `build`, `lint`, `analyse`,
  `typecheck` e testes;
- criar `Dockerfile.vercel`, `Caddyfile`, `vercel.json` e `.dockerignore` sem
  secrets na imagem;
- criar separação inicial `app/Domain`, `app/Application`, `app/Adapters`,
  `app/Http`, `resources/js` e `tests`;
- documentar versão/política PHP e bootstrap limpo;
- criar `AGENTS.md` local com comandos e limites do projeto.

## Escopo excluído

- Google Sign-In, Supabase/PostgreSQL e API externa;
- provisionamento de projetos ou secrets;
- schema persistente, vínculo e UI final;
- deploy de staging/production.

## Dependências

- Fase 001 e ADR 001;
- nenhuma dependência de código existente.

## Arquivos e símbolos prováveis

- `composer.json`, `composer.lock`, `package.json`, `package-lock.json`;
- `vite.config.ts`, `tsconfig.json`, `phpunit.xml`, `pest.php`, configs de lint;
- `Dockerfile.vercel`, `Caddyfile`, `vercel.json`, `.dockerignore`;
- `AGENTS.md`, `app/Domain/`, `app/Application/`, `app/Adapters/`,
  `app/Http/`, `resources/js/`;
- `tests/` e `README.md`/documentação operacional.

## Passos de implementação

1. Fixar PHP 8.5 e escolher versões compatíveis de Laravel, Inertia, React,
   TypeScript e Vite.
2. Gerar `composer.lock` e `package-lock.json`; confirmar install limpo.
3. Ativar tipagem estrita, análise estática, lint e scripts normativos.
4. Criar página Inertia mínima e `/health` que provam execução Laravel/Vite.
5. Criar imagem FrankenPHP seguindo guia oficial e testar `public/index.php`.
6. Verificar que domínio não importa Laravel, Inertia, Supabase ou Vercel.
7. Registrar comandos e limites no `AGENTS.md` local.

## Testes e comandos de validação

```text
composer install
npm ci
composer run lint
composer run analyse
composer run test:unit
npm run typecheck
npm run build
docker build -f Dockerfile.vercel -t crownpilot-vercel-ci .
```

Repetir `composer install` e `npm ci` em clone limpo e confirmar ausência de
dependência em arquivos locais não versionados. Iniciar container e testar
`/health` sem incluir `.env`, token ou credencial na imagem.

## Definição de pronto

- clone limpo executa Composer, npm, lint, análise, type-check e build;
- `composer.lock` e `package-lock.json` são versionados;
- PHP 8.5/runtime está documentado e compatibilidade Vercel está em gate;
- imagem FrankenPHP/PHP 8.5 resolve e constrói;
- `Dockerfile.vercel` constrói imagem e `/health` responde `200`;
- `Caddyfile` limita document root a `public/`;
- `vercel.json` declara service container e rewrite catch-all;
- `.dockerignore` exclui secrets, dependências locais e testes;
- estrutura e comandos estão documentados;
- nenhum secret ou token entra no bundle;
- diff não contém dependência Vercel-specific obrigatória.

## Riscos e cuidados

- Não adicionar service role, database password ou API token no client durante
  scaffold.
- Não transformar Laravel/Inertia layout em contrato de domínio.
- Não instalar biblioteca de testes sem script/uso planejado.
- Manter escopo no bootstrap; Auth e vínculo pertencem às subtarefas seguintes.
