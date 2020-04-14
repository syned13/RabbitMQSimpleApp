package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
	"github.com/subosito/gotenv"
)

var port string

type errorResponse struct {
	Message string `json:"message"`
}

type MessagePayload struct {
	Message string `json:"message"`
}

func init() {
	gotenv.Load()
	port = os.Getenv("PORT")
}

func main() {
	fmt.Println("Hello world")
	fmt.Println("PORT:" + port)

	conn, err := amqp.Dial(os.Getenv("CLOUDAMQP_URL"))

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	router := mux.NewRouter()

	router.HandleFunc("/", index).Methods(http.MethodGet)
	router.HandleFunc("/message", sendMessage(ch, q)).Methods(http.MethodPost)

	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	failOnError(err, "Failed to declare a queue")

	fmt.Println("Serving on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

func sendMessage(channel *amqp.Channel, queue amqp.Queue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		messagePayload := MessagePayload{}
		err := json.NewDecoder(r.Body).Decode(&messagePayload)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "invalid body")
		}

		fmt.Println("Sending message...")

		body := messagePayload.Message
		err = channel.Publish(
			"",         // exchange
			queue.Name, // routing key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})

		failOnError(err, "Failed to publish a message")
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// RespondWithError responds with a json with the given status code and message
func RespondWithError(w http.ResponseWriter, statusCode int, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse{errorMessage})
}
