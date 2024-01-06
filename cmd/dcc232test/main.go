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
	flag.StringVar(&cmdArgs.PortName, "port", "/dev/ttyS1", "Port name")
}

func main() {
	flag.Parse()

	c := &serial.Config{Name: cmdArgs.PortName, Baud: 19200}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}

	pins, err := NewPins()
	if err != nil {
		log.Fatal(err)
	}

	var p0, p1 dcc232.Packet
	p0 = p0.SpeedAndDirection(36, 0, true, dcc232.SpeedSteps128)
	p1 = p1.SpeedAndDirection(36, 45, true, dcc232.SpeedSteps128)
	ep0 := dcc232.EncodePacket(p0, nil)
	ep1 := dcc232.EncodePacket(p1, nil)

	start := time.Now()
	isEnabled, _ := pins.EnableStatus()
	for {
		sec := int64(time.Since(start).Seconds())
		msec := int64(time.Since(start).Seconds() * 10)
		ep := ep0
		if sec%2 == 0 {
			ep = ep1
		}
		if sec%10 == 0 {
			if err := pins.Enable(true); err != nil {
				log.Printf("Enable(true) failed: %s\n", err)
			}
		} else if sec%5 == 0 {
			if err := pins.Enable(false); err != nil {
				log.Printf("Enable(false) failed: %s\n", err)
			}
		}
		if msec%5 == 0 {
			if x, err := pins.EnableStatus(); err != nil {
				log.Printf("EnableStatus() failed: %s\n", err)
			} else {
				if x != isEnabled {
					isEnabled = x
					log.Printf("Enabled: %v\n", isEnabled)
				}
			}
		}

		if _, err := s.Write(ep); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Millisecond * 10)
	}
}
