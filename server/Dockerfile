FROM golang:latest
ENV GOPATH /go

COPY . /go/src/app

WORKDIR /go/src/app

RUN go mod tidy

RUN go build -o main .

CMD ["./main"]