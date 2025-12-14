package main

import (
	"io"
	"log"
	"net"
)

// Global variable to hold the active control connection
var ctrlConn net.Conn

func main() {
	// 1. Setup Control Listener (9000)
	ctrlListener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Control listener on :9000")

	// 2. Setup Data Listener (9002)
	dataListener, err := net.Listen("tcp", ":9002")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Data listener on :9002")

	// 3. Setup Public Listener (9001)
	pubListener, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Public listener on :9001. Ready for external traffic.")

	// --- HANDLE CONTROL CONNECTION IN A GOROUTINE ---
	// This allows the public listener loop to start immediately,
	// and the server to potentially re-accept the control connection if it closes.
	go func() {
		for {
			conn, err := ctrlListener.Accept()
			if err != nil {
				log.Println("Control listener accept error:", err)
				continue
			}
			log.Println("Client connected on control channel.")
			ctrlConn = conn // Set the global connection

			// Wait until the client closes this connection (or it errors)
			// A simple block is to read from it until EOF/error
			io.Copy(io.Discard, ctrlConn)

			// Connection closed/errored - reset global and loop to re-accept
			log.Println("Control connection closed. Waiting for client reconnection...")
			ctrlConn = nil
		}
	}()
	// ---------------------------------------------------

	// Main loop for accepting public connections
	for {
		externalConn, err := pubListener.Accept()
		if err != nil {
			log.Println("Error accepting external:", err)
			continue
		}

		// Check if the control connection is active before proceeding
		if ctrlConn == nil {
			log.Println("External request received but control channel is down. Rejecting.")
			externalConn.Close()
			continue
		}

		log.Println("External user connected. Requesting new data channel.")

		// Request a new data connection from the client
		_, err = ctrlConn.Write([]byte("new\n"))
		if err != nil {
			log.Println("Control write failed:", err)
			externalConn.Close()
			continue
		}

		// Wait for the client to dial back to the new data listener on :9002
		dataConn, err := dataListener.Accept()
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
