package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	queueName, err := os.Hostname()
	userName := flag.String("u", "guest", "username")
	userPass := flag.String("P", "guest", "password")
	hostName := flag.String("h", "localhost", "host")
	vhost := flag.String("v", "/", "vhost")
	port := flag.Int("p", 5672, "port")
	flag.Parse()
	connAddress := "amqp://" + *userName + ":" + *userPass + "@" + *hostName + ":" + strconv.Itoa(*port) + *vhost
	log.Printf("Connecting to: %s with queue name: %s", connAddress, queueName)

	conn, err := amqp.Dial(connAddress)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := "alive"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf("Sent: %s", body)
	failOnError(err, "Failed to publish a message")

	time.Sleep(1 * time.Second)

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			log.Printf("Received: %s", d.Body)
			os.Exit(0)
		}
		os.Exit(1)
	}()

}
