Memo
==

A CQRS+ES study project written in Go to manage reminders

### Setup
```bash
go mod tidy
```

### Test

`cp .env .env.local` and set your DSN

`cp docker/docker-compose.override.dist.yml docker/docker-compose.override.yml` and set your DB port

``` bash
make upd
make prepare-db
go test -v ./...
```