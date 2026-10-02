# Meta strategy comparison — Clash Royale

- **Task:** `001-04`
- **Observed at:** `2026-10-02`
- **Status:** `completed with constraints`
- **Authority:** [Clash Royale API developer portal](https://developer.clashroyale.com/)
- **Operational route:** RoyaleAPI Proxy, server-side, replacing only API host.
- **Sanitization:** player tags, names, clan tags, card identifiers, timestamps and
  raw payloads are omitted. Counts, ranges and shapes are aggregate observations.

## Decision boundary

The MVP needs three separate contexts:

1. trophy-road/Arena bands, including intermediate players;
2. advanced ladder, useful as a separate high-skill context;
3. Path of Legend/ranked, never merged with trophy-road data.

The probes below show that the API can expose candidate decks and context from
profiles, clan members and encountered opponents. They do **not** establish a
representative Arena meta. No strategy in this sample demonstrates sufficient
mid-ladder coverage to claim “meta da Arena”.

## Evidence boundary and sources

- **Documented:** the official portal is the API authority; API keys are bound to
  allowed IPs; candidate surfaces and authentication are recorded in
  `evidences/api-surface.md`.
- **Observed:** the response shapes, counts, status codes and cache headers below
  came from controlled requests through the operational proxy.
- **Inferred:** request projections, deduplication approach and graph-growth
  bounds. They are not limits or service-level guarantees.
- **Pending:** direct-origin confirmation, authenticated official schema,
  pagination semantics for every candidate endpoint, and permission to retain or
  redistribute opponent/clan-derived data.

Official and policy sources consulted on `2026-10-02`:

- [API portal](https://developer.clashroyale.com/)
- [API docs](https://developer.clashroyale.com/api-docs/index.html)
- [Fan Content Policy](https://supercell.com/en/fan-content-policy/), last updated
  September 27, 2023
- [Terms of Service](https://supercell.com/en/terms-of-service/), effective
  November 6, 2024
- [RoyaleAPI Proxy documentation](https://docs.royaleapi.com/proxy.html), fetch
  returned `403` in this environment; proxy key handling, retention, limits, SLA
  and Supercell compatibility remain unresolved.

The proxy is transport only. It is not a meta dataset, API license, Supercell
endorsement or commercial authorization. No third-party meta API with a verified
API contract and redistribution license was selected.

## Controlled probe summary

All requests used the local rotated credential without printing or persisting its
value. No production collection or crawler ran.

| Probe | Result | Sanitized observation |
| --- | --- | --- |
| `GET /v1/locations` | `200` | 262 locations; `paging` present; cache observed between 138s and 225s. |
| `GET /v1/locations/57000000/rankings/players` | `404` | Player leaderboard unavailable through selected proxy route in this run. |
| `GET /v1/locations/{regional}/rankings/players` | `404` | Same result for one regional location; not treated as proof that official endpoint is removed. |
| `GET /v1/locations/57000000/rankings/clans` | `200` | 999 clan rows; `paging` present; clan-score range 120,504–139,921; cache 60s. |
| `GET /v1/locations/{regional}/rankings/clans` | `200` | 999 clan rows; `paging` present; clan-score range 66,456–138,838; cache 60s. |
| `GET /v1/locations/57000000/pathoflegend/players` | `404` | No ranked leaderboard sample available through selected proxy route. |
| `GET /v1/clans?name={sanitized}` | `200` | 640 search results; `paging` present; cache 60s. |
| `GET /v1/clans/{tag}/members` | `200` | 50 members in each sampled response; `paging` present; cache 120s. |
| Six member profile lookups from two ranked clans | `200` | 6/6 profiles resolved; trophies were 14,000 in both three-profile samples; current deck lengths were 7 or 8. |
| Five member profile lookups from one sampled player clan | `200` | 5/5 profiles resolved; trophies were 14,000; current deck length was 8. |
| Battle logs for three sanitized profiles | `200` | 6, 30 and 30 entries; 6, 30 and 30 unique opponent rows. PvP sample A covered team trophies 2,960–3,110 and opponents 2,964–3,090; Path of Legend sample B covered ranked values 1,713–1,921; sample C covered 14,987–15,000 in `trail`/`unknown`. |

The 999-row leaderboard responses were not expanded. Only the first two clan
records were used to select two member samples. This is a probe of endpoint shape,
not a dataset ingestion run.

## Strategy matrix

| Strategy | Coverage | Mid-ladder | Bias | Requests | Freshness | Terms | Dependency | Verdict |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| Global/regional player rankings | `rankings/players` returned `404` for global and one regional probe; no deck coverage measured. Clan rankings returned 999 rows but do not contain player decks. | **Unproven**; no observed player sample. | Expected top/active selection bias even if route becomes available. | At least 1 ranking request per location plus 1 profile request per selected player; exact page/limit contract pending. | Clan ranking cache 60s; player ranking freshness unknown. | Official surface candidate; direct origin and schema pending. | API key IP allowlist plus current proxy. | **Reject as sole source.** Never equate top ranking with Arena meta. |
| Path of Legend/ranked | `pathoflegend/players` returned `404`; battle log exposed 30 `pathOfLegend` entries for one advanced active profile. | **Not demonstrated**; ranked is a different context, not mid-ladder trophy-road coverage. | Advanced/active-player bias; mode must remain separate. | One battle-log request per tracked player per sync; leaderboard cost unavailable. | Battle log cache observed at 60s; no ranked leaderboard freshness. | Official candidate, contract pending; proxy observations not official contract. | API route availability and proxy. | **Useful context only.** Do not use as Arena substitute. |
| Clan search + clan members | Search returned 640 rows; two ranked-clan member calls returned 50 each; 11 sampled member profiles resolved. | **Not demonstrated:** all sampled member profiles had 14,000 trophies. | Clan activity, recruitment and geography bias; ranked clans strongly top-biased. | Per clan: 1 search + 1 detail + 1 members + `K` profile calls; `K=50` implies ~53 requests before refresh. | Search 60s, members 120s, profiles 6–60s observed. | Official candidate endpoints; retention/use of member data and proxy terms unresolved. | Proxy and API key allowlist. | **Supplement only.** Can enrich bounded cohorts, cannot prove Arena prevalence. |
| Controlled expansion via battle-log opponents | 6/30/30 entries produced 6/30/30 unique opponent rows; opponent decks, arena and mode fields were present in observed entries. | **Partial candidate access, not coverage:** sample follows tracked players and their activity. | Ego-network, activity, mode and opponent-selection bias; no representative prevalence. | One seed battle-log request yields candidates already; one-hop log expansion costs up to `U` additional requests (`U=30` observed); two hops can reach up to 900 before dedup. | Battle log cache observed at 60s; list window and eviction unresolved. | Opponent data retention, redistribution and proxy handling unresolved; no mass expansion allowed. | Battle-log window, polling and proxy. | **Preferred candidate fallback.** Bounded, user/context-seeded only; no global meta claim. |
| Third-party API with clear license/terms | No candidate source with verified API, license and redistribution terms was identified. RoyaleAPI Proxy only transported official-shaped requests. | **Unknown.** | Depends on vendor sampling and methodology; cannot assess without a source. | Unknown; must be supplied by vendor contract. | Unknown; vendor SLA/freshness pending. | **Not accepted** without explicit API/data/license terms and Supercell compatibility review. | Potential vendor lock-in and opaque retention. | **Do not use in MVP.** Re-open only with a documented license. |
| Hybrid: profile + battle log + bounded clan cohort | Combines observed profile context, encountered decks and optional small clan cohorts; preserves `PvP`, `pathOfLegend`, `trail` and `unknown` partitions. | **Honest only as personalized/contextual candidates.** No representative Arena distribution observed. | Makes bias visible but does not remove it. | Initial player sync: 1 profile + 1 battle log. Optional clan cohort: ~53 requests/clan for 50 profile enrichments. Opponent expansion is opt-in and bounded. | Uses observed endpoint caches; freshness is snapshot freshness, not a meta SLA. | Same unresolved API/proxy/data-retention gates; no third-party dataset dependency. | Proxy now; replaceable fixed-egress adapter later. | **Preferred MVP acquisition strategy with constraints.** Recommend decks from collection plus observed context; do not label result “Arena meta”. |

## Deduplication and stability by strategy

| Strategy | Deduplication | Stability |
| --- | --- | --- |
| Global/regional player rankings | Player tag within ranking response; combine with location/season when schema is available. | **Pending:** player endpoint returned `404`; paging and season identity are unconfirmed. |
| Path of Legend/ranked | Player tag + ranked mode/season; battle-log fallback requires opaque battle fingerprint. | **Low/pending:** leaderboard unavailable; battle log has no universal battle ID. |
| Clan search + clan members | Clan tag + player tag; deduplicate search hits before member/profile enrichment. | **Medium:** `paging` and 50-member samples observed; membership churn and official paging contract remain unresolved. |
| Battle-log opponents | Versioned opaque battle fingerprint; player tag prevents repeated opponent enrichment within bounded batch. | **Low/medium:** no stable `battleId`; observed ordering and window are not contracts. |
| Third-party API with clear license/terms | Vendor record ID plus source/version; exact key depends on a future contract. | **Unknown:** no source with verified API/license/terms selected. |
| Hybrid | Cross-source player tag plus source context; battle fingerprints remain separate from player/deck identity. | **Medium with constraints:** inherits proxy, endpoint, mode and battle-window risks from component strategies. |

## Candidate context and deduplication rules

- Preserve raw `arena`, `trophies`, `gameMode`, `type` and ranked fields when
  present. Never merge `PvP`, Path of Legend, `trail` or `unknown` only because
  each row contains eight cards.
- A profile or clan member can be deduplicated by player tag inside a controlled
  ingest batch. Tags remain operational identifiers and must not enter versioned
  evidence.
- Battle-log entries had no universal stable `battleId`; use the versioned opaque
  fingerprint already proposed in `battlelog-findings.md`. This is application
  idempotency, not official identity.
- A future `DeckCandidate` needs source, fetched time, mode/context, trophy/Arena
  snapshot and provenance. Deck prevalence must not be inferred from candidate
  count without sampling weights and a documented population.

## Request projections

These are planning calculations, not load tests or limits:

| Operation | Formula | Example |
| --- | --- | --- |
| User profile + battle log | `2 × tracked players` for initial snapshot | 2 requests/player; 5-minute polling means 12 calls/hour, 288/day, or 2,016/4,032/8,640 calls/player over 7/14/30 days before retries. Cache hits may reduce upstream work but do not change this request schedule. |
| Clan cohort | `3 + K` per clan (`search + detail + members + profiles`) | `K=50` → ~53 requests/clan; cache can reduce repeated reads but does not guarantee it. |
| Opponent candidates | `1 + U` for seed log plus one log per unique opponent | `U=30` observed → ~31 requests for one-hop log expansion; not executed. |
| Two-hop graph | `1 + U + U×U` upper-bound shape before dedup | With `U=30`, up to 931 requests; this is why graph expansion is prohibited in MVP. |
| Rankings enrichment | `1 + K` per ranking context/page | Player ranking route was not available in probe; page and enrichment limits remain unknown. |

For 7/14/30 days of active-player battle-log polling at 5 minutes, the planning
order is **2,016/4,032/8,640 requests per player**, excluding retries and outages.
This is a client-call projection; proxy/upstream cache behavior can change actual
origin work. It is not an ingestion recommendation for global data.

## Product gate and recommendation

**Meta por Arena is not resolved.** No tested strategy demonstrated reasonable,
observed intermediate-ladder coverage. The PvP battle-log sample did include one
intermediate band (2,960–3,110 trophies), but only six entries from one profile;
the 14,000-trophy result in all 11 clan-member profile samples shows top
concentration in this probe, consistent with selection bias rather than coverage.

For Phase 005, accept the following constrained product boundary:

1. Preferred acquisition: hybrid, user/context-seeded acquisition using profile
   snapshots and battle-log candidates, with optional bounded clan enrichment.
2. Product language: **Best Decks for Your Collection**, contextualized by the
   player's observed trophies/Arena/mode and recent encounters.
3. Explicitly unsupported claim: representative “best decks for the Arena” or
   global meta prevalence, until a licensed or independently validated source
   supplies intermediate coverage.
4. Fallback options retained: segment by observable trophy/mode context, obtain a
   licensed dataset, or redesign acquisition. No crawler, scraping or mass graph
   expansion is approved.

## Residual risks

1. Direct official-origin probes and authenticated official schemas remain blocked
   by the API-key IP allowlist/session boundary.
2. Proxy key handling, retention, limits, SLA and compatibility with Supercell
   developer/API agreements remain unresolved.
3. Endpoint availability and response shapes can change; observed `404` responses
   are not proof of official endpoint removal.
4. Opponent and clan-member data may create privacy, retention and permitted-use
   obligations not closed by API access alone.
5. A small controlled sample cannot establish prevalence, causal strength or
   freshness SLA; the observed intermediate PvP band is not population coverage.
