
migrate_import:
	go get github.com/golang-migrate/migrate/v4/database/postgres

migrate_create:
	migrate create -ext sql -dir migrations -seq login

migrate_up:
	migrate -database postgres://postgres:3115@localhost:5432/login?sslmode=disable -path ./migrations up

migrate_down:
	migrate -database postgres://postgres:3115@localhost:5432/login?sslmode=disable -path ./migrations up

swag-gen:
  ~/go/bin/swag init -g ./api/handler.go -o api/docs force 1