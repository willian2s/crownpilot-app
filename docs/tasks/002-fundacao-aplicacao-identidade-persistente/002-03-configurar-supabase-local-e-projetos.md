# 002-03 — Configurar Supabase local e projetos

- **Ticker:** `002`
- **Número:** `03`
- **Status:** `pending`

## Objetivo e resultado esperado

Configurar Supabase CLI/Docker e separação de projetos para que local, staging e
production sejam reproduzíveis sem compartilhar dados ou credenciais.

## Requisitos cobertos

- Supabase Auth/PostgreSQL local com CLI/Docker;
- projetos Supabase separados;
- Google Sign-In via Supabase Auth em staging;
- migrations SQL e RLS versionadas;
- PostgreSQL sem acesso público irrestrito;
- ausência de produção no ambiente local/preview.

## Escopo incluído

- `supabase/config.toml` e configuração local sem secrets;
- migrations SQL e policies RLS versionadas;
- projetos/refs separados para staging e production, sem hardcode de
  credenciais;
- configuração inicial de provider Google no Supabase Auth staging;
- documentação de criação/configuração manual que não puder ser automatizada;
- criação de Supabase staging/production em `sa-east-1`, após confirmar plano.

## Escopo excluído

- repository Eloquent de vínculo e handlers de aplicação;
- dados reais de jogadores;
- acesso client-side ao Data API/PostgREST;
- snapshots, raw API, cache, índices ou tabelas futuras.

## Dependências

- `002-01` e `002-02`;
- acesso administrativo aos projetos Supabase, fornecido fora do Git;
- domínio fixo de staging disponível para Auth.

## Arquivos e símbolos prováveis

- `supabase/config.toml`, `supabase/migrations/`, `supabase/tests/`;
- `tests/Rls/` e configuração Supabase CLI/Docker;
- documentação de setup de projetos;
- `SupabaseEnvironment`, `LocalStackConfig`, `assertProjectIsolation`.

## Passos de implementação

1. Instalar/pinar Supabase CLI e confirmar Docker conforme política do repositório.
2. Executar `supabase init` e `supabase start` com dados descartáveis.
3. Criar migrations iniciais, RLS deny-by-default e harness de testes.
4. Associar refs locais, staging e production sem commitar secrets.
5. Configurar Google provider e redirect allowlist somente em staging/production.
6. Adicionar guardas de project ref/issuer e ambiente para impedir cross-environment.
7. Provisionar staging/production somente em `sa-east-1`, após validar plano,
   runtime PHP, PostgreSQL, SSL e pooler.

## Testes e comandos de validação

```text
supabase start
supabase db reset
composer run test:integration
npm run test:rls
npx supabase test db
```

Testar que local usa stack Supabase/Docker, preview não usa Supabase real e cada
projeto staging/production é distinto.

## Definição de pronto

- Supabase CLI/Docker inicia por comando documentado;
- migrations e RLS estão versionadas e deny-by-default;
- testes RLS passam para anônimo e autenticado;
- refs/projetos de staging e production são distintos;
- Google Sign-In de staging usa hostname fixo e redirect allowlist;
- nenhum secret, token ou projeto production é usado por local/preview;
- Supabase staging/production usa `sa-east-1`, se disponível no plano;
- adapter PHP conecta ao PostgreSQL local e ao projeto correto;
- production só existe após gates de região, runtime PHP/Vercel e secrets.

## Riscos e cuidados

- Supabase URL/anon key podem ser públicas, mas devem ser específicas por ambiente.
- service role, database password e JWT secret nunca entram em `.env.example`,
  bundle ou migrations.
- Não usar RLS como substituto da autorização Laravel.
- Não provisionar tabelas de snapshot por conveniência.
