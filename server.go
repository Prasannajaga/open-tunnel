package main

import (
	"io"
	"log"
	"net"
)

func main() {
	// 1. Control listener: client connects here for persistent commands
	ctrlListener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Waiting for client control connection on :9000...")

	// Accept the single control connection
	ctrlConn, err := ctrlListener.Accept()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Client connected on control channel")
	// Note: Close the listener if you only want ONE client, otherwise keep it for multiple clients/tunnels

	// 2. Data listener: Client connects back to this for data channel
	dataListener, err := net.Listen("tcp", ":9002") // *** NEW LISTENER ***
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Data listener on :9002")

	// 3. Public listener: External users connect here
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

		// Wait for the client to dial back to the new data listener on :9002
		dataConn, err := dataListener.Accept() // *** USE NEW LISTENER ***
		if err != nil {
			log.Println("Data connection failed:", err)
			externalConn.Close()
			continue
		}

		go func() {
			defer dataConn.Close()
			defer externalConn.Close()

			log.Println("Starting proxy pipe...")
			go io.Copy(dataConn, externalConn)
			io.Copy(externalConn, dataConn)
		}()
	}
}
