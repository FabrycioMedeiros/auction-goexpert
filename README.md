# Auction GoExpert

Sistema de leilões em Go com **fechamento automático** via Goroutines.

## Como funciona

Quando um leilão é criado, uma Goroutine é iniciada em background que aguarda a duração configurada e, ao expirar, atualiza o status do leilão para `Completed` (fechado) no banco de dados.

## Variáveis de ambiente

| Variável           | Descrição                                             | Padrão   |
|--------------------|-------------------------------------------------------|----------|
| `AUCTION_DURATION` | Duração do leilão (ex: `30s`, `5m`, `1h`)            | `5m`     |
| `AUCTION_INTERVAL` | Alternativa a `AUCTION_DURATION` (mesma semântica)   | —        |
| `MONGODB_URL`      | URL de conexão com o MongoDB                          | —        |
| `MONGODB_DB`       | Nome do banco de dados                                | —        |

O arquivo `cmd/auction/.env` contém as variáveis padrão usadas pelo Docker Compose.

## Rodando com Docker

```bash
docker-compose up --build
```

A API estará disponível em `http://localhost:8080`.

## Rodando localmente

**Pré-requisitos:** Go 1.21+, MongoDB em execução.

```bash
# Ajuste as variáveis no arquivo .env
cp cmd/auction/.env.example cmd/auction/.env

go run cmd/auction/main.go
```

## Rodando os testes

Os testes de integração utilizam [testcontainers-go](https://github.com/testcontainers/testcontainers-go) e requerem Docker em execução.

```bash
go test ./internal/infra/database/auction/ -v -timeout 120s
```

O teste `TestCreateAuctionAutoClose`:
1. Sobe um container MongoDB via testcontainers
2. Cria um leilão com `AUCTION_DURATION=3s`
3. Aguarda 5 segundos
4. Verifica que o status mudou automaticamente para `Completed`

## Endpoints

| Método | Rota                     | Descrição                      |
|--------|--------------------------|--------------------------------|
| POST   | `/auction`               | Cria um novo leilão            |
| GET    | `/auction`               | Lista leilões (filtros opcionais) |
| GET    | `/auction/:auctionId`    | Busca leilão por ID            |
| POST   | `/bid/:auctionId`        | Cria um lance                  |
| GET    | `/bid/:auctionId`        | Lista lances de um leilão      |
| GET    | `/auction/winner/:auctionId` | Busca o lance vencedor     |
