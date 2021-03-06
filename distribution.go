package goperfstat

import (
	"errors"
	"fmt"
	"github.com/montanaflynn/stats"
	"math/rand"
	"time"
)

type SampleDistribution struct {
	rate    float64
	Samples stats.Float64Data
	r       *rand.Rand
}

func NewSampleDistribution(rate float64, capacity int) *SampleDistribution {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	samples := make([]float64, capacity)
	samples = samples[:0]
	return &SampleDistribution{rate: rate, r: r, Samples: samples}
}

func (d *SampleDistribution) Sample(v float64) {
	f := d.r.Float64()
	if f < d.rate {
		d.Samples = append(d.Samples, v)
	}
}

func TakeSample(context *PerfContext, id int, rate float64, capacity int, v float64) {
	if context == nil {
		context = globalContext
	}

	context.distMu.Lock()
	defer context.distMu.Unlock()

	stats, err := GetSampleStats(context, id)
	if err != nil {
		stats = NewSampleDistribution(rate, capacity)
		context.distributions[id] = stats
	}
	stats.Sample(v)
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
