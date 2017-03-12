package goperfstat

import "testing"

func TestDistribution(t *testing.T) {
	InitGlobalPerfContext()
	globalContext.RegDistId("Sample", 1)

	t.Run("CanSampleData", func(t *testing.T) {
		iters := 10000
		rate := 0.25
		c := int(float64(iters) * rate)
		id := 1
		for x := 0; x < iters; x++ {
			TakeSample(nil, id, rate, c, float64(x))
		}
	})

	globalContext.Report()
	InitGlobalPerfContext()
}
