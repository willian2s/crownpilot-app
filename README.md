# CrownPilot

> **Play the right deck. Upgrade the right cards.**

CrownPilot is a personalized Clash Royale companion focused on one question:

> **What is the best deck I can play right now, in my current Arena, with the account I actually have?**

Instead of showing the same global meta to every player, CrownPilot is designed to combine the player's current context — Arena / trophy range, card collection, card levels, Evolutions / Heroes when available, and current meta data — to produce recommendations that are useful for that specific account.

## The idea

Most Clash Royale tools are excellent at answering questions such as:

- What decks are popular right now?
- What is the win rate of this deck?
- What are top players using?
- What cards does this player have?

CrownPilot aims to connect those pieces and answer a different set of questions:

- **What should I play right now?**
- **Which competitive decks can I already build?**
- **Which decks am I one or two upgrades away from?**
- **What should I upgrade next to unlock the most useful decks?**
- **How does the meta around my Arena / trophy range change what is good for me?**

The goal is to make the player's account the center of the experience.

## Product loop

CrownPilot is planned around four layers.

### Discover

**Best Decks for You**

Recommend proven decks using the player's:

- current Arena / trophy range;
- owned cards;
- card levels;
- available Evolutions / Heroes;
- relevant meta data;
- confidence and sample size of the underlying data.

The first recommendation engine should be deterministic and based on observed decks, not on decks invented by an LLM.

### Progress

**Upgrade Planner**

Show which upgrades have the highest impact on the player's account.

Examples:

- decks that are ready to play now;
- decks that are one upgrade away;
- upgrades that unlock several competitive decks;
- targets worth preparing for the next Arena / progression step.

A useful mental model is:

```text
NOW    -> what can I play today?
NEXT   -> what should I upgrade next?
TARGET -> what am I building toward?
```

### Optimize

Analyze the player's existing decks and compare them with proven variants and the matchups that matter in the player's environment.

This comes after the recommendation engine is trustworthy.

### Coach

AI can eventually explain CrownPilot's deterministic analysis in natural language, answer questions about the player's account, and help interpret performance.

AI is intentionally a layer on top of the product rather than the source of truth.

> **Deterministic by default. AI when it adds value.**

## Product principles

### Player-first personalization

A globally strong deck is not automatically the best deck for a specific player.

Recommendations should account for the real state of the player's account.

### Persistent identity

The player should submit a Clash Royale Player Tag once as a public-profile
reference; this does not verify account ownership.

After signing in on another device, CrownPilot should recover the public profile
reference linked by that user. The link remains `unverified`; CrownPilot does not
claim ownership or control of the external account.

### Data-driven recommendations

The core recommendation engine should start from decks with real evidence behind them and rank those candidates for the player.

### Explainability

Every recommendation should be able to answer:

> **Why is CrownPilot recommending this to me?**

The recommendation score must remain understandable and testable.

### AI is optional

Features that can be solved deterministically should not require an LLM.

This keeps the product:

- predictable;
- testable;
- faster;
- cheaper to operate;
- useful even without an AI provider.

## What CrownPilot is not

CrownPilot is not intended to be:

- an in-match overlay;
- an elixir or cycle tracker;
- a bot or gameplay automation tool;
- a replacement client for Clash Royale;
- an LLM that invents eight-card decks without data behind them.

The initial focus is account intelligence, deck discovery, progression, and post-game analysis.

## Project status

CrownPilot is currently in the **discovery and foundation** stage.

The first milestone is to validate the available data, establish the account model, and prove that a useful personalized recommendation engine can be built before expanding the product.

See the complete roadmap:

- [Project Roadmap](docs/roadmap/crownpilot-roadmap.md)

## Documentation model

The project follows a separation between planning and implementation:

- **Roadmap** — what we are building, in what order, and why;
- **Specs** — how a specific phase or feature will be implemented;
- **ADRs** — architectural decisions that must remain explicit over time;
- **Repository state** — the implementation is the final source of truth for what actually exists.

Future implementation specs should be created only when the previous phase has enough evidence to define them responsibly.

## Commercialization

CrownPilot should not assume that a conventional paid SaaS model is automatically allowed.

Before implementing subscriptions, paid AI features, or other paywalled functionality, the project must explicitly validate the commercial model against Supercell's current Fan Content Policy and any applicable developer agreements.

This is a product and architecture gate, not a task to postpone until launch.

## Disclaimer

CrownPilot is an unofficial fan project and is not endorsed by Supercell.

Clash Royale, Supercell, and related game assets and trademarks belong to their respective owners. Any use of game-related assets or data must follow Supercell's applicable policies and developer agreements.

See: [Supercell Fan Content Policy](https://supercell.com/en/fan-content-policy/)

## License

A project license has not been selected yet.
