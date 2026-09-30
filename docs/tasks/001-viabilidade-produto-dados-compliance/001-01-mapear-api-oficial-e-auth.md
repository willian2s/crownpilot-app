# 001-01 — Mapear API oficial e autenticação

- **Ticker:** `001`
- **Número:** `01`
- **Status:** `pending`

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

Preencher ao executar:

- **Status final:**
- **Data:**
- **Fontes consultadas:**
- **Probes executados:**
- **Resultados:**
- **Desvios:**
- **Riscos residuais:**
