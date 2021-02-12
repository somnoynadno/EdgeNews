FROM golang:alpine
LABEL maintainer="Alexander Zorkin"

ENV ENV PRODUCTION
WORKDIR /app
ADD . /app/

RUN go build -o main .
CMD ["/app/main"]
EXPOSE 8080
