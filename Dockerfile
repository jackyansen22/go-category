# ========================
# Build stage
# ========================
FROM golang:1.23-alpine AS builder

WORKDIR /app

# cache dependency
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# build dari standard layout
RUN go build -o server ./cmd/api


# ========================
# Runtime stage
# ========================
FROM alpine:3.19

WORKDIR /app

# copy binary
COPY --from=builder /app/server .

# expose port
EXPOSE 8080

# run
CMD ["/app/server"]

