package main

import (
	"encoding/json"
	"fmt"
	"github.com/ksooo091/go-rabbitmq/mailSender"
	amqp "github.com/rabbitmq/amqp091-go"
	viper "github.com/spf13/viper"
	"log"
)

type Mail struct {
	UserMail string `json:"user_mail"`
	MailType string `json:"mail_type"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	viper.SetConfigFile("config/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		viper.SetConfigFile("config.yaml")
		err = viper.ReadInConfig()
		if err != nil {
			fmt.Printf("Error reading config file: %s\n", err)
			return
		}
	}

	// 설정에서 값을 읽어옴
	mqFullURL := viper.GetString("rabbitmq.fullURL")
	senderMail := viper.GetString("smtp.mail")
	senderPw := viper.GetString("smtp.password")

	conn, err := amqp.Dial(mqFullURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello1", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

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

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			mailSend(d.Body,senderMail, senderPw)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func mailSend(jsonString []byte,senderMail, senderPw string) {
	var mailData Mail

	err := json.Unmarshal(jsonString, &mailData)
	fmt.Println(mailData)
	// EXCEPTION
	if err != nil {
		fmt.Println("Failed to json.Unmarshal", err)
	}else {
		mailSender.SendMail(mailData.UserMail,mailData.MailType, senderMail, senderPw)
	}

}
