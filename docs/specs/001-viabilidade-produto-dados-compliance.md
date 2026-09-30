# 001 — Viabilidade de produto, dados e compliance

## Ticker

`001`

## Contexto

CrownPilot nasce como companion personalizado para Clash Royale com uma pergunta
central: **qual é o melhor deck que este jogador consegue jogar agora, na Arena
e no estado real da conta que possui?**

O roadmap definiu como primeira fase uma etapa de discovery obrigatória antes de
implementar o recommendation engine. A infraestrutura-base do projeto já foi
definida pelo responsável:

- **Google Sign-In via Firebase Authentication** para identidade CrownPilot;
- **Cloud Firestore** como banco principal;
- **Vercel** como runtime/deploy inicial;
- **portabilidade de runtime** como requisito: domínio e integrações não podem
  depender de APIs exclusivas da Vercel sem uma camada substituível.

Na `main`, o repositório contém somente documentação de produto e planejamento:
`README.md`, roadmap, contexto histórico, ADR, esta spec e as oito subtarefas da
fase. Não existe aplicação, integração com a API, evidência de probe, modelo
persistente, dataset de meta, pipeline de ingestão ou autenticação implementada.

A viabilidade do produto depende de quatro provas independentes:

1. **conta:** a API precisa expor dados suficientes para representar coleção,
   níveis, Arena/progressão e Evolutions/Heroes;
2. **histórico:** precisamos entender exatamente o que o battle log oferece,
   por quanto tempo e com quais semânticas;
3. **meta:** precisamos de uma forma sustentável e permitida de construir ou
   consumir amostras de decks por contexto competitivo, especialmente fora do
   topo da ladder;
4. **compliance:** uso de API, dados, assets, marca e eventual monetização devem
   respeitar os termos atuais da Supercell.

A documentação oficial em `developer.clashroyale.com` é a autoridade para a
API. Como parte do discovery, observações de fontes secundárias podem orientar
hipóteses, mas nenhum comportamento de campo, limite ou endpoint será tratado
como contrato sem confirmação por documentação oficial ou probe controlado da
API real.

A Fan Content Policy pública consultada em 29/09/2026 informa "Last updated:
September 27, 2023" e estabelece limites relevantes para uso comercial. Isso
torna monetização um gate explícito, principalmente para qualquer plano Pro ou
feature paga baseada em IA.

## Objetivo

Eliminar as principais incógnitas que podem inviabilizar ou alterar o produto e
encerrar a fase com uma decisão objetiva:

- **GO:** dados e permissões sustentam o MVP planejado;
- **GO WITH CONSTRAINTS:** MVP é viável com limitações explícitas e roadmap
  ajustado;
- **NO-GO / REDESIGN:** a hipótese atual não é sustentável e precisa mudar antes
  da Fase 002.

Ao final, a equipe deve saber com evidência:

1. o que uma Player Tag permite conhecer;
2. o que a API não oferece;
3. como obter contexto de meta sem violar contratos;
4. qual freshness e custo de requests esperar;
5. quais boundaries comerciais e de marca se aplicam;
6. qual contrato mínimo de dados a Fase 002 pode assumir.

## Requisitos

### Funcionais de discovery

- Mapear endpoints relevantes do portal oficial e seus contratos observados.
- Validar Player Tag real e encoding.
- Investigar a semântica de vínculo da Player Tag: seleção de perfil público vs.
  ownership verificado, incluindo mecanismo oficial de verificação se existir e
  for aplicável.
- Validar payload de perfil e coleção em amostra controlada.
- Confirmar representação atual de níveis, Arena, troféus, current deck,
  Evolutions e Heroes.
- Identificar campos deprecated, instáveis, opcionais e derivados.
- Validar battle log, tamanho de janela e campos suficientes para decks e
  resultados.
- Avaliar estratégias para construir dataset de meta por faixa competitiva.
- Medir comportamento de cache, erro e rate limiting sem stress test.
- Registrar limites de custo operacional de sync e ingestão.
- Revisar Fan Content Policy, Terms of Service e developer agreements aplicáveis.
- Produzir contratos provisórios `PlayerSnapshot v0` e `CardCollection v0`.
- Encerrar com decisão de viabilidade e handoff explícito.

### Não funcionais

- Nenhum token/API key deve ser versionado.
- Probes devem respeitar rate limits e `Retry-After` quando houver.
- Nenhuma coleta deve usar scraping de cliente do jogo, engenharia reversa,
  interceptação de tráfego ou automação proibida.
