package tracereporter

import (
	"fmt"
	"time"

	"github.com/rcrowley/go-metrics"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	"github.com/najeira/measure"
	"github.com/yudppp/isutools/utils/measurereporter"
)

// Reporter .
type Reporter interface {
	GetServiceName() string
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
	spansMap := make(map[string][]mocktracer.Span, len(reporters))
	for _, v := range reporters {
		spansMap[v.GetServiceName()] = make([]mocktracer.Span, 0)
	} 
	for _, span := range spans {
		serviceName := getServiceName(span)
		serviceSpans, ok := spansMap[serviceName]
		if !ok {
			fmt.Printf("Unknow service.name: %s\n", serviceName)
			continue
		}
		spansMap[serviceName] = append(serviceSpans, span)
	}

	for _, reporter := range reporters {
		serviceName := reporter.GetServiceName()
		spans, ok := spansMap[serviceName]
		if !ok {
			continue
		}
		reporter.Report(spans)
	}
}

// SimpleReport .
type SimpleReport struct {
	serviceName string
	sortKey string
}

// NewSimpleReport .
// sortKey is key | count | sum | min | max | avg | rate | p95
func NewSimpleReport(serviceName, sortKey string) Reporter {
	return &SimpleReport{serviceName: serviceName, sortKey: sortKey}
}

// GetServiceName is implement Reporter
func (r *SimpleReport) GetServiceName() string {
	return r.serviceName
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
	sortKey := r.sortKey
	if sortKey == "" {
		sortKey = measure.Sum
	}
	result.SortDesc(sortKey)
	measurereporter.Send(fmt.Sprintf("%s.csv", r.GetServiceName()), result)
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