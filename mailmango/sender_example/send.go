// JUST A EXAMPLE

package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	sendMessageBroker()
}
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
func sendMessageBroker() {
	// conn, _ := amqp.Dial("amqp://root:root@localhost:5672")
	// defer conn.Close()
	// ch, _ := conn.Channel()
	// defer ch.Close()
	// // q, _ := ch.QueueDeclare(
	// // 	"mail",       // name
	// // 	true,         // durable
	// // 	false,        // delete when unused
	// // 	false,        // exclusive
	// // 	false,        // no-wait
	// // 	nil,          // arguments
	// // )
	// b, _ := json.Marshal(map[string]interface{}{
	// 	"to":      "smirnasrollahi@email.com", // change for tests
	// 	"subject": "Email confirmation!",
	// 	"body": `<p>` + "Erick, " + `Please confirm your email by clicking the link below:</p>
	// 			<p> <a> https://` + "test.com" + `/api/auth/verify/signup/` + "14510" + `/` + "321sad65124e1298712313215!@e" + `</a></p>`,
	// })
	// ch.Publish(
	// 	"",
	// 	"mail",
	// 	false,
	// 	false,
	// 	amqp.Publishing{
	// 		DeliveryMode: amqp.Persistent,
	// 		ContentType:  "application/json",
	// 		Body:         b,
	// 	},
	// )
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"mailSys", // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	b, _ := json.Marshal(map[string]interface{}{
		"to":      "smirnasrollahi@gmail.com", // change for tests
		"subject": "Email confirmation!",
		"body":    "hi",
	})

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(b),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", b)
}
