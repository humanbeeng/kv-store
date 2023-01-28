FROM golang:1.19-alpine as build-env

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /kv-store

EXPOSE 3000

CMD ["/kv-store"]
