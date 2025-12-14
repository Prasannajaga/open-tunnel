package main

import (
	"io"
	"log"
	"net"
)

var ctrlConn net.Conn

func main() {

	ctrlListener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Control listener on :9000")

	dataListener, err := net.Listen("tcp", ":9002")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Data listener on :9002")

	pubListener, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Public listener on :9001. Ready for external traffic.")

	go func() {
		for {
			conn, err := ctrlListener.Accept()
			if err != nil {
				log.Println("Control listener accept error:", err)
				continue
			}
			log.Println("Client connected on control channel.")
			ctrlConn = conn

			// Wait until the client closes this connection (or it errors)
			// A simple block is to read from it until EOF/error
			io.Copy(io.Discard, ctrlConn)

			log.Println("Control connection closed. Waiting for client reconnection...")
			ctrlConn = nil
		}
	}()

	for {
		externalConn, err := pubListener.Accept()
		if err != nil {
			log.Println("Error accepting external:", err)
			continue
		}

		if ctrlConn == nil {
			log.Println("External request received but control channel is down. Rejecting.")
			externalConn.Close()
			continue
		}

		log.Println("External user connected. Requesting new data channel.")

		_, err = ctrlConn.Write([]byte("new\n"))
		if err != nil {
			log.Println("Control write failed:", err)
			externalConn.Close()
			continue
		}

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
