# ===============================
# Makefile voor ISE Docker builds
# ===============================

# Variabelen
PYTHON_IMAGE := strategy-python:latest
RUST_IMAGE   := strategy-rust-add-age:latest

# -------------------------------
# Python Docker image
# -------------------------------
docker-python:
	@echo "🔹 Building Python Docker image..."
	@docker build -t $(PYTHON_IMAGE) -f dockerfiles/python/Dockerfile . \
	|| { echo "❌ Python Docker build failed"; exit 1; }

# Rebuild Python image zonder cache
rebuild-python:
	@echo "🔹 Rebuilding Python Docker image (no cache)..."
	@docker build --no-cache -t $(PYTHON_IMAGE) -f dockerfiles/python/Dockerfile . \
	|| { echo "❌ Python Docker rebuild failed"; exit 1; }

# -------------------------------
# Rust Docker image
# -------------------------------
docker-rust:
	@echo "🔹 Building Rust Docker image..."
	@docker build -t $(RUST_IMAGE) -f dockerfiles/rust/Dockerfile . \
	|| { echo "❌ Rust Docker build failed"; exit 1; }

# Rebuild Rust image zonder cache
rebuild-rust:
	@echo "🔹 Rebuilding Rust Docker image (no cache)..."
	@docker build --no-cache -t $(RUST_IMAGE) -f dockerfiles/rust/Dockerfile . \
	|| { echo "❌ Rust Docker rebuild failed"; exit 1; }

# -------------------------------
# Build alle images
# -------------------------------
docker-all: docker-python docker-rust
	@echo "✅ All Docker images built"

# Rebuild alle images
rebuild-all: rebuild-python rebuild-rust
	@echo "✅ All Docker images rebuilt (no cache)"

# -------------------------------
# Clean local Docker images
# -------------------------------
docker-clean:
	@docker rmi -f $(PYTHON_IMAGE) $(RUST_IMAGE) || true
	@echo "🧹 Docker images removed"

# -------------------------------
# Run Go pipeline (debug mode)
# -------------------------------
run:
	go run . --debug

.PHONY: docker-python rebuild-python docker-rust rebuild-rust docker-all rebuild-all docker-clean run