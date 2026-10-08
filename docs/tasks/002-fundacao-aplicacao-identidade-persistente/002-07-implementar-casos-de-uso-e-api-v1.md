# 002-07 — Implementar casos de uso e API v1

- **Ticker:** `002`
- **Número:** `07`
- **Status:** `pending`

Esta task permanece bloqueada até ADR 005 ser aprovada e `002-13` -> `002-14` ->
`002-15` concluírem gates verdes. O baseline ASP.NET/.NET permanece histórico;
implementação deve seguir `net/http`/`ServeMux` e o ajuste Go de `002-15`.

## Objetivo e resultado esperado

Orquestrar casos de uso de identidade e vínculo e expor API REST `/api/v1`, com
autorização server-side, contrato OpenAPI completo e comunicação honesta de
perfil público não verificado. A UI fica na Task
`002-08-entregar-frontend-de-identidade-e-vinculo.md`.

## Requisitos cobertos

- uma Player Tag primária;
- API Go com `net/http`/`ServeMux`, consumível pelo frontend Vite;
- recuperação em outro dispositivo;
- troca/desvinculação explícitas;
- erros de input/provider compreensíveis via ProblemDetails;
- contrato OpenAPI gerado e validado sem arquivo manual concorrente;
- exclusão inicial de dados próprios;
- disclaimer de conteúdo não oficial.

## Escopo incluído

- handlers finos `net/http` para read/link/replace/unlink/delete;
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

- ADR 005 aprovada, `002-13-bootstrap-http-config-openapi-go.md`,
  `002-14-autenticacao-firebase-go.md`, `002-15-persistencia-cutover-remocao-dotnet.md`,
  `002-04-implementar-google-sign-in-e-firebase-bearer.md` como histórico e
  `002-05-implementar-port-e-adapter-de-lookup.md`;
  `002-06-modelar-persistencia-repositories-e-rls.md`;
- contrato HTTP/ProblemDetails e política de observabilidade da spec;
- frontend/API de `002-01-bootstrap-toolchain.md`;
- hostname fixo de Staging para validação posterior.

## Arquivos e símbolos prováveis

- `internal/httpapi/`;
- `internal/identity/` e `internal/playerlink/`;
- DTOs Go de request/response e mapeamento de erros;
- caso de uso `DeleteOwnData`;
- disclaimer em layout/footer acessível.

## Passos de implementação

1. Implementar casos de uso de leitura, vínculo, replace, unlink e exclusão.
2. Validar tag no Application e executar lookup server-side antes da escrita.
3. Preservar vínculo anterior em falha externa e proteger replace por
   `expectedVersion`, retornando `409` sem last-write-wins.
4. Exigir `auth_time` recente para exclusão e mapear reautenticação necessária.
5. Expor endpoints finos e DTOs sem entidades de persistência ou payload externo.
6. Gerar tipos Go com `oapi-codegen` a partir de `api/openapi/v1.json`, manter
   `openapi-typescript` no frontend e documentar todos os status aplicáveis.
7. Adicionar autorização e API tests para A/B, UID arbitrário e todos os estados
   do provider, incluindo `Retry-After`, `traceId` e redaction.

## Testes e comandos de validação

```text
go test ./internal/identity/... ./internal/playerlink/...
go test ./internal/httpapi/...
npm run openapi:check --prefix frontend
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
