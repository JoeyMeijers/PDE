FROM rust:1.71 as builder
WORKDIR /app

COPY rust_functions/Cargo.toml .
COPY rust_functions/src ./src

RUN cargo build --release

FROM debian:bookworm-slim
WORKDIR /app

COPY --from=builder /app/target/release/ .