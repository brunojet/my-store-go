# PR Adjustments — Oportunidades no módulo `infra`

Data: 2026-01-11

## Contexto
Este documento resume oportunidades de melhoria identificadas no módulo `app/infra` (config/http/observability/persistence/server), com foco em arquitetura, design patterns, acoplamento, consistência e manutenção.

O objetivo é gerar um PR com mudanças pequenas, seguras e de alto impacto, mantendo a arquitetura Ports & Adapters/Clean já presente.

## Achados (o que está bom)
- Boa separação de contratos e adapters (Ports & Adapters): `contracts.Router`, `contracts.Context`, `config/contracts.Source`.
- Seleção de drivers (Strategy/Factory) para HTTP e DB.
- Composition root claro em `app/cmd/server/*`.
- Middleware/observability configurável e reaproveitável.

## Oportunidades priorizadas

### P0 — Remover duplicidade de CRUD handler (reduz bifurcação de padrão)
**Sintoma**
- Existe um CRUD handler framework-agnostic em `app/infra/http/ports/crud_handler.go` (baseado em `contracts.Context`).
- Também existe um CRUD handler dependente de Gin em `app/infra/http/adapters/gin/crud_handler.go`.

**Risco**
- Bifurcação de padrão (times começam a copiar o “legado”), manutenção duplicada, inconsistência de mapeamento de erro, e regressões por divergência.

**Ação recomendada**
- Identificar se o handler legado Gin é usado.
  - Se não for usado: remover arquivo e ajustar imports/referências.
  - Se for usado: migrar call sites para o handler agnóstico e então remover.

**Critérios de aceite**
- Apenas um caminho oficial para CRUD handlers.
- Rotas existentes continuam funcionando (Gin e Chi).

---

### P0 — Tornar `DatabaseManager` thread-safe (consistência com outros Managers)
**Sintoma**
- `HTTPManager` e `ObservabilityManager` possuem `sync.Mutex` e cache interno.
- `DatabaseManager` mantém `m.db` e `Connector`, mas não tem sincronização.

**Risco**
- Em cenários concorrentes (startup paralelo, health checks, workers), pode haver double-open, data race, ou close durante uso.

**Ação recomendada**
- Adicionar `sync.Mutex` em `DatabaseManager` e proteger `Open/OpenAndMigrate/Close`.
- Manter API pública igual.

**Critérios de aceite**
- `DatabaseManager` fica alinhado ao padrão de lifecycle dos outros managers.

---

### P1 — Alinhar validação de driver entre `HTTPManager` e `Observability`
**Sintoma**
- `HTTPManager` falha se `Driver` não é suportado.
- `BuildMiddlewares` em `observability` assume Gin como default quando driver é desconhecido.

**Risco**
- Comportamento divergente: configuração inválida pode “parecer funcionar” parcialmente, mascarando erro.

**Ação recomendada**
- Validar `Driver.IsSupported()` antes de construir middlewares, ou fazer `BuildMiddlewares` retornar erro em driver inválido.
- Alternativa: normalizar driver no layer de config (composition root) e validar explicitamente.

**Critérios de aceite**
- Configuração inválida falha de forma previsível e consistente.

---

### P1 — Clarificar flags de observability (`OBS_ACCESS_LOG` vs `OBS_TELEMETRY`)
**Sintoma**
- A config possui alias `OBS_ACCESS_LOG` para controlar `Telemetry`, mas há também middleware de access log separado.

**Risco**
- Semântica confusa de flags; dificulta operação e troubleshooting.

**Ação recomendada**
- Opção A (mais clara): adicionar flag dedicada `OBS_ACCESS_LOG` para habilitar `AccessLog` e manter `OBS_TELEMETRY` para tracing/metrics.
- Opção B (mínima): manter alias apenas por compatibilidade, mas documentar claramente e/ou planejar depreciação.

**Critérios de aceite**
- Flags refletem o que realmente habilitam.

---

### P2 — Endurecer defaults de CORS para produção (documentação/guardrails)
**Sintoma**
- Defaults permitem `AllowOrigins=["*"]` com `CORS_ENABLED=true`.

**Risco**
- Em produção pode ser permissivo demais (dependendo do contexto do produto).

**Ação recomendada**
- Documentar claramente defaults (dev) e recomendações de produção.
- Opcional: condicionar defaults por `ENV` (ex.: dev vs prod) no composition root.

**Critérios de aceite**
- Recomendações operacionais claras para produção.

## Escopo recomendado do PR (sugestão)
- PR 1 (pequeno e seguro):
  - Remover/migrar CRUD handler legado do Gin.
  - Adicionar `Mutex` no `DatabaseManager`.
- PR 2 (comportamental/operacional):
  - Unificar validação de driver e clarificar flags de observability.
  - Revisar documentação/defaults de CORS.

## Checklist de execução
- [ ] Confirmar uso do arquivo `app/infra/http/adapters/gin/crud_handler.go` (grep/references).
- [ ] Implementar mudanças mantendo APIs públicas estáveis.
- [ ] Rodar testes existentes (`*_test.go`) e ajustar apenas o necessário.
- [ ] Atualizar README/documentação onde houver configuração por env.

## Observações
- Evitar “refactors grandes”; priorizar mudanças cirúrgicas com cobertura de teste e diffs pequenos.
