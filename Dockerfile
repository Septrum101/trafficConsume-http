FROM golang:alpine AS builder

WORKDIR /build
COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod tidy
RUN go build -trimpath -v -ldflags "-X main.date=$(date -Iseconds)" -o trafficConsume

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /build/trafficConsume /app/trafficConsume
COPY --from=builder /build/urls.txt /app/urls.txt
ENTRYPOINT ["/app/trafficConsume"]