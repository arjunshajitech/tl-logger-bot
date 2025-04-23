FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /tl-logger-bot

EXPOSE 3333

CMD ["/tl-logger-bot"]