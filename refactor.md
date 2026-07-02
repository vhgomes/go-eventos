## Code Review — vhgomes-eventos

### ⚠️ Problemas Encontrados

---

#### [CRÍTICO] — Missing `return` após erro de bind no login

**Localização:** `cmd/api/auth.go:46-50`

**Problema:**
```go
func (app *application) login(c *gin.Context) {
    var auth loginRequest
    if err := c.ShouldBindJSON(&auth); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
    }  // ❌ Falta 'return'
    existingUser, err := app.models.Users.GetByEmail(auth.Email) // auth vazio ou inválido
```
Erro de validação é logado no response, mas a função continua executando. Se `auth.Email` for `""` (string vazia), o `GetByEmail` retorna `nil, nil` (não erro), e a linha `bcrypt.CompareHashAndPassword` tentará comparar `""` com o hash — resultando em `nil` pointer dereference ou panic.

**Impacto:** 
- Request com JSON malformado causa **panic em produção**.
- Testes de integração não detectam porque não testam edge case de bind failure.

**Correção:** Adicione `return` imediatamente após o `c.JSON`.

---

#### [CRÍTICO] — `GetUserFromContext` retorna usuário vazio sem erro

**Localização:** `cmd/api/context.go:6-15`

**Problema:**
```go
func (app *application) GetUserFromContext(c *gin.Context) *database.User {
    contextUser, exists := c.Get("user")
    if !exists {
        return &database.User{} // ❌ Retorna Id=0
    }
    // ...
}
```
Se o middleware `AuthMiddleware` não executar (ex: rota pública acidentalmente usando o helper), ou se a conversão falhar silenciosamente, `OwnerId` receberá `0` em `createEvent` (linha 15). O banco com `FOREIGN KEY` rejeita `owner_id=0` se a constraint estiver ativa, mas SQLite permite `0` mesmo com FK se não estiver rigorosa — potencialmente criando órfãos.

**Impacto:** 
- Em PostgreSQL (migração futura), `INSERT` falha com `foreign key violation`.
- Debugging confuso: evento criado com `ownerId=0`, ninguém consegue editá-lo depois.

**Correção:** 
- Mude assinatura para `func (app *application) GetUserFromContext(c *gin.Context) (*database.User, error)`.
- Retorne erro explícito quando user não existir ou não for `*database.User`.
- Em todos os handlers, trate o erro com `http.StatusUnauthorized`.

---

#### [IMPORTANTE] — SQLite em projeto que vai para currículo

**Localização:** Todo o pacote `internal/database`

**Problema:**
SQLite é excelente para desenvolvimento local, mas em entrevistas técnicas sênior, **é um sinal de inexperiência em produção**. Vagas pedem PostgreSQL, MySQL ou DynamoDB.

- Sem suporte a `pgx` (pool de conexões, prepared statements nativos).
- `database/sql` com queries manuais — propenso a erros de digitação e sem type-safety.
- Sem `sqlc` (já está no seu stack! — você deveria estar usando).

**Impacto no currículo:** Recrutadores técnicos veem SQLite e já descartam "produção-ready". Não é aceitável para vaga de pleno/sênior.

**Correção:** 
- Migrar para PostgreSQL (RDS na AWS).
- Integrar `sqlc` + `pgx` — type-safe, geração de código, 2x mais rápido que GORM.

---

#### [IMPORTANTE] — Falta graceful shutdown

**Localização:** `cmd/api/server.go:8-18`

**Problema:**
`server.ListenAndServe()` bloqueia indefinidamente. Se o serviço receber `SIGTERM` (K8s, ECS, deploy), ele mata o processo abruptamente. Conexões abertas, requisições em andamento e goroutines são interrompidas a meio caminho.

**Impacto:** 
- Em produção, escalonamento horizontal ou deploys causam **requests truncados** (5xx para clientes).
- Transações DB não fechadas corretamente.
- O orquestrador (Kubernetes) dá 30s antes de forçar kill — sem graceful, você desperdiça esse tempo.

**Correção:** 
Use `server.Shutdown(ctx)` com `context.WithTimeout` capturando `os.Interrupt` e `syscall.SIGTERM`.

---

