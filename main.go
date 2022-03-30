package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mackerelio/mackerel-client-go"
)

const (
	NATURE_DEVICE_URL = "https://api.nature.global/1/devices"

	NATURE_ACCESS_TOKEN_ENV   = "NATURE_ACCESS_TOKEN"
	MACKEREL_API_KEY_ENV      = "MACKEREL_API_KEY"
	MACKEREL_SERVICE_NAME_ENV = "MACKEREL_SERVICE_NAME"
)

type Device struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	NewestEvent struct {
		Hu struct {
			Val       int       `json:"val"`
			CreatedAt time.Time `json:"created_at"`
		} `json:"hu"`
		Te struct {
			Val       float64   `json:"val"`
			CreatedAt time.Time `json:"created_at"`
		} `json:"te"`
	} `json:"newest_events"`
}

func mustGetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("%s must be set", key)
	}

	return value
}

func main() {
	log.Println("start")

	ctx := context.Background()

	natureToken := mustGetEnv(NATURE_ACCESS_TOKEN_ENV)
	mackerelApiKey := mustGetEnv(MACKEREL_API_KEY_ENV)
	mackerelServiceName := mustGetEnv(MACKEREL_SERVICE_NAME_ENV)

	req, err := http.NewRequestWithContext(ctx, "GET", NATURE_DEVICE_URL, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", natureToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var devices []Device
	err = json.NewDecoder(resp.Body).Decode(&devices)
	if err != nil {
		log.Fatal(err)
	}
	if len(devices) != 1 {
		log.Fatal("device count is not 1")
	}

	device := devices[0]
	metricValues := []*mackerel.MetricValue{
		{
			Name:  fmt.Sprintf("%s.temperature", device.Id),
			Time:  device.NewestEvent.Te.CreatedAt.Unix(),
			Value: device.NewestEvent.Te.Val,
		},
		{
			Name:  fmt.Sprintf("%s.humidity", device.Id),
			Time:  device.NewestEvent.Hu.CreatedAt.Unix(),
			Value: device.NewestEvent.Hu.Val,
		},
	}

	mackerelClient := mackerel.NewClient(mackerelApiKey)
	if err := mackerelClient.PostServiceMetricValues(mackerelServiceName, metricValues); err != nil {
		log.Fatal(err)
	}

	log.Println("finished")
}
