# ----------- Build stage -----------
FROM golang:latest AS build

WORKDIR /app

# Install migrate with sqlite3 driver
RUN GOBIN=/usr/local/bin go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Copy Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . ./

# Build binary
RUN CGO_ENABLED=1 go build -o app ./cmd/main.go

# ----------- Runtime stage -----------
FROM debian:bookworm-slim

WORKDIR /app

# Install runtime dependencies
RUN apt-get update && apt-get install -y ca-certificates curl sqlite3 && rm -rf /var/lib/apt/lists/*

# Copy binary and migrations
COPY --from=build /app/app .
COPY --from=build /usr/local/bin/migrate /usr/local/bin/migrate
COPY internal/db/migrations ./migrations

# Create data directory for SQLite
RUN mkdir -p /data

# Environment defaults
ENV DB_PATH=/data/nerdy_link_manager.db
ENV PORT=8080

# Expose port
EXPOSE 8080

# Run migrations then start app
CMD migrate -path ./migrations -database "sqlite3://$DB_PATH" up || true && ./app
