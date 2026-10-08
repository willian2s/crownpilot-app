# ADR 005 — Go, React/Vite, Firebase Auth e PostgreSQL/Supabase

- **Status:** `accepted`
- **Data:** 2026-10-08
- **Substitui:** ADR 004 somente nas decisões de linguagem, runtime, organização
  do backend, autenticação server-side, persistência, OpenAPI e verificação
  operacional.

## Contexto

O baseline executável atual usa ASP.NET Core, quatro projetos .NET, Firebase Admin
SDK para .NET, EF Core/Npgsql, Docker e testes xUnit. As tasks `002-01` a `002-04`
foram concluídas como histórico desse baseline. Não há dados legados, usuários
persistidos ou deploy que exijam dual-write; a migração pode ser substituição
controlada.

A spec 002 permanece fonte normativa dos requisitos funcionais. Esta ADR decide
somente a stack e seus boundaries. O plano de migração está dividido em três
tasks, cada uma com gate próprio:

- `002-13` — bootstrap, HTTP, configuração e OpenAPI;
- `002-14` — autenticação Firebase;
- `002-15` — persistência, cutover e remoção do .NET.

Com esta ADR aceita, `002-13` está liberada para execução. `002-14` e `002-15`
permanecem `pending` até suas dependências e gates próprios serem concluídos.

## Decisões

### Backend e módulos

O backend será um único processo Go, um único deploy e uma imagem OCI. O layout
proposto é:

```text
cmd/crownpilot-api/              composição e lifecycle do processo
cmd/crownpilot-migrate/          job separado de migrations
internal/httpapi/                HTTP, middleware, DTOs e ProblemDetails
internal/config/                 carga e validação de configuração
internal/identity/               identidade, casos de uso e ports
internal/playerlink/             vínculo, casos de uso e ports
internal/platform/firebase/      adapter Firebase Admin Go
internal/platform/postgres/      pgx, sqlc, transações e RLS
internal/platform/clashroyale/   adapter do provider externo
internal/observability/          slog, request ID, métricas e redaction
api/openapi/                     fonte única OpenAPI e pacote de embed
```

`internal/` protege os módulos. `playerlink` pode importar contratos públicos do
módulo `identity`; `identity` não importa `playerlink`. Ambos podem depender
somente de contratos internos estáveis e biblioteca padrão, nunca de HTTP,
Firebase, PostgreSQL ou adapters concretos. Adapters são montados somente em
`cmd`.

Persistência usará `pgxpool` e `sqlc`, com SQL revisável e migrations SQL
versionadas por `goose` em job separado. SQL de segurança continua separado para RLS,
roles, grants, extensões e objetos de plataforma; somente esse SQL habilita RLS e
cria policies. Supabase continua hospedagem PostgreSQL, não backend da aplicação.

O toolchain Go será fixado em `1.27.2`. A compatibilidade da versão escolhida
com Firebase Admin Go e ferramentas pinadas será verificada no bootstrap.

### Roteamento HTTP

Usar somente `net/http` com `http.ServeMux`, incluindo padrões de método e path
disponíveis desde Go 1.22. Middlewares serão funções que envolvem
`http.Handler`.

Motivo: há poucas rotas, `ServeMux` é suficiente desde Go 1.22 e a escolha tem
valor didático ao tornar matching, composição e lifecycle explícitos. `chi` foi
avaliado como alternativa idiomática para middleware e rotas aninhadas, mas não
será dependência do primeiro corte. Frameworks maiores permanecem fora.

A porta local canônica permanece `5080`, igual ao `launchSettings.json` do .NET.
`5089` continua somente override isolado dos scripts de smoke; em hosting, o
processo respeita `PORT` fornecida pelo ambiente.

### OpenAPI

OpenAPI será spec-first. O arquivo canônico será `api/openapi/v1.json`, única
fonte do contrato. `oapi-codegen` gerará os tipos/interfaces Go e
`openapi-typescript` gerará os tipos do frontend. Artefatos gerados serão
versionados no Git.

`/openapi/v1.json` será servido pelo próprio processo a partir da fonte JSON
embutida por um pacote Go sob `api/openapi`, usando `go:embed`. Não haverá
`cmd/crownpilot-openapi`, conversor próprio ou segundo JSON gerado: servir a
fonte canônica inalterada elimina cópia, binário e possibilidade de drift na
inicialização. `/docs`, quando habilitado pelo ambiente, consumirá essa rota.

O CI executará os geradores pinados e falhará quando `git diff` detectar drift
entre a fonte, os artefatos Go e os tipos frontend. O contrato funcional continua
referenciado na spec 002.

