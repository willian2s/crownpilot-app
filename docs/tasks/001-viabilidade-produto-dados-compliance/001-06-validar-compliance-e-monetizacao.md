# 001-06 — Validar compliance e monetização

- **Ticker:** `001`
- **Número:** `06`
- **Status:** `pending`

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

Registrar URL e data de consulta.

## Escopo incluído

- uso de marca e assets;
- disclaimer obrigatório;
- domínio/nome do produto;
- restrições a bots/mods/automation;
- limites de fan apps;
- ads;
- donations;
- coaching;
- assinatura/paywall;
- AI Coach;
- developer/API data terms;
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

## Questões obrigatórias

1. O app gratuito de guia/análise é permitido no formato proposto?
2. Quais assets podemos usar sem modificar?
3. Qual disclaimer deve aparecer e onde?
4. `crownpilot` evita uso indevido de trademark em domínio/handle?
5. Assinatura por funcionalidades é proibida, permitida ou exige autorização?
6. "AI Coach" pago pode ser coaching segundo a política ou isso continua
   dependente de avaliação da Supercell?
7. Ads e donations são alternativas explicitamente permitidas?
8. Algum developer agreement restringe armazenamento/redistribuição de API data?
9. Precisamos contatar Supercell antes de qualquer cobrança?

## Regra de decisão

Ambiguidade comercial não vira "permitido".

Usar estes estados:

- **allowed by published policy**;
- **allowed with constraints**;
- **requires explicit approval / clarification**;
- **not allowed**;
- **unresolved**.

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
- não assumir que exceção de coaching cobre software automaticamente.

## Definição de pronto

O roadmap comercial possui boundaries explícitos e dúvidas que exigem contato
externo estão documentadas como bloqueio, não escondidas.

## Testes e comandos de validação

- consultar cada fonte oficial na data registrada e comparar a interpretação
  com o texto vigente;
- classificar cada tópico somente como `allowed by published policy`, `allowed
  with constraints`, `requires explicit approval / clarification`, `not allowed`
  ou `unresolved`;
- revisar que billing, plano Pro e feature essencial paga não sejam tratados como
  permitidos por inferência;
- executar `git diff --check` e remover dados pessoais desnecessários.

## Riscos e cuidados

- políticas podem mudar;
- coaching é definido pela Supercell e pode exigir interpretação;
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
