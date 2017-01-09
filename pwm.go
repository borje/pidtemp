package main

import (
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
	p.period = period
}

func (p *Pwm) SetDutyCycle(dc float64) {
	if dc > 1 {
		dc = 1
	}
	p.dutyCycle = dc
}

func (p *Pwm) Start() {
	go func() {
		isRunning := true
		p.stopWg.Add(1)
		p.runLock.Lock()
		p.running = isRunning
		p.runLock.Unlock()
		for isRunning {
			onTime := float64(p.period.Nanoseconds()) * p.dutyCycle
			p.sw.On()
			time.Sleep(time.Duration(onTime))

			p.sw.Off()
			offTime := p.period - time.Duration(onTime)
			time.Sleep(offTime)

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
