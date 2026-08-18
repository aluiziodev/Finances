# Finances

API para registrar e consultar faturas financeiras a partir de arquivos CSV de extratos.

## Visão geral

O projeto recebe um arquivo CSV contendo registros financeiros, valida a estrutura da fatura, calcula o total automaticamente e salva os dados em um banco de dados relacional. A aplicação também permite listar todas as faturas e consultar uma específica por ID.

## Funcionalidades

- Upload de CSV com transações da fatura
- Cálculo automático do valor total
- Persistência de faturas e transações
- Consulta de todas as faturas e de uma fatura específica por ID

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
- Arquivo `.env` com as variáveis de ambiente

## Variáveis de ambiente

Crie um arquivo `.env` na raiz do projeto com o conteúdo abaixo:

```env
DATABASE_URL= # URL do seu banco de dados PostgreSQL
API_PORT= # porta em que a API será executada
```

## Banco de dados

O script de criação das tabelas está em `sql/database.sql`.

```sql
CREATE TABLE fatura (
    id TEXT PRIMARY KEY,
    description TEXT NOT NULL,
    status TEXT NOT NULL,
    total NUMERIC(10,2) NOT NULL,

    CONSTRAINT status_check
        CHECK (status IN ('pending', 'paid'))
);

CREATE TABLE bill (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    date TEXT NOT NULL,
    amount NUMERIC(10,2) NOT NULL,
    fatura TEXT NOT NULL,

    CONSTRAINT fk_fatura
        FOREIGN KEY (fatura)
        REFERENCES fatura(id)
        ON DELETE CASCADE
);
```


## Endpoints

| Método | Rota | Descrição | Exemplo |
| --- | --- | --- | --- |
| `POST` | `/fatura` | Cria uma nova fatura a partir de um CSV e dos dados da fatura | `curl -X POST http://localhost:8080/fatura -F 'data={"description":"Fatura de julho","status":"pending"}' -F 'csv=@template/fatura.csv'` |
| `GET` | `/fatura` | Lista todas as faturas cadastradas | `curl http://localhost:8080/fatura` |
| `GET` | `/fatura/{id}` | Busca uma fatura específica pelo ID | `curl http://localhost:8080/fatura/123e4567-e89b-12d3-a456-426614174000` |

### Payload da criação

A requisição de criação deve usar `multipart/form-data` com dois campos:

- `data`: JSON com `description` e `status`
- `csv`: arquivo CSV com as transações

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

## Próximas features

- Implementar testes automatizados
- Criar uma interface de usuário e fazer deploy em um ambiente de produção
- Criar histórico de alterações
- Incluir autenticação
- Suportar exportação em outros formatos
- Incluir novas funcionalidades financeiras
