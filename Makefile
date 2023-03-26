.PHONY: generate migrate-create dev tests gen

generate:
	go generate ./...

gen: generate

migrate-create: generate
	atlas migrate diff $(CLI_ARGS) \
		--dir "file://ent/migrate/migrations" \
		--to "ent://ent/schema" \
		--dev-url "docker://postgres/15/test?search_path=public"

dev:
	go run ./cmd/main.go

tests:
	go test -v -parallel=20 -covermode atomic -coverprofile=coverage.out ./...