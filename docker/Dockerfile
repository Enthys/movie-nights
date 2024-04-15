FROM golang:1.22-alpine

WORKDIR /home/app

COPY . .

RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/a-h/templ/cmd/templ@latest

RUN apk add --update npm

WORKDIR /home/app/client
RUN npm install

EXPOSE 80

ENTRYPOINT ["air"]