`/docs` continuará usando Swagger UI somente em Local e Staging, apontando para
`/openapi/v1.json`. Swagger UI é visualização, não fonte nem substituto do gate
`openapi-check.mjs`, que permanece no CI.

### Credenciais Firebase

Usar Firebase Admin Go SDK. A inicialização montará credencial de service account
em memória a partir exatamente destas variáveis de ambiente:

```text
FIREBASE_ADMIN_PROJECT_ID
FIREBASE_ADMIN_CLIENT_EMAIL
FIREBASE_ADMIN_PRIVATE_KEY
```

Não haverá arquivo JSON de service account no repositório, na imagem OCI ou no
hosting. A opção será construída com a API atual de
`google.golang.org/api/option`:

```text
option.WithAuthCredentialsJSON(option.ServiceAccount, inMemoryJSON)
```

`option.WithCredentialsJSON` está depreciada por aceitar configuração sem
validação de tipo; `option.WithServiceAccountFile` também está depreciada.
`WithAuthCredentialsJSON` é a alternativa específica de tipo documentada
atualmente.

`FirebaseAuthenticationOptions` manterá os três campos acima. `.env.example`
manterá as três chaves vazias, somente como documentação local; valores reais
ficam em `.env` local ignorado pelo Git ou em secrets separados por ambiente no
hosting.

Regras obrigatórias:

- tratar `\n` escapado em `FIREBASE_ADMIN_PRIVATE_KEY` como newline antes da
  montagem;
- validar configuração tudo-ou-nada antes de inicializar o SDK;
- fora de Local, ausência de qualquer variável encerra o processo antes do
  listener;
- em Local com Firebase Emulator, as três variáveis podem ficar vazias;
- valores nunca aparecem em logs, erros, bundle ou artefatos versionados.

### Testes de arquitetura

`depguard` do `golangci-lint` será a forma principal de bloquear imports proibidos,
incluindo a regra `identity` não importar `playerlink`. `go list` permanecerá
somente onde provar algo que `depguard` não cobre, como composição do binário,
dependências transitivas relevantes ou confirmação de que adapters concretos só
entram pelo composition root. Não manter dois testes equivalentes.

### Decisões dependentes

Registro das perguntas numeradas da versão anterior:

1. **Versão Go e hosting:** Go `1.27.2` decidida. A autenticação administrativa
   usa as três variáveis de secret runtime já decididas; o hosting não precisa
   fornecer ADC ou workload identity. ADC/workload identity permanece somente como
   alternativa avaliada, não como requisito ou pendência.
2. **Router:** decidida; somente `net/http` + `ServeMux`, com `chi` apenas como
   alternativa avaliada.
3. **Migration tool e rollback:** decididos; usar `goose` em job separado. Em
   Preview, Staging e Production as migrations são forward-only: correções usam
   nova migration e `down` nunca roda em pipeline ou job de deploy. Migrations
   `down` existem somente para desenvolvimento local. Mudanças destrutivas usam
   expand/contract em releases separadas, e cada migration precisa ser compatível
   com a versão anterior da aplicação. Rollback da aplicação reimplanta a imagem
   anterior sobre o schema atual; rollback de dados usa backup/restore conforme
   runbook operacional. Não há rollback automático do banco por `down`.
4. **OpenAPI:** decidida; spec-first JSON, geradores pinados, artefatos no Git e
   CI contra drift.
5. **Firebase credential boundary:** decidida; três variáveis, credencial em
   memória, sem JSON em arquivo/imagem/hosting e com `WithAuthCredentialsJSON`.
6. **Revogação:** a política seletiva está decidida. Rotas autenticadas usam
   `VerifyIDToken` por padrão, sem checagem de revogação. Operações sensíveis usam
   `VerifyIDTokenAndCheckRevoked`; a lista inicial contém `DELETE /api/v1/me`.
   Billing, admin e novas ações destrutivas só entram nessa lista por decisão
   registrada. A checagem é declarada na composição de cada rota, preservando
   `net/http` + `ServeMux`, por meio de um wrapper/middleware específico de auth
   sensível, separado do wrapper padrão. Token revogado retorna `401` indistinguível
   das demais falhas de autenticação.

   Timeout ou indisponibilidade de comunicação com o Firebase durante a checagem
   sensível falha fechado com `503` e o novo code `authentication_unavailable`,
   usando ProblemDetails genérico sem detalhe do provider. `provider_unavailable`
   não é reutilizado: na spec 002, ele representa estado do lookup do Clash
   Royale. O contrato OpenAPI planejado para `DELETE /api/v1/me` documentará
   `503` com `code: authentication_unavailable`. A exigência de `auth_time`
   presente e recente permanece como proteção adicional conforme a spec 002.
