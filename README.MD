# Song library

## Start project
```bash
make clear && make pack && make run
```
or
```bash
docker network create song_network
docker compose build
docker compose up -d
```

## Make swagger docs
```bash
go-swagger3 --module-path . --main-file-path ./cmd/server/server.go --output ./docs/swagger.json --schema-without-pkg
```
Can check url `http://localhost:8001/api/v1/docs/swagger/`

## Migrate up
```bash
migrate -database "dbUriConnString" -path ./migrations up 
```
or build and up project on docker