# 002-04 — Implementar Firebase Auth Google e bearer

- **Ticker:** `002`
- **Número:** `04`
- **Status:** `pending`

## Objetivo e resultado esperado

Implementar Google Sign-In no React e validação server-side do Firebase ID Token
na API ASP.NET Core. O resultado é um `AuthContext` seguro, sem sessão cookie,
UID arbitrário ou credencial Supercell.

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
- middleware/handler de autenticação ASP.NET Core em Infrastructure;
- integração suportada para validação Firebase, sem JWT manual;
- `AuthenticatedSubject` e resolução para `CrownPilotUserId` via abstração;
- caso de uso idempotente `EnsureCrownPilotUser`, com tratamento de corrida no
  primeiro login e falha de banco sem autenticação parcial;
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

- `002-01`, `002-02` e `002-03`;
- Firebase Emulator/fixtures e projetos por ambiente;
- configuração de Google provider em Staging, fora do Git.

## Arquivos e símbolos prováveis

- `frontend/src/auth/firebase.ts`, `AuthProvider.tsx`;
- `src/Infrastructure/Authentication/FirebaseAuthentication.cs`;
- `src/Api/Program.cs`, `AuthenticationOptions`, `AuthContext`;
- `src/Application/Identity/IUserIdentityResolver.cs`;
- endpoints/controllers de sessão mínima e testes de API/auth.

## Passos de implementação

1. Configurar Firebase Web SDK e Google provider por ambiente.
2. Implementar login/refresh/logout sem persistir token manualmente além do SDK.
3. Enviar somente ID Token bearer por HTTPS às rotas da API.
4. Configurar validação suportada na Infrastructure com project ID/issuer allowlist.
5. Derivar Firebase UID do token verificado e resolver ID interno via abstração;
   rejeitar UID em body/query/header.
6. Retornar `401` sem detalhes internos e manter autorização em políticas/casos
   de uso posteriores.
7. Provar login em Staging; usar Emulator/fixtures nos testes locais/CI.

## Testes e comandos de validação

```text
dotnet test --filter Category=Authentication
dotnet test --filter Category=Authorization
npm run lint
npm run typecheck
npm run build
```

Cobrir token ausente, inválido, expirado, issuer/audience/project incorretos,
assinatura/`kid` rotacionado, `sub` vazio, usuário A/B e logout. Verificar bundle
sem service account, database password, token externo ou código server-only.

## Definição de pronto

- Google Sign-In funciona no projeto Firebase de Staging;
- API aceita somente Firebase ID Token verificável em bearer;
- UID é derivado server-side e não é chave de domínio;
- authentication não concede autorização sobre usuário B;
- outro dispositivo recupera a mesma identidade após login;
- API não usa cookie/sessão paralela nesta fase;
- `401`, CORS e configuração por ambiente são testados;
- bundle não contém secrets nem refresh token enviado à API.

## Riscos e cuidados

- Não aceitar issuer/audience de projeto diferente.
- Não implementar validação criptográfica manual.
- Não enviar Google token, refresh token ou Firebase token ao provider Clash Royale.
- Não usar variáveis públicas para service account ou senha de banco.
- Não persistir e-mail/nome sem requisito de domínio.
