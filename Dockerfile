FROM golang:1.15-alpine as builder

WORKDIR /github.com/links-japan/kakaku

ADD . .

# if you are in China, please uncomment this line to setup golang proxy
#RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download
RUN CGO_ENABLED=0 go build -o worker_cmd ./cmd/kakaku
RUN CGO_ENABLED=0 go build -o server_cmd ./cmd/server

RUN GRPC_HEALTH_PROBE_VERSION=v0.3.2 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

FROM alpine:3.12.0 as runner

COPY --from=builder /github.com/links-japan/kakaku/worker_cmd /bin/worker_cmd
COPY --from=builder /github.com/links-japan/kakaku/server_cmd /bin/server_cmd
