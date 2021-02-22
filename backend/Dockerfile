FROM ubuntu:latest
LABEL maintainer="Alexander Zorkin"

ENV ENV PRODUCTION
RUN apt-get update && apt-get install -y firefox golang-go default-jdk curl wget

WORKDIR /app
ADD . /app/

RUN go build -o main .
CMD ["/app/main"]

EXPOSE 8080
