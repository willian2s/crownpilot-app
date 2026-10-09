# CrownPilot — Resumo de contexto do projeto

**Data:** 29 de setembro de 2026  
**Origem:** conversa de descoberta e planejamento inicial do produto

Este documento registra o contexto que levou ao desenho atual do CrownPilot.

Ele é um **resumo histórico de produto**. Para intenção futura, precedem:
`docs/decisions/`, `docs/specs/`, `docs/roadmap/` e `docs/tasks/`. Para estado já
existente, a implementação em `main` é a fonte de verdade. Divergências exigem
investigação; planejamento não deve fingir que entrega foi implementada.

---

# 1. Origem da ideia

A conversa começou com a busca por um companion para Clash Royale que oferecesse uma experiência semelhante ao SnapComplete no Marvel Snap.

O requisito não era overlay durante a partida.

A experiência desejada era:

- possuir uma conta própria no companion;
- vincular a Player Tag uma única vez;
- acessar o produto em outro dispositivo sem precisar lembrar/copiar a tag;
- sincronizar coleção e progressão;
- conhecer níveis, Evolutions e Heroes disponíveis;
- cruzar a conta com decks competitivos;
- mostrar decks que o jogador realmente consegue montar;
- recomendar progressão e upgrades;
- no futuro, oferecer análise/coaching contextualizado.

A principal lacuna percebida nos produtos existentes foi a fragmentação: diferentes ferramentas resolvem partes do problema, mas poucas tratam **a conta do jogador como centro permanente da experiência**.

---

# 2. Proposta central do CrownPilot

A pergunta que define o produto é:

> **Qual é o melhor deck que eu consigo jogar agora, na minha Arena, com a conta que eu realmente tenho?**

A proposta não é simplesmente listar os melhores decks globais.

O CrownPilot deve considerar, conforme os dados permitirem:

- Arena;
- faixa de troféus;
- coleção;
- níveis das cartas;
- Evolutions;
- Heroes;
- decks reais observados;
- força do deck no contexto competitivo relevante;
- confiança e tamanho da amostra;
- futuramente, desempenho histórico do próprio jogador.

A frase que resume o posicionamento é:

> **Best decks for you.**

A intenção é que dois jogadores na mesma Arena possam receber recomendações diferentes porque possuem contas diferentes.

---

# 3. Nome do projeto

O nome escolhido foi:

> **CrownPilot**

Tagline inicial:

> **Play the right deck. Upgrade the right cards.**

Repositório:

`willian2s/crownpilot-app`

---

# 4. Modelo de produto

O produto foi organizado em quatro camadas.

## 4.1 Discover

Responder:

> **O que eu deveria jogar agora?**

Feature principal:

### Best Decks for You

O sistema parte de decks reais observados e os ranqueia para a conta do jogador.

A primeira versão não deve pedir a um LLM para inventar oito cartas.

---

## 4.2 Progress

Responder:

> **O que eu deveria melhorar agora?**

Feature principal:

### Upgrade Planner

Modelo de experiência:

- **NOW** — decks que o jogador consegue usar hoje;
- **NEXT** — decks que estão a um ou poucos upgrades de ficar prontos;
- **TARGET** — decks relevantes que fazem sentido perseguir como objetivo.

A ideia mais importante é medir **Upgrade Impact**.

Exemplo conceitual:

- melhorar Fisherman pode liberar quatro decks competitivos;
- melhorar outra carta pode afetar apenas um deck pouco relevante.

O produto deve conseguir recomendar o upgrade com maior impacto na conta, não simplesmente a carta mais popular.

---

## 4.3 Optimize

Responder:

> **Como posso melhorar o deck que eu já jogo?**

A primeira versão de otimização deve buscar **variantes reais/comprovadas**, não gerar substituições arbitrárias.

Exemplo:

```text
Current
A B C D E F G H

Recommended variant
A B C D E F X H
```

A recomendação deve considerar:

- evidência da variante;
- nível da carta alternativa;
- matchups;
- contexto competitivo;
- custo de oportunidade.

---

## 4.4 Coach

Responder perguntas em linguagem natural sobre os resultados calculados pelo CrownPilot.

Exemplos:

- Por que este deck é melhor para mim?
- Por que devo melhorar esta carta?
- Por que estou perdendo para determinado arquétipo?
- O que mudou nas minhas últimas partidas?
- Qual deve ser meu próximo objetivo?

A IA não é o core do produto.

---

# 5. Princípio central de arquitetura

