package goperfstat

import (
	"fmt"
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

func (p *PerfContext) Report() {
	for id, perf := range p.functions {

		var timesSummary string
		if len(perf.times) > 0 {
			min, _ := perf.times.Min()
			max, _ := perf.times.Max()
			fifty, _ := perf.times.Percentile(50)
			ninety, _ := perf.times.Percentile(90)
			ninenine, _ := perf.times.Percentile(99)
			timesSummary = fmt.Sprintf("; 'min,50,90,99,max': [%v, %v, %v, %v, %v]", min, fifty, ninety, ninenine, max)
		}

		name, found := p.FuncId2Name[id]
		if found {
			fmt.Printf("%v: { count: %v%v }\n", name, perf.count, timesSummary)
		} else {
			fmt.Printf("id %v: { count: %v%v }\n", id, perf.count, timesSummary)
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
