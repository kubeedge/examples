# syntax=docker/dockerfile:1
FROM ghcr.io/astral-sh/uv:python3.10-bookworm-slim

WORKDIR /app

# Copy project files
COPY . .

# Install dependencies using uv
RUN uv sync --frozen --no-cache

EXPOSE 8000

CMD ["uv", "run", "uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
