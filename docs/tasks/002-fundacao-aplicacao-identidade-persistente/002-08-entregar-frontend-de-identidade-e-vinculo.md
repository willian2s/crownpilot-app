# 002-08 — Entregar frontend de identidade e vínculo

- **Ticker:** `002`
- **Número:** `08`
- **Status:** `pending`

Esta task permanece bloqueada até ADR 005 ser aprovada e `002-13` -> `002-14` ->
`002-15` concluírem gates verdes. O frontend permanece React/Vite, mas consome o
contrato/tipos Go definidos no ajuste de `002-15`.

## Objetivo e resultado esperado

Entregar experiência React/Vite para login, vínculo, consulta, troca,
desvinculação e exclusão de dados CrownPilot, consumindo somente API `/api/v1` e
comunicando honestamente `public_profile`/`unverified`.

## Requisitos cobertos

- Firebase Web SDK para login, refresh, logout e token bearer;
- API client tipado e separado de componentes;
- estados loading, empty, success e error;
- mensagens para `400`, `401`, `403`, `404`, `409`, `422` quando aplicável, `429`,
  `500`, `503` e token expirado;
- confirmação antes de replace, unlink e delete;
- reautenticação Google quando exclusão exigir `auth_time` recente;
- disclaimer de conteúdo não oficial;
- testes de componentes, API client e fluxos críticos sem backend real.

## Escopo incluído

- páginas/componentes de login, vínculo, vínculo existente, troca e unlink;
- hooks/services para sessão e chamadas autenticadas;
- models/types derivados do contrato OpenAPI, sem entidades de persistência;
- mensagens “Perfil público salvo — ownership não verificado” sem linguagem de
  ownership, controle ou exclusividade;
- reload/outro dispositivo recuperando vínculo após login;
- caminho funcional para apagar dados próprios sem apagar Firebase/Google;
- acessibilidade estrutural básica, foco e feedback de formulário.

## Escopo excluído

- revisão arquitetural/visual abrangente, que pertence à Task
  `002-09-revisar-arquitetura-frontend-e-ux-visual.md`;
- dashboard de coleção, Arena, sync ou snapshot;
- múltiplas tags, social, clãs, billing ou coaching;
- chamada direta do browser à Clash Royale API, PostgreSQL ou Supabase Data API.

## Dependências

- ADR 005 aprovada, `002-13-bootstrap-http-config-openapi-go.md`,
  `002-14-autenticacao-firebase-go.md`, `002-15-persistencia-cutover-remocao-dotnet.md`,
  `002-04-implementar-google-sign-in-e-firebase-bearer.md` como histórico,
  `002-07-implementar-casos-de-uso-e-api-v1.md` e frontend do `002-01`;
- contrato OpenAPI e ProblemDetails da Task
  `002-07-implementar-casos-de-uso-e-api-v1.md`;
- hostname fixo de Staging para validação posterior.

## Arquivos e símbolos prováveis

- `frontend/src/pages/Login.tsx`, `PlayerLink.tsx`;
- `frontend/src/auth/`, `frontend/src/api/`, `frontend/src/hooks/`;
- `frontend/src/models/` e testes de React Testing Library/Vitest;
- layout/footer com disclaimer acessível.

## Passos de implementação

1. Integrar estado de autenticação Firebase e API client bearer.
2. Implementar leitura e formulário de vínculo com validação client-side apenas
   para UX; servidor continua autoridade.
3. Implementar replace, unlink e delete com confirmação e reautenticação quando
   receber `reauthentication_required`.
4. Mapear ProblemDetails para mensagens úteis sem expor UID, tag, URL ou provider.
5. Cobrir reload, logout/login, novo dispositivo simulado e falhas controladas.
6. Entregar estados acessíveis e deixar inventário de arquitetura/UX para
   `002-09-revisar-arquitetura-frontend-e-ux-visual.md`.

## Testes e comandos de validação

```text
npm run test:unit --prefix frontend
npm run lint --prefix frontend
npm run typecheck --prefix frontend
npm run build --prefix frontend
```

Testar componentes importantes, auth state, API client, loading/empty/error/
success, `400`/`401`/`403`/`404`/`409`/`422` quando aplicável/`429`/`500`/`503`,
token expirado, confirmação destrutiva e reautenticação. Usar mocks e fixtures;
não chamar Firebase ou Clash Royale live na CI normal.

## Definição de pronto

- usuário autentica e informa tag uma vez;
- API client envia somente Firebase ID Token bearer;
- reload e outro dispositivo recuperam vínculo após login;
- replace valida nova tag no servidor antes de substituir antiga;
- falha externa preserva vínculo anterior;
- unlink e exclusão funcionam de forma explícita e idempotente;
- UI nunca afirma ownership;
- estados vazios, erro, loading e sucesso são acessíveis;
- testes frontend, typecheck, lint e build passam;
- nenhum secret, token de refresh, payload externo ou código server-only chega ao
  bundle.

## Riscos e cuidados

- Não colocar autorização crítica somente no React.
- Não logar Player Tag real no client-side analytics.
- Não transformar esta task em redesign sem inventário; revisão ocorre em
  `002-09-revisar-arquitetura-frontend-e-ux-visual.md`.
