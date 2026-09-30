# 001-05 — Medir operação, cache e rate limits

- **Ticker:** `001`
- **Número:** `05`
- **Status:** `pending`

## Requisitos cobertos

- cache, latência, erros, `Retry-After`, rate behavior e requisitos de token/IP;
- custo de requests para sync de jogador, histórico, catálogo e meta;
- impacto operacional em 1k, 10k e 100k usuários sem executar carga real.

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
- validar compatibilidade do modelo de egress com Vercel;
- registrar se Static IPs da Vercel são necessários/adequados ou se um egress
  gateway substituível será necessário;
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
- troca da decisão de Vercel como deploy inicial;
- implementação definitiva de gateway/worker.

## Dependências

- 001-01 para auth, endpoints e headers básicos;
- 001-02 e 001-03 para custo de perfil e histórico;
- 001-04 para projetar custo de aquisição do meta.

## Arquivos e símbolos prováveis

- `evidences/operational-findings.md`;
- headers `Cache-Control` e `Retry-After`, códigos HTTP e variável
  `CLASH_ROYALE_API_TOKEN`;
- operações conceituais `syncPlayer`, `ingestMetaBatch` e `refreshCatalog`;
  nenhum job ou adapter existe na `main`.

## Passos de execução

1. Coletar headers de endpoints já usados nas tasks anteriores.
2. Registrar `Cache-Control`, `Retry-After` quando existir e códigos.
3. Separar requests user-driven de background ingestion.
4. Definir fórmula simples de requests por sync.
5. Projetar três cenários de escala.
6. Identificar quais resultados podem ser compartilhados/cacheados globalmente.
7. Verificar o impacto de IP allowlist no deploy inicial em Vercel.
8. Se egress estático for necessário, comparar:
   - Vercel Static IPs;
   - egress gateway externo/substituível.
9. Registrar requisitos que futura infraestrutura precisa satisfazer sem
   acoplar o domínio ao provider.

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

## Testes e comandos de validação

- reutilizar respostas das tasks anteriores, sem provocar `429` deliberadamente;
- registrar headers e status de probes seguros em evidência sanitizada;
- calcular requests por operação e projetar os três cenários de escala,
  marcando-os como estimativas;
- revisar `Retry-After`, separação user-driven/background e `git diff --check`.

## Riscos e cuidados

- IP-bound token pode exigir Vercel Static IPs ou egress gateway dedicado;
- dependência direta de networking proprietário da Vercel criaria lock-in;
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
