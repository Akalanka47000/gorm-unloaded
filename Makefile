# Uses docker compose to start the infrastructure for running all examples
infra-up:
	@echo "Starting local infrastructure..."
	@COMPOSE_BAKE=true docker compose -f docker-compose.yml up -d

# Teardown the infrastructure, removing containers but keeping volumes for data persistence. 
infra-down:
	@docker compose-f docker-compose.yml down

# Clean up the infrastructure by tearing down and removing all containers and associated volumes to ensure
# a completely clean state. 
infra-clean:
	@echo "Cleaning up local infrastructure..."
	@docker compose -f docker-compose.yml down -v

# Rebuilds the infrastructure by tearing down and bringing it back up.
# Removes volumes to ensure a clean state.
infra-rebuild:
	@echo "Rebuilding infrastructure..."
	@make infra-clean
	@make infra-up