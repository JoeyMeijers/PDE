docker-python:
	docker build -t strategy-python:latest ./python_functions

docker-rust:
	docker build -t strategy-rust-add-age:latest ./rust_functions
