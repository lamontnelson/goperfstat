package goperfstat

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
