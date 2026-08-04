ZOT_COMPOSE     := e2e/zot/compose.yaml
COMPOSE_PROJECT ?= local
CONTAINER_RUNTIME 	?= podman

.PHONY: start-zot stop-zot e2e-tests e2e-test/%

start-zot:
	$(CONTAINER_RUNTIME) compose -p $(COMPOSE_PROJECT) -f $(ZOT_COMPOSE) up -d

stop-zot:
	$(CONTAINER_RUNTIME) compose -p $(COMPOSE_PROJECT) -f $(ZOT_COMPOSE) down -v

e2e-tests:

e2e-test/%:
