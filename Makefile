.PHONY: generate migrate-create dev tests gen

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

generate:
	go generate ./...

gen: generate

migrate-create: generate
	atlas migrate diff $(filter-out $@,$(MAKECMDGOALS)) \
		--dir "file://ent/migrate/migrations" \
		--to "ent://ent/schema" \
		--dev-url "docker://postgres/15/test?search_path=public"

migrate-apply:
	atlas migrate apply \
      --dir "file://ent/migrate/migrations" \
      --url $(DATABASE_URL)

dev:
	go run ./cmd/main.go

tests:
	go test -parallel=20 -covermode atomic -coverprofile=coverage.out ./...

build:
	rm ./build-out || true
	go build -ldflags="-s -w" -o build-out cmd/main.go
	upx -9 -q ./build-out