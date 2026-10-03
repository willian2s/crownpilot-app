# 002-09 — Validar staging, deploy, E2E, smoke e handoff

- **Ticker:** `002`
- **Número:** `09`
- **Status:** `pending`

## Objetivo e resultado esperado

Provar a fundação em staging e production controlados, com Preview efêmero,
Google Sign-In real em staging, E2E do fluxo crítico, smoke pós-deploy e handoff
explícito para a próxima fase.

## Requisitos cobertos

- Preview sem Supabase Auth real;
- staging com hostname fixo e projeto separado;
- E2E de login/vínculo/recuperação/troca/unlink;
- production somente pela `main`;
- deploy Vercel somente após prova de runtime Laravel/PHP;
- deploy pelo container oficial com `Dockerfile.vercel`, FrankenPHP, `Caddyfile`
  e `vercel.json`;
- imagem e runtime padronizados em PHP 8.5;
- fallback para host PHP first-class sem alterar contratos;
- CI GH Actions verde para Pull Request e push em `main`;
- smoke em staging e production;
- riscos e gates remanescentes documentados.

## Escopo incluído

- Preview Deployment com build e smoke sem Supabase real;
- ambiente staging com hostname fixo e variáveis próprias;
- spike de Laravel/Inertia/Vite, `pdo_pgsql`, SSL, sessão e PostgreSQL no container
  escolhido;
- `.dockerignore` validado e imagem sem `.env*`, secrets ou testes;
- redirect allowlist desse host no Supabase Auth staging;
- deploy production a partir de `main`, em `sa-east-1`, após gates PHP, Supabase
  e secrets aprovados;
- E2E com conta Google/test data de staging;
- smoke controlado pós-deploy em staging e production;
- evidência de isolamento de projetos e ausência de secrets no client;
- registro de handoff, débitos e bloqueios da Fase 003.

## Escopo excluído

- liberar sync completo ou persistência de Player Snapshot;
- crawler, polling global, meta Arena ou recommendation engine;
- billing, AI Coach ou aprovação comercial;
- migração de dados legados inexistentes.

## Dependências

- `002-01` a `002-08` concluídas;
- projetos Supabase e disponibilidade de `sa-east-1` confirmados;
- hostname staging, Google provider e secrets próprios disponíveis;
- branch `staging` e `main` protegidas conforme workflow do repositório.

## Arquivos e símbolos prováveis

- configuração Vercel ou host PHP fallback, sem credenciais versionadas;
- `tests/e2e/`, `tests/smoke/` e fixtures de staging;
- documentação de release/rollback;
- `docs/tasks/002.../002-00-overview.md` e registro de execução;
- handoff final na spec ou evidência da fase.

## Passos de implementação

1. Criar Preview e comprovar build sem Auth real.
2. Provar `Dockerfile.vercel` com FrankenPHP/PHP 8.5, `Caddyfile`, `vercel.json`
   e `/health` localmente e no GH Actions.
3. Provar Laravel/Inertia/Vite, sessão e PostgreSQL no container Vercel ou host
   PHP first-class.
4. Configurar staging fixo, Supabase separado e Google Sign-In.
5. Executar E2E do fluxo completo em staging.
6. Corrigir isolamento, acessibilidade, erros e redaction encontrados.
7. Promover somente via `main` para production em `sa-east-1`.
8. Executar smoke controlado sem expor dados de teste.
9. Registrar resultados, riscos residuais e condição de handoff para Fase 003.

## Testes e comandos de validação

```text
composer install
npm ci
php artisan test
docker build -f Dockerfile.vercel -t crownpilot-vercel-ci .
npm run build
npm run test:e2e
npm run smoke
```

E2E deve confirmar login, vínculo, reload/outro dispositivo, troca, falha de
provider sem perda, unlink e logout/login. Smoke deve confirmar Auth, endpoint de
vínculo e isolamento de ambiente; não deve fazer ingestão ou polling.

## Definição de pronto

- Preview funciona sem Google Sign-In real e sem Supabase production;
- runtime Laravel/PHP, sessão e PostgreSQL/Supabase passam no host escolhido;
- local, GH Actions e container usam PHP 8.5;
- imagem Vercel/FrankenPHP e `/health` passam no GH Actions;
- staging tem host fixo, projeto separado e login Google funcional;
- E2E crítico passa em staging;
- production só recebe deploy da `main` e usa `sa-east-1`;
- smoke passa em staging e production;
- nenhuma credencial/token aparece no browser ou logs;
- overview/spec registram 9/9 subtarefas, critérios, decisões, riscos e handoff;
- Fase 003 continua explicitamente bloqueada até gates da Fase 001 serem
  reabertos e aprovados.

## Riscos e cuidados

- Não autorizar hosts Preview efêmeros no Supabase Auth.
- Não compartilhar Supabase production com staging ou Preview.
- Não declarar Vercel compatível com Laravel sem evidência de runtime, extensões,
  sessão e PostgreSQL.
- Se Vercel falhar, usar host PHP first-class sem alterar domínio/contratos.
- Não usar smoke para mascarar falta de E2E.
- Não declarar Fase 003 liberada somente porque identidade e deploy passaram.
