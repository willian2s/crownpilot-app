# 002-14 — Autenticação Firebase no backend Go

- **Ticker:** `002`
- **Número:** `14`
- **Status:** `pending`

## Objetivo e resultado esperado

Portar autenticação server-side para Firebase Admin Go SDK, mantendo Firebase UID
como subject externo e `CrownPilotUserID` como identidade interna. A task entrega
credenciais em memória, validação fail-closed e testes de authentication separados
de authorization.

## Requisitos cobertos

- Firebase Admin Go SDK oficial, sem JWT manual;
- `FirebaseAuthenticationOptions` com
  `FIREBASE_ADMIN_PROJECT_ID`, `FIREBASE_ADMIN_CLIENT_EMAIL` e
  `FIREBASE_ADMIN_PRIVATE_KEY`;
- `.env.example` com as três chaves vazias, sem valores reais;
- JSON de service account construído em memória;
- tratamento de `\n` escapado em private key;
- validação tudo-ou-nada antes do listener;
- Local com Emulator permitindo as três variáveis vazias;
- ausência de qualquer variável fora de Local encerrando o processo;
- nenhum valor secreto em logs, bundle, imagem ou Git;
- bearer Firebase, `AuthenticatedSubject` e boundary de autorização.

## Escopo incluído

- adapter `internal/platform/firebase` e composição no `cmd/crownpilot-api`;
- inicialização com `firebase.NewApp` e
  `option.WithAuthCredentialsJSON(option.ServiceAccount, inMemoryJSON)`;
- usar API atual de `google.golang.org/api/option`; não usar a depreciada
  `option.WithCredentialsJSON` nem credencial por arquivo;
- normalizar newline escapado somente na representação em memória;
- configurar project ID/issuer por ambiente e defesa em profundidade após
  verificação oficial;
- manter fixtures/emulator somente em Local e testes;
- atualizar mapa de equivalência de authentication e configuração;
- preservar comportamento funcional de authentication/authorization definido na
  spec 002, sem colocar requisitos de domínio no middleware.

## Escopo excluído

- criação de usuário em banco, repositories, RLS ou persistência final;
- casos de uso de vínculo e delete;
- revogação global por request, que permanece decisão de segurança posterior;
- arquivo JSON, volume ou secret file de service account.

## Dependências

- `002-13-bootstrap-http-config-openapi-go.md` com startup e configuração base;
- [ADR 005](../../decisions/005-go-react-vite-firebase-postgresql.md);
- ADR 005 aprovada;
- perguntas pendentes fechadas ou adiadas formalmente dentro da ADR aprovada;
  aprovação da ADR é pré-condição obrigatória;
- Firebase Emulator para Local ou fixtures controladas de teste;
- projetos e secrets reais de Staging somente fora do Git.

## Mapa de equivalência desta task

| Baseline .NET | Destino Go da migração |
|---|---|
| `FirebaseAuthenticationOptions.cs` | opções Go mantendo os três campos `FIREBASE_ADMIN_*`. |
| `FirebaseTokenVerifier.cs` | adapter `internal/platform/firebase` usando Admin Go. |
| `ContractAuthentication.cs` | middleware bearer que produz `AuthenticatedSubject`. |
| `AuthContracts.cs`, `FirebaseTokenContracts.cs` | contratos em `internal/identity`; somente subject verificado cruza boundary. |
| `.env.example` e User Secrets | chaves vazias no exemplo local; `.env` ignorado localmente; secrets separados por ambiente no hosting. |
| testes de Firebase/auth | testes Go de middleware/adapter e contract tests HTTP; fixtures fora do binário de produção. |

## Passos de implementação

1. Validar presença das três variáveis como conjunto, sem registrar valores.
2. Montar JSON de service account em memória, convertendo `\n` para newline.
3. Inicializar Admin Go antes do listener fora de Local com Emulator.
4. Verificar token com SDK oficial, conferir project/issuer e produzir somente
   `AuthenticatedSubject`.
5. Manter authentication distinta de authorization e não resolver/persistir
   usuário no middleware.
6. Cobrir ausência, parcialidade, newline, projeto errado, issuer, expiração,
   assinatura, `kid`, `sub` e UID arbitrário.
7. Executar gate verde e registrar evidências antes de `002-15`.

## Gate verde obrigatório

```text
go test ./...
go test -race ./...
go vet ./...
golangci-lint run
npm run typecheck --prefix frontend
npm run lint --prefix frontend
npm run test:unit --prefix frontend
npm run build --prefix frontend
npm run openapi:check --prefix frontend
```

Também deve haver teste de processo comprovando que configuração ausente fora de
Local falha antes de bindar a porta. Emulator/Google real de Staging pode ser
pré-condição operacional posterior, mas não pode ser simulado como gate verde.

## Definição de pronto

- Admin Go usa credencial montada em memória das três variáveis;
- `WithCredentialsJSON` não é usado; `WithAuthCredentialsJSON` é usado com tipo
  `option.ServiceAccount`;
- Local com Emulator aceita credenciais vazias; outros ambientes falham fechado
  se qualquer campo estiver ausente;
- `FirebaseAuthenticationOptions` e `.env.example` preservam o mapa decidido;
- nenhum segredo aparece em logs, erros, frontend, imagem ou arquivo versionado;
- `401` permanece indistinguível para token ausente, inválido, expirado, issuer,
  audience, projeto, assinatura ou `kid` incorretos;
- UID verificado não é identidade de domínio nem entrada confiável do request;
- fixtures ficam somente em testes/Local e não são compiladas no binário de
  produção;
- testes de authentication e authorization permanecem separados.

## Riscos e cuidados

- Não aceitar configuração parcial nem fallback silencioso para ADC fora de Local;
  ADC/workload identity é somente alternativa avaliada, não dependência do hosting.
- Não logar JSON montado, private key, client email ou project ID secreto.
- Não conceder autorização a partir de claims arbitrárias do provider.
- Conta administrativa ampla é risco aceito na ADR 005; least privilege fica como
  melhoria futura sem mudança de código.
