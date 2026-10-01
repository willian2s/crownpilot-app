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

Um proxy third-party pode ser usado como transporte operacional enquanto o ambiente
não possui IP fixo, por decisão explícita e reversível do projeto. Deve ser
documentado separadamente da rota oficial. Proxy operacional não prova endorsement
da Supercell, contrato oficial, ownership ou autorização comercial. Termos,
privacidade, retenção, tratamento de API key, limites e compatibilidade contratual
continuam gates próprios e débitos rastreados; uso do proxy não libera browser,
billing ou qualquer modelo comercial.

No caso do RoyaleAPI Proxy selecionado para o momento, a integração é substituição
de host (`api.clashroyale.com` → `proxy.royaleapi.dev`) mantendo paths `/v1/...`,
query e autenticação Bearer nos GETs observados. A documentação oficial da
Supercell continua autoridade para endpoints e contratos; o proxy não altera o
significado desses endpoints.

A Fan Content Policy pública consultada em 29/09/2026 informa "Last updated:
September 27, 2023" e estabelece limites relevantes para uso comercial. Isso
torna monetização um gate explícito, principalmente para qualquer plano Pro ou
feature paga baseada em IA.

Atualização documental em 01/10/2026: a leitura da política oficial confirmou que
Fan Content é não comercial por padrão e que não é permitido cobrar taxa de
qualquer tipo, inclusive por funcionalidades in-app, sem aprovação expressa da
Supercell. Ads, donations e coaching aparecem como exceções específicas, mas com
condições próprias. Essa exceção não deve ser interpretada como autorização
automática para software coaching ou AI Coach. A decisão abaixo é um gate de
discovery, não parecer jurídico, e deve ser confirmada contra os acordos de
developer/API aplicáveis.

## Objetivo

Eliminar as principais incógnitas que podem inviabilizar ou alterar o produto e
encerrar a fase com uma decisão objetiva:

- **GO:** dados e permissões sustentam o MVP planejado;
- **GO WITH CONSTRAINTS:** MVP é viável com limitações explícitas e roadmap
  ajustado;
- **GO WITH CONSTRAINTS / APPROVAL DEPENDENCY:** dados, operação e core técnico
  são viáveis, mas sustentabilidade depende de modelo que exige aprovação
  expressa; billing permanece bloqueado até essa aprovação;
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
- Classificar cada modelo comercial com `ALLOWED`, `ALLOWED WITH CONDITIONS`,
  `REQUIRES EXPLICIT APPROVAL`, `NOT ALLOWED` ou `UNRESOLVED`, sempre com fonte
  oficial e data de consulta.
- Separar acesso técnico à API, uso permitido de dados/assets e autorização
  comercial; uma API key ou resposta 2xx não autoriza cobrança.
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

## Monetização e Fan Content Policy

### Fontes oficiais e regra de interpretação

Fontes consultadas em `2026-10-01`:

