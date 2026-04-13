FROM golang:1.24-alpine AS builder

WORKDIR /build
COPY backend/go.mod backend/go.sum* ./
RUN go mod download 2>/dev/null || true
COPY backend/ .
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -o casamia .

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=builder /build/casamia .

# Frontend static files
COPY frontend/ /app/static/

# Seed images from frontend
RUN mkdir -p /app/uploads /app/seed-images
COPY frontend/images/ /app/seed-images/

ENV PORT=3000
ENV STATIC_DIR=/app/static
ENV UPLOAD_DIR=/app/uploads

EXPOSE 3000

CMD ["./casamia"]
