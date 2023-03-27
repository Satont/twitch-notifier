ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: generate migrate-create dev tests gen

generate:
	go generate ./...
	echo $(TELEGRAM_BOT_ADMINS)

gen: generate

migrate-create: generate
	atlas migrate diff $(CLI_ARGS) \
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