FROM ubuntu:latest

LABEL authors="coole"

RUN apt-get update && apt-get install -y ca-certificates golang

WORKDIR /app

COPY . .

RUN go build -o myapp

ENTRYPOINT ["./myapp"]
