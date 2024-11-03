# Imersão Fullcycle 19 - Plataforma de streaming de vídeos

## Descrição

Repositório do Django (admin dos vídeos)

## Requerimentos

Na aula 01, precisamos instalar o Python, agora temos um container Python no Docker, então não é mais necessário instalar o Python na máquina.

## Rodar a aplicação

Levante os containers do PostgreSQL e PGAdmin:

```bash
docker-compose up -d
```

Entre no container do Django:

```bash
docker-compose exec django bash
```

Instale as dependências:

```bash
pipenv install
```

A partir daqui, precisamos sempre rodar os comandos dentro do ambiente virtual do Pipenv:

```bash
pipenv shell
```

Rode as migrações do Django:

```bash
python manage.py migrate
```

Crie um superusuário:

```bash
python manage.py createsuperuser
```

Rode o servidor:

```bash
python manage.py runserver
```

Acesse o admin em [http://localhost:8000/admin]().


## Dados de testes

A aplicação já possui dados de testes, rode o comando abaixo para carregá-los:

```bash
python manage.py flush && python manage.py loaddata initial_data.json
```

O comando `flush` limpa o banco de dados e o `loaddata` carrega os dados de testes.

## Configurar /etc/hosts

O RabbitMQ está sendo executado no `docker-compose.yaml` da aplicação Golang, assim como o Django está em outro `docker-compose.yaml`, os containers estão em redes diferentes.
Usaremos a estratégia do `host.docker.internal` para comunicação entre os containers.

Para isto é necessário configurar um endereços que todos os containers Docker consigam acessar.

Acrescente no seu /etc/hosts (para Windows o caminho é C:\Windows\system32\drivers\etc\hosts):
```
127.0.0.1 host.docker.internal
```
Em todos os sistemas operacionais é necessário abrir o programa para editar o *hosts* como Administrator da máquina ou root.

Obs.: Se estiver usando o Docker Desktop, pode ser que o `host.docker.internal` já esteja configurado, então remova a linha do arquivo hosts e acrescente a recomendada acima.


## Consumer do RabbitMQ

Para rodar o consumer do RabbitMQ, entre no container do Django:

```bash
docker-compose exec django bash
```

Rode os consumers:

```bash
python manage.py consumer_upload_chunks_to_external_storage
python manage.py consumer_register_processed_video_path
```
