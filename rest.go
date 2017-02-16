package main

import (
	"log"
	"net/http"

	"github.com/unrolled/render"
	//"github.com/unrolled/render"
)

func initRest(p *Pwm) {
	r := render.New(render.Options{})
	mux := http.NewServeMux()

	mux.HandleFunc("/dutycycle", func(w http.ResponseWriter, req *http.Request) {
		dutycyclePercentage := p.GetDutyCycle() * 100.0
		log.Println("/dutycycle was delivered ", dutycyclePercentage)
		r.JSON(w, http.StatusOK, dutycyclePercentage)
	})

	http.ListenAndServe(":8081", mux)
}
