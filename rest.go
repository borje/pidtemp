package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/felixge/pidctrl"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/unrolled/render"
)

func initRest(pwm *Pwm, pid *pidctrl.PIDController) {
	r := render.New(render.Options{IndentJSON: true})
	router := mux.NewRouter()

	router.HandleFunc("/dutycycle", func(w http.ResponseWriter, req *http.Request) {
		dutycyclePercentage := pwm.GetDutyCycle() * 100.0
		log.Println("/dutycycle was delivered ", dutycyclePercentage)
		r.JSON(w, http.StatusOK, dutycyclePercentage)
	})

	router.HandleFunc("/target/{value}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		if f, err := strconv.ParseFloat(vars["value"], 64); err == nil {
			log.Println("received REST call to set Target temp to", f)
			viper.Set("pid.target", f)
			pid.Set(f)
		} else {
			log.Println("Failed to parse given Target temp as float", vars["value"])
		}
	})

	router.HandleFunc("/config", func(w http.ResponseWriter, req *http.Request) {
		p, i, d := pid.PID()
		values := map[string]float64{}
		values["p"], values["i"], values["d"] = p, i, d
		values["target"] = viper.GetFloat64("pid.target")
		values["dutycycle"] = pwm.GetDutyCycle()
		r.JSON(w, http.StatusOK, values)
		log.Println(p, i, d)
	})

	http.ListenAndServe(":8081", router)
}
