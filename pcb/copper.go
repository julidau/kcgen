package pcb

import (
	"errors"
	"fmt"
	"math"

	"github.com/nsf/sexp"
)

// XY represents a point in space.
type XY struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (xy XY) Distance(xy2 XY) float64 {
	return math.Sqrt(math.Pow(xy2.X-xy.X, 2) + math.Pow(xy2.Y-xy.Y, 2))
}

// XYX represents a point in 3d space.
type XYZ struct {
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Z        float64 `json:"z,omitempty"`
	ZPresent bool    `json:"z_present,omitempty"`
	Unlocked bool    `json:"unlocked,omitempty"`
}

// ViaType represents a via Type
type ViaType int

// Via Types
const (
	ViaThrough = ViaType(iota)
	ViaMicro
	ViaBlind
)

// Via represents a via.
type Via struct {
	At       XY       `json:"position"`
	Size     float64  `json:"size"`
	Drill    float64  `json:"drill,omitempty"`
	Layers   []string `json:"layers"`
	NetIndex int      `json:"net_index"`

	StatusFlags string `json:"status_flags"`
	ViaType     ViaType

	order int
}

type ZoneHatch struct {
	Mode  string  `json:"mode"`
	Pitch float64 `json:"pitch"`
}

type ZoneConnectPads struct {
	Clearance float64 `json:"clearance"`
	Mode      string  `json:"mode"`
}

type ZoneFill struct {
	IsFilled           bool    `json:"is_filled"`
	Mode               string  `json:"mode"`
	Segments           int     `json:"segments"`
	ThermalGap         float64 `json:"thermal_gap"`
	ThermalBridgeWidth float64 `json:"thermal_bridge_width"`
}

type ZoneKeepout struct {
	TracksAllowed     bool `json:"tracks_allowed"`
	ViasAllowed       bool `json:"vias_allowed"`
	CopperPourAllowed bool `json:"copperpour_allowed"`
}

// Zone represents a zone.
type Zone struct {
	IsKeepout bool     `json:"is_keepout"`
	NetNum    int      `json:"net_num"`
	NetName   string   `json:"net_name"`
	Layers    []string `json:"layers"`

	Tstamp   string `json:"tstamp"`
	Priority int    `json:"priority"`

	Hatch ZoneHatch `json:"hatch"`

	ConnectPads ZoneConnectPads `json:"connect_pads"`

	Fill ZoneFill `json:"fill"`

	Keepout ZoneKeepout `json:"keepout"`

	MinThickness        float64 `json:"min_thickness"`
	FilledAreaThickness bool    `json:"filled_area_thickness"`

	Polys     [][]XY `json:"polys,omitempty"`
	BasePolys [][]XY `json:"base_polys,omitempty"`

	order int
}

// Track represents a PCB track.
type Track struct {
	Start    XY      `json:"start"`
	End      XY      `json:"end"`
	Width    float64 `json:"width"`
	Layer    string  `json:"layer"`
	NetIndex int     `json:"net_index"`

	Tstamp string `json:"tstamp"`

	StatusFlags string `json:"status_flags"`

	order int
}

func parseVia(n sexp.Helper, ordering int) (Via, error) {
	v := Via{order: ordering}
	for x := 1; x < n.MustNode().NumChildren(); x++ {
		c := n.Child(x)
		if c.IsList() {
			switch c.Child(0).MustString() {
			case "size":
				v.Size = c.Child(1).MustFloat64()
			case "drill":
				v.Drill = c.Child(1).MustFloat64()
			case "net":
				v.NetIndex = c.Child(1).MustInt()
			case "at":
				v.At.X = c.Child(1).MustFloat64()
				v.At.Y = c.Child(2).MustFloat64()
			case "status":
				v.StatusFlags = c.Child(1).MustString()
			case "layers":
				for j := 1; j < c.MustNode().NumChildren(); j++ {
					v.Layers = append(v.Layers, c.Child(j).MustString())
				}
			}
		} else {
			// must be via type
			switch t := c.MustString(); t {
			case "blind":
				v.ViaType = ViaBlind
			case "micro":
				v.ViaType = ViaMicro
			default:
				return v, errors.New("via invalid type " + t)
			}
		}
	}
	return v, nil
}

