# Learning 001 — Minimal API composition root

Bootstrap uses ASP.NET Core Minimal API because this task needs only a small
health endpoint and one generated contract route. `Program.cs` is the composition
root: it registers ProblemDetails, health checks, OpenAPI, and the Infrastructure
boundary before building the process. It does not open a database connection or
validate a Firebase token; Local contract fixtures only prove bearer pipeline
behavior.

`MapHealthChecks("/health/live")` uses an empty predicate, so liveness answers
whether the process is serving without turning provider availability into process
liveness. `AddOpenApi` and `MapOpenApi` produce the JSON consumed by contract and
smoke checks; there is no second YAML contract.

Later tasks can add adapters and use cases behind the existing layer references.
They must preserve this composition-root rule: HTTP concerns stay in `Api`,
provider concerns stay in `Infrastructure`, and `Domain` remains dependency-free.
