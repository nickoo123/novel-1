FROM golang:1.12.2-alpine3.9 AS builder

MAINTAINER william

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update && apk --no-cache add build-base tzdata git

COPY . /code
WORKDIR /code
RUN export GOPROXY=https://goproxy.cn
RUN go mod tidy && go mod vendor && \
    CGO_ENABLED=1 go build -a

RUN pwd && ls

FROM alpine

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN mkdir -p /go && export GOTMPDIR=/go

WORKDIR /go
COPY --from=builder /code/novel /go
COPY --from=builder /code/conf /go/conf
COPY --from=builder /code/lang /go/lang
COPY --from=builder /code/static /go/static
RUN mkdir -p /go/static/sitemap
COPY --from=builder /code/views /go/views
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

EXPOSE 8081

CMD ["./novel"]