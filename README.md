# SLACK-CLONE-BACKEND

## For development

This is project use MakeFile

### create `config.yml`

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

### Init Migrate Postgres

```sh
make migrate-init
```

### Migrate Postgres

```sh
make migrate
```

### Run server

serve on `http://localhost:8080`

```sh
make run
```