7. **Nome e ordem SDD:** decidida; migração é `002-13` a `002-15`, e `002-01` a
   `002-04` permanecem histórico concluído.
8. **Porta local canônica:** decidida; usar `5080`, igual ao .NET. `5089` fica
   override de smoke e `PORT` continua autoridade em hosting.
9. **Destino de `/docs`:** decidida; manter Swagger UI em Local/Staging consumindo
   `/openapi/v1.json`; Preview/Production permanecem desabilitados.
10. **Pooler Supabase:** decidido; o projeto Supabase exclusivo de desenvolvimento
    e todos os ambientes hospedados usam Session pooler. Conexão direta fica
    restrita ao PostgreSQL local do Docker Compose, usado para testes/reset
    descartáveis, e não é usada no projeto Supabase dev nem em hosting. Transaction
    pooler não será usado neste corte; se
    adotado no futuro, exige protocolo simples sem prepared statements no `pgx` e
    prova do contexto RLS dentro da transação. O limite de conexões configurado no
    `pgxpool` precisa caber no limite de conexões do Session pooler do plano usado.
11. **Tooling frontend:** decidida; Swagger UI permanece como visualização e
    `openapi-check.mjs` permanece gate de CI; não haverá segundo contrato.
12. **Região e data residency:** API no Render, região Virgínia (`us-east`), e
    PostgreSQL no Supabase, região `us-east-1` (Northern Virginia). A mesma região
    reduz a latência das várias idas e voltas de cada operação (`BEGIN`, `SET LOCAL`
    de RLS, queries e `COMMIT`) e a costa leste dos EUA oferece a menor
    latência para usuários no Brasil entre as regiões disponíveis no Render.
    `sa-east-1` foi descartada porque o Render não tem região na América do Sul.
    Oregon e Ohio foram avaliadas; Ohio (`us-east-2`), recomendada pelo Supabase,
    era alternativa válida com latência ligeiramente maior. Dados pessoais ficam
    fora do Brasil; a transferência internacional precisa constar na política de
    privacidade antes de dados reais. Backups e egress continuam pendentes e são
    classificados abaixo, sem alterar o status da ADR.
13. **RLS automático do Supabase:** decidido; permanece desabilitado no projeto.
    `database/security/` é a única fonte autorizada a habilitar RLS e criar
    policies. Um teste obrigatório consulta `pg_class.relrowsecurity` e falha se
    qualquer tabela do schema `crownpilot` estiver sem RLS ativo. O teste roda no
    CI contra PostgreSQL local e é gate obrigatório antes de Staging.

### Pendências remanescentes e classificação

Estas classificações não alteram o status `accepted` da ADR:

| Pendência | Classificação | Gate obrigatório e justificativa |
|---|---|---|
| Backups e restauração | Não bloqueia a ADR; gate antes de dados reais | `002-12` deve fechar retenção, backup, restore testado e runbook operacional antes do handoff/qualquer dado real. É garantia operacional de recuperação, não escolha de stack ou boundary. |
| Egress/provider live | Não bloqueia a ADR; gate antes de dados reais | `002-12` deve validar egress permitido, rota escolhida e comportamento controlado em Staging antes de dados reais. `002-15` pode registrar o boundary/configuração, mas não fecha a prova live; o adapter mantém egress substituível. Isso preserva os gates de API data/egress da Fase 001. |

Backups e egress não estão confirmados como resolvidos nesta ADR. Sem os gates
correspondentes, não liberar dados reais nem tratar a Fase 003 como liberada.

## Opções avaliadas

### Roteamento

| Opção | Avaliação |
|---|---|
| `net/http` + `ServeMux` | Escolhida: stdlib, poucas rotas, suficiente desde Go 1.22 e didática. |
| `chi` | Alternativa avaliada: composição idiomática e matching útil, mas dependência adicional sem necessidade atual. |
| Gin/Echo | Rejeitada: abstração, acoplamento e convenções excessivos para o escopo. |

### OpenAPI

| Opção | Avaliação |
|---|---|
| code-first | Rejeitada: contrato nasce dos handlers e aumenta risco de drift sem uma fonte revisável. |
| spec-first JSON + `oapi-codegen` + `openapi-typescript` | Escolhida: uma fonte explícita, artefatos derivados versionados e drift verificável no CI. |
| spec-first YAML convertido por binário próprio | Rejeitada: conversor e artefato duplicam responsabilidade; fonte JSON embutida é mais simples. |
| annotations/reflection | Rejeitada: segunda linguagem no código e cobertura menos explícita. |

### Credenciais Firebase

