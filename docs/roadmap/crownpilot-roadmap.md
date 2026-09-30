# Roadmap — CrownPilot

> **Play the right deck. Upgrade the right cards.**

**Atualizado em:** 30 de setembro de 2026

Este documento define **o que construir**, **em qual ordem**, as principais dependências, os boundaries do produto e os grandes marcos do CrownPilot.

O roadmap **não substitui specs detalhadas**. Specs definem como uma fase será implementada. O roadmap define direção, sequência, dependências e critérios para liberar a próxima etapa.

---

# 1. Visão do produto

CrownPilot é um companion personalizado para Clash Royale.

A pergunta central do produto é:

> **Qual é o melhor deck que eu consigo jogar agora, na minha Arena, com a conta que eu realmente tenho?**

A experiência deve partir da conta do jogador, e não de um ranking genérico.

Quando os dados permitirem, a recomendação combina:

- Arena e faixa de troféus;
- coleção de cartas;
- níveis;
- Evolutions / Heroes disponíveis;
- decks observados em partidas reais;
- força do deck no meta relevante;
- tamanho e confiança da amostra;
- histórico e desempenho individual em fases posteriores.

O produto evolui em quatro camadas:

- **DISCOVER** — O que eu deveria jogar?
- **PROGRESS** — O que eu deveria melhorar?
- **OPTIMIZE** — Como melhorar o que eu já jogo?
- **COACH** — Por que isso faz sentido para mim?

## 1.1 Proposta de valor

O CrownPilot não deve responder apenas:

> Quais são os melhores decks do meta?

Ele deve responder:

> Quais são os melhores decks **para mim**?

Dois jogadores na mesma Arena podem receber recomendações diferentes porque possuem coleções, níveis, Evolutions / Heroes e progressão diferentes.

## 1.2 Hipótese principal

O produto gera valor recorrente quando consegue transformar:

**conta do jogador + contexto competitivo + decks reais + progressão**

em:

**recomendações acionáveis e explicáveis.**

O objetivo não é ser apenas um deck finder.

O objetivo é se tornar a camada de inteligência da conta do jogador.

---

# 2. Princípios arquiteturais e de produto

## 2.1 Deterministic by default

> **Deterministic by default. AI when it adds value.**

Toda feature que puder ser resolvida de forma determinística deve funcionar sem LLM.

Isso inclui inicialmente:

- Best Decks for You;
- Deck Readiness;
- Arena / Trophy Fit;
- Collection Fit;
- Upgrade Impact;
- Upgrade Planner;
- comparação de variantes;
- métricas de matchup.

IA deve interpretar e explicar resultados do engine, não substituir o engine.

## 2.2 A conta é o centro do produto

O usuário vincula sua Player Tag uma única vez.

Depois de autenticado em qualquer dispositivo, o CrownPilot deve recuperar a conta vinculada sem exigir que o usuário abra o jogo para copiar a tag novamente.

Nunca armazenar credenciais da conta Supercell.

## 2.3 Recomendações precisam de evidência

A primeira versão não deve gerar oito cartas do zero.

O recommendation engine começa com **decks reais observados** e ranqueia esses candidatos para o jogador.

Deck generation só entra em exploração depois que Discover e Optimize estiverem validados.

## 2.4 Arena é contexto, não apenas filtro

Arena / faixa de troféus influencia a relevância.

O produto não deve assumir que o deck mais forte no topo da ladder é automaticamente a melhor recomendação para um jogador em outra faixa competitiva.

## 2.5 Toda recomendação deve ser explicável

Para qualquer deck recomendado, o sistema precisa responder:

> Por que este deck está sendo recomendado para este jogador?

O score deve ser decomponível em fatores observáveis.

## 2.6 Dados vencidos não são meta

Balance changes e mudanças de temporada podem invalidar dados antigos rapidamente.

Toda estatística competitiva deve carregar:

- período;
- população;
- tamanho da amostra;
- contexto;
- confiança.

## 2.7 Infraestrutura-base e portabilidade

Decisões iniciais do projeto:

