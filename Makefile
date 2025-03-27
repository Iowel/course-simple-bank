postgres:
	sudo docker run --name postgres_container -p 5434:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=1234 -d postgres:16.3


createdb:
	sudo docker exec -it postgres_container createdb --username=root --owner=root simple_bank

dropdb:
	sudo docker exec -it postgres_container dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:1234@localhost:5434/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:1234@localhost:5434/simple_bank?sslmode=disable" -verbose down

sqlc:
	sqlc generate

pqInstall:
	go get github.com/lib/pq

testifyInstall:
	go get github.com/stretchr/testify

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/Iowel/course-simple-bank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc pqinstall testifyInstall test server mock