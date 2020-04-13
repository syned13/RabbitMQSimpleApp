package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

var port string

func init() {
	gotenv.Load()
	port = os.Getenv("PORT")
}

func main() {
	fmt.Println("Hello world")

	router := mux.NewRouter()

	router.HandleFunc("/", index).Methods(http.MethodGet)
	router.HandleFunc("/message", sendMessage).Methods(http.MethodPost)

	log.Fatal(http.ListenAndServe(port, router))
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Print("sending message...")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
