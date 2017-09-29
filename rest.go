package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	//"github.com/unrolled/render"
	"github.com/felixge/pidctrl"
)

func initRest(p *Pwm, pid *pidctrl.PIDController) {
	r := render.New(render.Options{})
	router := mux.NewRouter()

	router.HandleFunc("/dutycycle", func(w http.ResponseWriter, req *http.Request) {
		dutycyclePercentage := p.GetDutyCycle() * 100.0
		log.Println("/dutycycle was delivered ", dutycyclePercentage)
		r.JSON(w, http.StatusOK, dutycyclePercentage)
	})

	router.HandleFunc("/target/{value}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)

		log.Println("received REST call to set Target temp to", vars["value"])
	})

	http.ListenAndServe(":8081", router)
}
