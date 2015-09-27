package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/st3v/cfkit/service"
	"github.com/streadway/amqp"
)

const queueName = "testapp"

var rabbit *service.RabbitMQ

func main() {
	var err error
	if rabbit, err = service.Rabbit(); err != nil {
		log.Fatalf("Error getting RabbitMQ service: %s", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/", getMessageHandler).Methods("GET")
	router.HandleFunc("/", postMessageHandler).Methods("POST")

	http.Handle("/", router)

	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(addr, nil))
}

func postMessageHandler(rw http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Fprint(rw, "Invalid request body")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	ch, q, err := declareQueue()
	if err != nil {
		fmt.Fprint(rw, err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer ch.Close()

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: 2, // persistent
			ContentType:  "text/plain",
			Body:         body,
		},
	)
	if err != nil {
		fmt.Fprintf(rw, "Error publishing message: %s", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func getMessageHandler(rw http.ResponseWriter, req *http.Request) {
	ch, q, err := declareQueue()
	if err != nil {
		fmt.Fprint(rw, err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		fmt.Fprintf(rw, "Error registering consumer", err)
		rw.WriteHeader(http.StatusInternalServerError)
	}

	select {
	case msg := <-msgs:
		fmt.Fprint(rw, string(msg.Body))
		return
	case <-time.After(500 * time.Millisecond):
		rw.WriteHeader(http.StatusNotFound)
		return
	}
}

func declareQueue() (*amqp.Channel, *amqp.Queue, error) {
	conn, err := rabbit.Dial()
	if err != nil {
		return nil, nil, fmt.Errorf("Error opening AMQP connection: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("Error opening channel: %s", err)
	}

	q, err := ch.QueueDeclare(
		"testapp", // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		ch.Close()
		return nil, nil, fmt.Errorf("Error declaring queue: %s", err)
	}

	return ch, &q, nil
}
