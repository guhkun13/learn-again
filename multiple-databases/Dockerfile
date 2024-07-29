FROM golang:alpine

RUN apk update && apk add --no-cache git gcc g++

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=1 GOOS=linux go build -o main


EXPOSE 80

ENTRYPOINT ["/app/main"]