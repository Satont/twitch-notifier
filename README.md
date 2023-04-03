# Twitch Notifier

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/satont/twitch-notifier)
[![Coverage Status](https://coveralls.io/repos/github/Satont/twitch-notifier/badge.svg)](https://coveralls.io/github/Satont/twitch-notifier)

Bot for sending twitch streams notifications in telegram.

# Development

Download dependencies

```bash
go mod download
```

### Requirements

- Golang `1.19+`

### Generate

After clone/on first setup/on schema change - you should run 

```bash
make generate
```

### Testing

```bash
make tests
```

### Running

```bash
make dev
```

## Database schemas and migrations

### Writing schemas

All schemas located in `./ent/schema` directory, but also we are using internal structures. Internal structures located in `internal/db/db_models`. So you should change both of them.

After changing any schema in `/ent/schema` folder, you should regenerate data via `make generate`

### Migrations

#### Requirements

- [atlasgo cli](https://atlasgo.io/getting-started#installation)
- Docker

### Create

```bash
make migrate-create somecoolname
```

### Apply

```bash
make migrate-apply
```
