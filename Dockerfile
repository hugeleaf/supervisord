FROM golang:alpine AS builder

COPY . /src

WORKDIR /src

RUN GOPROXY=https://goproxy.io,direct go build -o /usr/bin/supervisord ./supervisor

FROM scratch

COPY --from=builder /usr/bin/supervisord /usr/bin/supervisord

ENTRYPOINT ["/usr/local/bin/supervisord"]
