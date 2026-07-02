# 📅 Go Eventos — Event Management API

API RESTful para gerenciamento de eventos e participantes, desenvolvida em **Go** com o framework **Gin**, seguindo os princípios de **Clean Architecture** (domain → service → repository → handler). Permite registro e autenticação de usuários via **JWT**, além de criação, atualização, exclusão e inscrição de participantes em eventos.

---

## 🚀 Tecnologias Utilizadas

- [Go](https://golang.org/) `1.23.4`
- [Gin Web Framework](https://github.com/gin-gonic/gin) — roteamento HTTP
- [PostgreSQL](https://www.postgresql.org/) — banco de dados relacional
- [golang-migrate](https://github.com/golang-migrate/migrate) — migrações de banco de dados
- [golang-jwt/jwt](https://github.com/golang-jwt/jwt) — autenticação via token JWT
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) — hash de senhas
- [Swaggo](https://github.com/swaggo/swag) + [gin-swagger](https://github.com/swaggo/gin-swagger) — documentação interativa da API
- [Zap](https://github.com/uber-go/zap) — logging estruturado
- [godotenv](https://github.com/joho/godotenv) — carregamento de variáveis de ambiente

---

## 🏗️ Arquitetura

O projeto segue uma separação em camadas inspirada em Clean Architecture:

```
cmd/api            → ponto de entrada da aplicação (main.go)
internal/domain     → entidades de negócio e regras de validação
internal/service     → regras de negócio e orquestração de casos de uso
internal/repository  → contratos (interfaces) e implementação em PostgreSQL
internal/handlers    → controllers HTTP (Gin)
internal/middleware  → autenticação JWT
internal/pkg         → utilitários (config, errors, logger)
migrations           → scripts SQL versionados (golang-migrate)
docs                 → documentação Swagger gerada automaticamente
```

### Entidades principais

| Entidade   | Descrição                                                        |
|------------|-------------------------------------------------------------------|
| `User`     | Usuário do sistema (email, nome, senha com hash)                  |
| `Event`    | Evento (nome, descrição, data, local, dono/`ownerId`)              |
| `Attendee` | Relação entre um usuário e um evento (inscrição)                  |

Regras de negócio relevantes:
- Apenas o **dono do evento** pode atualizar, excluir ou gerenciar participantes.
- Não é permitido criar dois eventos com o mesmo nome na mesma data.
- Não é permitido inscrever o mesmo participante duas vezes no mesmo evento.
- Senhas são armazenadas com hash `bcrypt`; tokens JWT expiram em 72 horas.

---

## 🔐 Autenticação

A API utiliza autenticação via **Bearer Token (JWT)**. Após login, envie o token no header:

```
Authorization: Bearer <seu_token>
```

Endpoints protegidos usam o middleware `AuthMiddleware`, que valida o token e injeta o usuário autenticado no contexto da requisição.

---

## 📚 Documentação da API (Swagger)

Com a aplicação em execução, a documentação interativa fica disponível em:

```
http://localhost:8080/swagger/index.html
```

---

## 🛣️ Endpoints

### 📂 Públicos

| Método | Rota                                  | Descrição                                   |
|--------|----------------------------------------|----------------------------------------------|
| GET    | `/health`                              | Healthcheck da aplicação                     |
| GET    | `/api/v1/events`                       | Lista todos os eventos                       |
| GET    | `/api/v1/events/:id`                   | Detalhes de um evento específico             |
| GET    | `/api/v1/events/:id/attendees`         | Lista os participantes de um evento          |
| GET    | `/api/v1/attendees/:id/events`         | Lista os eventos em que um usuário participa |
| POST   | `/api/v1/auth/register`                | Registra um novo usuário                     |
| POST   | `/api/v1/auth/login`                   | Autentica um usuário e retorna um token JWT  |

### 🔐 Protegidos (requer autenticação)

| Método | Rota                                            | Descrição                          |
|--------|--------------------------------------------------|-------------------------------------|
| POST   | `/api/v1/events`                                 | Cria um novo evento                 |
| PUT    | `/api/v1/events/:id`                             | Atualiza um evento existente        |
| DELETE | `/api/v1/events/:id`                             | Remove um evento                    |
| POST   | `/api/v1/events/:id/attendees/:userId`           | Adiciona um participante ao evento  |
| DELETE | `/api/v1/events/:id/attendees/:userId`           | Remove um participante do evento    |

---

## ⚙️ Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto (carregado automaticamente via `godotenv`):

| Variável           | Obrigatória | Padrão | Descrição                                                       |
|--------------------|:-----------:|:------:|-------------------------------------------------------------------|
| `PORT`             | Não         | `8080` | Porta em que o servidor HTTP será iniciado                        |
| `DATABASE_URL`     | **Sim**     | —      | String de conexão PostgreSQL, ex: `postgres://user:pass@localhost:5432/db?sslmode=disable` |
| `JWT_SECRET`       | **Sim**     | —      | Chave secreta para assinatura do JWT (mínimo 32 caracteres)       |
| `DB_MAX_OPEN_CONNS`| Não         | `25`   | Número máximo de conexões abertas com o banco                     |
| `DB_MAX_IDLE_CONNS`| Não         | `5`    | Número máximo de conexões ociosas com o banco                     |
| `LOG_LEVEL`        | Não         | `info` | Nível de log (`debug`, `info`, `warn`, `error`)                   |
| `LOG_OUTPUT`       | Não         | `logs.json` | Arquivo de saída dos logs                                    |

---

## 🛠️ Como Executar Localmente

### Pré-requisitos

- [Go 1.23+](https://golang.org/dl/)
- Uma instância [PostgreSQL](https://www.postgresql.org/) em execução
- [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) (opcional, para rodar as migrations)
- [swag CLI](https://github.com/swaggo/swag) (opcional, para regerar a documentação Swagger)

### Passo a passo

1. Clone o repositório:

   ```bash
   git clone https://github.com/vhgomes/go-eventos.git
   cd go-eventos
   ```

2. Crie o arquivo `.env` com as variáveis descritas acima:

   ```bash
   cp .env.example .env   # ajuste os valores conforme seu ambiente
   ```

3. Instale as dependências:

   ```bash
   make tidy
   ```

4. Execute as migrations do banco de dados:

   ```bash
   make migrate-up
   ```

5. (Opcional) Regere a documentação Swagger:

   ```bash
   make swag
   ```

6. Execute a aplicação:

   ```bash
   make run
   ```

7. Acesse:

   - API: `http://localhost:8080/api/v1`
   - Swagger UI: `http://localhost:8080/swagger/index.html`
   - Healthcheck: `http://localhost:8080/health`

---

## 🧰 Comandos disponíveis (Makefile)

O projeto conta com um `Makefile` para facilitar as tarefas do dia a dia. Rode `make help` para listar todos os comandos.

### Aplicação

| Comando           | Descrição                                       |
|--------------------|---------------------------------------------------|
| `make run`          | Executa a aplicação localmente                    |
| `make build`        | Compila o binário em `bin/go-eventos`             |
| `make test`         | Executa os testes                                 |
| `make test-cover`   | Executa os testes e gera relatório de cobertura (`coverage.html`) |

### Qualidade de código

| Comando        | Descrição                                  |
|-----------------|----------------------------------------------|
| `make tidy`      | Organiza e baixa as dependências do `go.mod` |
| `make fmt`       | Formata o código (`go fmt`)                  |
| `make vet`       | Analisa o código em busca de erros comuns    |
| `make lint`      | Executa `fmt` + `vet`                        |

### Banco de dados (migrations)

| Comando                                   | Descrição                                             |
|---------------------------------------------|---------------------------------------------------------|
| `make migrate-up`                            | Aplica todas as migrations pendentes                     |
| `make migrate-down`                          | Reverte a última migration aplicada                      |
| `make migrate-version`                       | Mostra a versão atual das migrations                     |
| `make migrate-force VERSION=1`               | Força a versão das migrations (uso em caso de *dirty state*) |
| `make migrate-create NAME=nome_da_migration` | Cria um novo par de arquivos de migration (`up`/`down`)  |
| `make db-reset`                              | Reverte todas as migrations e aplica novamente do zero   |

### Documentação e Docker

| Comando              | Descrição                                      |
|-----------------------|---------------------------------------------------|
| `make swag`            | Regenera a documentação Swagger em `docs/`         |
| `make docker-build`    | Builda a imagem Docker da aplicação                |
| `make docker-run`      | Sobe a aplicação via Docker usando o `.env`        |

### Utilitários

| Comando        | Descrição                                          |
|-----------------|-------------------------------------------------------|
| `make help`      | Lista todos os comandos disponíveis                   |
| `make clean`     | Remove binário, relatórios de cobertura e logs gerados |

> ℹ️ Os comandos `migrate-*` e `docker-run` utilizam automaticamente as variáveis definidas no arquivo `.env` (via `DATABASE_URL`).

---

## 📁 Estrutura de Diretórios

```
.
├── cmd/api/main.go
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal/
│   ├── domain/
│   │   ├── attendee.go
│   │   ├── event.go
│   │   └── user.go
│   ├── handlers/
│   │   ├── attendee_handler.go
│   │   ├── context.go
│   │   ├── error.go
│   │   ├── events_handler.go
│   │   └── user_handler.go
│   ├── middleware/
│   │   └── middleware.go
│   ├── pkg/
│   │   ├── config/config.go
│   │   ├── errors/errors.go
│   │   └── logger/logger.go
│   ├── repository/
│   │   ├── postgres/
│   │   │   ├── attendees_repo.go
│   │   │   ├── events_repo.go
│   │   │   └── users_repo.go
│   │   └── interface.go
│   └── service/
│       ├── attendees_service.go
│       ├── events_service.go
│       └── user_service.go
├── migrations/
│   ├── 000001_create_users_table.up.sql / .down.sql
│   ├── 000002_create_events_table.up.sql / .down.sql
│   └── 000003_create_attendees_table.up.sql / .down.sql
└── go.mod
```

---

## 🧪 Tratamento de Erros

A API centraliza o tratamento de erros de domínio em `internal/pkg/errors`, mapeando-os para códigos HTTP apropriados:

| Erro de domínio     | Código HTTP |
|----------------------|:-----------:|
| `ErrNotFound`         | 404          |
| `ErrUnauthorized`      | 401          |
| `ErrForbidden`         | 403          |
| `ErrConflict`          | 409          |
| `ErrInvalidData`       | 400          |
| *(demais erros)*      | 500          |

---

## ✍️ Autor

Desenvolvido por [Victor Hugo Gomes](https://github.com/vhgomes)
