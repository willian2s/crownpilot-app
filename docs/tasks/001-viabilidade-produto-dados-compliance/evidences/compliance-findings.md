# Compliance findings — Fan Content, API e monetização

- **Task:** `001-06`
- **Consulta desta execução:** `2026-10-02`
- **Classificação:** descoberta documental; não constitui parecer jurídico
- **Regra:** ambiguidade permanece bloqueio. Acesso técnico à API, disclaimer ou
  uso de assets não equivalem a autorização comercial.

## Fontes oficiais consultadas

| Fonte | URL | Versão/data visível | Consulta | Escopo observado |
| --- | --- | --- | --- | --- |
| Fan Content Policy | https://supercell.com/en/fan-content-policy/ | `Last updated: September 27, 2023` | `2026-10-02` | Fan Content, cobrança, ads, donations, coaching, assets, marca, disclaimer, domínios e software proibido. |
| Terms of Service | https://supercell.com/en/terms-of-service/ | `Effective Date: November 6, 2024` | `2026-10-02` | Uso do Service, propriedade intelectual, restrições a automação/mods, suspensão e políticas incorporadas. |
| Privacy Policy | https://supercell.com/en/privacy-policy/ | `Effective Date: March 11th, 2026` | `2026-10-02` | Confirma que Player Tag é dado coletado pela Supercell; não define obrigações próprias do CrownPilot nem autoriza retenção de API data. |
| Clash Royale API developer portal | https://developer.clashroyale.com/ | portal/documentação consultados; agreements dependem de sessão autenticada | `2026-10-02` | API key, Bearer JWT, allowlist de IP e existência da documentação; não foi obtido contrato autenticado nesta execução. |
| API documentation route | https://developer.clashroyale.com/api-docs/index.html | Swagger shell; specification depende de sessão | `2026-10-02` | Fonte oficial para contratos de endpoint; paths, termos adicionais, retenção e redistribuição permanecem pendentes sem sessão aplicável. |

O texto visível da Fan Content Policy permaneceu consistente com a versão já
registrada em `2026-10-01`. O portal oficial não foi tratado como prova de que
um acordo autenticado foi aceito ou revisado. As páginas da RoyaleAPI e seus
termos são dependência operacional third-party já registrada em
`evidences/api-surface.md`, não fonte de autorização Supercell.

## Matriz de classificação