| Opção | Avaliação |
|---|---|
| ADC/workload identity | Boa quando hosting fornece identidade; permanece alternativa avaliada, não requisito ou contrato desta aplicação. |
| Arquivo JSON de service account | Rejeitada: risco de cópia, montagem e exposição na imagem/hosting. |
| Três variáveis + JSON em memória | Escolhida: compatível com hosting por secrets, sem arquivo persistente e com validação explícita. |
| Verificador JWT próprio com Project ID e certificados públicos | Rejeitada: transfere criptografia, cache de `kid`, rotação e validação para CrownPilot. |

### Acesso a dados

| Opção | Avaliação |
|---|---|
| `pgx` + `sqlc` | Escolhida: SQL explícito, tipos gerados e transações/RLS visíveis. |
| `pgx` direto | Viável, mas repete mapeamentos e aumenta erro quando o domínio crescer. |
| ORM | Rejeitada: esconde SQL/RLS e cria ownership concorrente do schema. |

### Garantia de imports

| Opção | Avaliação |
|---|---|
| `depguard` como regra principal | Escolhida: lint no CI e regras declarativas próximas do restante do lint. |
| teste exclusivo com `go list` | Rejeitada como mecanismo principal: mais código de teste e duplicação de regras simples. Mantida apenas para lacunas de composição/transitividade. |

## Consequências

### Positivas

- backend, imagem e lifecycle menores e portáveis;
- `net/http`, `context`, `testing` e `slog` deixam boundaries e comportamento
  explícitos;
- contrato OpenAPI revisável, gerado para Go/frontend e servido sem binário
  auxiliar;
- credencial Firebase não depende de arquivo no filesystem;
- SQL visível facilita revisão de PostgreSQL, transações e RLS;
- ausência de dados legados permite cutover sem dual-write.

### Negativas

- perde-se integração pronta de DI, OpenAPI e migrations do .NET;
- manutenção combina Go, Node, PostgreSQL e geradores pinados;
- `ServeMux` e middleware exigem disciplina local;
- `pgx`/`sqlc` exigem disciplina de SQL e geração;
- Firebase Admin Go ainda exige credencial válida fora do Local com Emulator;
- a migração precisa provar paridade antes de remover artefatos .NET.

## Riscos

| Risco | Mitigação/aceite |
|---|---|
| Configuração Firebase parcial ou exposta | validação tudo-ou-nada antes do listener; newline escapado tratado; nenhum valor em logs ou artefatos. |
| Conta Firebase com permissões administrativas amplas | Risco aceito nesta ADR. Trocar por conta sem papéis IAM amplos é melhoria futura sem mudança de código. |
| Hosting não fornecer secrets corretamente | três secrets separados por ambiente; startup fail-closed fora de Local. ADC/workload identity não é dependência. |
| Rota sensível omitir revogação | política declarada por rota, wrapper específico para auth sensível, lista inicial registrada e testes de composição. |
| Firebase indisponível durante revogação sensível | falha fechada; `503` com `authentication_unavailable` em ProblemDetails genérico, sem detalhe do provider; `provider_unavailable` permanece exclusivo do lookup do Clash Royale. |
| Drift entre OpenAPI, handlers e frontend | fonte JSON única, geradores versionados e CI falhando em `git diff`. |
| `identity` e `playerlink` criarem dependência circular | `depguard` bloqueia `identity -> playerlink`; `playerlink -> identity` é a única direção permitida entre os dois. |
| RLS perder contexto no pool | transação com `SET LOCAL`, rollback/commit e testes de isolamento. |
| Pooler Supabase exceder limite de conexões | dimensionar `pgxpool` dentro do limite do Session pooler do plano; validar no projeto Supabase dev e em cada ambiente hospedado. |
| Pooler Supabase em modo transação com pgx | fora deste corte; se adotado, exige desativar prepared statements, protocolo simples e prova de contexto RLS na transação. |
| RLS automático ou policy ausente | auto-RLS permanece desabilitado; SQL de `database/security/` é fonte única; teste de `pg_class.relrowsecurity` bloqueia CI/Staging. |
| Migrations duplicarem schema ou quebrarem rollback | job separado, forward-only compartilhado, expand/contract destrutivo, compatibilidade com versão anterior e restore por runbook. |
| Fakes ou fixtures entrarem na imagem | fixtures somente em testes/Local; scan de bundle e imagem antes do cutover. |
| Região ou transferência internacional incompatível | Render e Supabase na Virgínia; registrar dados fora do Brasil na política de privacidade antes de dados reais. |
| Backups ou egress não comprovados | não bloqueiam stack ADR, mas são gates obrigatórios de `002-12` antes de dados reais; permanecerem pendentes mantém liberação bloqueada. |
| Corte Go quebrar scripts e frontend | manter interfaces públicas de scripts, smoke e `/api/v1`; gates verdes em cada task de migração. |