- Firebase Authentication com Google para a conta CrownPilot;
- Cloud Firestore como banco principal;
- Vercel como plataforma inicial de deploy.

Vercel não deve se tornar um boundary do domínio.

O core da aplicação e integrações devem permanecer portáveis para outro runtime
sem reescrita do domínio. Serviços exclusivos da Vercel podem ser usados apenas
quando isolados atrás de adapters ou quando existir estratégia clara de
substituição.

A integração com a Clash Royale API deve permanecer separada do runtime,
especialmente por possíveis requisitos de egress/IP allowlist.

## 2.8 Segurança e ambientes fazem parte da fundação

A fundação deve nascer segura por padrão:

- segredos de integração permanecem somente no servidor;
- acesso client-side ao Firestore exige Security Rules versionadas e testadas;
- acesso server-side ao Firestore usa IAM/Admin SDK com privilégio mínimo;
- local, preview e produção não devem compartilhar dados/segredos de forma
  acidental;
- Firebase Emulator Suite deve ser usado quando trouxer isolamento e testes
  reproduzíveis;
- decisões difíceis de reverter, como localização do Firestore, precisam ser
  tomadas e registradas antes de provisionar produção.

A estratégia exata de client vs. server access será definida na Fase 002, mas
não é permitido tratar autenticação como autorização.

## 2.9 Compliance é boundary arquitetural

CrownPilot é um companion de análise e coaching.

Não fazem parte do produto:

- automação de gameplay;
- bots;
- cliente alternativo;
- interferência no jogo;
- overlay competitivo em tempo real;
- ferramentas destinadas a obter vantagem injusta durante a partida.

Monetização também não deve ser assumida como livre de restrições.

A Fan Content Policy vigente da Supercell deve ser tratada como requisito do produto:

https://supercell.com/en/fan-content-policy/

---

# 3. Convenções de status

| Status | Significado |
|---|---|
| ✅ | Concluído e validado |
| 🚧 | Em andamento |
| ⬜ | Planejado |
| ⏸️ | Bloqueado por decisão ou dependência |
| 🧪 | Discovery / investigação |

Uma fase só é concluída quando:

- entregas previstas estão implementadas;
- critérios de aceite foram validados;
- decisões relevantes estão documentadas;
- riscos e débitos remanescentes estão explícitos;
- existe handoff claro para a próxima fase.

---

# 4. Fases

## 001 — Viabilidade de produto, dados e compliance 🚧

### Objetivo

Validar se os dados e permissões disponíveis sustentam o produto antes de comprometer a arquitetura.

### Escopo

Investigar e documentar:

- acesso à API oficial;
- consulta de jogador por Player Tag;
- campos disponíveis no perfil;
- Arena e troféus;
- coleção e níveis;
- representação disponível de Evolutions / Heroes;
- battle log;
- cards e locations;
- rankings ou outras superfícies úteis para amostragem;
- rate limits;
- limites de histórico;
- estratégia de cache;
- restrições de uso dos dados;
- políticas aplicáveis;
- semântica de vínculo da Player Tag e possibilidade de prova de ownership.

### Perguntas críticas

1. Conseguimos reconstruir de forma confiável o estado da coleção necessário para Deck Readiness?
2. Conseguimos determinar Arena / faixa competitiva sem heurísticas frágeis?
3. Como obter uma amostra representativa de decks por faixa?
4. O battle log é suficiente para histórico contínuo por snapshots?
5. Precisaremos de fonte externa complementar para meta?
6. Quais dados econômicos não estão disponíveis e exigiriam entrada manual?
7. Quais limites de uso impactam sync e ingestão em escala?
8. Vincular uma Player Tag significa selecionar um perfil público ou precisamos
   provar ownership? Existe mecanismo oficial aplicável e ele é necessário para
   o MVP?

### Compliance comercial

Antes de qualquer assinatura, paywall ou cobrança por feature, validar explicitamente o modelo comercial contra a política vigente e demais acordos aplicáveis.

A política pública atual da Supercell descreve fan content em geral como não comercial e lista exceções específicas como anúncios, doações e coaching.

**Não interpretar isso automaticamente como autorização para qualquer plano Pro.**

