FROM golang:1.14.6-alpine as base

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --update --no-cache ca-certificates git
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.io,direct

WORKDIR /app
COPY app .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o main
RUN chmod +x main
CMD ["./main"]