FROM golang:1.20-alpine AS builder
WORKDIR /app
ENV CGO_ENABLED 1
RUN apk add gcc && apk --no-cache --update add build-base
COPY . .
RUN go build main.go

FROM alpine
LABEL org.opencontainers.image.authors="alireza7@gmail.com"
ENV TZ=Asia/Tehran
WORKDIR /app

RUN apk add ca-certificates tzdata && mkdir bin

# Download latest rule files
ADD https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geoip.dat \
    https://github.com/Loyalsoldier/v2ray-rules-dat/releases/latest/download/geosite.dat \
    bin/

COPY --from=builder  /app/main /app/x-ui
VOLUME [ "/etc/x-ui" ]
CMD [ "./x-ui" ]