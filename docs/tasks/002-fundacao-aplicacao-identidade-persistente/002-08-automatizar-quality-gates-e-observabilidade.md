# 002-08 — Automatizar quality gates e observabilidade

- **Ticker:** `002`
- **Número:** `08`
- **Status:** `pending`

## Objetivo e resultado esperado

Automatizar gates de Pull Request e instrumentar o mínimo operacional para
diagnosticar Auth, vínculo e provider sem registrar dados sensíveis.

## Requisitos cobertos

- CI com lint, type-check, unit, integration, contract e build;
- GitHub Actions em Pull Requests e em todo push para `main`;
- build/smoke de `Dockerfile.vercel` sem secrets;
- RLS tests no Supabase CLI/Docker;
- E2E/smoke com comandos definidos;
- logs e métricas de identidade/lookup;
- chamadas live fora da CI normal;
- redaction de PII/secrets.

## Escopo incluído

- workflow CI em `.github/workflows/ci.yml`;
- triggers `pull_request` e `push` restrito a `main`;
- instalação reprodutível com `composer install` e `npm ci`;
- jobs para lint, typecheck, unit, integration, contract, Rules e build;
- job de `docker build -f Dockerfile.vercel` e smoke do endpoint `/health`;
- execução de Supabase CLI/Docker, migrations e teardown nos testes que precisarem;
- scripts `test:e2e` e `smoke` definidos, mesmo que staging-only;
- logger/métricas com request ID, ambiente, resultado e latência;
- testes/scan para secrets, URLs com tag e imports server-only no bundle.

## Escopo excluído

- provisionamento de produção pelo CI;
- polling ou probes live automáticos;
- dashboards avançados, alertas pagos ou analytics de produto;
- logging de payloads para depuração.

## Dependências

- `002-01` a `002-07`;
- secrets de CI somente para emuladores/fixtures;
- decisão de provider de CI compatível com branch/PR.

## Arquivos e símbolos prováveis

- `.github/workflows/ci.yml`;
- `tests/Unit`, `tests/Feature`, `tests/Contract`, `tests/Rls`;
- `app/Observability/Logger.php`, `Metrics.php`, `Redaction.php`;
- scripts em `composer.json`/`package.json`, `supabase/config.toml` e documentação
  de gates.

## Passos de implementação

1. Definir ordem e falha dos jobs de PR e push em `main`.
2. Rodar Supabase CLI/Docker, migrations e fixtures sem credenciais reais.
3. Adicionar métricas de Auth, link, erro, latência e ambiente.
4. Aplicar redaction antes de serializar qualquer log.
5. Testar que `CLASH_ROYALE_API_TOKEN`, secrets Supabase/PHP e tag não aparecem em
   bundle, output de erro ou logs.
6. Construir imagem Vercel com Docker/FrankenPHP sem injetar secrets na imagem.
7. Documentar gates staging/production sem misturá-los aos gates PR/main.

## Testes e comandos de validação

```text
composer run lint
composer run analyse
composer run test:unit
composer run test:integration
npm run lint
npm run typecheck
npm run test:unit
npm run test:integration
npm run test:contract
npm run test:rls
npx supabase test db
docker build -f Dockerfile.vercel -t crownpilot-vercel-ci .
npm run build
```

Validar workflow em Pull Request e em push para `main`; verificar que falhas em
qualquer gate bloqueiam merge/release. Confirmar ausência de chamadas live à API
externa.

## Definição de pronto

- CI executa todos os gates mínimos em PR;
- GitHub Actions executa os mesmos gates em cada push para `main`;
- RLS tests usam Supabase CLI/Docker e migrations versionadas;
- Docker image de Vercel constrói e `/health` passa sem secrets;
- fixtures cobrem provider sem rede;
- logs não contêm tag, e-mail, token, IP, URL real ou payload;
- métricas distinguem resultado/latência/ambiente;
- comandos E2E/smoke têm owner e pré-condição de staging documentados.

## Riscos e cuidados

- Não marcar CI verde por ignorar teste de Rules ou build.
- Não usar secret de produção em PR/Preview.
- Não considerar Vercel deploy validado sem build/smoke do container.
- Não logar erro HTTP completo se ele puder conter URL/payload sensível.
- Não transformar observabilidade em retenção de API data.
