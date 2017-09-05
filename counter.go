package goperfstat

import (
	"errors"
	"fmt"
)

type Counter struct {
	id    int
	Count float64
}

func NewCounter(id int) *Counter {
	return &Counter{id: id, Count: 0}
}

func (c *Counter) Inc(v float64) {
	c.Count += v
}

func getCounterUnlocked(context *PerfContext, id int) (*Counter, error) {
	var counter *Counter
	var found bool

	if context == nil {
		context = globalContext
	}

	if counter, found = context.counters[id]; !found {
		return nil, errors.New(fmt.Sprintf("counter %v not found", id))
	}

	return counter, nil
}

func GetCounter(context *PerfContext, id int) (*Counter, error) {
	if context == nil {
		context = globalContext
	}

	context.counterMu.Lock()
	defer context.counterMu.Unlock()
	return getCounterUnlocked(context, id)
}

func Count(context *PerfContext, id int, v float64) {
	var counter *Counter

	if context == nil {
		context = globalContext
	}

	context.counterMu.Lock()
	defer context.counterMu.Unlock()

	if counter, _ = getCounterUnlocked(context, id); counter == nil {
		counter = NewCounter(id)
		context.counters[id] = counter
	}
	counter.Inc(v)
}
