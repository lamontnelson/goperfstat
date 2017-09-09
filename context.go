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
	mu            sync.Mutex
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

func (p *PerfIdRegistry) RegFunc(name string) int {
	p.mu.Lock()
	defer p.mu.Unlock()
	id := p.newId()
	p.FuncId2Name[id] = name
	return id
}

func (p *PerfIdRegistry) RegDist(name string) int {
	p.mu.Lock()
	defer p.mu.Unlock()
	id := p.newId()
	p.SampleId2Name[id] = name
	return id
}

func (p *PerfIdRegistry) RegCounter(name string) int {
	p.mu.Lock()
	defer p.mu.Unlock()
	id := p.newId()
	p.Counter2Name[id] = name
	return id
}

func (p *PerfIdRegistry) newId() int {
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
