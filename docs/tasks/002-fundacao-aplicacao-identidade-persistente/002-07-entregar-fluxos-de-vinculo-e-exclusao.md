# 002-07 — Entregar fluxos de vínculo e exclusão

- **Ticker:** `002`
- **Número:** `07`
- **Status:** `pending`

## Objetivo e resultado esperado

Entregar API JSON e experiência React para consultar, criar, trocar, desvincular
e excluir dados CrownPilot, com comunicação honesta de perfil público não
verificado.

## Requisitos cobertos

- uma Player Tag primária;
- API ASP.NET Core consumida pelo frontend Vite;
- recuperação em outro dispositivo;
- troca/desvinculação explícitas;
- erros de input/provider compreensíveis via ProblemDetails;
- contrato OpenAPI gerado e validado sem arquivo manual concorrente;
- exclusão inicial de dados próprios;
- disclaimer de conteúdo não oficial.

## Escopo incluído

- telas/estados de login, vínculo, vínculo existente, troca e unlink;
- controllers/minimal endpoints finos para read/link/replace/unlink/delete;
- paths `/api/v1/me/player-link` e `/api/v1/me` conforme contrato da spec;
- DTOs JSON sem entidades de persistência expostas;
- confirmação antes de troca, desvinculação e exclusão;
- mensagem “Perfil público salvo — ownership não verificado”;
- loading, vazio, `401`, `403`, `404`, `409`, `429`, `503` e token expirado;
- reload/outro dispositivo recuperando vínculo após login;
- caminho funcional para apagar usuário e documentos próprios.
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

- `frontend/src/pages/Login.tsx`, `PlayerLink.tsx`;
- `src/Api/Endpoints/PlayerLinkEndpoints.cs` ou controllers equivalentes;
- `src/Application/PlayerLinks/`;
- `PlayerLinkRequest`, `PlayerLinkResponse`, `MapLinkError`;
- `DeleteOwnData` use case;
- disclaimer em layout/footer acessível.

## Passos de implementação

1. Expor leitura do estado vinculado do usuário autenticado.
2. Validar tag no client para UX e novamente no Application/API.
3. Executar lookup server-side e persistir somente após `resolved`.
4. Fazer replace somente após confirmação e preservar vínculo anterior em falha.
5. Implementar unlink e exclusão com autorização, confirmação e idempotência.
   Exigir `auth_time` recente e retornar `reauthentication_required` quando
   necessário.
6. Tratar `401`, reload, logout/login e novo dispositivo.
7. Exibir `public_profile`/`unverified` sem linguagem de ownership.
8. Adicionar disclaimer legível de fan content não oficial.

## Testes e comandos de validação

```text
dotnet test --filter Category=Api
npm run test:unit
npm run build
```

Cobrir ausência inicial, sucesso, tag inválida, perfil inexistente, provider
indisponível, rate limit, replace falho, unlink repetido e delete repetido.

## Definição de pronto

- usuário autentica e informa tag uma vez;
- API retorna somente dados do usuário autenticado;
- reload e outro dispositivo recuperam vínculo após login;
- replace valida nova tag antes de substituir antiga;
- falha externa preserva vínculo anterior;
- unlink e exclusão funcionam sem reentrada indevida;
- UI nunca diz que usuário possui ou controla perfil;
- estados vazios/erro/loading são acessíveis e build passa.
- OpenAPI documenta request/response, authn/authz, ProblemDetails e status;
- exclusão exige reautenticação recente e remove somente dados CrownPilot.

## Riscos e cuidados

- Não expor nome, coleção ou resposta raw para “enriquecer” tela.
- Não confundir exclusão local com exclusão no provider externo.
- Não enviar Player Tag real para analytics/logs client-side.
- Não criar promessa de sync a partir de validação única.
- Não colocar autorização crítica somente no React.
