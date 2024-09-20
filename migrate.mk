# migrate options
MIGRATE_DRV=postgres
POSTGRES_HOST?=localhost
MIGRATE_DSN?=${PG_DSN}
MIGRATE_DSN?=host=localhost user=goorthauer password=postgres dbname=hmb_bot port=5432 sslmode=disable
MIGRATE_DIR=./migrations/migrate
GOOSE_BASE_CMD=goose -dir ${MIGRATE_DIR} ${MIGRATE_DRV} "${MIGRATE_DSN}"

migration-up:		## Migrate the DB to the most recent version available
	${GOOSE_BASE_CMD} up

migration-down: 	## Roll back the version by 1
	${GOOSE_BASE_CMD} down

migration-reset:	## Roll back all migrations
	${GOOSE_BASE_CMD} reset

.PHONY: migration-up migration-down migration-reset


migration-save-scheme:	## Generate current DDL from playground
migration-save-scheme:
	docker exec $(shell basename $(shell pwd))-postgres-local pg_dump -U postgres --schema-only paas_db > ${PWD}/migrations/structure.sql
	docker exec $(shell basename $(shell pwd))-postgres-local pg_dump -U postgres --data-only --table=goose_db_version paas_db > ${PWD}/migrations/applied_migrations.sql

.PHONY: migration-save-scheme
