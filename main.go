package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
		pid.SetPID(viper.GetFloat64("pid.p"), viper.GetFloat64("pid.i"), 0)
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

func main() {
	readConfigFile()
	pid = pidctrl.NewPIDController(20, .05, 0)
	pid.SetOutputLimits(0, 100)
	pid.Set(viper.GetFloat64("pid.target"))
	for {

		_, temp := getTemp()
		//output := pid.UpdateDuration(temp, time.Minute*10)
		output := pid.Update(temp)
		log.Println("PID output: ", output)
		log.Println("Temperature is: ", temp)

		time.Sleep(time.Second * 60)
		//d2 := time.Duration(viper.GetInt("pwm.period")) * time.Second
		//fmt.Println(d2)
		//time.Sleep(time.Duration(viper.GetInt("pwm.period")) * time.Second)
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
