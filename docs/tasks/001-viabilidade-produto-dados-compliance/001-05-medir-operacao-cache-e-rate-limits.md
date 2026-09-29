# 001-05 — Medir operação, cache e rate limits

- **Ticker:** `001`
- **Número:** `05`
- **Status:** `planned`

## Objetivo e resultado esperado

Transformar a estratégia de dados em ordem de grandeza operacional antes de
escolher infraestrutura.

## Escopo incluído

- registrar headers de cache relevantes;
- observar latência básica;
- confirmar comportamento de 403/404/429/5xx quando ocorrer de forma natural ou
  com probes seguros;
- honrar `Retry-After`;
- verificar restrições de token/IP;
- estimar requests de:
  - player sync;
  - battle history;
  - card catalog;
  - meta ingestion;
- projetar cenários de 1k, 10k e 100k usuários ativos sem executar carga real;
- produzir `evidences/operational-findings.md`.

## Escopo excluído

- load test;
- tentativa de descobrir limite por saturação;
- bypass de quota;
- rotação de IP/token para aumentar throughput;
- decisão de cloud/provider.

## Passos de execução

1. Coletar headers de endpoints já usados nas tasks anteriores.
2. Registrar `Cache-Control`, `Retry-After` quando existir e códigos.
3. Separar requests user-driven de background ingestion.
4. Definir fórmula simples de requests por sync.
5. Projetar três cenários de escala.
6. Identificar quais resultados podem ser compartilhados/cacheados globalmente.
7. Registrar requisitos que futura infraestrutura precisa satisfazer.

## Princípio de segurança operacional

Não provocar 429 deliberadamente para "medir o limite".

Se o limite oficial não for publicado, registrar como desconhecido e trabalhar
com observação conservadora.

## Saída mínima

`evidences/operational-findings.md` deve mostrar:

- request path;
- cache observado/documentado;
- custo por operação;
- cadence candidata;
- failure behavior;
- impacto em 1k/10k/100k;
- mitigação.

## Definição de pronto

A Fase 002 conhece as restrições de integração sem estar presa a um deploy
específico.

## Riscos e cuidados

- IP-bound token pode excluir alguns modelos serverless;
- cache do provedor pode mudar;
- custo de meta ingestion pode dominar custo de user sync;
- números projetados devem ser marcados como estimativas, não SLOs.

## Registro de execução

- **Status final:**
- **Cache observado:**
- **Rate behavior:**
- **Request model:**
- **Infra constraints:**
- **Riscos residuais:**
