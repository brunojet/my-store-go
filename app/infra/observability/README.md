# Observability (infra)

This folder hosts HTTP observability primitives (request id, access logs, etc).

## Route labels (low cardinality)

Telemetry and access logs must avoid using the raw URL path as a `route` label.
When the router template/pattern is missing (e.g. 404/unmatched), the route is
recorded as a fixed value: `unmatched`.

- `requestid`: request correlation id helpers stored in `context.Context`
- `ginmw`: Gin middlewares
- `httpmw`: net/http middlewares