| Topic | Policy/source | Current text summary | Status | Product impact | Required action |
| --- | --- | --- | --- | --- | --- |
| Fan app gratuito de guia/análise | Fan Content Policy, seção “Fan Content Purposes only”, consulta `2026-10-02` | Guias e guide apps online não comerciais aparecem como exemplos, desde que exibam, identifiquem ou discutam produtos Supercell e cumpram a política. | `ALLOWED WITH CONDITIONS` | MVP gratuito é uma hipótese viável, não uma autorização irrestrita. | Limitar escopo a guia/análise; revisar assets, marca, dados, disclaimer e policies antes de publicar. |
| Conteúdo gratuito | Fan Content Policy, regra “Be non-commercial” | Gratuidade evita fee, mas não remove as demais regras de Fan Content, marca, assets ou conteúdo proibido. | `ALLOWED WITH CONDITIONS` | Não atrelar feature essencial a pagamento nem usar gratuidade como justificativa para uso fora da política. | Manter MVP sem billing e repetir revisão quando escopo mudar. |
| Ads | Fan Content Policy, regra de monetização e FAQ sobre promotions/ads | Ads são exceção geralmente permitida, sujeitos a leis, regras, developer policies e ausência de impressão de patrocínio da Supercell. | `ALLOWED WITH CONDITIONS` | Provedor, criativo e placement criam gates próprios; não são aprovados por esta matriz sozinhos. | Validar provedor, jurisdição, criativos e copy; remover qualquer sugestão de endorsement. |
| Donations | Fan Content Policy, regra de monetização | Donations precisam ser puras e não podem estar ligadas a feature especial, IAP ou qualquer benefício. | `ALLOWED WITH CONDITIONS` | Doação não pode desbloquear AI Coach, analytics, limites, prioridade ou conteúdo. | Se adotada, separar fluxo e benefício zero; revisar copy e provedor. |
| Coaching humano | Fan Content Policy, FAQ “What is coaching?” | Coaching é training/guidance; exemplos incluem coaching pessoal/online e venda de base layouts. Supercell reserva a definição final. | `ALLOWED WITH CONDITIONS` | Serviço humano pode ser hipótese condicionada, sem apresentar interpretação como aprovação. | Definir conteúdo, operador e limites; confirmar termos aplicáveis antes de oferta. |
| Software coaching | Fan Content Policy, regra de coaching e reserva de interpretação | A política descreve treinamento/guidance, mas não declara que software automatizado é coberto. | `REQUIRES EXPLICIT APPROVAL` | Não lançar nem cobrar software coaching por analogia com coaching humano. | Obter esclarecimento/aprovação escrita e rastreável da Supercell. |
| AI/software coaching | Fan Content Policy, mesmas seções | Não há autorização específica para AI/software coaching. “Similar activities” não fecha o enquadramento. | `REQUIRES EXPLICIT APPROVAL` | Não posicionar AI Coach como exceção automaticamente permitida. | Solicitar avaliação específica, incluindo dados, assets, outputs e automação. |
| AI Coach gratuito | Fan Content Policy, regras de Fan Content e não comercialidade | Ausência de cobrança resolve apenas fee; não resolve software/AI coaching, dados, assets ou enquadramento do produto. | `UNRESOLVED` | Não comprometer AI Coach gratuito no MVP como permitido. | Tratar como experimento bloqueado até esclarecimento de enquadramento. |
| Assinatura/paywall | Fan Content Policy, regra “Be non-commercial” | Cobrar fee de qualquer tipo, inclusive por funcionalidades in-app, exige aprovação expressa. | `REQUIRES EXPLICIT APPROVAL` | Billing, plano Pro e paywall ficam bloqueados. | Não implementar billing; obter aprovação expressa antes de desenhar cobrança. |
| Premium features | Fan Content Policy, regra sobre in-app functionalities | A vedação alcança cobrança por funcionalidades; não existe exceção geral para premium features. | `REQUIRES EXPLICIT APPROVAL` | Não reservar analytics, limites ou recursos essenciais a pagantes. | Manter paridade gratuita até aprovação específica. |
| Analytics pago | Fan Content Policy, regra sobre fee/in-app functionalities | Analytics pago é cobrança por acesso/funcionalidade quando baseado no Fan Content. | `REQUIRES EXPLICIT APPROVAL` | Plano pago e upgrade de analytics bloqueados. | Reavaliar com escopo de dados e feature descritos em pedido de aprovação. |
| AI Coach pago | Fan Content Policy, regras de fee e coaching | Exceção de coaching não cobre automaticamente software/AI pago; cobrança continua vedada sem aprovação. | `REQUIRES EXPLICIT APPROVAL` | Não criar preço, paywall ou promessa comercial. | Exigir aprovação escrita específica antes de billing ou marketing. |
| Sponsorship | Fan Content Policy, regra de marca e FAQ sobre ads | Promoção pode existir sob condições, mas Fan Content não pode sugerir que Supercell patrocina, cria ou endossa o produto. | `ALLOWED WITH CONDITIONS` | Patrocinador deve ser terceiro e não pode parecer parceria oficial. | Aprovar contrato, disclosure, criativos e uso de marcas separadamente. |
| SaaS pago | Fan Content Policy, regra geral de não comercialidade | Política não cria exceção geral para SaaS tradicional; cobrança por acesso/feature permanece fora do permitido sem aprovação. | `REQUIRES EXPLICIT APPROVAL` | Não assumir assinatura B2C/B2B como caminho liberado. | Pedir autorização com modelo, público, dados e funcionalidades detalhados. |
| Assets não modificados | Fan Content Policy, “Fan Content Purposes only” e “Respect Supercell’s brand” | Assets podem ser usados para exibir, identificar e discutir produtos Supercell, sem imitar branding ou criar impressão oficial. | `ALLOWED WITH CONDITIONS` | Inventário de assets e revisão visual são obrigatórios. | Usar somente assets necessários, sem alteração ou imitação de logo/trademark. |
| Assets modificados | Fan Content Policy, “You can’t modify Supercell Assets without our express permission” | Modificação exige permissão expressa. | `REQUIRES EXPLICIT APPROVAL` | Não editar, recortar para criar novo logo, recolorir ou derivar assets sem aprovação. | Preservar originais ou obter autorização escrita rastreável. |
| Domínio/handle `crownpilot` sem trademark Supercell | Fan Content Policy, “No domains or the like” e regra de marca | Nome sem marca/nome de jogo não é automaticamente proibido, mas não pode imitar branding nem sugerir endorsement. | `ALLOWED WITH CONDITIONS` | `crownpilot` é preferência compatível sob revisão de confusão; não é garantia jurídica. | Revisar domínio, handle, logo, title e copy antes de uso público. |
| Domínio/handle com trademark Supercell ou nome de jogo | Fan Content Policy, “No domains or the like” | Endereços que incluem trademarks Supercell ou nomes de jogos exigem acordo escrito separado. | `REQUIRES EXPLICIT APPROVAL` | Não registrar ou publicar domínio/handle com “Supercell” ou “Clash Royale” como marca do produto. | Obter acordo escrito separado antes do registro/uso. |
| Disclaimer | Fan Content Policy, “Insert disclaimers” | Aviso de não oficialidade deve ser legível e aparecer em conexão com o Fan Content. | `ALLOWED WITH CONDITIONS` | Ausência ou placement inadequado bloqueia publicação com assets. | Exibir, em inglês ou versão substancialmente similar legível: “This material is unofficial and is not endorsed by Supercell. For more information see Supercell's Fan Content Policy: www.supercell.com/fan-content-policy.”; colocar junto das telas/páginas e materiais que usam Fan Content. |
| Bots, mods, cheats, hacks, automação, private servers e software não autorizado | Fan Content Policy, “No bad stuff”; Terms of Service, seção 1.1 | Policy proíbe Fan Content que promova/contenha esses usos; Terms também proíbem software que modifique/interfira no Service. | `NOT ALLOWED` | Fora do produto, conteúdo, suporte, marketing e monetização. | Rejeitar integrações, copy, tutoriais e links com esses usos. |
| Account trading/boosting | Fan Content Policy FAQ; Terms of Service, seção 1.1 | Venda/troca de contas e boosting são proibidos. | `NOT ALLOWED` | Não oferecer, facilitar ou anunciar serviço de conta. | Bloquear conteúdo e fluxos relacionados. |
| Acesso técnico à API | Developer portal/API docs, consulta `2026-10-02`; probes em `2026-10-01` | API key, Bearer JWT, allowlist e resposta `2xx` demonstram no máximo acesso técnico ao escopo concedido; proxy foi observado como transporte operacional third-party. | `ALLOWED WITH CONDITIONS` | Uso continua server-side, sem expor key; proxy não vira contrato Supercell. | Revisar acordo autenticado, origem/egress, limites, segurança, retenção e termos do proxy. |
| Autorização comercial derivada de API key, endpoint ou `2xx` | Fan Content Policy + developer portal; consulta `2026-10-02` | Credencial, endpoint, IP allowlist ou sucesso HTTP não concedem permissão comercial; cobrança por Fan Content requer aprovação expressa. | `REQUIRES EXPLICIT APPROVAL` | Billing permanece bloqueado mesmo com acesso técnico funcionando. | Fechar autorização comercial em gate separado e guardar aprovação rastreável. |
| Armazenamento/redistribuição de API data | Developer/API agreements autenticados; portal consultado `2026-10-02` | Agreement aplicável e termos adicionais da criação/uso da key não ficaram disponíveis nesta sessão; limites públicos de armazenamento, cache e redistribuição não foram fechados. | `UNRESOLVED` | Não prometer retenção histórica, exportação pública ou redistribuição até revisão. | Acessar agreements como responsável autenticado; registrar cláusulas, retenção, cache, sublicença e deleção. |
| Player Tag associada à conta CrownPilot | Terms of Service/Privacy Policy incorporados + developer/API agreements; revisão `2026-10-02` | A associação cria tratamento de dado no produto; fontes consultadas não fecham por si só base, retenção, direitos, privacy policy ou autorização para armazenar/redistribuir resposta da API. | `UNRESOLVED` | Pode ser associação privada de perfil público, nunca ownership; privacy review continua necessária. | Minimizar dados, manter server-side/privado, definir retenção e criar privacy policy antes de produção. |

