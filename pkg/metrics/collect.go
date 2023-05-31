package metrics

import (
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"golang.org/x/exp/slices"
)

// Collect metrics.
func Collect[C prometheus.Collector](c C) map[string]float64 {
	var (
		metrics []prometheus.Metric
		result  map[string]float64
		mChan   = make(chan prometheus.Metric)
		done    = make(chan struct{})
	)

	go func() {
		for m := range mChan {
			metrics = append(metrics, m)
		}
		close(done)
	}()

	c.Collect(mChan)
	close(mChan)
	<-done

	for _, m := range metrics {
		pb := &dto.Metric{}
		if err := m.Write(pb); err != nil {
			panic(fmt.Errorf("error happened while collecting metrics: %w", err))
		}
		if result == nil {
			result = make(map[string]float64)
		}
		result[buildKey(pb.GetLabel())] += extractValue(pb)
	}
	return result
}

func buildKey(labels []*dto.LabelPair) string {
	var result strings.Builder
	slices.SortFunc(labels, func(a, b *dto.LabelPair) bool {
		if *a.Name == *b.Name {
			return *a.Value < *b.Name
		}
		return *a.Name < *b.Name
	})
	for idx, entry := range labels {
		if idx > 0 {
			result.WriteString(",")
		}
		result.WriteString(*entry.Name)
		result.WriteString("=")
		result.WriteString(*entry.Value)
	}
	return result.String()
}

func extractValue(pb *dto.Metric) float64 {
	if pb.Gauge != nil {
		return pb.Gauge.GetValue()
	}
	if pb.Counter != nil {
		return pb.Counter.GetValue()
	}
	if pb.Untyped != nil {
		return pb.Untyped.GetValue()
	}
	if pb.Histogram != nil {
		return float64(*pb.Histogram.SampleCount)
	}
	panic(fmt.Errorf("collected a non-gauge/counter/untyped metric: %s", pb))
}
