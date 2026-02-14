FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/sprayer-api ./cmd/api
RUN go build -o /app/sprayer-cli ./cmd/cli

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/sprayer-api .
COPY --from=builder /app/sprayer-cli .
COPY prompts/ ./prompts/

# Install chrome/chromium dependencies if browser scraping is needed in container
# Using chromium
RUN apk add --no-cache chromium

ENV GO_ROD_BIN=/usr/bin/chromium-browser

EXPOSE 8080

CMD ["./sprayer-api"]
