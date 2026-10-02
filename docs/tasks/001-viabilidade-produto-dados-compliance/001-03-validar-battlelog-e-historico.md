# 001-03 — Validar battle log e histórico

- **Ticker:** `001`
- **Número:** `03`
- **Status:** `completed with constraints`

## Requisitos cobertos

- janela, paginação/cursor, ordenação e modos do battle log;
- campos de batalha, decks, resultado, contexto competitivo e Evo/Hero;
- viabilidade de snapshots, polling, deduplicação e limites para histórico.

## Objetivo e resultado esperado

Determinar o que o battle log oficial permite reconstruir e qual estratégia
seria necessária para histórico pós-partida confiável.

## Escopo incluído

- medir quantidade de batalhas retornadas;
- confirmar se existe paginação/cursor;
- mapear modos presentes;
- identificar timestamp, resultado, decks, tower/support cards, Evo/Hero usado,
  trophies/trophy change e opponent quando disponíveis;
- verificar ordenação;
- repetir chamada em momentos diferentes se necessário para observar janela;
- estimar polling necessário para evitar gaps em jogadores muito ativos;
- avaliar deduplicação/idempotência;
- produzir `evidences/battlelog-findings.md`.

## Escopo excluído

- replay;
- posição/timing de cartas;
- tracking in-match;
- polling de produção;
- armazenamento definitivo.

## Dependências

- 001-01;
- semântica de cards iniciada em 001-02.

## Arquivos e símbolos prováveis

- `evidences/battlelog-findings.md`;
- endpoint `/v1/players/{tag}/battlelog` e campos de batalha observados;
- chave candidata de deduplicação e contrato conceitual de histórico; não há
  pipeline ou persistência implementados na `main`.

## Passos de execução

1. Consultar endpoint para mais de um perfil.
2. Registrar tamanho e response shape.
3. Confirmar presença/ausência de paginação.
4. Catalogar tipos/modes observados.
5. Definir chave candidata para deduplicar batalha sem assumir ID inexistente.
6. Verificar se decks e formas Evo/Hero são observáveis por participante.
7. Estimar cenário de perda de histórico com polling 1m/5m/15m/1h.
8. Registrar quais análises futuras são possíveis e quais não são.

## Perguntas de saída

- Podemos reconstruir histórico próprio por snapshots?
- Existe risco real de perder partidas entre syncs?
- Quais modos devem entrar no MVP e quais precisam ser filtrados?
- Battle log ajuda a descobrir candidatos de meta?
- Há informação suficiente para futuro matchup analysis?

## Validação

A evidência precisa distinguir:

- documentado;
- observado;
- inferido;
- não disponível.

## Definição de pronto

Existe um contrato operacional honesto para battle history e não prometemos
telemetria que a API não fornece.

## Testes e comandos de validação

- repetir o endpoint em mais de um perfil e em horários distintos, respeitando
  rate limit;
- comparar tamanho, ordenação e shape das respostas sem usar payload pessoal
  bruto como fixture versionada;
- simular por cálculo os cadences 1m/5m/15m/1h, identificando estimativas como
  tal;
- revisar `git diff --check` e a classificação documentado/observado/inferido.

## Riscos e cuidados

- battle log curto pode exigir polling mais frequente;
- modos especiais podem ter shapes diferentes;
- listas que não têm oito cartas não devem ser tratadas automaticamente como
  deck;
- evitar coletar histórico em escala nesta fase.

## Registro de execução

### Execução em `2026-10-02`

- **Status final:** `completed with constraints`.
- **Arquivos alterados:**
  `evidences/battlelog-findings.md`; este registro; overview para marcar somente
  `001-03` e recalcular progresso.
- **Decisões:** tratar battle log como lista observada de até 30 entradas, sem
  assumir paginação ou ID estável; separar `PvP` e `pathOfLegend`; ingerir
  `trail`/`unknown` com quarentena; usar fingerprint opaco para idempotência; não
  prometer histórico completo nem Hero explícito.
- **Janela observada:** três perfis retornaram 6, 30 e 30 entradas. Amostra de
  30 entradas cobriu aproximadamente 96 minutos em perfil ativo e 6–7 dias em
  perfis de menor atividade. Todas as 66 entradas vieram ordenadas por
  `battleTime` descendente.
- **Modos observados:** `PvP` (6), `pathOfLegend` (30), `trail` (29) e
  `unknown` (1). Semântica oficial de `trail`/`unknown` permanece pendente.
- **Campos úteis:** `battleTime`; `type`; `gameMode`; `arena`; `team`;
  `opponent`; oito `cards[]` observados por participante; `supportCards[]`;
  `crowns`; `trophyChange`; `startingTrophies`; dados de torres; campos opcionais
  de Evolution; `eventTag` não universal. Não houve campo explícito de resultado,
  winner, battle ID ou Hero deployment.
- **Comandos executados:** probes `curl` server-side via
  `proxy.royaleapi.dev` para três tags locais sanitizadas, com token carregado
  somente de `.env.local`; resumo estrutural Python em diretório temporário;
  repetição do perfil ativo após `65s`; cálculo Python das cadences `1m/5m/15m/1h`;
  `git diff --check`; validação estrutural Python do ticker, subtarefas e checklist;
  confirmação de ausência de toolchain de aplicação.
- **Resultados/evidências:** cinco requests retornaram `200`; resposta foi array
  JSON; amostra repetida manteve 30/30 fingerprints; `Cache-Control` observado
  entre `max-age=21` e `60`; fingerprint candidato foi único em todas as
  entradas observadas. Estimativa baseada em 29 intervalos/96 minutos indica
  aproximadamente 0,30, 1,51, 4,54 e 18,14 batalhas esperadas por sync em
  `1m/5m/15m/1h`; risco modelado de exceder 30 em 1h foi ~0,37%, sem garantia.
  Evidência completa em `evidences/battlelog-findings.md`.
- **Desvios:** rota oficial direta continua bloqueada por allowlist; observações
  são via proxy operacional e não contrato oficial. Não houve espera para medir
  eviction real nem coleta de produção; repetição de 65s não observou nova
  batalha.
- **Limitações:** paginação/cursor, limite rígido, semântica oficial de modos,
  resultado explícito, ID estável, Hero e contrato direto permanecem
  `unavailable`/`unresolved`. Histórico próprio só é best-effort por polling.
- **Riscos residuais:** gaps durante outage ou cadência insuficiente; bursts acima
  da amostra; retenção/privacy e termos do proxy; uso e retenção de dados de
  oponentes; mudança de shape/modes; risco de interpretar `iconUrls.heroMedium`
  como Hero usado.
- **Revisão independente:** solicitada após implementação e após correções de
  rastreabilidade, classificação de ausência e cálculo de polling; revisão aprovou
  subtarefa sem blockers. Nenhuma subtarefa seguinte foi iniciada.
