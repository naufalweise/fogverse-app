FROM golang:1.25.0 AS build
WORKDIR /src
COPY . .
RUN go mod download
RUN go build -v -o /bin/main

CMD ["/bin/main"]