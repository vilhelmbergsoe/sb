FROM golang:1.18-alpine

WORKDIR /usr/src/app

RUN apk add gcc musl-dev

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -v -o /usr/local/bin/sb ./...

CMD ["sb"]
