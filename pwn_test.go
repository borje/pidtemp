package main

import (
	"fmt"
	"testing"
	"time"
)

type FakeSwitch struct {
	DurationNs int64
	startTime  time.Time
}

func (s *FakeSwitch) On() {
	//fmt.Println("On")
	s.startTime = time.Now()
}
func (s *FakeSwitch) Off() {
	//fmt.Println("Off")
	diff := time.Since(s.startTime)
	fmt.Println("Adding: ", diff)
	s.DurationNs += diff.Nanoseconds()
}

func TestFoo(t *testing.T) {
	var fakeSwitch FakeSwitch
	pwm := NewPwm(&fakeSwitch)
	pwm.SetPeriod(10 * time.Second)
	pwm.SetDutyCycle(.5)
	start := time.Now()
	pwm.Start()
	time.Sleep(time.Second * 3)
	fmt.Println("Stopping")
	pwm.Stop()
	fmt.Println("Stopped")
	duration := time.Since(start)
	ratio := float64(fakeSwitch.DurationNs) / float64(duration.Nanoseconds())
	fmt.Println("On Time   : ", fakeSwitch.DurationNs)
	fmt.Println("Total Time: ", duration.Nanoseconds())
	fmt.Println("Ratio: ", ratio)
}

func TestPeriodChange() {
	//change period during "on time" and the ratio will be errorneous
}
