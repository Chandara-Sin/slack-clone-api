# SLACK-CLONE-BACKEND

## For development

This is project use MakeFile

### create `config.yml`

```yaml
app:
  host: localhost
  port: "8080"
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
  key: x-api-key
    public: <public-key-vlaue>
mail:
  sender: <sender-email>
  key: <sendgrid-key>
```

### Use Docker + Postgres

```sh
make up
```

### Down Docker + Postgres

```sh
make down
```

### Init Migrate Postgres - Start Project

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