### Entregas

- mapa das fontes de dados;
- matriz campo → fonte → freshness → limitação;
- spike de integração;
- modelo inicial de Player Snapshot;
- modelo inicial de Card Collection;
- proposta de aquisição do dataset de meta;
- registro dos principais riscos;
- decisão documentada sobre viabilidade do MVP.

### Critérios de aceite

- [ ] Player Tag válida pode ser resolvida de forma reproduzível;
- [ ] campos do MVP estão classificados como disponíveis, derivados, opcionais ou indisponíveis;
- [ ] estratégia de meta possui uma fonte inicial viável;
- [ ] limites de rate e histórico estão documentados;
- [ ] boundary de compliance está documentado;
- [ ] semântica de vínculo/ownership da Player Tag está decidida;
- [ ] nenhuma dependência crítica do MVP continua baseada apenas em suposição.

### Handoff

A Fase 002 só fixa contratos persistentes depois que a Fase 001 confirmar quais dados realmente existem.

---

## 002 — Fundação da aplicação e identidade persistente ⬜

### Objetivo

Criar a fundação executável do CrownPilot, estabelecer seus boundaries de
segurança e entregar a identidade persistente que elimina o atrito de informar a
Player Tag em cada dispositivo.

### Escopo

#### Bootstrap da aplicação

- inicializar a aplicação e escolher o framework web;
- definir package manager e lockfile autoritativo;
- TypeScript strict e configuração de build;
- lint/format e convenções de código;
- estrutura inicial de diretórios e boundaries;
- `AGENTS.md` e comandos operacionais do repositório;
- test runner mínimo;
- CI com gates de lint, type-check, testes e build;
- `.env.example` sem secrets;
- estratégia local / preview / production;
- deploy inicial na Vercel sem dependência obrigatória de serviços proprietários.

#### Firebase

- Firebase Authentication com Google;
- Cloud Firestore como banco principal;
- decidir e registrar localização do Firestore antes de provisionar produção;
- Firebase CLI e Emulator Suite para desenvolvimento/testes quando aplicável;
- definir boundary de acesso ao Firestore:
  - client SDK + Security Rules; ou
  - server SDK/Admin + IAM;
  - ou combinação explicitamente documentada;
- versionar e testar Security Rules para qualquer acesso client-side;
- garantir que preview/local não usem produção por acidente.

#### Identidade CrownPilot

- usuário CrownPilot;
- uma Player Tag primária no MVP;
- vínculo persistente com Player Tag;
- validação de existência da tag;
- aplicar a decisão da Fase 001 sobre perfil público vs. ownership verificado;
- tag inválida / jogador inexistente;
- troca e desvinculação controladas;
- recuperação do vínculo em outro dispositivo;
- estratégia inicial para exclusão dos dados da conta;
- observabilidade inicial.

### Boundary da fase

A Fase 002 pode consultar a Clash Royale API para validar o vínculo, mas **não**
sincroniza nem persiste coleção, níveis, Arena, battle history ou Player
Snapshot completo. Isso começa na Fase 003.

O token da Clash Royale API nunca é exposto ao browser.

### Requisitos

- nenhuma credencial da Supercell é solicitada;
- Firebase Authentication identifica o usuário CrownPilot;
- Player Tag é vínculo de domínio, com semântica definida pela Fase 001;
- nova sessão recupera a tag vinculada;
- mudança/desvinculação da tag é explícita;
- autenticação não é tratada como autorização;
- falhas da API externa não invalidam a identidade local;
- runtime/domain não dependem de API proprietária da Vercel.

### Critérios de aceite

- [ ] bootstrap pode ser reproduzido a partir do repositório limpo;
- [ ] lint, type-check, testes e build possuem comandos definidos e passam;
- [ ] CI executa os gates mínimos;
- [ ] ambientes local/preview/production estão separados e documentados;
- [ ] localização do Firestore está decidida antes do banco de produção;
- [ ] Firestore não possui acesso público irrestrito;
- [ ] Security Rules/IAM refletem o boundary escolhido e possuem validação;
- [ ] usuário consegue criar sessão com Google;
- [ ] usuário vincula uma Player Tag uma vez;
- [ ] outro dispositivo recupera o vínculo após login;
- [ ] usuário consegue trocar/desvincular a tag;
- [ ] erros de integração não causam perda do vínculo;
- [ ] nenhum token da Supercell ou credencial de serviço chega ao client;
- [ ] deploy inicial na Vercel funciona sem tornar Vercel parte do domínio.

