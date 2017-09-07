package goperfstat

import "testing"
import "fmt"

func TestCounter(t *testing.T) {
	id := 0
	initGlobalPerfContext()
	globalIdRegistry.RegCounterId("x", id)

	t.Run("CounterStartAt0", func(t *testing.T) {
		Count(nil, id, 0)
		c, _ := GetCounter(nil, id)
		if c.Count != 0 {
			t.Fatalf("counter initialized to %v, not 0", c)
		}
	})

	t.Run("CanCount", func(t *testing.T) {
		Count(nil, 0, 10)
		c, _ := GetCounter(nil, id)
		if c.Count != 10 {
			t.Fatalf("invalid count")
		}
	})

	fmt.Println(globalContext.Report())
}
