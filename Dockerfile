FROM golang:1.27.3-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/crypt-lab-1 ./src


FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache python3 py3-matplotlib
COPY --from=builder /app/crypt-lab-1 .
COPY scripts/ /app/scripts/
CMD ["sh", "-c","/app/crypt-lab-1 && python3 /app/scripts/plot.py"]