### Dependências

- Fase 001.

### Handoff

Entregar aplicação reproduzível, identidade estável, Firestore seguro e
boundaries suficientes para a Fase 003 implementar o sync da conta.

---

## 003 — Sync da conta e Player Dashboard ⬜

### Objetivo

Transformar a Player Tag em um estado de conta útil dentro do CrownPilot.

### Escopo

Sincronizar e normalizar, conforme disponibilidade confirmada na Fase 001,
sempre através de integração server-side com a Clash Royale API:

- definir política de freshness e gatilhos de refresh;
- separar payload bruto, normalização e snapshot de domínio;
- tratar concorrência/retry para impedir snapshots inconsistentes;
- definir retenção mínima necessária antes de acumular histórico;

Sincronizar:

- perfil;
- Arena;
- troféus;
- coleção;
- níveis;
- Evolutions / Heroes;
- deck atual quando disponível;
- timestamps de atualização.

### Dashboard inicial

O dashboard deve responder:

- quem é este jogador?
- onde ele está?
- qual o estado da coleção?
- quando os dados foram atualizados?
- quais dados não conseguimos observar?

### Arquitetura

Separar resposta externa de domínio interno:

**External API Response → Normalization → Player Snapshot**

O domínio do CrownPilot não deve depender diretamente do formato bruto do provedor.

### Freshness

Toda informação sincronizada possui semântica clara de freshness.

Snapshot antigo não deve parecer estado atual.

### Critérios de aceite

- [ ] sync completo do jogador funciona;
- [ ] coleção é normalizada para o recommendation engine;
- [ ] Arena / troféus são persistidos com timestamp;
- [ ] dados ausentes não quebram o perfil;
- [ ] freshness fica visível;
- [ ] sync repetido é idempotente;
- [ ] token da Clash Royale API permanece server-side;
- [ ] payload bruto e Player Snapshot normalizado possuem boundary explícito;
- [ ] política de freshness/retry está documentada.

### Dependências

- Fases 001 e 002.

### Handoff

Player Snapshot se torna o contrato de entrada da personalização.

---

## 004 — Dataset de meta e canonicalização de decks ⬜

### Objetivo

Construir uma base confiável de decks reais e contexto competitivo.

### Escopo

- estratégia de amostragem de jogadores / batalhas;
- ingestão de battle data;
- canonicalização de deck;
- representação de Evolutions / Heroes;
- agrupamento de decks equivalentes;
- frequência;
- win rate quando calculável;
- tamanho da amostra;
- faixa de Arena / troféus;
- janela temporal;
- confidence;
- decaimento após mudanças relevantes;
- versionamento de catálogo/temporada/balance context;
- estratégia de retenção e agregação;
- validação do custo e do modelo de armazenamento no Firestore;
- possibilidade de componente analítico especializado se Firestore deixar de
  ser adequado para dados brutos em escala.

### Regra central

Um deck só entra no recommendation engine como candidato quando existe evidência suficiente para tratá-lo como deck observado, e não combinação arbitrária.

### Estatística

Nenhuma taxa deve ser apresentada sem:

- população;
- período;
- tamanho da amostra;
- contexto competitivo.

### Critérios de aceite

- [ ] decks possuem identidade canônica;
- [ ] dataset distingue contexto competitivo;
- [ ] métricas carregam janela temporal;
- [ ] amostras pequenas são penalizadas ou descartadas;
- [ ] balance changes podem reduzir a relevância dos dados antigos;
- [ ] recommendation engine consulta candidatos eficientemente;
- [ ] mudanças de temporada/balance podem separar ou invalidar amostras;
- [ ] retenção e agregação possuem estratégia explícita;
- [ ] custo/query pattern foi validado para Firestore ou existe boundary para
  store analítico especializado.

