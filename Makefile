.DEFAULT_GOAL := help

.PHONY: install lint format typecheck complexity-check \
        check lint-check format-check ci build clean \
        start-zot stop-zot help

help: ## Show available targets
	@awk 'BEGIN {FS = ":.*?##"; printf "Usage:\n  make <target>\n\nTargets:\n"} \
	/^[a-zA-Z_-]+:.*?##/ { printf "  %-20s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

install: ## Sync dependencies from uv.lock
	uv sync --frozen

lint: ## Run ruff with auto-fix
	uv run ruff check --fix .

format: ## Format code with ruff
	uv run ruff format .
	uv run ruff check --fix .

typecheck: ## Run ty type checker
	uv run ty check

complexity-check: ## Run complexipy (max complexity 15)
	uv run complexipy

build: ## Build package with uv
	uv build

lint-check: ## Run ruff without fixes (CI mode)
	uv run ruff check .

format-check: ## Check formatting without fixes (CI mode)
	uv run ruff format --check .

ci: ## Run full CI pipeline locally
	install lint-check format-check typecheck complexity-check build

clean: ## Remove caches and build artifacts
	rm -rf .ruff_cache .complexipy_cache dist build
	find . -type d -name __pycache__ -exec rm -rf {} + 2>/dev/null || true

start-zot: ## Start local Zot OCI registry (podman)
	$(MAKE) -C test start-zot

stop-zot: ## Stop local Zot OCI registry (podman)
	$(MAKE) -C test stop-zot
