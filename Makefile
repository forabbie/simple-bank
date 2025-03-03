
postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=changeme -d postgres

createdb:
	docker exec -it postgres createdb --username=root --owner=root simplebank

dropdb:
	docker exec -it postgres dropdb simplebank

migrateup:
	migrate -path database/migration -database "postgresql://root:changeme@localhost:5432/simplebank?sslmode=disable" -verbose up

migrateup1:
	migrate -path database/migration -database "postgresql://root:changeme@localhost:5432/simplebank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path database/migration -database "postgresql://root:changeme@localhost:5432/simplebank?sslmode=disable" -verbose down

migratedown1:
	migrate -path database/migration -database "postgresql://root:changeme@localhost:5432/simplebank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -destination database/mock/store.go github.com/forabbie/vank-app/database/sqlc Store

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test server mock
