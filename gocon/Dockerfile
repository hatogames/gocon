# 1️⃣ Build Stage
FROM golang:1.25 AS builder

WORKDIR /app

# Module laden
COPY go.mod ./
RUN go mod download

# Rest des Codes kopieren und bauen
COPY . .
RUN go build -o server .

# 2️⃣ Final Stage – kleineres Image
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]