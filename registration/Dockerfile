FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY *.go .
RUN go build -o /auth-server

ENV SERVER_PORT=$SERVER_PORT
ENV DB_LOGIN=$DB_LOGIN
ENV DB_NAME=$DB_NAME
ENV DB_PASSWORD=$DB_PASSWORD
ENV DB_HOST=$DB_HOST
ENV DB_PORT=$DB_PORT

EXPOSE 8080

CMD [ "/auth-server" ]