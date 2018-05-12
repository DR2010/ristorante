package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

var clientlist []string

// Echo is for something
func Echo(ws *websocket.Conn) {
	var errreceive error
	var errsend error

	fmt.Println("Start of Echo func.")

	for {
		var reply string

		errreceive = websocket.Message.Receive(ws, &reply)

		if errreceive != nil {
			fmt.Println("Can't receive")
			break
		}

		// if reply is not empty, message has been received
		fmt.Println("Received back from client: " + reply)

		msg := reply
		fmt.Println("Sending to client: " + msg)

		errsend = websocket.Message.Send(ws, msg)

		if errsend != nil {
			fmt.Println("Can't send")
			break
		}
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {
	http.HandleFunc("/", serveHome)
	http.Handle("/ws", websocket.Handler(Echo))
	fmt.Println("ListenAndServe:1234")

	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
