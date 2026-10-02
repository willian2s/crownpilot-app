# Battle log findings — Clash Royale API

- **Task:** `001-03`
- **Observed at:** `2026-10-02`
- **Status:** `completed with constraints`
- **Authority:** [Clash Royale API developer portal](https://developer.clashroyale.com/)
- **Operational route:** RoyaleAPI Proxy, server-side, replacing only API host.
- **Sanitization:** tags, names, clan values, card names/IDs, timestamps and raw
  payloads are omitted. Counts and shapes are aggregate observations.

## Confidence and source boundary

The official portal is authoritative for the endpoint contract. The public
documentation session available during this run did not expose an authenticated
Swagger schema. Therefore:

- **Documented:** the official portal is the API source and requires authenticated
  requests; the project already recorded IP/egress binding for API keys.
- **Observed:** response shapes, counts, fields, modes and cache headers below,
  obtained through the selected proxy with three sanitized profiles.
- **Inferred:** polling estimates, deduplication fingerprint and product policy
  derived from observed count/window behavior.
- **Not available:** direct-origin confirmation, official pagination contract,
  stable battle ID, replay/timing data and explicit Hero deployment semantics.

Proxy observations do not prove official-origin behavior, Supercell endorsement,
ownership verification or commercial authorization. Proxy terms, key handling,
retention, limits and SLA remain separate risks from `001-01`.

## Probe summary

| Sequence | Sanitized sample | Status | Entries | Observed `type` values | Span/order | Cache-Control |
| ---: | --- | ---: | ---: | --- | --- | --- |
| 1 | A — low activity | 200 | 6 | `PvP` | ~5.9 days; descending | `max-age=22` |
| 2 | B — recent active | 200 | 30 | `pathOfLegend` | ~96 minutes; descending | `max-age=21` |
| 3 | C — mixed/low activity | 200 | 30 | `trail` (29), `unknown` (1) | ~6.8 days; descending | `max-age=22` |
| 4 | B — repeat after 65s | 200 | 30 | `pathOfLegend` | ~96 minutes; descending | `max-age=49` |
| 5 | B — repeat second call | 200 | 30 | `pathOfLegend` | ~96 minutes; descending | `max-age=60` |

The initial three calls returned HTTP `200`, JSON arrays and `battleTime` on
every entry. B was queried twice more after 65 seconds: each repeat returned 30
entries, with 30/30 overlap against the initial B candidate keys and the same
observed time span. No new match occurred during the interval; this proves
neither freshness guarantee nor a fixed time window.

## Window, pagination and ordering

- **Observed size:** 6 or 30 entries. A 30-entry response was the largest shape
  seen; this is an observation, not a documented hard maximum.
- **Window:** count-based behavior is more likely than fixed-duration behavior:
  low-activity samples reached roughly 6–7 days, while the recent active sample
  covered roughly 96 minutes with the same 30-entry count.
- **Pagination/cursor:** no `paging`, cursor, `next` link or pagination query was
  observed. Official pagination behavior remains **unresolved** without the
  authenticated schema. Do not implement backfill or assume a cursor from this
  evidence.
- **Ordering:** all 66 entries were ordered newest first by `battleTime` in the
  three-profile probe. Treat ordering as observed, not as an immutable contract.
- **Backfill:** not established. Polling cannot be assumed to recover entries
  after they leave the observed list; an official cursor/backfill contract still
  requires confirmation.

## Observed response shape

Each entry was an object with these top-level fields or structures:

| Area | Observed fields | Classification |
| --- | --- | --- |
| Time/type | `battleTime`, `type` | available in observed shape; semantics of `type` require mode mapping |
| Competitive context | `gameMode.id/name`, `arena.id/name/rawName`, `leagueNumber`, `isLadderTournament`, `isHostedMatch`, `deckSelection` | observed; optionality by mode remains unresolved |
| Sides | `team[]`, `opponent[]`; one participant per side in every observed entry | observed; do not assume all modes are 1v1 |
| Decks | `cards[]` on each participant; exactly 8 rows in observed entries | observed; validate cardinality before treating as deck |
| Support | `supportCards[]` on each participant; one row in observed entries | observed and separate from `cards[]` |
| Towers/result signals | `crowns`, `kingTowerHitPoints`, `princessTowersHitPoints`, participant `trophyChange`, `startingTrophies` | available as raw signals; no explicit `result`/`winner` field observed |
| Opponent | participant tag/name, optional clan, cards and support cards | observed; sanitize and apply retention policy |
| Card detail | `id`, `name`, `level`, `maxLevel`, `rarity`, `elixirCost`, optional `evolutionLevel`, `maxEvolutionLevel`, `starLevel`, `iconUrls` | observed; Evo/Hero semantics not fully resolved |
| Event metadata | `eventTag` present on all C entries and absent on A/B | optional; not a universal battle ID |

`globalRank` was present as null in observed participant rows. Tower hit-point
arrays were nullable and varied in length, so consumers must preserve null versus
array rather than assume fixed tower count.

## Modes and MVP treatment

Observed `type` values were `PvP`, `pathOfLegend`, `trail` and `unknown`.
Observed values alone do not establish official semantic labels for `trail` or
`unknown`.

| Mode bucket | MVP treatment | Reason |
| --- | --- | --- |
| `PvP` | include, partition by raw type/game mode | observed regular battle shape |
| `pathOfLegend` | include as separate competitive context | must not mix with trophy-road aggregates |
| `trail` | ingest raw, exclude from default aggregates until mapped | semantics not confirmed |
| `unknown` | ingest raw, quarantine from analysis | cannot infer mode from field name |

Never combine modes only because both contain eight card rows. Keep raw
`type/gameMode/arena` alongside any normalized mode.

## Deck, result, Evo/Hero and opponent coverage

- **Decks:** team and opponent participant rows exposed eight `cards[]` entries
  in every observed entry, plus a separate one-row `supportCards[]`. This is
  sufficient for coarse deck snapshots in this sample. It does not authorize a
  universal eight-card assumption for unseen modes or future shapes.
- **Result:** no explicit `result`, `winner` or `battleResult` field was observed.
  `crowns` and `trophyChange` exist per participant and support a derived result
  classification, but that classification is **inferred**, not an API contract.
- **Towers:** crown counts and king/princess tower hit-point fields are present,
  with nullable/variable tower arrays.
- **Evolution:** `evolutionLevel` and `maxEvolutionLevel` appeared on subsets of
  participant card rows. Presence is not proof of active Evolution deployment;
  preserve raw optional fields.
- **Hero:** no explicit `heroId`, `heroLevel`, `heroOwned` or deployment marker was
  observed. Some card rows exposed `iconUrls.heroMedium`, but that asset field
  cannot establish Hero identity or use. Hero remains **unresolved/unavailable as
  an explicit battle semantic**.
- **Opponent:** opponent participant, cards, support cards, trophy change and
  context were observed. This enables coarse matchup snapshots, not replay,
  timing, placement or in-match telemetry.

## Deduplication and idempotency

No stable top-level `id` or `battleId` was observed. `eventTag` was optional and
was not universal. Candidate inferred key uses versioned canonical JSON before
hashing:

```text
canonical_json({
  "v": 1,
  "battleTime": battleTime,
  "type": type,
  "gameModeId": gameMode.id,
  "teamTags": sorted(team participant tags),
  "opponentTags": sorted(opponent participant tags)
}) -> sha256
```

The three samples had unique `battleTime` values and unique candidate keys. A
65-second repeat of sample B retained the same 30/30 keys. Store only an opaque
fingerprint where possible; do not expose or version participant tags merely to
deduplicate. If participant tags are unavailable in a future shape, fall back to
the same canonical JSON strategy with `battleTime + type + gameMode.id + sorted
card IDs per side + arena.id`, mark the key lower-confidence and retain collision
telemetry. This is an application idempotency strategy, not an official battle
identity.

## Polling estimate

Sample B contained 30 entries across approximately 96 minutes, which has 29
observed intervals and an estimated rate of roughly 18.1 entries/hour. The table
models overflow of an observed
30-entry response using a Poisson assumption. It is an estimate, not a load test
or guarantee.

| Poll cadence | Expected entries between syncs | Estimated probability of >30 entries between syncs* | Assessment |
| ---: | ---: | ---: | --- |
| 1 minute | 0.30 | ~0% | strongest gap protection; higher request cost |
| 5 minutes | 1.51 | ~0% | recommended MVP cadence for active tracked players |
| 15 minutes | 4.54 | ~0% | likely sufficient at observed rate, less burst margin |
| 1 hour | 18.14 | ~0.37% | not reliable for high-activity history; burst risk material |

\*Probability is only a calculation from sample B rate and a 30-entry cap. It
does not model sessions, bursts, outages, cache behavior or a different player.
At a rate above the observed sample, overflow risk increases directly. A sync
gap longer than the returned window can lose battles permanently.

## Product answers and contract

- **Can we reconstruct own history from snapshots?** Yes, conditionally. Polling
  plus opaque-key upserts can build a best-effort private history while entries
  remain in the response. No arbitrary historical backfill can be assumed.
- **Is loss between syncs real?** Yes, under the observed list behavior. The
  response returned at most 30 entries in these samples and exposed no cursor;
  the official hard limit remains unresolved. More entries than observed between
  syncs, an outage, or an unobserved mode transition can create irreversible gaps.
- **Which modes enter MVP?** Keep raw all modes, analyze `PvP` and
  `pathOfLegend` separately, and filter/quarantine `trail`/`unknown` pending
  semantic confirmation.
- **Does battle log discover meta candidates?** It supplies encountered deck
  samples and context, but has strong player/activity/opponent-selection bias. It
  is not sufficient alone for representative Arena meta prevalence.
- **Is there enough for future matchup analysis?** Enough for coarse deck,
  opponent, mode, arena, tower and derived outcome analysis. Not enough for replay,
  card timing, placement, sequencing, causal matchup claims or explicit Hero
  deployment.
- **Can snapshots be stored?** Technically yes, but retention, privacy, proxy
  terms and permitted data use remain gates. Store normalized minimum fields,
  source/fetched time and dedup fingerprint; avoid raw payload retention by
  default.

## Residual limitations

1. Direct official-origin probes remain blocked by API-key IP allowlist; all live
   observations here use selected proxy transport.
2. Authenticated official schema, pagination behavior, hard entry limit and mode
   semantics remain pending.
3. No stable battle ID or explicit result/winner field was observed.
4. Hero representation and deployment remain unresolved.
5. Polling estimates use one active sample and must not be promoted to an SLA.
6. No production polling, persistence or scale collection was performed.
