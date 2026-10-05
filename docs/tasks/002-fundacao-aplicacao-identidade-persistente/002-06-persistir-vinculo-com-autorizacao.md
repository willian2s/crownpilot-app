# 002-06 — Persistir vínculo com autorização

- **Ticker:** `002`
- **Número:** `06`
- **Status:** `pending`

## Objetivo e resultado esperado

Persistir somente usuário CrownPilot e vínculo primário de perfil público usando
EF Core + Npgsql, com ID interno, autorização por usuário autenticado e sem
criar snapshot de API.

## Requisitos cobertos

- PostgreSQL server-side hospedado no Supabase;
- `CrownPilotDbContext` e mappings EF Core;
- migrations EF Core para schema da aplicação;
- `firebase_uid` externo separado do UUID interno;
- usuário A isolado de B por authorization server-side;
- `public_profile` + `unverified`;
- replace protegido contra concorrência e falhas;
- RLS/grants como defesa adicional, sem duplicar schema.

## Escopo incluído

- entidade `CrownPilotUser` com UUID interno e `firebase_uid` único;
- entidade `PrimaryPlayerLink` com `user_id` como primary key/FK;
- repository baseado no `CrownPilotUserId` resolvido no backend;
- operações read, create/replace, unlink e delete idempotente;
- transaction/concurrency token/precondition para replace;
- timestamps, `schema_version` e `last_validated_at`;
- migration EF Core e SQL posterior de RLS/grants;
- policies sem pressupor claims nativas de provider externo;
- validação que impede snapshot/raw/cache/coleção.

## Escopo excluído

- tabelas globais de players/decks/cards;
- payload raw, nome externo, Arena, coleção, battle history ou cache;
- unicidade global de Player Tag;
- TTL, retenção de API data ou freshness garantida;
- Data API client-side permissiva;
- migration automática no startup de múltiplas réplicas.

## Dependências

- `002-02`, `002-03`, `002-04` e `002-05`;
- contrato `PrimaryPlayerLink` da spec;
- conexão PostgreSQL local/Staging aprovada;
- [ADR 004](../../decisions/004-aspnet-core-react-vite-firebase-postgresql.md).

## Arquivos e símbolos prováveis

- `src/Domain/Identity/CrownPilotUserId.cs`;
- `src/Domain/Identity/PrimaryPlayerLink.cs`;
- `src/Application/Ports/IUserRepository.cs`;
- `src/Application/Ports/IPlayerLinkRepository.cs`;
- `src/Infrastructure/Persistence/CrownPilotDbContext.cs`;
- `src/Infrastructure/Persistence/Configurations/`;
- `src/Infrastructure/Persistence/Migrations/`;
- `database/security/` e `tests/Persistence`, `tests/Integration`, `tests/Rls`.

## Passos de implementação

1. Modelar ID interno, usuário, vínculo e invariantes no Domain.
2. Configurar mappings EF Core, FK, unique `firebase_uid`, `NOT NULL` checks,
   timestamps e token `version` de concorrência.
3. Gerar migration EF Core; não criar as mesmas tabelas em SQL de segurança.
4. Implementar repositories usando ID resolvido do `AuthContext`; ignorar UID do
   body/query.
5. Implementar read/insert/replace após lookup `resolved`, preservando vínculo
   antigo quando provider falha.
6. Proteger replace com `version`/ETag precondition, transação e mapear conflito
   para `409` sem last-write-wins silencioso.
7. Implementar unlink/delete idempotentes somente para dados próprios.
8. Aplicar e testar RLS/grants após migration EF, mantendo authorization obrigatória.

## Testes e comandos de validação

```text
dotnet test --filter Category=Unit
dotnet test --filter Category=Application
dotnet test --filter Category=Persistence
dotnet test --filter Category=Integration
dotnet test --filter Category=Rls
```

Cobrir usuário A/B, anônimo, contexto RLS ausente/adulterado, UID adulterado,
FK/unique/NOT NULL/check constraints, ETag/version ausente ou incorreto,
replace concorrente, falha antes da escrita, unlink repetido e delete repetido.

## Definição de pronto

- schema mínimo possui somente usuário CrownPilot e vínculo primário;
- EF Core é fonte das tabelas/constraints e migration aplica em PostgreSQL;
- SQL de RLS/grants não duplica schema;
- UUID interno é usado em authorization e relacionamentos;
- A não lê/altera B;
- falha de provider não remove vínculo anterior;
- replace concorrente não produz estado silenciosamente incorreto;
- unlink/delete são idempotentes;
- nenhum snapshot/raw/API data é persistido;
- testes de aplicação, persistência, integração e RLS passam.

## Riscos e cuidados

- Não usar Firebase UID como PK interna do domínio.
- Não alegar RLS como substituto da authorization da API.
- Não usar `auth.uid()` sem bridge explicitamente testada.
- Não criar índices/TTL “para o futuro”.
- Não impor unicidade global: uma tag pública pode ser vinculada por mais de um usuário.
