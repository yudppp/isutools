package tracereporter

import (
	"fmt"
	"time"
	"strings"

	"github.com/rcrowley/go-metrics"
	"github.com/yudppp/isutools/utils/slackcat"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
)

// HTTPReuestOperationName .
const HTTPReuestOperationName = "http.request"

type httpRequestRepoter struct {
	*Config
}

// NewHTTPRequestReport .
func NewHTTPRequestReport(opts ...Option) Reporter {
	cfg := new(Config)
	defaults(cfg)
	for _, fn := range opts {
		fn(cfg)
	}
	return &httpRequestRepoter{cfg}
}

// GetConfig is implement Reporter
func (r *httpRequestRepoter) GetConfig() *Config {
	return r.Config
}

// Report is implement Reporter
func (r *httpRequestRepoter) Report(spans []mocktracer.Span) {
	httpRequestSpans := make([]mocktracer.Span, 0,len(spans))
	httpRequestMap := make(map[uint64]string, len(spans))
	for _, span := range spans {
		if span.OperationName() != HTTPReuestOperationName {
			continue
		}
		resourceName := getResourceName(span)
		httpRequestMap[span.SpanID()] = resourceName
		httpRequestSpans = append(httpRequestSpans, span)
	}
	childSpansMap := make(map[string][]mocktracer.Span, len(spans))
	for _, span := range spans {
		parentID := span.ParentID()
		if parentID == 0 {
			continue
		}
		parentResource, ok := httpRequestMap[parentID]
		if !ok {
			continue
		}
		childSpans, ok := childSpansMap[parentResource]
		if ok {
			childSpansMap[parentResource] = append(childSpans, span)
		} else {
			childSpansMap[parentResource] = []mocktracer.Span{span}
		}
	}


	if len(httpRequestSpans) == 0 {
		slackcat.SendText("http.request.txt", "empty")
		return
	}

	parentMetrics := loadMesure(httpRequestSpans)
	var b strings.Builder
	remainCount := len(parentMetrics)
	for resourceName, parentMetric := range parentMetrics {
		remainCount--
		parentCount :=  parentMetric.Count()
		fmt.Fprintf(&b, "%s (Count=%v,Avg=%.2f[ms],Sum=%.2f[ms])\n", resourceName, parentCount, parentMetric.Mean() / float64(time.Millisecond), float64(parentMetric.Sum()) / float64(time.Millisecond))
		childSpans, ok := childSpansMap[resourceName]
		if !ok {
			fmt.Fprintf(&b, "")
			continue
		}
		childMetrics := loadMesure(childSpans)
		for resourceName, metrics := range childMetrics {
			count := float64(metrics.Count()) / float64(parentCount)
			avg := metrics.Mean() / float64(time.Millisecond) * count
			fmt.Fprintf(&b, "- %s (Count=%.2f[/req],Avg=%.2f[ms/req])\n", resourceName, count, avg)
		}
		if remainCount != 0 {
			fmt.Fprintf(&b, "")
		}
	}
	slackcat.SendText("http.request.txt", b.String())
}

func loadMesure(spans []mocktracer.Span) map[string]metrics.Timer {
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
	return metricsMap
}
