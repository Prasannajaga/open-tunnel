package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net"
	"strconv"
)

const serverIP = "34.133.55.212"
const controlPort = "9000"
const externalPort = "9001"
const dataPort = "9002"
const defaultPort = 8080

var (
	localPort  int
	localAddr  string
	serverAddr string
	dataAddr   string
)

func init() {

	flag.IntVar(&localPort, "port", 8080, "The local port to expose (e.g., 8080)")
	flag.Parse()

	if len(flag.Args()) > 0 {
		p, err := strconv.Atoi(flag.Args()[0])
		if err == nil {
			localPort = p
		} else {
			log.Fatalf("Invalid port argument: %s", flag.Args()[0])
		}
	} else {
		localPort = defaultPort
	}

	localAddr = "localhost:" + strconv.Itoa(localPort)
	serverAddr = serverIP + ":" + controlPort
	dataAddr = serverIP + ":" + dataPort
}

func main() {

	ctrlConn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatalf("Oops! The tunnel server at %s is not responding or unreachable. Please check the server status and IP address (%s). Original error: %v", serverAddr, serverIP, err)
	}

	exposedURL := "http://" + serverIP + ":" + externalPort
	log.Printf("Success! Tunnel client connected. Your local service is now publicly exposed at: %s", exposedURL)
	log.Println("Listening for incoming connections...")

	scanner := bufio.NewScanner(ctrlConn)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "new" {
			// New incoming connection → open new local + new server data channel
			go handleProxy()
		}
	}
	log.Println("Control connection closed by server or client. Exiting.")
}

func handleProxy() {
	dataConn, err := net.Dial("tcp", dataAddr)
	if err != nil {
		log.Println("Failed to open data channel to server:", err)
		return
	}

	localConn, err := net.Dial("tcp", localAddr)
	if err != nil {
		log.Printf("Could not reach your local service at %s. Please ensure the application is running and listening on port %d. Error: %v", localAddr, localPort, err)
		dataConn.Close()
		return
	}
	log.Printf("Connected to local service %s. Starting proxying traffic...", localAddr)

	defer dataConn.Close()
	defer localConn.Close()

	go func() {
		_, err := io.Copy(dataConn, localConn)
		if err != nil && err != io.EOF {
			log.Printf("Data pipe error: %v", err)
		}
	}()

	_, err = io.Copy(localConn, dataConn)
	if err != nil && err != io.EOF {
		log.Printf("Data pipe error: %v", err)
	}

	log.Println("Data pipe closed for this connection.")
}
