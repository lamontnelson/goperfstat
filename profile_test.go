package goperfstat

import "testing"
import "time"

func TestProfile(t *testing.T) {
	InitGlobalPerfContext()

	t.Run("CanCountCalls", func(t *testing.T) {
		id := 0
		CountCalls(nil, id)
		CountCalls(nil, id)
		if globalContext.functions[id].count != 2 {
			t.Fatalf("expected 2 calls; got %v", globalContext.functions[id].count)
		}
	})

	t.Run("CanTimeCalls", func(t *testing.T) {
		id := 1
		sleepDuration := 10 * time.Millisecond
		c := TimeCountCalls(nil, id)
		time.Sleep(sleepDuration)
		End(nil, id)
		if len(c.times) != 1 {
			t.Fatalf("expected %v measurements; got %v", 1, len(c.times))
		}
		d := time.Duration(c.times[0]) * time.Nanosecond
		if d < sleepDuration {
			t.Fatalf("recorded duration is %v; expected at least %v", d, sleepDuration)
		}
	})

	t.Run("CanUseLocalContext", func(t *testing.T) {
		id := 0
		localContext := NewPerfContext()
		CountCalls(localContext, id)
		CountCalls(localContext, id)
		if localContext.functions[id].count != 2 {
			t.Fatalf("expected 2 calls; got %v", localContext.functions[id].count)
		}
	})

	InitGlobalPerfContext()
}
