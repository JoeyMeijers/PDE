# Base image
FROM python:3.11-slim

# Werkdirectory
WORKDIR /app

# Kopieer code
COPY main.py .
COPY funcs/ ./funcs/

# Install dependencies (geen in dit voorbeeld)
# RUN pip install ...

# Default CMD
CMD ["python", "main.py"]
