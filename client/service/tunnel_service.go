package service

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"

	"opentunnel/client/config"
	"opentunnel/client/constants"
	"opentunnel/client/utils"
)

type TunnelService struct {
	cfg       *config.Config
	ctrlConn  net.Conn
	localAddr string
}

func NewTunnelService(cfg *config.Config) *TunnelService {
	return &TunnelService{
		cfg:       cfg,
		localAddr: utils.BuildLocalAddress(cfg.LocalPort),
	}
}

func (s *TunnelService) Connect() error {
	conn, err := net.Dial("tcp", s.cfg.ServerAddr())
	if err != nil {
		return fmt.Errorf(constants.ErrServerUnreachable, s.cfg.ServerAddr())
	}
	s.ctrlConn = conn
	return nil
}

func (s *TunnelService) Listen() {
	log.Printf("Success! Tunnel client connected. Your local service is now publicly exposed at: %s", s.cfg.ExposedURL())
	log.Println("Listening for incoming connections...")

	scanner := bufio.NewScanner(s.ctrlConn)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "new" {
			go s.handleProxy()
		}
	}
	log.Println("Control connection closed by server or client. Exiting.")
}

func (s *TunnelService) handleProxy() {
	dataConn, err := net.Dial("tcp", s.cfg.DataAddr())
	if err != nil {
		log.Println(constants.ErrDataChannelFailed+":", err)
		return
	}

	localConn, err := net.Dial("tcp", s.localAddr)
	if err != nil {
		log.Printf(constants.ErrLocalServiceDown+". Error: %v", s.localAddr, s.cfg.LocalPort, err)
		dataConn.Close()
		return
	}
	log.Printf("Connected to local service %s. Starting proxying traffic...", s.localAddr)

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
