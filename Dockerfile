#Build
FROM golang:1.18.2-alpine3.15 as builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

#Run
FROM alpine:3.14
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .

CMD [ "/app/main" ]
