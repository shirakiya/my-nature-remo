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
	"github.com/pkg/errors"
)

const (
	NatureDeviceURL = "https://api.nature.global/1/devices"

	NatureAccessTokenEnv   = "NATURE_ACCESS_TOKEN"
	MackerelAPIKeyEnv      = "MACKEREL_API_KEY" //nolint:gosec // env var name, not a credential
	MackerelServiceNameEnv = "MACKEREL_SERVICE_NAME"
)

//nolint:tagliatelle // field names follow the Nature Remo API response format
type Device struct {
	ID          string `json:"id"`
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

func run(ctx context.Context) error {
	natureToken := mustGetEnv(NatureAccessTokenEnv)
	mackerelAPIKey := mustGetEnv(MackerelAPIKeyEnv)
	mackerelServiceName := mustGetEnv(MackerelServiceNameEnv)

	req, err := http.NewRequestWithContext(ctx, "GET", NatureDeviceURL, nil)
	if err != nil {
		return errors.Wrap(err, "error was found NewRequestWithContext")
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", natureToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "error in requesting to Nature")
	}
	defer resp.Body.Close()

	var devices []Device

	err = json.NewDecoder(resp.Body).Decode(&devices)
	if err != nil {
		return errors.Wrap(err, "decode error")
	}

	if len(devices) != 1 {
		return errors.New("device count is not 1")
	}

	device := devices[0]
	metricValues := []*mackerel.MetricValue{
		{
			Name:  fmt.Sprintf("%s.temperature", device.ID),
			Time:  device.NewestEvent.Te.CreatedAt.Unix(),
			Value: device.NewestEvent.Te.Val,
		},
		{
			Name:  fmt.Sprintf("%s.humidity", device.ID),
			Time:  device.NewestEvent.Hu.CreatedAt.Unix(),
			Value: device.NewestEvent.Hu.Val,
		},
	}

	mackerelClient := mackerel.NewClient(mackerelAPIKey)
	if err := mackerelClient.PostServiceMetricValues(mackerelServiceName, metricValues); err != nil {
		return errors.Wrap(err, "error in requesting to Mackerel")
	}

	return nil
}

func main() {
	log.Println("start")

	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatalf("%+v", err)
	}

	log.Println("finished as success")
}
