# ADR 002 — Registro histórico de stack e deploy

- **Status:** superseded
- **Data:** 2026-10-02
- **Substituída por:** [ADR 004](004-aspnet-core-react-vite-firebase-postgresql.md)

## Nota de precedência

Este ADR é histórico e não normativo. A decisão de usar ASP.NET Core,
React/Vite, Firebase Authentication e PostgreSQL/Supabase está consolidada na
ADR 004.

## Registro histórico

O planejamento anterior propôs backend PHP com integração server-driven,
frontend React/Vite, autenticação Firebase e Firestore, com deploy containerizado
em Vercel. Essa proposta não foi implementada e foi reaberta antes do bootstrap.

## Princípios preservados

- React + TypeScript + Vite continuam sendo o frontend oficial;
- Vercel pode hospedar frontend estático sem ser dependência do backend;
- estado durável não depende de filesystem efêmero;
- CI valida build, testes, container e smoke antes de release;
- Staging usa hostname fixo e Production usa configuração própria;
- domínio não conhece provider de deploy ou integração externa.

Não usar este ADR para implementar framework, sessão, persistência ou runtime.
