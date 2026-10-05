# 002-08 — Automatizar quality gates e observabilidade

- **Ticker:** `002`
- **Número:** `08`
- **Status:** `pending`

## Objetivo e resultado esperado

Automatizar gates de Pull Request/main e instrumentar o mínimo operacional para
diagnosticar authentication, vínculo e provider sem registrar dados sensíveis.

## Requisitos cobertos

- CI .NET com restore, build, analyzers e `dotnet test`;
- CI frontend com npm ci, lint, type-check, testes e build;
- testes de persistência PostgreSQL, EF migrations e RLS;
- contract tests do lookup sem chamadas live;
- build/smoke da imagem ASP.NET Core sem secrets;
- E2E/smoke com comandos e pré-condições claros;
- logs/métricas com redaction e correlation ID;
- gates separados para PR, main e Staging.

## Escopo incluído

- workflow `.github/workflows/ci.yml` para `pull_request` e push em `main`;
- jobs .NET, frontend, persistence/RLS, contract, image e health;
- PostgreSQL local/descartável e Firebase Emulator/fixtures sem produção;
- scripts de E2E e smoke staging-only;
- logger/métricas com request ID, ambiente, resultado e latência;
- scans para secrets, URLs com tag e imports server-only no bundle;
- documentação de owner/pré-condição dos gates Staging/Production.

## Escopo excluído

- provisionamento de Production pelo CI;
- polling/probes live automáticos;
- dashboards avançados ou analytics de produto;
- logging de payloads para depuração;
- chamadas live à Clash Royale API na CI normal.

## Dependências

- `002-01` a `002-07` para fechamento; gates mínimos começam em `002-01`;
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
3. Subir PostgreSQL descartável, aplicar EF migrations, RLS e fixtures.
4. Executar testes de token Firebase com Emulator/fixtures assinadas.
5. Adicionar métricas de authentication, link, erro, latência e ambiente.
6. Aplicar redaction antes de serializar logs ou exceptions.
7. Construir imagem Docker e testar `/health/live`/`ready` sem secret na imagem.
8. Criar workflow protegido/manual de Staging com owner, aprovação, migration
   job e smoke antes de promoção para `main`.
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
npm run test:rls
npm run build
docker build -t crownpilot-api:ci .
curl --fail http://localhost:8080/health/live
```

Validar que qualquer falha bloqueia merge/release e que nenhuma chamada live ou
secret de Production ocorre em PR/Preview.

## Definição de pronto

- CI executa gates .NET/frontend em PR e em cada push para `main`;
- persistence/RLS/contract tests usam ambientes descartáveis/fixtures;
- imagem ASP.NET Core constrói e health smoke passa sem secrets;
- testes de authentication/authorization e redaction são obrigatórios;
- logs não contêm tag, e-mail, UID, token, IP, URL real ou payload;
- métricas distinguem resultado, latência e ambiente;
- E2E/smoke possuem owner e pré-condição de Staging documentados.

## Riscos e cuidados

- Não marcar CI verde omitindo migrations, RLS ou build frontend.
- Não usar secret real em PR/Preview.
- Não declarar provider externo saudável por teste sem rede.
- Não logar exception HTTP completa quando contiver URL/payload sensível.
- Não transformar observabilidade em retenção de API data.