- Evidências versionadas devem ser sanitizadas: evitar Player Tags, nomes,
  tokens, IPs ou payloads pessoais brutos quando não forem necessários.
- Toda afirmação de comportamento deve registrar fonte e data.
- Diferenciar claramente fato observado, documentação oficial e hipótese.

## Fora de escopo

- Reabrir as decisões já tomadas de Firebase Authentication, Cloud Firestore e
  Vercel como deploy inicial sem evidência concreta de incompatibilidade.
- Escolher framework web, ORM ou detalhes de modelagem física do Firestore.
- Construir login CrownPilot.
- Criar dashboard.
- Implementar pipeline de produção.
- Persistir dados reais de jogadores em banco de aplicação.
- Implementar recommendation engine.
- Definir pesos de Fit Score.
- Criar Upgrade Planner.
- Criar AI Coach.
- Implementar billing.
- Solicitar credenciais Supercell do jogador.
- Overlay, tracking em tempo real, bot, mod ou automação de gameplay.
- Scraping de sites de terceiros sem API/licença explícita.
- Consultoria jurídica conclusiva.

## Comportamento atual encontrado e baseline da main

Inspeção da `main` no commit `242f749ded8287629ad4635fdf1958664f1840d7` encontrou:

- `README.md`, `docs/roadmap/crownpilot-roadmap.md` e o contexto histórico;
- ADR aceito em `docs/decisions/001-firebase-firestore-vercel-portable.md`;
- esta spec, overview e oito subtarefas da Fase 001;
- nenhum `package.json`, lockfile, runtime, aplicação ou dependência;
- nenhum diretório `evidences/`, probe, fixture, CI, teste ou configuração de
  Firebase/Vercel versionado;
- Firebase Authentication + Cloud Firestore são decisões de infraestrutura;
- Vercel é o deploy inicial, com portabilidade obrigatória;
- não há `AGENTS.md`, regras locais, CI, testes ou framework web definido;
- não há credencial de API versionada;
- nenhuma capacidade da Fase 001 foi validada e o veredito ainda não existe.

Logo, esta spec deve corrigir o baseline documental sem transformar detalhes
abertos em arquitetura definitiva. Deve respeitar as decisões registradas em
`docs/decisions/001-firebase-firestore-vercel-portable.md`.

## Fontes e hierarquia de confiança

### Nível 1 — Autoritativas

- Portal oficial da Clash Royale API:
  https://developer.clashroyale.com/
- Fan Content Policy:
  https://supercell.com/en/fan-content-policy/
- Terms of Service:
  https://supercell.com/en/terms-of-service/
- demais developer agreements exibidos no portal autenticado.

### Nível 2 — Evidência direta

- respostas reais da API obtidas com token próprio;
- headers HTTP;
- códigos de erro;
- comportamento observado em probes controlados.

### Nível 3 — Secundárias

- SDKs e wrappers open source;
- documentação comunitária;
- RoyaleAPI e outros companions;
- discussões da comunidade.

Nível 3 pode gerar hipótese, nunca substituir validação quando a decisão afetar
contrato do produto.

## Perguntas críticas

### Vínculo e ownership

- A Player Tag pública é suficiente para o produto ou "minha conta" exige prova
  de ownership?
- Existe mecanismo oficial atual para verificar ownership/player token?
- Esse mecanismo possui scope/restrição que inviabiliza uso comum?
- Se não houver verificação, quais features precisam ser descritas como
  acompanhamento de perfil público em vez de propriedade confirmada?
- O MVP suporta uma Player Tag primária por usuário; múltiplas tags ficam fora
  do escopo inicial.

### Perfil e coleção

- `GET /players/{playerTag}` expõe a coleção completa atual?
- Nível retornado é suficiente para comparar readiness entre cartas de raridades
  diferentes ou precisa de normalização?
- Como o payload atual representa Evolution e Hero ownership?
- A mesma propriedade muda de semântica entre coleção, current deck e battle
  log?
- Arena e trophies representam o contexto que queremos ou existem múltiplas
  progressões/modos que exigem separação?
- Existe Collection Level / King Tower atual confiável e desde quando?
- Quais campos antigos continuam presentes mas não podem ser usados?

### Battle log

- Quantas partidas recentes são retornadas hoje?
- Existe paginação ou cursor?
- Quais modos aparecem?
- Deck, Evolution/Hero utilizado, resultado, trophy change e opponent estão
  disponíveis?
- O endpoint permite reconstruir histórico apenas por polling?
- Qual frequência mínima de polling evita lacunas para jogadores muito ativos?
- É permitido e operacionalmente sustentável armazenar snapshots para histórico
  próprio?

