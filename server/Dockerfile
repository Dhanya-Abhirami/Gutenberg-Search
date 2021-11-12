FROM golang:alpine

WORKDIR /server

RUN  go mod init server

RUN go get github.com/elastic/go-elasticsearch/v7

COPY ./ ./

CMD ["go", "run", "main.go"]

EXPOSE 8080