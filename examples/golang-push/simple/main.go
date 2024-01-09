package main

import (
	"context"
	"log"
	"os"

	"github.com/grafana/pyroscope-go"

	"pushsimple/cputimer"
)

//go:noinline
func work(n int) {
	// revive:disable:empty-block this is fine because this is a example app, not real production code
	for i := 0; i < n; i++ {
	}
	// revive:enable:empty-block
}

func fastFunction(c context.Context) {
	pyroscope.TagWrapper(c, pyroscope.Labels("function", "fast"), func(c context.Context) {
		work(20000000)
	})
}

var functionCPUTime = cputimer.NewCPUTimerVec(cputimer.Opts{
	Name: "slow_function_cpu_time_total",
}, []string{"function"})

var slowFunctionCPUTime = functionCPUTime.WithLabelValues("slow")

func slowFunction(c context.Context) {
	slowFunctionCPUTime.Do(c, func(c context.Context) {
		work(80000000)
	})
}

func main() {
	serverAddress := os.Getenv("PYROSCOPE_SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = "http://localhost:4040"
	}
	_, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: "simple.golang.app",
		ServerAddress:   serverAddress,
		Logger:          pyroscope.StandardLogger,
	})
	if err != nil {
		log.Fatalf("error starting pyroscope profiler: %v", err)
	}
	pyroscope.TagWrapper(context.Background(), pyroscope.Labels("foo", "bar"), func(c context.Context) {
		for {
			fastFunction(c)
			slowFunction(c)
		}
	})
}
