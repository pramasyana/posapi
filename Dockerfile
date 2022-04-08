FROM golang:1.17 as builder
WORKDIR /go/src/github.com/ryanpramasyana/posapi/
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
WORKDIR /home/ryanpramasyana/
ADD .env ./
COPY --from=builder /go/src/github.com/ryanpramasyana/posapi/app ./

CMD ./app