## Respostas às questões obrigatórias

1. **App gratuito de guia/análise:** sim, como hipótese `ALLOWED WITH CONDITIONS`
   para Fan Content não comercial e limitado; não é aprovação de todos os dados,
   assets ou integrações.
2. **Assets sem modificar:** somente para exibir, identificar e discutir produtos
   Supercell, sem imitar branding e dentro das demais regras. Modificação requer
   aprovação expressa.
3. **Disclaimer:** usar o aviso de não oficialidade indicado na matriz, legível e
   em conexão com cada área/material de Fan Content.
4. **`crownpilot`:** não contém trademark Supercell/nome de jogo e pode ser usado
   sob condição de não confundir ou sugerir endorsement; revisão visual e de copy
   ainda é obrigatória.
5. **Assinatura por funcionalidades:** `REQUIRES EXPLICIT APPROVAL`; não liberar
   billing ou plano Pro por inferência.
6. **AI Coach pago:** `REQUIRES EXPLICIT APPROVAL`; exceção de coaching não resolve
   software/AI coaching.
7. **Ads e donations:** `ALLOWED WITH CONDITIONS`; ads seguem leis/policies e não
   podem sugerir endorsement; donations não podem conceder benefício.
8. **Coaching:** coaching humano é `ALLOWED WITH CONDITIONS`; software/AI não recebe
   o mesmo tratamento e fica em `REQUIRES EXPLICIT APPROVAL` ou `UNRESOLVED` quando
   gratuito.
