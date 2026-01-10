# Prompt: API de Aplicativos (end-to-end) — Gin + GORM + SQLite em memória

Você é um assistente de engenharia de software. Implemente end-to-end a **API de Aplicativos** (CRUD) em Go, indo da requisição HTTP até o banco, usando:

- **Gin** para a camada HTTP (adapter)
- **GORM** para persistência
- **SQLite em memória** (para desenvolvimento/testes)

A implementação deve respeitar o estilo **monólito modular** e manter o `core` reutilizável: o `core` **não deve depender de Gin**. O adapter HTTP (Gin) chama serviços/use cases do `core`.

## Contexto do repositório

- Repo: https://github.com/brunojet/my-store-go
- Basepath da aplicação: `app/`
- Módulos de domínio existem como pastas (`store-*`) e devem continuar **auto contidos**.
- O `core` é a peça comum/reutilizável (pode virar uma lib no futuro).

## Requisitos funcionais

Criar a entidade **Aplicativo** com os campos:

- `nome` (string)
- `descricao` (string)
- `codigo_parceiro_externo` (string)
- Auditoria:
  - `created_at`
  - `updated_at`

Endpoints (5):

- `GET /v1/apps` (lista)
- `GET /v1/apps/:id` (detalhe)
- `POST /v1/apps` (criar)
- `PATCH /v1/apps/:id` (atualizar parcial)
- `DELETE /v1/apps/:id` (remover)

## Restrições e decisões arquiteturais

- O `core` **não importa Gin**. Nenhum `*gin.Context` em assinaturas do `core`.
- Use `context.Context` no `core`.
- A camada HTTP (Gin) faz:
  - parse/validação básica do request
  - chama o service/use case do `core`
  - formata resposta HTTP
- A camada de persistência usa GORM com SQLite em memória.
- Para SQLite in-memory funcionar bem com GORM e múltiplas conexões, use DSN:
  - `file::memory:?cache=shared`
  - e configure pool para `SetMaxOpenConns(1)`.

## Organização de pastas (sugestão)

Obs.: ajuste se você já tiver convenções melhores, mas mantenha as dependências apontando “para dentro”.

- `app/core/models/apps/app.go` (modelo)
- `app/core/dtos/apps/app_dto.go` (DTOs request/response)
- `app/core/repositories/apps/repository.go` (interface/port)
- `app/core/repositories/apps/sqlite/repository.go` (implementação GORM/SQLite)
- `app/core/services/apps/service.go` (use case)
- `app/cmd/server/main.go` (wire-up: conecta DB, cria repo/service, registra rotas)
- `app/infra/http/gin/apps/routes.go` (router/adapter Gin)
- `app/infra/http/gin/apps/handlers.go` (handlers Gin)

Se preferir manter tudo “core autocontido”, o adapter Gin pode ficar em `app/core/infra/http/gin/...`, mas **o mais importante** é: `core` não deve depender de Gin, apenas o adapter depende do core.

## Modelo e auditoria

- Use `gorm.Model` (inclui `ID uint`, `CreatedAt`, `UpdatedAt`, `DeletedAt`).
- Como o requisito pede `created_at`/`updated_at`, você pode:
  - usar `CreatedAt`/`UpdatedAt` internamente
  - expor no JSON como `created_at` e `updated_at` via DTO.

## Contratos (DTOs)

- `POST /v1/apps` request:
  - `nome` (obrigatório)
  - `descricao` (opcional)
  - `codigo_parceiro_externo` (obrigatório)

- `PATCH /v1/apps/:id` request:
  - campos opcionais (somente os presentes no JSON devem ser atualizados)

- Responses:
  - `id`
  - `nome`, `descricao`, `codigo_parceiro_externo`
  - `created_at`, `updated_at`

## Regras mínimas (para o exemplo)

- Validar que `nome` e `codigo_parceiro_externo` não são vazios.
- `GET /v1/apps/:id` retorna `404` se não existir.
- `DELETE` idempotente: pode retornar `204` mesmo se já não existir (ou `404` — escolha uma e documente).

## Erros e status codes

- `POST` sucesso: `201` + body
- `GET` lista/detalhe sucesso: `200` + body
- `PATCH` sucesso: `200` + body
- `DELETE` sucesso: `204` sem body
- validação: `400`
- não encontrado: `404`
- erro inesperado: `500`

## Implementação (passo a passo)

1. Adicionar dependências:
   - `github.com/gin-gonic/gin`
   - `gorm.io/gorm`
   - `gorm.io/driver/sqlite`

2. Criar modelo GORM em `core/models/apps`.

3. Criar DTOs em `core/dtos/apps`.

4. Criar interface do repositório (port) em `core/repositories/apps`:
   - `List(ctx)`
   - `Get(ctx, id)`
   - `Create(ctx, *App)`
   - `Patch(ctx, id, patch)` (ou `Update` com campos opcionais)
   - `Delete(ctx, id)`

5. Criar implementação SQLite/GORM em `core/repositories/apps/sqlite`.
   - Use `*gorm.DB` injetado.

6. Criar service/use case em `core/services/apps`.
   - Recebe repo interface.
   - Aplica validações e regras.

7. Criar adapter HTTP (Gin) em `app/infra/http/gin/apps`.
   - Bind JSON → chama service → responde.

8. Atualizar `app/cmd/server/main.go`.
   - Inicializa SQLite in-memory
   - `AutoMigrate(&models.App{})`
   - Sobe Gin na porta `:8080` (ou env `PORT`)
   - Registra rotas `/v1/apps`.

9. Adicionar um arquivo de testes mínimo (opcional, mas recomendado):
   - teste do repositório com SQLite em memória
   - teste de handler com `httptest` (ou smoke test)

## Checklist de validação

- `go test ./...` dentro de `app/`
- `curl` (exemplos):
  - `POST /v1/apps` cria
  - `GET /v1/apps` lista
  - `PATCH /v1/apps/:id` atualiza
  - `DELETE /v1/apps/:id` remove

## Saída esperada

- Liste os arquivos criados/modificados.
- Mostre o código (ou trechos principais) de:
  - modelo
  - interface do repositório
  - service
  - handlers Gin
  - main (wire-up)
- Explique como a separação core/adapter evita acoplamento com Gin.
