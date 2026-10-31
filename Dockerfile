FROM golang:l.27.3-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/crypt-lab-1 ./src


FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/crypt-lab-1 .

CMD ["/app/crypt-lab-1"]