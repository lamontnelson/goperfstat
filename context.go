package goperfstat

type PerfContext struct {
	functions map[int]*FuncPerf
}

func NewPerfContext() *PerfContext {
	return &PerfContext{functions: make(map[int]*FuncPerf)}
}

var globalContext *PerfContext

func InitGlobalPerfContext() {
	globalContext = NewPerfContext()
}

func init() {
	InitGlobalPerfContext()
}
