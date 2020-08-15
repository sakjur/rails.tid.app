package trafikverket

import (
	"bytes"
	"context"
	"sort"
	"time"
)

const longestDurationBetweenStops = 6 * time.Hour

func (c Client) trainAnnouncements() ([]byte, error) {
	return c.query(query{
		ObjectType:    "TrainAnnouncement",
		SchemaVersion: "1.6",
		Filter: xmlFilter{
			Filters: []filter{
				eq{Name: "Advertised", Value: "true"},
				eq{Name: "InformationOwner", Value: "SJ"},
				gt{Name: "AdvertisedTimeAtLocation", Value: duration(-longestDurationBetweenStops)},
				lt{Name: "AdvertisedTimeAtLocation", Value: duration(longestDurationBetweenStops)},
			},
		},
		Include: []include{
			{Attribute: "AdvertisedTimeAtLocation"},
			{Attribute: "TimeAtLocationWithSeconds"},
			{Attribute: "EstimatedTimeAtLocation"},
			{Attribute: "ActivityType"},
			{Attribute: "AdvertisedTrainIdent"},
			{Attribute: "LocationSignature"},
			{Attribute: "TrainOwner"},
			{Attribute: "TechnicalTrainIdent"},
		},
	})
}

type Stops struct {
	Arrival   time.Time
	Departure time.Time
	Station   Station
}

func (s Stops) IsNone() bool {
	return s.Arrival.IsZero() && s.Departure.IsZero()
}

func (c Client) Trains() (map[string][]Stops, error) {
	req, err := c.trainAnnouncements()
	if err != nil {
		return nil, err
	}

	body, err := c.request(context.Background(), bytes.NewReader(req))
	if err != nil {
		return nil, err
	}

	json, err := unmarshalResponse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	stations, err := c.Stations()
	if err != nil {
		return nil, err
	}

	type key struct {
		trainNumber string
		station     string
	}

	trains := map[key]Stops{}
	for _, r := range json.Response.Result {
		for _, t := range r.TrainAnnouncement {
			trainNumber := t.AdvertisedTrainIdent
			station := stations[t.LocationSignature]

			rawTime := firstNonEmpty(t.TimeAtLocationWithSeconds, t.EstimatedTimeAtLocation, t.AdvertisedTimeAtLocation)
			tm, err := time.Parse(time.RFC3339Nano, rawTime)
			if err != nil {
				return nil, err
			}

			k := key{trainNumber, station.ID}
			stop := trains[k]
			stop.Station = station

			switch t.ActivityType {
			case "Avgang":
				stop.Departure = tm
			case "Ankomst":
				stop.Arrival = tm
			default:
				continue
			}

			trains[k] = stop
		}
	}

	sortedTrains := map[string][]Stops{}
	for key, stop := range trains {
		if curr, exists := sortedTrains[key.trainNumber]; exists {
			sortedTrains[key.trainNumber] = append(curr, stop)
		} else {
			sortedTrains[key.trainNumber] = []Stops{stop}
		}
	}

	for _, stops := range sortedTrains {
		sort.Slice(stops, func(i, j int) bool {
			a, b := stops[i], stops[j]

			tA := a.Arrival
			if tA.IsZero() {
				tA = a.Departure
			}

			tB := b.Arrival
			if tB.IsZero() {
				tB = b.Departure
			}

			return tA.Before(tB)
		})
	}

	return sortedTrains, nil
}

func firstNonEmpty(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}
