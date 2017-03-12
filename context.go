package goperfstat

import (
	"fmt"
	"github.com/montanaflynn/stats"
)

type PerfContext struct {
	functions     map[int]*FuncPerf
	distributions map[int]*SampleDistribution

	FuncId2Name   map[int]string
	SampleId2Name map[int]string
}

func NewPerfContext() *PerfContext {
	return &PerfContext{
		functions:     make(map[int]*FuncPerf),
		distributions: make(map[int]*SampleDistribution),
		FuncId2Name:   make(map[int]string),
		SampleId2Name: make(map[int]string),
	}
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

func (p *PerfContext) Report() {

	fmt.Printf("Samples:\n")
	for id, samples := range p.distributions {
		summary := sampleSummary(samples.Samples)
		name, found := p.SampleId2Name[id]
		if found {
			fmt.Printf("\t%v: { count: %v%v }\n", name, len(samples.Samples), summary)
		} else {
			fmt.Printf("\tid %v: { count: %v%v }\n", id, len(samples.Samples), summary)
		}
	}

	fmt.Printf("Functions:\n")
	for id, perf := range p.functions {
		summary := sampleSummary(perf.times)
		name, found := p.FuncId2Name[id]
		if found {
			fmt.Printf("\t%v: { count: %v%v }\n", name, perf.count, summary)
		} else {
			fmt.Printf("\tid %v: { count: %v%v }\n", id, perf.count, summary)
		}
	}
}

func (p *PerfContext) RegFuncId(name string, id int) {
	p.FuncId2Name[id] = name
}

func (p *PerfContext) RegDistId(name string, id int) {
	p.SampleId2Name[id] = name
}

var globalContext *PerfContext

func InitGlobalPerfContext() {
	globalContext = NewPerfContext()
}

func init() {
	InitGlobalPerfContext()
}
