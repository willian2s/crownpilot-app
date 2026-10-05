# ADR 001 — Registro histórico de infraestrutura

- **Status:** superseded
- **Data:** 2026-09-29
- **Substituída por:** [ADR 004](004-aspnet-core-react-vite-firebase-postgresql.md)

## Nota de precedência

Este documento registra decisão exploratória anterior. Não orienta a
implementação da Fase 002. A arquitetura vigente está na ADR 004.

## Registro histórico

A decisão original considerou Firebase Authentication, Cloud Firestore e Vercel
como infraestrutura inicial. Também registrou portabilidade de domínio,
ausência de credenciais Supercell, adapter para a integração Clash Royale e
ausência de dependência obrigatória de serviços proprietários de hosting.

As partes de autenticação e banco foram substituídas pela ADR 004. Portabilidade,
isolamento da integração externa e Vercel como opção operacional permanecem
princípios válidos, reinterpretados pela arquitetura oficial.

Não há runtime ou dados legados para migrar. Não usar este ADR para escolher
framework, banco, ORM, sessão ou deployment da Fase 002.
