FROM golang:1.19 as build

WORKDIR /go/src/app

COPY go.mod . 
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -installsuffix nocgo -o app-bin .

FROM alpine:edge
RUN apk upgrade --update-cache --available && \
    apk add openssl && \
    apk add ca-certificates && \
    rm -rf /var/cache/apk/*

ADD https://letsencrypt.org/certs/isrgrootx1.pem.txt /usr/local/share/ca-certificates/isrgrootx1.pem
ADD https://letsencrypt.org/certs/trustid-x3-root.pem.txt /usr/local/share/ca-certificates/trustid-x3-root.pem

RUN cd /usr/local/share/ca-certificates \
    && openssl x509 -in isrgrootx1.pem -inform PEM -out isrgrootx1.crt \
    && openssl x509 -in trustid-x3-root.pem -inform PEM -out trustid-x3-root.crt \
    && update-ca-certificates
COPY  --from=build /go/src/app/app-bin ./

ENTRYPOINT ["./app-bin"]
