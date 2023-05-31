package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

// Collect metrics.
func Collect[C prometheus.Collector](c C) []*dto.Metric {
	tgt := make(chan prometheus.Metric)
	var result []*dto.Metric
	go func() {
		for metric := range tgt {
			elem := &dto.Metric{}
			result = append(result, elem)
			_ = metric.Write(elem)
		}
	}()
	defer close(tgt)
	c.Collect(tgt)
	return result
}
