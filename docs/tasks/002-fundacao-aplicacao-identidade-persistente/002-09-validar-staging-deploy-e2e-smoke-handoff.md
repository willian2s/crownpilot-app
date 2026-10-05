# 002-09 — Validar Staging, deploy, E2E, smoke e handoff

- **Ticker:** `002`
- **Número:** `09`
- **Status:** `pending`

## Objetivo e resultado esperado

Provar a fundação em Staging e Production controlados, mantendo Preview
efêmero, Google Sign-In real em Staging, E2E do fluxo crítico, smoke pós-deploy
e handoff explícito para a próxima fase.

## Requisitos cobertos

- Preview sem Firebase Auth real;
- Staging com hostname fixo, projeto Firebase e banco Supabase separados;
- E2E de login/vínculo/recuperação/troca/unlink;
- Production somente por promoção controlada de `main`;
- API ASP.NET Core em imagem Docker/OCI portátil;
- frontend Vite estático hospedado opcionalmente em Vercel;
- migrations controladas e health endpoints;
- CI verde, smoke e riscos remanescentes documentados.

## Escopo incluído

- Preview com build e smoke sem login real;
- API em host containerizado compatível e frontend `dist/` em host estático;
- `.dockerignore` validado e imagem sem `.env*`, secrets ou testes;
- Firebase project de Staging com Google Sign-In e allowlist do hostname fixo;
- Supabase project/database de Staging separado, com SSL e conexão validados;
- migration job controlado antes da aplicação, sem auto-migration concorrente;
- E2E com conta Google/test data de Staging;
- smoke controlado pós-deploy em Staging e Production;
- evidência de isolamento, redaction e ausência de secrets no client;
- registro de handoff, débitos e bloqueios da Fase 003.

## Escopo excluído

- liberar sync completo ou persistência de Player Snapshot;
- crawler, polling global, meta Arena ou recommendation engine;
- billing, AI Coach ou aprovação comercial;
- migração de dados legados inexistentes;
- tornar Vercel requisito da API.

## Dependências

- `002-01` a `002-08` concluídas;
- projetos Firebase/Supabase e região aprovados;
- hostname Staging, Google provider e secrets próprios disponíveis;
- branch `staging` e `main` protegidas conforme workflow.

## Arquivos e símbolos prováveis

- configuração do host containerizado e frontend estático;
- `tests/E2E/`, `tests/Smoke/` e fixtures de Staging;
- release/rollback e migration runbook;
- overview/spec e registro de execução;
- health checks `/health/live` e `/health/ready`.

## Passos de implementação

1. Criar Preview e comprovar build/smoke sem login real.
2. Construir imagem ASP.NET Core e provar health localmente e no CI.
3. Publicar frontend Vite em host estático, sem dependência exclusiva de Vercel.
4. Configurar Staging fixo, Firebase separado, Supabase separado e secrets
   próprios.
5. Aplicar migrations EF Core via job controlado e SQL RLS/grants posterior.
6. Executar E2E do fluxo completo em Staging.
7. Corrigir isolamento, acessibilidade, erros, CORS e redaction.
8. Promover somente via `main` para Production após aprovação dos gates.
9. Executar smoke controlado e registrar handoff/riscos residuais.

## Testes e comandos de validação

```text
dotnet restore
dotnet test --configuration Release
npm ci
npm run build
docker build -t crownpilot-api:release .
npm run test:e2e
npm run smoke
```

E2E deve confirmar login, vínculo, reload/outro dispositivo, replace, falha de
provider sem perda, unlink e logout/login. Smoke confirma build, health,
authentication, endpoint de vínculo e isolamento de ambiente; não faz ingestão ou
polling.

## Definição de pronto

- Preview funciona sem Google Sign-In real e sem projetos de Production;
- API containerizada e frontend estático passam build e health;
- migrations são aplicadas por job controlado;
- Staging tem host fixo, Firebase/Supabase separados e login Google funcional;
- E2E crítico passa em Staging;
- Production só recebe promoção aprovada da `main`;
- smoke passa sem expor dados de teste ou secrets;
- overview/spec registram 9/9 subtarefas, decisões, riscos e handoff;
- Fase 003 continua bloqueada até gates da Fase 001 serem reabertos.

## Riscos e cuidados

- Não adicionar hosts Preview efêmeros ao allowlist de login real.
- Não compartilhar banco, Firebase ou secrets entre ambientes.
- Não declarar API pronta por validar apenas hospedagem do frontend.
- Não aplicar migration automaticamente em múltiplas réplicas.
- Não usar smoke para mascarar falta de E2E.
- Não declarar Fase 003 liberada somente porque identidade e deploy passaram.
