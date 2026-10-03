# 002-07 — Entregar fluxos de vínculo e exclusão

- **Ticker:** `002`
- **Número:** `07`
- **Status:** `pending`

## Objetivo e resultado esperado

Entregar experiência e endpoints do ciclo de vida do vínculo: consultar, criar,
trocar, desvincular e excluir dados CrownPilot, com comunicação honesta de
perfil público não verificado.

## Requisitos cobertos

- uma Player Tag primária;
- recuperação em outro dispositivo;
- troca/desvinculação explícitas;
- erros de input/provider compreensíveis;
- exclusão inicial de dados próprios;
- disclaimer de conteúdo não oficial.

## Escopo incluído

- telas/estados de login, vínculo, vínculo existente, troca e unlink;
- Controllers/actions Laravel finos para read/link/replace/unlink/delete;
- confirmação antes de troca, desvinculação e exclusão;
- mensagem equivalente a “Perfil público salvo — ownership não verificado”;
- estados loading, vazio, `404`, `429`, `503`, sessão expirada e conflito;
- logout/login e reload recuperando vínculo;
- caminho funcional para apagar usuário e documentos próprios.

## Escopo excluído

- claim de “minha conta”, ownership ou exclusividade;
- múltiplas tags, social, clãs ou notificações;
- dashboard de coleção, Arena ou sincronização completa;
- cobrança, assets não aprovados ou coaching.

## Dependências

- `002-04`, `002-05` e `002-06`;
- contrato HTTP e política de observabilidade da spec;
- domínio fixo de staging para validação posterior.

## Arquivos e símbolos prováveis

- `resources/js/Pages/Login.tsx`, `resources/js/Pages/PlayerLink.tsx`;
- `app/Http/Controllers/PlayerLinkController.php`;
- `routes/web.php` e `routes/api.php` quando boundary API for necessário;
- `LinkPlayerTagForm`, `PlayerLinkStatus`, `mapLinkError`;
- `deleteAccount`/`deleteOwnData` use case;
- disclaimer em layout/footer apropriado.

## Passos de implementação

1. Implementar leitura do estado vinculado do usuário atual.
2. Implementar submit de tag com validação e feedback sem revelar provider.
3. Fazer replace somente após `resolved` e confirmação explícita.
4. Implementar unlink e exclusão com confirmação e idempotência.
5. Tratar sessão expirada e reload/outro dispositivo.
6. Exibir `public_profile`/`unverified` sem linguagem de ownership.
7. Adicionar disclaimer legível de fan content não oficial.

## Testes e comandos de validação

```text
npm run test:unit
npm run test:integration
npm run build
```

Cobrir ausência inicial, sucesso, tag inválida, perfil inexistente, provider
indisponível, rate limit, troca que falha, unlink repetido e delete repetido.

## Definição de pronto

- usuário autentica e informa tag uma vez;
- reload e outro dispositivo recuperam vínculo após login;
- troca valida nova tag antes de substituir antiga;
- falha externa preserva vínculo anterior;
- unlink e exclusão funcionam sem reentrada de dados indevida;
- UI nunca diz que usuário possui ou controla perfil;
- estados vazios/erro/loading são acessíveis e não quebram build.

## Riscos e cuidados

- Não expor nome, coleção ou resposta raw do provider para “enriquecer” tela.
- Não confundir exclusão local com exclusão no Clash Royale/provider.
- Não enviar Player Tag real para analytics/logs client-side.
- Não criar promessa de sync a partir de validação única.
