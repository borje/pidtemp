package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/felixge/pidctrl"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type TempHistory struct {
	tempHistory []struct {
		D  string  `json:"d"`
		Hu string  `json:"hu"`
		Te float64 `json:"te"`
	} `json:"result"`
	Status string `json:"status"`
	Title  string `json:"title"`
}

func readConfigFile() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	// Watch for modification
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		pid.Set(viper.GetFloat64("pid.target"))
		pid.SetPID(viper.GetFloat64("pid.p"), viper.GetFloat64("pid.i"), viper.GetFloat64("pid.d"))
		pwm.SetPeriod(time.Duration(viper.GetInt("pwm.period")) * time.Second)
		//duration, err := time.ParseDuration(viper.GetString("pwm.period"))
		//if err == nil {
		//pwm.SetPeriod(duration)
		//} else {
		//log.Println("error in pwm.period in config file")
		//}
		log.Println("Config changed")
	})

	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
		log.Fatal("Failed to load config file.")
	}

	i := viper.GetInt("domotics.id")
	hostname := viper.GetString("domotics.hostname")
	fmt.Println("hostname: ", hostname)
	fmt.Println("id: ", i)
}

var pid *pidctrl.PIDController
var pwm *Pwm

func roundOutput(output float64) float64 {
	if output < 0.1 {
		return 0
	} else if output > 0.9 {
		return 1
	} else {
		return output
	}

}

func main() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime | log.Lmicroseconds)
	f, err := os.OpenFile("pidtemp.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("error opening file: ", err)
		return
	}
	log.SetOutput(io.MultiWriter(os.Stderr, f))
	defer f.Close()

	readConfigFile()

	// Initiate PID Controller
	pid = pidctrl.NewPIDController(viper.GetFloat64("pid.p"), viper.GetFloat64("pid.i"), viper.GetFloat64("pid.d"))
	pid.SetOutputLimits(0, 100)
	pid.Set(viper.GetFloat64("pid.target"))

	// Initiate PWM
	powerSwitch := &DomoSwitch{Host: viper.GetString("domotics.hostname"), Id: 9000}
	pwm = NewPwm(powerSwitch)
	pwm.SetPeriod(time.Duration(viper.GetInt("pwm.period")) * time.Second)
	var temp float64
	err, temp = getTemp()
	// Keep try until we get a temperature
	for err != nil {
		log.Println(err)
		log.Println("Sleep and retry")
		time.Sleep(time.Minute)
		err, temp = getTemp()
	}

	pwm.SetDutyCycle(pid.Update(temp) / 100.0) // Set dutycycle for firsrt cycle
	pwm.Start()
	go initRest(pwm)
	defer pwm.Stop()

	for {

		err, temp := getTemp()
		// if cant get temperature, we should perhaps keep the same output
		for err != nil {
			log.Println(err)
			log.Println("Sleep and retry")
			time.Sleep(time.Minute)
			err, temp = getTemp()
		}

		output := pid.Update(temp) / 100.0
		// if the error is small and output is mostly constant,
		// largen the PWM period. multiply with 1.5 up to 30 minutes
		// the error gets to large, go down to configures period.

		roundedOutput := roundOutput(output)
		pwm.SetDutyCycle(roundedOutput)
		log.Println("Temperature is: ", temp)
		log.Println("PID output: ", output)

		// This will get off sync from the PWM, so this could be a problem
		//sleepTime := time.Second * 60)
		sleepTime := time.Duration(viper.GetInt("pwm.period")) * time.Second
		log.Println("Sleeping for ", sleepTime)
		time.Sleep(time.Duration(viper.GetInt("pwm.period")) * time.Second)
	}
}

func getTempHistory() {
	resp, err := http.Get("http://pilight:8080/json.htm?type=graph&sensor=temp&idx=62&range=day")
	if err != nil {
		fmt.Println(err)
	}

	var tempHistory TempHistory
	d := json.NewDecoder(resp.Body)
	err = d.Decode(&tempHistory)
	l := len(tempHistory.tempHistory)
	fmt.Println("size: ", l)
	//for _, v := range tempdata.Result {
	//fmt.Println(v.D, v.Te)
	//}
	last := &tempHistory.tempHistory[l-1]
	fmt.Println(last.D, last.Te)
	resp.Body.Close()
}
