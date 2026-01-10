# Store Provider OpenAPI

- Spec: [openapi.yaml](openapi.yaml)
- Schemas/components: [schemas.yaml](schemas.yaml)

## How to view

- Swagger UI (Docker):
  - Serve this folder (or the YAML files) via any static file server.
  - Point Swagger UI to `openapi.yaml`.

## Notes

This contract matches the current implementation in:
- routes: `store-provider/routers/apps_routes.go`
- handler plumbing: `infra/http/ports/CRUDHandler`
- DTOs: `core/dtos/app_dto.go`

Error response shape is always:

```json
{"error":"..."}
```
