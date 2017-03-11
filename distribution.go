package goperfstat

import (
	"errors"
	"fmt"
	"github.com/montanaflynn/stats"
	"math/rand"
	"time"
)

type SampleDistribution struct {
	rate    float32
	Samples stats.Float64Data
	r       *rand.Rand
}

func NewSampleDistribution(rate float32, capacity int) *SampleDistribution {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	samples := make([]float64, capacity)
	samples = samples[:0]
	return &SampleDistribution{rate: rate, r: r, Samples: samples}
}

func (d *SampleDistribution) Sample(v float64) {
	f := d.r.Float32()
	if f < d.rate {
		d.Samples = append(d.Samples, v)
	}
}

func GetSampleStats(context *PerfContext, id int) (*SampleDistribution, error) {
	if context == nil {
		context = globalContext
	}

	d, found := context.distributions[id]
	if !found {
		return nil, errors.New(fmt.Sprintf("no sample stats for %v", id))
	}
	return d, nil
}
