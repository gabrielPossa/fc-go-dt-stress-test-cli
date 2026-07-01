# fc-go-dt-stress-test-cli

CLI em Go para realizar **testes de carga** em serviços web. Você informa a URL,
o número total de requisições e o nível de concorrência, e ao final é exibido um
relatório com o tempo total, a quantidade de requests e a distribuição dos
códigos de status HTTP.

## Parâmetros

| Flag            | Descrição                                | Obrigatório | Default |
|-----------------|------------------------------------------|-------------|---------|
| `--url`         | URL do serviço a ser testado             | Sim         | —       |
| `--requests`    | Número total de requisições              | Não         | `1`     |
| `--concurrency` | Número de chamadas simultâneas           | Não         | `1`     |

> O total de `--requests` é sempre cumprido de forma exata, independentemente do
> valor de `--concurrency`.

## Executando com Docker (recomendado)

### 1. Build da imagem

```bash
docker build -t stress-test .
```

### 2. Executar o teste

```bash
docker run stress-test --url=http://google.com --requests=1000 --concurrency=10
```

## Executando localmente (sem Docker)

Requer Go 1.26+.

```bash
go run ./cmd/stress-test --url=http://google.com --requests=1000 --concurrency=10
```

Ou compilando o binário:

```bash
go build -o stress-test ./cmd/stress-test
./stress-test --url=http://google.com --requests=1000 --concurrency=10
```

## Exemplo de saída

```text
Total de requests: 1000
Requests com status 200: 976
Falhas (erro de conexão/timeout): 4
Tempo total: 3.21 s
Throughput: 311.53 req/s
Status codes:
200: 976
404: 12
500: 8
```

### Como ler o relatório

- **Total de requests** — quantidade total efetivamente executada.
- **Requests com status 200** — respostas bem-sucedidas.
- **Falhas** — requisições que não obtiveram resposta HTTP (timeout, DNS,
  conexão recusada, etc.). Não possuem código de status.
- **Tempo total** — duração da execução, do primeiro disparo ao último retorno.
- **Throughput** — vazão média (requests por segundo).
- **Status codes** — distribuição dos demais códigos HTTP retornados (404, 500, …).
