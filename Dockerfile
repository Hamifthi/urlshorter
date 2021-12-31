FROM golang:1.17

RUN apt update && apt upgrade -y && \
    apt install -y git \
    make openssh-client

WORKDIR /go/src/urlshortner
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 8080

CMD ["go", "run", "./cmd/app/main.go"]