package tracereporter

import (
	"fmt"
	"time"

	"github.com/rcrowley/go-metrics"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	"github.com/najeira/measure"
)

// Start .
func Start(duration time.Duration, sortKey string, getKey func(span mocktracer.Span) string) {
	mt := mocktracer.Start()
	go func() {
		time.Sleep(duration)
		spans := mt.FinishedSpans()
		metricsMap := make(map[string]metrics.Timer)
		for _, span := range spans {
			key := getKey(span)
			if key == "" {
				continue
			}
			t, ok := metricsMap[key]
			if !ok {
				t = metrics.NewTimer()
				metricsMap[key] = t
			}
			t.Update(span.FinishTime().Sub(span.StartTime()))
		}
		result := make(measure.StatsSlice, 0, len(metricsMap))
		for key, t := range metricsMap {
			stats := measure.Stats{
				Key:   key,
				Count: t.Count(),
				Sum:   float64(t.Sum()) / float64(time.Millisecond),
				Min:   float64(t.Min()) / float64(time.Millisecond),
				Max:   float64(t.Max()) / float64(time.Millisecond),
				Avg:   t.Mean() / float64(time.Millisecond),
				Rate:  t.Rate1(),
				P95:   t.Percentile(0.95) / float64(time.Millisecond),
			}
			result = append(result, stats)
		}
		if sortKey == "" {
			sortKey = "Sum"
		}
		result.SortDesc(sortKey)
		fmt.Println("Key, Count, Sum, Min, Sum, Max, Avg, Rate, P95")
		for _, row := range result {
			fmt.Println(row.Key, row.Count, row.Sum, row.Min, row.Sum, row.Max, row.Avg, row.Rate, row.P95)
		}
		mt.Stop()
	}()
}
