# ADR 003 — Registro histórico de Supabase Auth e persistência

- **Status:** superseded
- **Data:** 2026-10-02
- **Substituída por:** [ADR 004](004-aspnet-core-react-vite-firebase-postgresql.md)

## Nota de precedência

Este documento registra uma decisão intermediária que não deve orientar a
implementação. A ADR 004 separa Firebase Authentication de
PostgreSQL/Supabase: Supabase não é backend da aplicação nem provedor da
identidade CrownPilot.

## Registro histórico

O planejamento intermediário usava Supabase Auth, JWT/JWKS, PostgreSQL e RLS.
Permanecem válidos os princípios de isolamento por ambiente, autorização
server-side, browser sem acesso direto aos dados CrownPilot e RLS como defesa
adicional.

## Partes substituídas

- Supabase Auth dá lugar a Firebase Authentication/Google Sign-In;
- tokens Firebase são bearer ID tokens validados pelo ASP.NET Core;
- schema passa a ser gerido por migrations EF Core;
- SQL separado cobre somente RLS, grants e objetos de plataforma;
- backend e frontend não dependem de runtime ou sessão PHP.

PostgreSQL no Supabase permanece como infraestrutura de banco. `sa-east-1`
continua condicionada à disponibilidade do plano e à revisão de residência,
backups e subprocessadores. Não há dados legados ou dual-write.
