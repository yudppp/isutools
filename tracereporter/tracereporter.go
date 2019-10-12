package tracereporter

import (
	"fmt"
	"time"

	"github.com/rcrowley/go-metrics"
	"github.com/yudppp/isutools/utils/measurereporter"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"

	"github.com/najeira/measure"
)

// Start .
func Start(duration time.Duration, serviceName, sortKey string, getKey func(span mocktracer.Span) string) {
	mt := mocktracer.Start()
	go func() {
		time.Sleep(duration)
		spans := mt.FinishedSpans()
		metricsMap := make(map[string]metrics.Timer)
		for _, span := range spans {
			tags := span.Tags()
			spanServiceName, ok := tags["service.name"]
			if ok && serviceName != "" && spanServiceName != spanServiceName {
				continue
			}
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
		measurereporter.Send(fmt.Sprintf("%s.csv", serviceName), result)
		mt.Stop()
	}()
}
