postgres:
	docker run --name postgres12 -p 5436:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1234 -d postgres:latest

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root small_bank

dropdb:
	docker exec -it postgres12 dropdb small_bank

migrateup:
	 migrate -path db/migration/ -database "postgresql://root:1234@localhost:5436/small_bank?sslmode=disable" -verbose up

migratedown:
	 migrate -path db/migration/ -database "postgresql://root:1234@localhost:5436/small_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/Bakhram74/small_bank/db/sqlc  Store

.PHONY:createdb postgres dropdb migrateup migratedown sqlc test server mock