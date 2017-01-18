package main

import (
	"log"
	"sync"
	"time"
)

type Pwm struct {
	period    time.Duration
	dutyCycle float64
	sw        Switch
	runLock   sync.RWMutex
	running   bool
	stopWg    sync.WaitGroup
}

type Switch interface {
	On()
	Off()
}

func NewPwm(s Switch) *Pwm {
	return &Pwm{sw: s}
}

func (p *Pwm) SetPeriod(period time.Duration) {
	log.Println("Setting period: ", period)
	p.period = period
}

func (p *Pwm) SetDutyCycle(dc float64) {
	if dc > 1 {
		dc = 1
	}
	p.dutyCycle = dc
}

func (p *Pwm) Start() {
	if p.period == 0 {
		log.Fatal("PWM period is not set")
	}
	go func() {
		isRunning := true
		p.stopWg.Add(1)
		p.runLock.Lock()
		p.running = isRunning
		p.runLock.Unlock()
		for isRunning {
			period := p.period
			onTime := float64(period.Nanoseconds()) * p.dutyCycle
			if onTime > 0 {
				p.sw.On()
				log.Println("PWM: On time is ", time.Duration(onTime))
				time.Sleep(time.Duration(onTime))
			}

			offTime := period - time.Duration(onTime)
			if offTime > 0 {
				p.sw.Off()
				log.Println("PWM: Off time is ", time.Duration(offTime))
				time.Sleep(offTime)
			}

			p.runLock.RLock()
			isRunning = p.running
			p.runLock.RUnlock()
		}
		p.stopWg.Done()
	}()
}

func (p *Pwm) Stop() {
	p.runLock.Lock()
	p.running = false
	p.runLock.Unlock()
	p.stopWg.Wait()
}
