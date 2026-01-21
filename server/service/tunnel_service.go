package service

import (
	"io"
	"log"
	"net"
	"sync"

	"opentunnel/server/config"
	"opentunnel/server/constants"
)

type TunnelService struct {
	cfg          *config.Config
	ctrlConn     net.Conn
	ctrlMux      sync.Mutex
	ctrlListener net.Listener
	dataListener net.Listener
	pubListener  net.Listener
}

func NewTunnelService(cfg *config.Config) *TunnelService {
	return &TunnelService{cfg: cfg}
}

func (s *TunnelService) Start() error {
	var err error

	s.ctrlListener, err = net.Listen("tcp", ":"+s.cfg.ControlPort)
	if err != nil {
		return err
	}
	log.Printf("Control listener on :%s", s.cfg.ControlPort)

	s.dataListener, err = net.Listen("tcp", ":"+s.cfg.DataPort)
	if err != nil {
		return err
	}
	log.Printf("Data listener on :%s", s.cfg.DataPort)

	s.pubListener, err = net.Listen("tcp", ":"+s.cfg.ExternalPort)
	if err != nil {
		return err
	}
	log.Printf("Public listener on :%s. Ready for external traffic.", s.cfg.ExternalPort)

	go s.acceptControlConnections()
	s.acceptPublicConnections()

	return nil
}

func (s *TunnelService) acceptControlConnections() {
	for {
		conn, err := s.ctrlListener.Accept()
		if err != nil {
			log.Println(constants.ErrControlAccept+":", err)
			continue
		}
		log.Println("Client connected on control channel.")

		s.ctrlMux.Lock()
		s.ctrlConn = conn
		s.ctrlMux.Unlock()

		io.Copy(io.Discard, conn)

		log.Println("Control connection closed. Waiting for client reconnection...")

		s.ctrlMux.Lock()
		s.ctrlConn = nil
		s.ctrlMux.Unlock()
	}
}

func (s *TunnelService) acceptPublicConnections() {
	for {
		externalConn, err := s.pubListener.Accept()
		if err != nil {
			log.Println(constants.ErrExternalAccept+":", err)
			continue
		}

		s.ctrlMux.Lock()
		ctrl := s.ctrlConn
		s.ctrlMux.Unlock()

		if ctrl == nil {
			log.Println(constants.ErrControlDown + ". Rejecting.")
			externalConn.Close()
			continue
		}

		log.Println("External user connected. Requesting new data channel.")

		_, err = ctrl.Write([]byte("new\n"))
		if err != nil {
			log.Println(constants.ErrControlWrite+":", err)
			externalConn.Close()
			continue
		}

		dataConn, err := s.dataListener.Accept()
		if err != nil {
			log.Println(constants.ErrDataConnection+":", err)
			externalConn.Close()
			continue
		}

		go s.proxyConnections(dataConn, externalConn)
	}
}

func (s *TunnelService) proxyConnections(dataConn, externalConn net.Conn) {
	defer dataConn.Close()
	defer externalConn.Close()

	log.Println("Starting proxy pipe...")
	go io.Copy(dataConn, externalConn)
	io.Copy(externalConn, dataConn)
}
