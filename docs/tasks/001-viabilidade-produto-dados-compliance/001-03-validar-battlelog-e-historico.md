# 001-03 — Validar battle log e histórico

- **Ticker:** `001`
- **Número:** `03`
- **Status:** `planned`

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

## Riscos e cuidados

- battle log curto pode exigir polling mais frequente;
- modos especiais podem ter shapes diferentes;
- listas que não têm oito cartas não devem ser tratadas automaticamente como
  deck;
- evitar coletar histórico em escala nesta fase.

## Registro de execução

- **Status final:**
- **Janela observada:**
- **Modos observados:**
- **Campos úteis:**
- **Limitações:**
- **Riscos residuais:**
