FROM golang:1.24-alpine

RUN mkdir /app
WORKDIR /app

COPY go.mod /app
COPY go.sum /app

RUN go mod download

COPY cmd /app/cmd
COPY internal /app/internal
COPY pkg /app/pkg
COPY main.go /app


RUN go build -o /srv/main /app

CMD ["/srv/main", "server"]
