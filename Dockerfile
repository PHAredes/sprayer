FROM python:3.11-alpine

# Install minimal dependencies
RUN apk add --no-cache \
    git \
    build-base \
    libffi-dev \
    openssl-dev

# Set working directory
WORKDIR /app

# Copy requirements file
COPY requirements.txt .

# Install Python dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Create directories for your templates and outputs
RUN mkdir -p /app/cv_templates /app/email_templates /app/outputs

# Environment variables will be passed at runtime
ENV PYTHONUNBUFFERED=1

# Keep container running for interactive use
CMD ["tail", "-f", "/dev/null"]
