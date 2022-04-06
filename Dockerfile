FROM golang:1.17-alpine AS base

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o build/main main.go
ENV RUN_ENV=docker_dev

EXPOSE 3030

CMD [ "build/main" ]




FROM alpine:latest as prod

COPY --from=base /app/build/main /usr/local/bin/pos
EXPOSE 3030

ENTRYPOINT ["/usr/local/bin/pos"]


