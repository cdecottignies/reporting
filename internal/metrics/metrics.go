package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	nameSpace = "ns"
	subSystem = "example"
)

var HelloCounter = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: nameSpace,
	Subsystem: subSystem,
	Name:      "hello_total",
	Help:      "The number of hello",
}, []string{"scope"})
