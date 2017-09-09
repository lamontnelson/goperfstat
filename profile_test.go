package goperfstat

import "testing"
import "time"
import "fmt"

func TestProfile(t *testing.T) {
	initGlobalPerfContext()
	idc := globalIdRegistry.RegFunc("CountCalls")
	idf := globalIdRegistry.RegFunc("TimeCalls")

	t.Run("CanCountCalls", func(t *testing.T) {
		TimeFuncCall(nil, idc, time.Now())
		TimeFuncCall(nil, idc, time.Now())
		if globalContext.functions[idc].count != 2 {
			t.Fatalf("expected 2 calls; got %v", globalContext.functions[idc].count)
		}
	})

	t.Run("CanTimeCalls", func(t *testing.T) {
		var st time.Time
		id := idf
		sleepDuration := 10 * time.Millisecond

		st = time.Now()
		time.Sleep(sleepDuration)
		TimeFuncCall(nil, id, st)

		st = time.Now()
		time.Sleep(sleepDuration)
		TimeFuncCall(nil, id, st)

		st = time.Now()
		time.Sleep(sleepDuration)
		TimeFuncCall(nil, id, st)

		expectedLen := 3
		times := globalContext.functions[id].times
		l := len(times)
		if l != expectedLen {
			t.Fatalf("expected %v measurements; got %v", expectedLen, l)
		}
		d := time.Duration(times[0]) * time.Nanosecond
		if d < sleepDuration {
			t.Fatalf("recorded duration is %v; expected at least %v", d, sleepDuration)
		}
	})

	t.Run("CanUseLocalContext", func(t *testing.T) {
		id := idf
		localContext := NewPerfContext()
		st := time.Now()
		TimeFuncCall(localContext, id, st)
		TimeFuncCall(localContext, id, st)
		if localContext.functions[id].count != 2 {
			t.Fatalf("expected 2 calls; got %v", localContext.functions[id].count)
		}
	})

	fmt.Println(globalContext.Report())
	initGlobalPerfContext()
}
