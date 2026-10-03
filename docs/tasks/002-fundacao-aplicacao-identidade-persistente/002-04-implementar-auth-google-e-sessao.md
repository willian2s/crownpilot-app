# 002-04 — Implementar Auth Google e sessão

- **Ticker:** `002`
- **Número:** `04`
- **Status:** `pending`

## Objetivo e resultado esperado

Implementar identidade CrownPilot com Google Sign-In e boundary de sessão que
permita ao backend reconhecer o usuário sem receber UID arbitrário ou credencial
Supercell.

## Requisitos cobertos

- sessão Supabase Auth/Laravel;
- recuperação em outro dispositivo;
- verificação server-side do access JWT via JWKS e sessão Laravel segura;
- autenticação separada de autorização;
- falhas de Auth sem exposição de secrets;
- Web SDK limitado a Auth.

## Escopo incluído

- inicialização `@supabase/supabase-js` por ambiente;
- Google OAuth, refresh/logout e observação de Auth state;
- envio de access JWT ao endpoint de troca Laravel;
- verificação por adapter PHP Supabase/JWKS e criação de sessão/`AuthContext` com
  UUID `sub` derivado;
- respostas `401` consistentes;
- proteção para não inicializar Data API/PostgREST para dados CrownPilot;
- tratamento de popup/redirect e estados loading/error.

## Escopo excluído

- Player Tag, lookup ou persistência de vínculo;
- custom claims, roles administrativas ou ownership;
- autenticação Supercell;
- billing ou múltiplos providers.

## Dependências

- `002-01`, `002-02` e `002-03`;
- projetos, refs e credenciais Supabase por ambiente;
- Google provider autorizado em staging.

## Arquivos e símbolos prováveis

- `resources/js/supabase/auth-client.ts`;
- `app/Adapters/Supabase/AuthJwtVerifier.php`;
- `app/Application/Auth/ExchangeSupabaseJwt.php`;
- `AuthContext`, `requireAuth`, `UnauthorizedError`;
- `app/Http/Controllers/AuthController.php`, `routes/web.php`;
- `resources/js/Pages/Login.tsx` ou equivalente Inertia;
- fixtures/mocks de Auth e testes de bundle.

## Passos de implementação

1. Configurar Google provider e redirect allowlist no Supabase de cada ambiente.
2. Implementar Google OAuth, refresh/logout e estado de carregamento.
3. Enviar access JWT somente para endpoints CrownPilot via HTTPS.
4. Verificar JWT no PHP via JWKS, validar claims, derivar UUID `sub` e criar sessão
   Laravel segura.
5. Fazer handlers rejeitarem UID/body claims e tokens inválidos.
6. Testar reload e novo dispositivo como novo cliente com mesma identidade.

## Testes e comandos de validação

```text
composer run test:unit
composer run test:integration
npm run lint
npm run typecheck
npm run build
```

Cobrir token ausente, inválido, expirado, usuário diferente e sessão válida.
Verificar bundle para ausência de adapter PHP/Supabase server-side e secrets.

## Definição de pronto

- usuário entra e sai com Google em staging;
- backend aceita somente Supabase access JWT verificável e cria sessão segura;
- UUID `sub` é sempre derivado server-side;
- outro dispositivo recupera a mesma identidade após login;
- browser não acessa Data API/PostgREST nem recebe secret server-side;
- `401` não revela detalhes internos;
- testes unit/integration/build passam.

## Riscos e cuidados

- Não persistir e-mail/nome Google sem necessidade de domínio.
- Não tratar login como autorização para ler outro UID.
- Não enviar Google ID token ou Supabase JWT ao provider Clash Royale.
- Não usar `VITE_` para service role, database password, JWT secret ou API token.
