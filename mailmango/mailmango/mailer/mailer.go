package mailer

import (
	"encoding/json"
	"log"
	"math"
	"net"
	"strconv"
	"time"

	"gopkg.in/mail.v2"
)

type MailerService struct {
	HostPort string
	User     string
	Passcode string
}
type message struct {
	To      string `json:"to,omitempty"`
	Subject string `json:"subject,omitempty"`
	Body    string `json:"body,omitempty"`
}

func (ms *MailerService) SendMail(jsonBody []byte) {
	var msg message

	if err := json.Unmarshal(jsonBody, &msg); err != nil {
		log.Fatal(err)
	}

	m := mail.NewMessage()

	m.SetHeader("From", ms.User)
	m.SetHeader("To", msg.To)
	m.SetHeader("Subject", msg.Subject)

	m.SetBody("text/html", msg.Body)
	host, port_str, _ := net.SplitHostPort(ms.HostPort)
	port_number, _ := strconv.Atoi(port_str)

	var err error = nil

	for i := 0; i < 10; i++ {
		d := mail.NewDialer(host, port_number, ms.User, ms.Passcode)
		if err = d.DialAndSend(m); err == nil {
			break
		}
		timeToWait := math.Pow(2, float64(i))
		log.Println("waited ", timeToWait, "seconds to retry")
		time.Sleep(time.Duration(timeToWait) * time.Second)
	}
	if err != nil {
		log.Fatal(err)
	}
}
