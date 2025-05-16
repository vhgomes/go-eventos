# 📅 Event Management API

Uma API RESTful para gerenciamento de eventos e participantes, desenvolvida em Go utilizando o framework Gin. A API permite o registro e autenticação de usuários, bem como a criação, atualização e exclusão de eventos, além da adição de participantes a esses eventos.

---

## 🚀 Tecnologias Utilizadas

- [Go](https://golang.org/)
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Swagger UI com Swag](https://github.com/swaggo/gin-swagger)

---

## 🔐 Autenticação

Alguns endpoints são protegidos por middleware de autenticação. Para acessá-los, é necessário estar autenticado via token JWT.

---

## 📚 Documentação da API

A documentação interativa via Swagger está disponível em:

```
http://localhost:8080/swagger/index.html
```

---

## 🛣️ Endpoints

### 📂 Públicos

- `GET /api/v1/events`  
  Lista todos os eventos.

- `GET /api/v1/events/:id`  
  Retorna os detalhes de um evento específico.

- `GET /api/v1/events/:id/attendees`  
  Lista os participantes de um evento.

- `GET /api/v1/attendees/:id/events`  
  Lista os eventos que um participante está inscrito.

- `POST /api/v1/auth/register`  
  Registra um novo usuário.

- `POST /api/v1/auth/login`  
  Autentica um usuário e retorna um token JWT.

---

### 🔐 Protegidos (Requer Autenticação)

- `POST /api/v1/events`  
  Cria um novo evento.

- `PUT /api/v1/events/:id`  
  Atualiza os dados de um evento existente.

- `DELETE /api/v1/events/:id`  
  Remove um evento.

- `POST /api/v1/events/:id/attendees/:userId`  
  Adiciona um participante ao evento.

- `DELETE /api/v1/events/:id/attendees/:userId`  
  Remove um participante do evento.

---

## 🛠️ Como Executar Localmente

1. Clone o repositório:

   ```bash
   git clone https://github.com/seu-usuario/seu-repo.git
   cd seu-repo
   ```

2. Instale as dependências:

   ```bash
   go mod tidy
   ```

3. Gere a documentação Swagger (se necessário):

   ```bash
   swag init
   ```

4. Execute a aplicação:

   ```bash
   go run main.go
   ```

## ✍️ Autor

Desenvolvido por [Victor Hugo Gomes](https://github.com/vhgomes)