Foi definido o princípio:

> **Deterministic by default. AI when it adds value.**

O recommendation engine, readiness, upgrade impact, métricas e rankings devem ser determinísticos sempre que possível.

Fluxo desejado:

```text
Player Data
    |
    v
Recommendation / Analytics Engine
    |
    v
Structured Insight
    |
    +--> UI
    |
    +--> AI Coach
```

A IA recebe resultados já estruturados e serve para explicar, sintetizar e conversar sobre eles.

Isso traz:

- menor custo;
- comportamento previsível;
- testabilidade;
- facilidade para trocar provider de IA;
- continuidade do produto mesmo sem LLM.

---

# 6. Como os decks devem ser escolhidos

A primeira versão não deve criar decks do zero.

Pipeline conceitual:

```text
partidas/decks reais
        |
        v
canonicalização
        |
        v
meta por contexto
        |
        v
compatibilidade com a conta
        |
        v
ranking personalizado
```

Componentes conceituais do futuro `PlayerDeckFit`:

- Meta Strength;
- Collection Compatibility;
- Level Readiness;
- Arena Relevance;
- Confidence.

Posteriormente pode entrar:

- Personal Performance / Playstyle Fit.

A fórmula final ainda não foi definida e dependerá dos dados reais.

---

# 7. Arena e faixa competitiva

Uma contribuição importante para a ideia foi perceber que Arena não deve ser apenas um filtro de cartas desbloqueadas.

O objetivo é usar o contexto competitivo do jogador.

Em vez de:

```text
meta global
-> filtrar pela coleção
```

o objetivo é chegar mais perto de:

```text
meta relevante para a faixa do jogador
+ coleção
+ níveis
+ Evolutions/Heroes
-> Best Decks for You
```

Entretanto, isso virou um **gate técnico da Fase 001**.

Rankings oficiais podem representar principalmente o topo da ladder e não necessariamente a mid-ladder.

Portanto, o produto só deve alegar "meta da sua Arena" se houver estratégia de aquisição de dados com cobertura e confiança suficientes.

Caso contrário, deve:

- usar outra segmentação, como trophy bands;
- consumir fonte licenciada;
- ou reduzir honestamente a promessa.

---

# 8. Identidade persistente

Um dos principais requisitos de UX é eliminar o problema de precisar copiar a Player Tag em cada dispositivo.

Fluxo desejado:

```text
Google Sign-In
      |
      v
CrownPilot User
      |
      v
Player Tag vinculada uma única vez
      |
      v
Conta Clash Royale conhecida em qualquer dispositivo
```

Não haverá armazenamento de senha ou credencial Supercell.

A Player Tag é um vínculo público do domínio.

---

# 9. Infraestrutura definida no planejamento original

> **Nota histórica:** esta seção registra decisões anteriores à ADR 004. Ela não
> é fonte de verdade para implementação da Fase 002. A stack vigente está em
> `docs/decisions/004-aspnet-core-react-vite-firebase-postgresql.md`.

Foi escolhida a seguinte infraestrutura-base:

## Firebase Authentication

- Google Sign-In como provider inicial;
- identidade da conta CrownPilot;
- independente da identidade Supercell.

## Cloud Firestore (substituído na Fase 002)

Banco principal da aplicação.

Deverá armazenar progressivamente, conforme as specs:

- usuário CrownPilot;
- Player Tag vinculada;
- snapshots normalizados;
- coleção;
- histórico permitido;
- dados derivados;
- recomendações;
- metadados de sync.

Firestore foi um vendor lock-in aceito no planejamento original, mas foi
substituído por PostgreSQL via EF Core + Npgsql, hospedado inicialmente no
Supabase. Não criar Firestore na Fase 002.

A modelagem de collections ainda não foi definida porque deve nascer dos contratos reais descobertos na Fase 001.

## Vercel (alternativa estática; hosting original substituído)

Era plataforma inicial de deploy/runtime. Na Fase 002, Render é o hosting inicial
da API Docker e do frontend estático preferencial; Vercel permanece alternativa
de frontend estático, nunca backend obrigatório.

A decisão é operacional, não um boundary arquitetural.

A aplicação deve ser migrável para outro runtime sem reescrever o domínio.

Evitar dependência obrigatória de:

- Vercel KV/Postgres;
- Blob;
- Edge Config;
- Queues/Workflow;
- Cron;
- APIs proprietárias do Edge/runtime.

