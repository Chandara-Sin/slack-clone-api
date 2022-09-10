# SLACK-CLONE-BACKEND

## For development

This is project use MakeFile

### create `config-dev.yml`

```yaml
app:
  port:
postgres:
  host:
  user:
  password:
  dbname:
  port:
jwt:
  secret:
```

### Use Docker + Postgres

```sh
make up
```

### Down Docker + Postgres

```sh
make down
```

### Run server

serve on `http://localhost:8080`

```sh
make run
```

### Run format

```sh
make format
```