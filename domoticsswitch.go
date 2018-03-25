package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
)

type DomoSwitch struct {
	Id   int
	Host string
}

func createSwitchCmd(cmd string) url.URL {
	u := url.URL{Scheme: "http",
		Host: viper.GetString("domotics.hostname"),
		Path: "json.htm",
	}
	q := u.Query()
	q.Set("type", "command")
	q.Set("param", "switchlight")
	q.Set("idx", "105")
	q.Set("switchcmd", cmd)
	u.RawQuery = q.Encode()
	return u
}

func (s *DomoSwitch) On() {
	u := createSwitchCmd("On")
	resp, err := http.Get(u.String())
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	log.Println("Switch ON")
}

func (s *DomoSwitch) Off() {
	u := createSwitchCmd("Off")
	resp, err := http.Get(u.String())
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	log.Println("Switch OFF")

}
