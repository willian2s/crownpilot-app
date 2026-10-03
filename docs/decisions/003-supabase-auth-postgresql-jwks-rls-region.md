# ADR 003 — Supabase Auth, PostgreSQL, JWT/JWKS, RLS e região

- **Status:** Accepted
- **Data:** 2026-10-02
- **Escopo:** identidade, persistência, ambientes e região da Fase 002
- **Fase:** 002 — Fundação da aplicação e identidade persistente
- **Supersede parcialmente:** [ADR 001](001-firebase-firestore-vercel-portable.md)
- **Atualiza:** [ADR 002](002-laravel-inertia-react-vite-firestore-region.md)

## Contexto

Laravel/PHP não possui Firebase Admin SDK oficial. A integração exigiria
dependências PHP de terceiros para Auth e Firestore, além de um client Google
Cloud separado. Isso aumenta risco de compatibilidade, especialmente em
FrankenPHP/Vercel e `ext-grpc`.

O Supabase oferece PostgreSQL, conexão PHP padrão, migrations, transações e uma
[integração Laravel documentada](https://supabase.com/docs/guides/getting-started/quickstarts/laravel).
A decisão mantém Laravel + Inertia + React +
TypeScript + Vite, Vercel via Docker/FrankenPHP e GitHub Actions.

## Decisão

### Identidade

- Supabase Auth como provider de identidade CrownPilot;
- Google como provider OAuth;
- React usa `@supabase/supabase-js` somente para login, refresh e logout;
- Laravel recebe Supabase access JWT, valida assinatura via JWKS e cria sessão
  Laravel segura;
- Laravel nunca recebe Google ID token diretamente;
- `sub` do JWT é UUID do usuário Supabase e origina `crownpilotUserId`;
- `iss`, `aud`, `exp`, `nbf`, `alg`, `kid` e project ref são validados;
- JWKS é separado por ambiente e cacheado com refresh quando `kid` desconhecido;
- signing keys assimétricas são obrigatórias; JWT simétrico não pode ser aceito
  silenciosamente.

Fluxo:

```text
React -> Supabase Auth Google -> Supabase access JWT
     -> Laravel valida JWT/JWKS -> Laravel session
     -> autorização -> Eloquent/PostgreSQL
```

### Persistência

Supabase PostgreSQL é banco principal. `supabase/migrations/` é fonte única do
schema e RLS. Eloquent fornece models, casts, queries e repositories, mas não
mantém migrations equivalentes em `database/migrations/`.

Schema mínimo:

```sql
create table public.crownpilot_users (
  id uuid primary key references auth.users(id),
  schema_version text not null,
  created_at timestamptz not null,
  updated_at timestamptz not null
);

create table public.primary_player_links (
  user_id uuid primary key references public.crownpilot_users(id),
  schema_version text not null,
  player_tag text not null,
  subject_type text not null check (subject_type = 'public_profile'),
  ownership_status text not null check (ownership_status = 'unverified'),
  linked_at timestamptz not null,
  updated_at timestamptz not null,
  last_validated_at timestamptz null
);
```

`player_tag` não possui unicidade global. `auth.users` é gerenciado pelo Supabase;
aplicação não escreve diretamente nele.

### Authorization e RLS

Browser acessa Supabase somente para Auth. Browser não acessa Data API/PostgREST
para dados CrownPilot nesta fase.

RLS fica habilitado e deny-by-default nas tabelas expostas. Grants `anon` e
`authenticated` devem ser revogados ou limitados. Laravel continua responsável
por autorização A/B porque conexão PostgreSQL privilegiada pode ignorar RLS.
RLS é defesa adicional, não substituto de autorização server-side.

Nenhuma `service_role`, database password ou secret Supabase chega ao browser.

### Região e projetos

Staging e production usarão região Supabase **South America (São Paulo),
`sa-east-1`**, conforme [regiões oficiais](https://supabase.com/docs/guides/platform/regions),
somente após confirmar disponibilidade no plano/organização.
Local usa Supabase CLI/Docker e não depende de região cloud.

Se `sa-east-1` não estiver disponível, provisionamento fica bloqueado; não usar
região genérica `Americas` como fallback silencioso. DPA, backups,
subprocessadores, Auth/logs e requisitos de residência devem ser revisados antes
de dados reais.

### Ambientes

- local: Supabase CLI/Docker, dados descartáveis;
- Preview: sem Supabase real/Auth real, mocks e build;
- staging: projeto Supabase separado, domínio fixo e Google OAuth real;
- production: projeto Supabase separado, secrets próprios e deploy pela `main`.

### Deploy e CI

Deploy continua Vercel via `Dockerfile.vercel`, FrankenPHP/PHP 8.5, `Caddyfile`,
`vercel.json` e `.dockerignore`. Estado durável permanece PostgreSQL/Supabase;
filesystem de container é efêmero.

GitHub Actions é CI autoritativo em `pull_request` e `push` para `main`. Gates:

- Composer/npm, lint, análise estática e type-check;
- testes PHP/JS, contrato e E2E conforme ambiente;
- `supabase start`, migrations e testes RLS;
- build Docker e smoke `/health`;
- nenhuma chamada live à Clash Royale API.

## Alternativas consideradas

### Firebase Auth + Firestore

Não adotado para esta implementação. Firebase continua referência histórica da
Fase 001, mas a ausência de Admin SDK PHP oficial e o custo de adapters Auth +
Firestore tornam Supabase mais coerente com Laravel.

### Laravel Auth + Supabase somente como PostgreSQL

Não adotado inicialmente. Duplicaria o provider OAuth/identidade. Supabase Auth
fica responsável por Google OAuth; Laravel mantém sessão e autorização da
aplicação.

### Supabase Data API no browser

Não adotado. Aumentaria superfície pública e misturaria RLS/Data API com o
boundary Laravel. Pode ser reavaliado somente por necessidade concreta.

### Migrations duplicadas Laravel + Supabase

Não adotado. Duas fontes de schema criariam drift. `supabase/migrations/` é a
fonte única; Eloquent é camada de acesso.

### Região genérica `Americas`

Não adotada. Pode escolher região diferente de São Paulo sem decisão explícita.

## Consequências

### Positivas

- Laravel usa PostgreSQL/Eloquent, transações e migrations naturais;
- elimina dependência de Firebase Admin PHP não oficial;
- RLS adiciona defesa de banco;
- Auth Google continua gerenciado;
- domínio permanece isolado via ports/adapters;
- região e projetos ficam explícitos.

### Trade-offs

- Supabase passa a ser dependência de Auth, database e operação;
- JWT/JWKS, rotação de keys e sessão exigem implementação cuidadosa;
- RLS não protege automaticamente conexão Eloquent privilegiada;
- Supabase CLI/Docker aumenta peso local/CI;
- `sa-east-1` precisa estar disponível no plano;
- migração futura de Firebase, se houver dados externos, exigirá projeto próprio.

## Gates obrigatórios

- `sa-east-1` disponível antes de provisionar staging/production;
- projetos Supabase separados;
- OAuth Google com redirect allowlist por ambiente;
- JWKS assimétrico validado com issuer/audience/expiry/kid;
- `pdo_pgsql`, SSL e pooler validados no container PHP 8.5;
- migrations únicas aplicam em `supabase db reset`;
- RLS/grants testados para anônimo e autenticado;
- Laravel impede usuário A de acessar dados de B;
- `service_role`, database password e JWT secrets não aparecem no browser;
- exclusão diferencia dados CrownPilot de `auth.users`;
- CI passa em PR e push `main` antes de release.

## Rollback e revisão

Não há dados Firebase remotos ou runtime implementado neste baseline. Não criar
dual-write. Se Supabase falhar antes de provisionamento, reabrir ADR 001/002 sem
migração. Depois de provisionar, rollback de aplicação deve preservar schema;
mudanças SQL incompatíveis exigem migration reversível ou cutover controlado.

Revisar esta ADR se disponibilidade regional, DPA, segurança JWT, RLS/Eloquent,
limites de conexão, custo ou requisito de Auth invalidar a decisão.
