FROM golang:1.15

WORKDIR /somewhere

ADD go.mod .
ADD go.sum .

RUN go mod download

ADD main.go .

RUN go build -o test-etcd .

ENTRYPOINT ["./test-etcd"]
