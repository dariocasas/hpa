# stage 1

FROM golang:1.19 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s -w' -o /app/webserver 

# stage 2

FROM alpine:latest
# FROM scratch

WORKDIR /app/

RUN apk update && \
    apk upgrade && \
    apk add bash 

COPY --from=builder /app/webserver ./

CMD ["./webserver"] 