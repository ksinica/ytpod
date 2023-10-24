FROM golang:1.21.3-alpine3.18 AS builder

WORKDIR /usr/src/app

COPY . .

RUN go install github.com/ksinica/ytpod

FROM alpine:3.18

COPY --from=builder /go/bin/ytpod /usr/bin/ytpod

ENTRYPOINT ["/usr/bin/ytpod"]