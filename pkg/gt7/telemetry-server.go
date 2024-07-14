package gt7

import (
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"net"
	"time"
)

const (
	heartbeatPort = "33739"
	serverPort    = "33740"
)

func sendHeartBeat(conn *net.UDPConn) {
	heartbeatMsg := []byte("A")
	_, err := conn.Write(heartbeatMsg)
	if err != nil {
		log.DefaultLogger.Warn("SendHeartBeat", "Error sending heartbeat", err)
	}
}

func RunTelemetryServer(playstationIP string, ch chan TelemetryFrame, errCh chan error, hbChan chan *net.UDPConn, strChan chan *net.UDPConn) {
	// Heartbeat connection setup
	if playstationIP == "" {
		playstationIP = "192.168.1.5"
	}
	heartbeatAddr, err := net.ResolveUDPAddr("udp4", net.JoinHostPort(playstationIP, heartbeatPort))
	if err != nil {
		errCh <- fmt.Errorf("heartbeat address resolution failed: %v", err)
		return
	}

	heartbeatConn, err := net.DialUDP("udp4", nil, heartbeatAddr)
	if err != nil {
		errCh <- fmt.Errorf("heartbeat connection failed: %v", err)
		return
	}
	defer heartbeatConn.Close()
	hbChan <- heartbeatConn

	// Server connection setup
	serverAddr, err := net.ResolveUDPAddr("udp4", net.JoinHostPort("", serverPort))
	if err != nil {
		errCh <- fmt.Errorf("server address resolution failed: %v", err)
		return
	}

	serverConn, err := net.ListenUDP("udp4", serverAddr)
	if err != nil {
		errCh <- fmt.Errorf("server listener failed: %v", err)
		return
	}
	defer serverConn.Close()
	strChan <- serverConn

	log.DefaultLogger.Info("Starting telemetry server for Gran Turismo 7", "PlaystationIP", playstationIP)

	buffer := make([]byte, 4096)
	lastHeartbeatTime := time.Now()

	for {
		// Set read deadline to prevent blocking indefinitely
		err = serverConn.SetReadDeadline(time.Now().Add(time.Second))
		if err != nil {
			log.DefaultLogger.Warn("SetReadDeadline failed", "err", err.Error())
			return
		}

		n, _, err := serverConn.ReadFromUDP(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// It's just a timeout, continue and send heartbeat if needed
			} else {
				log.DefaultLogger.Warn("ReadFromUDP failed", "err", err.Error())
				return
			}
		}

		if n > 0 {
			fmt.Printf("Read %v bytes\n", n)
			packetBuffer := buffer[0:n]
			p, err := ReadPacket(packetBuffer)
			if err != nil {
				log.DefaultLogger.Warn("ReadPacket failed", "err", err.Error())
				return
			}
			ch <- *p
		}

		// Send heartbeat if a second has passed since the last one
		if time.Since(lastHeartbeatTime) >= time.Second {
			sendHeartBeat(heartbeatConn)
			lastHeartbeatTime = time.Now()
		}
	}
}
