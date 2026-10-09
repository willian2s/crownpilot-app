# Sanitized fixtures

No application fixture is defined in `002-03`. Domain tables and their
sanitized fixtures belong to `002-06`; `db:fixtures` therefore fails closed
until a versioned `.sql` fixture exists. It never reports an empty success.

Fixtures must contain no real Firebase UID, Player Tag, provider payload,
credential, or production data. They run only after EF Core migrations and
versioned security SQL.
