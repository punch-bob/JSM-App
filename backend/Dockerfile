FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY *.go .
RUN go build -o /simple-rest-server

ENV SERVER_PORT=$SERVER_PORT
ENV DB_LOGIN=$DB_LOGIN
ENV DB_NAME=$DB_NAME
ENV DB_PASSWORD=$DB_PASSWORD

EXPOSE 8080

CMD [ "/simple-rest-server" ]