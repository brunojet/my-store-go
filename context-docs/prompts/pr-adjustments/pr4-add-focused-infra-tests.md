# PR4 — Add focused unit tests for infra HTTP and observability helpers

## Objetivo
Adicionar testes unitários pequenos e rápidos para partes críticas do módulo `app/infra`, cobrindo regras de borda e prevenindo regressões (especialmente em CORS, adaptação de paths e policy de route para observabilidade).

## Contexto
Há boa cobertura de params/persistence, mas faltam testes diretos para helpers do `infra/http` e `infra/observability`.

## Escopo
- Criar testes de funções puras/helpers (sem servidor real, sem banco real).
- Manter suíte rápida (`go test ./...`).

## Mudanças propostas
1. **Testar CORS**
   - Alvo: `app/infra/http/ports.BuildCORSHeaders`
   - Casos mínimos:
     - Origin vazio => ok=false
     - AllowOrigins=["*"] e AllowCredentials=false => `AllowOrigin="*"`, `VaryOrigin=false`
     - AllowOrigins explícitas => `AllowOrigin=origin`, `VaryOrigin=true`
     - MaxAge negativo/zero => coerção adequada

2. **Testar tradução de path (Gin -> Chi)**
   - Alvo: `app/infra/http/adapters/nethttp.translatePath`
   - Casos:
     - "/apps/:id" -> "/apps/{id}"
     - path vazio
     - múltiplos params

3. **Testar policy de route de baixa cardinalidade (se existir do PR2)**
   - Se um helper for criado em `app/infra/observability/ports`, criar testes para:
     - template/pattern presente
     - template/pattern ausente => retorna constante (ex.: "unmatched")

## Critérios de aceite
- Novos testes cobrem casos de borda relevantes.
- `go test ./...` passa.

## Fora de escopo
- Testes de integração com frameworks (Gin/Chi) ou bancos.
- Testar handlers de módulos fora de `app/infra`.
