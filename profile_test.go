package goperfstat

import "testing"
import "time"

func TestCanCountCalls(t *testing.T) {
	CountCalls(0)
	CountCalls(0)
	if perfStats.functions[0].count != 2 {
		t.Fatalf("expected 2 calls; got %v", perfStats.functions[0].count)
	}
}

func TestCanTimeCalls(t *testing.T) {
	sleepDuration := 10 * time.Millisecond
	c := TimeCountCalls(0)
	time.Sleep(sleepDuration)
	End(0)
	if len(c.times) != 1 {
		t.Fatalf("expected %v measurements; got %v", 1, len(c.times))
	}
	d := time.Duration(c.times[0]) * time.Nanosecond
	if d < sleepDuration {
		t.Fatalf("recorded duration is %v; expected at least %v", d, sleepDuration)
	}
}
