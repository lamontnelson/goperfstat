package goperfstat

import "testing"
import "time"

func TestReportData(t *testing.T) {

	t.Run("CanGenerateJSON", func(t *testing.T) {
		initGlobalPerfContext()
		globalIdRegistry.RegDistId("foo", 1)
		globalIdRegistry.RegFuncId("bar", 2)
		globalIdRegistry.RegCounterId("baz", 3)
		TimeFuncCall(nil, 2, time.Now().Add(-1*time.Second))
		TimeFuncCall(nil, 2, time.Now().Add(-1*time.Second))
		TimeFuncCall(nil, 2, time.Now().Add(-1*time.Second))
		TakeSample(nil, 1, 1.0, 1, float64(123))
		TakeSample(nil, 1, 1.0, 1, float64(123))
		TakeSample(nil, 1, 1.0, 1, float64(123))
		Count(nil, 3, 10.0)
		t.Logf("%v\n", string(GlobalPerfContext().ReportJson()))
	})

	t.Run("CanGenerateReport", func(t *testing.T) {
		initGlobalPerfContext()
		globalIdRegistry.RegDistId("foo", 1)
		globalIdRegistry.RegFuncId("bar", 2)
		globalIdRegistry.RegCounterId("baz", 3)
		TimeFuncCall(nil, 2, time.Now().Add(-1*time.Second))
		TimeFuncCall(nil, 2, time.Now().Add(-1*time.Second))
		TimeFuncCall(nil, 2, time.Now().Add(-1*time.Second))
		TakeSample(nil, 1, 1.0, 1, float64(123))
		TakeSample(nil, 1, 1.0, 1, float64(123))
		TakeSample(nil, 1, 1.0, 1, float64(123))
		Count(nil, 3, 10.0)

		r := GlobalPerfContext().ReportData()
		if c, found := r.Counters["baz"]; !found {
			t.Fatal("counter not found")
		} else {
			t.Log(c)
			if c != 10.0 {
				t.Fatal("counter wrong value")
			}
		}

		if f, found := r.Functions["bar"]; !found {
			t.Fatal("function not found")
		} else {
			t.Logf("%+v\n", f)
			if f.Count != 3 {
				t.Fatal("missing function data")
			}
		}

		if d, found := r.Distributions["foo"]; !found {
			t.Fatal("distribution not found")
		} else {
			t.Logf("%+v\n", d)
			if d.Count != 3 {
				t.Fatal("missing distribution data")
			}
			if d.Average != 123 {
				t.Fatal("distribution average incorrect")
			}
		}
	})
}