Recursos Vercel podem ser usados no futuro quando estiverem isolados atrás de adapters substituíveis.

Decisão formal:

`docs/decisions/001-firebase-firestore-vercel-portable.md`

---

# 10. Integração com a API do Clash Royale

A integração deve ser isolada do domínio.

Modelo conceitual:

```text
Domain / Application
        |
        v
ClashRoyaleClient
        |
        v
HTTP Adapter
        |
        v
api.clashroyale.com
```

Motivos:

- permitir troca de runtime;
- isolar autenticação/token;
- controlar cache e retry;
- lidar com mudanças de API;
- resolver possível necessidade de egress/IP estático sem contaminar o domínio.

A Fase 001 deve validar os requisitos reais atuais da API.

---

# 11. Jobs e sincronização

Jobs devem ser modelados como operações da aplicação independentes do scheduler.

Exemplos:

```text
syncPlayer(...)
refreshCardCatalog(...)
ingestMetaBatch(...)
```

Vercel Cron pode futuramente disparar essas funções, mas não deve ser a implementação da regra.

Assim, o scheduler pode ser trocado sem alterar o domínio.

---

# 12. IA e monetização

A conclusão inicial foi que features com custo variável de IA precisam de controle econômico.

Entretanto, a arquitetura escolhida reduz bastante a dependência de IA.

## Core gratuito/determinístico

Pode incluir:

- conta persistente;
- sync da coleção;
- Best Decks for You;
- Arena/trophy recommendations;
- readiness;
- Upgrade Planner;
- estatísticas.

## Possível camada de AI Coach

Pode incluir:

- explicações aprofundadas;
- análise de histórico;
- perguntas sobre a conta;
- coaching;
- planos personalizados.

O desenho comercial inicialmente imaginado considerava:

- plano gratuito;
- plano pago;
- franquia/créditos de IA em vez de uso ilimitado.

Porém, isso **não está aprovado como modelo comercial**.

---

# 13. Gate de compliance e monetização

A Fan Content Policy da Supercell impõe restrições importantes.

O projeto não deve presumir que um SaaS tradicional com assinatura/paywall é permitido.

A Fase 001 possui uma task específica para determinar o estado de cada possibilidade:

- `ALLOWED`;
- `ALLOWED WITH CONDITIONS`;
- `REQUIRES EXPLICIT APPROVAL`;
- `NOT ALLOWED`;
- `UNRESOLVED`.

Pontos a validar:

- fan app gratuito;
- uso de assets;
- disclaimer;
- anúncios;
- donations;
- coaching;
- AI Coach;
- assinatura;
- armazenamento/redistribuição dos dados da API.

Até essa análise ser concluída:

> **billing permanece bloqueado.**

### Atualização documental em 1º de outubro de 2026

A leitura da Fan Content Policy oficial confirmou que Fan Content é não comercial
por padrão e que não é permitido cobrar taxa de qualquer tipo, inclusive por
funcionalidades in-app, sem aprovação expressa da Supercell. Ads, donations e
coaching são exceções condicionais, não autorização geral para assinatura,
premium features, software coaching ou AI Coach.

Para o CrownPilot, a ausência de aprovação de assinatura não gera NO-GO automático.
O produto pode permanecer tecnicamente viável com operação gratuita ou modelo
permitido sob condições. Se a sustentabilidade depender de modelo que exige
aprovação, o resultado deve ser `GO WITH CONSTRAINTS / APPROVAL DEPENDENCY`, com a
receita bloqueada até aprovação rastreável. Acesso à API, API key ou ownership
verification não equivalem a autorização comercial.

Fontes oficiais consultadas em `2026-10-01`:

