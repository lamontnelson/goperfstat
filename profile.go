package goperfstat

import (
	"errors"
	"fmt"
	"github.com/montanaflynn/stats"
	"time"
)

var perfStats *PerfInfo

type PerfInfo struct {
	functions map[int]*FuncPerf
}

func NewPerfInfo() *PerfInfo {
	return &PerfInfo{functions: make(map[int]*FuncPerf)}
}

type FuncPerf struct {
	id int
	// number of times function was called
	count uint64
	// time measured in nanoseconds
	times stats.Float64Data
	start time.Time
}

func GetFuncPerf(id int) (*FuncPerf, error) {
	var perf *FuncPerf
	var ok bool
	if perf, ok = perfStats.functions[id]; !ok {
		return nil, errors.New(fmt.Sprintf("no perf info for %v", id))
	}
	return perf, nil
}

func CountCalls(id int) *FuncPerf {
	perf, err := GetFuncPerf(id)
	if err != nil {
		times := make(stats.Float64Data, 10)
		times = times[:0]
		perf = &FuncPerf{id: id, times: times}
		perfStats.functions[id] = perf
	}

	perf.count++
	return perf
}

func TimeCountCalls(id int) *FuncPerf {
	perf := CountCalls(id)
	perf.start = time.Now()
	return perf
}

func End(id int) {
	perf, err := GetFuncPerf(id)
	if err != nil {
		panic("perf record not found")
	}

	d := time.Since(perf.start)
	perf.times = append(perf.times, float64(d/time.Nanosecond))
}

func init() {
	perfStats = NewPerfInfo()
}
