package column

import (
	"errors"
	"math"
	"sort"
)

type TrayReading struct {
	Tray        int
	Temperature float64
	Pressure    float64
}

type Profile struct{ readings map[int]TrayReading }

func NewProfile() *Profile { return &Profile{readings: make(map[int]TrayReading)} }

func (p *Profile) Put(reading TrayReading) error {
	if reading.Tray < 1 || reading.Tray > 200 {
		return errors.New("tray outside configured column")
	}
	if math.IsNaN(reading.Temperature) || math.IsNaN(reading.Pressure) {
		return errors.New("tray reading is not numeric")
	}
	p.readings[reading.Tray] = reading
	return nil
}

func (p *Profile) Ordered() []TrayReading {
	keys := make([]int, 0, len(p.readings))
	for tray := range p.readings {
		keys = append(keys, tray)
	}
	sort.Ints(keys)
	out := make([]TrayReading, 0, len(keys))
	for _, tray := range keys {
		out = append(out, p.readings[tray])
	}
	return out
}

func (p *Profile) Gradient() (float64, error) {
	values := p.Ordered()
	if len(values) < 2 {
		return 0, errors.New("profile needs at least two trays")
	}
	first := values[0]
	last := values[len(values)-1]
	return (last.Temperature - first.Temperature) / float64(last.Tray-first.Tray), nil
}

func (p *Profile) Stable(maxDelta float64) bool {
	values := p.Ordered()
	if len(values) < 2 {
		return false
	}
	for index := 1; index < len(values); index++ {
		delta := values[index].Temperature - values[index-1].Temperature
		if delta > maxDelta || delta < -maxDelta {
			return false
		}
	}
	return true
}

func (p *Profile) Count() int { return len(p.readings) }
