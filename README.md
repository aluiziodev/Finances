# Finances

API para registrar e consultar faturas a partir de arquivos CSV de extrato.

## Visão geral

O projeto recebe um arquivo CSV contendo registros financeiros, valida a estrutura da fatura, calcula o total automaticamente e salva os dados em um banco de dados relacional. A aplicação também permite listar todas as faturas e consultar uma específica por ID.

## Funcionalidades

- Recebe CSV no formato de extrato (ex.: Nubank)
- Calcula automaticamente o total da fatura
- Persiste faturas e transações no PostgreSQL
- Endpoints para listagem e consulta por ID

## Tecnologias

- Go
- PostgreSQL
- gocsv
- UUID
- net/http
- godotenv

## Estrutura do projeto

```text
. 
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── categorizer/
│   ├── connection/
│   ├── controllers/
│   ├── models/
│   ├── parser/
│   ├── repository/
│   ├── response/
│   ├── router/
│   └── service/
├── sql/
│   └── database.sql
├── go.mod
├── go.sum
├── README.md
└── .env  # Criar o .env na raiz do projeto
```

## Requisitos

- Go 1.26+
- PostgreSQL
- Arquivo `.env` na raiz com as variáveis abaixo

## Variáveis de ambiente

Crie um arquivo `.env` na raiz do projeto com o conteúdo abaixo:

```env
DATABASE_URL= # URL do seu banco de dados PostgreSQL
API_PORT= # porta em que a API será executada
```

## Banco de dados

O script de criação das tabelas está em `sql/database.sql`. 


Execute o `sql/database.sql` no seu banco antes de rodar a aplicação.

## Endpoints

| Método | Rota | Descrição |
| --- | --- | --- |
| `POST` | `/fatura` | Cria uma nova fatura a partir de um CSV e dos dados da fatura (multipart/form-data) |
| `GET` | `/fatura` | Lista todas as faturas cadastradas |
| `GET` | `/fatura/{id}` | Busca uma fatura específica pelo ID |
| `DELETE` | `/fatura/{id}` | Deleta a fatura (remove lançamentos relacionados) |
| `GET` | `/fatura/{id}/parcelado` | Retorna apenas lançamentos com `method = "parcelado"` |
| `GET` | `/fatura/{id}/fixo` | Retorna apenas lançamentos com `method = "fixo"` |
| `GET` | `/fatura/{id}/category` | Retorna lançamentos filtrados por categoria |

Exemplo de `curl` para criar uma fatura (note o campo `data` com JSON e o arquivo CSV em `template/`):

```bash
curl -X POST http://localhost:8080/fatura \
  -F 'data={"description":"Fatura de julho","status":"pending"}' \
  -F 'csv=@template/Nubank_2026-08-14.csv'
```

Exemplo de `curl` para filtrar por categoria:

```bash
curl "http://localhost:8080/fatura/123/category?category=varejo"
```

### Payload da criação

O campo `data` deve ser um JSON com `description` e `status`.

Exemplo de `data`:

```json
{
  "description": "Fatura de julho",
  "status": "pending"
}
```

Resposta esperada:

```json
{
  "id": "uuid",
  "description": "Fatura de julho",
  "status": "pending",
  "total": 123.45,
  "bills": [
    {
      "id": "uuid",
      "date": "2026-08-06",
      "title": "titulo da conta",
      "amount": 8.5
    }
  ]
}
```

## Formato do CSV

O CSV aceita valores no seguinte formato:

```csv
date,title,amount
2026-08-06,conta1,"18,33"
2026-08-06,conta2,"8,50"
2026-07-14,conta3,"461,69"
```

Formato adotado pelo Nubank, por exemplo.

## Observações importantes

- A API usa `multipart/form-data` para receber o CSV e o JSON em uma única requisição
- O status válido é `pending` ou `paid`
- O projeto foi pensado para controle financeiro de faturas importadas de extratos

## Como executar

1. Configure um banco PostgreSQL
2. Crie o arquivo `.env`
3. Execute o script SQL para criar as tabelas
4. Inicie a aplicação:

```bash
go run ./cmd/api
```

A API ficará disponível em:

```text
http://localhost:8080
```

## Testes

O repositório contém testes unitários. Para executá-los:

```bash
go test -v ./...
```

Para executar testes de um pacote específico:

```bash
go test ./internal/parser -v
```

## Configuração CORS 

O projeto já inclui um middleware simples de CORS no arquivo [internal/router/router.go](internal/router/router.go). Atualmente ele permite todas as origens e métodos, o que facilita o desenvolvimento front-end:

```go
w.Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
```

Se você quiser restringir ao front-end em produção, altere `"*"` para a origem do seu front-end e, se necessário, habilite credenciais:

```go
w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
w.Header().Set("Access-Control-Allow-Credentials", "true")
```

Lembre-se de ajustar os cabeçalhos permitidos em `Access-Control-Allow-Headers` caso o front-end envie cabeçalhos personalizados.

## Próximas features

- Incluir novas funcionalidades financeiras
- Criar histórico de alterações
- Incluir autenticação
- Suportar exportação em outros formatos
