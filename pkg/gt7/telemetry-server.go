package gt7

import (
	"fmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"net"
)

const (
	serverPort    = "33740"
	playstationIP = "192.168.1.8"
)

func RunTelemetryServer(ch chan TelemetryFrame, errCh chan error) {
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

	for {
		n, _, err := connection.ReadFromUDP(buffer)
		if err != nil {
			errCh <- err
			return
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
