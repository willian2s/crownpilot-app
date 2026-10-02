# 001-06 — Validar compliance e monetização

- **Ticker:** `001`
- **Número:** `06`
- **Status:** `completed with constraints`

## Requisitos cobertos

- Fan Content Policy, Terms of Service e developer/API agreements;
- uso de marca, assets, dados, disclaimer, fan app, ads, donations e coaching;
- status verificável de assinatura, paywall e AI Coach, com billing bloqueado em
  caso de ambiguidade.

## Objetivo e resultado esperado

Definir boundaries explícitos para API, marca, assets, dados e eventual
monetização antes que decisões de produto criem dependência de um modelo não
permitido.

## Fontes mínimas

Revisar versões vigentes de:

- Fan Content Policy;
- Terms of Service;
- developer policies/agreements do portal;
- termos adicionais apresentados na criação/uso de API key.

Fontes oficiais mínimas e versão observada em `2026-10-01`:

- [Fan Content Policy](https://supercell.com/en/fan-content-policy/) — página
  indica `Last updated: September 27, 2023`;
- [Terms of Service](https://supercell.com/en/terms-of-service/) — página indica
  `Effective Date: November 6, 2024`;
- [Clash Royale API developer portal](https://developer.clashroyale.com/) —
  contratos, termos e acordos exibidos ao responsável autenticado.

Registrar URL e data de consulta.

## Escopo incluído

- uso de marca e assets;
- disclaimer obrigatório;
- domínio/nome do produto;
- restrições a bots/mods/automation;
- limites de fan apps;
- ads;
- donations;
- coaching humano;
- software coaching;
- AI/software coaching;
- assinatura/paywall;
- premium features;
- AI Coach pago;
- analytics pago;
- sponsorship;
- SaaS pago;
- AI Coach;
- developer/API data terms;
- distinção entre acesso técnico à API e autorização comercial;
- armazenamento de Player Tag associado a conta CrownPilot;
- necessidade futura de privacy policy.

## Escopo excluído

- parecer jurídico;
- criação de empresa;
- cobrança;
- termos de uso próprios definitivos;
- política de privacidade final.

## Dependências

- documentação pública atualizada e acesso aos developer/API agreements;
- pode ocorrer em paralelo com 001-02 a 001-05;
- resultado deve estar disponível antes de 001-08 e de qualquer decisão de
  billing.

## Arquivos e símbolos prováveis

- `evidences/compliance-findings.md`;
- URLs e datas das políticas consultadas, sem copiar trechos extensos;
- boundaries de produto `AI Coach`, `billing`, `Player Tag` e disclaimer; não há
  código comercial ou de identidade na `main`.

## Passos de execução

1. Consultar Fan Content Policy, Terms of Service e agreements aplicáveis.
2. Registrar URL, data e resumo curto por tópico.
3. Classificar cada uso com um dos estados permitidos.
4. Traduzir restrições em boundaries de produto e ações necessárias.
5. Marcar billing como bloqueado quando faltar autorização ou clareza.

## Descoberta incorporada ao escopo

A leitura inicial da Fan Content Policy deve ser tratada como requisito de
investigação, não como encerramento desta task. A política descreve Fan Content
como não comercial por padrão e veda cobrança de qualquer taxa, inclusive por
funcionalidades in-app, salvo aprovação expressa da Supercell. Ads, donations e
coaching aparecem como exceções, mas com condições específicas. A exceção de
coaching não resolve automaticamente software coaching ou AI Coach.

Matriz inicial que `001-06` deve confirmar contra a política vigente, Terms of
Service e developer/API agreements:

| Modelo | Classificação inicial conservadora | Evidência/decisão exigida |
| --- | --- | --- |
| Conteúdo gratuito / fan app | `ALLOWED WITH CONDITIONS` | Confirmar escopo de Fan Content, assets, marca, disclaimer e dados. |
| Ads | `ALLOWED WITH CONDITIONS` | Confirmar leis, developer policies, criativos e ausência de endorsement da Supercell. |
| Donations | `ALLOWED WITH CONDITIONS` | Confirmar que são puras, sem feature, IAP ou benefício associado. |
| Coaching humano | `ALLOWED WITH CONDITIONS` | Delimitar training/guidance e conteúdo efetivamente prestado. |
| Software coaching | `REQUIRES EXPLICIT APPROVAL` | Política não declara cobertura automática; obter esclarecimento/aprovação. |
| AI/software coaching | `REQUIRES EXPLICIT APPROVAL` | Não tratar exceção de coaching como autorização automática. |
| AI Coach gratuito | `UNRESOLVED` | Gratuidade não resolve o enquadramento de software/AI coaching, dados e assets. |
| Assinatura/paywall | `REQUIRES EXPLICIT APPROVAL` | Cobrança por acesso/feature fica bloqueada até aprovação expressa. |
| Premium features / analytics pago | `REQUIRES EXPLICIT APPROVAL` | Verificar se são funcionalidades in-app ou cobrança equivalente. |
| AI Coach pago | `REQUIRES EXPLICIT APPROVAL` | Exigir enquadramento específico; não implementar billing como premissa. |
| Sponsorship | `ALLOWED WITH CONDITIONS` | Não sugerir patrocínio, criação ou endorsement da Supercell. |
| SaaS pago | `REQUIRES EXPLICIT APPROVAL` | Não assumir exceção geral para modelo SaaS tradicional. |
| Assets não modificados | `ALLOWED WITH CONDITIONS` | Confirmar uso limitado a exibir, identificar e discutir produtos Supercell. |
| Assets modificados | `REQUIRES EXPLICIT APPROVAL` | A política exige permissão expressa para modificação. |
| Domínio/handle sem trademark Supercell | `ALLOWED WITH CONDITIONS` | Confirmar ausência de confusão com endorsement e de imitação de branding. |
| Domínio/handle com trademark Supercell ou nome de jogo | `REQUIRES EXPLICIT APPROVAL` | Exige acordo escrito separado segundo a política. |
| Bots, mods, automação, private servers e software não autorizado | `NOT ALLOWED` | Registrar como boundary proibido, não como modelo comercial disponível. |

Se a fonte não responder de modo suficiente ao caso concreto, usar `UNRESOLVED`
em vez de inferir permissão. `NOT ALLOWED` só deve ser usado quando houver
vedação explícita. O artefato `evidences/compliance-findings.md` deve registrar
status por modelo, URL, data, resumo da fonte, impacto e ação necessária.

### Gate API access versus autorização comercial

API key, autenticação, IP/egress, endpoint documentado ou eventual resposta `2xx`
demonstram acesso técnico à API no escopo concedido. Não demonstram permissão
para cobrar, vender funcionalidades, usar software/AI coaching como exceção ou
operar SaaS pago. Esta task deve fechar os dois gates separadamente e apontar
qualquer dependência de aprovação à Fase 001 e à decisão posterior da Fase 012.

## Questões obrigatórias

1. O app gratuito de guia/análise é permitido no formato proposto?
2. Quais assets podemos usar sem modificar?
3. Qual disclaimer deve aparecer e onde?
4. `crownpilot` evita uso indevido de trademark em domínio/handle?
5. Assinatura por funcionalidades é proibida, permitida ou exige autorização?
6. "AI Coach" pago pode ser coaching segundo a política ou isso continua
   dependente de avaliação da Supercell?
7. Ads e donations são alternativas permitidas sob quais condições?
8. Coaching humano, software coaching e AI/software coaching recebem o mesmo
   tratamento ou permanecem `UNRESOLVED`/`REQUIRES EXPLICIT APPROVAL`?
9. Analytics pago, premium features e SaaS pago são cobrança por funcionalidades?
10. Sponsorship pode ser usado sem sugerir endorsement da Supercell?
11. Algum developer agreement restringe armazenamento/redistribuição de API data?
12. API access pode ser confundido com autorização comercial em algum contrato?
13. Quais disclaimers são obrigatórios e onde devem aparecer?
14. Existem restrições específicas para assets, nome, branding, domínio e handles?
15. Precisamos contatar Supercell antes de qualquer cobrança?

## Regra de decisão

Ambiguidade comercial não vira "permitido".

Usar exatamente estes estados:

- **`ALLOWED`**;
- **`ALLOWED WITH CONDITIONS`**;
- **`REQUIRES EXPLICIT APPROVAL`**;
- **`NOT ALLOWED`**;
- **`UNRESOLVED`**.

## Evidência obrigatória

`evidences/compliance-findings.md` com:

| Topic | Policy/source | Current text summary | Status | Product impact | Required action |
| --- | --- | --- | --- | --- | --- |

Não copiar longos trechos protegidos; resumir e linkar a fonte.

## Gate de monetização

Até esta task indicar base suficiente:

- não implementar billing;
- não prometer plano Pro;
- não atrelar feature essencial a pagamento;
- não assumir que exceção de coaching cobre software automaticamente;
- não tratar API key, acesso a endpoint ou ownership verification como autorização
  comercial;
- não liberar modelo classificado como `REQUIRES EXPLICIT APPROVAL` até existir
  aprovação expressa e rastreável.

## Definição de pronto

O roadmap comercial possui boundaries explícitos; cada modelo possui status,
fonte, data, condições e ação; dúvidas que exigem contato externo estão
documentadas como bloqueio, não escondidas. O resultado também informa se o MVP
pode operar gratuitamente e se um eventual veredito é
`GO WITH CONSTRAINTS / APPROVAL DEPENDENCY`, sem converter isso em aprovação.

## Testes e comandos de validação

- consultar cada fonte oficial na data registrada e comparar a interpretação
  com o texto vigente;
- classificar cada tópico somente como `ALLOWED`, `ALLOWED WITH CONDITIONS`,
  `REQUIRES EXPLICIT APPROVAL`, `NOT ALLOWED` ou `UNRESOLVED`;
- revisar que billing, plano Pro e feature essencial paga não sejam tratados como
  permitidos por inferência;
- revisar separadamente acesso técnico à API e autorização comercial;
- confirmar que ads, donations e coaching possuem condições próprias;
- confirmar que software/AI coaching não foi tratado como automaticamente coberto
  pela exceção de coaching;
- executar `git diff --check` e remover dados pessoais desnecessários.

## Riscos e cuidados

- políticas podem mudar;
- coaching é definido pela Supercell e pode exigir interpretação;
- política não esclarece automaticamente software coaching ou AI Coach;
- acesso à API não prova autorização para monetização;
- donations podem ser incompatíveis com qualquer benefício vinculado;
- API agreement autenticado pode conter cláusulas não visíveis publicamente;
- associar login a Player Tag cria tratamento de dados que merece privacy review.

## Registro de execução

- **Status final:**
- **Políticas/data de consulta:**
- **Modelo gratuito:**
- **Monetização:**
- **AI Coach:**
- **Aprovação externa necessária:**
- **Riscos residuais:**

### Atualização documental em `2026-10-01`

- **Status:** permanece `pending`; esta atualização não encerra a task nem cria
  `compliance-findings.md` final.
- **Arquivos alterados:** esta task, spec da Fase 001, roadmap, overview, task
  `001-01`, task `001-08` e contexto histórico.
- **Decisões incorporadas:** cobrança por funcionalidades, assinatura/paywall,
  premium features, analytics pago, AI Coach pago, software coaching e SaaS pago
  ficam em `REQUIRES EXPLICIT APPROVAL` até aprovação expressa; AI Coach gratuito
  fica `UNRESOLVED`; ads, donations, coaching humano e sponsorship ficam
  `ALLOWED WITH CONDITIONS` sujeitos à confirmação detalhada; assets modificados,
  domínios/handles com trademarks e usos de bots/mods/automação ficam
  explicitamente cobertos como `REQUIRES EXPLICIT APPROVAL` ou `NOT ALLOWED`.
- **Comandos/evidências:** consulta das fontes oficiais em `2026-10-01`, busca
  cruzada de referências e `git diff --check`; não houve cobrança, billing ou
  probe autenticado da API.
- **Desvios:** developer/API agreements autenticados e enquadramento final de
  software/AI coaching ainda não foram fechados; a matriz é requisito de
  investigação e baseline conservador, não veredito final da task.
- **Riscos residuais:** policy pode mudar; Supercell reserva interpretação de
  coaching; API access não prova autorização comercial; donations não podem
  conceder benefícios; aprovação externa pode ser necessária.

### Execução final em `2026-10-02`

- **Status final:** `completed with constraints`; evidência final criada em
  `evidences/compliance-findings.md`.
- **Arquivos alterados:** `evidences/compliance-findings.md`, esta task e
  `001-00-overview.md`.
- **Fontes/data de consulta:** Fan Content Policy consultada em `2026-10-02`,
  com `Last updated: September 27, 2023`; Terms of Service consultados em
  `2026-10-02`, com `Effective Date: November 6, 2024`; Privacy Policy consultada
  em `2026-10-02`, com `Effective Date: March 11th, 2026`; portal/API docs oficial
  consultado em `2026-10-02`. A versão visível da Fan Content Policy permaneceu
  consistente com o registro de `2026-10-01`.
- **Decisões:** fan app/MVP gratuito é `ALLOWED WITH CONDITIONS`; ads,
  donations, coaching humano e sponsorship são `ALLOWED WITH CONDITIONS`; assets
  não modificados e `crownpilot` são condicionais; cobrança, assinatura,
  paywall, premium features, analytics pago, SaaS, AI Coach pago, software/AI
  coaching e assets modificados são `REQUIRES EXPLICIT APPROVAL`; AI Coach
  gratuito e storage/redistribuição de API data permanecem `UNRESOLVED`; bots,
  mods, automação, private servers, account trading e software não autorizado
  são `NOT ALLOWED`.
- **Boundaries:** MVP pode operar gratuitamente apenas sem billing, sem promessa
  de ownership e sob minimização/privacidade de Player Tag, disclaimer legível,
  assets permitidos e revisão dos agreements; API access e autorização comercial
  permanecem gates separados.
- **Monetização:** resultado da task é `GO WITH CONSTRAINTS / APPROVAL DEPENDENCY`;
  não é aprovação jurídica ou Supercell. Billing, plano Pro e feature essencial
  paga continuam bloqueados.
- **Comandos/resultados/evidências:** consultas oficiais via `webfetch` e `curl`
  retornaram HTTP `200` nas cinco URLs registradas; validação documental via
  `python3` confirmou ticker/número/status, checklist único com oito itens,
  progresso `6/8`, 25 linhas classificadas somente com estados permitidos e 15/15
  respostas; busca via `rg` não encontrou credenciais, Player Tags, IPs ou e-mails
  nos três arquivos alterados; `git diff --check` terminou sem saída/erro.
- **Revisão final:** corrigida a classificação da autorização comercial derivada
  de API access para `REQUIRES EXPLICIT APPROVAL`, alinhando a matriz ao gate de
  cobrança; Privacy Policy incluída como fonte do boundary de Player Tag, sem
  tratá-la como autorização para o tratamento feito pelo CrownPilot.
- **Desvios:** developer/API agreements autenticados não estavam disponíveis
  nesta sessão; por isso storage, cache, retenção e redistribuição de API data
  permanecem `UNRESOLVED`, e o portal não foi tratado como autorização comercial.
- **Aprovação externa necessária:** sim, antes de qualquer cobrança, modelo SaaS,
  premium feature ou software/AI coaching; também antes de assets modificados ou
  domínio/handle com trademark/nome de jogo.
- **Riscos residuais:** policy e Terms podem mudar; interpretação de coaching é
  reservada à Supercell; termos autenticados, proxy, retenção de API data e
  privacy review do vínculo Player Tag permanecem pendentes. `001-08` deve
  consumir este gate sem converter dependência em autorização.
- **Revisão independente:** aprovada com ressalvas; confirmou estrutura SDD,
  matriz, questões obrigatórias, gates comerciais separados e ausência de
  bloqueadores adicionais. Riscos de agreements autenticados, privacidade e
  interpretação de coaching permanecem conforme registrado acima.
