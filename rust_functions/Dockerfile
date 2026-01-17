# Stage 1: build
FROM rust:1.71 as builder
WORKDIR /app
COPY funcs/add_age_plus_one.rs .
RUN rustc add_age_plus_one.rs -o add_age_plus_one

# Stage 2: runtime
FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/add_age_plus_one .

# Zorg dat de executable uitvoerbaar is
RUN chmod +x add_age_plus_one

# Default CMD
CMD ["./add_age_plus_one"]
