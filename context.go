package goperfstat

import (
	"sync"
	"time"
)

type PerfContext struct {
	functions     map[int]*FuncPerf
	distributions map[int]*SampleDistribution
	counters      map[int]*Counter
	counterMu     sync.Mutex
	funcMu        sync.Mutex
	distMu        sync.Mutex
	startTime     time.Time
}
type PerfIdRegistry struct {
	FuncId2Name   map[int]string
	SampleId2Name map[int]string
	Counter2Name  map[int]string
	nextId        int
}

func NewPerfContext() *PerfContext {
	return &PerfContext{
		functions:     make(map[int]*FuncPerf),
		distributions: make(map[int]*SampleDistribution),
		counters:      make(map[int]*Counter),
		startTime:     time.Now(),
	}
}

func (p *PerfContext) StartTime() {
	p.startTime = time.Now()
}

func (p *PerfIdRegistry) RegFuncId(name string, id int) int {
	p.FuncId2Name[id] = name
	return id
}

func (p *PerfIdRegistry) RegDistId(name string, id int) int {
	p.SampleId2Name[id] = name
	return id
}

func (p *PerfIdRegistry) RegCounterId(name string, id int) int {
	p.Counter2Name[id] = name
	return id
}

func (p *PerfIdRegistry) NextId() int {
	res := p.nextId
	p.nextId++
	return res
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
