# 002-02 — Fixar boundaries e ambientes

- **Ticker:** `002`
- **Número:** `02`
- **Status:** `pending`

## Objetivo e resultado esperado

Registrar decisões arquiteturais duradouras e matriz operacional de local,
preview, staging e production antes de provisionar dados reais. O resultado é um
boundary explícito entre browser, Laravel/PHP, Supabase, provider externo e
Vercel.

## Requisitos cobertos

- portabilidade fora de Vercel;
- PHP/Laravel versus runtime serverless Vercel;
- adapter PHP/JWT/JWKS/Eloquent e Supabase Auth;
- PostgreSQL regional em `sa-east-1`, condicionado ao plano;
- projetos e secrets separados;
- preview sem Supabase Auth real;
- staging fixo e production pela `main`.

## Escopo incluído

- confirmar Laravel + Inertia + React + Vite, Web SDK somente Auth e adapter PHP
  server-side;
- validar compatibilidade de Laravel/PHP com Vercel antes de tratá-la como host
  de production;
- definir RLS/grants deny-by-default e autorização backend;
- registrar `sa-east-1` para staging e production se disponível;
- criar ADRs para boundary de acesso, região/topologia e persistência mínima;
- definir matriz de variáveis por ambiente e startup fail-closed;
- definir quais features são mockadas/desabilitadas em preview.

## Escopo excluído

- criar projetos Supabase ou configurar OAuth;
- implementar handlers, repository ou UI;
- escolher egress definitivo para produção;
- liberar sync, retenção, billing ou polling.

## Dependências

- `002-01`;
- ADR 001;
- ADR 003;
- veredito da Fase 001;
- disponibilidade de `sa-east-1` no plano/organização;
- spike de runtime PHP/Vercel ainda pendente.

## Arquivos e símbolos prováveis

- `docs/decisions/002-laravel-inertia-react-vite-firestore-region.md`;
- `docs/decisions/003-supabase-auth-postgresql-jwks-rls-region.md`;
- `.env.example`, documentação de ambientes e `AGENTS.md`;
- `config/environment.php` ou equivalente;
- `EnvironmentName`, `RuntimeConfig`, `SupabaseProjectConfig`.

## Passos de implementação

1. Desenhar fluxo Browser → Supabase Auth → JWT/JWKS → Laravel session →
   Eloquent/PostgreSQL.
2. Registrar que UUID só vem de JWT verificado e que RLS não substitui
   autorização Laravel.
3. Confirmar `sa-east-1`, paridade staging/production e custo de migração futura.
4. Definir allowlist de project IDs e hosts por ambiente.
5. Definir nomes de variáveis públicas e server-only, sem valores reais.
6. Registrar rollback e mudança futura de região como cutover de novo projeto.

## Testes e comandos de validação

```text
composer run lint
composer run analyse
npm run lint
npm run typecheck
npm run test:unit
```

Revisão documental deve confirmar que não há secret real, que preview não aponta
para Supabase e que produção não pode ser provisionada sem região disponível e
registrada.

## Definição de pronto

- ADRs registram boundary, região/topologia e persistência mínima;
- matriz local/preview/staging/production tem Supabase Auth/PostgreSQL, dados, hosts e
  secrets explicitamente separados;
- configuração inválida falha fechado;
- Vercel não aparece como dependência de domínio;
- `sa-east-1` está disponível e registrada antes de qualquer provisionamento;
- gate Vercel/PHP e fallback PHP first-class estão documentados.

## Riscos e cuidados

- Não trocar `sa-east-1` sem nova decisão explícita e plano de migração.
- Não chamar acesso autenticado de autorização.
- Não registrar IDs, emails ou secrets reais em ADR.
- Não criar ADR para proxy como contrato permanente; transporte permanece adapter.
