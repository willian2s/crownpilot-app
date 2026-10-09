# 002-11 — Automatizar CI, OCI e gates de release

- **Ticker:** `002`
- **Número:** `11`
- **Status:** `pending`

Esta task permanece bloqueada até ADR 005 ser aprovada e `002-13` -> `002-14` ->
`002-15` concluírem gates verdes. Referências .NET abaixo descrevem baseline
histórico; implementação deve seguir Go, `depguard` e o ajuste de `002-15`.

## Objetivo e resultado esperado

Automatizar gates de Pull Request/main e preparar candidate Staging com imagem
OCI imutável, sem secrets, mantendo observabilidade já definida na Task
`002-10-instrumentar-observabilidade-health-e-redaction.md`.

## Requisitos cobertos

- CI Go com `go mod download`, build, `golangci-lint`, `go test` e `go test -race`;
- CI frontend com npm ci, lint, type-check, testes e build;
- testes de persistência PostgreSQL, migrations `goose` e RLS;
- contract tests do lookup sem chamadas live;
- build/smoke da imagem Go sem secrets;
- E2E/smoke com comandos e pré-condições claros;
- validação do documento OpenAPI gerado e drift de contrato;
- checks de observabilidade, health e redaction da Task
  `002-10-instrumentar-observabilidade-health-e-redaction.md`;
- gates separados para PR, main e Staging.

## Escopo incluído

- workflow `.github/workflows/ci.yml` para `pull_request`, staging candidate e
  promoção em `main`;
- `container-smoke.mjs` deve consumir `CROWNPILOT_CONTAINER_IMAGE` quando
  fornecida, sem reconstruir outra imagem, e validar `/health/live` e
  `/health/ready`;
- geração do bundle/job final de migrations a partir do schema atual, separado do
  container de runtime;
- jobs Go, frontend, persistence/RLS, contract, image e health;
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

- `002-05` a `002-10` e `002-13` a `002-15` para fechamento;
  gates mínimos começam em `002-01-bootstrap-toolchain.md`;
- Docker/PostgreSQL e fixtures disponíveis em CI;
- secrets reais somente em ambientes controlados de Staging.

## Arquivos e símbolos prováveis

- `.github/workflows/ci.yml`;
- packages `internal/`, contract tests, persistence/RLS e métricas;
- `internal/observability/` e scripts `go`, `npm`, `database` e documentação de release.

## Passos de implementação

1. Definir jobs e falhas para PR e push em `main`.
2. Executar restore/build/test Go e frontend sem secrets reais.
3. Subir PostgreSQL descartável, aplicar migrations `goose`, RLS e fixtures; executar
   checks de pool/reset e isolamento A/B.
4. Executar testes de token Firebase com Emulator/fixtures assinadas.
5. Executar `depguard` para provar que `cmd` compõe, `httpapi` traduz transporte,
   módulos não importam providers/HTTP e `identity` não importa `playerlink`;
   usar `go list` somente para lacunas de composição/transitividade.
6. Construir imagem Docker e testar `/health/live`/`ready` sem secret na imagem.
7. Gerar e versionar como artefato de candidate o bundle/job final de migrations;
   não reutilizar bundle produzido antes da Task
   `002-06-modelar-persistencia-repositories-e-rls.md`.
8. Criar workflow protegido/manual de Staging com owner, aprovação, migration
   job, mesmo digest OCI e smoke antes de promoção para `main`.
9. Documentar gates Staging/Production fora dos gates de PR/main.

## Testes e comandos de validação

```text
go mod download
go build ./cmd/crownpilot-api
go test ./...
go test -race ./...
golangci-lint run
npm ci --prefix frontend
npm run lint --prefix frontend
npm run typecheck --prefix frontend
npm run test:unit --prefix frontend
npm run test:contract --prefix frontend
npm run build --prefix frontend
go test ./internal/platform/postgres/...
npm run test:rls --prefix frontend
npm run openapi:check --prefix frontend
docker build -t crownpilot-api:ci .
CROWNPILOT_CONTAINER_IMAGE=crownpilot-api:ci npm run smoke:container --prefix frontend
```

`smoke:container` deve iniciar a imagem, aguardar readiness/liveness, falhar
fechado em resposta inválida e remover o container mesmo em erro. Não usar
processo residual do host como substituto do container validado.

Validar que qualquer falha bloqueia merge/release e que nenhuma chamada live ou
secret de Production ocorre em PR/Preview.

## Definição de pronto

- CI executa gates Go/frontend em PR; staging candidate publica digest OCI e
  `main` promove o mesmo digest sem rebuild;
- persistence/RLS/contract tests usam ambientes descartáveis/fixtures;
- imagem Go constrói e health smoke passa sem secrets;
- testes de authentication/authorization e redaction são obrigatórios;
- checks de observabilidade e redaction da Task
  `002-10-instrumentar-observabilidade-health-e-redaction.md` são obrigatórios;
- `depguard` impede dependências invertidas e SDKs em módulos de domínio;
- métricas distinguem resultado, latência e ambiente sem dados sensíveis;
- E2E/smoke possuem owner e pré-condição de Staging documentados.

## Riscos e cuidados

- Não marcar CI verde omitindo migrations, RLS ou build frontend.
- Não usar secret real em PR/Preview.
- Não declarar provider externo saudável por teste sem rede.
- Não logar exception HTTP completa quando contiver URL/payload sensível.
- Não transformar observabilidade em retenção de API data.
