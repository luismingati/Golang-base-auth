package gen

//go:generate sqlc generate -f ./internal/store/pg/sqlc.yml

//go:generate go run ./cmd/tools/tern/tern.go
