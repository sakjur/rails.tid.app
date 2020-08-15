package trafikverket

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
)

type jsonResponse struct {
	Response struct {
		Result []struct {
			TrainStations     []jsonTrainStations     `json:"TrainStation"`
			TrainAnnouncement []jsonTrainAnnouncement `json:"TrainAnnouncement"`
		} `json:"RESULT"`
	} `json:"RESPONSE"`
}

type jsonTrainStations struct {
	AdvertisedShortLocationName string `json:"AdvertisedShortLocationName"`
	CountryCode                 string `json:"CountryCode"`
	Geometry                    struct {
		WGS84 string `json:"WGS84"`
	} `json:"Geometry"`
	LocationSignature string `json:"LocationSignature"`
}

type jsonTrainAnnouncement struct {
	AdvertisedTimeAtLocation  string `json:"AdvertisedTimeAtLocation"`
	TimeAtLocationWithSeconds string `json:"TimeAtLocation"`
	EstimatedTimeAtLocation   string `json:"EstimatedTimeAtLocation"`
	TrainOwner                string `json:"TrainOwner"`
	ActivityType              string `json:"ActivityType"`
	AdvertisedTrainIdent      string `json:"AdvertisedTrainIdent"`
	LocationSignature         string `json:"LocationSignature"`
}

func unmarshalResponse(r io.Reader) (jsonResponse, error) {
	raw, err := ioutil.ReadAll(r)
	if err != nil {
		return jsonResponse{}, err
	}

	resp := jsonResponse{}
	err = json.Unmarshal(raw, &resp)
	return resp, err
}

var pointRE = regexp.MustCompile(`POINT \(([\d.]+) ([\d.]+)\)`)

func (j jsonTrainStations) AsStation() (Station, error) {
	geometry := pointRE.FindStringSubmatch(j.Geometry.WGS84)

	if len(geometry) < 3 {
		return Station{}, fmt.Errorf("could not parse coordinates for station. Raw coordinates: %s", j.Geometry.WGS84)
	}

	lon, err := strconv.ParseFloat(geometry[1], 64)
	if err != nil {
		return Station{}, err
	}

	lat, err := strconv.ParseFloat(geometry[2], 64)
	if err != nil {
		return Station{}, err
	}

	return Station{
		Name:        j.AdvertisedShortLocationName,
		ID:          j.LocationSignature,
		CountryCode: j.CountryCode,
		Location: Location{
			Lon: lon,
			Lat: lat,
		},
	}, nil
}

func PointBetween(ratio float64, from, to Station) Location {
	dLat := (to.Location.Lat - from.Location.Lat) * ratio
	dLon := (to.Location.Lon - from.Location.Lon) * ratio

	return Location{
		Lat: from.Location.Lat + dLat,
		Lon: from.Location.Lon + dLon,
	}
}
