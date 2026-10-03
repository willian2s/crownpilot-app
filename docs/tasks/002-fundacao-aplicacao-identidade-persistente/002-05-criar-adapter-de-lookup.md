# 002-05 — Criar adapter de lookup

- **Ticker:** `002`
- **Número:** `05`
- **Status:** `pending`

## Objetivo e resultado esperado

Criar boundary server-side substituível para validar a existência de perfil
público por Player Tag, sem acoplar domínio ao proxy, host ou payload da API.

## Requisitos cobertos

- `ClashRoyaleClient` server-side;
- token nunca exposto ao browser;
- normalização e encoding de tag;
- estados externos distintos;
- timeout/retry conservador;
- fixtures sem chamadas live na CI.

## Escopo incluído

- port `resolvePublicProfile(playerTag)`;
- normalização determinística e `%23` somente no path HTTP;
- implementação HTTP com host configurado no servidor;
- mapeamento de `200`, `404`, `403`, `429`, `5xx`, timeout e configuração
  ausente;
- `Retry-After` quando presente, backoff limitado quando aplicável;
- redaction de tag, token, URL e payload em logs;
- fixtures sanitizadas de contrato.

## Escopo excluído

- persistência de resposta, cache, snapshot ou catálogo;
- sync de coleção/Arena/battle log;
- escolha definitiva de proxy/egress;
- polling, retry infinito ou expansão de perfis.

## Dependências

- `002-01`, `002-02` e `002-04` para runtime/config/auth;
- evidências de API da Fase 001;
- segredo próprio de staging, se smoke externo for explicitamente autorizado.

## Arquivos e símbolos prováveis

- `app/Contracts/ClashRoyaleClient.php`;
- `app/Adapters/ClashRoyale/HttpClient.php`;
- `NormalizedPlayerTag`, `ResolvePublicProfileResult`, `ProviderError`;
- `tests/Contract/fixtures/`;
- configuração server-only de `CLASH_ROYALE_API_TOKEN` e provider host.

## Passos de implementação

1. Definir resultado discriminado sem campos de snapshot.
2. Validar input antes de montar request e rejeitar host fornecido pelo usuário.
3. Construir path com tag encoded e headers server-side.
4. Mapear status/timeout sem expor detalhes do provider.
5. Implementar retry somente para casos permitidos, com teto e jitter.
6. Instrumentar duração/status categorizado sem PII.
7. Criar fixtures para resposta resolvida e falhas.

## Testes e comandos de validação

```text
composer run test:unit
npm run test:contract
npm run typecheck
npm run lint
```

Confirmar que CI não faz request live e que bundle client não contém adapter,
token, credencial Supabase/PHP, host real ou URL com tag.

## Definição de pronto

- somente server-side chama provider;
- `resolvePublicProfile` retorna todos os estados previstos;
- `404` não sofre retry automático;
- `429` respeita `Retry-After` ou backoff com teto;
- resposta externa não é retornada como raw ao browser nem persistida;
- fixtures cobrem status e payload incompleto;
- logs são redacted e testes passam.

## Riscos e cuidados

- Proxy observado não é contrato oficial nem autorização de uso de dados.
- Não afirmar ownership após `200`.
- Não tratar `Cache-Control` como freshness ou quota.
- Não colocar tag real em URL de log, métrica ou erro.