### Dependências

- Fase 001;
- Fase 002 para fundação executável, configuração e boundaries de integração.

### Handoff

Entregar catálogo de candidatos com evidência suficiente para personalização.

---

## 005 — Recommendation Engine v1: Best Decks for You ⬜

### Objetivo

Entregar a primeira feature que define o CrownPilot.

> **Best Decks for You**

### Entrada

- Player Snapshot;
- Deck Candidates;
- Competitive Context.

### Saída

Lista ordenada de decks acompanhada do motivo da recomendação.

### Fit Score

O score deve começar simples e decomponível.

Componentes conceituais:

- Meta Strength;
- Collection Compatibility;
- Level Readiness;
- Arena Relevance;
- Confidence.

A fórmula final não é definida neste roadmap. A spec da fase decide pesos, normalização e forma de combinação com dados reais.

### Readiness

Cada recomendação deve possuir classificação compreensível:

- READY TO PLAY;
- 1 UPGRADE AWAY;
- 2 UPGRADES AWAY;
- MISSING EVOLUTION / HERO;
- NOT READY.

### Regras

- não usar LLM para calcular ranking;
- não gerar deck arbitrário;
- mesma entrada produz mesma saída;
- ausência de dados reduz confiança;
- ranking global não substitui contexto do jogador.

### Critérios de aceite

- [ ] engine recebe Player Snapshot normalizado;
- [ ] candidatos vêm de decks observados;
- [ ] score é determinístico;
- [ ] score é decomponível;
- [ ] readiness é calculada;
- [ ] usuário entende por que um deck foi recomendado;
- [ ] Arena / faixa competitiva influencia a relevância quando os dados suportarem isso.

### Dependências

- Fases 003 e 004.

### Resultado

O CrownPilot passa a responder sua pergunta central.

---

## 006 — Progress: Upgrade Planner ⬜

### Objetivo

Responder:

> **O que eu deveria melhorar agora?**

### Experiência

- **NOW** — decks que posso jogar hoje;
- **NEXT** — decks que ficam prontos com pequenas melhorias;
- **TARGET** — decks relevantes que fazem sentido perseguir como objetivo.

### Upgrade Impact

Cada upgrade deve ser avaliado pelo impacto que gera.

Exemplo conceitual:

- uma carta sobe de nível;
- quatro decks passam para READY;
- sete decks melhoram Level Readiness.

Isso deve pesar mais do que um upgrade que melhora apenas um deck pouco relevante.

### Boundary

Se Gold, Wild Cards ou outros recursos econômicos não estiverem disponíveis na fonte de dados, o engine não inventa disponibilidade.

Entrada manual pode ser considerada em spec futura.

### Critérios de aceite

- [ ] cada recomendação de upgrade possui impacto explicável;
- [ ] NOW / NEXT / TARGET derivam do estado real da conta;
- [ ] nenhuma recomendação assume recursos desconhecidos;
- [ ] usuário navega de upgrade para decks impactados;
- [ ] usuário navega de deck para upgrades necessários.

### Dependências

- Fase 005.

### Resultado

O CrownPilot deixa de ser apenas discovery e passa a guiar progressão.

---

## 007 — MVP Experience e Beta ⬜

### Objetivo

Transformar os engines em uma experiência coerente e utilizável regularmente.

### Onboarding

**Sign in → Link Player Tag → Sync Account → Best Decks for You**

### Home

Priorizar:

- estado atual da conta;
- Arena / troféus;
- Best Decks for You;
- NOW / NEXT / TARGET;
- próximo upgrade recomendado;
- freshness.

### Produto

- mobile-first / responsive;
- loading e erro consistentes;
- estados vazios;
- cache;
- sync controlado;
- deep links quando fizer sentido;
- analytics de produto;
- logs e observabilidade;
- privacidade;
- fluxo de exclusão de conta/dados;
- revisão de autorização e Security Rules/IAM;
- avaliar App Check se houver acesso client-side direto ao Firebase;
- disclaimer de conteúdo não oficial.

