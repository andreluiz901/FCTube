# Imersão Fullcycle 19 - Plataforma de streaming de vídeos

## Descrição

Repositório do Golang (conversor de vídeos)

## Rodar a aplicação

Levante os containers:

```bash
docker-compose up -d
```

Através do PGAdmin, crie as tabelas do arquivo `db.sql` no banco de dados `converter`.

Acesse o container o rodando o comando:

```bash
docker compose exec go_app_dev bash
```

Use os comandos `cmd/splitchunks/main.go` e `cmd/videoconverter/main.go` para rodar a aplicação, conforme a aula.