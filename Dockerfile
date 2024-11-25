FROM ubuntu:latest

LABEL authors="coolerSAI"

RUN apt-get update && apt-get install -y ca-certificates golang

WORKDIR /app

COPY . .

RUN go build -o myapp

ENTRYPOINT ["./myapp"]
