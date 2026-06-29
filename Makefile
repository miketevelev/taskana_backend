include .env
export

export PROJECT_ROOT=$(shell pwd)

env-up:
	@docker compose up -d taskana-db

env-down:
	@docker compose down taskana-db

env-cleanup:
	@read -p "Make cleanup volume for taskana? [y/n]: " ans; \
	if [ "$$ans" = "y" ]; then \
	  docker compose down taskana-db taskana-port-forwarder && \
	  rm -rf ${PROJECT_ROOT}/out/pgdata && \
	  echo "All volume files is deleted"; \
	else \
	  echo "Cleanup is decline"; \
	fi

env-port-forward:
	@docker compose up -d taskana-port-forwarder

env-port-close:
	@docker compose down taskana-port-forwarder

# Create a new migration version
# make migrate-create seq=new_version
migrate-create:
	@if [ -z "$(seq)" ]; then \
  	  echo "Missing parameter 'seq', for example: make migrate-create seq=init"; \
  	  exit 1; \
  	fi; \
	docker compose run --rm taskana-db-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
      	  echo "Missing parameter 'action', for example: make migrate-action action=up"; \
      	  exit 1; \
	fi; \
	docker compose run --rm taskana-db-migrate \
    	-path /migrations \
    	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@taskana-db:5432/${POSTGRES_DB}?sslmode=disable \
    	"$(action)"

backend-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/backend/main.go

backend-deploy:
	@docker compose up -d --build taskana

backend-undeploy:
	@docker compose down taskana

logs-cleanup:
	@read -p "Make cleanup logs for taskana? [y/n]: " ans; \
	if [ "$$ans" = "y" ]; then \
	  rm -rf ${PROJECT_ROOT}/out/logs && \
	  echo "All logs files is deleted"; \
	else \
	  echo "Logs cleanup is decline"; \
	fi
