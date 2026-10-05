# 002-07 — Implementar casos de uso e API v1

- **Ticker:** `002`
- **Número:** `07`
- **Status:** `pending`

## Objetivo e resultado esperado

Orquestrar casos de uso de identidade e vínculo e expor API REST `/api/v1`, com
autorização server-side, contrato OpenAPI completo e comunicação honesta de
perfil público não verificado. A UI fica na Task 002-08.

## Requisitos cobertos

- uma Player Tag primária;
- API ASP.NET Core consumível pelo frontend Vite;
- recuperação em outro dispositivo;
- troca/desvinculação explícitas;
- erros de input/provider compreensíveis via ProblemDetails;
- contrato OpenAPI gerado e validado sem arquivo manual concorrente;
- exclusão inicial de dados próprios;
- disclaimer de conteúdo não oficial.

## Escopo incluído

- controllers/minimal endpoints finos para read/link/replace/unlink/delete;
- casos de uso para read/link/replace/unlink/delete usando ports;
- paths `/api/v1/me/player-link` e `/api/v1/me` conforme contrato da spec;
- DTOs JSON sem entidades de persistência expostas;
- respostas com `400`, `401`, `403`, `404`, `409`, `422` quando aplicável, `429`,
  `500`, `503` e `ProblemDetails` sem vazamento;
- `auth_time` recente para exclusão, com distinção entre claim ausente, malformada,
  futura e expirada;
- OpenAPI documentando método, authn/authz, parâmetros, body, response, status e
  exemplos úteis;
- preservação do vínculo anterior em falha de provider e `409` em concorrência;
- autenticação recente (`auth_time`) para exclusão, sem excluir Firebase/Google.

## Escopo excluído

- claim de “minha conta”, ownership ou exclusividade;
- múltiplas tags, social, clãs ou notificações;
- dashboard de coleção, Arena ou sync completo;
- cobrança, assets não aprovados ou coaching;
- resposta raw ou nome do provider externo.

## Dependências

- `002-04`, `002-05` e `002-06`;
- contrato HTTP/ProblemDetails e política de observabilidade da spec;
- frontend/API de `002-01`;
- hostname fixo de Staging para validação posterior.

## Arquivos e símbolos prováveis

- `src/Api/Endpoints/PlayerLinkEndpoints.cs` ou controllers equivalentes;
- `src/Application/PlayerLinks/`;
- `PlayerLinkRequest`, `PlayerLinkResponse`, `MapLinkError`;
- `DeleteOwnData` use case;
- disclaimer em layout/footer acessível.

## Passos de implementação

1. Implementar casos de uso de leitura, vínculo, replace, unlink e exclusão.
2. Validar tag no Application e executar lookup server-side antes da escrita.
3. Preservar vínculo anterior em falha externa e proteger replace por
   `expectedVersion`, retornando `409` sem last-write-wins.
4. Exigir `auth_time` recente para exclusão e mapear reautenticação necessária.
5. Expor endpoints finos e DTOs sem entidades EF ou payload externo.
6. Gerar OpenAPI code-first único e documentar todos os status aplicáveis.
7. Adicionar autorização e API tests para A/B, UID arbitrário e todos os estados
   do provider, incluindo `Retry-After`, `traceId` e redaction.

## Testes e comandos de validação

```text
dotnet test --filter Category=Application
dotnet test --filter Category=Api
npm run openapi:check
```

Cobrir ausência inicial, sucesso, input inválido/malformado, perfil inexistente,
provider indisponível/mal configurado, rate limit, replace falho, conflito,
unlink repetido, delete repetido, `401`/`403`, `422` quando aplicável, erro
inesperado `500` e claim `auth_time` inválida.

## Definição de pronto

- casos de uso usam somente identidade resolvida no backend;
- API retorna somente dados do usuário autenticado;
- replace valida nova tag antes de substituir antiga;
- falha externa preserva vínculo anterior;
- unlink e exclusão funcionam sem reentrada indevida;
- OpenAPI documenta request/response, authn/authz, ProblemDetails e status;
- exclusão exige reautenticação recente e remove somente dados CrownPilot.

## Riscos e cuidados

- Não expor nome, coleção ou resposta raw para “enriquecer” tela.
- Não confundir exclusão local com exclusão no provider externo.
- Não enviar Player Tag real para analytics/logs client-side.
- Não criar promessa de sync a partir de validação única.
- Não colocar autorização crítica somente no React.
