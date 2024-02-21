package main

import (
	"errors"
	"fmt"
	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
	"log"
	"time"
)

func getPorts() ([]*enumerator.PortDetails, error) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return nil, err
	}
	if len(ports) == 0 {
		return nil, errors.New("no ports found")
	}
	return ports, nil
}
func main() {
	for {
		ports, err := getPorts()
		if err != nil {
			continue
		}
		var targetPort string
		for _, port := range ports {
			//fmt.Printf("Found port: %s\n", port.Name)
			if port.IsUSB {
				//	fmt.Printf("   USB ID     %s:%s\n", port.VID, port.PID)
				//	fmt.Printf("   USB serial %s\n", port.SerialNumber)
				if port.VID == "2e8a" && port.PID == "000a" {
					targetPort = port.Name
					break
				}
			}
		}
		if targetPort == "" {
			continue
		}
		// Open the first serial port detected at 9600bps N81
		mode := &serial.Mode{
			BaudRate: 115200,
			Parity:   serial.NoParity,
			DataBits: 8,
			StopBits: serial.OneStopBit,
		}
		port, err := serial.Open(targetPort, mode)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		// Read and print the response

		buff := make([]byte, 100)
		for {
			// Reads up to 100 bytes
			n, err := port.Read(buff)
			if err != nil {
				log.Println(err)
				break
			}
			if n == 0 {
				fmt.Println("\nEOF")
				break
			}

			fmt.Printf("%s", string(buff[:n]))
		}
	}
}
