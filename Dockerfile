# build
FROM golang:1.14-alpine as builder

RUN apk update && apk add tzdata \
    && cp /usr/share/zoneinfo/Asia/Bangkok /etc/localtime \
    && echo "Asia/Bangkok" >  /etc/timezone \
    && apk del tzdata

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go clean --cache

RUN go build \
    -ldflags "-X main.buildcommit=$(cat .git/refs/heads/develop) -X main.buildtime=$(date +%Y%m%d.%H%M%S)" \
    -o goapp main.go

# ---------------------------------------------------------

# run
FROM alpine:latest

RUN apk update && apk add tzdata \
    && cp /usr/share/zoneinfo/Asia/Bangkok /etc/localtime \
    && echo "Asia/Bangkok" >  /etc/timezone \
    && apk del tzdata

WORKDIR /app

COPY --from=builder /app/goapp .

CMD ["./goapp"]
