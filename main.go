package main

import (
	"fmt"
	"log"

	"assignment/config"
	"assignment/rabbitmq"

	"github.com/streadway/amqp"
)

func main() {
	fmt.Println("### start worker ###")

	cfg, err := config.NewConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	connect, err := amqp.Dial(cfg.RabbitmqURL)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	consumer := rabbitmq.NewConsumer(connect)

	forever := make(chan bool)
	go func() {
		err := consumer.StartConsumer("source")
		if err != nil {
			log.Fatal(err)
		}
	}()
	<-forever
}
