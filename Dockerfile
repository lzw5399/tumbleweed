# build stage
FROM golang:1.16 as builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" .

RUN mkdir publish && cp workflow publish && \
    mkdir publish/config && \
    cp src/config/appsettings.yaml publish/config/

FROM alpine:3.12

WORKDIR /app

COPY --from=builder /app .

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk update

# set timezone to Asia/Shanghai
RUN apk update && \
    apk add tzdata && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone
ENV TZ Asia/Shanghai

ENV APP_ENV=Production \
    PORT=5000

EXPOSE 5000

ENTRYPOINT ["./workflow"]