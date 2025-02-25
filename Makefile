postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=changeme -d postgres

createdb:
	docker exec -it postgres createdb --username=root --owner=root simplebank

dropdb:
	docker exec -it postgres dropdb simplebank

migrateup:
	migrate -path database/migration -database "postgresql://root:changeme@localhost:5432/simplebank?sslmode=disable" -verbose up

migratedown:
	migrate -path database/migration -database "postgresql://root:changeme@localhost:5432/simplebank?sslmode=disable" -verbose down

.PHONY: postgres createdb dropdb migrateup migratedown