### Meta

- Rankings oficiais permitem descobrir jogadores suficientes por faixa?
- Como obter amostra representativa de arenas intermediárias?
- Clan/member discovery é adequado e permitido?
- Expansão via oponentes do battle log é tecnicamente e contratualmente aceitável?
- Qual viés cada estratégia introduz?
- Qual request cost para manter 7/14/30 dias de dados?
- Existe fonte third-party com API/licença explícita que seja melhor do que
  coletar internamente?
- Se não houver amostra confiável por Arena, qual fallback de produto é honesto?

### Operação

- Como tokens de API são emitidos e restritos atualmente?
- Existe allowlist de IP e como isso afeta serverless?
- Quais limites oficiais ou observados se aplicam?
- Quais endpoints retornam `Cache-Control` e por quanto tempo?
- Como 403, 404, 429, 5xx e manutenção se comportam?
- Quanto custa em requests sincronizar 1, 1k, 10k e 100k jogadores?

### Compliance

- Quais usos da API e assets são permitidos para fan content?
- Qual disclaimer precisa existir?
- CrownPilot pode usar nome/domínio sem marca Supercell?
- Quais assets podem ou não ser modificados?
- Um AI Coach pago pode ser classificado como coaching?
- Se isso não estiver claro, precisamos de aprovação expressa antes de billing?
- Que dados de usuário passam a ser nossos quando vinculamos login CrownPilot a
  Player Tag e qual política de privacidade será necessária depois?

## Baseline de infraestrutura

### Firebase Authentication

- Google é o provider inicial da conta CrownPilot.
- Auth do CrownPilot é independente da conta Supercell.
- Player Tag é um vínculo de domínio, não identidade de autenticação.
- Nenhuma credencial Supercell será armazenada.

### Cloud Firestore

Firestore será o banco principal para dados da aplicação, incluindo
progressivamente:

- identidade/vínculo do jogador;
- snapshots normalizados;
- histórico coletado quando permitido;
- dados derivados de recomendação;
- metadados de sync.

A Fase 001 **não** define ainda collections, índices, TTLs ou granularidade final.
Essas decisões dependem do volume e dos contratos v0 produzidos pelo discovery.

### Vercel com portabilidade

Vercel é a plataforma inicial de deploy, mas o core não pode depender de:

- Vercel KV;
- Vercel Postgres;
- Vercel Blob;
- Edge Config;
- Queues/Workflow;
- Cron;
- APIs de runtime proprietárias;

como requisito obrigatório do domínio.

Recursos Vercel podem ser usados no futuro como adapters operacionais quando
houver fallback ou boundary explícito.

Integrações externas devem ficar atrás de contratos próprios. Em particular, a
API da Supercell deve ser acessada por um `ClashRoyaleClient`/adapter equivalente
para que egress, hosting ou provider possam mudar sem reescrever domínio.

### Egress da Clash Royale API

Se a criação/uso de token oficial confirmar allowlist por IP, a Fase 001 deve
validar a consequência para Vercel. Em 29/09/2026, Vercel documenta Static IPs
para planos Pro+; isso é opção operacional, não contrato arquitetural.

Se custo ou restrição tornar Static IP inadequado, o acesso à API deve poder ser
movido para um egress service/gateway com IP estável sem alterar o restante da
aplicação.

## Abordagem escolhida

### Discovery orientado a evidência

A fase é executada em quatro trilhas que convergem:

1. **API surface** — documentação + probes;
2. **data semantics** — player, collection e battle log;
3. **meta acquisition** — cobertura, viés e custo;
4. **policy/compliance** — boundaries do produto e monetização.

Não haverá implementação de UI.

### Probes mínimos e reproduzíveis

Os experimentos devem poder ser repetidos sem uma stack definida.

Pode-se usar `curl`, shell ou um script descartável/local. Se um script
versionado trouxer valor real, ele deve ficar isolado em `scripts/discovery/`
e não ser tratado como fundação da stack futura.

Credencial esperada localmente:

`CLASH_ROYALE_API_TOKEN`

Nunca registrar o valor no repositório, logs versionados ou evidências.

### Amostragem de jogadores

Não confiar em um único perfil.

A validação deve usar uma amostra pequena e intencionalmente diversa, obtida por
fontes públicas permitidas, cobrindo quando possível:

- jogador em faixa intermediária;
- jogador avançado;
- perfil com Evolutions;
- perfil com Hero;
- perfil sem clan;
- diferentes modos/progressões.