9. **Analytics/premium/SaaS:** `REQUIRES EXPLICIT APPROVAL` quando cobrados como
   acesso ou funcionalidade.
10. **Sponsorship:** `ALLOWED WITH CONDITIONS`, desde que seja claramente de terceiro
    e não sugira criação, patrocínio ou endorsement da Supercell.
11. **API data:** storage/redistribuição permanece `UNRESOLVED` até revisão do
    agreement autenticado e das condições do proxy.
12. **API access versus comercial:** são gates independentes; key, endpoint, IP
    allowlist e `2xx` não autorizam cobrança.
13. **Disclaimers:** aviso legível em conexão com Fan Content; placement definitivo
    deve cobrir app, páginas e materiais de marketing relevantes.
14. **Assets/nome/branding/domínio/handles:** assets modificados e endereços com
    trademark/nome de jogo requerem aprovação; `crownpilot` segue condicional.
15. **Contato antes de cobrança:** sim. Qualquer fee, paywall, premium feature,
    SaaS ou AI/software coaching comercial fica bloqueado até aprovação expressa e
    rastreável.

## Boundaries de produto e gate de monetização

### Pode avançar, sob condições

- MVP de guia/análise gratuito, read-only e sem promessa de ownership;
- dados e Player Tag minimizados, privados por padrão, sem exposição da API key;
- assets não modificados usados apenas para exibir, identificar e discutir produtos;
- disclaimer legível em conexão com Fan Content;
- ads, donations, coaching humano e sponsorship somente após revisão própria das
  condições de cada modelo.

### Bloqueado

- billing, assinatura, paywall, plano Pro e premium features;
- analytics pago e SaaS pago;
- AI Coach pago e qualquer software/AI coaching apresentado como exceção;
- AI Coach gratuito como feature liberada por padrão, pois permanece `UNRESOLVED`;
- bots, mods, automação, private servers, account trading, boosting e software não
  autorizado;
- assets modificados e domínios/handles com trademark Supercell ou nome de jogo sem
  acordo escrito separado.

O resultado comercial desta task é **GO WITH CONSTRAINTS / APPROVAL DEPENDENCY**:
um MVP gratuito pode continuar em discovery sob os boundaries acima, mas toda
receita dependente de aprovação expressa permanece bloqueada. Isso não é parecer
jurídico, não é aprovação Supercell e não substitui o veredito final de `001-08`.

## Pendências e ações rastreáveis

1. Responsável deve revisar no portal autenticado os developer/API agreements e os
   termos apresentados na criação/uso da key.
2. Registrar se esses termos permitem armazenamento, cache, retenção, derivação,
   exposição e redistribuição de API data; até então, status permanece
   `UNRESOLVED`.
3. Solicitar avaliação/aprovação escrita para qualquer fee, SaaS, premium feature,
   software/AI coaching ou AI Coach pago antes de produto, marketing ou billing.
4. Fazer privacy review do vínculo Player Tag–CrownPilot e publicar privacy policy
   antes de operação com contas reais.
5. Revalidar URLs, versões e classificações quando a política ou o produto mudar.
