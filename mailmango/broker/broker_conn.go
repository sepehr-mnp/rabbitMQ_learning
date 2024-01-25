package broker

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type WorkerMQ struct {
	Mailer    MailerService
	HostPort  string
	QueueName string
}
type MailerService interface {
	SendMail([]byte)
}

func (wmq *WorkerMQ) StartConsume() {
	conn, err := amqp.Dial(wmq.HostPort)
	logErr(err, "connection failed")
	defer conn.Close()

	ch, err := conn.Channel()
	logErr(err, "channel not connected")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		wmq.QueueName, // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	logErr(err, "fail on declare queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	logErr(err, "fail on register consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			// TODO: implement validators
			wmq.Mailer.SendMail(d.Body)
			log.Printf("Mail sent to: %s", d.Body)
			//d.Ack(false)
		}
	}()
	log.Printf("[*] Waiting for messages. To exit press CTRL+C")

	<-forever
}

func logErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", err, msg)
	}
}
