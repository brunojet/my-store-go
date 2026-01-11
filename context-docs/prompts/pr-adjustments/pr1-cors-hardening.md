# PR1 — Harden CORS defaults and add production guardrails

## Objetivo
Tornar o comportamento padrão de CORS **seguro para produção** no módulo `app/infra/http`, mantendo compatibilidade com a arquitetura atual (contracts/ports/adapters) e preservando a capacidade de habilitar CORS explicitamente por configuração.

## Contexto
Hoje o CORS é aplicado via:
- Gin: `app/infra/http/adapters/gin/cors.go`
- net/http (chi): `app/infra/http/adapters/nethttp/cors.go`

A configuração trafega como `types.CORSConfig` (`app/infra/http/types/cors.go`) e é interpretada no bootstrap (ex.: `app/cmd/server/config.go`). O risco é o default permissivo (e.g. origens amplas) ir para produção sem intenção.

## Escopo
- Ajustar defaults e/ou guardrails para evitar que CORS fique implicitamente aberto em produção.
- Melhorar clareza de documentação e exemplos de configuração.

## Mudanças propostas (cirúrgicas)
1. **Defaults production-safe**
   - Opção A (preferida): `CORS_ENABLED` default = `false`.
   - Opção B: manter `CORS_ENABLED=true`, mas exigir `CORS_ALLOW_ORIGINS` explícito (sem `*`) quando `CORS_ALLOW_CREDENTIALS=true`.

2. **Guardrails de configuração**
   - Garantir que `AllowCredentials=true` nunca seja combinado com `AllowOrigins=["*"]`.
   - Comportamento esperado: se houver combinação inválida, o middleware **não** aplica headers CORS (e opcionalmente loga/warn).

3. **Documentação**
   - Documentar guardrails e exemplos em `app/infra/http/types/cors.go` (docstring) e/ou em um README do infra (se existir um lugar padrão).

## Critérios de aceite
- Com config padrão (sem env vars): **nenhuma resposta** inclui `Access-Control-Allow-Origin`.
- Com `CORS_ENABLED=true` e origens explícitas: somente origens permitidas recebem `Access-Control-Allow-Origin`.
- Preflight `OPTIONS` continua respondendo `204`.
- `go test ./...` passa.

## Fora de escopo
- Introduzir dependências novas (bibliotecas de CORS).
- Refatorar o design de `contracts.Context`/`contracts.Router`.

## Sugestão de testes
- Adicionar/expandir testes unitários para `BuildCORSHeaders` em `app/infra/http/ports/cors.go` cobrindo:
  - Origin vazio
  - AllowOrigins=["*"] com/sem credentials
  - Vary: Origin
  - MaxAge
