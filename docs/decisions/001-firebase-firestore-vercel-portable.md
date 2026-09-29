# 001 — Firebase Authentication, Cloud Firestore e Vercel portátil

- **Status:** accepted
- **Data:** 2026-09-29
- **Escopo:** infraestrutura-base do CrownPilot

## Contexto

O CrownPilot precisa de identidade persistente entre dispositivos, banco para
estado da aplicação e uma plataforma simples de deploy.

O responsável definiu:

- Firebase Authentication com Google;
- Cloud Firestore como banco principal;
- Vercel como plataforma inicial de deploy;
- capacidade de migrar o runtime para outro provider sem reescrever o domínio.

## Decisão

### Autenticação

Usar **Firebase Authentication** com Google como provider inicial da conta
CrownPilot.

A autenticação CrownPilot não deve ser confundida com autenticação Supercell.

O vínculo será conceitualmente:

```text
Firebase User
    |
    +-- CrownPilot Profile
            |
            +-- Clash Royale Player Tag
```

Não solicitar ou armazenar credenciais da conta Supercell.

### Banco

Usar **Cloud Firestore** como banco principal.

Firestore é uma dependência aceita do produto. Não existe requisito de tornar o
banco agnóstico a provider nesta decisão.

Ainda ficam para specs futuras:

- collections;
- document shapes;
- índices;
- TTL;
- granularidade de snapshots;
- retenção de battle history;
- estratégia de agregações/meta;
- custos e limites.

Essas decisões devem partir dos contratos de dados e volume medidos na Fase 001.

### Deploy/runtime

Usar **Vercel** como plataforma inicial.

Vercel é uma decisão operacional, não um boundary do domínio.

Código de domínio não deve depender diretamente de serviços proprietários da
Vercel para funcionar.

Evitar como dependência obrigatória:

- Vercel KV/Postgres/Blob;
- Edge Config;
- Vercel Queues/Workflow;
- Vercel Cron;
- APIs proprietárias de runtime/edge.

Esses serviços podem ser adotados no futuro se ficarem isolados e substituíveis.

### Integrações externas

A API da Supercell deve ficar atrás de um adapter/port dedicado, por exemplo:

```text
Domain/Application
       |
       v
ClashRoyaleClient
       |
       +-- HTTP implementation
               |
               +-- api.clashroyale.com
```

Assim, hosting, networking e egress podem mudar sem alterar regras de domínio.

### Jobs e scheduling

Jobs futuros devem ser modelados como operações invocáveis independentemente do
scheduler.

Exemplo:

```text
syncPlayer(playerTag)
ingestMetaBatch(...)
refreshCatalog()
```

Vercel Cron pode disparar uma operação no futuro, mas não deve ser a própria
implementação da regra.

### Egress/IP allowlist

A Fase 001 deve confirmar os requisitos atuais de token da Clash Royale API.

Se for necessário IP estável:

1. Vercel Static IPs pode ser usado como primeira opção quando custo/plano fizer
   sentido;
2. caso contrário, o adapter da Supercell pode ser executado atrás de um egress
   gateway/worker com IP estável;
3. nenhum desses caminhos deve vazar para o domínio.

## Consequências

### Positivas

- login Google simples e persistente;
- Firestore combina com snapshots/documentos do domínio;
- deployment inicial rápido;
- liberdade para mover compute sem trocar auth/banco;
- provider-specific networking fica isolado.

### Trade-offs

- Firestore é vendor lock-in aceito;
- queries/agregações de meta precisam respeitar o modelo do Firestore;
- ingestão analítica pesada pode futuramente exigir componente especializado;
- IP allowlist pode adicionar custo ou componente de egress;
- portabilidade exige disciplina nos boundaries desde o início.

## Não decidido

Esta decisão não escolhe:

- Next.js ou outro framework;
- ORM;
- estrutura de collections;
- região do Firestore;
- estratégia de índices;
- scheduler;
- fila;
- runtime Node/Edge;
- arquitetura do pipeline de meta.

## Regra de revisão

Reabrir esta decisão somente se houver incompatibilidade concreta de:

- custo;
- limites;
- compliance;
- escala;
- requisitos de networking;
- capacidade funcional.

Preferência abstrata por outro provider não é motivo suficiente.
