package main

import (
	"github.com/ecc1/gpio"
)

type Pins struct {
	enReq    gpio.OutputPin
	enStatus gpio.InputPin
}

const (
	enPin    = 1
	enReqPin = 0
)

func NewPins() (*Pins, error) {
	activeLow := false
	initialValue := false
	enReq, err := gpio.Output(enReqPin, activeLow, initialValue)
	if err != nil {
		return nil, err
	}
	enStatus, err := gpio.Input(enPin, activeLow)
	if err != nil {
		return nil, err
	}
	return &Pins{
		enReq:    enReq,
		enStatus: enStatus,
	}, nil
}

func (p *Pins) EnableStatus() (bool, error) {
	return p.enStatus.Read()
}

func (p *Pins) Enable(value bool) error {
	return p.enReq.Write(value)
}
