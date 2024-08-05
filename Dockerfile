# syntax=docker/dockerfile:1

FROM golang:latest

WORKDIR /app/

COPY . .

RUN go build -mod=vendor -o /go-store

EXPOSE 8000
EXPOSE 50051

CMD [ "/go-store" ]