package goperfstat

import (
	"github.com/montanaflynn/stats"
	"sync"
	"time"
)

const (
	ARR_SIZE = 10000
)

type FuncPerf struct {
	id int
	// number of times function was called
	count uint64
	// time measured in nanoseconds
	times stats.Float64Data

	mu sync.Mutex
}

func TimeFuncCall(context *PerfContext, id int, start time.Time) {
	if context == nil {
		context = globalContext
	}

	context.funcMu.Lock()
	var perf *FuncPerf
	if fp, found := context.functions[id]; found {
		perf = fp
	} else {
		times := make(stats.Float64Data, ARR_SIZE)
		times = times[:0]
		perf = &FuncPerf{id: id, times: times}
		context.functions[id] = perf
	}
	context.funcMu.Unlock()

	perf.mu.Lock()
	defer perf.mu.Unlock()
	d := time.Since(start)
	perf.times = append(perf.times, float64(d/time.Nanosecond))
	perf.count++
}