### Métricas iniciais

Medir:

- Player Tag vinculada;
- sync concluído;
- recomendações geradas;
- deck visualizado;
- deck copiado / aberto quando aplicável;
- upgrade recomendado visualizado;
- retorno após mudança de progressão.

Metas numéricas entram na spec de beta depois de existir baseline.

### Critérios de aceite

- [ ] fluxo completo funciona em dispositivo novo após login;
- [ ] recomendações aparecem sem reentrada da tag;
- [ ] falha de provedor degrada de forma controlada;
- [ ] observabilidade mínima existe;
- [ ] eventos críticos são medidos;
- [ ] compliance e disclaimer estão visíveis;
- [ ] usuário possui caminho funcional para excluir sua conta e dados próprios;
- [ ] revisão de segurança do acesso ao Firestore foi concluída;
- [ ] MVP funciona sem IA.

### Dependências

- Fases 002–006.

### Marco

Primeira versão capaz de validar utilidade real com usuários.

---

## 008 — Histórico, Performance e Matchups ⬜

### Objetivo

Adicionar contexto do próprio jogador às recomendações.

### Escopo

- snapshots contínuos de battle history;
- retenção de histórico próprio quando permitido;
- classificação de decks / arquétipos;
- performance por deck;
- performance por matchup;
- tendências;
- período selecionável;
- confiança por volume.

### Boundary

A análise é pós-partida.

Não existe objetivo de assistência competitiva em tempo real.

### Critérios de aceite

- [ ] histórico pode persistir além da janela externa quando permitido;
- [ ] partidas são classificáveis;
- [ ] métricas mostram sample size;
- [ ] usuário consegue distinguir tendência de ruído.

### Dependências

- Fases 001, 003 e 004.

---

## 009 — Optimize My Deck ⬜

### Objetivo

Responder:

> **Como posso melhorar o deck que eu já jogo?**

### Estratégia

A primeira versão procura **variantes comprovadas**.

A recomendação considera:

- nível da alternativa na conta;
- evidência da variante;
- matchups relevantes;
- custo de oportunidade;
- contexto competitivo.

### Não fazer

Não pedir a um LLM para simplesmente inventar uma substituição.

### Critérios de aceite

- [ ] variante possui evidência;
- [ ] substituição é compatível com a coleção;
- [ ] motivo da troca é explicável;
- [ ] comparação mostra trade-offs.

### Dependências

- Fases 004, 005 e 008.

---

## 010 — Personalização por desempenho do jogador ⬜

### Objetivo

Fazer contas semelhantes receberem recomendações diferentes quando o histórico justificar isso.

### Escopo

- performance pessoal;
- arquétipos mais utilizados;
- consistência;
- matchups;
- confiança mínima;
- ajuste limitado do ranking.

### Regra

Personalização histórica não domina o recommendation engine com amostra pequena.

Sem evidência suficiente, prevalece:

**account + meta**

Com evidência suficiente:

**account + meta + player performance**

### Dependências

- Fases 008 e 009.

---

## 011 — AI Coach ⬜

### Objetivo

Adicionar linguagem natural e coaching sobre resultados já calculados pelo CrownPilot.

### Arquitetura

**Player Data → Recommendation / Analytics Engines → Structured Insight → AI Coach**

O modelo não deve reconstruir o domínio a partir de dados brutos quando o backend já consegue produzir resposta estruturada.

### Casos de uso

- Por que este deck é melhor para mim?
- Por que devo melhorar esta carta?
- O que mudou nas minhas últimas partidas?
- Por que estou sofrendo contra determinado arquétipo?
- Qual é meu plano para chegar ao próximo objetivo?

### Saúde operacional

A feature deve possuir:

- budget por chamada;
- limite de contexto;
- cache quando seguro;
- telemetria de custo;
- limites de uso;
- fallback sem IA;
- isolamento do provedor.

### Monetização

Qualquer plano pago, crédito ou franquia de IA depende do gate comercial das Fases 001 e 012.

Não construir billing primeiro e procurar justificativa depois.

### Critérios de aceite

