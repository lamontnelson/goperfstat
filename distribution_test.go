package goperfstat

import "testing"
import "fmt"

func TestDistribution(t *testing.T) {
	initGlobalPerfContext()
	globalIdRegistry.RegDistId("Sample", 1)

	t.Run("CanSampleData", func(t *testing.T) {
		iters := 10000
		rate := 0.25
		c := int(float64(iters) * rate)
		id := 1
		for x := 0; x < iters; x++ {
			TakeSample(nil, id, rate, c, float64(x))
		}
	})

	fmt.Println(globalContext.Report())
	initGlobalPerfContext()
}
