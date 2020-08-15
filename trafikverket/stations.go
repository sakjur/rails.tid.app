package trafikverket

import (
	"bytes"
	"context"
)

func (c Client) allStations() ([]byte, error) {
	return c.query(
		query{
			ObjectType:    "TrainStation",
			SchemaVersion: "1",
			Filter: xmlFilter{
				Filters: []filter{
					eq{
						Name:  "Advertised",
						Value: "true",
					},
				},
			},
			Include: []include{
				{Attribute: "AdvertisedShortLocationName"},
				{Attribute: "Geometry.WGS84"},
				{Attribute: "LocationSignature"},
				{Attribute: "CountryCode"},
			},
		},
	)
}

type Station struct {
	Name        string
	ID          string
	CountryCode string
	Location    Location
}

type Location struct {
	Lon float64
	Lat float64
}

func (l Location) IsZero() bool {
	return l.Lon == 0 && l.Lat == 0
}

func (c Client) cacheStations() (Client, error) {
	req, err := c.allStations()
	if err != nil {
		return c, err
	}

	body, err := c.request(context.Background(), bytes.NewReader(req))
	if err != nil {
		return c, err
	}

	json, err := unmarshalResponse(bytes.NewReader(body))
	if err != nil {
		return c, err
	}

	stations := map[string]Station{}
	for _, r := range json.Response.Result {
		for _, s := range r.TrainStations {
			station, err := s.AsStation()
			if err != nil {
				return c, err
			}

			stations[station.ID] = station
		}
	}

	c.stations = stations
	return c, nil
}

func (c Client) Stations() (map[string]Station, error) {
	if c.stations == nil {
		var err error
		c, err = c.cacheStations()
		if err != nil {
			return nil, err
		}
	}
	return c.stations, nil
}
