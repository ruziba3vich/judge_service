proto:
	./generate_protos.sh

MIGRATIONS_DIR=./migrations

create-migration:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq judge_service
