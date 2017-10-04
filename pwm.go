package pidtemp

import (
	"log"
	"sync"
	"time"
)

type Pwm struct {
	period    time.Duration
	dutyCycle float64
	sw        Switch
	rwLock    sync.RWMutex
	quit      chan int
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
	p.rwLock.Lock()
	p.period = period
	p.rwLock.Unlock()
}

func (p *Pwm) SetDutyCycle(dc float64) {
	if dc > 1 {
		dc = 1
	}
	p.rwLock.Lock()
	p.dutyCycle = dc
	p.rwLock.Unlock()
}

func (p *Pwm) GetDutyCycle() float64 {
	p.rwLock.RLock()
	defer p.rwLock.RUnlock()
	return p.dutyCycle
}

func (p *Pwm) createOffTimer() *time.Timer {
	onTime := time.Duration(float64(p.period.Nanoseconds()) * p.dutyCycle)
	if onTime > 0 {
		p.sw.On()
		log.Println("PWM: On time is ", onTime)
		if onTime < p.period {
			return time.NewTimer(onTime)
		}
	}
	timer := time.NewTimer(time.Minute)
	timer.Stop()
	return timer
}

func (p *Pwm) Start() {
	if p.period == 0 {
		log.Fatal("PWM period is not set")
	}

	go func() {
		p.quit = make(chan int)
		period := p.period
		cycleTick := time.NewTicker(period)
		offTimer := p.createOffTimer()

		for {
			select {
			case <-cycleTick.C:
				if p.dutyCycle == 0 {
					p.sw.Off()
				}
				if p.period != period {
					log.Println("New period set in PWM")
					period = p.period
					cycleTick = time.NewTicker(period)
				}
				offTimer = p.createOffTimer()
			case <-offTimer.C:
				p.sw.Off()
			case <-p.quit:
				cycleTick.Stop()
				// Decide if the status should be off or on
				return

			}
		}
	}()
}

func (p *Pwm) Stop() {
	close(p.quit)
}
