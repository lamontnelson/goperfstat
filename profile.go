package goperfstat

import (
	"errors"
	"fmt"
	"github.com/montanaflynn/stats"
	"time"
)

type FuncPerf struct {
	id int
	// number of times function was called
	count uint64
	// time measured in nanoseconds
	times stats.Float64Data
	start time.Time
}

func FunctionStats(context *PerfContext, id int) (*FuncPerf, error) {
	var perf *FuncPerf
	var ok bool

	if context == nil {
		context = globalContext
	}

	if perf, ok = context.functions[id]; !ok {
		return nil, errors.New(fmt.Sprintf("no perf info for %v", id))
	}
	return perf, nil
}

func CountCalls(context *PerfContext, id int) *FuncPerf {
	if context == nil {
		context = globalContext
	}

	perf, err := FunctionStats(context, id)
	if err != nil {
		times := make(stats.Float64Data, 10)
		times = times[:0]
		perf = &FuncPerf{id: id, times: times}
		context.functions[id] = perf
	}

	perf.count++
	return perf
}

func TimeCountCalls(context *PerfContext, id int) *FuncPerf {
	if context == nil {
		context = globalContext
	}

	perf := CountCalls(context, id)
	perf.start = time.Now()
	return perf
}

func End(context *PerfContext, id int) {
	if context == nil {
		context = globalContext
	}

	perf, err := FunctionStats(context, id)
	if err != nil {
		panic("perf record not found")
	}

	d := time.Since(perf.start)
	perf.times = append(perf.times, float64(d/time.Nanosecond))
}
