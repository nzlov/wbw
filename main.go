package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
)

var (
	lat  = flag.String("lat", "", "latitude")
	lon  = flag.String("lon", "", "longitude")
	lang = flag.String("lang", "zh", "longitude")
)

func main() {
	flag.Parse()
	if *lat == "" {
		panic("missing latitude")
	}
	if *lon == "" {
		panic("missing longitude")
	}

	resp, err := http.Get(fmt.Sprintf("https://weather.sl.al/?lat=%s&lon=%s&lang=%s", *lat, *lon, *lang))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	w := Weather{}
	if err := json.Unmarshal(data, &w); err != nil {
		panic(err)
	}

	hs := ``
	for _, v := range w.Hourly {
		hs += fmt.Sprintf("%s %d°C %s\r", v.Time, v.Temperature, v.RainProbability)
	}

	ds := ``
	for _, v := range w.Daily {
		ds += fmt.Sprintf("%s %d-%d°C %s\r", v.Date, v.Low, v.High, v.RainProbability)
	}

	m := map[string]any{
		"text":    fmt.Sprintf("%s %d°C", w.Current.Description, w.Current.Temperature),
		"tooltip": fmt.Sprintf(`%s空气质量:%s%s%s`, w.Location.City+"\r", w.Current.AirQuality.Category+"\r", hs, ds),
	}
	data, err = json.Marshal(&m)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
}

type Weather struct {
	Lat      string   `json:"lat"`
	Lon      string   `json:"lon"`
	Alert    string   `json:"alert"`
	Location Location `json:"location"`
	Current  Current  `json:"current"`
	Sun      Sun      `json:"sun"`
	Hourly   []Hourly `json:"hourly"`
	Daily    []Daily  `json:"daily"`
}

type Current struct {
	Temperature int64      `json:"temperature"`
	FeelsLike   int64      `json:"feelsLike"`
	Description string     `json:"description"`
	AirQuality  AirQuality `json:"airQuality"`
}

type AirQuality struct {
	Category  string `json:"category"`
	Statement string `json:"statement"`
}

type Daily struct {
	Date             string  `json:"date"`
	High             int64   `json:"high"`
	DayDescription   string  `json:"dayDescription"`
	RainProbability  string  `json:"rainProbability"`
	Low              int64   `json:"low,omitempty"`
	NightDescription *string `json:"nightDescription,omitempty"`
}

type Hourly struct {
	Time            string `json:"time"`
	Temperature     int64  `json:"temperature"`
	RainProbability string `json:"rainProbability"`
}

type Location struct {
	City   string `json:"city"`
	Region string `json:"region"`
}

type Sun struct {
	Duration string `json:"duration"`
	Sunrise  string `json:"sunrise"`
	Sunset   string `json:"sunset"`
}
