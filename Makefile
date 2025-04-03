DB_URL=postgresql://root:1234@localhost:5434/simple_bank?sslmode=disable

postgres:
	sudo docker run --name postgres_container --network bank-network -p 5434:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1234 -d postgres:16.3

createdb:
	sudo docker exec -it postgres_container createdb --username=root --owner=root simple_bank

dropdb:
	sudo docker exec -it postgres_container dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)


sqlc:
	sqlc generate

pqInstall:
	go get github.com/lib/pq

testifyInstall:
	go get github.com/stretchr/testify

test:
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/Iowel/course-simple-bank/db/sqlc Store

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
    proto/*.proto
	statik -src=./doc/swagger -dest=./doc -f

evans:
	evans --host localhost --port 9090 -r repl

redis:
	docker run --name redis -p 6379:6379 -d redis:7.4-alpine

.PHONY: postgres createdb dropdb migrateup migratedown sqlc pqinstall testifyInstall test server mock migratedown1 migrateup1 db_docs db_schema proto evans redis new_migration