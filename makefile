# Makefile

dev-start:
	@echo "Starting the DEV server..."
	@docker-compose -f docker-compose-dev-local.yml build
	@docker-compose -f docker-compose-dev-local.yml up

dev-stop:
	@echo "Stopping the DEV server..."
	@docker-compose -f docker-compose-dev-local.yml down -v

CICD-start:
	@echo "Starting the DEV server..."
	@docker-compose -f docker-compose-dev-CICD.yml build
	@docker-compose -f docker-compose-dev-CICD.yml up

CICD-stop:
	@echo "Stopping the DEV server..."
	@docker-compose -f docker-compose-dev-CICD.yml down -v



