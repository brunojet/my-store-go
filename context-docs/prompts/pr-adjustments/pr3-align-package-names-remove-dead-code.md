# PR3 — Align package names with folder structure and remove dead code

## Objetivo
Reduzir atrito cognitivo no módulo `app/infra/http`:
1) alinhar o nome do pacote dos adapters Gin com o diretório (`adapters/gin`), e
2) remover código aparentemente não utilizado que cria duplicidade de “formas oficiais” de construir runtime.

## Contexto
- Arquivos em `app/infra/http/adapters/gin/` declaram `package base`, o que é inconsistente com o path e força imports/aliases confusos.
- Em `app/infra/http/adapters/nethttp/router_adapter.go` existe `ChiApp`/`NewChiApp`, mas o runtime canônico parece ser `NewRuntime`.

## Escopo
- Renomear o pacote dos adapters Gin.
- Remover `ChiApp` se não houver uso real.

## Mudanças propostas (cirúrgicas)
1. **Renomear package dos adapters Gin**
   - Trocar `package base` por `package gin` ou `package ginadapter` em todos os arquivos de `app/infra/http/adapters/gin/`:
     - `context_adapter.go`
     - `cors.go`
     - `router_adapter.go`
     - `runtime.go`
   - Atualizar imports e referências em:
     - `app/infra/http/http_manager.go`
     - qualquer outro consumidor interno de `infra`.

2. **Remover `ChiApp` não utilizado**
   - Confirmar que `NewChiApp` não tem usos.
   - Remover `type ChiApp` e `func NewChiApp()` de `app/infra/http/adapters/nethttp/router_adapter.go`.

## Critérios de aceite
- `go test ./...` passa.
- Imports ficam coerentes com o path (sem alias inesperado como `base`).
- Não há símbolos não utilizados no módulo.

## Fora de escopo
- Mudanças em APIs públicas de `contracts.Router` / `contracts.Context`.
- Refatorar `HTTPManager`.

## Notas
- Se o pacote for renomeado para `gin`, garantir que não conflite com `github.com/gin-gonic/gin` nos imports.
