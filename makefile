# Makefile

dev-start:
	@echo "Starting the DEV server..."
	@docker-compose -f docker-compose-dev-local.yml build
	@docker-compose -f docker-compose-dev-local.yml up

dev-stop:
	@echo "Stopping the DEV server..."
	@docker-compose -f docker-compose-dev-local.yml down -v
	@docker container prune -f



