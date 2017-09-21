package goperfstat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/montanaflynn/stats"
	"math"
	"strconv"
	"time"
)

var (
	SummaryPercentiles = []float64{50, 90, 99, 99.9, 100}
)

type SafeFloat64 float64

func (f SafeFloat64) MarshalJSON() ([]byte, error) {
	if math.IsNaN(float64(f)) {
		return json.Marshal(nil)
	}
	return json.Marshal(float64(f))
}

type DistributionSummary struct {
	Percentiles map[string]SafeFloat64
	Min         SafeFloat64
	Max         SafeFloat64
	Average     SafeFloat64
	Stddev      SafeFloat64
	Variance    SafeFloat64
	Count       uint64
}

type PerfReport struct {
	Counters      map[string]float64
	Functions     map[string]DistributionSummary
	Distributions map[string]DistributionSummary
}

func (p *PerfContext) ReportJson() []byte {
	b, e := json.Marshal(p.ReportData())
	if e != nil {
		fmt.Printf("%v", e)
	}
	return b
}

func (p *PerfContext) ReportData() PerfReport {
	r := PerfReport{
		Counters:      make(map[string]float64),
		Functions:     make(map[string]DistributionSummary),
		Distributions: make(map[string]DistributionSummary),
	}

	for id, counter := range p.counters {
		name, found := globalIdRegistry.Counter2Name[id]
		if !found {
			name = fmt.Sprintf("$counter_%v", id)
		}
		r.Counters[name] = counter.Count
	}

	for id, samples := range p.distributions {
		summary := CalculateDistributionSummary(samples.Samples)
		name, found := globalIdRegistry.SampleId2Name[id]
		if !found {
			name = fmt.Sprintf("sample_%v", id)
		}
		r.Distributions[name] = summary
	}

	for id, perf := range p.functions {
		summary := CalculateDistributionSummary(perf.times)
		name, found := globalIdRegistry.FuncId2Name[id]
		if !found {
			name = fmt.Sprintf("function_%v", id)
		}
		r.Functions[name] = summary
	}

	return r
}

func CalculateDistributionSummary(samples stats.Float64Data) DistributionSummary {
	var summary DistributionSummary
	summary.Count = uint64(samples.Len())
	summary.Percentiles = make(map[string]SafeFloat64)
	if len(samples) > 0 {
		var f float64
		f, _ = samples.Min()
		summary.Min = SafeFloat64(f)

		f, _ = samples.Max()
		summary.Max = SafeFloat64(f)

		f, _ = samples.Mean()
		summary.Average = SafeFloat64(f)

		f, _ = samples.StandardDeviation()
		summary.Stddev = SafeFloat64(f)

		f, _ = samples.Variance()
		summary.Variance = SafeFloat64(f)

		for _, percentile := range SummaryPercentiles {
			key := strconv.FormatFloat(percentile, 'f', -1, 32)

			f, _ = samples.Percentile(percentile)
			summary.Percentiles[key] = SafeFloat64(f)
		}
	}
	return summary
}

func sampleSummary(samples stats.Float64Data) string {
	var summary string
	if len(samples) > 0 {
		min, _ := samples.Min()
		max, _ := samples.Max()
		fifty, _ := samples.Percentile(50)
		ninety, _ := samples.Percentile(90)
		ninenine, _ := samples.Percentile(99)
		summary = fmt.Sprintf("; 'min,50,90,99,max': [%v, %v, %v, %v, %v]", min, fifty, ninety, ninenine, max)
	}
	return summary
}

func (p *PerfContext) Report() string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%v elapsed\n", time.Since(p.startTime)))
	buf.WriteString("Counters:\n")
	for id, counter := range p.counters {
		name, found := globalIdRegistry.Counter2Name[id]
		if found {
			buf.WriteString(fmt.Sprintf("\t%v: %v\n", name, counter.Count))
		} else {
			buf.WriteString(fmt.Sprintf("\tcounter_%v: %v\n", id, counter.Count))
		}
	}

	buf.WriteString("Samples:\n")
	for id, samples := range p.distributions {
		summary := sampleSummary(samples.Samples)
		name, found := globalIdRegistry.SampleId2Name[id]
		if found {
			buf.WriteString(fmt.Sprintf("\t%v: { count: %v%v }\n", name, len(samples.Samples), summary))
		} else {
			buf.WriteString(fmt.Sprintf("\tsample_%v: { count: %v%v }\n", id, len(samples.Samples), summary))
		}
	}

	buf.WriteString("Functions:\n")
	for id, perf := range p.functions {
		summary := sampleSummary(perf.times)
		name, found := globalIdRegistry.FuncId2Name[id]
		if found {
			buf.WriteString(fmt.Sprintf("\t%v: { count: %v%v }\n", name, perf.count, summary))
		} else {
			buf.WriteString(fmt.Sprintf("\tfunction_%v: { count: %v%v }\n", id, perf.count, summary))
		}
	}

	return buf.String()
}
