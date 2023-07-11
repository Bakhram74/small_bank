postgres:
	docker run --name postgres12 --network bank_network  -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1234 -d postgres:latest

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root small_bank

dropdb:
	docker exec -it postgres12 dropdb small_bank

migrateup:
	 migrate -path db/migration/ -database "postgresql://root:1234@localhost:5436/small_bank?sslmode=disable" -verbose up

migrateup1:
	 migrate -path db/migration/ -database "postgresql://root:1234@localhost:5436/small_bank?sslmode=disable" -verbose up 1

migratedown:
	 migrate -path db/migration/ -database "postgresql://root:1234@localhost:5436/small_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration/ -database "postgresql://root:1234@localhost:5436/small_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/Bakhram74/small_bank/db/sqlc  Store

.PHONY:createdb postgres dropdb migrateup migrateup1  migratedown1 sqlc test server mock