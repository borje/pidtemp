package main

import "log"

type DomoSwitch struct {
	Id   int
	Host string
}

func (s *DomoSwitch) On() {
	log.Println("Switch ON")
}

func (s *DomoSwitch) Off() {
	log.Println("Switch OFF")

}