#### [IMPORTANTE] — JWT Secret hardcoded com fallback inseguro

**Localização:** `cmd/api/main.go:25`

**Problema:**
`jwtSecret: env.GetEnvString("JWT_SECRET", "some-secret-123456")` — se esquecer de setar no `.env`, produção usa chave pública previsível. Qualquer atacante pode forjar tokens.

**Impacto:** 
- **Vulnerabilidade de segurança grave** em ambiente produtivo.
- Ferramentas de SAST (SonarQube, Trivy) apontam como alta severidade.

**Correção:** 
- Remova fallback. Se `JWT_SECRET` não existir, `log.Fatal` no `main()`.
- Use `JWT_SECRET` com tamanho mínimo (32 chars) — valide no startup.

---

#### [IMPORTANTE] — Erro de autenticação no middleware retorna 401, mas não aborta corretamente em um caso

**Localização:** `cmd/api/middleware.go:38-44`

**Problema:**
```go
claims, ok := token.Claims.(jwt.MapClaims)
if !ok {
    c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
    // ❌ Faltou c.Abort() e return aqui
}
userId := claims["userId"].(float64) // claims é nil, panic
```
Se a conversão falha, `claims` é `nil`, e a linha seguinte causa **panic** (type assertion em nil).

**Impacto:** Idêntico ao primeiro — panic em produção com token malformado.

**Correção:** Adicione `c.Abort(); return` dentro do `if !ok`.

---

#### [MELHORIA] — Inconsistência nas mensagens de erro (maiúsculas/minúsculas)

**Localização:** Vários arquivos, ex: `auth.go:25` usa `"Error:"`, `events.go:11` usa `"error"`.

**Impacto:** Dificulta parsing de erros no frontend e logs estruturados. Padronize para `"error"` minúsculo em toda a API.

---

#### [MELHORIA] — Naming inconsistency em `getAtteendeesForEvent` (typo)

**Localização:** `cmd/api/events.go:160`

**Problema:** `getAtteendeesForEvent` → deveria ser `getAttendeesForEvent`. Typos em nomes públicos comprometem legibilidade e Swagger.

---

### ✅ Pontos Positivos (que demonstram boa base)

- **Context timeout em todas as queries DB** (3s) — evita conexões acumuladas.
- **Uso de `defer rows.Close()`** em todos os scans.
- **Separação de models por domínio** (Users, Events, Attendees) — já estrutura razoável.
- **Swagger integrado** com `swaggo` — documentação básica pronta.
- **Middleware de autenticação** com JWT separado das rotas — boa separação de responsabilidades.
- **Migrations** com `golang-migrate` — já tem pipeline de evolução de schema.

---

## Roadmap para Portfólio CV-Ready

Para transformar este projeto de "CRUD básico de faculdade" para **"Arquitetura Cloud-Native Production-Ready"**, siga a tasklist abaixo.  
Use a tag **[CV-READY]** para commits que agregam valor direto ao seu currículo.

---

### 🏷️ Tags de Modificação

| Tag | Significado |
|-----|-------------|
| `[CV-READY]` | Melhoria que gera **entrevista técnica diferencial** |
| `[INFRA]` | Mudanças em Docker, AWS, Terraform, CI/CD |
| `[ARCH]` | Mudanças arquiteturais (padrões, desacoplamento) |
| `[OBS]` | Observabilidade (logs, métricas, tracing) |
| `[SEC]` | Segurança e hardening |
| `[DB]` | Migrações e otimizações de banco |

---

### FASE 0 — Correções Críticas (1 dia)

| Prioridade | Task | Tag | Detalhe Técnico |
|------------|------|-----|-----------------|
| 🔴 | Fix `login` missing return | `[SEC]` | Adicionar `return` após `c.JSON` no erro de bind. |
| 🔴 | Fix `AuthMiddleware` missing return | `[SEC]` | Adicionar `c.Abort(); return` na conversão de claims. |
| 🔴 | Refactor `GetUserFromContext` para retornar erro | `[ARCH]` | Assinatura `(*User, error)`. Handlers devem tratar erro e retornar 401. |
| 🔴 | Remover fallback de `JWT_SECRET` | `[SEC]` | Fatal no `main()` se variável ausente. Exigir env. |
| 🟠 | Padronizar respostas de erro (`"error"` minúsculo) | `[MELHORIA]` | Refatorar todas as `gin.H{"Error":}` para `gin.H{"error":}`. |
| 🟠 | Corrigir typo `getAtteendeesForEvent` | `[MELHORIA]` | Renomear para `getAttendeesForEvent`. Atualizar rotas. |

