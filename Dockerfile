FROM golang:1.15-alpine as builder

WORKDIR /github.com/links-japan/kakaku

ADD . .

# if you are in China, please uncomment this line to setup golang proxy
#RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download
RUN CGO_ENABLED=0 go build -o kakaku_cmd ./cmd/main.go

FROM alpine:3.12.0 as runner

COPY --from=builder /github.com/links-japan/kakaku/kakaku_cmd ./kakaku_cmd
