# 002-04 — Implementar Google Sign-In e Firebase bearer

- **Ticker:** `002`
- **Número:** `04`
- **Status:** `completed`

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

## Implementação e evidências

- **Status:** `completed`
- **Arquivos alterados:** `.env.example`, `frontend/.env.example`, `frontend/package.json`,
  `frontend/package-lock.json`, `frontend/src/App.css`, `frontend/src/App.tsx`,
  `frontend/src/main.tsx`, `frontend/src/api/client.ts`,
  `frontend/src/api/client.test.ts`, `frontend/src/auth/firebase.ts`,
  `frontend/src/auth/AuthContext.ts`, `frontend/src/auth/AuthProvider.tsx`,
  `frontend/src/auth/AuthProvider.test.tsx`, `frontend/src/auth/useAuth.ts`,
  `src/Api/Authentication/ContractAuthentication.cs`, `src/Api/Program.cs`,
  `src/Api/Configuration/RuntimeOptions.cs`,
  `src/Api/Configuration/RuntimeOptionsValidator.cs`, `src/Api/Api.csproj`,
  `src/Application/Identity/AuthContracts.cs`,
  `src/Application/Identity/FirebaseTokenContracts.cs`,
  `src/Infrastructure/Authentication/FirebaseAuthenticationOptions.cs`,
  `src/Infrastructure/Authentication/FirebaseTokenVerifier.cs`,
  `src/Infrastructure/DependencyInjection.cs`,
  `src/Infrastructure/Infrastructure.csproj`, `tests/Api/FirebaseAuthenticationTests.cs`,
  `tests/Api/RuntimeOptionsTests.cs`, `tests/Application/BoundaryTests.cs`.
- **Decisões:** Firebase Admin SDK concentra assinatura, issuer, audience,
  expiração, `sub`, `kid` e rotação; adapter aplica allowlist de project/issuer.
  Render pode fornecer `FIREBASE_ADMIN_PROJECT_ID`,
  `FIREBASE_ADMIN_CLIENT_EMAIL` e `FIREBASE_ADMIN_PRIVATE_KEY` como secrets;
  chave privada aceita newline literal ou escapado e nunca é exposta ao frontend.
  Desenvolvimento local usa .NET User Secrets com mesmo naming do Render; não
  exige export manual a cada execução.
  Fixtures contratuais permanecem somente em Local. Application recebeu somente
  ports/fake de `EnsureCrownPilotUser` e contrato do verifier; handler não resolve
  nem persiste identidade. Claims do provider não viram permissões da aplicação.
  Frontend usa SDK para lifecycle/refresh/logout e envia apenas ID Token bearer,
  com `credentials: omit`, URL explícita em builds de produção e HTTPS fora de
  hosts locais em desenvolvimento.
- **Desvios:** testes de API usam `IFirebaseTokenVerifier` fake para não chamar
  Firebase real; execução de Emulator/Google real fica para staging/handoff em
  `002-12`, conforme escopo.
- **Comandos executados e resultados:**
  - `dotnet build CrownPilot.sln --configuration Release` — passou, 0 warnings,
    0 errors.
  - `dotnet test CrownPilot.sln --configuration Release` — passou, 50 testes.
  - `dotnet test CrownPilot.sln --configuration Release --filter Category=Authentication` —
    passou, 16 testes de API.
  - `dotnet test CrownPilot.sln --configuration Release --filter Category=Authorization` —
    passou, 1 teste de autorização de API.
  - `npm ci --prefix frontend` — instalação reproduzível concluída; npm reportou
    4 vulnerabilidades high transientes, sem aplicar `audit fix --force`.
  - `npm run lint --prefix frontend` — passou com `--max-warnings=0`.
  - `npm run typecheck --prefix frontend` — passou.
  - `npm run test:unit --prefix frontend` — passou, 3 arquivos/8 testes.
  - `npm run build --prefix frontend` — passou; bundle gerado sem credenciais
    server-side, refresh token ou service account.
  - `npm run openapi:check --prefix frontend` — passou; contrato gerado validado.
  - `npm audit --prefix frontend --omit=dev --audit-level=high` — bloqueado por
    4 vulnerabilidades high transitivas em `@grpc/grpc-js`; correção automática
    exige downgrade breaking do Firebase e não foi aplicada.
  - `rg -n -i 'service.?account|private.?key|refresh.?token|database.?password|firebase-admin|GoogleCredential' frontend/dist` —
    nenhuma ocorrência.
- **Riscos residuais:** rotação de chaves e login Google real dependem de
  Emulator/Staging, fora da execução local desta task. Vulnerabilidades high
  transientes do grafo npm precisam triagem antes dos gates de release.