- [Supercell Fan Content Policy](https://supercell.com/en/fan-content-policy/),
  `Last updated: September 27, 2023`;
- [Supercell Terms of Service](https://supercell.com/en/terms-of-service/),
  `Effective Date: November 6, 2024`;
- [Clash Royale API developer portal](https://developer.clashroyale.com/).

---

# 14. O que o CrownPilot não é

Não faz parte da proposta:

- overlay durante partida;
- elixir tracker;
- cycle tracker em tempo real;
- bot;
- mod;
- automação de gameplay;
- cliente alternativo;
- qualquer ferramenta destinada a obter vantagem injusta durante a partida.

O foco é:

- inteligência da conta;
- discovery;
- progressão;
- análise pós-partida;
- coaching.

---

# 15. Roadmap definido

O roadmap atual está em:

`docs/roadmap/crownpilot-roadmap.md`

Sequência principal:

```text
001 Viabilidade de produto, dados e compliance
        |
        +--> 002 Identidade persistente
        |        |
        |        v
        |    003 Player Sync
        |
        +--> 004 Dataset de meta
                 |
                 +---------+
                           v
                 005 Best Decks for You
                           |
                           v
                 006 Upgrade Planner
                           |
                           v
                 007 MVP / Beta
                           |
                           v
                 008 History & Matchups
                           |
                           v
                 009 Optimize My Deck
                           |
                           v
                 010 Player Personalization
                           |
                           v
                 011 AI Coach

001 + evidência do Beta + policy
        |
        v
012 Commercialization & Scale
```

---

# 16. Fase 001

Spec:

`docs/specs/001-viabilidade-produto-dados-compliance.md`

Tasks:

`docs/tasks/001-viabilidade-produto-dados-compliance/`

A fase possui oito subtarefas:

1. mapear API oficial e autenticação;
2. validar player profile e collection;
3. validar battle log e histórico;
4. validar aquisição do meta;
5. medir operação, cache e rate limits;
6. validar compliance e monetização;
7. definir contratos de dados v0;
8. fechar gates e handoff.

Resultado obrigatório:

- **GO**;
- **GO WITH CONSTRAINTS**;
- **NO-GO / REDESIGN**.

---

# 17. Principais hipóteses que ainda precisam ser provadas

## Dados da conta

Precisamos confirmar na API real:

- coleção completa;
- níveis;
- Arena;
- trophies;
- current deck;
- Evolution ownership;
- Hero ownership;
- progressão atual;
- campos deprecated ou instáveis.

## Battle log

Precisamos saber:

- tamanho da janela;
- paginação;
- modos;
- representação do deck;
- resultado;
- opponent;
- possibilidade de snapshots contínuos;
- frequência necessária para evitar gaps.

## Dataset de meta

É provavelmente o maior risco técnico.

Precisamos descobrir uma estratégia sustentável para obter decks relevantes em diferentes faixas.

Possibilidades a investigar:

- rankings;
- Path of Legend;
- clans/members;
- expansão via opponents;
- fontes third-party licenciadas;
- estratégia híbrida.

## Operação

Precisamos medir:

- cache;
- rate behavior;
- custo em requests;
- egress/IP allowlist;
- compatibilidade com Vercel;
- necessidade eventual de gateway com IP estável.

---

# 18. Contratos de dados esperados

A Fase 001 deve produzir versões provisórias de:

- `PlayerSnapshotV0`;
- `CardCollectionEntryV0`;
- `CompetitiveContextV0`;
- `CurrentDeckV0`.

Eles devem separar:

- dado bruto;
- dado normalizado;
- dado derivado;
- freshness;
- provenance;
- optionality.

Especial atenção deve existir para:

- Evolution capability;
- Evolution/Hero ownership;
- Evolution/Hero deployment.

Esses conceitos não devem ser confundidos.

---

# 19. Padrão de documentação

Foi adotado o padrão já utilizado no projeto Reserva Clara.

Estrutura:

```text
docs/
├── decisions/
├── roadmap/
├── specs/
└── tasks/
```

Cada fase possui:

- uma spec principal em `docs/specs/`;
- uma pasta de tasks;
- `00-overview.md`;
- subtarefas numeradas;
- critérios de aceite;
- registro de execução;
- handoff.

Resumo de responsabilidades:

- **Roadmap:** o que e em qual ordem;
- **Spec:** como a fase será implementada/executada;
- **Task:** unidade executável;
- **Decision/ADR:** decisão arquitetural duradoura;
- **Context:** histórico e motivação;
- **main:** fonte final de verdade do que realmente existe; não substitui decisões
  normativas sobre intenção futura.

---

# 20. Próximo passo

O próximo passo planejado é executar:

> **001-01 — Mapear API oficial e autenticação**

Antes de implementar a aplicação, precisamos responder com evidência:

1. o que conseguimos conhecer a partir de uma Player Tag;
2. o que conseguimos conhecer do meta relevante;
3. quanto custa operacionalmente obter esses dados;
4. quais limites técnicos e comerciais existem;
5. se esses dados sustentam de fato **Best Decks for You**.

Somente depois disso a arquitetura de domínio e a modelagem persistente deveriam
começar a ser congeladas. A referência histórica a Firestore foi substituída por
PostgreSQL via EF Core + Npgsql na Fase 002; não orienta a implementação atual.
