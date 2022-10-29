FROM golang:latest
WORKDIR /usr/src/app
COPY go.mod go.sum ./
COPY ./ ./
RUN apt-get update
RUN apt-get -y install postgresql-client
RUN go mod download && go mod verify
RUN go build -o main .
CMD ["./main"]