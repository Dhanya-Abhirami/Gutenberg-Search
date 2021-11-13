FROM golang:alpine

WORKDIR /server

RUN  go mod init server

RUN go get github.com/elastic/go-elasticsearch/v7

RUN go get github.com/gin-gonic/gin

COPY ./ ./

RUN go build

CMD ["go", "run", "."]

EXPOSE 8080