// go run amqp_server.go
// based on https://www.rabbitmq.com/tutorials/tutorial-six-go.html

package main

import (
	"../image"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func get_response(request map[string]string, err error) map[string]string {
	fmt.Println(request)
	if err != nil {
		return map[string]string{"message": "error in contrast enhancement. check image format"}
	}
	return map[string]string{"message": "success!"}
}

func main() {
	conn, err := amqp.Dial("amqp://localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"rpc_queue", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var request map[string]string
			json.Unmarshal(d.Body, &request)
			failOnError(err, "Failed to convert body")

			log.Printf("%s", request)
			err := image.ContrastEnhancement(request["image"], false)
			response := get_response(request, err)
			body, _ := json.Marshal(response)

			err = ch.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          body,
				})
			failOnError(err, "Failed to publish a message")

			d.Ack(false)
		}
	}()

	log.Printf("Awaiting RPC requests")
	<-forever
}
