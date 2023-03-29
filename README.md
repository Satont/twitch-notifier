# Twitch Notifier

[![Coverage Status](https://coveralls.io/repos/github/Satont/twitch-notifier/badge.svg)](https://coveralls.io/github/Satont/twitch-notifier)

Bot for sending twitch streams notifications in telegram.

# Development

Download dependencies

```bash
go mod download
```

### Requirements

- Golang 1.20+

### Testing


```bash
make test
```

### Running

```bash
make dev
```

## Database schemas and migrations

We're using [ent](https://entgo.io/) for database schemas.

### Writing schemas

All schemas located in `./ent/schema` directory.

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
