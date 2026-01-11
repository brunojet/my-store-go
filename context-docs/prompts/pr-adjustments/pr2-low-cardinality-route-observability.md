# PR2 — Prevent high-cardinality route labels in observability middleware

## Objetivo
Evitar **alta cardinalidade** em métricas/traces/logs gerados pelos middlewares de observabilidade (Gin e net/http) quando a rota não tem template/pattern disponível, mantendo utilidade operacional e consistência entre drivers.

## Contexto
Os middlewares de telemetria e access log estão em:
- Gin: `app/infra/observability/adapters/ginmw/*`
- net/http: `app/infra/observability/adapters/httpmw/*`

Hoje, quando não há template/pattern (ex.: 404, rotas não registradas), há fallback para `URL.Path`. Isso pode explodir cardinalidade (ids no path, etc.) em métricas e traces.

## Escopo
- Definir e aplicar uma política única de seleção de `route` (baixa cardinalidade) em telemetry e, se aplicável, no access log.
- Manter compatibilidade com o design atual.

## Mudanças propostas (cirúrgicas)
1. **Política de route de baixa cardinalidade**
   - Quando houver template/pattern: usar template/pattern (baixo cardinalidade).
   - Quando não houver: usar um valor fixo, ex.: `"unmatched"` ou `"unknown"`.

2. **Aplicar a policy nos middlewares**
   - Gin: `app/infra/observability/adapters/ginmw/telemetry.go`
   - net/http: `app/infra/observability/adapters/httpmw/telemetry.go`
   - (Opcional) Access log:
     - Gin: `app/infra/observability/adapters/ginmw/access_log.go`
     - net/http: `app/infra/observability/adapters/httpmw/access_log.go`

3. **Centralização (opcional, recomendada)**
   - Criar helper pequeno em `app/infra/observability/ports` para escolher route e/ou para formatar campos de observabilidade, para evitar drift entre Gin e net/http.

## Critérios de aceite
- Para rotas não registradas (sem template/pattern): métricas/traces usam `route="unmatched"` (ou equivalente) e **não** incluem o path real como label.
- Para rotas registradas: continuam usando template/pattern.
- `go test ./...` passa.

## Fora de escopo
- Trocar implementação de telemetry para OpenTelemetry/exporters.
- Alterar contratos do core (`app/core/telemetry`).

## Sugestão de testes
- Teste unitário do helper de route (se criado) cobrindo:
  - template/pattern presente
  - template/pattern vazio
  - caminhos dinâmicos simulados