- [ ] AI Coach usa resultados estruturados;
- [ ] resposta pode ser rastreada até dados do engine;
- [ ] indisponibilidade do modelo não quebra o core;
- [ ] custo por interação é observável;
- [ ] hallucination não altera scores ou fatos persistidos;
- [ ] modelo comercial, se existir, está validado.

### Dependências

- Fases 005–010 conforme o caso;
- gate de compliance.

---

## 012 — Comercialização e escala ⬜

### Objetivo

Definir um modelo sustentável somente depois de provar utilidade e validar o que é permitido.

### Gate obrigatório

Antes de cobrar:

- revisar Fan Content Policy vigente;
- revisar acordos da API / developer;
- confirmar enquadramento das features;
- registrar a decisão;
- buscar autorização adicional quando necessária.

### Princípio

Arquitetura de custos deve existir mesmo se o produto continuar gratuito.

Especialmente para IA, observar:

- cost per active user;
- cost per AI interaction;
- cache hit rate;
- provider spend;
- usage distribution.

### Regra

Não assumir no roadmap que assinatura SaaS tradicional é permitida.

### Dependências

- utilidade comprovada no beta;
- métricas reais;
- gate jurídico / policy;
- custos medidos.

---

# 5. Explorações futuras

Estas ideias não fazem parte do MVP.

## 5.1 Build Me a Deck

Geração de novos candidatos considerando:

- cartas obrigatórias;
- sinergias;
- função das cartas;
- custo de elixir;
- matchups;
- coleção;
- níveis;
- dados históricos.

Qualquer deck gerado deve ser validado contra dados antes de ser apresentado como competitivo.

## 5.2 Alertas de progressão

Exemplos:

- nova Arena detectada;
- novos decks ficaram disponíveis;
- upgrade alterou ranking;
- mudança de meta afetou recomendações.

## 5.3 Planejamento avançado

- múltiplos objetivos;
- custo de progressão;
- caminhos alternativos;
- comparação entre estratégias de upgrade.

Depende da disponibilidade real de dados econômicos.

---

# 6. Mapa de dependências

Fluxo principal:

**001 Viabilidade / Dados / Compliance**

A partir dela:

**002 Fundação + Identidade**

Depois da fundação, 003 e 004 podem avançar em paralelo:

- **003 Player Sync**
- **004 Meta Dataset**

As Fases 003 e 004 convergem em:

**005 Best Decks for You → 006 Upgrade Planner → 007 MVP / Beta**

Depois:

**008 History & Matchups → 009 Optimize My Deck → 010 Player Personalization → 011 AI Coach**

Comercialização depende de:

**001 + evidência do Beta + policy validation → 012 Commercialization & Scale**

---

# 7. Sequência crítica do MVP

1. **001** — provar dados e boundaries;
2. **002** — bootstrap, segurança e identidade persistente;
3. **003 + 004 em paralelo** — Player Snapshot e dataset competitivo;
4. **005** — Best Decks for You;
5. **006** — Upgrade Planner;
6. **007** — MVP / Beta.

O maior risco técnico não é UI.

É provar que conseguimos obter e manter dados suficientes para:

1. compreender a conta;
2. construir candidatos de meta confiáveis;
3. contextualizar por Arena / faixa;
4. explicar a recomendação.

Por isso a Fase 001 vem antes de decisões profundas de stack.

---

# 8. Grandes marcos

## Marco A — Evidence

**Fase 001**

Sabemos quais dados temos, o que não temos e como construir o dataset necessário.

## Marco B — Connected Account

**Fases 002–003**

Aplicação possui fundação reproduzível; usuário entra em qualquer dispositivo e
o CrownPilot conhece e sincroniza sua conta.

## Marco C — Competitive Dataset

**Fase 004**

Temos decks reais, contexto e confiança suficientes para alimentar o engine.

## Marco D — Best Decks for You

**Fase 005**

A proposta central funciona.

## Marco E — Progression Intelligence

**Fase 006**

O produto recomenda não apenas o que jogar, mas o que melhorar.

## Marco F — Usable Beta

**Fase 007**

