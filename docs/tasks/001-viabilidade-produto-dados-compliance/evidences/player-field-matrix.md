# Player field matrix — profile and collection

- **Task:** `001-02`
- **Observed at:** `2026-10-02`
- **Status:** `completed with constraints`
- **Samples:** three sanitized public profiles queried server-side through the
  currently selected RoyaleAPI Proxy route. Profile A was intermediate and had
  no clan object; profiles B/C were advanced, had clan objects, and exposed
  different collection/progression coverage. Real tags, names and clan values
  are intentionally omitted.
- **Authority:** [Clash Royale API developer portal](https://developer.clashroyale.com/)
- **Operational route:** `https://proxy.royaleapi.dev`, host substitution only;
  proxy observations do not promote third-party behavior to an official API
  contract.

## Probe summary

| Probe | Result | Sanitized observation |
| --- | --- | --- |
| `GET /v1/cards` | `200` | `items[123]`, `supportItems[4]`; catalog item IDs and names present; `maxLevel` varies by rarity; `maxEvolutionLevel` present on 55/123 items. |
| `GET /v1/players/{tag}` profile A | `200` | `cards[73]`, `currentDeck[8]`, no clan object; `evolutionLevel` on 3 collection cards and 0 deck entries; support cards present. |
| `GET /v1/players/{tag}` profile B | `200` | `cards[123]`, `currentDeck[8]`, clan object; `evolutionLevel` on 45 collection cards and 4 deck entries; support cards present. |
| `GET /v1/players/{tag}` profile C | `200` | `cards[123]`, `currentDeck[8]`, clan and league-statistics objects; `evolutionLevel` on 39 collection cards and 7 deck entries; support cards present. |

Selected response headers: JSON content type on all probes; `Cache-Control`
was `max-age=7` for catalog and `max-age=36` for all three profile responses
in this run. Cache values are observations, not freshness guarantees.

## Field matrix

| Requirement | Endpoint | Raw field | Semantics | Optionality | Freshness | MVP status | Evidence |
| --- | --- | --- | --- | --- | --- | --- | --- |
| Card capability/catalog | `GET /v1/cards` | `items[]` | Catalog capability rows; 123 rows observed. | Required endpoint; row shape observed | Proxy cache observed; catalog refresh policy unresolved | `available` | Three profile card IDs were compared against 123 catalog IDs. |
| Card identity | `GET /v1/cards`, `GET /v1/players/{tag}` | `id`, `name` | Stable-looking card identity fields; names omitted from evidence. | Present on every observed catalog/card/deck row | Snapshot-time | `available` | `id` integer and `name` string on all observed rows. |
| Current card level | `GET /v1/players/{tag}` | `cards[].level` | Integer level returned for each observed collection row. | Present on 73/73, 123/123 and 123/123 rows | Snapshot-time | `available` | Observed ranges differed by profile; do not infer readiness score. |
| Card level ceiling | `GET /v1/cards`, `GET /v1/players/{tag}` | `maxLevel` | Ceiling associated with card/rarity; not one global scale. | Present on all observed catalog/card/deck rows | Snapshot-time/catalog dependent | `available` | Catalog ceilings varied across common, rare, epic, legendary and champion rows. |
| Card quantity | `GET /v1/players/{tag}` | `cards[].count` | Numeric quantity field; zero values observed. Exact ownership semantics are not established by this probe. | Present on all observed collection/deck rows | Snapshot-time | `unresolved` | Never treat row presence or nonzero count as verified account ownership without contract evidence. |
| Collection coverage | `GET /v1/players/{tag}` vs `GET /v1/cards` | `cards[]` | Returned collection representation is profile-dependent: 73/123 catalog IDs for A and 123/123 for B/C. | Required array in all three samples; completeness not guaranteed | Snapshot-time | `unresolved` | MVP must preserve returned rows and coverage metadata; cannot claim every catalog card is owned. |
| Evolution capability | `GET /v1/cards`, `GET /v1/players/{tag}` | `maxEvolutionLevel` | Maximum evolution level/capability field where emitted. | Optional: 55/123 catalog rows; 44/73, 55/123 and 55/123 profile rows | Snapshot-time/catalog dependent | `available` | Presence is not equivalent to player ownership or deployment. |
| Evolution state/ownership signal | `GET /v1/players/{tag}` | `cards[].evolutionLevel` | Optional integer on subset of card rows; ownership meaning needs product contract confirmation. | Optional: 3, 45 and 39 rows in A/B/C | Snapshot-time | `unresolved` | Values 1..3 observed. No explicit `evolutionOwned` boolean was observed. |
| Current deck | `GET /v1/players/{tag}` | `currentDeck[]` | Eight deployed card rows in each sample; IDs were subsets of `cards[]`. | Present in all samples; support deck separate | Snapshot-time | `available` | Deployment must remain separate from collection. |
| Deck level/quantity | `GET /v1/players/{tag}` | `currentDeck[].level`, `currentDeck[].count` | Deck-context values can differ from matching collection values. | Present on all observed deck rows | Snapshot-time | `available` | Profile A differed on level/count for all eight matched deck IDs; B/C matched in this comparison. Do not merge by ID without preserving source context. |
| Evolution deployment | `GET /v1/players/{tag}` | `currentDeck[].evolutionLevel` | Optional evolution field emitted on deployed rows in B/C; active deployment semantics are not proven. | Optional: 0/8, 4/8 and 7/8 | Snapshot-time | `unresolved` | Field presence indicates row data only; do not claim active Evolution deployment. |
| Hero/capability representation | `GET /v1/cards`, `GET /v1/players/{tag}` | hero-specific fields | No explicit `hero*` field was observed in catalog, cards, deck or top-level profile fields. Champion rarity exists but is not proof of Hero semantics. | Unavailable in observed shape | Unresolved | `unavailable` | Do not synthesize Hero ownership/deployment from rarity or card names. |
| Support capability | `GET /v1/cards` | `supportItems[]` | Separate support catalog with id/name/maxLevel/rarity/icon fields. Relationship to Hero product semantics was not proven. | Present in catalog response | Catalog snapshot | `available` | Four support items observed; no hero marker. |
| Support ownership/state | `GET /v1/players/{tag}` | `supportCards[]` | Separate support-card rows with id/name/level/maxLevel/count/rarity/icon fields. | Present in all three samples; global optionality unresolved | Snapshot-time | `available` | No `evolutionLevel`, `heroLevel`, `heroId` or `heroOwned` field observed. |
| Support deployment | `GET /v1/players/{tag}` | `currentDeckSupportCards[]` | Separate deployed support-card rows. | Present in all samples; one row each | Snapshot-time | `available` | Must not be conflated with `currentDeck[]` or Hero ownership. |
| Competitive context | `GET /v1/players/{tag}` | `arena.id`, `arena.name`, `arena.rawName` | Arena object gives current competitive context labels/identifier. | Present in all samples | Snapshot-time | `available` | Values and names omitted; exact mode semantics remain to be separated where needed. |
| Trophy context | `GET /v1/players/{tag}` | `trophies`, `bestTrophies` | Current and historical-high trophy integers. | Present in all samples | Snapshot-time | `available` | Sufficient for basic context; not a substitute for full progression model. |
| Account progression | `GET /v1/players/{tag}` | `collectionLevel`, `kingTowerLevel` | Current progression integers observed on all three profiles. | Present in all samples | Snapshot-time | `available` | Use as raw fields; no cross-account readiness normalization in this task. |
| Mode/season progression | `GET /v1/players/{tag}` | `progress`, `*PathOfLegend*`, `leagueStatistics` | Mode/season-specific and dynamic progression structures. | Optional/nullable; league statistics absent in A/B and present in C | Snapshot-time and season dependent | `optional` | Preserve only after separate contract review; do not assume stable keys. |
| Clan context | `GET /v1/players/{tag}` | `clan` | Optional clan object with badge/name/tag when present. | Absent/null in A; object in B/C | Snapshot-time | `optional` | Clan is not required for profile/collection MVP. |
| Legacy progression | `GET /v1/players/{tag}` | `legacyTrophyRoadHighScore` | Historical/legacy-named field observed in all samples. | Present in all samples | Historical; semantic freshness unresolved | `unresolved` | No deprecation announcement was established; name and semantics require explicit validation before use. |

## Answers to validation questions

- **Possuímos todos os cards do jogador?** Não como contrato geral. A amostra
  retornou 73/123 catalog IDs for profile A and 123/123 for B/C. `count=0`
  occurred, so array membership and quantity cannot be promoted to verified
  ownership semantics.
- **Sabemos o nível utilizável?** Temos raw `level` and `maxLevel` for observed
  card/deck rows. This is enough for raw snapshot; cross-rarity normalization and
  readiness remain out of scope.
- **Evo/Hero necessário está disponível?** Evolution capability is available as
  optional raw fields, but ownership and active deployment remain `unresolved`.
  No explicit Hero ownership or deployment field was found; Hero is
  `unavailable` in the observed shape.
- **Sabemos contexto competitivo básico?** Yes: `arena`, `trophies`,
  `bestTrophies`, `collectionLevel` and `kingTowerLevel` were observed.
- **Progressão indisponível?** Stable raw progression is available, but dynamic
  `progress`, Path of Legend and league structures need separate contracts.
- **Ownership é necessário para MVP?** Public read-only profile association is
  sufficient when represented as `subjectType: public_profile` and
  `ownershipStatus: unverified`. Features claiming account ownership, exclusive
  control, actions or notifications remain blocked until an official mechanism
  is confirmed.

## Classification and limitations

- No field was confirmed deprecated by this probe. `legacyTrophyRoadHighScore`
  is flagged as legacy/unstable, not declared deprecated.
- `available` means returned consistently in this sample; it does not mean
  officially documented or semantically complete.
- `unresolved` is retained wherever names, rarity or numeric presence could lead
  to an ownership/Hero inference not proven by the observed payload.
- The proxy resolved current egress, but its key handling, retention, limits,
  SLA and compatibility with Supercell agreements remain separate risks.
