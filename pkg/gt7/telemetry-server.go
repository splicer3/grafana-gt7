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
	playstationIP = "192.168.1.8"
)

func sendHeartBeat(conn *net.UDPConn) {
	heartbeatMsg := []byte("A")
	_, err := conn.Write(heartbeatMsg)
	if err != nil {
		log.DefaultLogger.Warn("SendHeartBeat", "Error sending heartbeat", err)
	}
}

func RunTelemetryServer(ch chan TelemetryFrame, errCh chan error) {
	heartbeatPort := fmt.Sprintf("%s:%s", playstationIP, heartbeatPort)
	sHeartbeat, err := net.ResolveUDPAddr("udp5", heartbeatPort)
	if err != nil {
		log.DefaultLogger.Warn("Heartbeat address resolution not working.")
		return
	}
	heartbeatConn, err := net.DialUDP("udp5", nil, sHeartbeat)
	if err != nil {
		log.DefaultLogger.Warn("Heartbeat not working.")
		return
	}
	sendHeartBeat(heartbeatConn)
	port := fmt.Sprintf("%s:%s", playstationIP, serverPort)
	s, err := net.ResolveUDPAddr("udp4", port)
	if err != nil {
		errCh <- err
		return
	}

	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		errCh <- err
		return
	}
	log.DefaultLogger.Info("Starting telemetry server for Gran Turismo 7")

	defer connection.Close()
	buffer := make([]byte, 4096)

	lastTimeSent := time.Now()

	for {
		n, _, err := connection.ReadFromUDP(buffer)
		if err != nil {
			errCh <- err
			return
		}

		if time.Now().After(lastTimeSent.Add(time.Second)) {
			sendHeartBeat(heartbeatConn)
		}

		fmt.Printf("Read %v bytes\n", n)

		packetBuffer := buffer[0:n]
		p, err := ReadPacket(packetBuffer)
		if err != nil {
			errCh <- err
			return
		}

		ch <- *p
	}
}
