FROM golang:alpine as builder

LABEL maintainer="sergio <raven1901@mail.ru>"

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
 
RUN go mod download 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./main"]
