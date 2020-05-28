FROM golang:1.10
WORKDIR /app
COPY ./rmq_healthz.go .
RUN go get github.com/streadway/amqp && go build rmq_healthz.go

FROM rabbitmq:3.8.4-management
COPY --from=0 /app/rmq_healthz /

