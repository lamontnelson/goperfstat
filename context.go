package goperfstat

import (
	"fmt"
	"github.com/montanaflynn/stats"
	"time"
)

type PerfContext struct {
	functions     map[int]*FuncPerf
	distributions map[int]*SampleDistribution
	counters      map[int]*Counter
	startTime     time.Time
}
type PerfIdRegistry struct {
	FuncId2Name   map[int]string
	SampleId2Name map[int]string
	Counter2Name  map[int]string
}

func NewPerfContext() *PerfContext {
	return &PerfContext{
		functions:     make(map[int]*FuncPerf),
		distributions: make(map[int]*SampleDistribution),
		counters:      make(map[int]*Counter),
		startTime:     time.Now(),
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

func (p *PerfContext) StartTime() {
	p.startTime = time.Now()
}

func (p *PerfContext) Report() {
	fmt.Printf("%v elapsed\n", time.Since(p.startTime))

	fmt.Printf("Counters:\n")
	for id, counter := range p.counters {
		name, found := globalIdRegistry.Counter2Name[id]
		if found {
			fmt.Printf("\t%v: %v\n", name, counter.Count)
		} else {
			fmt.Printf("\tcounter_%v: %v\n", id, counter.Count)
		}
	}

	fmt.Printf("Samples:\n")
	for id, samples := range p.distributions {
		summary := sampleSummary(samples.Samples)
		name, found := globalIdRegistry.SampleId2Name[id]
		if found {
			fmt.Printf("\t%v: { count: %v%v }\n", name, len(samples.Samples), summary)
		} else {
			fmt.Printf("\tsample_%v: { count: %v%v }\n", id, len(samples.Samples), summary)
		}
	}

	fmt.Printf("Functions:\n")
	for id, perf := range p.functions {
		summary := sampleSummary(perf.times)
		name, found := globalIdRegistry.FuncId2Name[id]
		if found {
			fmt.Printf("\t%v: { count: %v%v }\n", name, perf.count, summary)
		} else {
			fmt.Printf("\tfunction_%v: { count: %v%v }\n", id, perf.count, summary)
		}
	}
}

func (p *PerfIdRegistry) RegFuncId(name string, id int) {
	p.FuncId2Name[id] = name
}

func (p *PerfIdRegistry) RegDistId(name string, id int) {
	p.SampleId2Name[id] = name
}

func (p *PerfIdRegistry) RegCounterId(name string, id int) {
	p.Counter2Name[id] = name
}

var globalContext *PerfContext
var globalIdRegistry *PerfIdRegistry

func initGlobalPerfContext() {
	globalContext = NewPerfContext()
}

func InitGlobalPerfContext() {
	initGlobalPerfContext()
}
func initIdRegistry() {
	globalIdRegistry = &PerfIdRegistry{
		FuncId2Name:   make(map[int]string),
		SampleId2Name: make(map[int]string),
		Counter2Name:  make(map[int]string),
	}
}

func GlobalPerfContext() *PerfContext {
	return globalContext
}

func IdRegistry() *PerfIdRegistry {
	return globalIdRegistry
}

func init() {
	initIdRegistry()
	initGlobalPerfContext()
}
