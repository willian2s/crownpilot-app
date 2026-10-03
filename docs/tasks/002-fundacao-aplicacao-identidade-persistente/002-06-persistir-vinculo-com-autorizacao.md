# 002-06 — Persistir vínculo com autorização

- **Ticker:** `002`
- **Número:** `06`
- **Status:** `pending`

## Objetivo e resultado esperado

Persistir somente identidade CrownPilot e vínculo primário de perfil público,
com autorização derivada do UID verificado e sem criar snapshot de API.

## Requisitos cobertos

- PostgreSQL server-side via Eloquent e adapter Supabase;
- schema mínimo versionado;
- usuário A isolado de B;
- `public_profile` + `unverified`;
- troca sem perda e concorrência protegida;
- RLS/grants sem acesso público irrestrito.

## Escopo incluído

- documento `users/{uid}` mínimo;
- tabela `primary_player_links` com `user_id` como primary key;
- repository com `user_id` derivado exclusivamente de `AuthContext`;
- operações read, create/replace, unlink e delete idempotente;
- transaction/locking/precondition para replace concorrente;
- timestamps e `schemaVersion`;
- validação para impedir campos de snapshot/raw/cache.

## Escopo excluído

- collections globais de players/decks/cards;
- payload raw, nome externo, Arena, coleção, battle history ou cache;
- unicidade global de Player Tag;
- TTL, índices compostos e retenção de API data;
- Data API client-side permissiva.

## Dependências

- `002-02`, `002-03`, `002-04` e `002-05`;
- `sa-east-1` e projetos Supabase;
- contrato `PrimaryPlayerLink` da spec.

## Arquivos e símbolos prováveis

- `app/Domain/Identity/PrimaryPlayerLink.php`;
- `app/Application/PlayerLink/PlayerLinkRepository.php`;
- `app/Adapters/Supabase/PostgresPlayerLinkRepository.php`;
- `PrimaryPlayerLinkRepository`, `AuthContext`, `OwnershipStatus`;
- `firestore.rules`, `tests/integration/player-link/`.

## Passos de implementação

1. Definir schema mínimo e validação de versão.
2. Derivar document path do UID verificado; ignorar UID do body/query.
3. Implementar leitura de ausência sem fabricar perfil.
4. Implementar replace após lookup resolvido e preservar o antigo em falha.
5. Proteger concorrência com transaction/precondition e mapear `409`.
6. Implementar unlink/delete idempotentes e remoção somente dos dados próprios.
7. Verificar que RLS e conexão Eloquent não contornam autorização de aplicação.

## Testes e comandos de validação

```text
composer run test:unit
composer run test:integration
npm run test:rls
npx supabase test db
```

Cobrir usuário A/B, anônimo, UID adulterado, replace concorrente, falha antes da
escrita, unlink repetido e delete repetido.

## Definição de pronto

- apenas duas tabelas mínimas por usuário são necessárias;
- nenhum snapshot ou raw API data é persistido;
- UUID `sub` é derivado e autorizado server-side;
- A não lê/altera B;
- falha de provider não remove vínculo anterior;
- replace concorrente não produz estado silenciosamente incorreto;
- unlink/delete são idempotentes;
- integration e RLS tests passam.

## Riscos e cuidados

- Conexão Eloquent privilegiada pode ignorar RLS: repository e casos de uso
  Laravel devem aplicar autorização.
- `lastValidatedAt` não deve ser apresentado como freshness do provider.
- Não criar índice/TTL “para o futuro”.
- Não impor unicidade global: mesma tag pode ser perfil público vinculado por mais
  de um usuário CrownPilot.
