package pidtemp

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
	fmt.Println("starting test")
	var fakeSwitch FakeSwitch
	pwm := NewPwm(&fakeSwitch)
	pwm.SetPeriod(10 * time.Second)
	pwm.SetDutyCycle(.5)
	start := time.Now()
	pwm.Start()
	time.Sleep(time.Second * 10)
	fmt.Println("Stopping")
	pwm.Stop()
	fmt.Println("Stopped")
	duration := time.Since(start)
	ratio := float64(fakeSwitch.DurationNs) / float64(duration.Nanoseconds())
	fmt.Println("On Time   : ", fakeSwitch.DurationNs)
	fmt.Println("Total Time: ", duration.Nanoseconds())
	fmt.Println("Ratio: ", ratio)
	if ratio < .45 || ratio > .55 {
		t.Error("Wrong ratio")
	}
}

func TestPeriodChange(t *testing.T) {
	//change period during "on time" and the ratio will be errorneous
	fmt.Println("starting test")
	var fakeSwitch FakeSwitch
	pwm := NewPwm(&fakeSwitch)
	pwm.SetPeriod(2 * time.Second)
	pwm.SetDutyCycle(.5)
	start := time.Now()
	pwm.Start()
	time.Sleep(time.Second * 3)
	pwm.SetPeriod(10 * time.Second)
	time.Sleep(time.Second * 11)
	fmt.Println("Stopping")
	pwm.Stop()
	fmt.Println("Stopped")
	duration := time.Since(start)
	ratio := float64(fakeSwitch.DurationNs) / float64(duration.Nanoseconds())
	fmt.Println("On Time   : ", fakeSwitch.DurationNs)
	fmt.Println("Total Time: ", duration.Nanoseconds())
	fmt.Println("Ratio: ", ratio)

}