---

### FASE 1 — Fundação de Produção (3 dias)

| Prioridade | Task | Tag | Detalhe Técnico |
|------------|------|-----|-----------------|
| 🔵 | Migrar SQLite → PostgreSQL | `[DB]` | Mudar driver para `pgx`, ajustar queries (placeholders `$1` já estão ok para PG). Adicionar `pgxpool` com connection pool configurável. |
| 🔵 | Integrar `sqlc` + `pgx` | `[DB]` | Escrever queries `.sql` no `sqlc/`, gerar models type-safe. Substituir `database/sql` manual. **Isso é o que empresas tier-1 esperam.** |
| 🔵 | Adicionar Graceful Shutdown | `[ARCH]` | Capturar `SIGINT`/`SIGTERM` + `server.Shutdown` com timeout de 30s. Adicionar log de shutdown. |
| 🔵 | Adicionar Health Checks | `[OBS]` | Endpoints `/health` (liveness) e `/ready` (readiness). Checar conectividade com DB. |
| 🟢 | Adicionar logging estruturado (Zap) | `[OBS]` | Substituir `log.Printf` por `zap.L().Info()` com campos (`request_id`, `user_id`, `latency`). Adicionar middleware para log de requisições HTTP. |

---

### FASE 2 — Containerização e Infraestrutura (2 dias)

| Prioridade | Task | Tag | Detalhe Técnico |
|------------|------|-----|-----------------|
| 🟣 | Dockerfile multi-stage | `[INFRA]` | Build estágio 1 (`golang:1.23-alpine`), estágio 2 (`alpine`). Copiar binário estático. Pino versão (`golang:1.23.4-alpine3.20`). |
| 🟣 | Docker Compose com dependências | `[INFRA]` | Serviços: `app`, `postgres:16-alpine`, `redis:7-alpine` (cache), `prometheus` (opcional). Healthchecks no Compose. |
| 🟣 | Configuração via ambiente | `[INFRA]` | Usar `.env` + Viper ou Godotenv. DB_URL, REDIS_URL, JWT_SECRET. |

---

### FASE 3 — Observabilidade e Resiliência (2 dias)

| Prioridade | Task | Tag | Detalhe Técnico |
|------------|------|-----|-----------------|
| 🟠 | Integrar Prometheus metrics | `[OBS]` | Adicionar `promhttp` + métricas customizadas: `http_requests_total` (por rota, status), `http_request_duration_seconds` (histogram), `db_query_duration` (histogram). Endpoint `/metrics`. |
| 🟠 | Adicionar Cache com Redis | `[ARCH]` | Cachear `GetAllEvents` por 60s (TTL) com invalidação no `Create/Update/Delete`. Usar `go-redis/v9`. Demonstra conhecimento de caching distribuído. |
| 🟠 | Circuit Breaker para DB (opcional) | `[ARCH]` | Usar `sony/gobreaker` nas queries DB para fallback rápido se DB estiver lento. |
| 🟢 | Adicionar `request_id` no contexto | `[ARCH]` | Middleware que gera UUID e injeta no `c.Request.Context()`. Logs e spans carregam esse ID para correlação. |

---

### FASE 4 — Cloud-Native com AWS (3 dias) — **DIFERENCIAL COMPETITIVO**

