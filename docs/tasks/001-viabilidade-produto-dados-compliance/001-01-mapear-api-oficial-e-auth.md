# 001-01 — Mapear API oficial e autenticação

- **Ticker:** `001`
- **Número:** `01`
- **Status:** `completed`

## Requisitos cobertos

- superfície oficial, autenticação, encoding de Player Tag e erros;
- requisito de IP/egress e existência de verificação oficial de ownership;
- fontes necessárias para perfil, battle log, catálogo e descoberta de meta.

## Objetivo e resultado esperado

Estabelecer a superfície oficial da API necessária ao MVP e provar que conseguimos
fazer chamadas autenticadas reproduzíveis sem versionar segredo.

Ao final deve existir um mapa de endpoints, auth, encoding, erros e fontes
autoritativas suficiente para executar as próximas tasks.

## Escopo incluído

- revisar portal `developer.clashroyale.com`;
- registrar processo atual de criação/uso de API key sem registrar token;
- confirmar header de autorização e restrições de origem/IP;
- confirmar encoding de Player Tag;
- mapear endpoints de:
  - players;
  - mecanismo oficial de player/ownership verification, se disponível;
  - battle log;
  - cards;
  - locations/rankings;
  - leaderboards;
  - clans/members se relevantes para discovery;
- registrar response shape/paginação declarados;
- executar um probe mínimo de sucesso e probes seguros de erro;
- criar `evidences/api-surface.md`.

## Escopo excluído

- crawler;
- banco;
- SDK próprio;
- wrapper genérico;
- stress test;
- implementação web.

## Dependências

- conta de developer/API key criada pelo responsável;
- acesso de rede a `api.clashroyale.com`.

## Arquivos e símbolos prováveis

- `evidences/api-surface.md`;
- documentação oficial, endpoints `/v1/players`, `/v1/players/{tag}/battlelog`,
  `/v1/cards`, `/v1/locations` e superfícies de clans/rankings quando aplicável;
- variável local `CLASH_ROYALE_API_TOKEN`; nenhum símbolo de produção existe.

## Passos de execução

1. Consultar documentação oficial e registrar data.
2. Criar/configurar token fora do repositório.
3. Confirmar variável local `CLASH_ROYALE_API_TOKEN`.
4. Executar chamada autenticada mínima.
5. Testar uma Player Tag válida com `#` corretamente encoded.
6. Registrar status/shape de endpoints candidatos sem coletar massa de dados.
7. Confirmar paginação apenas onde aplicável.
8. Confirmar se existe endpoint/scope oficial para verificação de ownership e
   quais pré-requisitos ele possui; não solicitar token de jogador nesta task se
   não for necessário para mapear o contrato.
9. Registrar códigos de erro seguros: tag inexistente, token ausente e parâmetro
   inválido quando isso não gerar carga indevida.
10. Sanitizar qualquer evidência antes de commit.

## Evidências obrigatórias

`evidences/api-surface.md` deve conter:

- endpoint;
- finalidade;
- auth;
- parâmetros;
- paginação;
- response shape;
- source URL;
- observed/documented;
- observed at;
- status para MVP: required / useful / out.

## Validação

- ao menos uma chamada 2xx real;
- ao menos um caso de tag encoding confirmado;
- nenhum token em arquivos ou output versionado;
- todos os endpoints críticos classificados.

## Testes e comandos de validação

- `test -n "$CLASH_ROYALE_API_TOKEN"` sem imprimir o valor;
- executar probes `curl` autenticados usando a variável de ambiente, com headers
  salvos apenas em arquivo temporário sanitizado;
- `git diff --check` e busca por segredo antes de versionar evidências;
- revisar manualmente cada endpoint contra fonte oficial e data de consulta.

## Definição de pronto

- auth e superfície mínima deixaram de ser suposição;
- task 001-02 pode executar probes em perfis reais;
- limitações ainda desconhecidas estão explicitamente marcadas.

## Riscos e cuidados

- portal oficial é dinâmico e pode não ser facilmente capturável; registrar
  manualmente o contrato visto no portal quando necessário;
- IP allowlist pode afetar ambiente local/cloud — observar sem decidir infra;
- não confiar em SDK antigo como autoridade.

## Registro de execução

### Execução inicial em `2026-10-01`

