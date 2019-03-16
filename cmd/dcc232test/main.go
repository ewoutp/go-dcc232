package main

import (
	"flag"
	"log"
	"time"

	"github.com/ewoutp/go-dcc232"
	"github.com/tarm/serial"
)

var (
	cmdArgs struct {
		PortName string
	}
)

func init() {
	flag.StringVar(&cmdArgs.PortName, "port", "", "Port name")
}

func main() {
	c := &serial.Config{Name: cmdArgs.PortName, Baud: 19200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	p := dcc232.SpeedAndDirection(36, 45, true, dcc232.SpeedSteps128)
	ep := dcc232.EncodePacket(p)

	for {
		if _, err := s.Write(ep); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Millisecond * 10)
	}
}
