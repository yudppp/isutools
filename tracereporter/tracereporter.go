package tracereporter

import (
	"fmt"
	"time"

	"github.com/najeira/measure"
	"github.com/rcrowley/go-metrics"
	"github.com/yudppp/isutools/utils/measurereporter"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
)

const allServiceName = "all"

// Reporter .
type Reporter interface {
	GetConfig() *Config
	Report([]mocktracer.Span)
}

// Start .
func Start(duration time.Duration, reporters ...Reporter) {
	mt := mocktracer.Start()
	go func() {
		time.Sleep(duration)
		spans := mt.FinishedSpans()
		reportAll(reporters, spans)
		mt.Stop()
	}()
}

func reportAll(reporters []Reporter, spans []mocktracer.Span) {
	for _, reporter := range reporters {
		cfg := reporter.GetConfig()
		filterdSpans := make([]mocktracer.Span, 0, len(spans))
		for _, span := range spans {
			if cfg.serviceName != "" {
				serviceName := getServiceName(span)
				if cfg.serviceName != serviceName {
					continue
				}
			}
			filterdSpans = append(filterdSpans, span)
		}
		reporter.Report(spans)
	}
}

// SimpleReport .
type SimpleReport struct {
	*Config
}

// NewSimpleReport .
func NewSimpleReport(opts ...Option) Reporter {
	cfg := new(Config)
	defaults(cfg)
	for _, fn := range opts {
		fn(cfg)
	}
	return &SimpleReport{cfg}
}

// GetConfig is implement Reporter
func (r *SimpleReport) GetConfig() *Config {
	return r.Config
}

// Report is implement Reporter
func (r *SimpleReport) Report(spans []mocktracer.Span) {
	metricsMap := make(map[string]metrics.Timer)
	for _, span := range spans {
		resourceName := getResourceName(span)
		if resourceName == "" {
			continue
		}
		t, ok := metricsMap[resourceName]
		if !ok {
			t = metrics.NewTimer()
			metricsMap[resourceName] = t
		}
		t.Update(span.FinishTime().Sub(span.StartTime()))
	}
	result := make(measure.StatsSlice, 0, len(metricsMap))
	for resourceName, t := range metricsMap {
		stats := measure.Stats{
			Key:   resourceName,
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
	result.SortDesc(r.sortKey)
	measurereporter.Send(fmt.Sprintf("%s.csv", r.serviceName), result)
}

// getServiceName .
func getServiceName(span mocktracer.Span) string {
	tags := span.Tags()
	value, ok := tags["service.name"]
	if !ok {
		return ""
	}
	return fmt.Sprint(value)
}

// getResourceName .
func getResourceName(span mocktracer.Span) string {
	tags := span.Tags()
	value, ok := tags["resource.name"]
	if !ok {
		return ""
	}
	return fmt.Sprint(value)
}
