FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /vertex-api

EXPOSE 8080

CMD ["/vertex-api"] 