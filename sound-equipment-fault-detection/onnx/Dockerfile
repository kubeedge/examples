# Use Python official base image
FROM python:3.10-slim

# Set working directory
WORKDIR /app

# Copy the rest of the application code
COPY . /app

# Install Python dependencies and clear cache
RUN pip install --upgrade pip \
    && pip install --no-cache-dir onnxruntime librosa numpy