| Fonte | Versão/data visível | Uso nesta spec |
| --- | --- | --- |
| [Supercell Fan Content Policy](https://supercell.com/en/fan-content-policy/) | `Last updated: September 27, 2023` | Regra de Fan Content, cobrança, exceções de ads/donations/coaching, disclaimer, assets, marca e domínios. |
| [Supercell Terms of Service](https://supercell.com/en/terms-of-service/) | `Effective Date: November 6, 2024` | Termos incorporados, propriedade intelectual, uso não comercial do Service e possibilidade de suspensão/alteração. |
| [Clash Royale API developer portal](https://developer.clashroyale.com/) | consulta `2026-10-01` | Acesso técnico, API key, IP/egress e acordos exibidos no portal; não substitui autorização comercial. |

A Fan Content Policy é a fonte principal para o enquadramento do conteúdo. A
política afirma que guias e guide apps online não comerciais podem ser Fan Content
quando respeitam suas regras. Também afirma que não se pode cobrar fee de qualquer
tipo, incluindo funcionalidades in-app, salvo aprovação expressa da Supercell, e
lista ads, donations e coaching como exceções geralmente permitidas. “Geralmente”
não elimina as condições da política nem resolve o enquadramento de um produto
automatizado.

`ALLOWED` significa permitido diretamente pela fonte, sem condição material
adicional identificada nesta fase. `ALLOWED WITH CONDITIONS` significa que a fonte
descreve a possibilidade, mas impõe limites explícitos. `REQUIRES EXPLICIT
APPROVAL` significa que cobrança ou uso não pode avançar com base apenas na
política publicada. `NOT ALLOWED` significa vedação publicada. `UNRESOLVED`
significa que a fonte não responde suficientemente ao caso concreto; enquanto
permanecer assim, o produto não pode tratá-lo como permitido.

### Classificação documental inicial para o CrownPilot

Esta matriz registra a leitura conservadora da política publicada em `2026-10-01`;
não encerra `001-06` nem substitui a revisão dos developer/API agreements. A
classificação final deve ser registrada em `evidences/compliance-findings.md`.

| Pergunta/modelo | Status | Conclusão baseada na fonte oficial | Consequência para o produto |
| --- | --- | --- | --- |
| Fan Content permitido | `ALLOWED WITH CONDITIONS` | Guias e guide apps online não comerciais são exemplos de Fan Content permitido quando exibem, identificam ou discutem produtos Supercell e cumprem toda a política. | O core pode ser desenhado como conteúdo/análise não comercial, sujeito a assets, marca, disclaimer, dados e demais termos. |
| Conteúdo gratuito | `ALLOWED WITH CONDITIONS` | Gratuito não remove limites de Fan Content, assets, marca, conteúdo proibido ou developer policies. | Gratuidade não substitui revisão de compliance. |
| Cobrança direta por funcionalidades | `REQUIRES EXPLICIT APPROVAL` | A política veda cobrar fee de qualquer tipo, inclusive por funcionalidades in-app, salvo aprovação expressa. | Assinatura, paywall e venda de features não entram como capacidade liberada. |
| Ads | `ALLOWED WITH CONDITIONS` | Ads são exceção geralmente permitida, sujeitos a leis, regras, developer policies e proibição de sugerir patrocínio/endorsement da Supercell. | Avaliar provedor, criativos, jurisdição e integração individualmente. |
| Donations | `ALLOWED WITH CONDITIONS` | Donations devem ser puramente donations e não podem estar ligadas a features especiais, IAPs ou qualquer benefício. | Doação não pode desbloquear AI, analytics, limites ou prioridade. |
| Coaching humano | `ALLOWED WITH CONDITIONS` | FAQ define coaching como training/guidance e cita personal/online coaching e venda de base layouts; Supercell reserva a definição do enquadramento. | Exige definição operacional, conteúdo compatível e revisão dos termos aplicáveis. |
| Software coaching | `REQUIRES EXPLICIT APPROVAL` | A política não afirma que software automatizado é coaching permitido; “similar activities” e a reserva de interpretação não resolvem o caso. | Não cobrar nem posicionar como exceção sem esclarecimento/aprovação escrita. |
| AI Coach gratuito | `UNRESOLVED` | A ausência de cobrança remove a questão da taxa, mas a política não confirma que coaching automatizado por software/AI seja Fan Content permitido. | Não tratar como automaticamente permitido; validar enquadramento de produto, dados e assets separadamente. |
| AI Coach pago | `REQUIRES EXPLICIT APPROVAL` | A exceção de coaching não cobre automaticamente produto de software/AI pago. | AI Coach pago pode existir como hipótese técnica futura, mas monetização e enquadramento ficam bloqueados. |
| Assinatura/paywall | `REQUIRES EXPLICIT APPROVAL` | É cobrança por acesso/funcionalidade e cai na regra de fee de qualquer tipo. | Não assumir plano Pro; não implementar billing como premissa. |
| Venda de funcionalidades/premium features | `REQUIRES EXPLICIT APPROVAL` | A vedação menciona expressamente in-app functionalities; não há autorização geral para recursos premium. | Advanced Analytics, AI Coach ou limites pagos permanecem bloqueados. |
| SaaS pago | `REQUIRES EXPLICIT APPROVAL` | A política não cria exceção geral para SaaS; o produto proposto continua sendo avaliado como Fan Content quando usa assets/marca/conteúdo Supercell. | Não tratar modelo SaaS tradicional como permitido por padrão. |
| Sponsorship | `ALLOWED WITH CONDITIONS` | A política permite ads/promoções sob condições, mas proíbe criar impressão de que Supercell patrocina, cria ou endossa o Fan Content. | Patrocínio deve ser de terceiros e revisado; uso de marca Supercell em material comercial exige cuidado adicional/aprovação quando aplicável. |
| Assets | `ALLOWED WITH CONDITIONS` | Assets devem ser usados para exibir, identificar e discutir produtos; não podem ser modificados sem permissão expressa nem usados para imitar logos/trademarks. | Inventariar assets, evitar alterações e bloquear material que pareça oficial. |
| Assets modificados | `REQUIRES EXPLICIT APPROVAL` | A política não permite modificar Supercell Assets sem permissão expressa. | Não alterar assets sem aprovação rastreável. |
| Nome, branding e domínio sem trademark Supercell | `ALLOWED WITH CONDITIONS` | Não criar impressão de endorsement e não imitar logos, trademarks ou elementos dos produtos. | Revisar `CrownPilot`, title, logos e copy antes de publicação. |
| Domínio, conta social ou endereço com trademark Supercell/nome de jogo | `REQUIRES EXPLICIT APPROVAL` | A política exige acordo escrito separado para esses endereços. | Não registrar ou publicar sem acordo escrito. |
| Disclaimer | `ALLOWED WITH CONDITIONS` | Deve existir aviso de não oficialidade, ou substancialmente similar, legível e em conexão com o Fan Content. | Preservar no produto: “This material is unofficial and is not endorsed by Supercell. For more information see Supercell's Fan Content Policy: www.supercell.com/fan-content-policy.” |
| Bots, mods, automação, private servers e software não autorizado | `NOT ALLOWED` | A política proíbe Fan Content que promova ou contenha esses usos. | Permanecem fora do produto e do marketing. |

### Diferença entre acesso à API e autorização comercial

São gates independentes:

1. **API access:** credencial, IP/egress, endpoints, limites e contratos do
   developer portal. A existência de API key, endpoint documentado ou resposta
   `2xx` prova no máximo acesso técnico autorizado naquele escopo.
2. **Fan Content/data use:** permissão para exibir, identificar, discutir,
   armazenar ou redistribuir conteúdo/dados conforme política e acordos
   aplicáveis.
3. **Commercial authorization:** permissão para cobrar ou explorar um modelo
   comercial específico. Não é concedida automaticamente pelo acesso à API, pelo
   uso do disclaimer ou pela exceção de coaching.

Enquanto developer/API agreements não forem revisados na sessão aplicável, a
classificação comercial não pode ser elevada por inferência. Qualquer modelo que
dependa de aprovação expressa permanece bloqueado até existir registro rastreável
da aprovação.

### Efeito no GO/NO-GO

A impossibilidade de cobrar assinatura, por si só, não determina `NO-GO`. O
veredito da Fase 001 deve avaliar separadamente:

1. viabilidade técnica;
2. disponibilidade e qualidade dos dados;
3. viabilidade operacional;
4. compliance e limites de uso;
5. modelo sustentável permitido ou estratégia explícita para operar sem
   monetização inicialmente.

Se dados e operação forem viáveis, mas sustentabilidade depender de modelo que
exige aprovação da Supercell, o resultado correto é
`GO WITH CONSTRAINTS / APPROVAL DEPENDENCY`. Isso bloqueia a receita dependente de
aprovação sem bloquear automaticamente a investigação ou um MVP gratuito. A Fase
012 deve consumir esta decisão e nunca funcionar como autorização implícita.

## Baseline de infraestrutura

### Firebase Authentication

- Google é o provider inicial da conta CrownPilot.
- Auth do CrownPilot é independente da conta Supercell.
- Player Tag é um vínculo de domínio, não identidade de autenticação.
- Nenhuma credencial Supercell será armazenada.

### Associação de perfil público sem ownership

O MVP pode permitir que usuário CrownPilot salve uma Player Tag e associe o
perfil público resolvido à própria conta CrownPilot, sem afirmar que usuário é
proprietário da conta Clash Royale.

Invariantes:

- associação pertence ao `crownpilotUserId` derivado server-side;
- domínio usa `subjectType: public_profile` e `ownershipStatus: unverified`;
- UI deve dizer **“Perfil público salvo — ownership não verificado”**;
- perfil não concede ações, notificações, exclusividade, impersonation ou
  autorização sobre conta Clash Royale;
- mesma Player Tag pode aparecer associada a mais de um usuário CrownPilot;
- associação, snapshot e derivados ficam privados por padrão, com remoção e
  refresh controlados;
- `not_found`, `provider_unavailable`, `rate_limited` e `stale` são estados
  distintos de ownership;
- nenhum endpoint `verifytoken` é assumido; ownership permanece `UNRESOLVED`.

Essa decisão permite `GO WITH CONSTRAINTS` para lookup read-only de perfil público,
mas não fecha contratos de dados, privacidade, rate limits ou uso de proxy. Esses
itens continuam dependências de `001-02`, `001-05`, `001-06` e `001-07`.

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
- Tratar acesso técnico à API como autorização comercial: são decisões e fontes
  distintas; a primeira não prova a segunda.
- Interpretar a exceção de coaching como autorização automática para software,
  AI Coach, assinatura ou venda de funcionalidades: a política não resolve esse
  enquadramento e exige gate conservador.
- Tratar ads ou donations como permissão sem condições: ambas possuem limites
  próprios, e donations não podem desbloquear benefícios.

## Arquivos, módulos e contratos afetados

- `docs/roadmap/crownpilot-roadmap.md`: fonte dos gates e do handoff; atualizar
  somente se o veredito mudar a sequência ou promessa do produto.
- `docs/specs/001-viabilidade-produto-dados-compliance.md`: contrato de execução
  desta fase e critérios de aceite.
- `docs/tasks/001-viabilidade-produto-dados-compliance/`: oito unidades de
  discovery e seus registros de execução.
- `docs/tasks/001-viabilidade-produto-dados-compliance/001-01-mapear-api-oficial-e-auth.md`:
  relação entre acesso à API e autorização comercial, sem alterar o bloqueio
  existente.
- `docs/tasks/001-viabilidade-produto-dados-compliance/001-06-validar-compliance-e-monetizacao.md`:
  matriz de modelos comerciais, fontes oficiais e gate de aprovação.
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
| Exceção de coaching é aplicada a AI/software sem base | Classificar como `REQUIRES EXPLICIT APPROVAL` ou `UNRESOLVED`; não cobrar nem prometer a feature como permitida. |
| API access é confundido com autorização comercial | Manter gates separados e exigir revisão de Fan Content Policy e developer/API agreements. |
| Ads/donations são tratados como receita sem condições | Validar individualmente; donations nunca desbloqueiam benefícios. |
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
- [ ] monetização possui status explícito `ALLOWED`, `ALLOWED WITH CONDITIONS`,
  `REQUIRES EXPLICIT APPROVAL`, `NOT ALLOWED` ou `UNRESOLVED` por modelo;
- [ ] assinatura/paywall e cobrança por funcionalidades não são tratados como
  permitidos sem aprovação expressa;
- [ ] ads, donations, coaching humano, software coaching, AI Coach, sponsorship
  e SaaS foram analisados individualmente;
- [ ] o veredito separa viabilidade técnica, dados, operação, compliance e
  sustentabilidade, permitindo `GO WITH CONSTRAINTS / APPROVAL DEPENDENCY`;
- [ ] arquitetura e MVP possuem caminho sem monetização ainda não aprovada;
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
- Fan Content Policy consultada em 2026-10-01 continua sendo referência sujeita a
  mudança; a classificação não é parecer jurídico.
- API access, Fan Content/data use e commercial authorization são gates
  independentes.
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
