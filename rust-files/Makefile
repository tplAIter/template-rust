.DEFAULT_GOAL := help

.PHONY: help fmt lint test build run infra-up infra-down
help:
	@echo "make fmt lint test build run infra-up infra-down"
fmt:
	cargo fmt --all
lint:
	cargo clippy --all-targets -- -D warnings
test:
	cargo test
build:
	cargo build --release
run:
	@if [ -f .env ]; then set -a; . ./.env; set +a; fi; cargo run
infra-up:
	docker compose up -d postgres
infra-down:
	docker compose down