Tags e nomes devem ser redigidos nas evidências versionadas quando não agregarem
valor técnico.

### Meta como prova separada

"Conseguimos ler meu perfil" não implica "conseguimos recomendar o melhor deck
da minha Arena".

A Fase 001 só aprova a hipótese `Best Decks for You` se existir uma estratégia
reproduzível para gerar candidatos com contexto competitivo e confiança
suficiente.

## Alternativas descartadas

- Implementar aplicação, autenticação ou recommendation engine antes de provar
  os dados: criaria contratos e custo de infraestrutura baseados em suposições.
- Tratar SDKs, wrappers ou sites de terceiros como autoridade: só documentação
  oficial, acordos aplicáveis e probes controlados podem fechar contratos.
- Usar ranking global como sinônimo de meta da Arena: pode introduzir viés de
  topo e não prova cobertura de faixas intermediárias.
- Escolher framework, schema físico do Firestore ou scheduler nesta fase:
  decisões dependem dos contratos e volumes observados.
- Fazer stress test para descobrir rate limit: risco desnecessário; limites não
  publicados permanecem desconhecidos e recebem margem conservadora.

## Arquivos, módulos e contratos afetados

- `docs/roadmap/crownpilot-roadmap.md`: fonte dos gates e do handoff; atualizar
  somente se o veredito mudar a sequência ou promessa do produto.
- `docs/specs/001-viabilidade-produto-dados-compliance.md`: contrato de execução
  desta fase e critérios de aceite.
- `docs/tasks/001-viabilidade-produto-dados-compliance/`: oito unidades de
  discovery e seus registros de execução.
- `docs/tasks/001-viabilidade-produto-dados-compliance/evidences/`: artefatos
  sanitizados produzidos durante execução; ainda inexistente no baseline.
- Contratos conceituais `PlayerSnapshotV0`, `CardCollectionEntryV0`,
  `CompetitiveContextV0` e `CurrentDeckV0`; não há módulos de produção nem
  símbolos implementados para alterar.
- ADR `docs/decisions/001-firebase-firestore-vercel-portable.md`: decisão aceita
  que limita mudanças de infraestrutura nesta fase.

## Artefatos esperados

A execução deve produzir, dentro de
`docs/tasks/001-viabilidade-produto-dados-compliance/evidences/`, artefatos
sanitizados equivalentes a:

- `api-surface.md` — endpoints relevantes, auth, parâmetros e status;
- `player-field-matrix.md` — campo → fonte → semântica → required/optional →
  freshness → risco;
- `battlelog-findings.md` — janela, modos e limitações;
- `meta-strategy-comparison.md` — estratégias, cobertura, viés, custo e
  compliance;
- `operational-findings.md` — cache, erros, rate behavior e estimativas;
- `compliance-findings.md` — políticas, links, datas e gates;
- `data-contract-v0.md` — PlayerSnapshot/CardCollection provisórios;
- `phase-001-verdict.md` — GO / GO WITH CONSTRAINTS / NO-GO + handoff.

Nomes podem variar se a estrutura final ficar mais clara, mas as informações
não podem desaparecer.

## Contrato provisório desejado

A fase deve terminar com algo conceitualmente próximo de:

```text
PlayerSnapshotV0
  identity
    playerTag
    playerName?
  competitiveContext
    arena
    trophies
    bestTrophies?
    rankedContext?
  progression
    collectionLevel?
    kingTowerLevel?
  collection[]
    cardId
    level
    count?
    evolutionOwned
    heroOwned
  currentDeck?
  source
    fetchedAt
    endpoint
    schemaObservedAt
```

Isso é objetivo de descoberta, não schema final.

Campos só entram se forem observados/justificados.

## Estratégia de validação

Não existe ainda toolchain de aplicação; portanto os gates desta fase são
documentais e experimentais.

Validação mínima:

1. conferir portal oficial e políticas atuais;
2. executar probes autenticados reais sem expor token;
3. repetir endpoints críticos em mais de um perfil;
4. registrar payload shape, optionality e headers relevantes;
5. comparar resultados com hipóteses secundárias;
6. calcular ordem de grandeza de requests;
7. revisar inconsistências antes de concluir;
8. realizar review final dos artefatos contra esta spec.

Não realizar stress/load test contra a API oficial.

## Riscos e mitigação

