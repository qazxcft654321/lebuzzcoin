FROM golang:1.14

RUN go get -u github.com/FiloSottile/mkcert
RUN cd /go/src/github.com/FiloSottile/mkcert && go build -o /bin/mkcert

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY . .

RUN go mod download
RUN go build -o app

WORKDIR /app

RUN cp /build/app .
RUN cp /build/VERSION .
RUN cp /build/env.example ./.env
RUN rm -rf /build
RUN mkdir tls
RUN cd tls && mkcert -install && mkcert localhost

CMD ["./app"]