- **Status final:** `blocked` — `CLASH_ROYALE_API_TOKEN` ausente; dependência declarada pela task não está disponível.
- **Arquivos alterados:** `evidences/api-surface.md`; este registro; overview somente para registrar bloqueio.
- **Decisões:** classificar endpoint categories como candidate/pending, sem transformar shape não observado em contrato; tratar Player Tag como vínculo de perfil público até ownership ser confirmado; não versionar tag, nome, IP ou token.
- **Fontes consultadas:** portal oficial `https://developer.clashroyale.com/`; documentação `https://developer.clashroyale.com/api-docs/index.html`; API `https://api.clashroyale.com/v1/...`, consultados em `2026-10-01`.
- **Probes executados:** `test -n "$CLASH_ROYALE_API_TOKEN"`; requests sem Authorization para `/v1/cards`, `/v1/players/%23INVALIDTAG` e `/v1/players/not-a-player-tag`; headers/body mantidos somente em diretório temporário sanitizado.
- **Resultados:** API alcançável; três probes retornaram `403 accessDenied / Missing authorization`; JSON UTF-8 e `Cache-Control: public max-age=600` observados; encoding `%23` foi enviado, mas resolução positiva não pode ser confirmada sem auth; nenhum `2xx`.
- **Desvios:** não foi possível executar chamada autenticada mínima, confirmar header por sucesso, validar Player Tag, separar 404/invalid tag, confirmar paginação/shape de campos ou testar ownership. Não houve stress/rate-limit probe.
- **Revisão independente:** solicitada após implementação; confirmou status `blocked`, checklist overview `1 seção/8 itens`, ausência de segredo e necessidade de manter shapes/paginação como `pending` até sessão autenticada.
- **Riscos residuais:** token/IP allowlist podem bloquear local ou Vercel; contrato Swagger atual e endpoint de ownership continuam pendentes; `001-02` e `001-03` não devem iniciar probes reais até remover o bloqueio.

### Complemento documental: Fan Content Policy e autorização comercial

Este complemento não altera o status `blocked`, os probes, resultados ou
decisões registradas acima. Ele adiciona uma dependência de compliance descoberta
durante a revisão cruzada da task.

