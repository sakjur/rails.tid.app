package trafikverket

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

const trafikverketAPIURL = "https://api.trafikinfo.trafikverket.se/v2/data.json"

type Client struct {
	apiKey   string
	stations map[string]Station
	timeout  time.Duration
}

func NewClient() (Client, error) {
	apiKey, err := getAPIKey()
	if err != nil {
		return Client{}, err
	}

	return Client{
			apiKey:  apiKey,
			timeout: 5 * time.Second,
		},
		nil
}

func (c Client) request(ctx context.Context, body io.Reader) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		trafikverketAPIURL,
		body,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/xml")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)

	return respBody, err
}

func duration(dur time.Duration) string {
	if dur == 0 {
		return "$now"
	}

	totalSeconds := int(math.Abs(dur.Seconds()))

	var (
		seconds = totalSeconds % 60
		minutes = totalSeconds / 60 % 60
		hours   = totalSeconds / 3600 % 24
		days    = totalSeconds / (24 * 60 * 60)
		sign    = ""
	)

	if dur < 0 {
		sign = "-"
	}

	return fmt.Sprintf("$dateadd(%s%d.%d:%02d:%02d)", sign, days, hours, minutes, seconds)
}
