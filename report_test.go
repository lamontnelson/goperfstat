package goperfstat

import "testing"
import "time"

func TestReportData(t *testing.T) {
	initGlobalPerfContext()
	globalIdRegistry.RegDistId("foo", 1)
	globalIdRegistry.RegFuncId("bar", 2)
	globalIdRegistry.RegCounterId("baz", 3)

	t.Run("CanGenerateReport", func(t *testing.T) {
		TimeFuncCall(nil, 2, time.Now().Add(-1*time.Second))
		TimeFuncCall(nil, 2, time.Now().Add(-1*time.Second))
		TimeFuncCall(nil, 2, time.Now().Add(-1*time.Second))
		TakeSample(nil, 1, 1.0, 1, float64(123))
		TakeSample(nil, 1, 1.0, 1, float64(123))
		TakeSample(nil, 1, 1.0, 1, float64(123))
		Count(nil, 3, 10.0)

		r := GlobalPerfContext().ReportData()
		if c, found := r.counters["baz"]; !found {
			t.Fatal("counter not found")
		} else {
			//t.Log(c)
			if c != 10.0 {
				t.Fatal("counter wrong value")
			}
		}

		if f, found := r.functions["bar"]; !found {
			t.Fatal("function not found")
		} else {
			//t.Logf("%+v\n", f)
			if f.count != 3 {
				t.Fatal("missing function data")
			}
		}

		if d, found := r.distributions["foo"]; !found {
			t.Fatal("distribution not found")
		} else {
			//t.Logf("%+v\n", d)
			if d.count != 3 {
				t.Fatal("missing distribution data")
			}
			if d.average != 123 {
				t.Fatal("distribution average incorrect")
			}
		}
	})
}
