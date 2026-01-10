# Prompts (foundation)

Este diretório contém prompts base para montar o **chassis/esqueleto** do projeto em Go.

## Como usar

- Cada arquivo `.md` é um prompt autocontido.
- Copie o conteúdo do prompt para o seu LLM (ou para a ferramenta de prompts que você usar) e preencha as variáveis.

## Convenções

- Os prompts aqui **não** implementam regras de negócio.
- Objetivo: criar estrutura de pastas, arquivos mínimos (stubs), organização de pacotes, e documentação inicial.
- Manter consistência com o layout descrito no `README.md` do repositório.

## Arquivos

- `01-chassis-go-project.md`: prompt principal para gerar o esqueleto do projeto.
- `02-apps-api-end-to-end-gin-gorm-sqlite-memory.md`: prompt para implementar CRUD de aplicativos end-to-end (HTTP → service → repo → SQLite in-memory via GORM) com adapter desacoplado.
