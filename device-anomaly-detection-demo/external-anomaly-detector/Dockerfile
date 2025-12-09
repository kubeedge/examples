# syntax=docker/dockerfile:1
FROM ghcr.io/astral-sh/uv:python3.10-bookworm-slim

WORKDIR /app

# Copy project files
COPY . .

# Install dependencies using uv
RUN uv sync --frozen --no-cache

CMD ["uv", "run", "python", "main.py"]
