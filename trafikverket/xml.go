package trafikverket

import "encoding/xml"

type include struct {
	XMLName   struct{} `xml:"INCLUDE"`
	Attribute string   `xml:",chardata"`
}

type root struct {
	XMLName struct{} `xml:"REQUEST"`
	Login   auth
	Queries []query `xml:"QUERY"`
}

type auth struct {
	XMLName           struct{} `xml:"LOGIN"`
	AuthenticationKey string   `xml:"authenticationkey,attr"`
}

type xmlFilter struct {
	XMLName struct{} `xml:"Filter"`
	Filters []filter
}

type filter interface {
	isTrafikverketFilter()
}

type query struct {
	ObjectType    string `xml:"objecttype,attr"`
	SchemaVersion string `xml:"schemaversion,attr"`
	Filter        xmlFilter
	Include       []include
}

type eq struct {
	XMLName struct{} `xml:"EQ"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:"value,attr"`
}

func (e eq) isTrafikverketFilter() {}

type gt struct {
	XMLName struct{} `xml:"GT"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:"value,attr"`
}

func (e gt) isTrafikverketFilter() {}

type lt struct {
	XMLName struct{} `xml:"LT"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:"value,attr"`
}

func (e lt) isTrafikverketFilter() {}

func (c Client) query(queries ...query) ([]byte, error) {
	query := root{
		Login: auth{
			AuthenticationKey: c.apiKey,
		},
		Queries: queries,
	}

	return xml.MarshalIndent(query, "", "\t")
}

func (c Client) auth() auth {
	return auth{
		AuthenticationKey: c.apiKey,
	}
}
