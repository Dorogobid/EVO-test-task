FROM golang:latest
WORKDIR /usr/src/app
COPY ./ ./
RUN go build -o main .
CMD ["./main"]