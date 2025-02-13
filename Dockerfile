FROM ubuntu:latest

LABEL authors="coolerSAI"

RUN apt-get update && apt-get install -y \
    ca-certificates \
    golang \
    && apt-get clean

WORKDIR /app

COPY . .

RUN go build -o myapp

EXPOSE 8080

ENTRYPOINT ["./myapp"]