| Prioridade | Task | Tag | Detalhe Técnico |
|------------|------|-----|-----------------|
| 🔥 | Terraform para infra AWS | `[INFRA]` | Módulos: **RDS (PostgreSQL)** com subnet group, security group; **ElastiCache (Redis)**; **ECS Fargate** (ou EC2). Usar `aws-provider` v5. |
| 🔥 | CI/CD com GitHub Actions | `[INFRA]` | Workflow: `lint` (golangci-lint), `test` (unit + integration), `build`, `docker build & push` para ECR, `deploy` via `aws ecs update-service`. |
| 🔥 | Migrar dados para RDS | `[DB]` | Usar `golang-migrate` em container de init ou no entrypoint do app. | 
| 🔥 | SQS para notificações assíncronas | `[ARCH]` | No `addAttendeeToEvent`, publicar mensagem para **SQS** (ex: "user-X joined event-Y"). Criar um worker separado (ou goroutine) que consome da fila e envia e-mail (mock). **Demonstra event-driven architecture** — tópico central do seu perfil. |
| 🔥 | S3 para upload de imagens de eventos | `[INFRA]` | Endpoint `POST /events/:id/image` — upload para S3 bucket, retornar URL pública. Usar presigned URLs para segurança. |
| 🔥 | Adicionar OpenTelemetry (tracing) | `[OBS]` | Integrar com **AWS X-Ray** ou **Jaeger** via `otel` + `otlptrace` (exporter para Collector). Spans para HTTP, DB, Redis, SQS. **Gap identificado no seu perfil — preenche exatamente.** |

---

### FASE 5 — Testes e Qualidade (2 dias)

| Prioridade | Task | Tag | Detalhe Técnico |
|------------|------|-----|-----------------|
| 🟢 | Testes unitários com `testify` | `[CV-READY]` | Cobrir handlers (mocks das models), middleware (JWT parsing), utils. |
| 🟢 | Testes de integração com `testcontainers` | `[CV-READY]` | Subir PostgreSQL e Redis em containers para testes de API completa. |
| 🟢 | Adicionar `golangci-lint` no pre-commit | `[INFRA]` | Configuração com regras: `errcheck`, `govet`, `staticcheck`, `ineffassign`. |
| 🟢 | Adicionar `Dockerfile` de teste local | `[INFRA]` | Comando `docker-compose up` já deve rodar tudo com hot-reload (usar `air` ou `compiled`). |

---

### FASE 6 — Documentação e Apresentação (1 dia)

| Prioridade | Task | Tag | Detalhe Técnico |
|------------|------|-----|-----------------|
| 📘 | Atualizar README com arquitetura | `[CV-READY]` | Diagrama (Mermaid) mostrando: API Gateway → ECS Fargate → RDS/Redis/SQS. Incluir badges (CI, coverage, version). |
| 📘 | Escrever `ARCHITECTURE.md` | `[CV-READY]` | Decisões técnicas: *Por que sqlc? Por que SQS?* Justificativas para entrevistas. |
| 📘 | Swagger enriquecido com exemplos | `[CV-READY]` | Adicionar `@example` nas annotations para cada endpoint. |
| 📘 | Gravar demo (vídeo curto) | `[CV-READY]` | 3-min demo: criar evento, adicionar participante, ver no Grafana dashboards (métricas) e X-Ray (traces). |

---

## Ordem Recomendada de Implementação

1. **Fase 0** (correções críticas) — *prioridade máxima, 1 dia.*
2. **Fase 1** + **Fase 2** (PostgreSQL, Docker, Graceful) — *base sólida para produção.*
3. **Fase 4** (AWS + SQS) — *maior impacto no currículo.*
4. **Fase 3 + Fase 5** (Observabilidade + Testes) — *amadurecimento.*
5. **Fase 6** — *empacotamento final.*

---

## O que este projeto demonstrará no seu currículo

| Habilidade | Evidência |
|------------|-----------|
| **Go idiomático production-ready** | sqlc + pgx, structured logging, graceful shutdown, context propagation |
| **Arquitetura de microsserviços** | Separação de worker (SQS) e API, event-driven com filas |
| **Cloud AWS** | Terraform + RDS + ElastiCache + ECS + SQS + S3 + X-Ray |
| **Observabilidade** | Prometheus (métricas) + OpenTelemetry (traces) + Zap (logs) |
| **DevOps** | Docker multi-stage + GitHub Actions CI/CD + ECR |
| **Design patterns** | Repository pattern, Dependency Injection, Circuit Breaker, Caching Strategy |

Após finalizar, você terá um **case de entrevista completo** para falar sobre decisões de escalabilidade, resiliência e custo na AWS. Isso coloca você no nível pleno/sênior que empresas europeias e Fortune 500 buscam.
