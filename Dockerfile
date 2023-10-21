# build stage
FROM golang:1.21.0-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

#Getting rid of this as we're going to do it in our code.
#RUN apk add curl
#RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
#Getting rid of this as we're going to do it in our code.
#COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY start.sh .
#RUN chmod +x start.sh
COPY wait-for.sh .
#RUN chmod +x wait-for.sh
COPY db/migration ./db/migration

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]


