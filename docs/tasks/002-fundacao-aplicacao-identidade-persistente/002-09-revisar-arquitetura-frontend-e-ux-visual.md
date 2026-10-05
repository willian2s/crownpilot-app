# 002-09 — Revisar arquitetura frontend e UX visual

- **Ticker:** `002`
- **Número:** `09`
- **Status:** `pending`

## Objetivo e resultado esperado

Revisar o frontend React existente após a entrega dos fluxos de identidade e
vínculo. A task produz achados verificáveis, correções mínimas necessárias e um
registro de decisões; não é um placeholder genérico para “melhorar frontend”.

## Requisitos cobertos

- estrutura de diretórios, páginas, componentes, hooks e services;
- API client, models/types, estado local/global, autenticação e environment config;
- loading, cache quando necessário, erros, duplicação e acoplamento;
- hierarquia visual, navegação, responsividade mobile/desktop e consistência;
- acessibilidade, teclado, foco, formulários e mensagens ao usuário;
- testes de componentes e fluxos críticos;
- evidência visual em viewports definidos, sem introduzir biblioteca por preferência.

## Escopo incluído

- inventariar arquitetura real em `frontend/`, não somente documentação;
- rastrear fluxo Auth → token → API e fronteiras entre UI, hooks, services e types;
- revisar configuração por ambiente e tratamento de erros sem secrets;
- exercitar login, vínculo, replace, unlink e delete em estados success/error/
  loading/empty;
- revisar navegação, labels, foco, contraste, semântica, mensagens e responsive
  behavior em mobile e desktop;
- registrar duplicações, acoplamentos e inconsistências com prioridade;
- corrigir somente problemas necessários para os fluxos da Fase 002;
- registrar qualquer item futuro como risco/débito, sem ampliar escopo.

## Escopo excluído

- nova biblioteca de estado, UI ou data fetching sem justificativa concreta;
- redesign amplo, dashboard de jogo, cache especulativo ou abstrações vazias;
- mudanças de contrato API não justificadas pelo achado;
- implementação de features posteriores.

## Dependências

- `002-08` concluída;
- contrato `/api/v1` da `002-07`;
- frontend React/Vite executável do `002-01`.

## Arquivos e evidências prováveis

- `frontend/src/` e configuração Vite/TypeScript;
- relatório versionado em `docs/learning/` ou documentação da task, quando
  necessário para registrar decisão didática;
- screenshots/artefatos de revisão visual, sem dados reais ou secrets;
- testes de componentes e fluxos críticos.

## Passos de implementação

1. Mapear árvore, dependências, ownership dos componentes e fluxo de dados.
2. Revisar API client, auth lifecycle, environment config, estados e duplicação.
3. Executar fluxos críticos com fixtures em viewport mobile e desktop.
4. Verificar teclado, foco, labels, contraste, mensagens, empty/error/loading e
   responsividade.
5. Corrigir apenas achados de severidade que comprometam segurança, compreensão,
   acessibilidade ou manutenção dos fluxos atuais.
6. Registrar achados resolvidos, débitos aceitos e justificativa para não adicionar
   biblioteca/arquitetura nova.

## Testes e comandos de validação

```text
npm run test:unit
npm run lint
npm run typecheck
npm run build
npm run test:e2e
```

Quando browser automation não estiver disponível, registrar evidência equivalente
com testes de componentes, inspeção em viewports e screenshots sanitizados. Não
usar login real ou API Clash Royale live.

## Definição de pronto

- relatório aponta arquivos/componentes analisados e achados com severidade;
- arquitetura frontend, auth/API boundary, estado, config, erros e duplicação
  foram avaliados;
- mobile e desktop têm evidência de navegação, formulário e estados principais;
- acessibilidade de teclado, foco, labels e feedback foi verificada;
- correções necessárias estão implementadas sem biblioteca nova injustificada;
- testes, lint, typecheck e build passam;
- débitos residuais e riscos estão registrados, sem checklist adicional no overview.

## Riscos e cuidados

- Não chamar task concluída por apenas criar uma tela nova.
- Não confundir revisão visual com teste E2E de Staging.
- Não introduzir arquitetura global, cache ou design system sem necessidade concreta.
