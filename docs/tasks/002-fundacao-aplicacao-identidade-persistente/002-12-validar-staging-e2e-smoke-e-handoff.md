# 002-12 — Validar Staging, E2E, smoke e handoff

- **Ticker:** `002`
- **Número:** `12`
- **Status:** `pending`

## Objetivo e resultado esperado

Provar a fundação em Staging controlado, mantendo Preview efêmero, Google Sign-In
real em Staging, E2E do fluxo crítico, smoke pós-deploy e handoff explícito para
a próxima fase. Production é promoção posterior e controlada, não pré-requisito
para bootstrap local.

## Requisitos cobertos

- Preview sem Firebase Auth real;
- Staging com hostname fixo, projeto Firebase e banco Supabase separados;
- E2E de login/vínculo/recuperação/troca/unlink;
- Production somente por promoção controlada de `main`, quando houver aprovação;
- API ASP.NET Core em imagem Docker/OCI portátil;
- frontend Vite estático hospedado inicialmente no Render Static Site ou alternativa;
- migrations controladas e health endpoints;
- CI verde, smoke e riscos remanescentes documentados.

## Escopo incluído

- Preview com build e smoke sem login real;
- API em host containerizado compatível e frontend `dist/` em host estático;
- `.dockerignore` validado e imagem sem `.env*`, secrets ou testes;
- Firebase project de Staging com Google Sign-In e allowlist do hostname fixo;
- Supabase project/database de Staging separado, com SSL e conexão validados;
- bundle/job final de migrations produzido pela Task 002-11, executado antes da
  aplicação, sem auto-migration concorrente;
- E2E com identidade de teste controlada e bearer Firebase; smoke/manual confirma
  Google Sign-In real e authorized domain, sem credencial persistida no repositório;
- smoke controlado pós-deploy em Staging; Production somente após go/no-go;
- promoção do mesmo digest OCI validado em Staging;
- evidência de isolamento, redaction e ausência de secrets no client;
- registro de handoff, débitos e bloqueios da Fase 003.

## Escopo excluído

- liberar sync completo ou persistência de Player Snapshot;
- crawler, polling global, meta Arena ou recommendation engine;
- billing, AI Coach ou aprovação comercial;
- migração de dados legados inexistentes;
- tornar Vercel requisito da API.

## Dependências

- `002-01` a `002-11` concluídas;
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
3. Publicar frontend Vite no Render Static Site ou host compatível, sem dependência
   exclusiva de Vercel.
4. Configurar Staging fixo, Firebase separado, Supabase separado e secrets
   próprios.
5. Executar exatamente o bundle/job final produzido pela Task 002-11 e aplicar
   SQL RLS/grants posterior.
6. Executar E2E do fluxo completo em Staging.
7. Corrigir isolamento, acessibilidade, erros, CORS e redaction.
8. Registrar decisão explícita de go/no-go para Production; se aprovada, promover
   via `main` exatamente o digest OCI publicado e validado em Staging, aplicar
   migration job e executar smoke não destrutivo. Não reconstruir a imagem.
9. Registrar handoff/riscos residuais.

## Testes e comandos de validação

```text
dotnet restore
dotnet test --configuration Release
npm ci
npm run build
npm run test:e2e
npm run smoke
```

E2E deve confirmar login por identidade de teste, vínculo, reload/outro dispositivo,
replace, falha de provider sem perda, unlink e logout/login. Smoke confirma build,
health, authentication, endpoint de vínculo e isolamento de ambiente; não faz
 ingestão ou polling. O ambiente deve estar executando o digest OCI validado pelo
 candidate. Google UI real pode ser confirmação manual controlada, não
dependência frágil de automação de terceiro.

## Definição de pronto

- Preview funciona sem Google Sign-In real e sem projetos de Production;
- API containerizada e frontend estático passam build e health;
- migrations são aplicadas por job controlado;
- Staging tem host fixo, Firebase/Supabase separados e login Google funcional;
- E2E crítico passa em Staging;
- Production só recebe promoção aprovada da `main`, caso go/no-go exista;
- Staging e Production usam digest OCI promovível, não rebuild divergente;
- smoke de Staging passa sem expor dados de teste ou secrets; smoke de Production
  é obrigatório somente após provisionamento e go/no-go explícitos;
- overview/spec registram 12/12 subtarefas, decisões, riscos e handoff;
- Fase 003 continua bloqueada até gates da Fase 001 serem reabertos.

## Riscos e cuidados

- Não adicionar hosts Preview efêmeros ao allowlist de login real.
- Não compartilhar banco, Firebase ou secrets entre ambientes.
- Não declarar API pronta por validar apenas hospedagem do frontend.
- Não aplicar migration automaticamente em múltiplas réplicas.
- Não usar smoke para mascarar falta de E2E.
- Não declarar Fase 003 liberada somente porque identidade e deploy passaram.
