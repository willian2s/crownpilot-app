# 002 — Fundação da aplicação e identidade persistente

- **Status geral:** pending
- **Spec:** [002-fundacao-aplicacao-identidade-persistente.md](../../specs/002-fundacao-aplicacao-identidade-persistente.md)
- **Progresso:** 0/9 subtarefas concluídas

## Objetivo

Criar fundação reproduzível, boundaries seguros e identidade CrownPilot
persistente com Google, mantendo vínculo de uma Player Tag como perfil público
read-only e `ownershipStatus: unverified`.

## Checklist

- [ ] [002-01-bootstrap-toolchain.md](002-01-bootstrap-toolchain.md)
- [ ] [002-02-fixar-boundaries-e-ambientes.md](002-02-fixar-boundaries-e-ambientes.md)
- [ ] [002-03-configurar-supabase-local-e-projetos.md](002-03-configurar-supabase-local-e-projetos.md)
- [ ] [002-04-implementar-auth-google-e-sessao.md](002-04-implementar-auth-google-e-sessao.md)
- [ ] [002-05-criar-adapter-de-lookup.md](002-05-criar-adapter-de-lookup.md)
- [ ] [002-06-persistir-vinculo-com-autorizacao.md](002-06-persistir-vinculo-com-autorizacao.md)
- [ ] [002-07-entregar-fluxos-de-vinculo-e-exclusao.md](002-07-entregar-fluxos-de-vinculo-e-exclusao.md)
- [ ] [002-08-automatizar-quality-gates-e-observabilidade.md](002-08-automatizar-quality-gates-e-observabilidade.md)
- [ ] [002-09-validar-staging-deploy-e-handoff.md](002-09-validar-staging-deploy-e-handoff.md)

## Observações

- Fase 001 liberou somente bootstrap, identidade e vínculo privado read-only.
- Não persistir snapshot, coleção, Arena, battle history, cache ou payload raw.
- Não alegar ownership; usar `public_profile` + `unverified`.
- Stack definida: Laravel + Inertia + React + TypeScript + Vite, com Composer e
  npm mantendo lockfiles próprios.
- Supabase PostgreSQL/Auth definidos; região `sa-east-1` depende de disponibilidade
  no plano/organização.
- Vercel continua alvo inicial, mas Laravel/PHP exige spike de compatibilidade;
- deploy será via `Dockerfile.vercel`/FrankenPHP conforme guia oficial;
  host PHP first-class é fallback se o gate falhar.
- GitHub Actions executará CI em Pull Requests e em todo push para `main`.
- Fase 003 permanece bloqueada até reabertura dos gates de API data, retenção,
  operação, egress, privacidade e meta.
