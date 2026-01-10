# Prompt: Chassis (esqueleto) do projeto Go

Você é um assistente de engenharia de software. Gere o **esqueleto** (chassis) de um projeto em Go para uma “loja de aplicativos de PDV”, sem implementar regras de negócio.

## Contexto do projeto

- O repositório já tem um `README.md` com descrição do domínio, papéis e uma proposta de modularização.
- Repositório: https://github.com/brunojet/my-store-go
- Basepath da aplicação: `app/`.
- Módulos pretendidos:
  - `core/` (reutilizável): models, dtos, repositories, persistence/connectors/migrations.
  - `store-provider/`, `store-office/`, `store-consumer/`: cada um com `handlers/`, `routers/`, `services/`.
- Além disso, o `core` deve expor handlers/utilitários reutilizáveis em `core/handlers/`.
- Observação importante: existe `audit` em `models` e em `dtos` porque **DTO (API) nem sempre é igual ao modelo de dados**.

## Tarefa

1. Criar a árvore de diretórios seguindo o layout do README.
2. Adicionar arquivos `.gitkeep` em diretórios que ficariam vazios, para manter a estrutura no Git.
3. Criar arquivos Go **mínimos** (stubs) apenas quando fizer sentido para fixar a convenção do pacote:
   - `package` correto
   - tipos e interfaces vazios ou mínimos
   - sem lógica de negócio
4. Inicializar o módulo Go (`go.mod`) dentro de `app/`.
5. Definir uma convenção simples de roteamento HTTP e handlers (somente estrutura):
   - um arquivo de roteamento por módulo (`store-*/routers`)
   - handlers chamando services, services chamando repositories (interfaces)
6. Manter o repositório limpo e fácil de evoluir.

## Restrições

- Não implementar banco real, queries, endpoints completos, autenticação, etc.
- Não adicionar dependências externas sem necessidade.
- Preferir apenas `net/http` por enquanto.
- Código deve compilar (stubs coerentes), mesmo que não faça nada.

## Variáveis (preencher antes de executar)

- `{{MODULE_PATH}}`: ex.: `github.com/brunojet/my-store-go`
- `{{GO_VERSION}}`: ex.: `1.25.5` (ou a versão escolhida)

## Entregáveis

- Pastas criadas em `app/` conforme o README.
- `.gitkeep` em pastas vazias.
- `app/go.mod` com `module {{MODULE_PATH}}/app` e `go {{GO_VERSION}}`.
- Arquivos Go mínimos, por exemplo:
  - `app/core/handlers/base_handler.go`
  - `app/core/handlers/crud_handler.go`
  - `app/core/models/entity_interface.go`
  - `app/core/models/base_audit.go`
  - `app/core/dtos/dto_interface.go`
  - `app/core/dtos/base_audit.go`
  - `app/core/repositories/base_repository.go`
  - `app/store-provider/routers/routes.go`
  - `app/store-provider/handlers/health.go` (ou similar)
  - `app/store-provider/services/service.go`
  - (mesma ideia para office/consumer)
- Um `app/cmd/` com pelo menos um binário (`server`) que sobe um `http.Server` e registra rotas (mesmo que sejam só rotas de health).

## Checklist de validação

- `go test ./...` (ou pelo menos `go test` sem falhas) dentro de `app/`.
- Imports organizados (`gofmt`).
- Estrutura bate com o README.

## Saída esperada

- Liste os arquivos/pastas criados.
- Mostre os trechos principais do `go.mod` e do `main.go`.
- Aponte onde os próximos passos devem entrar (ex.: definição de contratos HTTP, DTOs reais, repos, migrations).
