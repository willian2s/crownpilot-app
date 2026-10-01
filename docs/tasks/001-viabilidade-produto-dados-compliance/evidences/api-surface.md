# API surface — Clash Royale API

- **Task:** `001-01`
- **Observed at:** `2026-10-01`
- **Overall status:** `completed with constraints`
- **Initial blocker:** `CLASH_ROYALE_API_TOKEN` was not configured in the first execution environment. No 2xx authenticated probe or positive tag resolution was possible in that run.
- **Current direct-route blocker:** token loads from `.env.local`, but current egress is not allowed by API key configuration.
- **Current proxy result:** authenticated cards, profile and battlelog probes returned 2xx through RoyaleAPI Proxy; route is selected provisionally for server-side discovery, with third-party conditions recorded below.
- **Authority:** [Clash Royale API developer portal](https://developer.clashroyale.com/)

## Confidence rules

- **Documented:** information visible in the official portal or its documentation route.
- **Observed:** result returned by a controlled request in this run.
- **Pending:** requires authenticated portal/API session; not treated as contract.
- No secondary SDK or community wrapper was used as authority.

## Account, key and authentication

| Item | Result | Evidence |
| --- | --- | --- |
| Account/key creation | Portal supports account registration, email verification and `Create New Key`. Key configuration includes description and allowed public IP addresses. Key configuration cannot be edited; configuration changes require a new key. | **Documented**, portal UI loaded on `2026-10-01`: [portal](https://developer.clashroyale.com/), [getting started](https://developer.clashroyale.com/#/getting-started) |
| Credential type | Portal states that each request requires a JSON Web Token. | **Documented**, [portal](https://developer.clashroyale.com/) |
| Authorization header | Expected request form is `Authorization: Bearer <CLASH_ROYALE_API_TOKEN>`. Acceptance is **pending on direct official route** because egress is not allowlisted; operational acceptance was observed through RoyaleAPI Proxy. | **Direct route pending; proxy observed separately**, [API documentation](https://developer.clashroyale.com/api-docs/index.html) |
| IP/egress | Key is bound to specified IP addresses. Portal rejects private IP ranges in key configuration and recommends a web server rather than browser calls. Current portal UI exposes up to five CIDR entries for a key. | **Documented**, portal UI/bundle, [portal](https://developer.clashroyale.com/) |
| Rate binding | Portal states token is bound to rate limitations and calls fail when limitations are exceeded. Numeric limit was not exposed without authenticated documentation session. | **Documented qualitatively; numeric value pending**, [portal](https://developer.clashroyale.com/) |
| Response encoding | API responses are JSON UTF-8 documents and use standard HTTP status codes; error responses are JSON too. | **Documented**, [portal](https://developer.clashroyale.com/) |
| Local secret (initial run) | `test -n "$CLASH_ROYALE_API_TOKEN"` returned false. No token was printed or written. | **Observed** |

## Candidate endpoint map

The table below describes the direct official route. Proxy observations are
recorded separately and are not promoted to official API contracts.

The documentation route is a Swagger UI served by the official portal. Its resource specification requires an authenticated portal session. Every route, parameter detail, pagination rule and response shape below is therefore a **pending direct-route candidate mapping**, not a confirmed contract. Proxy observations are recorded separately and are not promoted to official API contracts.

| Endpoint | Purpose | Auth | Parameters | Pagination | Response shape | Source URL | Evidence | MVP status |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| `GET /v1/players/{playerTag}` | Candidate public player profile, progression, collection and current deck | JWT; header confirmation pending on direct route | Candidate path tag; encode `#` as `%23`; exact contract pending | Pending authenticated schema | Pending authenticated schema; top-level type not confirmed | [API docs](https://developer.clashroyale.com/api-docs/index.html) | Direct route blocked by IP; proxy observation separate, not official contract | required |
| Official ownership verification surface | Endpoint not confirmed in Swagger reference available to this run; no path assumed | `UNRESOLVED` | Player token flow and endpoint, if any, require official confirmation | Pending | Pending | [API docs](https://developer.clashroyale.com/api-docs/index.html) | No `verifytoken` path should be treated as verified from current evidence | useful |
| `GET /v1/players/{playerTag}/battlelog` | Candidate recent battle history | JWT; header confirmation pending on direct route | Candidate path tag; exact contract pending | Pending authenticated schema | Pending authenticated schema; top-level type not confirmed | [API docs](https://developer.clashroyale.com/api-docs/index.html) | Direct route blocked by IP; proxy observation separate, not official contract | required |
| `GET /v1/cards` | Candidate card catalog and card capability metadata | JWT; header confirmation pending on direct route | Exact parameters pending | Pending authenticated schema | Pending authenticated schema; top-level type not confirmed | [API docs](https://developer.clashroyale.com/api-docs/index.html) | Direct route blocked by IP; proxy observation separate, not official contract | required |
| `GET /v1/locations` | Candidate location/region catalog | JWT; header confirmation pending | Exact parameters pending | Pending authenticated schema | Pending authenticated schema; top-level type not confirmed | [API docs](https://developer.clashroyale.com/api-docs/index.html) | Candidate route from official docs route; 2xx pending | useful |
| `GET /v1/locations/{locationId}` | Candidate location detail | JWT; header confirmation pending | Candidate path location ID; exact contract pending | Not applicable pending schema | Pending authenticated schema; top-level type not confirmed | [API docs](https://developer.clashroyale.com/api-docs/index.html) | Candidate route from official docs route; 2xx pending | useful |
| `GET /v1/locations/{locationId}/rankings/players` | Candidate player leaderboard by location | JWT; header confirmation pending | Candidate path location ID; ranking parameters pending | Pending authenticated schema | Pending authenticated schema; entry wrapper not confirmed | [API docs](https://developer.clashroyale.com/api-docs/index.html) | Candidate route from official docs route; 2xx pending | useful |
| `GET /v1/locations/{locationId}/rankings/clans` | Candidate clan leaderboard by location | JWT; header confirmation pending | Candidate path location ID; ranking parameters pending | Pending authenticated schema | Pending authenticated schema; entry wrapper not confirmed | [API docs](https://developer.clashroyale.com/api-docs/index.html) | Candidate route from official docs route; 2xx pending | useful |
| `GET /v1/locations/{locationId}/pathoflegend/players` | Candidate ranked/Path of Legends leaderboard | JWT; header confirmation pending | Candidate path location ID; parameters pending | Pending authenticated schema | Pending authenticated schema; entry wrapper not confirmed | [API docs](https://developer.clashroyale.com/api-docs/index.html) | Candidate surface; current availability pending | useful |
| `GET /v1/clans?name={name}` | Candidate small-sample clan discovery | JWT; header confirmation pending | Candidate clan-name query; matching/limit semantics pending | Pending authenticated schema | Pending authenticated schema; top-level type not confirmed | [API docs](https://developer.clashroyale.com/api-docs/index.html) | Candidate route from official docs route; 2xx pending | useful |
| `GET /v1/clans/{clanTag}` | Candidate clan metadata | JWT; header confirmation pending | Candidate path clan tag; encoding pending for `#` | Pending authenticated schema | Pending authenticated schema; top-level type not confirmed | [API docs](https://developer.clashroyale.com/api-docs/index.html) | Candidate route from official docs route; 2xx pending | useful |
| `GET /v1/clans/{clanTag}/members` | Candidate clan member discovery | JWT; header confirmation pending | Candidate path clan tag; exact contract pending | Pending authenticated schema | Pending authenticated schema; entry wrapper not confirmed | [API docs](https://developer.clashroyale.com/api-docs/index.html) | Candidate route from official docs route; 2xx pending | useful |

### Ownership decision pending

A public profile lookup must not be described as verified ownership. No official
ownership-verification endpoint is confirmed from the Swagger reference available
to this run; do not assume a `verifytoken` path or request a player token. MVP
language remains **public profile linked by user-provided Player Tag**, not
“account ownership confirmed”.

### Current product decision

- Save user-provided Player Tag as a private, read-only public-profile
  association.
- Use `subjectType: public_profile` and `ownershipStatus: unverified`.
- Derive CrownPilot user association server-side from Firebase Auth.
- Do not imply account ownership, exclusivity, identity, actions or notifications.
- Do not request player token or invent an ownership endpoint.

## Official documentation validation

- **Observed at:** `2026-10-01`
- `https://developer.clashroyale.com/` returned `200` and served the official
  documentation application shell.
- `https://developer.clashroyale.com/api-docs/index.html` loaded the same application,
  but the resource specification requires an authenticated portal session.
- Public bundle inspection did not expose endpoint literals or an OpenAPI spec;
  portal login flow references a session `swaggerUrl` and temporary API token.
- **Conclusion:** official source and session requirement are confirmed, but
  public access did not close endpoint paths, parameters, pagination or
  ownership verification. Keep these fields `PENDING` until authenticated portal review or
  an official sanitized export is available. No portal credential was collected.

## Safe probes executed (initial direct route)

All probes used temporary response/header files, did not print a token and did not collect player data.

| Probe | Result | Interpretation |
| --- | --- | --- |
| `GET https://api.clashroyale.com/v1/cards` without `Authorization` | `403`, `{"reason":"accessDenied","message":"Missing authorization"}` | API is reachable and rejects missing auth before resource access. |
| `GET https://api.clashroyale.com/v1/players/%23INVALIDTAG` without `Authorization` | `403`, same JSON error | `%23` was sent in path; auth gate prevents confirming tag parsing or existence. |
| `GET https://api.clashroyale.com/v1/players/not-a-player-tag` without `Authorization` | `403`, same JSON error | Invalid-tag status cannot be isolated without valid auth. |
| Response headers for all three probes | `Content-Type: application/json`; `Cache-Control: public max-age=600`; no `Retry-After` | Observed unauthenticated error response only; not a resource-cache contract. |

In the initial direct route, no 2xx result, valid tag resolution,
missing-token-vs-invalid-tag separation, 404, 429 or 5xx behavior was confirmed.
A deliberate rate-limit probe was not run.

## Retake after credential configuration

- **Observed at:** `2026-10-01`
- **Credential loading:** `.env.local` was loaded only by the local probe
  process. `CLASH_ROYALE_API_TOKEN` was present; its value was not printed or
  persisted.
- **Probe 1:** `GET /v1/cards` with `Authorization: Bearer` → `403`, JSON
  `reason=accessDenied.invalidIp`.
- **Probe 2:** `GET /v1/players/%23INVALIDTAG` with `Authorization: Bearer` →
  `403`, JSON `reason=accessDenied.invalidIp`.
- **Selected headers:** `Content-Type: application/json`; `Cache-Control: public
  max-age=600`.
- **Interpretation:** request reached API authorization and was rejected because
  current egress is not allowed by key configuration. This is not evidence of an
  invalid token or invalid `#` encoding. No IP was recorded in this evidence.
- **Remaining:** add current egress to key allowlist or run from an allowed
  environment, then repeat cards, permitted Player Tag and battle log probes.

## Third-party egress proxy decision

- **Lead:** [RoyaleAPI Proxy](https://docs.royaleapi.com/proxy.html), reported as
  an option for servers without fixed IP.
- **Verification:** automated fetch returned `403`/Cloudflare on `2026-10-01`;
  responsible operator supplied RoyaleAPI Terms of Service and Privacy Policy,
  both marked `Last Updated: 2024-10-11`.
- **Source locator:** exact URLs for supplied TOS/Privacy text were not included;
  attempts against `https://royaleapi.com/terms`,
  `https://royaleapi.com/privacy` and
  `https://royaleapi.com/terms-of-service` returned `403`/Cloudflare.
  Applicability to proxy endpoint remains pending.
- **Decision:** use proxy as the current operational transport until explicit
  replacement or fixed-egress server, because direct egress has no fixed
  allowlisted IP. Classification is `CURRENT OPERATIONAL DEPENDENCY —
  CONDITIONAL — NOT AUTHORIZATION`.
- **Privacy evidence:** supplied Privacy Policy declares collection of IP,
  geolocation, access times, URLs, clickstream and usage data, and sharing with
  service providers. Supplied text does not specify proxy API-key handling.
- **Boundary:** do not classify proxy access as official API authorization or
  commercial permission. Keep key out of browser and do not log token, real
  Player Tag or personal payloads.
- **Pending conditions:** key handling, retention, limits, SLA, security behavior
  and compatibility with Supercell agreements remain `UNRESOLVED` and tracked as
  operational debt. Current use remains server-side and does not authorize
  browser exposure, billing or commercial use.
- **Fallback:** evaluate own fixed-egress gateway if proxy cannot provide
  acceptable guarantees. Any proxy integration must remain server-side.

## Authenticated probes through RoyaleAPI Proxy

- **Observed at:** `2026-10-01`
- **Route:** replace `https://api.clashroyale.com` with
  `https://proxy.royaleapi.dev`, preserving `/v1/...` paths and Bearer
  authorization. This is host substitution only for the observed GETs; POST
  behavior, including any ownership-verification endpoint, remains unprobed.
- **Credential:** rotated key supplied by responsible operator, loaded from
  `.env.local`; value was not printed or persisted.
- **Data scope:** probes used only the responsible operator's supplied sample
  tag; no production user dataset was collected or persisted.

| Probe | Status | Sanitized result | Cache-Control |
| --- | --- | --- | --- |
| `GET /v1/cards` | `200` | JSON object; `items=123`, `supportItems=4` | `max-age=49` |
| `GET /v1/locations` | `200` | JSON object; `items=262`, `paging` field present | `max-age=404` |
| `GET /v1/players/%23<REDACTED_VALID_TAG>` | `200` | JSON object; `cards=123`; progression/deck fields present | `max-age=60` |
| `GET /v1/players/%23<REDACTED_VALID_TAG>/battlelog` | `200` | JSON array; 30 entries; battle/arena/deck/game fields present | `max-age=60` |
| `GET /v1/players/%23INVALIDTAG` | `404` | JSON error object with `reason` field | `public max-age=600` |

`200` responses confirm operational authentication, `%23` path encoding and
candidate response shapes through proxy. They do not prove direct-origin access,
official Supercell endorsement of the proxy, ownership verification, or
commercial authorization. Player values, real tag and public IP were not
recorded. Further use remains subject to conditions above. Responsible operator
confirmed revocation of the previous key on `2026-10-01`; no identifier was
recorded and confirmation was not independently verified.

## Reproduction after blocker removal

Set the secret outside Git, then run without shell tracing or outputting its value:

```sh
test -n "$CLASH_ROYALE_API_TOKEN"
curl --path-as-is --fail-with-body --silent --show-error \
  -H "Authorization: Bearer $CLASH_ROYALE_API_TOKEN" \
  "https://api.clashroyale.com/v1/players/%23<REDACTED_VALID_TAG>"
```

Repeat with `/v1/cards`, `/v1/players/%23<REDACTED_VALID_TAG>/battlelog` and a safe nonexistent tag. Save only sanitized status, selected headers and top-level shape.

## Consequence for next tasks

`001-02` and `001-03` may use current RoyaleAPI Proxy route server-side under the
documented controls while proxy terms/contract remain operational debt. They must
not treat proxy observations as official API contracts or authorization. Direct
egress remains an alternative for future migration.