func parseZone(n sexp.Helper, ordering int) (*Zone, error) {
	z := Zone{order: ordering}
	for x := 1; x < n.MustNode().NumChildren(); x++ {
		c := n.Child(x)
		switch c.Child(0).MustString() {
		case "net":
			z.NetNum = c.Child(1).MustInt()
		case "net_name":
			z.NetName = c.Child(1).MustString()
		case "layer":
			z.Layers = []string{c.Child(1).MustString()}
		case "layers":
			z.Layers = nil
			for j := 1; j < c.MustNode().NumChildren(); j++ {
				z.Layers = append(z.Layers, c.Child(j).MustString())
			}
		case "tstamp":
			z.Tstamp = c.Child(1).MustString()

		case "hatch":
			z.Hatch.Mode = c.Child(1).MustString()
			z.Hatch.Pitch = c.Child(2).MustFloat64()
		case "min_thickness":
			z.MinThickness = c.Child(1).MustFloat64()
		case "priority":
			z.Priority = c.Child(1).MustInt()
		case "filled_area_thickness":
			z.FilledAreaThickness = c.Child(1).MustString() == "yes"

		case "connect_pads":
			for y := 1; y < c.MustNode().NumChildren(); y++ {
				c2 := c.Child(y)
				if c2.IsList() {
					switch c2.Child(0).MustString() {
					case "clearance":
						z.ConnectPads.Clearance = c2.Child(1).MustFloat64()
					}
				} else {
					z.ConnectPads.Mode = c2.MustString()
				}
			}
		case "fill":
			for y := 1; y < c.MustNode().NumChildren(); y++ {
				c2 := c.Child(y)
				if c2.IsScalar() {
					switch c2.MustNode().Value {
					case "yes":
						z.Fill.IsFilled = true
					default:
						return nil, fmt.Errorf("unhandled scalar in zone fill: %v", c2.MustNode().Value)
					}
					continue
				}

				switch c2.Child(0).MustString() {
				case "mode":
					z.Fill.Mode = c2.Child(1).MustString()
				case "arc_segments":
					z.Fill.Segments = c2.Child(1).MustInt()
				case "thermal_gap":
					z.Fill.ThermalGap = c2.Child(1).MustFloat64()
				case "thermal_bridge_width":
					z.Fill.ThermalBridgeWidth = c2.Child(1).MustFloat64()
				}
			}

		case "keepout":
			z.IsKeepout = true
			for y := 1; y < c.MustNode().NumChildren(); y++ {
				c2 := c.Child(y)

				switch c2.Child(0).MustString() {
				case "tracks":
					z.Keepout.TracksAllowed = c2.Child(1).MustString() == "allowed"
				case "vias":
					z.Keepout.ViasAllowed = c2.Child(1).MustString() == "allowed"
				case "copperpour":
					z.Keepout.CopperPourAllowed = c2.Child(1).MustString() == "allowed"
				}
			}

		case "polygon":
			var points []XY
			for y := 1; y < c.Child(1).MustNode().NumChildren(); y++ {
				pt := c.Child(1).Child(y)
				ptType, err2 := pt.Child(0).String()
				if err2 != nil || ptType != "xy" {
					return nil, errors.New("zone.polygon point is not xy point")
				}
				points = append(points, XY{X: pt.Child(1).MustFloat64(), Y: pt.Child(2).MustFloat64()})
			}
			z.BasePolys = append(z.BasePolys, points)

		case "filled_polygon":
			var points []XY
			for y := 1; y < c.Child(1).MustNode().NumChildren(); y++ {
				pt := c.Child(1).Child(y)
				ptType, err2 := pt.Child(0).String()
				if err2 != nil || ptType != "xy" {
					return nil, errors.New("zone.filled_polygon point is not xy point")
				}
				points = append(points, XY{X: pt.Child(1).MustFloat64(), Y: pt.Child(2).MustFloat64()})
			}
			z.Polys = append(z.Polys, points)
		}
	}
	return &z, nil
}

func parseSegment(n sexp.Helper, ordering int) (Track, error) {
	t := Track{order: ordering}
	for x := 1; x < n.MustNode().NumChildren(); x++ {
		c := n.Child(x)
		switch c.Child(0).MustString() {
		case "width":
			t.Width = c.Child(1).MustFloat64()
		case "net":
			t.NetIndex = c.Child(1).MustInt()
		case "layer":
			t.Layer = c.Child(1).MustString()
		case "tstamp":
			t.Tstamp = c.Child(1).MustString()
		case "status":
			t.StatusFlags = c.Child(1).MustString()
		case "start":
			t.Start = XY{X: c.Child(1).MustFloat64(), Y: c.Child(2).MustFloat64()}
		case "end":
			t.End = XY{X: c.Child(1).MustFloat64(), Y: c.Child(2).MustFloat64()}
		}
	}
	return t, nil
}
