FROM golang:1.22-alpine

WORKDIR /home/app

COPY . .

RUN go install github.com/cosmtrek/air@latest

EXPOSE 80

ENTRYPOINT ["air"]
