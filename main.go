package main

import (
	"expvar"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

var requestCounter = expvar.NewInt("requests")
var connectionCounter = expvar.NewInt("connections")

func main() {
	fmt.Printf("Listening on *:8080\n")
	fmt.Printf("URLs:\n")
	fmt.Printf("\t/\tLoad test target\n")
	fmt.Printf("\t/get\tGet connection/request counters\n")
	fmt.Printf("\t/reset\tReset connection/request counters\n")
	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/reset", handleReset)
	http.HandleFunc("/get", handleGet)
	server := http.Server{
		Addr:      ":8080",
		ConnState: connectionChanged,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func handleGet(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Requests: %v\r\n", requestCounter.Value())
	fmt.Fprintf(writer, "Connections: %v\r\n", connectionCounter.Value())
}

func handleReset(writer http.ResponseWriter, request *http.Request) {
	requestCounter.Set(0)
	connectionCounter.Set(0)
	io.WriteString(writer, "OK\r\n")
}

func handleRequest(writer http.ResponseWriter, request *http.Request) {
	requestCounter.Add(1)
	io.WriteString(writer, "OK\r\n")
}

func connectionChanged(conn net.Conn, state http.ConnState) {
	if state == http.StateNew {
		connectionCounter.Add(1)
	}
}
