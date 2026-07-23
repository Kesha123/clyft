.PHONY: install lint format typecheck complexity-check \
        check lint-check format-check ci build clean

install:
	uv sync --frozen

lint:
	uv run ruff check --fix .

format:
	uv run ruff format .
	uv run ruff check --fix .

typecheck:
	uv run ty check

complexity-check:
	uv run complexipy

build:
	uv build

lint-check:
	uv run ruff check .

format-check:
	uv run ruff format --check .

ci: install lint-check format-check typecheck complexity-check build

clean:
	rm -rf .ruff_cache .complexipy_cache dist build
	find . -type d -name __pycache__ -exec rm -rf {} + 2>/dev/null || true
