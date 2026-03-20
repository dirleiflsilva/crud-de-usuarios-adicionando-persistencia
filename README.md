# CRUD de Usuários

API simples em Go com persistência em SQLite.

## Requisitos

- Go 1.21 ou superior

## Como executar

```bash
go run main.go
```

A API sobe em `http://localhost:8080`.

Na primeira execução, o arquivo `users.db` é criado automaticamente.

## Endpoints

- `POST /users`
- `GET /users`
- `GET /user?id={id}`
- `PUT /user?id={id}`
- `DELETE /user?id={id}`

## Exemplo de JSON

```json
{
  "first_name": "Joao",
  "last_name": "Silva",
  "biography": "Desenvolvedor Go"
}
```

## Exemplo com curl

Criar usuário:

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"first_name":"Joao","last_name":"Silva","biography":"Desenvolvedor Go"}'
```

Listar usuários:

```bash
curl http://localhost:8080/users
```

Buscar um usuário:

```bash
curl "http://localhost:8080/user?id=SEU_ID"
```

Atualizar um usuário:

```bash
curl -X PUT "http://localhost:8080/user?id=SEU_ID" \
  -H "Content-Type: application/json" \
  -d '{"first_name":"Joao","last_name":"Santos","biography":"Backend Go"}'
```

Remover um usuário:

```bash
curl -X DELETE "http://localhost:8080/user?id=SEU_ID"
```
