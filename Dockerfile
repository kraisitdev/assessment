FROM golang:1.19-alpine as build-base

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test -v --tags=unit ./...
RUN go build -o ./out/go-sample .

# ====================

FROM alpine:3.16.2

RUN apk update && apk upgrade
RUN apk add tzdata

COPY --from=build-base /app/out/go-sample /app/go-sample

CMD ["/app/go-sample"]