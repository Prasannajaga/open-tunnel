package main

import (
	"io"
	"log"
	"net"
)

func main() {

	clientListener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Waiting for client on :9000...")

	clientConn, err := clientListener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Client connected")

	// Accept external users
	publicListener, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Public listener on :9001")

	for {
		externalConn, err := publicListener.Accept()
		if err != nil {
			log.Println("Public accept error:", err)
			continue
		}
		log.Println("External user connected")

		// Forward external <-> client
		go io.Copy(clientConn, externalConn)
		go io.Copy(externalConn, clientConn)
	}
}
