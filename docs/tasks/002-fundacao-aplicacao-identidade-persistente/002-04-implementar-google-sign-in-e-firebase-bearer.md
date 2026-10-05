# 002-04 — Implementar Google Sign-In e Firebase bearer

- **Ticker:** `002`
- **Número:** `04`
- **Status:** `pending`

## Objetivo e resultado esperado

Implementar Google Sign-In no React e validação server-side do Firebase ID Token
na API ASP.NET Core. O resultado é um `AuthenticatedSubject` seguro e um bearer
validado, sem sessão cookie, UID arbitrário ou credencial Supercell.

## Requisitos cobertos

- Firebase Authentication/Google Sign-In;
- refresh/logout no SDK frontend;
- Firebase ID Token em `Authorization: Bearer`;
- validação de assinatura, issuer, audience, expiração, `sub`, `kid` e rotação;
- identificação CrownPilot separada de Firebase UID;
- authentication separada de authorization;
- `401` consistente e ausência de secrets no bundle.

## Escopo incluído

- cliente Firebase Web configurado por ambiente;
- login, refresh, logout e estados loading/error;
- middleware/handler de autenticação ASP.NET Core em API, com verifier Firebase
  atrás de adapter na Infrastructure;
- Firebase Admin SDK para .NET ou integração oficial equivalente, sem JWT manual;
- `AuthenticatedSubject` e contrato de resolução para `CrownPilotUserId` via
  abstração Application;
- contrato/fake de `EnsureCrownPilotUser`, sem conectar EF/RLS nesta task;
- configuração de project ID/issuer por ambiente;
- CORS exato e envio HTTPS do bearer;
- testes de token, rotação, project errado e UID adulterado.

## Escopo excluído

- Player Tag, lookup ou persistência final;
- cookies, sessão server-side ou dois esquemas simultâneos;
- custom claims administrativas, ownership ou billing;
- autenticação Supercell;
- envio de refresh token ao backend.

## Dependências

- `002-01-bootstrap-toolchain.md` e
  `002-02-estabelecer-boundaries-contrato-base-e-ambientes.md`;
- Firebase Emulator/fixtures e projetos por ambiente, configurados nesta task;
- configuração de Google provider em Staging, fora do Git.

## Arquivos e símbolos prováveis

- `frontend/src/auth/firebase.ts`, `AuthProvider.tsx`;
- `src/Infrastructure/Authentication/FirebaseTokenVerifier.cs`;
- `src/Api/Program.cs`, `AuthenticationOptions`, `AuthContext`;
- `src/Application/Identity/IUserIdentityResolver.cs`;
- componentes de authentication/authorization e testes de API/auth.

## Passos de implementação

1. Configurar Firebase Web SDK, Emulator/fixtures e Google provider por ambiente.
2. Implementar login/refresh/logout sem persistir token manualmente além do SDK.
3. Enviar somente ID Token bearer por HTTPS às rotas da API.
4. Configurar validação suportada na Infrastructure com project ID/issuer allowlist
   e credencial server-only; não persistir ou criar usuário no handler.
5. Derivar Firebase UID do token verificado e resolver ID interno via abstração;
   rejeitar UID em body/query/header.
6. Retornar `401` sem detalhes internos e manter autorização em políticas/casos
   de uso posteriores.
7. Usar Emulator/fixtures nos testes locais/CI; Google real fica para Task
   `002-12-validar-staging-e2e-smoke-e-handoff.md`.

## Testes e comandos de validação

```text
dotnet test --filter Category=Authentication
dotnet test --filter Category=Authorization
npm run lint
npm run typecheck
npm run build
```

Cobrir token ausente, inválido, expirado, issuer incorreto, audience incorreta,
projeto Firebase incorreto, assinatura inválida, `kid` rotacionado, `sub` vazio e
logout. Testar separadamente authentication válida sem authorization. Verificar
bundle sem service account, database password, token externo ou código server-only.

## Definição de pronto

- Google Sign-In usa Firebase Web SDK com configuração pública por ambiente;
- API aceita somente Firebase ID Token verificável em bearer;
- handler de authentication não escreve no banco nem cria identidade como efeito
  colateral;
- UID é derivado server-side e não é chave de domínio;
- authentication não concede autorização sobre usuário B;
- resolução/criação de identidade fica disponível como port/fake para
  `002-06-modelar-persistencia-repositories-e-rls.md`;
- API não usa cookie/sessão paralela nesta fase;
- `401`, CORS e configuração por ambiente são testados;
- bundle não contém secrets nem refresh token enviado à API.

## Riscos e cuidados

- Não aceitar issuer/audience de projeto diferente.
- Não implementar validação criptográfica manual.
- Não enviar Google token, refresh token ou Firebase token ao provider Clash Royale.
- Não usar variáveis públicas para service account ou senha de banco.
- Não persistir e-mail/nome sem requisito de domínio.
