package main

import (
	"io"
	"log"
	"net"
)

func main() {
	// Control channel: client connects here
	ctrlListener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Waiting for client control connection on :9000...")

	ctrlConn, err := ctrlListener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Client connected on control channel")

	// Public listener
	pubListener, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Public listener on :9001")

	for {
		externalConn, err := pubListener.Accept()
		if err != nil {
			log.Println("Error accepting external:", err)
			continue
		}
		log.Println("External user connected")

		// Request a new data connection from the client
		_, err = ctrlConn.Write([]byte("new\n"))
		if err != nil {
			log.Println("Control write failed:", err)
			externalConn.Close()
			continue
		}

		// The client will call net.Dial back to the server on :9002
		dataConn, err := ctrlListener.Accept()
		if err != nil {
			log.Println("Data connection failed:", err)
			externalConn.Close()
			continue
		}

		go func() {
			defer dataConn.Close()
			defer externalConn.Close()

			go io.Copy(dataConn, externalConn)
			io.Copy(externalConn, dataConn)
		}()
	}
}
