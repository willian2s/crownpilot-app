# 002-11 — Automatizar CI, OCI e gates de release

- **Ticker:** `002`
- **Número:** `11`
- **Status:** `pending`

## Objetivo e resultado esperado

Automatizar gates de Pull Request/main e preparar candidate Staging com imagem
OCI imutável, sem secrets, mantendo observabilidade já definida na Task 002-10.

## Requisitos cobertos

- CI .NET com restore, build, analyzers e `dotnet test`;
- CI frontend com npm ci, lint, type-check, testes e build;
- testes de persistência PostgreSQL, EF migrations e RLS;
- contract tests do lookup sem chamadas live;
- build/smoke da imagem ASP.NET Core sem secrets;
- E2E/smoke com comandos e pré-condições claros;
- validação do documento OpenAPI gerado e drift de contrato;
- checks de observabilidade, health e redaction da Task 002-10;
- gates separados para PR, main e Staging.

## Escopo incluído

- workflow `.github/workflows/ci.yml` para `pull_request`, staging candidate e
  promoção em `main`;
- geração do bundle/job final de migrations a partir do schema atual, separado do
  container de runtime;
- jobs .NET, frontend, persistence/RLS, contract, image e health;
- PostgreSQL local/descartável e Firebase Emulator/fixtures sem produção;
- scripts de E2E e smoke staging-only;
- validação de logger/métricas, health e redaction já instrumentados;
- liveness sem dependências e readiness com configuração/PostgreSQL, nunca lookup
  live;
- scans para secrets, URLs com tag e imports server-only no bundle;
- documentação de owner/pré-condição dos gates Staging/Production.

## Escopo excluído

- provisionamento de Production pelo CI;
- polling/probes live automáticos;
- dashboards avançados ou analytics de produto;
- logging de payloads para depuração;
- chamadas live à Clash Royale API na CI normal.

## Dependências

- `002-01` a `002-10` para fechamento; gates mínimos começam em `002-01`;
- Docker/PostgreSQL e fixtures disponíveis em CI;
- secrets reais somente em ambientes controlados de Staging.

## Arquivos e símbolos prováveis

- `.github/workflows/ci.yml`;
- `tests/Unit`, `Application`, `Integration`, `Persistence`, `Contract`, `Api`;
- `src/Infrastructure/Observability/Redaction.cs` e métricas;
- scripts `dotnet`, `npm`, `database` e documentação de release.

## Passos de implementação

1. Definir jobs e falhas para PR e push em `main`.
2. Executar restore/build/test .NET e frontend sem secrets reais.
3. Subir PostgreSQL descartável, aplicar EF migrations, RLS e fixtures; executar
   checks de pool/reset e isolamento A/B.
4. Executar testes de token Firebase com Emulator/fixtures assinadas.
5. Executar architecture tests para provar `API → Application → Domain` e
   `Infrastructure → Application/Domain`, sem providers em Domain/Application.
6. Construir imagem Docker e testar `/health/live`/`ready` sem secret na imagem.
7. Gerar e versionar como artefato de candidate o bundle/job final de migrations;
   não reutilizar bundle produzido antes da Task 002-06.
8. Criar workflow protegido/manual de Staging com owner, aprovação, migration
   job, mesmo digest OCI e smoke antes de promoção para `main`.
9. Documentar gates Staging/Production fora dos gates de PR/main.

## Testes e comandos de validação

```text
dotnet restore
dotnet build --configuration Release
dotnet test --configuration Release
npm ci
npm run lint
npm run typecheck
npm run test:unit
npm run test:contract
npm run build
dotnet test --filter Category=Persistence
dotnet test --filter Category=Rls
dotnet test --filter Category=Api
npm run openapi:check
docker build -t crownpilot-api:ci .
npm run smoke:container -- --image crownpilot-api:ci
```

`smoke:container` deve iniciar a imagem, aguardar readiness/liveness, falhar
fechado em resposta inválida e remover o container mesmo em erro. Não usar
processo residual do host como substituto do container validado.

Validar que qualquer falha bloqueia merge/release e que nenhuma chamada live ou
secret de Production ocorre em PR/Preview.

## Definição de pronto

- CI executa gates .NET/frontend em PR; staging candidate publica digest OCI e
  `main` promove o mesmo digest sem rebuild;
- persistence/RLS/contract tests usam ambientes descartáveis/fixtures;
- imagem ASP.NET Core constrói e health smoke passa sem secrets;
- testes de authentication/authorization e redaction são obrigatórios;
- checks de observabilidade e redaction da Task 002-10 são obrigatórios;
- architecture tests impedem dependências invertidas e SDKs em Domain/Application;
- métricas distinguem resultado, latência e ambiente sem dados sensíveis;
- E2E/smoke possuem owner e pré-condição de Staging documentados.

## Riscos e cuidados

- Não marcar CI verde omitindo migrations, RLS ou build frontend.
- Não usar secret real em PR/Preview.
- Não declarar provider externo saudável por teste sem rede.
- Não logar exception HTTP completa quando contiver URL/payload sensível.
- Não transformar observabilidade em retenção de API data.
