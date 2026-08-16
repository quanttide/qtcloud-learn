# 学习云 provider：多阶段构建。
# 数据持久化走 OSS（OSS_BUCKET 环境变量注入，FC 实例盘不持久），凭证经 FC 环境变量。

FROM golang:1.26-alpine AS build

RUN apk add --no-cache gcc musl-dev

WORKDIR /src

ARG GOPROXY=https://goproxy.cn,direct
ENV GOPROXY=${GOPROXY}

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/learn-server ./cmd/server

FROM alpine:3.20

RUN addgroup -S app && adduser -S -G app app
COPY --from=build /out/learn-server /usr/local/bin/learn-server

USER app
EXPOSE 8080

CMD ["/usr/local/bin/learn-server"]
