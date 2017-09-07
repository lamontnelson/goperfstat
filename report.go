package goperfstat

import (
	"bytes"
	"fmt"
	"github.com/montanaflynn/stats"
	"strconv"
	"time"
)

var (
	SummaryPercentiles = []float64{50, 90, 99, 99.9, 100}
)

type DistributionSummary struct {
	percentiles map[string]float64
	min         float64
	max         float64
	average     float64
	stddev      float64
	variance    float64
	count       uint64
}

type PerfReport struct {
	counters      map[string]float64
	functions     map[string]DistributionSummary
	distributions map[string]DistributionSummary
}

func (p *PerfContext) ReportData() PerfReport {
	r := PerfReport{
		counters:      make(map[string]float64),
		functions:     make(map[string]DistributionSummary),
		distributions: make(map[string]DistributionSummary),
	}

	for id, counter := range p.counters {
		name, found := globalIdRegistry.Counter2Name[id]
		if !found {
			name = fmt.Sprintf("$counter_%v", id)
		}
		r.counters[name] = counter.Count
	}

	for id, samples := range p.distributions {
		summary := CalculateDistributionSummary(samples.Samples)
		name, found := globalIdRegistry.SampleId2Name[id]
		if !found {
			name = fmt.Sprintf("sample_%v", id)
		}
		r.distributions[name] = summary
	}

	for id, perf := range p.functions {
		summary := CalculateDistributionSummary(perf.times)
		name, found := globalIdRegistry.FuncId2Name[id]
		if !found {
			name = fmt.Sprintf("function_%v", id)
		}
		r.functions[name] = summary
	}

	return r
}

func CalculateDistributionSummary(samples stats.Float64Data) DistributionSummary {
	var summary DistributionSummary
	summary.count = uint64(samples.Len())
	summary.percentiles = make(map[string]float64)
	if len(samples) > 0 {
		summary.min, _ = samples.Min()
		summary.max, _ = samples.Max()
		summary.average, _ = samples.Mean()
		summary.stddev, _ = samples.StandardDeviation()
		summary.variance, _ = samples.Variance()
		for _, percentile := range SummaryPercentiles {
			key := strconv.FormatFloat(percentile, 'f', -1, 32)
			summary.percentiles[key], _ = samples.Percentile(percentile)
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
