FROM golang:latest AS build

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . ./

RUN GOOS=linux GOARCH=amd64 go build -o go-ssh-app ./cmd

RUN ssh-keygen -t rsa -f id_rsa -N ""

FROM alpine:latest

COPY --from=build app/id_rsa app/go-ssh-app /

CMD ["/go-ssh-app"]
