# SLACK-CLONE-BACKEND

## For development

This is project use MakeFile

### create `config-dev.yml`

```yaml
app:
  host:
  port:
postgres:
  user:
  password:
  dbname:
  port:
redis:
  port:
  password:
jwt:
  secret:
api:
  key:
    public:
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