Existe experiência completa capaz de gerar aprendizado real.

## Marco G — Player Intelligence

**Fases 008–010**

Histórico e comportamento passam a melhorar recomendações.

## Marco H — Coach

**Fase 011**

IA adiciona interpretação sem substituir o engine.

## Marco I — Sustainable Product

**Fase 012**

Modelo operacional e comercial está validado técnica e legalmente.

---

# 9. Boundaries do MVP

O MVP inclui:

- bootstrap reproduzível e quality gates;
- Firebase Authentication com Google;
- Firestore seguro e ambientes separados;
- deploy inicial portável na Vercel;
- identidade CrownPilot;
- Player Tag persistente;
- sync da conta;
- Arena / troféus;
- coleção normalizada;
- dataset inicial de decks;
- Best Decks for You;
- Deck Readiness;
- Upgrade Planner;
- NOW / NEXT / TARGET;
- experiência responsiva;
- observabilidade básica.

O MVP não inclui:

- múltiplas Player Tags por usuário;
- AI Coach;
- geração livre de deck;
- overlay;
- tracking em tempo real durante partida;
- automação;
- billing;
- assinatura;
- social;
- clãs;
- torneios;
- app nativo obrigatório.

---

# 10. Roadmap, Specs, ADRs e repositório

## Roadmap

Responde:

> **O quê e em qual ordem?**

Este arquivo é o mapa de evolução.

## Specs

Respondem:

> **Como esta fase será implementada?**

Uma spec deve registrar, no mínimo:

- baseline;
- problema;
- objetivo;
- não objetivos;
- decisões;
- arquitetura;
- contratos;
- persistência;
- tasks ordenadas;
- critérios de aceite;
- riscos;
- testes;
- observabilidade;
- rollout;
- dependências;
- handoff.

## ADRs

Registram decisões arquiteturais que precisam sobreviver à implementação.

Exemplos futuros:

- estratégia de aquisição de meta;
- Player Snapshot;
- algoritmo do Fit Score;
- freshness;
- retenção de battle history;
- isolamento do provider de IA.

## Repositório

A implementação na main é a fonte final de verdade sobre o que existe.

Se roadmap, spec e código divergirem, investigar antes de continuar.

Não atualizar documentação para fingir que algo foi entregue.

---

# 11. Workflow por fase

Antes de planejar:

1. ler este roadmap;
2. inspecionar a main;
3. ler specs anteriores relevantes;
4. ler ADRs relacionados;
5. ler o handoff anterior;
6. confirmar baseline real;
7. identificar spikes necessários.

Durante a fase:

1. manter escopo explícito;
2. registrar decisões;
3. validar por checkpoints;
4. evitar antecipar fases futuras sem necessidade;
5. manter critérios de aceite verificáveis.

Ao finalizar:

1. validar testes;
2. validar comportamento real;
3. atualizar documentação;
4. registrar riscos e débitos;
5. produzir handoff.

---

# 12. Handoff padrão

Toda fase concluída termina com:

### Status

O que foi concluído e o que permaneceu fora.

### Resultado principal

Qual capacidade nova existe.

### Decisões

Quais decisões importantes foram tomadas.

### Arquitetura

Quais contratos e boundaries passam a importar.

### Validação

Quais testes, checks e evidências confirmam a entrega.

### Riscos e débitos

O que ainda pode afetar a próxima fase.

### Dependências

O que a próxima fase pode assumir como verdade.

### Próximo passo

Qual fase está liberada e por quê.

---

# 13. Próximo passo

A próxima fase ativa é:

> **001 — Viabilidade de produto, dados e compliance**

Antes de escolher stack definitiva ou implementar o recommendation engine, o CrownPilot precisa provar seus contratos fundamentais de dados.

O primeiro deliverable técnico deve responder com evidência:

1. O que conseguimos saber de forma confiável a partir de uma Player Tag?
2. O que conseguimos saber de forma confiável sobre o meta relevante?
3. Esses dois conjuntos de dados sustentam Best Decks for You?

Se a resposta for sim, a arquitetura deixa de ser hipótese e começa a ser produto.
