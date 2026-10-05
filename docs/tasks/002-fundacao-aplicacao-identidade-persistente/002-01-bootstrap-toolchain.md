# 002-01 — Bootstrap do toolchain

- **Ticker:** `002`
- **Número:** `01`
- **Status:** `pending`

## Objetivo e resultado esperado

Criar baseline reproduzível a partir de clone limpo com solution .NET, API ASP.NET
Core, frontend React + TypeScript + Vite, EF Core + Npgsql, testes e container
Docker. O resultado ainda não configura login real, banco remoto ou lookup.

## Requisitos cobertos

- SDK .NET pinado e solution versionada;
- .NET 10 LTS e Node.js 24 LTS pinados, com patches atualizados conforme suporte
  oficial vigente;
- camadas `API`, `Application`, `Domain` e `Infrastructure`;
- React, TypeScript strict e Vite independentes da API;
- EF Core e Npgsql referenciados somente na Infrastructure;
- testes .NET/frontend, lint, type-check, build e health;
- Docker/OCI portátil, sem dependência estrutural de Vercel;
- comandos de desenvolvimento e CI mínimo documentados.
- documentação Mac/Linux equivalente e aprendizado incremental em `docs/learning/`.

## Escopo incluído

- criar `global.json` com .NET 10, solution e projetos `.csproj`;
- fixar Node/npm/package manager por `.nvmrc`, `packageManager` ou equivalente;
- criar manifesto local de ferramentas `dotnet-ef` e registrar versões de Firebase
  CLI, Supabase CLI e Docker exigidas pelo setup;
- criar `src/Api`, `src/Application`, `src/Domain` e `src/Infrastructure`;
- configurar referências para impedir Infrastructure no Domain e API no Domain;
- criar frontend Vite separado, com `package.json`, lockfile e TypeScript strict;
- adicionar EF Core, Npgsql e tooling de migrations sem conectar a ambiente real;
- adicionar runner de testes .NET, testes frontend e endpoint health mínimo;
- definir scripts oficiais sem dependência de shell específico: `test:unit`,
  `test:contract`, `test:rls`, `test:e2e`, `smoke`, `smoke:container`,
  `openapi:check`, `typecheck`, `lint` e `build`;
- criar `Dockerfile` multi-stage para API e documentação de execução local;
- criar `.env.example`/settings sem secrets;
- criar `AGENTS.md` local com comandos, boundaries e limites de escopo;
- deixar workflow CI mínimo para restore, build, test, type-check e frontend build.

## Escopo excluído

- Google Sign-In ou validação real de Firebase ID Token;
- provisionamento de Firebase/Supabase, schema de negócio ou secrets;
- Player Tag, lookup, autorização de usuário ou UI final;
- deploy de Staging/Production;
- uso de SDK C# do Supabase.

## Dependências

- Fase 001 e [ADR 004](../../decisions/004-aspnet-core-react-vite-firebase-postgresql.md);
- nenhum runtime existente ou código legado.

## Arquivos e símbolos prováveis

- `global.json`, `CrownPilot.sln`, `src/*/*.csproj`;
- `src/Api/Program.cs`, `HealthEndpoints`, `DependencyInjection`;
- `src/Domain/`, `src/Application/`, `src/Infrastructure/`;
- `frontend/package.json`, `vite.config.ts`, `tsconfig.json`;
- `tests/Unit`, `tests/Application`, `frontend/src/**/*.test.ts`;
- `Dockerfile`, `.dockerignore`, `.env.example`, `.github/workflows/ci.yml`.

## Passos de implementação

1. Fixar .NET 10 LTS e Node.js 24 LTS, validando patches suportados na documentação
   oficial; registrar exceção somente se compatibilidade concreta exigir.
2. Criar solution/projetos e validar referências entre camadas.
3. Configurar frontend React/Vite independente, scripts e build `dist/`.
4. Adicionar EF Core/Npgsql na Infrastructure, manifesto `dotnet-ef` local e
   tooling sem migration de domínio ainda.
5. Criar health mínimo e teste unitário por camada sem provider externo.
6. Criar imagem Docker genérica da API, com configuração em runtime e filesystem
   efêmero.
7. Documentar comandos POSIX comuns a Mac/Linux e deixar CI mínimo executável em
   Pull Request. Scripts que precisem de lógica usam Node ou .NET, não Bash-only.

## Testes e comandos de validação

```text
dotnet --info
dotnet restore
dotnet build --configuration Release
dotnet test --configuration Release
npm ci
npm run typecheck
npm run lint
npm run build
docker build -t crownpilot-api:bootstrap .
docker run --rm -d --name crownpilot-api-bootstrap -p 8080:8080 crownpilot-api:bootstrap
curl --fail http://localhost:8080/health/live
docker stop crownpilot-api-bootstrap
```

Confirmar clone limpo, ausência de secrets na imagem/bundle e que Domain não
importa ASP.NET Core, EF Core, Npgsql, Firebase ou HTTP. `dotnet build` deve
gerar/verificar OpenAPI quando a API já possuir endpoints; a fundação não deve
manter YAML OpenAPI manual. `/health/live` não testa banco ou provider.

## Definição de pronto

- solution e SDK pinado restauram em clone limpo;
- Node/npm/package manager são pinados e reproduzidos em CI;
- API e frontend compilam separadamente;
- EF Core/Npgsql existem somente na camada prevista;
- testes .NET/frontend, type-check, lint e build possuem comandos reais;
- `/health/live` responde no processo e no container;
- Docker inicia sem secret embutido e sem depender de Vercel;
- CI mínimo bloqueia falha de restore/build/test/type-check/build;
- `AGENTS.md` documenta comandos e não permite avançar Auth/persistência nesta task;
- setup Mac/Linux usa fonte de verdade versionada (`global.json`, Node pinado,
  lockfile e `dotnet-tools.json`);
- Task registra primeiro conceito de aprendizado usado, sem criar curso paralelo;
- nenhum Auth, banco remoto, lookup ou fluxo de negócio é implementado nesta task.

## Riscos e cuidados

- Não criar abstração vazia sem consumidor previsto.
- Não adicionar pacote de provider ao Domain/Application.
- Não adicionar Firebase Admin ou configuração real antes da Task 04.
- Não aplicar migrations automaticamente no startup do container.
- Não criar arquivos de teste vazios apenas para satisfazer scripts.
