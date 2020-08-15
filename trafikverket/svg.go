package trafikverket

import (
	"io"
	"sort"
	"strconv"
	"time"

	svg "github.com/ajstarks/svgo"
)

func SVG(w io.Writer, trains map[string][]Stops) {
	canvas := svg.New(w)
	defer canvas.End()
	canvas.Start(700, 1000)

	type rails struct {
		A Station
		B Station
	}

	type candidate struct {
		arrivalAtEnd time.Time
		color        string
		pos          Location
	}

	candidates := map[rails][]candidate{}

	for trainNumber, stops := range trains {
		n, err := strconv.Atoi(trainNumber)
		if err != nil {
			n = 0
		}
		color := colors[n%len(colors)]

		var prev Stops
		for _, stop := range stops {
			shouldBreak, loc := lineForStop(stop, prev)

			if !prev.IsNone() {
				key := rails{A: prev.Station, B: stop.Station}
				cand := candidate{
					arrivalAtEnd: stop.Arrival,
					color:        color,
					pos:          loc,
				}
				if curr, exists := candidates[key]; exists {
					candidates[key] = append(curr, cand)
				} else {
					candidates[key] = []candidate{cand}
				}
			}

			prev = stop
			if shouldBreak {
				break
			}
		}
	}

	for stations, cands := range candidates {
		sort.Slice(cands, func(i, j int) bool {
			return cands[i].arrivalAtEnd.After(cands[j].arrivalAtEnd)
		})

		line := cands[0]
		x0, y0 := stations.A.Location.canvasPos()
		x1, y1 := line.pos.canvasPos()
		canvas.Line(x0, y0, x1, y1, "fill:none;stroke-width:3;stroke:"+line.color)
	}
}

func (l Location) canvasPos() (int, int) {
	x := l.Lon - 10
	y := 20 - (l.Lat - 53)

	return int(x * 50), int(y * 50)
}

func lineForStop(stop, prev Stops) (bool, Location) {
	var loc Location
	var shouldBreak bool
	if stop.Arrival.After(time.Now()) {
		dNow := time.Now().Sub(prev.Departure)
		dArrival := stop.Arrival.Sub(prev.Departure)
		ratio := dNow.Seconds() / dArrival.Seconds()
		loc = PointBetween(ratio, prev.Station, stop.Station)
		shouldBreak = true
	} else {
		loc = stop.Station.Location
	}

	return shouldBreak, loc
}

var colors = []string{
	"black",
	"aliceblue",
	"antiquewhite",
	"aqua",
	"aquamarine",
	"azure",
	"beige",
	"bisque",
	"blanchedalmond",
	"blue",
	"blueviolet",
	"brown",
	"burlywood",
	"cadetblue",
	"chartreuse",
	"chocolate",
	"coral",
	"cornflowerblue",
	"cornsilk",
	"crimson",
	"cyan",
	"darkblue",
	"darkcyan",
	"darkgoldenrod",
	"darkgray",
	"darkgreen",
	"darkgrey",
	"darkkhaki",
	"darkmagenta",
	"darkolivegreen",
	"darkorange",
	"darkorchid",
	"darkred",
	"darksalmon",
	"darkseagreen",
	"darkslateblue",
	"darkslategray",
	"darkslategrey",
	"darkturquoise",
	"darkviolet",
	"deeppink",
	"deepskyblue",
	"dimgray",
	"dimgrey",
	"dodgerblue",
	"firebrick",
	"floralwhite",
	"forestgreen",
	"fuchsia",
	"gainsboro",
	"ghostwhite",
	"gold",
	"goldenrod",
	"gray",
	"green",
	"greenyellow",
	"grey",
	"honeydew",
	"hotpink",
	"indianred",
	"indigo",
	"ivory",
	"khaki",
	"lavender",
	"lavenderblush",
	"lawngreen",
	"lemonchiffon",
	"lightblue",
	"lightcoral",
	"lightcyan",
	"lightgoldenrodyellow",
	"lightgray",
	"lightgreen",
	"lightgrey",
	"lightpink",
	"lightsalmon",
	"lightseagreen",
	"lightskyblue",
	"lightslategray",
	"lightslategrey",
	"lightsteelblue",
	"lightyellow",
	"lime",
	"limegreen",
	"linen",
	"magenta",
	"maroon",
	"mediumaquamarine",
	"mediumblue",
	"mediumorchid",
	"mediumpurple",
	"mediumseagreen",
	"mediumslateblue",
	"mediumspringgreen",
	"mediumturquoise",
	"mediumvioletred",
	"midnightblue",
	"mintcream",
	"mistyrose",
	"moccasin",
	"navajowhite",
	"navy",
	"oldlace",
	"olive",
	"olivedrab",
	"orange",
	"orangered",
	"orchid",
	"palegoldenrod",
	"palegreen",
	"paleturquoise",
	"palevioletred",
	"papayawhip",
	"peachpuff",
	"peru",
	"pink",
	"plum",
	"powderblue",
	"purple",
	"red",
	"rosybrown",
	"royalblue",
	"saddlebrown",
	"salmon",
	"sandybrown",
	"seagreen",
	"seashell",
	"sienna",
	"silver",
	"skyblue",
	"slateblue",
	"slategray",
	"slategrey",
	"snow",
	"springgreen",
	"steelblue",
	"tan",
	"teal",
	"thistle",
	"tomato",
	"turquoise",
	"violet",
	"wheat",
	"white",
	"whitesmoke",
	"yellow",
	"yellowgreen",
}
