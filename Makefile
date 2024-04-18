CURRENT_DIR=$(shell pwd)

build:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${APP} ${APP_CMD_DIR}/main.go

proto-gen:
	./scripts/gen-proto.sh	${CURRENT_DIR}
	ls genproto/*/*.pb.go | xargs -n1 -IX bash -c "sed -e '/bool/ s/,omitempty//' X > X.tmp && mv X{.tmp,}"


# migrate
.PHONY: create-migration
create-migration:
	migrate create -ext sql -dir migrations -seq "$(name)"

.PHONY: migrate-up
migrate-up:
	migrate -source file://migrations -database postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable up

.PHONY: migrate-down
migrate-down:
	migrate -source file://migrations -database postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable down

.PHONY: migration-version
migration-version:
	migrate -database file://migrations - database postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable -path migrations version 

.PHONY: migrate-dirty
migrate-dirty:
	migrate -path ./migrations/ -database file://migrations - database postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable force "$(number)"