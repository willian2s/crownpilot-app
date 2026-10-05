# 002-10 — Instrumentar observabilidade, health e redaction

- **Ticker:** `002`
- **Número:** `10`
- **Status:** `pending`

## Objetivo e resultado esperado

Instrumentar diagnóstico mínimo da API, autenticação, vínculo e provider sem
registrar tokens, secrets ou dados sensíveis, com health liveness/readiness,
correlation/request ID e métricas de baixa cardinalidade.

## Requisitos cobertos

- structured logging;
- liveness sem dependências externas;
- readiness com configuração e PostgreSQL conforme ambiente;
- correlation/request ID via `Activity`/primitives .NET;
- métricas básicas de request, latência, status e provider;
- redaction testável;
- ausência de payload externo completo em logs.

## Escopo incluído

- endpoints `/health/live` e `/health/ready`;
- middleware/diagnostics para request ID e correlation ID;
- categorias de resultado, latência, status e ambiente;
- redaction de Firebase ID/refresh tokens, passwords, service accounts, DB
  passwords, UID, Player Tag, URL com tag, payload externo e secrets;
- readiness fail-closed para configuração inválida, sem lookup Clash Royale;
- testes e documentação operacional dos sinais produzidos.

## Escopo excluído

- dashboards avançados, tracing/exporter obrigatório ou OpenTelemetry prematuro;
- analytics de produto;
- logs de payload para depuração;
- chamadas live automáticas ao provider.

## Dependências

- `002-01-bootstrap-toolchain.md`,
  `002-02-estabelecer-boundaries-contrato-base-e-ambientes.md`,
  `002-04-implementar-google-sign-in-e-firebase-bearer.md`,
  `002-05-implementar-port-e-adapter-de-lookup.md`,
  `002-07-implementar-casos-de-uso-e-api-v1.md` e frontend
  `002-08-entregar-frontend-de-identidade-e-vinculo.md`;
- configuração de banco da `002-03-preparar-postgresql-migrations-e-harness-rls.md`.

## Passos de implementação

1. Implementar liveness e readiness com contratos distintos.
2. Adicionar request/correlation ID e structured logging.
3. Instrumentar resultados e latências de auth, link, API e provider sem PII.
4. Aplicar redaction antes de serializar logs ou exceptions.
5. Testar ausência de tokens, secrets, UID, tag, URL e payload em logs/bundles.

## Testes e comandos de validação

```text
dotnet test --filter Category=Observability
dotnet test --filter Category=Api
npm run build
```

Verificar liveness sem banco/provider, readiness conforme configuração, IDs
estáveis, métricas de baixa cardinalidade e respostas sem stack trace/secrets.

## Definição de pronto

- liveness responde quando dependências estão indisponíveis;
- readiness falha de modo explícito quando configuração/banco não estão prontos;
- logs são estruturados, correlacionáveis e redacted;
- métricas não carregam tag, UID, token ou payload;
- testes de redaction e health passam;
- observabilidade não cria retenção de dados Clash Royale.

## Riscos e cuidados

- Não fazer readiness chamar provider externo.
- Não usar correlation ID como identidade de usuário.
- Não registrar exception HTTP completa sem sanitização.