| Risco | Mitigação |
| --- | --- |
| API não expõe dados necessários | Classificar como indisponível e ajustar MVP antes da Fase 002. |
| Campos mudaram após updates de 2026 | Priorizar live probes e registrar data de observação. |
| Evolution/Hero têm semântica contextual | Validar coleção, current deck e battle log separadamente. |
| Battle log é curto | Medir janela e estimar polling necessário; não prometer histórico antes disso. |
| Meta de Arena não é amostrável | Tratar como gate; considerar outra segmentação ou fonte licenciada. |
| Rate limit inviabiliza crawler | Medir custo, usar cache e separar sync de usuário de ingestão global. |
| Serverless incompatível com IP allowlist | Registrar como restrição de infraestrutura para fase posterior. |
| Fonte third-party vira dependência opaca | Exigir API/licença/termos claros e plano de fallback. |
| Monetização conflita com Fan Content Policy | Billing permanece bloqueado até interpretação/aprovação suficiente. |
| Evidências expõem dados desnecessários | Sanitizar tags, nomes, tokens, IPs e payloads antes de commit. |
| Discovery escolhe stack sem necessidade | Scripts são descartáveis; decisões de stack ficam para fase apropriada. |

## Critérios de aceite

- [ ] superfície relevante da API está mapeada com fonte e data;
- [ ] autenticação/token/IP requirements foram confirmados;
- [ ] Player Tag válida foi resolvida em probe real;
- [ ] vínculo de perfil público vs. ownership verificado possui decisão
  documentada e baseada em capacidade oficial observada;
- [ ] coleção e níveis foram observados em múltiplos perfis;
- [ ] Evolution/Hero ownership e deployment foram diferenciados ou marcados
  explicitamente como ainda desconhecidos;
- [ ] Arena/trophies e demais contextos de progressão foram classificados;
- [ ] battle log foi medido quanto a janela, modos, paginação e campos;
- [ ] existe estratégia de meta avaliada para faixas intermediárias;
- [ ] viés e request cost dessa estratégia estão documentados;
- [ ] cache/rate/error behavior foi medido sem stress test;
- [ ] Fan Content Policy e agreements aplicáveis foram revisados;
- [ ] monetização possui status explícito: allowed, constrained ou blocked
  pending approval;
- [ ] `PlayerSnapshot v0` e `CardCollection v0` estão documentados;
- [ ] dados necessários ao MVP estão classificados como available, derived,
  optional, unavailable ou unresolved;
- [ ] existe veredito GO / GO WITH CONSTRAINTS / NO-GO;
- [ ] handoff descreve exatamente o que a Fase 002 pode assumir.

## Ordem das subtarefas

1. [001-01-mapear-api-oficial-e-auth.md](../tasks/001-viabilidade-produto-dados-compliance/001-01-mapear-api-oficial-e-auth.md)
2. [001-02-validar-player-profile-e-collection.md](../tasks/001-viabilidade-produto-dados-compliance/001-02-validar-player-profile-e-collection.md)
3. [001-03-validar-battlelog-e-historico.md](../tasks/001-viabilidade-produto-dados-compliance/001-03-validar-battlelog-e-historico.md)
4. [001-04-validar-aquisicao-do-meta.md](../tasks/001-viabilidade-produto-dados-compliance/001-04-validar-aquisicao-do-meta.md)
5. [001-05-medir-operacao-cache-e-rate-limits.md](../tasks/001-viabilidade-produto-dados-compliance/001-05-medir-operacao-cache-e-rate-limits.md)
6. [001-06-validar-compliance-e-monetizacao.md](../tasks/001-viabilidade-produto-dados-compliance/001-06-validar-compliance-e-monetizacao.md)
7. [001-07-definir-contratos-de-dados-v0.md](../tasks/001-viabilidade-produto-dados-compliance/001-07-definir-contratos-de-dados-v0.md)
8. [001-08-fechar-gates-e-handoff.md](../tasks/001-viabilidade-produto-dados-compliance/001-08-fechar-gates-e-handoff.md)

## Premissas explícitas

- O portal oficial e live API podem mudar; data de observação é parte do
  contrato de discovery.
- A API key será criada/configurada manualmente pelo responsável e nunca
  compartilhada em documentação.
- Não assumimos que `expLevel`, `collectionLevel`, `kingTowerLevel`,
  `evolutionLevel` ou qualquer outro campo histórico mantenha semântica antiga.
- Não assumimos que ranking global represente Arena intermediária.
- Não assumimos que a exceção de "coaching" da Fan Content Policy autorize
  automaticamente AI Coach pago.
- Firebase Authentication, Cloud Firestore e Vercel como deploy inicial já são
  decisões; framework web, modelagem física e detalhes do runtime continuam
  abertos nesta fase.
