FROM ubuntu:latest
LABEL maintainer="Alexander Zorkin"

ENV TZ  Europe/Moscow
ENV ENV PRODUCTION

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apt-get update && apt-get install -y firefox golang-go default-jdk xvfb curl wget

WORKDIR /app
ADD . /app/

RUN go build -o main .
CMD ["/app/main"]

EXPOSE 8080