- **Fontes oficiais relevantes:** [Fan Content Policy](https://supercell.com/en/fan-content-policy/), consultada em `2026-10-01` e indicada pela página como atualizada em `2023-09-27`; [Terms of Service](https://supercell.com/en/terms-of-service/), com data efetiva indicada como `2024-11-06`; [Clash Royale API developer portal](https://developer.clashroyale.com/).
- **Regra comercial relevante:** Fan Content é não comercial por padrão; cobrança de qualquer taxa, inclusive por funcionalidades in-app, exige aprovação expressa da Supercell. Ads, donations e coaching são exceções condicionais, não autorização geral para assinatura, premium features, software coaching ou AI Coach.
- **Relação com esta task:** API key, autenticação, IP/egress, endpoint documentado ou eventual resposta `2xx` provam acesso técnico ao escopo da API, não autorização comercial para cobrar ou explorar um modelo específico.
- **Dependência/gate:** `001-06-validar-compliance-e-monetizacao.md` deve fechar a classificação oficial de assinatura, paywall, premium features, analytics pago, AI Coach, ads, donations, coaching, software/AI coaching, sponsorship e SaaS. Até lá, billing e modelos que dependam de aprovação permanecem bloqueados.
- **Ownership versus comercialização:** nenhum endpoint oficial de ownership foi
  confirmado na Swagger disponível; não assumir `verifytoken` nem solicitar
  player token. Perfil público, eventual ownership verificado e autorização de
  monetização continuam gates distintos.
- **Limite de escopo:** esta task registra a relação entre API access e policy; não executa a análise comercial completa nem altera a execução pendente de `001-06`.

### Revisão documental adicional em `2026-10-01`

- **Status:** permanece `blocked`; nenhuma evidência autenticada nova foi
  produzida e nenhum status de execução foi alterado.
- **Arquivos relacionados atualizados:** spec da Fase 001, roadmap, overview,
  task `001-06`, task `001-08` e contexto histórico; este arquivo recebeu
  somente o complemento acima.
- **Decisão:** manter três gates separados: API access; Fan Content/data use;
  commercial authorization. Acesso técnico à API não autoriza assinatura,
  paywall, venda de funcionalidades, software coaching ou AI Coach.
- **Fontes/evidência:** Fan Content Policy, Terms of Service e portal oficial da
  API consultados em `2026-10-01`; policy indica atualização em `2023-09-27` e
  Terms indica vigência em `2024-11-06`. A evidência autenticada continua sendo a
  registrada em `evidences/api-surface.md`.
- **Comandos/checks:** `git diff --check`; busca cruzada de referências de
  monetização e revisão dos links oficiais. Nenhum probe autenticado foi
  repetido nessa tentativa documental, pois `CLASH_ROYALE_API_TOKEN` ainda não
  estava disponível no processo.
- **Desvios:** análise comercial detalhada permanece em `001-06`; não foi
  incorporada como execução de API nesta task.
- **Riscos residuais:** agreements autenticados, escopo de uso de dados,
  ownership e resposta positiva da API continuam pendentes; modelos dependentes
  de aprovação permanecem bloqueados.
- **Revisão independente adicional:** aprovada com ressalva documental
  preexistente em `001-08`; confirmou preservação do bloqueio, histórico,
  evidências e ausência de autorização comercial implícita.

### Tentativa de retomada em `2026-10-01`

- **Status:** permanece `blocked`.
- **Comando:** `test -n "$CLASH_ROYALE_API_TOKEN"`.
- **Resultado:** variável não disponível no processo de execução desta sessão;
  nenhum probe autenticado foi executado e nenhum segredo foi impresso.
- **Próximo passo:** disponibilizar a variável exportada no ambiente do agente,
  fora do repositório, e repetir primeiro `/v1/cards` e depois os probes de
  Player Tag. O bloqueio anterior e todas as evidências permanecem válidos.

### Retomada com `.env.local` em `2026-10-01`

- **Status:** permanece `blocked` por allowlist de IP/egress.
- **Configuração:** `.env.local` foi carregado somente no processo do probe;
  `CLASH_ROYALE_API_TOKEN` estava presente e seu valor não foi impresso,
  versionado ou registrado.
- **Probes executados:** `GET /v1/cards` e
  `GET /v1/players/%23INVALIDTAG`, ambos com `Authorization: Bearer` e resposta
  mantida em arquivos temporários sanitizados.
- **Resultados:** ambos retornaram `HTTP 403`, JSON com `reason=accessDenied.invalidIp`
  e `message=Invalid authorization: API key does not allow access from current
  egress`; `Content-Type: application/json` e `Cache-Control: public max-age=600`
  foram observados. Nenhum `2xx` foi obtido.
- **Decisão:** autenticação chegou à API, mas o egress desta execução não está
  autorizado pela chave. Não classificar isso como token inválido nem como falha
  de encoding; ainda não há prova de `Player Tag` válida, ownership, shape ou
  paginação.
- **Próximo passo:** adicionar o egress permitido à configuração da API key ou
  executar probes a partir de ambiente com IP já permitido; depois repetir
  `/v1/cards`, uma Player Tag permitida e o battle log. Não registrar o IP
  público em evidência versionada.
- **Segurança:** `.env.example` foi sanitizado para conter apenas placeholder e
  `.env.local` passou a ser ignorado pelo Git. A credencial que apareceu fora do
  arquivo local deve ser revogada/rotacionada antes de qualquer novo probe; seu
  valor não foi registrado em task ou evidência.

### Retomada via RoyaleAPI Proxy em `2026-10-01`

- **Status:** `in_progress`; rota proxy autenticada funciona, mas contrato
  third-party, termos e egress direto continuam pendentes.
- **Configuração:** nova chave informada pelo responsável foi carregada de
  `.env.local`; nenhum valor foi impresso ou registrado.
- **Comandos/probes:** `curl` autenticado para `https://proxy.royaleapi.dev/v1/cards`,
  `https://proxy.royaleapi.dev/v1/players/%23<REDACTED_VALID_TAG>`,
  `https://proxy.royaleapi.dev/v1/players/%23<REDACTED_VALID_TAG>/battlelog` e
  `https://proxy.royaleapi.dev/v1/players/%23INVALIDTAG`, com respostas mantidas
  somente em temporários sanitizados.
- **Resultados:** `cards=200`, perfil=`200`, battlelog=`200`, tag inválida=`404`.
  Cards retornou objeto com `items=123` e `supportItems=4`; perfil retornou objeto
  com `cards=123` e campos de progressão/deck; battlelog retornou array com 30
  itens. `Content-Type` JSON foi observado; `Cache-Control` foi `max-age=49`
  para cards, `max-age=60` para perfil/battlelog e `public max-age=600` para
  erro de tag inválida.
- **Decisão:** encoding `%23` e header de autorização estão confirmados na rota
  proxy para esses endpoints; isso prova acesso operacional via proxy, não
  contrato oficial direto, ownership verificado ou autorização comercial.
- **Desvios:** não foi possível provar rota direta por causa da allowlist; o
  proxy permanece dependência third-party até revisão de termos, retenção,
  limites, SLA e compatibilidade com agreements da Supercell.
- **Próximo passo:** decidir proxy validado versus egress próprio, repetir probes
  a partir da rota escolhida e classificar endpoints críticos antes de liberar
  `001-02`/`001-03`.
- **Segurança:** revogar/rotacionar a chave usada nos probes exploratórios antes
  de qualquer nova chamada; revisar contrato, privacidade, retenção, logs,
  limites, SLA e compatibilidade com a Supercell antes de aceitar proxy.
- **Revisão independente:** confirmou shapes e status dos probes, mas não aprovou
  encerramento: proxy continua pending, tag foi sanitizada na evidência e a chave
  usada deve ser rotacionada antes de novos requests.

### Proxy de egress candidato

- **Lead recebido:** [RoyaleAPI Proxy](https://docs.royaleapi.com/proxy.html),
  indicado como alternativa para ambientes sem IP fixo.
- **Estado da evidência:** a URL retornou `403`/Cloudflare para esta execução;
  conteúdo, termos, tratamento da API key, retenção, limites e SLA não foram
  verificados. A referência permanece hipótese de infraestrutura, não contrato.
- **Gate original:** não aceitar proxy como dependência antes de revisar termos
  da RoyaleAPI, política de privacidade, retenção/logs, escopo de uso, limites,
  disponibilidade e compatibilidade com os developer/API agreements da
  Supercell.
- **Arquitetura:** se validado, proxy deve ser server-side e tratado como
  dependência third-party explícita, com fallback e observabilidade. Nunca expor
  API key no browser. API access via proxy continua separado de autorização
  comercial e de Fan Content Policy.
- **Desvio registrado:** probes exploratórios foram executados antes dessa
  revisão de confiança, a pedido do responsável, somente para medir conectividade
  e shape. Decisão posterior selecionou proxy como transporte operacional
  temporário, mantendo esses itens como débito e risco.
- **Alternativa:** avaliar egress próprio com IP fixo (NAT/gateway/VPS/serviço
  equivalente) caso proxy não ofereça garantias suficientes.

### Decisão provisória: RoyaleAPI Proxy

- **Decisão do responsável:** utilizar RoyaleAPI Proxy como transporte atual do
  CrownPilot até decisão explícita de troca ou disponibilidade de servidor com IP
  fixo, devido à ausência de IP fixo no ambiente.
- **Evidência operacional:** após rotação da chave, cards, locations, perfil e
  battlelog retornaram `200` via `https://proxy.royaleapi.dev`; tag inválida
  retornou `404`. Locations retornou objeto com `items=262` e `paging`.
- **Fontes third-party recebidas:** RoyaleAPI Terms of Service e Privacy Policy,
  ambas `Last Updated: 2024-10-11`. Privacy Policy declara coleta de IP,
  geolocalização, horários, URLs, clickstream e dados de uso, além de
  compartilhamento com service providers. O texto fornecido não especifica
  tratamento da API key pelo endpoint proxy.
- **Classificação operacional:** `CURRENT OPERATIONAL DEPENDENCY — CONDITIONAL —
  NOT AUTHORIZATION`; `UNRESOLVED` permanece para retenção e tratamento da API
  key, dados de jogadores, SLA, limites e compatibilidade contratual com
  Supercell.
- **Condições:** manter API key fora do browser; não registrar token, Player Tag
  real ou payload pessoal; limitar endpoints; observar cache; registrar erros;
  manter uso server-side, controlar acesso e logs, e acompanhar privacidade e
  contratos como débito explícito em `001-06`.
- **Limite:** proxy resolve egress, mas não concede API license, ownership,
  autorização de Fan Content ou autorização comercial. Billing continua sujeito
  aos gates da Fase 001.
- **Credencial anterior:** responsável confirmou revogação da chave anterior em
  `2026-10-01`; nenhum identificador ou valor foi registrado. Confirmação não foi
  verificada independentemente.
- **Fonte third-party:** TOS e Privacy Policy foram fornecidos pelo responsável,
  ambos marcados `Last Updated: 2024-10-11`. Tentativas de consulta a
  `https://royaleapi.com/terms`, `https://royaleapi.com/privacy` e
  `https://royaleapi.com/terms-of-service` retornaram `403`/Cloudflare;
  aplicabilidade específica ao endpoint proxy permanece pendente.

### Decisão de continuidade em `2026-10-01`

- RoyaleAPI Proxy deixa de ser apenas candidato e passa a ser rota operacional
  atual, até substituição deliberada ou egress próprio com IP fixo.
- A integração usa substituição de host: `https://api.clashroyale.com` vira
  `https://proxy.royaleapi.dev`; paths `/v1/...`, método, query e header Bearer
  permanecem os mesmos nos GETs observados. Isso não reescreve contratos da API.
- Essa decisão aceita dependência third-party temporária; não encerra os gates de
  privacidade, segurança, limites, SLA, Supercell, Fan Content ou monetização.
- Dados e requests devem permanecer server-side, com menor escopo possível,
  sem exposição da chave e sem transformar proxy em autorização comercial.
- O status da task permanece `in_progress`; troca futura de transporte deve ser
  registrada como decisão, não presumida como requisito imediato.

### Validação da documentação oficial em `2026-10-01`

- `https://developer.clashroyale.com/` respondeu `200` e carregou o shell oficial
  da documentação Supercell.
- `https://developer.clashroyale.com/api-docs/index.html` carregou a aplicação de
  documentação, mas o Swagger/resource specification depende de sessão
  autenticada do portal; a sessão disponível nesta execução não permitiu extrair
  o contrato completo.
- O bundle público não expôs paths de endpoint nem spec OpenAPI; o fluxo de login
  recebe `swaggerUrl` e token temporário de sessão. Não tentar contornar essa
  barreira nem registrar credenciais do portal.
- **Conclusão:** a documentação oficial confirma a fonte e a necessidade de
  sessão, mas não fecha publicamente paths, parâmetros, paginação ou ownership
  verification. Esses itens continuam `PENDING` até revisão autenticada ou
  artefato oficial exportado e sanitizado.

### Decisão de produto: salvar perfil por Player Tag

- Permitir associação read-only de Player Tag fornecida pelo usuário ao perfil
  público resolvido.
- Representar como `subjectType: public_profile` e
  `ownershipStatus: unverified`.
- Derivar `crownpilotUserId` do Firebase Auth no servidor; não persistir qualquer
  credencial Supercell.
- Exibir “Perfil público salvo — ownership não verificado”; nunca “minha conta”,
  `linkedAccount`, `verified` ou `owner`.
- Não exigir endpoint de ownership não confirmado; não solicitar player token.
- Bloquear ações na conta, exclusividade, impersonation, notificações ao jogador
  real e exposição pública da associação.
- Detalhar contrato, remoção, refresh, rate limit e privacidade nas tasks
  `001-02`, `001-05`, `001-06` e `001-07`; esta decisão não encerra essas tasks.

### Fechamento da task em `2026-10-01`

- **Status final:** `completed` com constraints; rota operacional atual é
  RoyaleAPI Proxy server-side.
- **Critérios atendidos:** houve respostas `2xx` autenticadas para cards, perfil,
  battlelog e locations; `%23` foi confirmado no Player Tag; erro de tag inválida
  retornou `404`; nenhum segredo, tag real, IP ou payload bruto foi versionado.
- **Superfície MVP classificada:** profile, cards e battlelog são `required`;
  locations é `useful`; rankings, clans e Path of Legends permanecem `useful`
  candidatos/deferidos; ownership permanece `UNRESOLVED` e nenhum `verifytoken`
  é assumido.
- **Arquivos/evidências:** este registro; `evidences/api-surface.md`; overview;
  spec da Fase 001; `.env.example` sanitizado e `.gitignore` protegendo
  `.env.local`.
- **Comandos/checks:** probes `curl` via `proxy.royaleapi.dev`; portal oficial e
  `api-docs/index.html` consultados; `git diff --check`; buscas de tag real,
  segredo e IP literal.
- **Desvios:** rota oficial direta permanece bloqueada por allowlist; proxy é
  dependência third-party operacional aceita até troca deliberada ou egress
  próprio com IP fixo. Swagger completo exige sessão do portal.
- **Riscos residuais:** termos/retensão/key handling do proxy, limites/SLA,
  contrato oficial completo, ownership, rate limits e paginação avançada seguem
  pendentes nas tasks correspondentes. Billing e autorização comercial não são
  afetados por este fechamento.
- **Handoff:** `001-02` e `001-03` podem usar proxy server-side sob controles
  documentados; nenhuma task seguinte foi marcada ou executada automaticamente.
