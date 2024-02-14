# SLACK-CLONE-BACKEND

[![run time: Makefile](https://img.shields.io/badge/Run_Time-Makefile-e63946.svg?style=flat-square)](https://github.com/prettier/prettier)

## For development

### Create `config.yml`

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
    public: <public-key>
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

```sh
make run
```

Serve on `http://localhost:8080`
