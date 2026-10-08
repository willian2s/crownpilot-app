# 002-05 — Implementar port e adapter de lookup

- **Ticker:** `002`
- **Número:** `05`
- **Status:** `pending`

Esta task permanece bloqueada até ADR 005 ser aprovada e `002-13` -> `002-14` ->
`002-15` concluírem gates verdes. Referências .NET abaixo descrevem baseline
histórico; implementação deve seguir o ajuste Go de `002-15`.

## Objetivo e resultado esperado

Criar boundary server-side substituível para validar existência de perfil público
por Player Tag, sem acoplar Domain/Application ao host, proxy ou payload da API.

## Requisitos cobertos

- port de lookup em `internal/playerlink`, sem host ou provider no módulo;
- adapter HTTP em `internal/platform/clashroyale` usando `http.Client`;
- token nunca exposto ao browser;
- normalização e encoding determinísticos;
- estados externos distintos e ProblemDetails posterior;
- códigos de erro estáveis e `Retry-After` sem vazar provider;
- timeout/retry conservador e fixtures sem rede na CI.

## Escopo incluído

- operação `ResolvePublicProfileAsync` com cancellation;
- normalização e `%23` somente no path HTTP;
- host configurado no servidor, nunca input do usuário;
- mapeamento de `200`, `404`, `403`, `429`, `5xx`, timeout e misconfiguration;
- `Retry-After` e backoff limitado quando aplicável;
- redaction de tag, token, URL e payload em logs;
- fixtures sanitizadas e testes de contrato.

## Escopo excluído

- persistência de resposta, cache, snapshot ou catálogo;
- sync de coleção/Arena/battle log;
- escolha definitiva de proxy/egress;
- polling, retry infinito ou expansão de perfis;
- dependência de authentication para executar o adapter.

## Dependências

- `002-01-bootstrap-toolchain.md` e
  `002-02-estabelecer-boundaries-contrato-base-e-ambientes.md`;
- evidências e constraints da Fase 001;
- gates verdes de `002-13`, `002-14` e `002-15`, que substituem o baseline .NET
  antes da implementação no backend Go;
- não depende de authentication, mas não inicia antes do cutover Go.

## Arquivos e símbolos prováveis

- `internal/playerlink/` e resultado de `ResolvePublicProfile`;
- `internal/platform/clashroyale/`;
- `internal/httpapi/` para mapeamento posterior;
- fixtures de contract test e opções server-only de provider.

## Passos de implementação

1. Definir resultado discriminado sem snapshot ou dados extras.
2. Validar input antes de montar request e rejeitar host do usuário.
3. Construir path encoded e headers server-side.
4. Mapear status/timeout sem expor detalhes do provider.
5. Aplicar retry somente em casos permitidos, com teto e cancellation.
6. Instrumentar duração/status categorizado sem PII.
7. Criar fixtures para resposta resolvida, incompleta e falhas.

## Testes e comandos de validação

```text
go test ./internal/playerlink/... ./internal/platform/clashroyale/...
npm run typecheck --prefix frontend
npm run lint --prefix frontend
```

Confirmar que CI não faz request live e que bundle client não contém adapter,
token, host real ou URL com tag.

## Definição de pronto

- somente Infrastructure chama provider;
- `ResolvePublicProfileAsync` retorna todos estados previstos;
- `404` não sofre retry automático;
- `429` respeita `Retry-After` ou teto definido;
- mapeamento expõe códigos como `player_not_found`, `provider_unavailable` e
  `provider_rate_limited`, sem URL ou payload externo;
- resposta externa não é raw para o browser nem persistida;
- fixtures cobrem status, timeout e payload incompleto;
- logs são redacted e testes passam.

## Riscos e cuidados

- Proxy observado não é contrato oficial nem autorização de uso de dados.
- Não afirmar ownership após `200`.
- Não tratar `Cache-Control` como freshness ou quota.
- Não colocar tag real em URL de log, métrica ou erro.
