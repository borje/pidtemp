package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
)

type TempData struct {
	ActTime    int    `json:"ActTime"`
	ServerTime string `json:"ServerTime"`
	Sunrise    string `json:"Sunrise"`
	Sunset     string `json:"Sunset"`
	Result     []struct {
		AddjMulti         float64 `json:"AddjMulti"`
		AddjMulti2        float64 `json:"AddjMulti2"`
		AddjValue         float64 `json:"AddjValue"`
		AddjValue2        float64 `json:"AddjValue2"`
		BatteryLevel      int     `json:"BatteryLevel"`
		CustomImage       int     `json:"CustomImage"`
		Data              string  `json:"Data"`
		Description       string  `json:"Description"`
		DewPoint          string  `json:"DewPoint"`
		Favorite          int     `json:"Favorite"`
		HardwareID        int     `json:"HardwareID"`
		HardwareName      string  `json:"HardwareName"`
		HardwareType      string  `json:"HardwareType"`
		HardwareTypeVal   int     `json:"HardwareTypeVal"`
		HaveTimeout       bool    `json:"HaveTimeout"`
		Humidity          int     `json:"Humidity"`
		HumidityStatus    string  `json:"HumidityStatus"`
		ID                string  `json:"ID"`
		LastUpdate        string  `json:"LastUpdate"`
		Name              string  `json:"Name"`
		Notifications     string  `json:"Notifications"`
		PlanID            string  `json:"PlanID"`
		PlanIDs           []int   `json:"PlanIDs"`
		Protected         bool    `json:"Protected"`
		ShowNotifications bool    `json:"ShowNotifications"`
		SignalLevel       string  `json:"SignalLevel"`
		SubType           string  `json:"SubType"`
		Temp              float64 `json:"Temp"`
		Timers            string  `json:"Timers"`
		Type              string  `json:"Type"`
		TypeImg           string  `json:"TypeImg"`
		Unit              int     `json:"Unit"`
		Used              int     `json:"Used"`
		XOffset           string  `json:"XOffset"`
		YOffset           string  `json:"YOffset"`
		Idx               string  `json:"idx"`
	} `json:"result"`
	Status string `json:"status"`
	Title  string `json:"title"`
}

func getTemp() (error, float64) {
	u := url.URL{
		Scheme: "http",
		Host:   viper.GetString("domotics.hostname"),
		Path:   "json.htm",
	}
	q := u.Query()
	q.Set("type", "devices")
	q.Set("rid", viper.GetString("domotics.id"))
	u.RawQuery = q.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		log.Println(err)
		return err, 0
	}
	defer resp.Body.Close()
	log.Println(resp.Status)
	d := json.NewDecoder(resp.Body)
	var tempData TempData
	err = d.Decode(&tempData)
	if err != nil {
		log.Fatal(err)
		return err, 0
	}
	if l := len(tempData.Result); l == 1 {
		return nil, tempData.Result[0].Temp
	}

	return nil, 0
	//result := tempData.Result
	//for _, row := range result {
	//fmt.Println("temp:       ", row.Temp)
	//fmt.Println("Last update:", row.LastUpdate)
	//}
	//fmt.Println("temp: ", tempData.Result)
	//return 1

}
